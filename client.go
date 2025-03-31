package token

import (
	"context"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/rpc"
)

const (
	// Devnet solana devnet rpc
	Devnet string = rpc.DevnetRPCEndpoint
	// Testnet solana testnet rpc
	Testnet string = rpc.TestnetRPCEndpoint
)

type Client struct {
	cli *client.Client
}

func NewClient(rpc string) *Client {
	return &Client{
		cli: client.NewClient(rpc),
	}
}

func (c *Client) IsHealthy() bool {
	block, err := c.cli.GetLatestBlockhash(context.Background())
	if err != nil {
		return false
	}

	return block.LatestValidBlockHeight > 0
}

func (c *Client) Client() *client.Client {
	return c.cli
}

func (c *Client) GetFaucet(ctx context.Context, acc *Account, amount float64) (string, error) {
	return c.cli.RequestAirdrop(ctx, acc.PublicKey().ToBase58(), ConvertToLamport(amount))
}
