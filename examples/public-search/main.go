package main

import (
	"fmt"

	polymarket "github.com/0xfakeSpike/polymarket-go"
)

func main() {
	c, err := polymarket.NewPublicClient()
	if err != nil {
		panic(err)
	}
	events, err := c.SearchEventsWithQuery("election")
	if err != nil {
		panic(err)
	}
	fmt.Println("matched events:", len(events))
}
