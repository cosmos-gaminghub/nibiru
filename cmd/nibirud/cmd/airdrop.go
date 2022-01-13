package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos-gaminghub/nibiru/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	"github.com/spf13/cobra"
)

const (
	MaxCap                 = 50000000000
	TotalGameAirdropAmount = 1000000000000 // 0.5% * 200000000
)

type Snapshot struct {
	TotalAtomAmount        sdk.Int `json:"total_atom_amount"`
	TotalGameAirdropAmount sdk.Int `json:"total_game_amount"`
	NumberAccounts         uint64  `json:"num_accounts"`

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
		Use:   "export-airdrop-snapshot [airdrop-to-denom] [first-input-snapshot-file] [second-input-snapshot-file] [input-games-file]",
		Short: "Export a quadratic fairdrop snapshot from a provided cosmos-sdk v0.42 genesis export",
		Long: `Export a quadratic fairdrop snapshot from a provided cosmos-sdk v0.42 genesis export
Example:
	nibirud export-airdrop-snapshot uatom ~/.nibiru/config/genesis.json ../snapshot.json --nibiru-supply=100000000000000
	- Check input genesis:
		file is at ~/.nibirud/config/genesis.json
	- Snapshot
		file is at "../snapshot.json"
`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			denom := args[0]
			firstGenesisFile := args[1]
			secondGenesisFile := args[2]
			snapshotOutput := args[3]

			var snapshot Snapshot
			snapshot.Accounts = make(map[string]SnapshotAccount)
			snapshot.TotalGameAirdropAmount = sdk.ZeroInt()

			snapshot = exportSnapShotFromGenesisFile(clientCtx, firstGenesisFile, denom, snapshotOutput, snapshot)
			snapshot = exportSnapShotFromGenesisFile(clientCtx, secondGenesisFile, denom, snapshotOutput, snapshot)

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

func getDenominator(snapshotAccs map[string]SnapshotAccount) sdk.Int {
	denominator := sdk.ZeroInt()
	for _, acc := range snapshotAccs {
		allAtoms := acc.AtomBalance.ToDec()
		denominator = denominator.Add(getMin(allAtoms).RoundInt())
	}
	return denominator
}

func exportSnapShotFromGenesisFile(clientCtx client.Context, genesisFile string, denom string, snapshotOutput string, snapshot Snapshot) Snapshot {
	appState, _, _ := genutiltypes.GenesisStateFromGenFile(genesisFile)
	bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
	stakingGenState := stakingtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
	authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

	snapshotAccs := make(map[string]SnapshotAccount)
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

	for _, account := range bankGenState.Balances {

		acc, ok := snapshotAccs[account.Address]
		if !ok {
			fmt.Printf("No account found for bank balance %s \n", account.Address)
			continue
		}
		balance := account.Coins.AmountOf(denom)

		acc.AtomBalance = acc.AtomBalance.Add(balance)
		acc.AtomUnstakedBalance = acc.AtomUnstakedBalance.Add(balance)

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

	denominator := getDenominator(snapshotAccs)
	totalBalance := sdk.ZeroInt()
	totalAtomBalance := sdk.NewInt(0)
	for address, acc := range snapshotAccs {
		allAtoms := acc.AtomBalance.ToDec()

		allAtomSqrt := getMin(allAtoms).RoundInt()

		if denominator.IsZero() {
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

		totalBalance = totalBalance.Add(acc.GameBalance)
		snapshotAccount, ok := snapshot.Accounts[address]
		if !ok {
			snapshot.Accounts[address] = acc
			totalAtomBalance = totalAtomBalance.Add(acc.AtomBalance)
		} else {
			if snapshotAccount.GameBalance.IsNil() {
				snapshotAccount.GameBalance = sdk.ZeroInt()
			}
			snapshotAccount.GameBalance = snapshotAccount.GameBalance.Add(acc.GameBalance)
			snapshotAccount.AtomBalance = snapshotAccount.AtomBalance.Add(acc.AtomBalance)
			snapshotAccount.AtomUnstakedBalance = snapshotAccount.AtomUnstakedBalance.Add(acc.AtomUnstakedBalance)
			snapshot.Accounts[address] = snapshotAccount

			totalAtomBalance = totalAtomBalance.Add(acc.AtomBalance)
		}
	}
	snapshot.TotalAtomAmount = totalAtomBalance
	snapshot.TotalGameAirdropAmount = snapshot.TotalGameAirdropAmount.Add(totalBalance)
	snapshot.NumberAccounts = snapshot.NumberAccounts + uint64(len(snapshot.Accounts))

	fmt.Printf("Complete read genesis file %s \n", genesisFile)
	fmt.Printf("# accounts: %d\n", len(snapshotAccs))
	fmt.Printf("atomTotalSupply: %s\n", totalAtomBalance.String())
	fmt.Printf("gameTotalSupply: %s\n", totalBalance.String())
	return snapshot
}

func ImportGenesisAccountsFromSnapshotCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-genesis-accounts-from-snapshot [input-snapshot-file] [input-games-file]",
		Short: "Import genesis accounts from fairdrop snapshot.json and an games.json",
		Long: `Import genesis accounts from fairdrop snapshot.json
		20% of airdrop amount is liquid in accounts.
		The remaining is placed in the claims module.
		Must also pass in an games.json file to airdrop genesis games
		Example:
		nibirud import-genesis-accounts-from-snapshot ../snapshot.json ../games.json
		- Check input genesis:
			file is at ~/.nibirud/config/genesis.json
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			// aminoCodec := clientCtx.LegacyAmino.Amino

			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)

			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			authGenState := authtypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)

			accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			// Read snapshot file
			snapshotInput := args[0]
			snapshotJSON, err := os.Open(snapshotInput)
			if err != nil {
				return err
			}
			defer snapshotJSON.Close()
			byteValue, _ := ioutil.ReadAll(snapshotJSON)
			snapshot := Snapshot{}
			err = json.Unmarshal(byteValue, &snapshot)
			if err != nil {
				return err
			}

			// Read ions file
			gameInput := args[1]
			gameJSON, err := os.Open(gameInput)
			if err != nil {
				return err
			}
			defer gameJSON.Close()
			byteValue2, _ := ioutil.ReadAll(gameJSON)
			var gameAmts map[string]int64
			err = json.Unmarshal(byteValue2, &gameAmts)
			if err != nil {
				return err
			}

			// get genesis params
			genesisParams := MainnetGenesisParams()
			nonAirdropAccs := make(map[string]sdk.Coins)

			for _, acc := range genesisParams.DistributedAccounts {
				nonAirdropAccs[acc.Address] = acc.GetCoins()
			}

			for addr, amt := range gameAmts {
				// set atom bech32 prefixes
				bech32Addr, err := app.ConvertBech32(addr)
				if err != nil {
					return err
				}

				address, err := sdk.AccAddressFromBech32(bech32Addr)
				if err != nil {
					return err
				}

				if val, ok := nonAirdropAccs[address.String()]; ok {
					nonAirdropAccs[address.String()] = val.Add(sdk.NewCoin("game", sdk.NewInt(amt).MulRaw(1_000_000)))
				} else {
					nonAirdropAccs[address.String()] = sdk.NewCoins(sdk.NewCoin("game", sdk.NewInt(amt).MulRaw(1_000_000)))
				}
			}

			// figure out normalizationFactor to normalize snapshot balances to desired airdrop supply
			normalizationFactor := genesisParams.AirdropSupply.ToDec().QuoInt(snapshot.TotalGameAirdropAmount)
			fmt.Printf("normalization factor: %s\n", normalizationFactor)

			// for each account in the snapshot
			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			liquidBalances := bankGenState.Balances
			supply := bankGenState.Supply
			for _, acc := range snapshot.Accounts {
				// read address from snapshot
				bech32Addr, err := app.ConvertBech32(acc.AtomAddress)
				if err != nil {
					return err
				}

				address, err := sdk.AccAddressFromBech32(bech32Addr)
				if err != nil {
					return err
				}

				// initial liquid amounts
				// We consistently round down to the nearest uosmo
				liquidCoins := sdk.NewCoins(sdk.NewCoin(genesisParams.NativeCoinMetadatas[0].Base, acc.GameBalance))

				if coins, ok := nonAirdropAccs[address.String()]; ok {
					liquidCoins = liquidCoins.Add(coins...)
					delete(nonAirdropAccs, address.String())
				}

				liquidBalances = append(liquidBalances, banktypes.Balance{
					Address: address.String(),
					Coins:   liquidCoins,
				})
				supply = supply.Add(liquidCoins...)

				// Add the new account to the set of genesis accounts
				baseAccount := authtypes.NewBaseAccount(address, nil, 0, 0)
				if err := baseAccount.Validate(); err != nil {
					return fmt.Errorf("failed to validate new genesis account: %w", err)
				}
				accs = append(accs, baseAccount)

			}

			// distribute remaining game to accounts not in fairdrop
			for addr, coin := range nonAirdropAccs {
				// read address from snapshot
				address, err := sdk.AccAddressFromBech32(addr)
				if err != nil {
					return err
				}

				liquidBalances = append(liquidBalances, banktypes.Balance{
					Address: address.String(),
					Coins:   coin,
				})
				supply = supply.Add(coin...)

				// Add the new account to the set of genesis accounts
				baseAccount := authtypes.NewBaseAccount(address, nil, 0, 0)
				if err := baseAccount.Validate(); err != nil {
					return fmt.Errorf("failed to validate new genesis account: %w", err)
				}
				accs = append(accs, baseAccount)
			}
			// auth module genesis
			accs = authtypes.SanitizeGenesisAccounts(accs)
			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs
			authGenStateBz, err := clientCtx.Codec.MarshalJSON(&authGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}
			appState[authtypes.ModuleName] = authGenStateBz

			// bank module genesis
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(liquidBalances)
			bankGenState.Supply = supply
			bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}
			appState[banktypes.ModuleName] = bankGenStateBz

			byteIBCTransfer, err := appState[ibctransfertypes.ModuleName].MarshalJSON()
			if err != nil {
				return fmt.Errorf("Error marshal ibc transfer: %w", err)
			}

			var ibcGenState ibctransfertypes.GenesisState
			err = ibctransfertypes.ModuleCdc.UnmarshalJSON(byteIBCTransfer, &ibcGenState)
			if err != nil {
				return fmt.Errorf("Error unmarshal ibc transfer: %w", err)
			}
			ibcGenState.Params = ibctransfertypes.NewParams(false, false)
			ibcGenStateBz, err := clientCtx.Codec.MarshalJSON(&ibcGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal ibc genesis state: %w", err)
			}
			appState[ibctransfertypes.ModuleName] = ibcGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}
			genDoc.AppState = appStateJSON

			err = genutil.ExportGenesisFile(genDoc, genFile)
			return err
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
