package keeper

import (
	"context"

	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) NFT(c context.Context, request *types.QueryNFTRequest) (*types.QueryNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	nft, err := k.GetNFT(ctx, request.DenomId, request.Id)
	if err != nil {
		return nil, sdkerrors.Wrapf(irismodtypes.ErrUnknownNFT, "invalid NFT %d from collection %s", request.Id, request.DenomId)
	}

	return &types.QueryNFTResponse{NFT: &types.BaseNFT{
		Id:    nft.GetID(),
		Name:  nft.GetName(),
		URI:   nft.GetURI(),
		Data:  nft.GetData(),
		Owner: nft.GetOwner().String(),
	}}, nil
}
