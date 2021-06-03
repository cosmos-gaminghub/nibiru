package keeper

import (
	"github.com/cosmos-gaminghub/nibiru/x/NFTMarket/types"
)

var _ types.QueryServer = Keeper{}
