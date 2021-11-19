package entity

import "sync"

type Blockchain struct {
	TransactionPool   []*Transaction
	Chain             []*Block
	BlockchainAddress string
	Port              uint16
	Mux               sync.Mutex
	Neighbors         []string
	MuxNeighbors      sync.Mutex
}
