package token

import (
	"crypto/ed25519"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
)

type Account struct {
	types.Account
}

func NewAccount() *Account {
	return &Account{types.NewAccount()}
}

func NewAccountFromHex(privHex string) (*Account, error) {
	acc, err := types.AccountFromHex(privHex)
	if err != nil {
		return nil, err
	}
	return &Account{acc}, nil
}

func NewAccountFromBase58(privBase58 string) (*Account, error) {
	acc, err := types.AccountFromBase58(privBase58)
	if err != nil {
		return nil, err
	}
	return &Account{acc}, nil
}

func (a *Account) PublicKey() common.PublicKey {
	return a.Account.PublicKey
}

func (a *Account) PrivateKey() ed25519.PrivateKey {
	return a.Account.PrivateKey
}
