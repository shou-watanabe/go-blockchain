package repository

import (
	"go-blockchain/wallet/domain/entity"
	"net/http"
)

type WalletServerRepository interface {
	Port(ws *entity.WalletServer) uint16
	Gateway(ws *entity.WalletServer) string
	Index(ws *entity.WalletServer, w http.ResponseWriter, req *http.Request)
	Wallet(wr WalletRepository, w http.ResponseWriter, req *http.Request)
	CreateTransaction(ws *entity.WalletServer, tr TransactionRepository, w http.ResponseWriter, req *http.Request)
	WalletAmount(ws *entity.WalletServer, w http.ResponseWriter, req *http.Request)
	Run(ws *entity.WalletServer, wr WalletRepository, tr TransactionRepository)
}
