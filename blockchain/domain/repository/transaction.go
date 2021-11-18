package repository

import "go-blockchain/blockchain/domain/entity"

type TransactionRepository interface {
	Print(t *entity.Transaction)
	MarshalJSON(t *entity.Transaction) ([]byte, error)
	UnmarshalJSON(t *entity.Transaction, data []byte) error
}
