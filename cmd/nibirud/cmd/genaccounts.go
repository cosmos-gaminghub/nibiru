package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos-gaminghub/nibiru/app"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
)

const (
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add a genesis account to genesis.json",
		Long: `Add a genesis account to genesis.json. The provided account must specify
the account address or key name and a list of initial coins. If a key name is given,
the address will be looked up in the local Keybase. The list of initial tokens must
contain valid denominations. Accounts may optionally be supplied with vesting parameters.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
				if err != nil {
					return err
				}

				// attempt to lookup address from Keybase if no address was provided
				kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf)
				if err != nil {
					return err
				}

				info, err := kb.Key(args[0])
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoinsNormalized(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse coins: %w", err)
			}

			vestingStart, err := cmd.Flags().GetInt64(flagVestingStart)
			if err != nil {
				return err
			}
			vestingEnd, err := cmd.Flags().GetInt64(flagVestingEnd)
			if err != nil {
				return err
			}
			vestingAmtStr, err := cmd.Flags().GetString(flagVestingAmt)
			if err != nil {
				return err
			}

			vestingAmt, err := sdk.ParseCoinsNormalized(vestingAmtStr)
			if err != nil {
				return fmt.Errorf("failed to parse vesting amount: %w", err)
			}

			// create concrete account type based on input parameters
			var genAccount authtypes.GenesisAccount

			balances := banktypes.Balance{Address: addr.String(), Coins: coins.Sort()}
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)

			if !vestingAmt.IsZero() {
				baseVestingAccount := authvesting.NewBaseVestingAccount(baseAccount, vestingAmt.Sort(), vestingEnd)

				if (balances.Coins.IsZero() && !baseVestingAccount.OriginalVesting.IsZero()) ||
					baseVestingAccount.OriginalVesting.IsAnyGT(balances.Coins) {
					return errors.New("vesting amount cannot be greater than total amount")
				}

				switch {
				case vestingStart != 0 && vestingEnd != 0:
					genAccount = authvesting.NewContinuousVestingAccountRaw(baseVestingAccount, vestingStart)

				case vestingEnd != 0:
					genAccount = authvesting.NewDelayedVestingAccountRaw(baseVestingAccount)

				default:
					return errors.New("invalid vesting parameters; must supply start and end time or end time")
				}
			} else {
				genAccount = baseAccount
			}

			if err := genAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

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

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}

			// Add the new account to the set of genesis accounts and sanitize the
			// accounts afterwards.
			accs = append(accs, genAccount)
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

			bankGenState := banktypes.GetGenesisStateFromAppState(clientCtx.Codec, appState)
			bankGenState.Balances = append(bankGenState.Balances, balances)
			bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)
			bankGenState.Supply = bankGenState.Supply.Add(balances.Coins...)

			bankGenStateBz, err := clientCtx.Codec.MarshalJSON(bankGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal bank genesis state: %w", err)
			}

			appState[banktypes.ModuleName] = bankGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Int64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Int64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
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
			fmt.Println(supply.AmountOf("game"))
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
			fmt.Println(supply.AmountOf("game"))

			var total sdk.Coins
			for _, balance := range liquidBalances {
				total = total.Add(balance.Coins...)
			}

			fmt.Println(total)

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
