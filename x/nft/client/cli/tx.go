package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	iriscli "github.com/irisnet/irismod/modules/nft/client/cli"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		iriscli.GetCmdIssueDenom(),
		iriscli.GetCmdMintNFT(),
		iriscli.GetCmdEditNFT(),
		iriscli.GetCmdTransferNFT(),
		iriscli.GetCmdBurnNFT(),
	)

	// this line is used by starport scaffolding # 1

	return cmd
}