package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// EditNFT updates an already existing NFT
// overwrite irismod, restrict changing token URI
func (k Keeper) EditNFT(
	ctx sdk.Context, denomID, tokenID, tokenNm,
	tokenURI, tokenData string, owner sdk.AccAddress,
) error {
	nft, err := k.GetNFT(ctx, denomID, tokenID)
	if err != nil {
		return err
	}

	if nft.GetURI() != tokenURI {
		return sdkerrors.Wrapf(irismodtypes.ErrInvalidTokenURI, "changing token URI(%s) is restricted", tokenURI)
	}

	return k.Keeper.EditNFT(ctx, denomID, tokenID, tokenNm, tokenURI, tokenData, owner)
}
