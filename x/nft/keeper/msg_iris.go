package keeper

import (
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
)

func (k Keeper) toIrisMsgIssueDenom(ctx sdk.Context, msg *types.MsgIssueDenom) (*irismodtypes.MsgIssueDenom, error) {
	return irismodtypes.NewMsgIssueDenom(
		msg.DenomId,
		msg.Name,
		msg.Schema,
		msg.Sender,
	), nil
}

func (k Keeper) toIrisMsgMintNFT(ctx sdk.Context, msg *types.MsgMintNFT) (*irismodtypes.MsgMintNFT, error) {
	tokenID, err := k.NewTokenID(ctx, msg.DenomId)
	if err != nil {
		return nil, err
	}

	return irismodtypes.NewMsgMintNFT(
		tokenID.ToIris(),
		msg.DenomId,
		msg.Name,
		msg.URI,
		msg.Data,
		msg.Sender,
		msg.Recipient,
	), nil
}

func (k Keeper) toIrisMsgEditNFT(ctx sdk.Context, msg *types.MsgEditNFT) (*irismodtypes.MsgEditNFT, error) {
	nft, err := k.GetNFT(ctx, msg.DenomId, msg.Id)
	if err != nil {
		return nil, err
	}

	return irismodtypes.NewMsgEditNFT(
		types.TokenID(msg.Id).ToIris(),
		msg.DenomId,
		msg.Name,
		nft.GetURI(),
		msg.Data,
		msg.Sender,
	), nil
}

func (k Keeper) toIrisMsgTransferNFT(ctx sdk.Context, msg *types.MsgTransferNFT) (*irismodtypes.MsgTransferNFT, error) {
	nft, err := k.GetNFT(ctx, msg.DenomId, msg.Id)
	if err != nil {
		return nil, err
	}

	return irismodtypes.NewMsgTransferNFT(
		types.TokenID(msg.Id).ToIris(),
		msg.DenomId,
		nft.GetName(),
		nft.GetURI(),
		nft.GetData(),
		msg.Sender,
		msg.Recipient,
	), nil
}

func (k Keeper) toIrisMsgBurnNFT(ctx sdk.Context, msg *types.MsgBurnNFT) (*irismodtypes.MsgBurnNFT, error) {
	return irismodtypes.NewMsgBurnNFT(
		msg.Sender,
		types.TokenID(msg.Id).ToIris(),
		msg.DenomId,
	), nil
}