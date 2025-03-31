package token

import (
	"context"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMint(t *testing.T) {
	ctx := context.Background()
	owner, err := NewAccountFromHex(os.Getenv("TEST_PRIVATE_KEY"))
	require.NoError(t, err)

	cli := NewClient(Devnet)
	require.True(t, cli.IsHealthy())

	tok := NewTokenManager(cli, owner, owner, common.PublicKey{})

	tests := []struct {
		Name   string
		Params MintParams
	}{
		{
			Name: "Basic mint token",
			Params: MintParams{
				Metadata: Metadata{
					Name:   "Solana Token 001",
					Symbol: "SMY1",
					URI:    "https://raw.githubusercontent.com/Ja7ad/sol-token/refs/heads/main/test/metadata.json",
					Supply: 1000,
				},
				DisableFutureMinting: false,
				EnableFreeze:         false,
			},
		},
		{
			Name: "Mint token with disableFutureMinting",
			Params: MintParams{
				Metadata: Metadata{
					Name:   "Solana Token 002",
					Symbol: "SMY2",
					URI:    "https://raw.githubusercontent.com/Ja7ad/sol-token/refs/heads/main/test/metadata.json",

					Supply: 2000,
				},
				DisableFutureMinting: true,
				EnableFreeze:         false,
			},
		},
		{
			Name: "Mint token with freeze token",
			Params: MintParams{
				Metadata: Metadata{
					Name:   "Solana Token 003",
					Symbol: "SMY3",
					URI:    "https://raw.githubusercontent.com/Ja7ad/sol-token/refs/heads/main/test/metadata.json",

					Supply: 3000,
				},
				DisableFutureMinting: false,
				EnableFreeze:         true,
			},
		},
		{
			Name: "Mint token with freeze token and disableFutureMinting",
			Params: MintParams{
				Metadata: Metadata{
					Name:   "Solana Token 004",
					Symbol: "SMY4",
					URI:    "https://raw.githubusercontent.com/Ja7ad/sol-token/refs/heads/main/test/metadata.json",

					Supply: 5000,
				},
				DisableFutureMinting: true,
				EnableFreeze:         true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tokenAddr, hash, err := tok.Mint(ctx, tt.Params)
			require.NoError(t, err)
			require.NotEmpty(t, tokenAddr)
			require.NotEmpty(t, hash)

			t.Logf("token: %s, tokenAddr: %s, hash: %s", tt.Params.Metadata.Name, tokenAddr.String(), hash)
		})
	}
}

func TestMintAndTransfer(t *testing.T) {
	ctx := context.Background()
	owner, err := NewAccountFromHex(os.Getenv("TEST_PRIVATE_KEY"))
	require.NoError(t, err)

	cli := NewClient(Devnet)
	require.True(t, cli.IsHealthy())

	tok := NewTokenManager(cli, owner, owner, common.PublicKey{})

	params := MintParams{
		Metadata: Metadata{
			Name:   "Solana Token Test Transfer Final",
			Symbol: "SMO",
			URI:    "https://raw.githubusercontent.com/Ja7ad/sol-token/refs/heads/main/test/metadata.json",
			Supply: 1000,
		},
		DisableFutureMinting: false,
		EnableFreeze:         false,
	}

	tokAddr, tx, err := tok.Mint(ctx, params)
	require.NoError(t, err)
	require.NotEmpty(t, tokAddr)
	require.NotEmpty(t, tx)

	t.Log("tokenAddr: ", tokAddr.String())
	t.Log("tx: ", tx)

	recipient := NewAccount()
	amount := 250.0

	txHash, err := tok.Transfer(ctx, TransferParams{
		Amount:            amount,
		Recipient:         recipient.PublicKey(),
		CheckTokenProgram: true,
	})
	require.NoError(t, err)
	require.NotEmpty(t, txHash)

	t.Log("txHash:", txHash)
	t.Log("recipient:", recipient.PublicKey())
}
