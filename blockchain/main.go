package main

import (
	"flag"
	"log"

	bir "go-blockchain/blockchain/infra/repository"
	wir "go-blockchain/wallet/infra/repository"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	bsr := bir.NewBlockchainServerRepository()
	bcr := bir.NewBlockchainRepository()
	br := bir.NewBlockRepository()
	wr := wir.NewWalletRepository()

	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Server")
	flag.Parse()
	bs := bir.NewBlockchainServer(uint16(*port))
	bsr.Run(bs, bcr, br, wr)
}
