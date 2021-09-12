package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// IssueDenom issue a new denom.
func (m msgServer) IssueDenom(goCtx context.Context, msg *types.MsgIssueDenom) (*types.MsgIssueDenomResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	denomID, err := m.Keeper.IssueDenomn(ctx, msg.Name, msg.Schema, sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			irismodtypes.EventTypeIssueDenom,
			sdk.NewAttribute(irismodtypes.AttributeKeyDenomID, fmt.Sprintf("%d", denomID)),
			sdk.NewAttribute(irismodtypes.AttributeKeyDenomName, msg.Name),
			sdk.NewAttribute(irismodtypes.AttributeKeyCreator, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, irismodtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgIssueDenomResponse{
		Id: denomID,
	}, nil
}

func (m msgServer) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	tokenID, err := m.Keeper.MintNFTn(ctx, msg.DenomId, msg.Name, msg.URI, msg.Data, recipient)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			irismodtypes.EventTypeMintNFT,
			sdk.NewAttribute(irismodtypes.AttributeKeyDenomID, fmt.Sprintf("%d", msg.DenomId)),
			sdk.NewAttribute(irismodtypes.AttributeKeyTokenID, fmt.Sprintf("%d", tokenID)),
			sdk.NewAttribute(irismodtypes.AttributeKeyTokenURI, msg.URI),
			sdk.NewAttribute(irismodtypes.AttributeKeyRecipient, msg.Recipient),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, irismodtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgMintNFTResponse{
		Id: tokenID,
	}, nil
}

func (m msgServer) EditNFT(goCtx context.Context, msg *types.MsgEditNFT) (*types.MsgEditNFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.EditNFTn(ctx, msg.DenomId, msg.Id, msg.Name, msg.Data, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			irismodtypes.EventTypeEditNFT,
			sdk.NewAttribute(irismodtypes.AttributeKeyDenomID, fmt.Sprintf("%d", msg.DenomId)),
			sdk.NewAttribute(irismodtypes.AttributeKeyTokenID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(irismodtypes.AttributeKeyOwner, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, irismodtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgEditNFTResponse{}, nil
}

func (m msgServer) TransferNFT(goCtx context.Context, msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	recipient, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.TransferOwnern(ctx, msg.DenomId, msg.Id, sender, recipient); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			irismodtypes.EventTypeTransfer,
			sdk.NewAttribute(irismodtypes.AttributeKeyDenomID, fmt.Sprintf("%d", msg.DenomId)),
			sdk.NewAttribute(irismodtypes.AttributeKeyTokenID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(irismodtypes.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(irismodtypes.AttributeKeyRecipient, msg.Recipient),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, irismodtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgTransferNFTResponse{}, nil
}

func (m msgServer) BurnNFT(goCtx context.Context, msg *types.MsgBurnNFT) (*types.MsgBurnNFTResponse, error) {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.Keeper.BurnNFTn(ctx, msg.DenomId, msg.Id, sender); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			irismodtypes.EventTypeBurnNFT,
			sdk.NewAttribute(irismodtypes.AttributeKeyDenomID, fmt.Sprintf("%d", msg.DenomId)),
			sdk.NewAttribute(irismodtypes.AttributeKeyTokenID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(irismodtypes.AttributeKeyOwner, msg.Sender),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, irismodtypes.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
	})

	return &types.MsgBurnNFTResponse{}, nil
}
