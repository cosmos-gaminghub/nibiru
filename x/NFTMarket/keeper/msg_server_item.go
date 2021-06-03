package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos-gaminghub/nibiru/x/NFTMarket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateItem(goCtx context.Context, msg *types.MsgCreateItem) (*types.MsgCreateItemResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var item = types.Item{
		Creator: msg.Creator,
		Denom:   msg.Denom,
		NFTID:   msg.NFTID,
		Price:   msg.Price,
		Fee:     msg.Fee,
		Detail:  msg.Detail,
		InSale:  msg.InSale,
	}

	id := k.AppendItem(
		ctx,
		item,
	)

	return &types.MsgCreateItemResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateItem(goCtx context.Context, msg *types.MsgUpdateItem) (*types.MsgUpdateItemResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var item = types.Item{
		Creator: msg.Creator,
		Id:      msg.Id,
		Denom:   msg.Denom,
		NFTID:   msg.NFTID,
		Price:   msg.Price,
		Fee:     msg.Fee,
		Detail:  msg.Detail,
		InSale:  msg.InSale,
	}

	// Checks that the element exists
	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the the msg sender is the same as the current owner
	if msg.Creator != k.GetItemOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetItem(ctx, item)

	return &types.MsgUpdateItemResponse{}, nil
}

func (k msgServer) DeleteItem(goCtx context.Context, msg *types.MsgDeleteItem) (*types.MsgDeleteItemResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.HasItem(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != k.GetItemOwner(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveItem(ctx, msg.Id)

	return &types.MsgDeleteItemResponse{}, nil
}
