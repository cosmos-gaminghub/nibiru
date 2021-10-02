package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/spf13/cobra"
)

const (
	MaxCap                 = 50000
	TotalGameAirdropAmount = 1000000 // 0.5% * 200000000
)

type Snapshot struct {
	TotalAtomAmount        sdk.Int `json:"total_atom_amount"`
	TotalGameAirdropAmount sdk.Int `json:"total_game_amount"`
	NumberAccounts         uint64  `json:"num_accounts"`
	Denominator            sdk.Int `json:"denominator"`

	Accounts map[string]SnapshotAccount `json:"accounts"`
}

// SnapshotAccount provide fields of snapshot per account
type SnapshotAccount struct {
	AtomAddress string `json:"atom_address"` // Atom Balance = AtomStakedBalance + AtomUnstakedBalance

	AtomBalance          sdk.Int `json:"atom_balance"`
	AtomOwnershipPercent sdk.Dec `json:"atom_ownership_percent"`

	AtomStakedBalance   sdk.Int `json:"atom_staked_balance"`
	AtomUnstakedBalance sdk.Int `json:"atom_unstaked_balance"` // AtomStakedPercent = AtomStakedBalance / AtomBalance
	AtomStakedPercent   sdk.Dec `json:"atom_staked_percent"`

	GameBalance sdk.Int `json:"game_balance"`
	Denominator sdk.Int `json:"denominator"`
}

type Account struct {
	Address       string `json:"address,omitempty"`
	AccountNumber uint64 `json:"account_number,omitempty"`
	Sequence      uint64 `json:"sequence,omitempty"`
}

