package keeper

import (
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
)

var _ types.QueryServer = Keeper{}
