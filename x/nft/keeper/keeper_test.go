package keeper

import (
	"testing"

	"github.com/cosmos-gaminghub/nibiru/testutil"
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	keeper := NewKeeper(
		codec.NewProtoCodec(registry),
		storeKey,
		memStoreKey,
		nil,
		nil,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return keeper, ctx
}

func TestNewTokenID(t *testing.T) {
	var (
		keeper, ctx = setupKeeper(t)
		owner       = testutil.CreateTestAddrs(1)[0]
		denomID     = "denomID"
	)

	keeper.IssueDenom(ctx, types.NewMsgIssueDenom(denomID, "name", "schema", owner.String()))

	for _, tc := range []struct {
		desc       string
		denomID    string
		prepare    func()
		expectedID types.TokenID
		err        error
	}{
		{
			desc:    "invalid denomID",
			denomID: "invalid",
			prepare: func() {},
			err:     irismodtypes.ErrInvalidDenom,
		},
		{
			desc:       "first id",
			denomID:    denomID,
			prepare:    func() {},
			expectedID: types.TokenID(types.MIN_TOKEN_ID),
		},
		{
			desc:    "second id",
			denomID: denomID,
			prepare: func() {
				keeper.MintNFT(ctx, types.NewMsgMintNFT(denomID, "name", "token-uri", "data", owner.String(), owner.String()))
			},
			expectedID: types.TokenID(types.MIN_TOKEN_ID + 1),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			tc.prepare()
			id, err := keeper.NewTokenID(ctx, tc.denomID)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedID, id)
			}
		})
	}
}

func TestIssueDenomMintEditTransferBurnNFT(t *testing.T) {
	var (
		keeper, ctx = setupKeeper(t)
		owner       = testutil.CreateTestAddrs(1)[0]
		recipient   = testutil.CreateTestAddrs(2)[1]
	)

	//------test IssueDenom()-------------
	denomID := "denomID"
	err := keeper.IssueDenom(ctx, types.NewMsgIssueDenom(denomID, "name", "schema", owner.String()))
	require.NoError(t, err)

	//------test MintNFT()-------------
	expectedTokenid := uint64(types.MIN_TOKEN_ID)
	tokenID, err := keeper.MintNFT(ctx, types.NewMsgMintNFT(denomID, "name", "token-uri", "data", owner.String(), owner.String()))
	require.NoError(t, err)
	require.Equal(t, expectedTokenid, tokenID)

	//------test EditNFT()-------------
	expectedData := "new-token-data"
	err = keeper.EditNFT(ctx, types.NewMsgEditNFT(denomID, tokenID, "new-name", expectedData, owner.String()))
	require.NoError(t, err)
	nft, _ := keeper.GetNFT(ctx, denomID, tokenID)
	require.Equal(t, expectedData, nft.GetData())

	//------test TransferNFT()-------------
	err = keeper.TransferNFT(ctx, types.NewMsgTransferNFT(denomID, tokenID, owner.String(), recipient.String()))
	require.NoError(t, err)
	nft, err = keeper.GetNFT(ctx, denomID, tokenID)
	require.NoError(t, err)
	require.Equal(t, nft.GetOwner(), recipient)

	//------test BurnNFT()-------------
	err = keeper.BurnNFT(ctx, types.NewMsgBurnNFT(recipient.String(), denomID, tokenID))
	require.NoError(t, err)
	_, err = keeper.GetNFT(ctx, denomID, tokenID)
	require.ErrorIs(t, err, irismodtypes.ErrUnknownCollection)
}
