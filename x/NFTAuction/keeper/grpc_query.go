package keeper

import (
	"github.com/cosmos-gaminghub/nibiru/x/NFTAuction/types"
)

var _ types.QueryServer = Keeper{}
