package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

type Snapshot struct {
	TotalGameAmount          sdk.Int `json:"total_game_amount"`
	TotalNibiruAirdropAmount sdk.Int `json:"total_nibiru_amount"`
	NumberAccounts           uint64  `json:"num_accounts"`

	Accounts map[string]SnapshotAccount `json:"accounts"`
}

// SnapshotAccount provide fields of snapshot per account
type SnapshotAccount struct {
	GameAddress string `json:"game_address"` // game Balance = GametakedBalance + GameUnstakedBalance

	GameBalance          sdk.Int `json:"game_balance"`
	GameOwnershipPercent sdk.Dec `json:"game_ownership_percent"`

	GameStakedBalance   sdk.Int `json:"game_staked_balance"`
	GameUnstakedBalance sdk.Int `json:"game_unstaked_balance"` // GameStakedPercent = GameStakedBalance / GameBalance
	GameStakedPercent   sdk.Dec `json:"game_staked_percent"`

	NibiruBalance      sdk.Int `json:"nibiru_balance"`           // OsmoBalance = sqrt( GameBalance ) * (1 + 1.5 * game staked percent)
	NibiruBalanceBase  sdk.Int `json:"nibiru_balance_base"`      // OsmoBalanceBase = sqrt(game balance)
	NibiruBalanceBonus sdk.Int `json:"nibiru_balance_bonus"`     // OsmoBalanceBonus = OsmoBalanceBase * (1.5 * game staked percent)
	NibiruPercent      sdk.Dec `json:"nibiru_ownership_percent"` // OsmoPercent = OsmoNormalizedBalance / TotalOsmoSupply
}

