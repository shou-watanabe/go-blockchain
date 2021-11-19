package repository

import (
	"crypto/ecdsa"

	"go-blockchain/blockchain/domain/entity"
	"go-blockchain/utils"
)

type BlockchainRepository interface {
	Chain(bc *entity.Blockchain) []*entity.Block
	Run(bc *entity.Blockchain, br BlockRepository)
	SetNeighbors(bc *entity.Blockchain)
	SyncNeighbors(bc *entity.Blockchain)
	StartSyncNeighbors(bc *entity.Blockchain)
	TransactionPool(bc *entity.Blockchain) []*entity.Transaction
	ClearTransactionPool(bc *entity.Blockchain)
	MarshalJSON(bc *entity.Blockchain) ([]byte, error)
	UnmarshalJSON(bc *entity.Blockchain, data []byte) error
	CreateBlock(bc *entity.Blockchain, nonce int, previousHash [32]byte) *entity.Block
	LastBlock(bc *entity.Blockchain) *entity.Block
	Print(br BlockRepository, tr TransactionRepository, bc *entity.Blockchain)
	CreateTransaction(bc *entity.Blockchain, sender string, recipient string, value float32,
		senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool
	AddTransaction(bc *entity.Blockchain, sender string, recipient string, value float32,
		senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool
	VerifyTransactionSignature(bc *entity.Blockchain,
		senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *entity.Transaction) bool
	CopyTransactionPool(bc *entity.Blockchain) []*entity.Transaction
	ValidProof(bc *entity.Blockchain, br BlockRepository, nonce int, previousHash [32]byte, transactions []*entity.Transaction, difficulty int) bool
	ProofOfWork(bc *entity.Blockchain, br BlockRepository) int
	Mining(bc *entity.Blockchain, br BlockRepository) bool
	StartMining(bc *entity.Blockchain, br BlockRepository)
	CalculateTotalAmount(bc *entity.Blockchain, blockchainAddress string) float32
	ValidChain(bc *entity.Blockchain, br BlockRepository, chain []*entity.Block) bool
	ResolveConflicts(bc *entity.Blockchain, br BlockRepository) bool
}
