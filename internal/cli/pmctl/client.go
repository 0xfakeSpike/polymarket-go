package pmctl

import (
	"fmt"
	"os"

	"github.com/0xfakeSpike/polymarket-go"
)

func newClientFromFlags(public bool, privateKeyHex string) (*polymarket.Client, error) {
	if public {
		return polymarket.NewPublicClient()
	}
	if privateKeyHex == "" {
		privateKeyHex = os.Getenv("PMCTL_PRIVATE_KEY")
	}
	if privateKeyHex == "" {
		return nil, fmt.Errorf("private client requires -private-key or env PMCTL_PRIVATE_KEY")
	}
	return polymarket.NewClient(privateKeyHex)
}
