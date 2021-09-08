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

		irismodkeeper.Keeper
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
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
		Keeper:   irismodkeeper.NewKeeper(cdc, storeKey),
		// this line is used by starport scaffolding # ibc/keeper/return
		accountKeeper: accountKeeper, bankKeeper: bankKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) irisKeeper() irismodkeeper.Keeper {
	return k.Keeper
}

func (k Keeper) GetNFTCount(ctx sdk.Context, denom string) uint64 {
	return uint64(len(k.GetNFTs(ctx, denom)))
}

// GetNFTn gets the the specified NFT
func (k Keeper) GetNFTn(ctx sdk.Context, denomID string, tokenID uint64) (irismodexported.NFT, error) {
	return k.irisKeeper().GetNFT(ctx, denomID, fmt.Sprintf("%d", tokenID))
}

// GetNFT return error
// Override irismod, restricting using this function
func (k Keeper) GetNFT(ctx sdk.Context, denomID, tokenID string) (nft irismodexported.NFT, err error) {
	err = sdkerrors.Wrap(types.ErrRestricted, "please use another GetNFTn function")
	return
}

// MintNFTn mints an NFT and manages the NFT's existence within Collections and Owners
func (k Keeper) MintNFTn(
	ctx sdk.Context, denomID, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) (count uint64, err error) {
	if !k.HasDenomID(ctx, denomID) {
		err = sdkerrors.Wrapf(irismodtypes.ErrInvalidDenom, "denom ID %s not exists", denomID)
		return
	}
	count = k.GetNFTCount(ctx, denomID)
	count++
	k.irisKeeper().MintNFT(ctx, denomID, fmt.Sprintf("%d", count), tokenNm, tokenURI, tokenData, owner)

	return
}

// MintNFT return error
// Override irismod, restricting using this function
func (k Keeper) MintNFT(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) error {
	return sdkerrors.Wrap(types.ErrRestricted, "please use another MintNFTn function")
}

// EditNFTn updates an already existing NFT
// Override irismod, restrict changing token URI
func (k Keeper) EditNFTn(
	ctx sdk.Context, denomID string, tokenID uint64, tokenNm string,
	tokenURI, tokenData string, owner sdk.AccAddress,
) error {
	tokenIDStr := fmt.Sprintf("%d", tokenID)
	nft, err := k.irisKeeper().GetNFT(ctx, denomID, tokenIDStr)
	if err != nil {
		return err
	}

	if nft.GetURI() != tokenURI {
		return sdkerrors.Wrapf(irismodtypes.ErrInvalidTokenURI, "changing token URI(%s) is restricted", tokenURI)
	}

	return k.irisKeeper().EditNFT(ctx, denomID, tokenIDStr, tokenNm, tokenURI, tokenData, owner)
}

// EditNFT return error
// Override irismod, restricting using this function
func (k Keeper) EditNFT(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) error {
	return sdkerrors.Wrap(types.ErrRestricted, "please use another EditNFTn function")
}
