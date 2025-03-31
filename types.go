package token

import "github.com/blocto/solana-go-sdk/common"

type MintParams struct {
	Metadata             Metadata
	DisableFutureMinting bool
	EnableFreeze         bool
}

type TransferParams struct {
	Recipient         common.PublicKey // Recipient receiver address
	Amount            float64
	CheckTokenProgram bool // CheckTokenProgram check token initialized on chain before create transfer program
}

type Metadata struct {
	Name   string
	Symbol string
	URI    string
	Supply float64
}
