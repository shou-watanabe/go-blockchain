package entity

type Transaction struct {
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}
