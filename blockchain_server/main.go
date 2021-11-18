package main

import (
	"flag"
	"go-blockchain/infra/repository"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	bsr := repository.NewBlockchainServerRepository()
	bcr := repository.NewBlockchainRepository()
	br := repository.NewBlockRepository()

	port := flag.Uint("port", 5000, "TCP Port Number for Blockchain Server")
	flag.Parse()
	bs := repository.NewBlockchainServer(uint16(*port))
	bsr.Run(bs, bcr, br)
}
