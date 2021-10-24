package wasm

import (
	"encoding/json"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos-gaminghub/nibiru/testutil"
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCustomMessageEncoder(t *testing.T) {
	var (
		encoder       = DefaultCustomEncoder()
		sender        = testutil.CreateTestAddrs(1)[0]
		msgIssueDenom = types.MsgIssueDenom{
			DenomId: "denomid",
			Name:    "name",
			Schema:  "schema",
			Sender:  sender.String(),
		}
		msgMintNFT = types.MsgMintNFT{
			DenomId:   "denomid",
			Name:      "name",
			URI:       "uri",
			Data:      "data",
			Sender:    sender.String(),
			Recipient: "recipient",
		}
		msgEditNFT = types.MsgEditNFT{
			DenomId: "denomid",
			Id:      1,
			Name:    "name",
			Data:    "data",
			Sender:  sender.String(),
		}
		msgTransferNFT = types.MsgTransferNFT{
			DenomId:   "denomid",
			Id:        1,
			Recipient: "recipient",
			Sender:    sender.String(),
		}
		msgBurnNFT = types.MsgBurnNFT{
			DenomId: "denomid",
			Id:      1,
			Sender:  sender.String(),
		}

		_denomMeg = types.DenomIssueMessage{
			DenomId: "denomid",
			Name:    "name",
			Schema:  "schema",
		}
		_denomIssueMsg   = types.NftDenomIssueMessage{IssueDenom: &_denomMeg}
		nftDenomIssueMsg = types.GameNftDenomIssueMessage{Nft: &_denomIssueMsg}
		_mintMeg         = types.MintMessage{
			DenomId:   "denomid",
			Name:      "name",
			URI:       "uri",
			Data:      "data",
			Recipient: "recipient",
		}
		_nftMintMsg = types.NftMintMessage{MintNft: &_mintMeg}
		nftMintMsg  = types.GameNftMintMessage{Nft: &_nftMintMsg}
		_editMeg    = types.EditMessage{
			DenomId: "denomid",
			Id:      1,
			Name:    "name",
			Data:    "data",
		}
		_nftEditMsg  = types.NftEditMessage{EditNft: &_editMeg}
		nftEditMsg   = types.GameNftEditMessage{Nft: &_nftEditMsg}
		_transferMeg = types.TransferMessage{
			DenomId:   "denomid",
			Id:        1,
			Recipient: "recipient",
		}
		_nftTransferMsg = types.NftTransferMessage{TransferNft: &_transferMeg}
		nftTransferMsg  = types.GameNftTransferMessage{Nft: &_nftTransferMsg}
		_burnMeg        = types.BurnMessage{
			DenomId: "denomid",
			Id:      1,
		}
		_nftBurnMsg = types.NftBurnMessage{BurnNft: &_burnMeg}
		nftBurnMsg  = types.GameNftBurnMessage{Nft: &_nftBurnMsg}
	)

	msgIssueDenomByte, err := json.Marshal(nftDenomIssueMsg)
	require.NoError(t, err)
	msgMintByte, err := json.Marshal(nftMintMsg)
	require.NoError(t, err)
	msgEditByte, err := json.Marshal(nftEditMsg)
	require.NoError(t, err)
	msgTransferByte, err := json.Marshal(nftTransferMsg)
	require.NoError(t, err)
	msgBurnByte, err := json.Marshal(nftBurnMsg)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc     string
		sender   sdk.AccAddress
		msg      json.RawMessage
		expected []sdk.Msg
		err      error
	}{
		{
			desc:   "issue denom",
			sender: sender,
			msg:    json.RawMessage(msgIssueDenomByte),
			expected: []sdk.Msg{
				&msgIssueDenom,
			},
		},
		{
			desc:   "mint nft",
			sender: sender,
			msg:    json.RawMessage(msgMintByte),
			expected: []sdk.Msg{
				&msgMintNFT,
			},
		},
		{
			desc:   "edit nft",
			sender: sender,
			msg:    json.RawMessage(msgEditByte),
			expected: []sdk.Msg{
				&msgEditNFT,
			},
		},
		{
			desc:   "transfer nft",
			sender: sender,
			msg:    json.RawMessage(msgTransferByte),
			expected: []sdk.Msg{
				&msgTransferNFT,
			},
		},
		{
			desc:   "burn nft",
			sender: sender,
			msg:    json.RawMessage(msgBurnByte),
			expected: []sdk.Msg{
				&msgBurnNFT,
			},
		},
		{
			desc:   "custom",
			sender: sender,
			msg:    json.RawMessage([]byte("custom")),
			err:    wasmtypes.ErrUnknownMsg,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			msgs, err := encoder.Encode(tc.sender, tc.msg)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, msgs)
			}
		})
	}
}
