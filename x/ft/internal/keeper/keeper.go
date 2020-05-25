package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/cosmos/gaia/x/ft/internal/types"
	selfexported "github.com/cosmos/gaia/x/lockproxy/exported"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the mint store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	paramSpace   params.Subspace
	authKeeper   types.AccountKeeper
	bankKeeper   types.BankKeeper
	supplyKeeper types.SupplyKeeper
	ccmKeeper    types.CrossChainManager
	selfexported.UnlockKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, ak types.AccountKeeper, bankKeeper types.BankKeeper, supplyKeeper types.SupplyKeeper, ccmKeeper types.CrossChainManager) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the crosschain module account has not been set")
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		authKeeper:   ak,
		supplyKeeper: supplyKeeper,
		ccmKeeper:    ccmKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) EnsureAccountExist(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	acct := k.authKeeper.GetAccount(ctx, addr)
	if acct == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("lockproxy: account %s does not exist", addr.String()))
	}
	return nil
}

func (k Keeper) MintCoins(ctx sdk.Context, toAcct sdk.AccAddress, amt sdk.Coins) sdk.Error {

	_, err := k.bankKeeper.AddCoins(ctx, toAcct, amt)
	if err != nil {
		panic(err)
	}

	// update total supply
	supply := k.supplyKeeper.GetSupply(ctx)
	supply = supply.Inflate(amt)

	k.supplyKeeper.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("minted coin:%s to account:%s ", amt.String(), toAcct.String()))

	return nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k Keeper) BurnCoins(ctx sdk.Context, fromAcct sdk.AccAddress, amt sdk.Coins) sdk.Error {

	_, err := k.bankKeeper.SubtractCoins(ctx, fromAcct, amt)
	if err != nil {
		panic(err)
	}

	// update total supply
	supply := k.supplyKeeper.GetSupply(ctx)
	supply = supply.Deflate(amt)
	k.supplyKeeper.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("burned coin:%s from account:%s ", amt.String(), fromAcct.String()))
	return nil
}
