package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"

	blockchainRequest "go-blockchain/blockchain/infra/http/request"
	"go-blockchain/blockchain/infra/http/response"
	"go-blockchain/utils"
	"go-blockchain/wallet/domain/entity"
	"go-blockchain/wallet/domain/repository"
	walletRequest "go-blockchain/wallet/infra/http/request"
)

const tempDir = "templates"

type walletServerRepository struct{}

func NewWalletServerRepository() repository.WalletServerRepository {
	return &walletServerRepository{}
}

func (wsr walletServerRepository) Port(ws *entity.WalletServer) uint16 {
	return ws.Port
}

func (wsr walletServerRepository) Gateway(ws *entity.WalletServer) string {
	return ws.Gateway
}

func (wsr walletServerRepository) Index(ws *entity.WalletServer, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(tempDir, "index.html"))
		t.Execute(w, "")
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (wsr walletServerRepository) Wallet(wr repository.WalletRepository, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := NewWallet()
		m, _ := wr.MarshalJSON(myWallet)
		io.WriteString(w, string(m[:]))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (wsr walletServerRepository) CreateTransaction(ws *entity.WalletServer, tr repository.TransactionRepository, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var t walletRequest.TransactionRequest
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
		privateKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			log.Println("ERROR: parse error")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		value32 := float32(value)

		w.Header().Add("Content-Type", "application/json")

		transaction := NewTransaction(privateKey, publicKey,
			*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
		signature := tr.GenerateSignature(transaction)
		signatureStr := signature.String()

		bt := &blockchainRequest.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32, Signature: &signatureStr,
		}
		m, _ := json.Marshal(bt)
		buf := bytes.NewBuffer(m)

		resp, _ := http.Post(wsr.Gateway(ws)+"/transactions", "application/json", buf)
		if resp.StatusCode == 201 {
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		}
		io.WriteString(w, string(utils.JsonStatus("fail")))
	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (wsr walletServerRepository) WalletAmount(ws *entity.WalletServer, w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("blockchain_address")
		endpoint := fmt.Sprintf("%s/amount", wsr.Gateway(ws))

		client := &http.Client{}
		bcsReq, _ := http.NewRequest("GET", endpoint, nil)
		q := bcsReq.URL.Query()
		q.Add("blockchain_address", blockchainAddress)
		bcsReq.URL.RawQuery = q.Encode()

		bcsResp, err := client.Do(bcsReq)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		if bcsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResp.Body)
			var bar response.AmountResponse
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				io.WriteString(w, string(utils.JsonStatus("fail")))
				return
			}

			m, _ := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})
			io.WriteString(w, string(m[:]))
		} else {
			io.WriteString(w, string(utils.JsonStatus("fail")))
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (wsr walletServerRepository) Run(ws *entity.WalletServer, wr repository.WalletRepository, tr repository.TransactionRepository) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		wsr.Index(ws, w, req)
	})
	http.HandleFunc("/wallet", func(w http.ResponseWriter, req *http.Request) {
		wsr.Wallet(wr, w, req)
	})
	http.HandleFunc("/wallet/amount", func(w http.ResponseWriter, req *http.Request) {
		wsr.WalletAmount(ws, w, req)
	})
	http.HandleFunc("/transaction", func(w http.ResponseWriter, req *http.Request) {
		wsr.CreateTransaction(ws, tr, w, req)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(wsr.Port(ws))), nil))
}
