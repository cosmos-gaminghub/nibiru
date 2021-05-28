package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	app "github.com/cosmos-gaminghub/nibiru/app"
	"github.com/cosmos-gaminghub/nibiru/cmd/nibirud/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {

	config := sdk.GetConfig()
	app.SetBech32AddressPrefixes(config)
	config.Seal()

	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
