# 🪙 sol-token

A Golang package for creating and managing **ERC20-like SPL tokens** on the **Solana blockchain**.

> Easily mint, transfer, and manage tokens with optional metadata and supply locking. Built with [`blocto/solana-go-sdk`](https://github.com/blocto/solana-go-sdk).

---

## ✨ Features

- ✅ Mint new SPL tokens (ERC20-style)
- ✅ Set name, symbol, URI, decimals, and custom metadata
- ✅ Transfer tokens between accounts
- ✅ Automatically create associated token accounts (ATA)
- ✅ Lock mint authority (disable future minting)
- ✅ Optional freeze authority support
- ✅ Lightweight and developer-friendly

---

## 🛠️ Installation

```bash
go get github.com/ja7ad/sol-token
```

---

## 🚀 Quick Start

```go
client := NewClient(Devnet)
owner := token.NewAccount()
payer := token.NewAccount()
tok := token.NewTokenManager(client, owner, payer)

mintAddr, txHash, err := tok.Mint(ctx, token.MintParams{
    Metadata: token.Metadata{
        Name:     "MyToken",
        Symbol:   "MTK",
        URI:      "https://example.com/token.json",
        Supply:   1000.0,
    },
    DisableFutureMinting: true,
})
```

Transfer tokens:

```go
_, err = tok.Transfer(ctx, TransferParams{
		Amount:            amount,
		Recipient:         recipient.PublicKey(),
		CheckTokenProgram: true,
	})
```

---

## 🔐 Security Notes

- Only store the **mint account private key** if you plan to mint more tokens later.
- If `DisableFutureMinting` is `true`, minting will be permanently disabled.

## 🌐 Resources

- [Solana SPL Token Docs](https://spl.solana.com/token)
- [Token Metadata Standard](https://docs.metaplex.com/programs/token-metadata/overview)
