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

func (k Keeper) GetDenomCount(ctx sdk.Context) uint64 {
	return uint64(len(k.GetDenoms(ctx)))
}

func (k Keeper) GetNFTCount(ctx sdk.Context, denomID uint64) uint64 {
	return uint64(len(k.GetNFTs(ctx, fmt.Sprintf("%d", denomID))))
}

// IssueDenomn issues a denom according to the given params
func (k Keeper) IssueDenomn(ctx sdk.Context,
	name, schema string, creator sdk.AccAddress) (uint64, error) {
	count := k.GetDenomCount(ctx)
	// must longer thant 3 length
	count += 100
	err := k.irisKeeper().IssueDenom(ctx, fmt.Sprintf("%d", count), name, schema, creator)
	return count, err
}

// IssueDenom return error
// Override irismod, restricting using this function
func (k Keeper) IssueDenom(ctx sdk.Context,
	id, name, schema string,
	creator sdk.AccAddress) error {
	return sdkerrors.Wrap(types.ErrRestricted, "please use another IssueDenomn function")
}

// GetNFTn gets the the specified NFT
func (k Keeper) GetNFTn(ctx sdk.Context, denomID, tokenID uint64) (irismodexported.NFT, error) {
	return k.irisKeeper().GetNFT(ctx, fmt.Sprintf("%d", denomID), fmt.Sprintf("%d", tokenID))
}

// GetNFT return error
// Override irismod, restricting using this function
func (k Keeper) GetNFT(ctx sdk.Context, denomID, tokenID string) (nft irismodexported.NFT, err error) {
	err = sdkerrors.Wrap(types.ErrRestricted, "please use another GetNFTn function")
	return
}

// MintNFTn mints an NFT and manages the NFT's existence within Collections and Owners
func (k Keeper) MintNFTn(
	ctx sdk.Context, denomID uint64, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) (count uint64, err error) {
	denomIDStr := fmt.Sprintf("%d", denomID)
	if !k.HasDenomID(ctx, denomIDStr) {
		err = sdkerrors.Wrapf(irismodtypes.ErrInvalidDenom, "denom ID %d not exists", denomID)
		return
	}
	count = k.GetNFTCount(ctx, denomID)
	// must longer thant 3 length
	count += 100
	k.irisKeeper().MintNFT(ctx, denomIDStr, fmt.Sprintf("%d", count), tokenNm, tokenURI, tokenData, owner)

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
	ctx sdk.Context, denomID, tokenID uint64, tokenNm, tokenData string, owner sdk.AccAddress,
) error {
	nft, err := k.GetNFTn(ctx, denomID, tokenID)
	if err != nil {
		return err
	}

	return k.irisKeeper().EditNFT(ctx, fmt.Sprintf("%d", denomID), fmt.Sprintf("%d", tokenID), tokenNm, nft.GetURI(), tokenData, owner)
}

// EditNFT return error
// Override irismod, restricting using this function
func (k Keeper) EditNFT(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) error {
	return sdkerrors.Wrap(types.ErrRestricted, "please use another EditNFTn function")
}

// TransferOwnern transfers the ownership of the given NFT to the new owner
func (k Keeper) TransferOwnern(
	ctx sdk.Context, denomID, tokenID uint64, srcOwner, dstOwner sdk.AccAddress,
) error {
	nft, err := k.GetNFTn(ctx, denomID, tokenID)
	if err != nil {
		return err
	}

	return k.irisKeeper().TransferOwner(ctx, fmt.Sprintf("%d", denomID), fmt.Sprintf("%d", tokenID), nft.GetName(), nft.GetURI(), nft.GetData(), srcOwner, dstOwner)
}

// TransferOwner return error
// Override irismod, restricting using this function
func (k Keeper) TransferOwner(
	ctx sdk.Context, denomID, tokenID, tokenNm, tokenURI,
	tokenData string, srcOwner, dstOwner sdk.AccAddress,
) error {
	return sdkerrors.Wrap(types.ErrRestricted, "please use another TransferOwnern function")
}

// BurnNFTn deletes a specified NFT
func (k Keeper) BurnNFTn(ctx sdk.Context, denomID, tokenID uint64, owner sdk.AccAddress) error {
	return k.irisKeeper().BurnNFT(ctx, fmt.Sprintf("%d", denomID), fmt.Sprintf("%d", tokenID), owner)
}

// BurnNFT return error
// Override irismod, restricting using this function
func (k Keeper) BurnNFT(ctx sdk.Context, denomID, tokenID string, owner sdk.AccAddress) error {
	return sdkerrors.Wrap(types.ErrRestricted, "please use another BurnNFTn function")
}
