package repository

import (
	"net/http"

	"go-blockchain/blockchain/domain/entity"
	"go-blockchain/wallet/domain/repository"
)

type BlockchainServerRepository interface {
	Port(bs *entity.BlockchainServer) uint16
	GetBlockchain(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, wr repository.WalletRepository) *entity.Blockchain
	GetChain(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, wr repository.WalletRepository, w http.ResponseWriter, req *http.Request)
	Transactions(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, wr repository.WalletRepository, w http.ResponseWriter, req *http.Request)
	Mine(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, wr repository.WalletRepository, w http.ResponseWriter, req *http.Request)
	StartMine(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, wr repository.WalletRepository, w http.ResponseWriter, req *http.Request)
	Amount(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, wr repository.WalletRepository, w http.ResponseWriter, req *http.Request)
	Consensus(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, wr repository.WalletRepository, w http.ResponseWriter, req *http.Request)
	Run(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, wr repository.WalletRepository)
}
