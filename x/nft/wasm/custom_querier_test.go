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
		recipient   = testutil.CreateTestAddrs(3)[2]
		msgDenom    = types.NewMsgIssueDenom("denom-id1", "name1", "schema1", sender.String())
		msgDenom2   = types.NewMsgIssueDenom("denom-id2", "name2", "schema2", sender2.String())
		tokenid     = uint64(1)
		msgMint     = types.NewMsgMintNFT(msgDenom.DenomId, "name", "uri", "data", sender.String(), recipient.String())
		err         error

		_denomQuery = types.WasmDenomQuery{DenomId: msgDenom.GetDenomId()}
		_denomReq   = types.WasmDenomReq{Denom: &_denomQuery}
		nftDenomReq = types.WasmGameDenomReq{Nft: &_denomReq}

		_noneDenomQuery = types.WasmDenomQuery{DenomId: "not-exist"}
		_noneDenomReq   = types.WasmDenomReq{Denom: &_noneDenomQuery}
		noneNftDenomReq = types.WasmGameDenomReq{Nft: &_noneDenomReq}

		_denomAllQuery = types.WasmDenomAllQuery{}
		_denomAllReq   = types.WasmDenomAllReq{DenomAll: &_denomAllQuery}
		nftDenomAllReq = types.WasmGameDenomAllReq{Nft: &_denomAllReq}

		_nftRuery = types.WasmNftQuery{DenomId: msgMint.GetDenomId(), Id: tokenid}
		_nftReq   = types.WasmNftReq{Nft: &_nftRuery}
		nftReq    = types.WasmGameNftReq{Nft: &_nftReq}

		_noDenomNftRuery = types.WasmNftQuery{DenomId: "not-exist", Id: tokenid}
		_noDenomNftReq   = types.WasmNftReq{Nft: &_noDenomNftRuery}
		noDenomNftReq    = types.WasmGameNftReq{Nft: &_noDenomNftReq}

		_noIdNftRuery = types.WasmNftQuery{DenomId: msgMint.GetDenomId(), Id: 100}
		_noIdNftReq   = types.WasmNftReq{Nft: &_noIdNftRuery}
		noIdNftReq    = types.WasmGameNftReq{Nft: &_noIdNftReq}
	)

	err = keeper.IssueDenom(ctx, msgDenom)
	require.NoError(t, err)
	err = keeper.IssueDenom(ctx, msgDenom2)
	require.NoError(t, err)
	tokenid, err = keeper.MintNFT(ctx, msgMint)
	require.NoError(t, err)

	require.NoError(t, err)
	reqDenomExistByte, err := json.Marshal(nftDenomReq)
	require.NoError(t, err)
	reqDenomNotExistByte, err := json.Marshal(noneNftDenomReq)
	require.NoError(t, err)
	reqDenomAllByte, err := json.Marshal(nftDenomAllReq)
	require.NoError(t, err)
	reqNftByte, err := json.Marshal(nftReq)
	require.NoError(t, err)
	reqNoDenomNftByte, err := json.Marshal(noDenomNftReq)
	require.NoError(t, err)
	reqNoIdNftByte, err := json.Marshal(noIdNftReq)
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
	nftByte, err := json.Marshal(irismodtypes.QueryNFTResponse{
		NFT: &irismodtypes.BaseNFT{
			Id:    types.TokenID(tokenid).ToIris(),
			Name:  msgMint.GetName(),
			URI:   msgMint.GetURI(),
			Data:  msgMint.GetData(),
			Owner: msgMint.GetRecipient(),
		},
	})
	require.NoError(t, err)
	emptyNftByte, err := json.Marshal(irismodtypes.QueryNFTResponse{})
	require.NoError(t, err)

	for _, tc := range []struct {
		desc     string
		request  json.RawMessage
		expected []byte
		err      error
	}{
		{
			desc:     "succeed to get denom",
			request:  json.RawMessage(reqDenomExistByte),
			expected: denomByte,
		},
		{
			desc:     "faild to get denom by invalid denomid",
			request:  json.RawMessage(reqDenomNotExistByte),
			expected: emptyDenomByte,
		},
		{
			desc:     "succeed to get denom all",
			request:  json.RawMessage(reqDenomAllByte),
			expected: denomAllByte,
		},
		{
			desc:     "succeed to get nft",
			request:  json.RawMessage(reqNftByte),
			expected: nftByte,
		},
		{
			desc:     "faild to get nft when invalid denomid",
			request:  json.RawMessage(reqNoDenomNftByte),
			expected: emptyNftByte,
		},
		{
			desc:     "faild to get nft when invalid tokneid",
			request:  json.RawMessage(reqNoIdNftByte),
			expected: emptyNftByte,
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
