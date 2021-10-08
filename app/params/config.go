package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	HumanCoinUnit = "GAME"
	BaseCoinUnit  = "game"
	GameExponent  = 6
)

func init() {
	RegisterDenoms()
}

func RegisterDenoms() {
	err := sdk.RegisterDenom(HumanCoinUnit, sdk.OneDec())
	if err != nil {
		panic(err)
	}
	err = sdk.RegisterDenom(BaseCoinUnit, sdk.NewDecWithPrec(1, GameExponent))
	if err != nil {
		panic(err)
	}
}
