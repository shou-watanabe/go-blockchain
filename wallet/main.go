package main

import (
	"flag"
	"log"

	"go-blockchain/wallet/infra/repository"
)

func init() {
	log.SetPrefix("Wallet Server: ")
}

func main() {
	wsr := repository.NewWalletServerRepository()
	wr := repository.NewWalletRepository()
	tr := repository.NewTransactionRepository()

	port := flag.Uint("port", 8080, "TCP Port Number for Wallet Server")
	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Blockchain Gateway")
	flag.Parse()

	ws := repository.NewWalletServer(uint16(*port), *gateway)
	wsr.Run(ws, wr, tr)
}
