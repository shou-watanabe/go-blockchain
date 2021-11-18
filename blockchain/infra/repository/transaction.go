package repository

import (
	"encoding/json"
	"fmt"
	"strings"

	"go-blockchain/blockchain/domain/entity"
	"go-blockchain/blockchain/domain/repository"
)

type transactionRepository struct{}

func NewTransactionRepository() repository.TransactionRepository {
	return &transactionRepository{}
}

func NewTransaction(sender string, recipient string, value float32) *entity.Transaction {
	return &entity.Transaction{sender, recipient, value}
}

func (tr *transactionRepository) Print(t *entity.Transaction) {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" sender_blockchain_address      %s\n", t.SenderBlockchainAddress)
	fmt.Printf(" recipient_blockchain_address   %s\n", t.RecipientBlockchainAddress)
	fmt.Printf(" value                          %.1f\n", t.Value)
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

func (tr *transactionRepository) UnmarshalJSON(t *entity.Transaction, data []byte) error {
	v := &struct {
		Sender    *string  `json:"sender_blockchain_address"`
		Recipient *string  `json:"recipient_blockchain_address"`
		Value     *float32 `json:"value"`
	}{
		Sender:    &t.SenderBlockchainAddress,
		Recipient: &t.RecipientBlockchainAddress,
		Value:     &t.Value,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
