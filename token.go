package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/associated_token_account"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
	"time"
)

const _defaultDecimals uint8 = 10

type Token struct {
	client *Client
	payer  *Account
	owner  *Account
	token  common.PublicKey
}

// NewTokenManager create token to mint or transfer exists token
// token field can nil if not exists.
func NewTokenManager(client *Client, owner, payer *Account, tokenAddr common.PublicKey) *Token {
	return &Token{
		client: client,
		payer:  payer,
		owner:  owner,
		token:  tokenAddr,
	}
}

func (t *Token) Mint(ctx context.Context, params MintParams) (tokenAddr common.PublicKey, txHash string, err error) {
	if t.token.Bytes() != nil {
		return common.PublicKey{}, "", errors.New("token already exists, cannot mint")
	}

	mint := NewAccount()
	cli := t.client.Client()

	rent, err := cli.GetMinimumBalanceForRentExemption(ctx, token.MintAccountSize)
	if err != nil {
		return common.PublicKey{}, "", err
	}

	instructions := make([]types.Instruction, 0)

	// Instructions
	createMint := system.CreateAccount(system.CreateAccountParam{
		From:     t.payer.PublicKey(),
		New:      mint.PublicKey(),
		Lamports: rent,
		Space:    token.MintAccountSize,
		Owner:    common.TokenProgramID,
	})

	instructions = append(instructions, createMint)

	mintInitParam := token.InitializeMintParam{
		Mint:       mint.PublicKey(),
		Decimals:   _defaultDecimals,
		MintAuth:   t.payer.PublicKey(),
		FreezeAuth: nil,
	}

	if params.EnableFreeze {
		freezer := t.payer.PublicKey()
		mintInitParam.FreezeAuth = &freezer
	}

	initMint := token.InitializeMint(mintInitParam)

	instructions = append(instructions, initMint)

	ownerATA, _, err := common.FindAssociatedTokenAddress(t.owner.PublicKey(), mint.PublicKey())
	if err != nil {
		return common.PublicKey{}, "", err
	}

	createATA := associated_token_account.Create(associated_token_account.CreateParam{
		Funder:                 t.payer.PublicKey(),
		Owner:                  t.owner.PublicKey(),
		Mint:                   mint.PublicKey(),
		AssociatedTokenAccount: ownerATA,
	})

	instructions = append(instructions, createATA)

	mintTo := token.MintTo(token.MintToParam{
		Mint:    mint.PublicKey(),
		To:      ownerATA,
		Auth:    t.payer.PublicKey(),
		Amount:  ConvertToDecimals(params.Metadata.Supply, _defaultDecimals),
		Signers: []common.PublicKey{},
	})

	instructions = append(instructions, mintTo)

	blockhash, err := cli.GetLatestBlockhash(ctx)
	if err != nil {
		return common.PublicKey{}, "", fmt.Errorf("failed to get blockhash: %w", err)
	}

	metadataSeeds := [][]byte{
		[]byte("metadata"),
		common.MetaplexTokenMetaProgramID.Bytes(),
		mint.PublicKey().Bytes(),
	}
	metadataAccount, _, err := common.FindProgramAddress(metadataSeeds, common.MetaplexTokenMetaProgramID)
	if err != nil {
		return common.PublicKey{}, "", fmt.Errorf("failed to derive metadata account: %w", err)
	}

	createMetadataInst := token_metadata.CreateMetadataAccountV3(token_metadata.CreateMetadataAccountV3Param{
		Metadata:        metadataAccount,
		Mint:            mint.PublicKey(),
		MintAuthority:   t.payer.PublicKey(),
		Payer:           t.payer.PublicKey(),
		UpdateAuthority: t.owner.PublicKey(),
		Data: token_metadata.DataV2{
			Name:                 params.Metadata.Name,
			Symbol:               params.Metadata.Symbol,
			Uri:                  params.Metadata.URI,
			SellerFeeBasisPoints: 0,
			Creators: &[]token_metadata.Creator{
				{
					Address:  t.owner.PublicKey(),
					Verified: true,
					Share:    100,
				},
			},
		},
		IsMutable: true,
	})

	instructions = append(instructions, createMetadataInst)

	if params.DisableFutureMinting {
		setAuth := token.SetAuthority(token.SetAuthorityParam{
			Account:  mint.PublicKey(),
			NewAuth:  nil,
			AuthType: token.AuthorityTypeMintTokens,
			Auth:     t.payer.PublicKey(),
		})
		instructions = append(instructions, setAuth)
	}

	// Create transaction
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{t.payer.Account, mint.Account},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        t.payer.PublicKey(),
			RecentBlockhash: blockhash.Blockhash,
			Instructions:    instructions,
		}),
	})
	if err != nil {
		return common.PublicKey{}, "", fmt.Errorf("failed to build transaction: %w", err)
	}

	// Send transaction
	sig, err := cli.SendTransaction(ctx, tx)
	if err != nil {
		return common.PublicKey{}, "", fmt.Errorf("failed to send transaction: %w", err)
	}

	t.token = mint.PublicKey()

	return mint.PublicKey(), sig, nil
}

func (t *Token) Transfer(ctx context.Context, params TransferParams) (string, error) {
	if t.token.Bytes() == nil {
		return "", errors.New("token is not initialized, please mint new token")
	}

	cli := t.client.Client()

	if params.CheckTokenProgram {
		// check token initialized
		if err := waitForAccount(ctx, cli, t.token, 30*time.Second); err != nil {
			return "", err
		}
	}

	senderATA := getATA(t.owner.PublicKey(), t.token)
	receiverATA := getATA(params.Recipient, t.token)

	instructions := make([]types.Instruction, 0)

	createReceiverATA := associated_token_account.Create(associated_token_account.CreateParam{
		Funder:                 t.payer.PublicKey(),
		Owner:                  params.Recipient,
		Mint:                   t.token,
		AssociatedTokenAccount: receiverATA,
	})
	instructions = append(instructions, createReceiverATA)

	transfer := token.Transfer(token.TransferParam{
		From:   senderATA,
		To:     receiverATA,
		Auth:   t.owner.PublicKey(),
		Amount: ConvertToDecimals(params.Amount, _defaultDecimals),
	})

	instructions = append(instructions, transfer)

	blockhash, err := cli.GetLatestBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get blockhash: %w", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{t.payer.Account, t.owner.Account},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        t.payer.PublicKey(),
			RecentBlockhash: blockhash.Blockhash,
			Instructions:    instructions,
		}),
	})
	if err != nil {
		return "", fmt.Errorf("failed to build transfer tx: %w", err)
	}

	sig, err := cli.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transfer tx: %w", err)
	}

	return sig, nil
}

func (t *Token) Token() common.PublicKey {
	return t.token
}

func waitForAccount(ctx context.Context, cli *client.Client, addr common.PublicKey, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			info, err := cli.GetAccountInfo(ctx, addr.ToBase58())
			if err != nil {
				continue
			}
			if len(info.Data) > 0 {
				return nil
			}
		}
	}
}

func getATA(owner, mint common.PublicKey) common.PublicKey {
	ata, _, _ := common.FindAssociatedTokenAddress(owner, mint)
	return ata
}
