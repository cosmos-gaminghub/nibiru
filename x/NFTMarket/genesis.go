package NFTMarket

import (
	"github.com/cosmos-gaminghub/nibiru/x/NFTMarket/keeper"
	"github.com/cosmos-gaminghub/nibiru/x/NFTMarket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the item
	for _, elem := range genState.ItemList {
		k.SetItem(ctx, *elem)
	}

	// Set item count
	k.SetItemCount(ctx, genState.ItemCount)

	// this line is used by starport scaffolding # ibc/genesis/init
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all item
	itemList := k.GetAllItem(ctx)
	for _, elem := range itemList {
		elem := elem
		genesis.ItemList = append(genesis.ItemList, &elem)
	}

	// Set the current count
	genesis.ItemCount = k.GetItemCount(ctx)

	// this line is used by starport scaffolding # ibc/genesis/export

	return genesis
}
