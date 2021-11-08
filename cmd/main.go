package main

import (
	"fmt"
	"go-blockchain/utils"
)

func main() {
	fmt.Println(utils.IsFoundHost("127.0.0.1", 5000))
}
