package cli

import (
	"context"
	"fmt"

	"github.com/cosmos-gaminghub/nibiru/x/nft/types"
	"github.com/cosmos/cosmos-sdk/client"
	iriscli "github.com/irisnet/irismod/modules/nft/client/cli"
	irismodtypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group nft queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		iriscli.GetCmdQueryDenom(),
		iriscli.GetCmdQueryDenoms(),
		iriscli.GetCmdQueryCollection(),
		iriscli.GetCmdQuerySupply(),
		iriscli.GetCmdQueryOwner(),
		GetCmdQueryNFT(),
	)

	// this line is used by starport scaffolding # 1

	return cmd
}

func GetCmdQueryNFT() *cobra.Command {
	cmd := iriscli.GetCmdQueryNFT()
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}

		if err := irismodtypes.ValidateDenomID(args[0]); err != nil {
			return err
		}

		tokenID, err := types.ToTokenID(args[1])
		if err != nil {
			return err
		}

		if err := irismodtypes.ValidateTokenID(tokenID.ToIris()); err != nil {
			return err
		}

		queryClient := irismodtypes.NewQueryClient(clientCtx)
		resp, err := queryClient.NFT(context.Background(), &irismodtypes.QueryNFTRequest{
			DenomId: args[0],
			TokenId: tokenID.ToIris(),
		})
		if err != nil {
			return err
		}
		return clientCtx.PrintProto(resp.NFT)
	}

	return cmd
}
