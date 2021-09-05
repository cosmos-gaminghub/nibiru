package keeper

import (
	"github.com/cosmos-gaminghub/nibiru/x/auction/types"
)

var _ types.QueryServer = Keeper{}
