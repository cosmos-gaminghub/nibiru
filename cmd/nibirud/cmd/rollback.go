package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	dbm "github.com/tendermint/tm-db"
)

// RollbackCmd update state last results hash from block
func RollbackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "rollback when app crash",
		Long:  "nibirud rollback",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := cfg.DefaultBaseConfig()
			dbType := dbm.BackendType(config.DBBackend)

			// Get BlockStore
			clientCtx := client.GetClientContextFromCmd(cmd)
			config.RootDir = clientCtx.HomeDir
			blockStoreDB, err := dbm.NewDB("blockstore", dbType, config.DBDir())
			if err != nil {
				return err
			}
			blockStore := store.NewBlockStore(blockStoreDB)
			storeBlockHeight := blockStore.Height()
			// Get StateStore
			stateDB, err := dbm.NewDB("state", dbType, config.DBDir())
			if err != nil {
				return err
			}
			stateStore := state.NewStore(stateDB)
			state, err := stateStore.Load()
			if err != nil {
				return err
			}
			block := blockStore.LoadBlock(storeBlockHeight)
			state.LastResultsHash = block.LastResultsHash
			err = stateStore.Save(state)
			return err
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
