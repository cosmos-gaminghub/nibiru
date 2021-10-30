package keeper

import (
	"github.com/cosmos-gaminghub/nibiru/x/signal/types"
)

var _ types.QueryServer = Keeper{}
