package main

import (
	"fmt"
	"os"

	polymarket "github.com/0xfakeSpike/polymarket-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run ./examples/orderbook <token_id>")
		os.Exit(2)
	}

	c, err := polymarket.NewPublicClient()
	if err != nil {
		panic(err)
	}
	book, err := c.GetOrderBook(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("bids=%d asks=%d\n", len(book.Bids), len(book.Asks))
}
