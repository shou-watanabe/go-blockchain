package entity

import "crypto/ecdsa"

type Transaction struct {
	SenderPrivateKey           *ecdsa.PrivateKey
	SenderPublicKey            *ecdsa.PublicKey
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}
