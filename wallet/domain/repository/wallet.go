package repository

import (
	"crypto/ecdsa"

	"go-blockchain/wallet/domain/entity"
)

type WalletRepository interface {
	PrivateKey(w *entity.Wallet) *ecdsa.PrivateKey
	PrivateKeyStr(w *entity.Wallet) string
	PublicKey(w *entity.Wallet) *ecdsa.PublicKey
	PublicKeyStr(w *entity.Wallet) string
	BlockchainAddress(w *entity.Wallet) string
	MarshalJSON(w *entity.Wallet) ([]byte, error)
}
