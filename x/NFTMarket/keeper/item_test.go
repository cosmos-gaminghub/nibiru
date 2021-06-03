package keeper

import (
	"testing"

	"github.com/cosmos-gaminghub/nibiru/x/NFTMarket/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func createNItem(keeper *Keeper, ctx sdk.Context, n int) []types.Item {
	items := make([]types.Item, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Id = keeper.AppendItem(ctx, items[i])
	}
	return items
}

func TestItemGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNItem(keeper, ctx, 10)
	for _, item := range items {
		assert.Equal(t, item, keeper.GetItem(ctx, item.Id))
	}
}

func TestItemExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNItem(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasItem(ctx, item.Id))
	}
}

func TestItemRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNItem(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveItem(ctx, item.Id)
		assert.False(t, keeper.HasItem(ctx, item.Id))
	}
}

func TestItemGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNItem(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllItem(ctx))
}

func TestItemCount(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNItem(keeper, ctx, 10)
	count := uint64(len(items))
	assert.Equal(t, count, keeper.GetItemCount(ctx))
}
