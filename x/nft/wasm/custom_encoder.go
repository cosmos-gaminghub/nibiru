package wasm

import (
	"encoding/json"
	"errors"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/davecgh/go-spew/spew"
)

type CustomMessageEncoder struct {
	Nft    func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error)
	Custom func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error)
}

func DefaultCustomEncoder() CustomMessageEncoder {
	return CustomMessageEncoder{
		Nft:    EncodeNftMsg,
		Custom: wasmkeeper.NoCustomMsg,
	}
}

func (e CustomMessageEncoder) Encode(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
	spew.Dump(msg)
	if msgs, err := e.Nft(sender, msg); err == nil {
		return msgs, nil
	} else if !errors.Is(err, types.ErrUnexpectedMsg) {
		return nil, err
	}

	return e.Custom(sender, msg)
}

func EncodeNftMsg(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
	var denom types.MsgIssueDenom
	if err := denom.Unmarshal(msg); err == nil {
		denom.Sender = sender.String()
		return []sdk.Msg{&denom}, nil
	}

	return nil, types.ErrUnexpectedMsg
}