// ExportAirdropSnapshotCmd generates a snapshot.json from a provided cosmos-sdk v0.36 genesis export.
func ExportAirdropSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-airdrop-snapshot [airdrop-to-denom] [input-genesis-file] [output-snapshot-json] --nibiru-supply=[osmos-genesis-supply]",
		Short: "Export a quadratic fairdrop snapshot from a provided cosmos-sdk v0.36 genesis export",
		Long: `Export a quadratic fairdrop snapshot from a provided cosmos-sdk v0.36 genesis export
Sample genesis file:
	https://raw.githubusercontent.com/cephalopodequipment/cosmoshub-3/master/genesis.json
Example:
	nibirud export-airdrop-snapshot game ~/.nibiru/config/genesis.json ../snapshot.json --nibiru-supply=100000000000000
	- Check input genesis:
		file is at ~/.nibirud/config/genesis.json
	- Snapshot
		file is at "../snapshot.json"
`,
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			denom := args[0]
			genesisFile := args[1]
			snapshotOutput := args[2]

			// Read genesis file
			appState, _, err := genutiltypes.GenesisStateFromGenFile(genesisFile)
			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			stakingGenState := stakingtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			// authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			fmt.Println(bankGenState.Balances)

			// Produce the map of address to total game balance, both staked and unstaked
			snapshotAccs := make(map[string]SnapshotAccount)

			totalGameBalance := sdk.NewInt(0)
			for _, account := range bankGenState.Balances {

				balance := account.Coins.AmountOf(denom)
				totalGameBalance = totalGameBalance.Add(balance)

				snapshotAccs[account.Address] = SnapshotAccount{
					GameAddress:         account.Address,
					GameBalance:         balance,
					GameStakedBalance:   balance,
					GameUnstakedBalance: sdk.ZeroInt(),
				}
			}

			for _, unbonding := range stakingGenState.UnbondingDelegations {
				address := unbonding.DelegatorAddress
				acc, ok := snapshotAccs[address]
				if !ok {
					panic("no account found for unbonding")
				}

				unbondingGames := sdk.NewInt(0)
				for _, entry := range unbonding.Entries {
					unbondingGames = unbondingGames.Add(entry.Balance)
				}

				acc.GameBalance = acc.GameBalance.Add(unbondingGames)
				acc.GameUnstakedBalance = acc.GameUnstakedBalance.Add(unbondingGames)

				snapshotAccs[address] = acc
			}

			// Make a map from validator operator address to the v036 validator type
			validators := make(map[string]stakingtypes.Validator)
			for _, validator := range stakingGenState.Validators {
				validators[validator.OperatorAddress] = validator
			}

			for _, delegation := range stakingGenState.Delegations {
				address := delegation.DelegatorAddress

				acc, ok := snapshotAccs[address]
				if !ok {
					panic("no account found for delegation")
				}

				val := validators[delegation.ValidatorAddress]
				stakedGames := delegation.Shares.MulInt(val.Tokens).Quo(val.DelegatorShares).RoundInt()

				acc.GameBalance = acc.GameBalance.Add(stakedGames)
				acc.GameStakedBalance = acc.GameStakedBalance.Add(stakedGames)

				snapshotAccs[address] = acc
			}

			totalNibiruBalance := sdk.NewInt(0)
			onePointFive := sdk.MustNewDecFromStr("1.5")

			for address, acc := range snapshotAccs {
				allGames := acc.GameBalance.ToDec()

				acc.GameOwnershipPercent = allGames.QuoInt(totalGameBalance)

				if allGames.IsZero() {
					acc.GameStakedPercent = sdk.ZeroDec()
					acc.NibiruBalanceBase = sdk.ZeroInt()
					acc.NibiruBalanceBonus = sdk.ZeroInt()
					acc.NibiruBalance = sdk.ZeroInt()
					snapshotAccs[address] = acc
					continue
				}

				stakedGame := acc.GameStakedBalance.ToDec()
				stakedPercent := stakedGame.Quo(allGames)
				acc.GameStakedPercent = stakedPercent

				baseNibiru, err := allGames.ApproxSqrt()
				if err != nil {
					panic(fmt.Sprintf("failed to root game balance: %s", err))
				}
				acc.NibiruBalanceBase = baseNibiru.RoundInt()

				bonusNibiru := baseNibiru.Mul(onePointFive).Mul(stakedPercent)
				acc.NibiruBalanceBonus = bonusNibiru.RoundInt()

				allOsmo := baseNibiru.Add(bonusNibiru)
				// nibiruBalance = sqrt( all games) * (1 + 1.5) * (staked game percent) =
				acc.NibiruBalance = allOsmo.RoundInt()

				if allGames.LTE(sdk.NewDec(1000000)) {
					acc.NibiruBalanceBase = sdk.ZeroInt()
					acc.NibiruBalanceBonus = sdk.ZeroInt()
					acc.NibiruBalance = sdk.ZeroInt()
				}

				totalNibiruBalance = totalNibiruBalance.Add(acc.NibiruBalance)

				snapshotAccs[address] = acc
			}

			// iterate to find Osmo ownership percentage per account
			for address, acc := range snapshotAccs {
				acc.NibiruPercent = acc.NibiruBalance.ToDec().Quo(totalNibiruBalance.ToDec())
				snapshotAccs[address] = acc
			}

			snapshot := Snapshot{
				TotalGameAmount:          totalGameBalance,
				TotalNibiruAirdropAmount: totalNibiruBalance,
				NumberAccounts:           uint64(len(snapshotAccs)),
				Accounts:                 snapshotAccs,
			}

			fmt.Printf("# accounts: %d\n", len(snapshotAccs))
			fmt.Printf("gameTotalSupply: %s\n", totalGameBalance.String())
			fmt.Printf("nibiruTotalSupply: %s\n", totalNibiruBalance.String())

			// export snapshot json
			snapshotJSON, err := json.MarshalIndent(snapshot, "", "    ")
			if err != nil {
				return fmt.Errorf("failed to marshal snapshot: %w", err)
			}

			err = ioutil.WriteFile(snapshotOutput, snapshotJSON, 0644)
			return err
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
