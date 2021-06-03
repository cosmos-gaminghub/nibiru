package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/spf13/cast"

	"github.com/cosmos-gaminghub/nibiru/x/NFTMarket/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func CmdCreateItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-item [Denom] [NFTID] [Price] [Fee] [Detail] [InSale]",
		Short: "Create a new item",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsDenom, err := cast.ToStringE(args[0])
			if err != nil {
				return err
			}
			argsNFTID, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}
			argsPrice, err := cast.ToStringE(args[2])
			if err != nil {
				return err
			}
			argsFee, err := cast.ToStringE(args[3])
			if err != nil {
				return err
			}
			argsDetail, err := cast.ToStringE(args[4])
			if err != nil {
				return err
			}
			argsInSale, err := cast.ToBoolE(args[5])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateItem(clientCtx.GetFromAddress().String(), argsDenom, argsNFTID, argsPrice, argsFee, argsDetail, argsInSale)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-item [id] [Denom] [NFTID] [Price] [Fee] [Detail] [InSale]",
		Short: "Update a item",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argsDenom, err := cast.ToStringE(args[1])
			if err != nil {
				return err
			}

			argsNFTID, err := cast.ToStringE(args[2])
			if err != nil {
				return err
			}

			argsPrice, err := cast.ToStringE(args[3])
			if err != nil {
				return err
			}

			argsFee, err := cast.ToStringE(args[4])
			if err != nil {
				return err
			}

			argsDetail, err := cast.ToStringE(args[5])
			if err != nil {
				return err
			}

			argsInSale, err := cast.ToBoolE(args[6])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateItem(clientCtx.GetFromAddress().String(), id, argsDenom, argsNFTID, argsPrice, argsFee, argsDetail, argsInSale)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteItem() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-item [id]",
		Short: "Delete a item by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteItem(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
