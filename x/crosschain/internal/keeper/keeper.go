package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/cosmos/gaia/x/crosschain/internal/types"
)

type Keeper interface {
	HeaderSyncKeeper
	LockProxyKeeper
	GetModuleAccount(ctx sdk.Context) exported.ModuleAccountI
}


// Keeper of the mint store
type CrossChainKeeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	paramSpace   params.Subspace
	authKeeper types.AccountKeeper
	supplyKeeper types.SupplyKeeper
}


// NewKeeper creates a new mint Keeper instance
func NewCrossChainKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, ak types.AccountKeeper, supplyKeeper types.SupplyKeeper) CrossChainKeeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the crosschain module account has not been set")
	}

	return CrossChainKeeper{
		cdc:          cdc,
		storeKey:     key,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
		authKeeper: ak,
		supplyKeeper: supplyKeeper,
	}
}

func (k CrossChainKeeper) GetModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k CrossChainKeeper) EnsureAccountExist(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	acct := k.authKeeper.GetAccount(ctx, addr)
	if acct == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("lockproxy: account %s does not exist", addr.String()))
	}
	return nil
}
