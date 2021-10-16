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

		_denomMeg = types.DenomMessage{
			DenomId: "denomid",
			Name:    "name",
			Schema:  "schema",
		}
		_denomIssueMsg   = types.DenomIssueMessage{IssueDenom: &_denomMeg}
		nftDenomIssueMsg = types.NftDenomIssueMessage{Nft: &_denomIssueMsg}
	)

	msgIssueDenomByte, err := json.Marshal(nftDenomIssueMsg)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc     string
		sender   sdk.AccAddress
		msg      json.RawMessage
		expected []sdk.Msg
		err      error
	}{
		{
			desc:   "nft",
			sender: sender,
			msg:    json.RawMessage(msgIssueDenomByte),
			expected: []sdk.Msg{
				&msgIssueDenom,
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
