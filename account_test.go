package token

import (
	"encoding/hex"
	"github.com/mr-tron/base58"
	"testing"

	"github.com/blocto/solana-go-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	acc := NewAccount()
	assert.NotNil(t, acc)
	assert.Len(t, acc.PrivateKey(), 64)
	assert.NotEqual(t, "", acc.PublicKey().ToBase58())
}

func TestNewAccountFromHex(t *testing.T) {
	original := NewAccount()
	hexPriv := hex.EncodeToString(original.PrivateKey())

	acc, err := NewAccountFromHex(hexPriv)
	assert.NoError(t, err)
	assert.Equal(t, original.PublicKey(), acc.PublicKey())
}

func TestNewAccountFromBase58(t *testing.T) {
	original := types.NewAccount()
	base58Priv := base58.Encode(original.PrivateKey)

	acc, err := NewAccountFromBase58(base58Priv)
	assert.NoError(t, err)
	assert.Equal(t, original.PublicKey, acc.PublicKey())
}

func TestNewAccountFromHex_Invalid(t *testing.T) {
	_, err := NewAccountFromHex("invalid-hex")
	assert.Error(t, err)
}

func TestNewAccountFromBase58_Invalid(t *testing.T) {
	_, err := NewAccountFromBase58("badbase58==")
	assert.Error(t, err)
}
