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
		sender2     = testutil.CreateTestAddrs(2)[1]
		msgDenom    = types.NewMsgIssueDenom("denom-id1", "name1", "schema1", sender.String())
		msgDenom2   = types.NewMsgIssueDenom("denom-id2", "name2", "schema2", sender2.String())
		err         error

		_denomQuery = types.WasmDenomQuery{DenomId: msgDenom.GetDenomId()}
		_denomReq   = types.WasmDenomReq{Denom: &_denomQuery}
		nftDenomReq = types.WasmNftDenomReq{Nft: &_denomReq}

		_noneDenomQuery = types.WasmDenomQuery{DenomId: "not-exist"}
		_noneDenomReq   = types.WasmDenomReq{Denom: &_noneDenomQuery}
		noneNftDenomReq = types.WasmNftDenomReq{Nft: &_noneDenomReq}

		_denomAllQuery = types.WasmDenomAllQuery{}
		_denomAllReq   = types.WasmDenomAllReq{DenomAll: &_denomAllQuery}
		nftDenomAllReq = types.WasmNftDenomAllReq{Nft: &_denomAllReq}
	)

	err = keeper.IssueDenom(ctx, msgDenom)
	require.NoError(t, err)

	err = keeper.IssueDenom(ctx, msgDenom2)
	require.NoError(t, err)

	require.NoError(t, err)
	reqDenomExistByte, err := json.Marshal(nftDenomReq)
	require.NoError(t, err)
	reqDenomNotExistByte, err := json.Marshal(noneNftDenomReq)
	require.NoError(t, err)
	reqDenomAllByte, err := json.Marshal(nftDenomAllReq)
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
	denomAllByte, err := json.Marshal(irismodtypes.QueryDenomsResponse{
		Denoms: []irismodtypes.Denom{
			irismodtypes.Denom{
				Id:      msgDenom.GetDenomId(),
				Name:    msgDenom.GetName(),
				Schema:  msgDenom.GetSchema(),
				Creator: msgDenom.GetSender(),
			},
			irismodtypes.Denom{
				Id:      msgDenom2.GetDenomId(),
				Name:    msgDenom2.GetName(),
				Schema:  msgDenom2.GetSchema(),
				Creator: msgDenom2.GetSender(),
			},
		},
	})
	require.NoError(t, err)

	for _, tc := range []struct {
		desc     string
		request  json.RawMessage
		expected []byte
		err      error
	}{
		{
			desc:     "get denom when exist",
			request:  json.RawMessage(reqDenomExistByte),
			expected: denomByte,
		},
		{
			desc:     "get denom when not exist",
			request:  json.RawMessage(reqDenomNotExistByte),
			expected: emptyDenomByte,
		},
		{
			desc:     "get denom all",
			request:  json.RawMessage(reqDenomAllByte),
			expected: denomAllByte,
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
