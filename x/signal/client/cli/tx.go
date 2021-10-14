package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos-gaminghub/nibiru/x/signal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		EmitEvent(),
	)
	return cmd
}

func EmitEvent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "emit [action] [address]",
		Short: "Emit signal event",
		Long: `Emit signal event
Example:
	nibirud tx signal emit "update signal" "game1..."
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			action := args[0]
			address := args[1]

			key, err := sdk.AccAddressFromBech32(address)
			if err != nil {
				return err
			}

			e := sdk.NewEvent("signal", sdk.NewAttribute("action", action))
			e = e.AppendAttributes(sdk.NewAttribute("sender", key.String()))

			em := sdk.NewEventManager()
			em.EmitEvent(e)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
