package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
)

// GetDistributionAccount returns the distribution ModuleAccount
func (k Keeper) GetModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) SetModuleAccount(ctx sdk.Context, supplyKeeper types.SupplyKeeper) {
	moduleAcc := k.GetModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	supplyKeeper.SetModuleAccount(ctx, moduleAcc)
}


