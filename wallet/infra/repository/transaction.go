package repository

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"go-blockchain/utils"
	"go-blockchain/wallet/domain/entity"
	"go-blockchain/wallet/domain/repository"
)

type transactionRepository struct{}

func NewTransactionRepository() repository.TransactionRepository {
	return &transactionRepository{}
}

func NewTransaction(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey,
	sender string, recipient string, value float32) *entity.Transaction {
	return &entity.Transaction{SenderPrivateKey: privateKey, SenderPublicKey: publicKey, SenderBlockchainAddress: sender, RecipientBlockchainAddress: recipient, Value: value}
}

func (tr *transactionRepository) GenerateSignature(t *entity.Transaction) *utils.Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.SenderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}

func (tr *transactionRepository) MarshalJSON(t *entity.Transaction) ([]byte, error) {
	return json.Marshal(struct {
		Sender    string  `json:"sender_blockchain_address"`
		Recipient string  `json:"recipient_blockchain_address"`
		Value     float32 `json:"value"`
	}{
		Sender:    t.SenderBlockchainAddress,
		Recipient: t.RecipientBlockchainAddress,
		Value:     t.Value,
	})
}
