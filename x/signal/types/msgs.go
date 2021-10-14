package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constant used to indicate that some field should not be updated
const (
	TypeMsgCreateSignal = "create_signal"
)

var (
	_ sdk.Msg = &MsgCreateSignal{}
)

// NewMsgIssueDenom is a constructor function for MsgSetName
func NewMsgSignal(action string, sender string) *MsgCreateSignal {
	return &MsgCreateSignal{
		Action: action,
		Sender: sender,
	}
}

// Route Implements Msg
func (msg MsgCreateSignal) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgCreateSignal) Type() string { return TypeMsgCreateSignal }

// ValidateBasic Implements Msg.
func (msg MsgCreateSignal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateSignal) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgCreateSignal) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
