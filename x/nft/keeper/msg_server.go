package keeper

import (
	"context"

	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	irismodkeeper "github.com/irisnet/irismod/modules/nft/keeper"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
)

type msgServer struct {
	Keeper
	irisServer irismodtypes.MsgServer
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		Keeper:     keeper,
		irisServer: irismodkeeper.NewMsgServerImpl(keeper.irisKeeper),
	}
}

var _ types.MsgServer = msgServer{}

// IssueDenom issue a new denom.
func (m msgServer) IssueDenom(goCtx context.Context, msg *types.MsgIssueDenom) (*types.MsgIssueDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	irisMsg, err := m.Keeper.toIrisMsgIssueDenom(ctx, msg)
	if err != nil {
		return nil, err
	}

	_, err = m.irisServer.IssueDenom(goCtx, irisMsg)
	if err != nil {
		return nil, err
	}

	denomID, _ := types.ToDenomID(irisMsg.Id)

	return &types.MsgIssueDenomResponse{
		Id: denomID.Uint64(),
	}, nil
}

func (m msgServer) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	irisMsg, err := m.Keeper.toIrisMsgMintNFT(ctx, msg)
	if err != nil {
		return nil, err
	}

	_, err = m.irisServer.MintNFT(goCtx, irisMsg)
	if err != nil {
		return nil, err
	}

	tokenID, _ := types.ToTokenID(irisMsg.Id)

	return &types.MsgMintNFTResponse{
		Id: tokenID.Uint64(),
	}, nil
}

func (m msgServer) EditNFT(goCtx context.Context, msg *types.MsgEditNFT) (*types.MsgEditNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	irisMsg, err := m.Keeper.toIrisMsgEditNFT(ctx, msg)
	if err != nil {
		return nil, err
	}

	_, err = m.irisServer.EditNFT(goCtx, irisMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgEditNFTResponse{}, nil
}

func (m msgServer) TransferNFT(goCtx context.Context, msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	irisMsg, err := m.Keeper.toIrisMsgTransferNFT(ctx, msg)
	if err != nil {
		return nil, err
	}

	_, err = m.irisServer.TransferNFT(goCtx, irisMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgTransferNFTResponse{}, nil
}

func (m msgServer) BurnNFT(goCtx context.Context, msg *types.MsgBurnNFT) (*types.MsgBurnNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	irisMsg, err := m.Keeper.toIrisMsgBurnNFT(ctx, msg)
	if err != nil {
		return nil, err
	}

	_, err = m.irisServer.BurnNFT(goCtx, irisMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgBurnNFTResponse{}, nil
}
