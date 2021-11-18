package repository

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"go-blockchain/blockchain/block"
	"go-blockchain/blockchain/domain/entity"
	"go-blockchain/blockchain/domain/repository"
	"go-blockchain/utils"
	"go-blockchain/wallet"
)

type blockchainServerRepository struct{}

func NewBlockchainServerRepository() repository.BlockchainServerRepository {
	return &blockchainServerRepository{}
}

var cache map[string]*entity.Blockchain = make(map[string]*entity.Blockchain)

func NewBlockchainServer(port uint16) *entity.BlockchainServer {
	return &entity.BlockchainServer{Port: port}
}

func (bsr *blockchainServerRepository) Port(bs *entity.BlockchainServer) uint16 {
	return bs.Port
}

func (bsr *blockchainServerRepository) GetBlockchain(bs *entity.BlockchainServer, bcr repository.BlockchainRepository, br repository.BlockRepository) *entity.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = NewBlockchain(br, bcr, minersWallet.BlockchainAddress(), bsr.Port(bs))
		cache["blockchain"] = bc
		log.Printf("private_key %v", minersWallet.PrivateKeyStr())
		log.Printf("publick_key %v", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address %v", minersWallet.BlockchainAddress())
	}
	return bc
}

func (bsr *blockchainServerRepository) GetChain(bs *entity.BlockchainServer, bcr repository.BlockchainRepository, br repository.BlockRepository, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bsr.GetBlockchain(bs, bcr, br)
		m, _ := bcr.MarshalJSON(bc)
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("ERROR: Invalid HTTP Method")

	}
}

func (bsr *blockchainServerRepository) Transactions(bs *entity.BlockchainServer, bcr repository.BlockchainRepository, br repository.BlockRepository, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bsr.GetBlockchain(bs, bcr, br)
		transactions := bcr.TransactionPool(bc)
		m, _ := json.Marshal(struct {
			Transactions []*entity.Transaction `json:"transactions"`
			Length       int                   `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m[:]))

	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t block.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)
		bc := bsr.GetBlockchain(bs, bcr, br)
		isCreated := bcr.CreateTransaction(bc, *t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = utils.JsonStatus("success")
		}
		io.WriteString(w, string(m))
	case http.MethodPut:
		decoder := json.NewDecoder(req.Body)
		var t block.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(s)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		publicKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		signature := utils.SignatureFromString(*t.Signature)
		bc := bsr.GetBlockchain(bs, bcr, br)
		isUpdated := bcr.AddTransaction(bc, *t.SenderBlockchainAddress,
			*t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isUpdated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}
		io.WriteString(w, string(m))
	case http.MethodDelete:
		bc := bsr.GetBlockchain(bs, bcr, br)
		bcr.ClearTransactionPool(bc)
		io.WriteString(w, string(utils.JsonStatus("success")))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bsr *blockchainServerRepository) Mine(bs *entity.BlockchainServer, bcr repository.BlockchainRepository, br repository.BlockRepository, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := bsr.GetBlockchain(bs, bcr, br)
		isMined := bcr.Mining(bc, br)

		var m []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bsr *blockchainServerRepository) StartMine(bs *entity.BlockchainServer, bcr repository.BlockchainRepository, br repository.BlockRepository, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		bc := bsr.GetBlockchain(bs, bcr, br)
		bcr.StartMining(bc, br)

		m := utils.JsonStatus("success")
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bsr *blockchainServerRepository) Amount(bs *entity.BlockchainServer, bcr repository.BlockchainRepository, br repository.BlockRepository, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		bc := bsr.GetBlockchain(bs, bcr, br)
		amount := bcr.CalculateTotalAmount(bc, blockchainAddress)

		ar := &block.AmountResponse{Amount: amount}
		m, _ := ar.MarshalJSON()

		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m[:]))

	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bsr *blockchainServerRepository) Consensus(bs *entity.BlockchainServer, bcr repository.BlockchainRepository, br repository.BlockRepository, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		bc := bsr.GetBlockchain(bs, bcr, br)
		replaced := bcr.ResolveConflicts(bc, br)

		w.Header().Add("Content-Type", "application/json")
		if replaced {
			io.WriteString(w, string(utils.JsonStatus("success")))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bsr *blockchainServerRepository) Run(bs *entity.BlockchainServer, bcr repository.BlockchainRepository, br repository.BlockRepository) {
	bc := bsr.GetBlockchain(bs, bcr, br)
	bcr.Run(bc, br)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		bsr.GetChain(bs, bcr, br, w, req)
	})
	http.HandleFunc("/transactions", func(w http.ResponseWriter, req *http.Request) {
		bsr.Transactions(bs, bcr, br, w, req)
	})
	http.HandleFunc("/mine", func(w http.ResponseWriter, req *http.Request) {
		bsr.Mine(bs, bcr, br, w, req)
	})
	http.HandleFunc("/mine/start", func(w http.ResponseWriter, req *http.Request) {
		bsr.StartMine(bs, bcr, br, w, req)
	})
	http.HandleFunc("/amount", func(w http.ResponseWriter, req *http.Request) {
		bsr.Amount(bs, bcr, br, w, req)
	})
	http.HandleFunc("/consensus", func(w http.ResponseWriter, req *http.Request) {
		bsr.Consensus(bs, bcr, br, w, req)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bsr.Port(bs))), nil))
}
