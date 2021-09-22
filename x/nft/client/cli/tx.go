package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	iriscli "github.com/irisnet/irismod/modules/nft/client/cli"
	"github.com/spf13/cobra"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
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
		GetCmdIssueDenom(),
		GetCmdMintNFT(),
		GetCmdEditNFT(),
		GetCmdTransferNFT(),
		GetCmdBurnNFT(),
	)

	// this line is used by starport scaffolding # 1

	return cmd
}

// GetCmdIssueDenom is the CLI command for an IssueDenom transaction
func GetCmdIssueDenom() *cobra.Command {
	cmd := iriscli.GetCmdIssueDenom()
	cmd.Example = fmt.Sprintf(
		"$ %s tx nft issue <denom-id> "+
			"--name=<denom-name> "+
			"--schema=<schema-content or path to schema.json> ",
		version.AppName,
	)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}

		denomName, err := cmd.Flags().GetString(iriscli.FlagDenomName)
		if err != nil {
			return err
		}
		schema, err := cmd.Flags().GetString(iriscli.FlagSchema)
		if err != nil {
			return err
		}
		optionsContent, err := ioutil.ReadFile(schema)
		if err == nil {
			schema = string(optionsContent)
		}

		msg := types.NewMsgIssueDenom(
			args[0],
			denomName,
			schema,
			clientCtx.GetFromAddress().String(),
		)
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
		return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
	}

	return cmd
}

// GetCmdMintNFT is the CLI command for a MintNFT transaction
func GetCmdMintNFT() *cobra.Command {
	cmd := iriscli.GetCmdMintNFT()
	cmd.Use = "mint [denom-id]"
	cmd.Example = fmt.Sprintf(
		"$ %s tx nft mint <denom-id> "+
			"--uri=<uri> "+
			"--recipient=<recipient> ",
		version.AppName,
	)
	cmd.Args = cobra.ExactArgs(1)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}

		var sender = clientCtx.GetFromAddress().String()

		recipient, err := cmd.Flags().GetString(iriscli.FlagRecipient)
		if err != nil {
			return err
		}

		recipientStr := strings.TrimSpace(recipient)
		if len(recipientStr) > 0 {
			if _, err = sdk.AccAddressFromBech32(recipientStr); err != nil {
				return err
			}
		} else {
			recipient = sender
		}

		tokenName, err := cmd.Flags().GetString(iriscli.FlagTokenName)
		if err != nil {
			return err
		}
		tokenURI, err := cmd.Flags().GetString(iriscli.FlagTokenURI)
		if err != nil {
			return err
		}
		tokenData, err := cmd.Flags().GetString(iriscli.FlagTokenData)
		if err != nil {
			return err
		}

		msg := types.NewMsgMintNFT(
			args[0],
			tokenName,
			tokenURI,
			tokenData,
			sender,
			recipient,
		)
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
		return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
	}

	return cmd
}

// GetCmdEditNFT is the CLI command for sending an MsgEditNFT transaction
func GetCmdEditNFT() *cobra.Command {
	cmd := iriscli.GetCmdEditNFT()
	cmd.Example = fmt.Sprintf(
		"$ %s tx nft edit <denom-id> <token-id> "+
			"--name=<token-name> "+
			"--data=<token-data> ",
		version.AppName,
	)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}

		tokenName, err := cmd.Flags().GetString(iriscli.FlagTokenName)
		if err != nil {
			return err
		}
		tokenData, err := cmd.Flags().GetString(iriscli.FlagTokenData)
		if err != nil {
			return err
		}
		tokenID, err := types.ToTokenID(args[1])
		if err != nil {
			return err
		}

		msg := types.NewMsgEditNFT(
			args[0],
			tokenID.Uint64(),
			tokenName,
			tokenData,
			clientCtx.GetFromAddress().String(),
		)
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
		return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
	}

	return cmd
}

// GetCmdTransferNFT is the CLI command for sending a TransferNFT transaction
func GetCmdTransferNFT() *cobra.Command {
	cmd := iriscli.GetCmdTransferNFT()
	cmd.Example = fmt.Sprintf(
		"$ %s tx nft transfer <recipient> <denom-id> <token-id>",
		version.AppName,
	)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}

		if _, err := sdk.AccAddressFromBech32(args[0]); err != nil {
			return err
		}
		tokenID, err := types.ToTokenID(args[2])
		if err != nil {
			return err
		}

		msg := types.NewMsgTransferNFT(
			args[1],
			tokenID.Uint64(),
			clientCtx.GetFromAddress().String(),
			args[0],
		)
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
		return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
	}

	return cmd
}

// GetCmdBurnNFT is the CLI command for sending a BurnNFT transaction
func GetCmdBurnNFT() *cobra.Command {
	cmd := iriscli.GetCmdBurnNFT()
	cmd.Args = cobra.ExactArgs(2)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}

		tokenID, err := types.ToTokenID(args[1])
		if err != nil {
			return err
		}

		msg := types.NewMsgBurnNFT(
			clientCtx.GetFromAddress().String(),
			args[0],
			tokenID.Uint64(),
		)
		if err := msg.ValidateBasic(); err != nil {
			return err
		}
		return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
	}

	return cmd
}
