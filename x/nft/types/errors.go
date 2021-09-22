package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/nft module sentinel errors
var (
	ErrRestricted = sdkerrors.Register(ModuleName, 1, "restricted")
	// this line is used by starport scaffolding # ibc/errors
)
