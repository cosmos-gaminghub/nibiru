package signal_test

import (
	"testing"

	keepertest "github.com/cosmos-gaminghub/nibiru/testutil/keeper"
	"github.com/cosmos-gaminghub/nibiru/x/signal"
	"github.com/cosmos-gaminghub/nibiru/x/signal/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SignalKeeper(t)
	signal.InitGenesis(ctx, *k, genesisState)
	got := signal.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// this line is used by starport scaffolding # genesis/test/assert
}
