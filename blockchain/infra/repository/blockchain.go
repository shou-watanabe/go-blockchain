package repository

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go-blockchain/blockchain/domain/entity"
	"go-blockchain/blockchain/domain/repository"
	"go-blockchain/blockchain/infra/http/request"
	"go-blockchain/utils"
)

const (
	MINING_DIFFICULTY = 3
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
	MINING_TIMER_SEC  = 20

	BLOCKCHAIN_PORT_RANGE_START      = 5000
	BLOCKCHAIN_PORT_RANGE_END        = 5003
	NEIGHBOR_IP_RANGE_START          = 0
	NEIGHBOR_IP_RANGE_END            = 1
	BLOCKCHIN_NEIGHBOR_SYNC_TIME_SEC = 20
)

type blockchainRepository struct{}

func NewBlockchainRepository() repository.BlockchainRepository {
	return &blockchainRepository{}
}

func NewBlockchain(br repository.BlockRepository, bcr repository.BlockchainRepository, blockchainAddress string, port uint16) *entity.Blockchain {
	b := &entity.Block{}
	bc := new(entity.Blockchain)
	bc.BlockchainAddress = blockchainAddress
	bcr.CreateBlock(bc, 0, br.Hash(b))
	bc.Port = port
	return bc
}

func (bcr *blockchainRepository) Chain(bc *entity.Blockchain) []*entity.Block {
	return bc.Chain
}

func (bcr *blockchainRepository) Run(bc *entity.Blockchain, br repository.BlockRepository) {
	bcr.StartSyncNeighbors(bc)
	bcr.ResolveConflicts(bc, br)
	bcr.StartMining(bc, br)
}

func (bcr *blockchainRepository) SetNeighbors(bc *entity.Blockchain) {
	bc.Neighbors = utils.FindNeighbors(
		utils.GetHost(), bc.Port,
		NEIGHBOR_IP_RANGE_START, NEIGHBOR_IP_RANGE_END,
		BLOCKCHAIN_PORT_RANGE_START, BLOCKCHAIN_PORT_RANGE_END)
	log.Printf("%v", bc.Neighbors)
}

func (bcr *blockchainRepository) SyncNeighbors(bc *entity.Blockchain) {
	bc.MuxNeighbors.Lock()
	defer bc.MuxNeighbors.Unlock()
	bcr.SetNeighbors(bc)
}

func (bcr *blockchainRepository) StartSyncNeighbors(bc *entity.Blockchain) {
	bcr.SyncNeighbors(bc)
	_ = time.AfterFunc(time.Second*BLOCKCHIN_NEIGHBOR_SYNC_TIME_SEC, func() { bcr.StartSyncNeighbors(bc) })
}

func (bcr *blockchainRepository) TransactionPool(bc *entity.Blockchain) []*entity.Transaction {
	return bc.TransactionPool
}

func (bcr *blockchainRepository) ClearTransactionPool(bc *entity.Blockchain) {
	bc.TransactionPool = bc.TransactionPool[:0]
}

func (bcr *blockchainRepository) MarshalJSON(bc *entity.Blockchain) ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*entity.Block `json:"chain"`
	}{
		Blocks: bc.Chain,
	})
}

func (bcr *blockchainRepository) UnmarshalJSON(bc *entity.Blockchain, data []byte) error {
	v := &struct {
		Blocks *[]*entity.Block `json:"chain"`
	}{
		Blocks: &bc.Chain,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}

func (bcr *blockchainRepository) CreateBlock(bc *entity.Blockchain, nonce int, previousHash [32]byte) *entity.Block {
	b := NewBlock(nonce, previousHash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, b)
	bc.TransactionPool = []*entity.Transaction{}
	for _, n := range bc.Neighbors {
		endpoint := fmt.Sprintf("http://%s/transactions", n)
		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}
	return b
}

func (bcr *blockchainRepository) LastBlock(bc *entity.Blockchain) *entity.Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bcr *blockchainRepository) Print(br repository.BlockRepository, tr repository.TransactionRepository, bc *entity.Blockchain) {
	for i, block := range bc.Chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i,
			strings.Repeat("=", 25))
		br.Print(block, tr)
	}
	fmt.Printf("%s\n", strings.Repeat("*", 25))
}

func (bcr *blockchainRepository) CreateTransaction(bc *entity.Blockchain, sender string, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bcr.AddTransaction(bc, sender, recipient, value, senderPublicKey, s)

	if isTransacted {
		for _, n := range bc.Neighbors {
			publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(),
				senderPublicKey.Y.Bytes())
			signatureStr := s.String()
			bt := &request.TransactionRequest{
				&sender, &recipient, &publicKeyStr, &value, &signatureStr}
			m, _ := json.Marshal(bt)
			buf := bytes.NewBuffer(m)
			endpoint := fmt.Sprintf("http://%s/transactions", n)
			client := &http.Client{}
			req, _ := http.NewRequest("PUT", endpoint, buf)
			resp, _ := client.Do(req)
			log.Printf("%v", resp)
		}
	}

	return isTransacted
}

