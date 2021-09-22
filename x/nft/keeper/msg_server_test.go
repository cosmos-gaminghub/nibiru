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

func TestMsgIssueDenomMintEditTransferBurnNFT(t *testing.T) {
	var (
		srv, ctx  = setupMsgServer(t)
		denomID   = "denomID"
		tokenID   = uint64(types.MIN_TOKEN_ID)
		owner     = testutil.CreateTestAddrs(1)[0].String()
		receipent = testutil.CreateTestAddrs(2)[1].String()
		err       error
	)

	//------test IssueDenom()-------------
	_, err = srv.IssueDenom(ctx, &types.MsgIssueDenom{
		DenomId: denomID,
		Name:    "name",
		Schema:  "schema",
		Sender:  owner,
	})
	require.NoError(t, err)

	//------test MintNFT()-------------
	resp, err := srv.MintNFT(ctx, &types.MsgMintNFT{
		DenomId:   denomID,
		Name:      "name",
		URI:       "uri",
		Data:      "data",
		Sender:    owner,
		Recipient: owner,
	})
	require.NoError(t, err)
	assert.Equal(t, tokenID, resp.Id)

	//------test EditNFT()-------------
	_, err = srv.EditNFT(ctx, &types.MsgEditNFT{
		DenomId: denomID,
		Id:      tokenID,
		Name:    "name2",
		Data:    "data2",
		Sender:  owner,
	})
	require.NoError(t, err)

	//------test TransferNFT()-------------
	_, err = srv.TransferNFT(ctx, &types.MsgTransferNFT{
		DenomId:   denomID,
		Id:        tokenID,
		Sender:    owner,
		Recipient: receipent,
	})
	require.NoError(t, err)

	//------test BurnNFT()-------------
	_, err = srv.BurnNFT(ctx, &types.MsgBurnNFT{
		DenomId: denomID,
		Id:      tokenID,
		Sender:  receipent,
	})
	require.NoError(t, err)
}
