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

func TestEdit(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	var (
		denomID      = "denom-id"
		tokenID      = "token-id"
		tokenURI     = "token-uri"
		newTokenData = "new-token-data"
		testAddrs    = testutil.CreateTestAddrs(2)
		denomCreator = testAddrs[0]
		nftOwner     = testAddrs[1]
		err          error
	)

	err = keeper.IssueDenom(ctx, denomID, "name", "schema", denomCreator)
	require.NoError(t, err)

	err = keeper.MintNFT(ctx, denomID, tokenID, "name", tokenURI, "data", nftOwner)
	require.NoError(t, err)

	type args struct {
		denomID string
		tokenID string
		nm      string
		uri     string
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
				denomID: "invalid-denomID",
				tokenID: "",
				nm:      "",
				uri:     "",
				data:    "",
				owner:   nftOwner,
			},
			err: irismodtypes.ErrUnknownCollection,
		},
		{
			desc: "not found nft by invalid tokenID",
			args: args{
				denomID: denomID,
				tokenID: "invalid-tokenID",
				nm:      "",
				uri:     "",
				data:    "",
				owner:   nftOwner,
			},
			err: irismodtypes.ErrUnknownCollection,
		},
		{
			desc: "attempt to change uri",
			args: args{
				denomID: denomID,
				tokenID: tokenID,
				nm:      "",
				uri:     "invalid-uri",
				data:    "",
				owner:   nftOwner,
			},
			err: irismodtypes.ErrInvalidTokenURI,
		},
		{
			desc: "valid",
			args: args{
				denomID: denomID,
				tokenID: tokenID,
				nm:      "",
				uri:     tokenURI,
				data:    newTokenData,
				owner:   nftOwner,
			},
			expectedTokenData: newTokenData,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err = keeper.EditNFT(ctx, tc.args.denomID, tc.args.tokenID, tc.args.nm, tc.args.uri, tc.args.data, tc.args.owner)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				nft, err := keeper.GetNFT(ctx, tc.args.denomID, tc.args.tokenID)
				require.NoError(t, err)
				require.Equal(t, tc.expectedTokenData, nft.GetData())
			}
		})
	}
}
