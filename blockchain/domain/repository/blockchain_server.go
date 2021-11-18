package repository

import (
	"net/http"

	"go-blockchain/blockchain/domain/entity"
)

type BlockchainServerRepository interface {
	Port(bs *entity.BlockchainServer) uint16
	GetBlockchain(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository) *entity.Blockchain
	GetChain(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, w http.ResponseWriter, req *http.Request)
	Transactions(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, w http.ResponseWriter, req *http.Request)
	Mine(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, w http.ResponseWriter, req *http.Request)
	StartMine(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, w http.ResponseWriter, req *http.Request)
	Amount(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, w http.ResponseWriter, req *http.Request)
	Consensus(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository, w http.ResponseWriter, req *http.Request)
	Run(bs *entity.BlockchainServer, bcr BlockchainRepository, br BlockRepository)
}
