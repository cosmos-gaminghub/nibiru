package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	irismodexported "github.com/irisnet/irismod/modules/nft/exported"
	irismodkeeper "github.com/irisnet/irismod/modules/nft/keeper"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey

		irisKeeper irismodkeeper.Keeper
		// this line is used by starport scaffolding # ibc/keeper/attribute

		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
	// this line is used by starport scaffolding # ibc/keeper/parameter
	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper,
) *Keeper {
	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		irisKeeper: irismodkeeper.NewKeeper(cdc, storeKey),
		// this line is used by starport scaffolding # ibc/keeper/return
		accountKeeper: accountKeeper, bankKeeper: bankKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) IrisKeeper() irismodkeeper.Keeper {
	return k.irisKeeper
}

func (k Keeper) NewDenomID(ctx sdk.Context) types.DenomID {
	return types.DenomID(k.GetDenomCount(ctx) + types.MIN_DENOM_ID)
}

func (k Keeper) NewTokenID(ctx sdk.Context, denomID uint64) (tokenID types.TokenID, err error) {
	if !k.irisKeeper.HasDenomID(ctx, types.DenomID(denomID).String()) {
		err = sdkerrors.Wrapf(irismodtypes.ErrInvalidDenom, "denom ID %s not exists", types.DenomID(denomID).String())
		return
	}

	tokenID = types.TokenID(k.GetNFTCount(ctx, denomID) + types.MIN_TOKEN_ID)
	return
}

func (k Keeper) GetDenomCount(ctx sdk.Context) uint64 {
	return uint64(len(k.irisKeeper.GetDenoms(ctx)))
}

func (k Keeper) GetNFTCount(ctx sdk.Context, denomID uint64) uint64 {
	return uint64(len(k.irisKeeper.GetNFTs(ctx, types.DenomID(denomID).String())))
}

// GetNFT gets the the specified NFT
func (k Keeper) GetNFT(ctx sdk.Context, denomID, tokenID uint64) (irismodexported.NFT, error) {
	return k.irisKeeper.GetNFT(ctx, types.DenomID(denomID).String(), types.TokenID(tokenID).String())
}

// IssueDeno issues a denom according to the given params
func (k Keeper) IssueDenom(ctx sdk.Context, msg *types.MsgIssueDenom) (uint64, error) {
	irisMsg, err := k.toIrisMsgIssueDenom(ctx, msg)
	if err != nil {
		return 0, err
	}

	creator, err := sdk.AccAddressFromBech32(irisMsg.Sender)
	if err != nil {
		return 0, err
	}

	if err = k.irisKeeper.IssueDenom(ctx, irisMsg.Id, irisMsg.Name, irisMsg.Schema, creator); err != nil {
		return 0, err
	}

	denomID, _ := types.ToDenomID(irisMsg.Id)

	return denomID.Uint64(), nil
}

// MintNFT mints an NFT and manages the NFT's existence within Collections and Owners
func (k Keeper) MintNFT(ctx sdk.Context, msg *types.MsgMintNFT) (uint64, error) {
	irisMsg, err := k.toIrisMsgMintNFT(ctx, msg)
	if err != nil {
		return 0, err
	}

	receipent, err := sdk.AccAddressFromBech32(irisMsg.Recipient)
	if err != nil {
		return 0, err
	}

	if err = k.irisKeeper.MintNFT(ctx, irisMsg.DenomId, irisMsg.Id, irisMsg.Name, irisMsg.URI, irisMsg.Data, receipent); err != nil {
		return 0, err
	}

	tokenID, _ := types.ToTokenID(irisMsg.Id)

	return tokenID.Uint64(), nil
}

// EditNFT updates an already existing NFT
func (k Keeper) EditNFT(ctx sdk.Context, msg *types.MsgEditNFT) error {
	irisMsg, err := k.toIrisMsgEditNFT(ctx, msg)
	if err != nil {
		return err
	}

	owner, err := sdk.AccAddressFromBech32(irisMsg.Sender)
	if err != nil {
		return err
	}

	return k.irisKeeper.EditNFT(ctx, irisMsg.DenomId, irisMsg.Id, irisMsg.Name, irisMsg.URI, irisMsg.Data, owner)
}

// TransferNFT transfers the ownership of the given NFT to the new owner
func (k Keeper) TransferNFT(ctx sdk.Context, msg *types.MsgTransferNFT) error {
	irisMsg, err := k.toIrisMsgTransferNFT(ctx, msg)
	if err != nil {
		return err
	}

	owner, err := sdk.AccAddressFromBech32(irisMsg.Sender)
	if err != nil {
		return err
	}

	receipent, err := sdk.AccAddressFromBech32(irisMsg.Recipient)
	if err != nil {
		return err
	}

	return k.irisKeeper.TransferOwner(ctx, irisMsg.DenomId, irisMsg.Id, irisMsg.Name, irisMsg.URI, irisMsg.Data, owner, receipent)
}

// BurnNFT deletes a specified NFT
func (k Keeper) BurnNFT(ctx sdk.Context, msg *types.MsgBurnNFT) error {
	irisMsg, err := k.toIrisMsgBurnNFT(ctx, msg)
	if err != nil {
		return err
	}

	owner, err := sdk.AccAddressFromBech32(irisMsg.Sender)
	if err != nil {
		return err
	}

	return k.irisKeeper.BurnNFT(ctx, irisMsg.DenomId, irisMsg.Id, owner)
}
