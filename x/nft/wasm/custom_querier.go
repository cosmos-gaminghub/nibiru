package wasm

import (
	"encoding/json"
	"errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	nftmodulekeeper "github.com/cosmos-gaminghub/nibiru/x/nft/keeper"
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
	// "github.com/davecgh/go-spew/spew"
)

type CustomQuerier struct {
	Nft    wasmkeeper.CustomQuerier
	Custom wasmkeeper.CustomQuerier
}

func DefaultCustomQuerier(k *nftmodulekeeper.Keeper) CustomQuerier {
	return CustomQuerier{
		Nft:    NftQuerier(k),
		Custom: wasmkeeper.NoCustomQuerier,
	}
}

func (q CustomQuerier) Querier() wasmkeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		if res, err := q.Nft(ctx, request); err == nil {
			return res, nil
		} else if !errors.Is(err, types.ErrUnexpectedReq) {
			return nil, err
		}

		return q.Custom(ctx, request)
	}
}

func NftQuerier(k *nftmodulekeeper.Keeper) wasmkeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		// query denom
		var denomReq types.WasmGameDenomReq
		if err := json.Unmarshal(request, &denomReq); err == nil && denomReq.Nft.Denom != nil {
			denom, err := k.IrisKeeper().GetDenom(ctx, denomReq.Nft.Denom.DenomId)
			if err != nil {
				// return err when an err other than "no denom exist" happen
				if !errors.Is(err, irismodtypes.ErrInvalidDenom) {
					return nil, err
				}
				return json.Marshal(irismodtypes.QueryDenomResponse{})
			}
			return json.Marshal(irismodtypes.QueryDenomResponse{
				Denom: &denom,
			})
		}

		// query all denom
		var denomAllReq types.WasmGameDenomAllReq
		if err := json.Unmarshal(request, &denomAllReq); err == nil && denomAllReq.Nft.DenomAll != nil {
			denoms := k.IrisKeeper().GetDenoms(ctx)
			return json.Marshal(irismodtypes.QueryDenomsResponse{
				Denoms: denoms,
			})
		}

		// query nft
		var nftReq types.WasmGameNftReq
		if err := json.Unmarshal(request, &nftReq); err == nil && nftReq.Nft.Nft != nil {
			nft, err := k.GetNFT(ctx, nftReq.Nft.Nft.DenomId, nftReq.Nft.Nft.Id)
			if err != nil {
				// return err when an err other than "unknown nft collection" happen
				if !errors.Is(err, irismodtypes.ErrUnknownCollection) {
					return nil, err
				}
				return json.Marshal(irismodtypes.QueryNFTResponse{})
			}
			return json.Marshal(irismodtypes.QueryNFTResponse{
				NFT: &irismodtypes.BaseNFT{
					Id:    nft.GetID(),
					Name:  nft.GetName(),
					URI:   nft.GetURI(),
					Data:  nft.GetData(),
					Owner: nft.GetOwner().String(),
				},
			})
		}

		return nil, types.ErrUnexpectedReq
	}
}
