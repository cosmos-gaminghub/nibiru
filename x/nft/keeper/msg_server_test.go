package keeper

import (
	"context"
	"testing"

	"github.com/cosmos-gaminghub/nibiru/testutil"
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	keeper, ctx := setupKeeper(t)
	return NewMsgServerImpl(*keeper), sdk.WrapSDKContext(ctx)
}

func TestMsgIssueDenom(t *testing.T) {
	var (
		srv, ctx = setupMsgServer(t)
	)
	_, err := srv.IssueDenom(ctx, &types.MsgIssueDenom{"denomid", "name", "schema", testutil.CreateTestAddrs(1)[0].String()})
	require.NoError(t, err)
}

func TestMsgMintNFT(t *testing.T) {
	var (
		srv, ctx   = setupMsgServer(t)
		denomID    = "denomID"
		owner      = testutil.CreateTestAddrs(1)[0].String()
		expectedId = uint64(types.MIN_TOKEN_ID)
	)
	srv.IssueDenom(ctx, &types.MsgIssueDenom{denomID, "name", "schema", owner})
	resp, err := srv.MintNFT(ctx, &types.MsgMintNFT{denomID, "name", "uri", "data", owner, owner})
	require.NoError(t, err)
	assert.Equal(t, expectedId, resp.Id)
}

func TestMsgEditNFT(t *testing.T) {
	var (
		srv, ctx = setupMsgServer(t)
		denomID  = "denomID"
		tokenid  = uint64(types.MIN_TOKEN_ID)
		owner    = testutil.CreateTestAddrs(1)[0].String()
	)
	srv.IssueDenom(ctx, &types.MsgIssueDenom{denomID, "name", "schema", owner})
	srv.MintNFT(ctx, &types.MsgMintNFT{denomID, "name", "uri", "data", owner, owner})
	_, err := srv.EditNFT(ctx, &types.MsgEditNFT{denomID, tokenid, "name2", "data2", owner})
	require.NoError(t, err)
}

func TestMsgTransferNFT(t *testing.T) {
	var (
		srv, ctx  = setupMsgServer(t)
		denomID   = "denomID"
		tokenid   = uint64(types.MIN_TOKEN_ID)
		owner     = testutil.CreateTestAddrs(1)[0].String()
		receipent = testutil.CreateTestAddrs(2)[1].String()
	)
	srv.IssueDenom(ctx, &types.MsgIssueDenom{denomID, "name", "schema", owner})
	srv.MintNFT(ctx, &types.MsgMintNFT{denomID, "name", "uri", "data", owner, owner})
	_, err := srv.TransferNFT(ctx, &types.MsgTransferNFT{denomID, tokenid, owner, receipent})
	require.NoError(t, err)
}

func TestMsgBurnNFT(t *testing.T) {
	var (
		srv, ctx = setupMsgServer(t)
		denomID  = "denomID"
		tokenid  = uint64(types.MIN_TOKEN_ID)
		owner    = testutil.CreateTestAddrs(1)[0].String()
	)
	srv.IssueDenom(ctx, &types.MsgIssueDenom{denomID, "name", "schema", owner})
	srv.MintNFT(ctx, &types.MsgMintNFT{denomID, "name", "uri", "data", owner, owner})
	_, err := srv.BurnNFT(ctx, &types.MsgBurnNFT{denomID, tokenid, owner})
	require.NoError(t, err)
}
