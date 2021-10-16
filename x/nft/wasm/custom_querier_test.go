package wasm

import (
	"encoding/json"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos-gaminghub/nibiru/testutil"
	"github.com/cosmos-gaminghub/nibiru/x/nft/keeper"
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

func setupKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	keeper := keeper.NewKeeper(
		codec.NewProtoCodec(registry),
		storeKey,
		memStoreKey,
		nil,
		nil,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return keeper, ctx
}

func TestCustomQuerier(t *testing.T) {
	var (
		keeper, ctx = setupKeeper(t)
		querier     = DefaultCustomQuerier(keeper).Querier()
		sender      = testutil.CreateTestAddrs(1)[0]
		msgDenom    = types.NewMsgIssueDenom("denom-id", "name", "schema", sender.String())

		_denomQuery     = types.DenomQuery{DenomId: msgDenom.GetDenomId()}
		_denomRequest   = types.DenomRequest{Denom: &_denomQuery}
		nftDenomRequest = types.NftDenomRequest{Nft: &_denomRequest}

		_noneDenomQuery     = types.DenomQuery{DenomId: "not-exist"}
		_noneDenomRequest   = types.DenomRequest{Denom: &_noneDenomQuery}
		noneNftDenomRequest = types.NftDenomRequest{Nft: &_noneDenomRequest}
	)

	err := keeper.IssueDenom(ctx, msgDenom)
	require.NoError(t, err)

	require.NoError(t, err)
	reqDenomExistByte, err := json.Marshal(nftDenomRequest)
	require.NoError(t, err)
	reqDenomNotExistByte, err := json.Marshal(noneNftDenomRequest)
	require.NoError(t, err)

	denomByte, err := json.Marshal(irismodtypes.QueryDenomResponse{
		Denom: &irismodtypes.Denom{
			Id:      msgDenom.GetDenomId(),
			Name:    msgDenom.GetName(),
			Schema:  msgDenom.GetSchema(),
			Creator: msgDenom.GetSender(),
		},
	})
	require.NoError(t, err)
	emptyDenomByte, err := json.Marshal(irismodtypes.QueryDenomResponse{})
	require.NoError(t, err)

	for _, tc := range []struct {
		desc     string
		request  json.RawMessage
		expected []byte
		err      error
	}{
		{
			desc:     "denom exist",
			request:  json.RawMessage(reqDenomExistByte),
			expected: denomByte,
		},
		{
			desc:     "denom not exist",
			request:  json.RawMessage(reqDenomNotExistByte),
			expected: emptyDenomByte,
		},
		{
			desc:    "custom",
			request: json.RawMessage([]byte("custom")),
			err:     wasmvmtypes.UnsupportedRequest{Kind: "custom"},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			msgs, err := querier(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, msgs)
			}
		})
	}
}