// ExportAirdropSnapshotCmd generates a snapshot.json from a provided cosmos-sdk v0.36 genesis export.
func ExportAirdropSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-airdrop-snapshot [airdrop-to-denom] [input-genesis-file] [output-snapshot-json] --nibiru-supply=[nibiru-genesis-supply]",
		Short: "Export a quadratic fairdrop snapshot from a provided cosmos-sdk v0.42 genesis export",
		Long: `Export a quadratic fairdrop snapshot from a provided cosmos-sdk v0.42 genesis export
Example:
	nibirud export-airdrop-snapshot uatom ~/.nibiru/config/genesis.json ../snapshot.json --nibiru-supply=100000000000000
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
			appState, _, _ := genutiltypes.GenesisStateFromGenFile(genesisFile)
			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			stakingGenState := stakingtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			snapshotAccs := make(map[string]SnapshotAccount, len(authGenState.GetAccounts()))
			for _, account := range authGenState.GetAccounts() {

				if account.TypeUrl == "/cosmos.auth.v1beta1.BaseAccount" {
					_, ok := account.GetCachedValue().(authtypes.GenesisAccount)
					if ok {
						var byteAccounts []byte
						// Reason is prefix is nibiru --> getAddress will be empty
						// Marshal construct and convert to new struct to get address
						byteAccounts, err := codec.NewLegacyAmino().MarshalJSON(account.GetCachedValue())
						if err != nil {
							fmt.Printf("No account found for bank balance %s \n", string(byteAccounts))
							fmt.Println(err.Error())
							continue
						}
						var accountAfter Account
						if err := codec.NewLegacyAmino().UnmarshalJSON(byteAccounts, &accountAfter); err != nil {
							continue
						}

						snapshotAccs[accountAfter.Address] = SnapshotAccount{
							AtomAddress:         accountAfter.Address,
							AtomBalance:         sdk.ZeroInt(),
							AtomUnstakedBalance: sdk.ZeroInt(),
							AtomStakedBalance:   sdk.ZeroInt(),
						}
					}
				}

			}

			// Produce the map of address to total atom balance, both staked and unstaked

			totalAtomBalance := sdk.NewInt(0)
			denominator := sdk.NewInt(0)
			for _, account := range bankGenState.Balances {

				acc, ok := snapshotAccs[account.Address]
				if !ok {
					fmt.Printf("No account found for bank balance %s \n", account.Address)
					continue
				}
				balance := account.Coins.AmountOf(denom)
				totalAtomBalance = totalAtomBalance.Add(balance)

				// sum all sqrt of min(balance, max cap)

				denomi := getMin(balance.ToDec()).RoundInt()
				acc.AtomBalance = acc.AtomBalance.Add(balance)
				acc.AtomUnstakedBalance = acc.AtomUnstakedBalance.Add(balance)
				acc.Denominator = denomi
				denominator = denominator.Add(denomi)

				snapshotAccs[account.Address] = acc

			}

			for _, unbonding := range stakingGenState.UnbondingDelegations {
				address := unbonding.DelegatorAddress
				acc, ok := snapshotAccs[address]
				if !ok {
					fmt.Printf("No account found for unbonding %s \n", address)
					continue
				}

				unbondingAtoms := sdk.NewInt(0)
				for _, entry := range unbonding.Entries {
					unbondingAtoms = unbondingAtoms.Add(entry.Balance)
				}

				acc.AtomBalance = acc.AtomBalance.Add(unbondingAtoms)
				acc.AtomUnstakedBalance = acc.AtomUnstakedBalance.Add(unbondingAtoms)

				snapshotAccs[address] = acc
			}

			// Make a map from validator operator address to the v42 validator type
			validators := make(map[string]stakingtypes.Validator)
			for _, validator := range stakingGenState.Validators {
				validators[validator.OperatorAddress] = validator
			}

			for _, delegation := range stakingGenState.Delegations {
				address := delegation.DelegatorAddress

				acc, ok := snapshotAccs[address]
				if !ok {
					fmt.Printf("No account found for delegation address %s \n", address)
					continue
				}

				val := validators[delegation.ValidatorAddress]
				stakedAtoms := delegation.Shares.MulInt(val.Tokens).Quo(val.DelegatorShares).RoundInt()

				acc.AtomBalance = acc.AtomBalance.Add(stakedAtoms)
				acc.AtomStakedBalance = acc.AtomStakedBalance.Add(stakedAtoms)

				snapshotAccs[address] = acc
			}

			for address, acc := range snapshotAccs {
				allAtoms := acc.AtomBalance.ToDec()

				allAtomSqrt := getMin(allAtoms).RoundInt()

				if denominator.Int64() == 0 {
					acc.AtomOwnershipPercent = sdk.NewInt(0).ToDec()
				} else {
					acc.AtomOwnershipPercent = allAtomSqrt.ToDec().QuoInt(denominator)
				}

				if allAtoms.IsZero() {
					acc.AtomStakedPercent = sdk.ZeroDec()
					acc.GameBalance = sdk.ZeroInt()
					snapshotAccs[address] = acc
					continue
				}

				stakedAtoms := acc.AtomStakedBalance.ToDec()
				stakedPercent := stakedAtoms.Quo(allAtoms)
				acc.AtomStakedPercent = stakedPercent

				acc.GameBalance = acc.AtomOwnershipPercent.MulInt(sdk.NewInt(TotalGameAirdropAmount)).RoundInt()

				snapshotAccs[address] = acc
			}

			snapshot := Snapshot{
				TotalAtomAmount:        totalAtomBalance,
				TotalGameAirdropAmount: sdk.NewInt(TotalGameAirdropAmount),
				NumberAccounts:         uint64(len(snapshotAccs)),
				Denominator:            denominator,
				Accounts:               snapshotAccs,
			}

			fmt.Printf("# accounts: %d\n", len(snapshotAccs))
			fmt.Printf("atomTotalSupply: %s\n", totalAtomBalance.String())
			fmt.Printf("gameTotalSupply: %s\n", sdk.NewInt(TotalGameAirdropAmount).String())

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

// compare balance with max cap
func getMin(balance sdk.Dec) sdk.Dec {
	if balance.GTE(sdk.NewDec(MaxCap)) {
		atomSqrt, err := sdk.NewInt(MaxCap).ToDec().ApproxSqrt()
		if err != nil {
			panic(fmt.Sprintf("failed to root atom balance: %s", err))
		}
		return atomSqrt
	} else {
		atomSqrt, err := balance.ApproxSqrt()
		if err != nil {
			panic(fmt.Sprintf("failed to root atom balance: %s", err))
		}
		return atomSqrt
	}
}