func (bcr *blockchainRepository) AddTransaction(bc *entity.Blockchain, sender string, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.TransactionPool = append(bc.TransactionPool, t)
		return true
	}

	if bcr.VerifyTransactionSignature(bc, senderPublicKey, s, t) {
		if bcr.CalculateTotalAmount(bc, sender) < value {
			log.Println("ERROR: Not enough balance in a wallet")
			return false
		}
		bc.TransactionPool = append(bc.TransactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false

}

func (bcr *blockchainRepository) VerifyTransactionSignature(bc *entity.Blockchain,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *entity.Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bcr *blockchainRepository) CopyTransactionPool(bc *entity.Blockchain) []*entity.Transaction {
	transactions := make([]*entity.Transaction, 0)
	for _, t := range bc.TransactionPool {
		transactions = append(transactions,
			NewTransaction(t.SenderBlockchainAddress,
				t.RecipientBlockchainAddress,
				t.Value))
	}
	return transactions
}

func (bcr *blockchainRepository) ValidProof(bc *entity.Blockchain, br repository.BlockRepository, nonce int, previousHash [32]byte, transactions []*entity.Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := entity.Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", br.Hash(&guessBlock))
	return guessHashStr[:difficulty] == zeros
}

func (bcr *blockchainRepository) ProofOfWork(bc *entity.Blockchain, br repository.BlockRepository) int {
	transactions := bcr.CopyTransactionPool(bc)
	previousHash := br.Hash(bcr.LastBlock(bc))
	nonce := 0
	for !bcr.ValidProof(bc, br, nonce, previousHash, transactions, MINING_DIFFICULTY) {
		nonce += 1
	}
	return nonce
}

func (bcr *blockchainRepository) Mining(bc *entity.Blockchain, br repository.BlockRepository) bool {
	bc.Mux.Lock()
	defer bc.Mux.Unlock()

	/*
		if len(bc.transactionPool) == 0 {
			return false
		}
	*/

	bcr.AddTransaction(bc, MINING_SENDER, bc.BlockchainAddress, MINING_REWARD, nil, nil)
	nonce := bcr.ProofOfWork(bc, br)
	previousHash := br.Hash(bcr.LastBlock(bc))
	bcr.CreateBlock(bc, nonce, previousHash)
	log.Println("action=mining, status=success")

	for _, n := range bc.Neighbors {
		endpoint := fmt.Sprintf("http://%s/consensus", n)
		client := &http.Client{}
		req, _ := http.NewRequest("PUT", endpoint, nil)
		resp, _ := client.Do(req)
		log.Printf("%v", resp)
	}

	return true
}

func (bcr *blockchainRepository) StartMining(bc *entity.Blockchain, br repository.BlockRepository) {
	bcr.Mining(bc, br)
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, func() {
		bcr.StartMining(bc, br)
	})
}

func (bcr *blockchainRepository) CalculateTotalAmount(bc *entity.Blockchain, blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.Chain {
		for _, t := range b.Transactions {
			value := t.Value
			if blockchainAddress == t.RecipientBlockchainAddress {
				totalAmount += value
			}

			if blockchainAddress == t.SenderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (bcr *blockchainRepository) ValidChain(bc *entity.Blockchain, br repository.BlockRepository, chain []*entity.Block) bool {
	preBlock := chain[0]
	currentIndex := 1
	for currentIndex < len(chain) {
		b := chain[currentIndex]
		if b.PreviousHash != br.Hash(preBlock) {
			return false
		}

		if !bcr.ValidProof(bc, br, br.Nonce(b), br.PreviousHash(b), br.Transactions(b), MINING_DIFFICULTY) {
			return false
		}

		preBlock = b
		currentIndex += 1
	}
	return true
}

func (bcr *blockchainRepository) ResolveConflicts(bc *entity.Blockchain, br repository.BlockRepository) bool {
	var longestChain []*entity.Block = nil
	maxLength := len(bc.Chain)

	for _, n := range bc.Neighbors {
		endpoint := fmt.Sprintf("http://%s/chain", n)
		resp, _ := http.Get(endpoint)
		if resp.StatusCode == 200 {
			var bcResp entity.Blockchain
			decoder := json.NewDecoder(resp.Body)
			_ = decoder.Decode(&bcResp)

			chain := bcr.Chain(&bcResp)

			if len(chain) > maxLength && bcr.ValidChain(bc, br, chain) {
				maxLength = len(chain)
				longestChain = chain
			}
		}
	}

	if longestChain != nil {
		bc.Chain = longestChain
		log.Printf("Resovle confilicts replaced")
		return true
	}
	log.Printf("Resovle conflicts not replaced")
	return false
}
