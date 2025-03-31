package token

import (
	"context"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_IsHealthy(t *testing.T) {
	client := NewClient(Devnet)
	ok := client.IsHealthy()
	assert.True(t, ok, "Solana Devnet should be healthy")
}

func TestClient_GetFaucet(t *testing.T) {
	ctx := context.Background()
	client := NewClient(Devnet)

	acc := NewAccount()
	amount := 1.0

	tx, err := client.GetFaucet(ctx, acc, amount)

	require.NoError(t, err)
	assert.NotEmpty(t, tx)

blc:
	balance, err := client.cli.GetBalance(ctx, acc.PublicKey().ToBase58())
	require.NoError(t, err)

	if balance <= 0 {
		goto blc
	}

	assert.Equal(t, balance, ConvertToLamport(amount))

	t.Log("account public key:", acc.PublicKey().ToBase58())
	t.Log("account balance:", amount)
	t.Log("account private key:", hex.EncodeToString(acc.PrivateKey()[:]))
}
