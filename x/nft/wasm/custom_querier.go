package wasm

import (
	"encoding/json"
	"errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	nftmodulekeeper "github.com/cosmos-gaminghub/nibiru/x/nft/keeper"
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/davecgh/go-spew/spew"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
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
		spew.Dump("CustomQuerier pass!", request)

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
		var denomReq irismodtypes.QueryDenomRequest
		if err := denomReq.Unmarshal(request); err == nil {
			denom, err := k.IrisKeeper().GetDenom(ctx, denomReq.DenomId)
			if err != nil {
				return nil, err
			}
			return json.Marshal(irismodtypes.QueryDenomResponse{
				Denom: &denom,
			})
		}

		return nil, types.ErrUnexpectedReq
	}
}
