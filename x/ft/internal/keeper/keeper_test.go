package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/gaia/x/btcx"
	"github.com/cosmos/gaia/x/ccm"
	"github.com/cosmos/gaia/x/ft/internal/types"
	"github.com/cosmos/gaia/x/headersync"
	"github.com/cosmos/gaia/x/lockproxy"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

type testInput struct {
	cdc              *codec.Codec
	ctx              sdk.Context
	authKeeper       auth.AccountKeeper
	paramKeeper      params.Keeper
	supplyKeeper     supply.Keeper
	bankKeeper       bank.Keeper
	headerSyncKeeper headersync.Keeper
	ccmKeeper        ccm.Keeper
	lockProxyKeeper  lockproxy.Keeper
	ftkeeper         Keeper
	btcxKeeper       btcx.Keeper
}

var (
	// module account permissions
	maccPerms = map[string][]string{
		types.ModuleName:     {supply.Burner, supply.Minter},
		btcx.ModuleName:      {supply.Burner, supply.Minter},
		lockproxy.ModuleName: {supply.Minter},
	}
)

func setupTestInput() testInput {
	newDb := dbm.NewMemDB()
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	auth.RegisterCodec(cdc)
	params.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	headersync.RegisterCodec(cdc)
	ccm.RegisterCodec(cdc)
	lockproxy.RegisterCodec(cdc)
	btcx.RegisterCodec(cdc)
	types.RegisterCodec(cdc)

	//cdc.RegisterInterface((*interface{})(nil), nil)

	authKeyStore := sdk.NewKVStoreKey(auth.StoreKey)
	paramKeyStore := sdk.NewKVStoreKey(params.ModuleName)
	supplyKeyStore := sdk.NewKVStoreKey(supply.ModuleName)
	ccmStoreKey := sdk.NewKVStoreKey(ccm.StoreKey)
	headerSyncStoreKey := sdk.NewKVStoreKey(headersync.StoreKey)
	lockProxyStoreKey := sdk.NewKVStoreKey(lockproxy.StoreKey)
	btcxStoreKey := sdk.NewKVStoreKey(btcx.StoreKey)
	ftStoreKey := sdk.NewKVStoreKey(types.StoreKey)

	tKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(newDb)
	ms.MountStoreWithDB(authKeyStore, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(paramKeyStore, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(supplyKeyStore, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(ccmStoreKey, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(headerSyncStoreKey, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(lockProxyStoreKey, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(btcxStoreKey, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(ftStoreKey, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(tKey, sdk.StoreTypeIAVL, newDb)

	ms.LoadLatestVersion()

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.AccAddress([]byte("moduleAcc")).String()] = true

	paramKeeper := params.NewKeeper(cdc, paramKeyStore, tKey, params.DefaultCodespace)

	authKeeper := auth.NewAccountKeeper(cdc, authKeyStore, paramKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(authKeeper, paramKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, moduleAccountAddrs())
	supplyKeeper := supply.NewKeeper(cdc, supplyKeyStore, authKeeper, bankKeeper, maccPerms)

	headersyncKeeper := headersync.NewKeeper(cdc, headerSyncStoreKey, paramKeeper.Subspace(headersync.DefaultParamspace))
	ccmKeeper := ccm.NewKeeper(cdc, ccmStoreKey, paramKeeper.Subspace(ccm.DefaultParamspace), headersyncKeeper, nil)
	lockproxyKeeper := lockproxy.NewKeeper(cdc, lockProxyStoreKey, paramKeeper.Subspace(lockproxy.DefaultParamspace), authKeeper, supplyKeeper, ccmKeeper)
	btcxKeeper := btcx.NewKeeper(cdc, btcxStoreKey, paramKeeper.Subspace(btcx.DefaultParamspace), authKeeper, bankKeeper, supplyKeeper, ccmKeeper)
	ftKeeper := NewKeeper(cdc, ftStoreKey, paramKeeper.Subspace(types.DefaultParamspace), authKeeper, bankKeeper, supplyKeeper, lockproxyKeeper, ccmKeeper)
	ccmKeeper.MountUnlockKeeperMap(map[string]ccm.UnlockKeeper{
		btcx.StoreKey:      btcxKeeper,
		types.StoreKey:     ftKeeper,
		lockproxy.StoreKey: lockproxyKeeper,
	})

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	var initialSupply supply.Supply

	supplyStore := ctx.KVStore(supplyKeyStore)
	b := supplyStore.Get(supply.SupplyKey)
	if b == nil {
		fmt.Printf("nil supply")
	} else {
		cdc.MustUnmarshalBinaryLengthPrefixed(b, &initialSupply)
	}

	var coins sdk.Coins
	if initialSupply.GetTotal().AmountOf("nativecoin").IsZero() {
		coins := append(coins, sdk.NewCoin("nativecoin", sdk.NewInt(0)))
		supplyKeeper.SetSupply(ctx, supply.NewSupply(coins))
	}
	lpModuleAcc := lockproxyKeeper.GetModuleAccount(ctx)
	if lpModuleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", lockproxy.ModuleName))
	}

	ftModuleAcc := ftKeeper.GetModuleAccount(ctx)
	if ftModuleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	ccmModuleAcc := ftKeeper.GetModuleAccount(ctx)
	if ccmModuleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", ccm.ModuleName))
	}
	return testInput{cdc: cdc, ctx: ctx, authKeeper: authKeeper, supplyKeeper: supplyKeeper, bankKeeper: bankKeeper, ccmKeeper: ccmKeeper, headerSyncKeeper: headersyncKeeper, lockProxyKeeper: lockproxyKeeper, btcxKeeper: btcxKeeper, ftkeeper: ftKeeper}
}

// ModuleAccountAddrs returns all the app's module account addresses.
func moduleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func Test_ftKeeper_CreateCoins(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx

	addr := sdk.AccAddress([]byte("acct1"))
	acc := input.authKeeper.NewAccountWithAddress(ctx, addr)

	creator := acc.GetAddress()
	coinsStr := "1000000000ont"
	coins, err := sdk.ParseCoins(coinsStr)
	if err != nil {
		t.Errorf("parsecoins error:%v", err)
	}
	assert.Equal(t, addr, creator)

	err = input.ftkeeper.CreateCoins(ctx, creator, coins)
	if err != nil {
		t.Errorf("CreateCoins error:%v", err)
	}

}
