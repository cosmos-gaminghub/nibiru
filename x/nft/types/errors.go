package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/nft module sentinel errors
var (
	ErrRestricted    = sdkerrors.Register(ModuleName, 101, "restricted")
	ErrUnexpectedMsg = sdkerrors.Register(ModuleName, 102, "unexpected msg")
	ErrUnexpectedReq = sdkerrors.Register(ModuleName, 103, "unexpected request")
	// this line is used by starport scaffolding # ibc/errors
)
