package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
)

// constant used to indicate that some field should not be updated
const (
	TypeMsgIssueDenom  = "issue_denom"
	TypeMsgTransferNFT = "transfer_nft"
	TypeMsgEditNFT     = "edit_nft"
	TypeMsgMintNFT     = "mint_nft"
	TypeMsgBurnNFT     = "burn_nft"
)

var (
	_ sdk.Msg = &MsgIssueDenom{}
	_ sdk.Msg = &MsgTransferNFT{}
	_ sdk.Msg = &MsgEditNFT{}
	_ sdk.Msg = &MsgMintNFT{}
	_ sdk.Msg = &MsgBurnNFT{}
)

// NewMsgIssueDenom is a constructor function for MsgSetName
func NewMsgIssueDenom(denomId, denomName, schema, sender string) *MsgIssueDenom {
	return &MsgIssueDenom{
		DenomId: denomId,
		Sender:  sender,
		Name:    denomName,
		Schema:  schema,
	}
}

// Route Implements Msg
func (msg MsgIssueDenom) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgIssueDenom) Type() string { return TypeMsgIssueDenom }

// ValidateBasic Implements Msg.
func (msg MsgIssueDenom) ValidateBasic() error {
	if err := irismodtypes.ValidateDenomID(msg.DenomId); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return irismodtypes.ValidateDenomName(msg.Name)
}

// GetSignBytes Implements Msg.
func (msg MsgIssueDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgIssueDenom) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgMintNFT is a constructor function for MsgMintNFT
func NewMsgMintNFT(
	denomID string, tokenName, tokenURI, tokenData, sender, recipient string,
) *MsgMintNFT {
	return &MsgMintNFT{
		DenomId:   denomID,
		Name:      tokenName,
		URI:       tokenURI,
		Data:      tokenData,
		Sender:    sender,
		Recipient: recipient,
	}
}

// Route Implements Msg
func (msg MsgMintNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMintNFT) Type() string { return TypeMsgMintNFT }

// ValidateBasic Implements Msg.
func (msg MsgMintNFT) ValidateBasic() error {
	if err := irismodtypes.ValidateDenomID(msg.DenomId); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid receipt address (%s)", err)
	}
	if err := irismodtypes.ValidateTokenURI(msg.URI); err != nil {
		return err
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMintNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgMintNFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgEditNFT is a constructor function for MsgSetName
func NewMsgEditNFT(
	denomID string, tokenID uint64, tokenName, tokenData, sender string,
) *MsgEditNFT {
	return &MsgEditNFT{
		DenomId: denomID,
		Id:      tokenID,
		Name:    tokenName,
		Data:    tokenData,
		Sender:  sender,
	}
}

// Route Implements Msg
func (msg MsgEditNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgEditNFT) Type() string { return TypeMsgEditNFT }

// ValidateBasic Implements Msg.
func (msg MsgEditNFT) ValidateBasic() error {
	if err := irismodtypes.ValidateDenomID(msg.DenomId); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return irismodtypes.ValidateTokenID(TokenID(msg.Id).ToIris())
}

// GetSignBytes Implements Msg.
func (msg MsgEditNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgEditNFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgTransferNFT is a constructor function for MsgSetName
func NewMsgTransferNFT(
	denomID string, tokenID uint64, sender, recipient string,
) *MsgTransferNFT {
	return &MsgTransferNFT{
		DenomId:   denomID,
		Id:        tokenID,
		Sender:    sender,
		Recipient: recipient,
	}
}

// Route Implements Msg
func (msg MsgTransferNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgTransferNFT) Type() string { return TypeMsgTransferNFT }

// ValidateBasic Implements Msg.
func (msg MsgTransferNFT) ValidateBasic() error {
	if err := irismodtypes.ValidateDenomID(msg.DenomId); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.Recipient); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}
	return irismodtypes.ValidateTokenID(TokenID(msg.Id).ToIris())
}

// GetSignBytes Implements Msg.
func (msg MsgTransferNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgTransferNFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgBurnNFT is a constructor function for MsgBurnNFT
func NewMsgBurnNFT(sender, denomID string, tokenID uint64) *MsgBurnNFT {
	return &MsgBurnNFT{
		Sender:  sender,
		DenomId: denomID,
		Id:      tokenID,
	}
}

// Route Implements Msg
func (msg MsgBurnNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgBurnNFT) Type() string { return TypeMsgBurnNFT }

// ValidateBasic Implements Msg.
func (msg MsgBurnNFT) ValidateBasic() error {
	if err := irismodtypes.ValidateDenomID(msg.DenomId); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return irismodtypes.ValidateTokenID(TokenID(msg.Id).ToIris())
}

// GetSignBytes Implements Msg.
func (msg MsgBurnNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBurnNFT) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}