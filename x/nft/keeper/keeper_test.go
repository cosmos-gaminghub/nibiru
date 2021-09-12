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

func TestIssueDenomnn(t *testing.T) {
	var (
		keeper, ctx     = setupKeeper(t)
		expectedDenomID = uint64(100)
	)
	id, err := keeper.IssueDenomn(ctx, "name", "schema", testutil.CreateTestAddrs(1)[0])
	require.NoError(t, err)
	require.Equal(t, expectedDenomID, id)
}

func TestIssueDenom(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	err := keeper.IssueDenom(ctx, "denomID", "name", "shema", testutil.CreateTestAddrs(1)[0])
	require.ErrorIs(t, err, types.ErrRestricted)
}

func TestGetNFTn(t *testing.T) {
	var (
		keeper, ctx = setupKeeper(t)
		owner       = testutil.CreateTestAddrs(1)[0]
	)
	denomID, err := keeper.IssueDenomn(ctx, "name", "schema", owner)
	require.NoError(t, err)
	tokenID, err := keeper.MintNFTn(ctx, denomID, "name", "token-uri", "data", owner)
	require.NoError(t, err)

	_, err = keeper.GetNFTn(ctx, denomID, tokenID)
	require.NoError(t, err)
}

func TestGetNFT(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	_, err := keeper.GetNFT(ctx, "denomID", "tokenID")
	require.ErrorIs(t, err, types.ErrRestricted)
}

func TestMintNFTn(t *testing.T) {
	var (
		keeper, ctx     = setupKeeper(t)
		owner           = testutil.CreateTestAddrs(1)[0]
		expectedTokenid = uint64(100)
	)
	denomID, err := keeper.IssueDenomn(ctx, "name", "schema", owner)
	require.NoError(t, err)

	tokenID, err := keeper.MintNFTn(ctx, denomID, "name", "token-uri", "data", owner)
	require.NoError(t, err)
	require.Equal(t, expectedTokenid, tokenID)
}

func TestMintNFT(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	err := keeper.MintNFT(ctx, "denomID", "tokenID", "nm", "uri", "data", testutil.CreateTestAddrs(1)[0])
	require.ErrorIs(t, err, types.ErrRestricted)
}

func TestEditNFTn(t *testing.T) {
	var (
		keeper, ctx  = setupKeeper(t)
		newTokenData = "new-token-data"
		testAddrs    = testutil.CreateTestAddrs(2)
		denomCreator = testAddrs[0]
		nftOwner     = testAddrs[1]
	)

	denomID, err := keeper.IssueDenomn(ctx, "name", "schema", denomCreator)
	require.NoError(t, err)

	tokenID, err := keeper.MintNFTn(ctx, denomID, "name", "token-uri", "data", nftOwner)
	require.NoError(t, err)

	type args struct {
		denomID uint64
		tokenID uint64
		nm      string
		data    string
		owner   sdk.AccAddress
	}

	for _, tc := range []struct {
		desc              string
		args              args
		expectedTokenData string
		err               error
	}{
		{
			desc: "not found nft by invalid denomID",
			args: args{
				denomID: 0,
				tokenID: tokenID,
				nm:      "",
				data:    "",
				owner:   nftOwner,
			},
			err: irismodtypes.ErrUnknownCollection,
		},
		{
			desc: "not found nft by invalid tokenID",
			args: args{
				denomID: denomID,
				tokenID: 0,
				nm:      "",
				data:    "",
				owner:   nftOwner,
			},
			err: irismodtypes.ErrUnknownCollection,
		},
		{
			desc: "valid",
			args: args{
				denomID: denomID,
				tokenID: tokenID,
				nm:      "",
				data:    newTokenData,
				owner:   nftOwner,
			},
			expectedTokenData: newTokenData,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err = keeper.EditNFTn(ctx, tc.args.denomID, tc.args.tokenID, tc.args.nm, tc.args.data, tc.args.owner)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				nft, err := keeper.GetNFTn(ctx, tc.args.denomID, tc.args.tokenID)
				require.NoError(t, err)
				require.Equal(t, tc.expectedTokenData, nft.GetData())
			}
		})
	}
}

func TestEditNFT(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	err := keeper.EditNFT(ctx, "denomID", "tokenID", "nm", "uri", "data", testutil.CreateTestAddrs(1)[0])
	require.ErrorIs(t, err, types.ErrRestricted)
}

func TestTransferOwnern(t *testing.T) {
	var (
		keeper, ctx = setupKeeper(t)
		owner       = testutil.CreateTestAddrs(1)[0]
		recipient   = testutil.CreateTestAddrs(2)[1]
	)
	denomID, err := keeper.IssueDenomn(ctx, "name", "schema", owner)
	require.NoError(t, err)

	tokenID, err := keeper.MintNFTn(ctx, denomID, "name", "token-uri", "data", owner)
	require.NoError(t, err)

	err = keeper.TransferOwnern(ctx, denomID, tokenID, owner, recipient)
	require.NoError(t, err)
	nft, err := keeper.GetNFTn(ctx, denomID, tokenID)
	require.NoError(t, err)
	require.Equal(t, nft.GetOwner(), recipient)
}

func TestTransferOwner(t *testing.T) {
	var (
		keeper, ctx      = setupKeeper(t)
		owner, recipient = testutil.CreateTestAddrs(1)[0], testutil.CreateTestAddrs(2)[1]
	)
	err := keeper.TransferOwner(ctx, "denomID", "tokenID", "nm", "uri", "data", owner, recipient)
	require.ErrorIs(t, err, types.ErrRestricted)
}

func TestBurnNFTn(t *testing.T) {
	var (
		keeper, ctx = setupKeeper(t)
		owner       = testutil.CreateTestAddrs(1)[0]
	)
	denomID, err := keeper.IssueDenomn(ctx, "name", "schema", owner)
	require.NoError(t, err)

	tokenID, err := keeper.MintNFTn(ctx, denomID, "name", "token-uri", "data", owner)
	require.NoError(t, err)

	err = keeper.BurnNFTn(ctx, denomID, tokenID, owner)
	require.NoError(t, err)
	_, err = keeper.GetNFTn(ctx, denomID, tokenID)
	require.ErrorIs(t, err, irismodtypes.ErrUnknownCollection)
}

func TestBurnNFT(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	err := keeper.BurnNFT(ctx, "denomID", "tokenID", testutil.CreateTestAddrs(1)[0])
	require.ErrorIs(t, err, types.ErrRestricted)
}
