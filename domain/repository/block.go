package repository

import "go-blockchain/domain/entity"

type BlockRepository interface {
	PreviousHash(b *entity.Block) [32]byte
	Nonce(b *entity.Block) int
	Transaction(b *entity.Block) []*entity.Transaction
	Print(b *entity.Block)
	Hash(b *entity.Block) [32]byte
	MarshalJSON(b *entity.Block) ([]byte, error)
	UnmarshalJSON(b *entity.Block, data []byte) error
}
