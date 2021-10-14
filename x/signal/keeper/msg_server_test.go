package keeper_test

import (
	"context"
	"testing"

	keepertest "github.com/cosmos-gaminghub/nibiru/testutil/keeper"
	"github.com/cosmos-gaminghub/nibiru/x/signal/keeper"
	"github.com/cosmos-gaminghub/nibiru/x/signal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.SignalKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
