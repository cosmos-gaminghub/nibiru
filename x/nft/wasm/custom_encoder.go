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
	spew.Dump("Encoder pass!", msg)
	if msgs, err := e.Nft(sender, msg); err == nil {
		return msgs, nil
	} else if !errors.Is(err, types.ErrUnexpectedMsg) {
		return nil, err
	}

	return e.Custom(sender, msg)
}

func EncodeNftMsg(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error) {
	var wasmIssueMsg types.GameNftDenomIssueMessage
	if err := json.Unmarshal(msg, &wasmIssueMsg); err == nil && wasmIssueMsg.Nft.IssueDenom != nil {
		return []sdk.Msg{&types.MsgIssueDenom{
			DenomId: wasmIssueMsg.Nft.IssueDenom.DenomId,
			Name:    wasmIssueMsg.Nft.IssueDenom.Name,
			Schema:  wasmIssueMsg.Nft.IssueDenom.Schema,
			Sender:  sender.String(),
		}}, nil
	}

	var wasmMintMsg types.GameNftMintMessage
	if err := json.Unmarshal(msg, &wasmMintMsg); err == nil && wasmMintMsg.Nft.MintNft != nil {
		return []sdk.Msg{&types.MsgMintNFT{
			DenomId:   wasmMintMsg.Nft.MintNft.DenomId,
			Name:      wasmMintMsg.Nft.MintNft.Name,
			URI:       wasmMintMsg.Nft.MintNft.URI,
			Data:      wasmMintMsg.Nft.MintNft.Data,
			Sender:    sender.String(),
			Recipient: wasmMintMsg.Nft.MintNft.Recipient,
		}}, nil
	}

	return nil, types.ErrUnexpectedMsg
}
