package repository

import (
	"go-blockchain/utils"
	"go-blockchain/wallet/domain/entity"
)

type TransactionRepository interface {
	GenerateSignature(t *entity.Transaction) *utils.Signature
	MarshalJSON(t *entity.Transaction) ([]byte, error)
}
