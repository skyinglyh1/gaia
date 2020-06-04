package keeper

import (
	"encoding/hex"
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
	"github.com/cosmos/gaia/x/ft"
	"github.com/cosmos/gaia/x/headersync/internal/types"
	polycommon "github.com/cosmos/gaia/x/headersync/poly-utils/common"
	polytype "github.com/cosmos/gaia/x/headersync/poly-utils/core/types"
	"github.com/cosmos/gaia/x/lockproxy"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

var header0 = "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010ae3a2d1cba9ed56653edab871d93f8a96294debb6169a62681552dfd6d0fc70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c8365b000000001dac2b7c00000000fd1a057b226c6561646572223a343239343936373239352c227672665f76616c7565223a22484a675171706769355248566745716354626e6443456c384d516837446172364e4e646f6f79553051666f67555634764d50675851524171384d6f38373853426a2b38577262676c2b36714d7258686b667a72375751343d222c227672665f70726f6f66223a22785864422b5451454c4c6a59734965305378596474572f442f39542f746e5854624e436667354e62364650596370382f55706a524c572f536a5558643552576b75646632646f4c5267727052474b76305566385a69413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a343239343936373239352c226e65775f636861696e5f636f6e666967223a7b2276657273696f6e223a312c2276696577223a312c226e223a372c2263223a322c22626c6f636b5f6d73675f64656c6179223a31303030303030303030302c22686173685f6d73675f64656c6179223a31303030303030303030302c22706565725f68616e647368616b655f74696d656f7574223a31303030303030303030302c227065657273223a5b7b22696e646578223a312c226964223a2231323035303238313732393138353430623262353132656165313837326132613265336132386439383963363064393564616238383239616461376437646437303664363538227d2c7b22696e646578223a322c226964223a2231323035303338623861663632313065636664636263616232323535326566386438636634316336663836663963663961623533643836353734316366646238333366303662227d2c7b22696e646578223a332c226964223a2231323035303234383261636236353634623139623930363533663665396338303632393265386161383366373865376139333832613234613665666534316330633036663339227d2c7b22696e646578223a342c226964223a2231323035303236373939333061343261616633633639373938636138613366313265313334633031393430353831386437383364313137343865303339646538353135393838227d2c7b22696e646578223a352c226964223a2231323035303234363864643138393965643264316363326238323938383261313635613065636236613734356166306337326562323938326436366234333131623465663733227d2c7b22696e646578223a362c226964223a2231323035303265623162616162363032633538393932383235363163646161613761616262636464306363666362633365373937393361633234616366393037373866333561227d2c7b22696e646578223a372c226964223a2231323035303331653037373966356335636362323631323335326665346132303066393964336537373538653730626135336636303763353966663232613330663637386666227d5d2c22706f735f7461626c65223a5b362c342c332c352c362c312c322c352c342c372c342c322c332c332c372c362c352c342c362c352c312c342c332c312c322c352c322c322c362c312c342c352c342c372c322c332c342c312c352c372c342c312c322c322c352c362c342c342c322c372c332c362c362c352c312c372c332c312c362c312c332c332c322c342c342c312c352c362c352c312c322c362c372c352c362c332c342c372c372c332c322c372c312c352c362c352c322c332c362c322c362c312c372c372c372c312c372c342c332c332c332c322c312c372c355d2c226d61785f626c6f636b5f6368616e67655f76696577223a36303030307d7d9fe171f3fe643eb1c188400b828ba184816fc9ac0000"
var header1 = "000000000000000000000000f7259d9da6edb2672055c4f0efd8729f921ff4f2ea6cfe2c632bf9137a8eabbc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d43e5bb4e452c5130a39a3fa2f4e738e84b8caba1ab8a525eb0c379224a0c48d6c6dba5e010000005a6108a580c36ac2fd0c017b226c6561646572223a322c227672665f76616c7565223a22424c48634b703946724866376b64383866685a3644724748314f735178726f795a6a66766165664d5546337673517a36764a654e2b3252657a524a515a396e686143554759645544745869533232355851584b773563413d222c227672665f70726f6f66223a223037366b5331617a4551714a6e61706774546e554e4b5131576649435755596a2f65554e693469714b46615a4c3345614b715338385855737241396267594152717a4763764c6635792f435a612f745653336e504a773d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d000000000000000000000000000000000000000005231205038b8af6210ecfdcbcab22552ef8d8cf41c6f86f9cf9ab53d865741cfdb833f06b231205028172918540b2b512eae1872a2a2e3a28d989c60d95dab8829ada7d7dd706d658231205031e0779f5c5ccb2612352fe4a200f99d3e7758e70ba53f607c59ff22a30f678ff23120502679930a42aaf3c69798ca8a3f12e134c019405818d783d11748e039de851598823120502eb1baab602c5899282561cdaaa7aabbcdd0ccfcbc3e79793ac24acf90778f35a0542011b86005fa58d8286db7873bf9f1f116b59757518b36568bd2fe3e4c52d80710bc026a25f8dd3b45aa609e1c0f9b01cf43f2b94d061c936862dcedfb5d3c125830f42011caf13504a2c253135307f440cfa7053d0c96268c20c882b19c85753a1e4cc72fb1344f3ef00535304d3ad908959d393c906548bb078c52f14c6fd60036193072242011b2574b5ee43fb9345e90c1e3c8269a49b4f8b45266ccd6e783ffb858a9766c96362df590aa2e89bc8c086ddb2a4c80dc43b9eae52cbb539f8ddccfa61e018293142011bd8cf6c36d04358ed8bc4055ae372a5302dc18a7b4e56959a1be01b3a20b831c94e04e5623518512cbb38d2b80d6e4c2bb3e246f0f2cd94251f0f2ba54475eb4142011b46e5c26aab0b23e0594f0769909b36c4c2f6a9ea6393a17ff680ea7a901e00f31bb8271d2c1d019486fe7d142f3ddc943d9bee3d71c890da5e66d0f20eb53b9c"
var header100 = "0000000000000000000000008f9160cead73841bf2f4ad2a46a83c254ba386f2a2688712b0a1d57c1fe6b6f800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000fc6398ae59f2a21aa3bf8daae71f45b9585a32d4d202c9dfaf2c2422b66c94771679ba5e6400000006b5b5dcabb4b61dfd0c017b226c6561646572223a352c227672665f76616c7565223a224247786d68385a676b784a7568684a6f653570504f4c3032684d654c516b4a71416b317458486f567967616751597266786d78347a6673575170662b7956755062466a5747452f59624a37736f6f67556b776d666543413d222c227672665f70726f6f66223a2230475a794a326547333858594d4d7955654968316165705a42544c5138635771707231473147524d2f68494748625772765671646d6f446d735577624359493032386e5357454e31653373727a497a574c38454531773d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000523120502468dd1899ed2d1cc2b829882a165a0ecb6a745af0c72eb2982d66b4311b4ef73231205038b8af6210ecfdcbcab22552ef8d8cf41c6f86f9cf9ab53d865741cfdb833f06b23120502482acb6564b19b90653f6e9c806292e8aa83f78e7a9382a24a6efe41c0c06f3923120502eb1baab602c5899282561cdaaa7aabbcdd0ccfcbc3e79793ac24acf90778f35a231205031e0779f5c5ccb2612352fe4a200f99d3e7758e70ba53f607c59ff22a30f678ff0542011b65d702be46bfb18442ccd1241e135f300844770801a17f64fef9beb6458bd7e6491411f4c61be30ac971685940e78637b1b08e55ceb2dc1ec47fd183773d58f142011c8d004f33500811ed4ed0fe1cf878485cc0ec03c10948dcf9e5c8acd7f1cc46bb1deb2d39a10745fdbe91429660f0b8f111d1fd045c65f85e7bd588db3cd2baff42011beb42902869f1233ecc8523c817a0111fe0cfa3c86294af3348ca3d712a44236460bf6ad9d40b9990a888a871f6cb6e0022041032e12aba82e581177a0acf233e42011c7ecfb10b6d0507ba06deaa07ae1fda50e1cc8287fedb2e9fba27e4dba9edf1470029e94be290b9c0f95708ea9e66c5c6603262695120ca8e8bb1ade9c660ecb942011b742daab4c62493f848039dedf5a96aa03f456f1db721840d9f62f5e11d521ae443ccb0447dc10cda51133e6f215169c32123fcde648a48518de3c39d3243f6f1"

type testInput struct {
	cdc              *codec.Codec
	ctx              sdk.Context
	authKeeper       auth.AccountKeeper
	paramKeeper      params.Keeper
	supplyKeeper     supply.Keeper
	bankKeeper       bank.Keeper
	headerSyncKeeper Keeper
	ccmKeeper        ccm.Keeper
	lockProxyKeeper  lockproxy.Keeper
	ftkeeper         ft.Keeper
	btcxKeeper       btcx.Keeper
}

var (
	// module account permissions
	maccPerms = map[string][]string{
		ft.ModuleName:        {supply.Burner, supply.Minter},
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
	ft.RegisterCodec(cdc)
	ccm.RegisterCodec(cdc)
	lockproxy.RegisterCodec(cdc)
	btcx.RegisterCodec(cdc)
	types.RegisterCodec(cdc)

	//cdc.RegisterInterface((*interface{})(nil), nil)

	authKeyStore := sdk.NewKVStoreKey(auth.StoreKey)
	paramKeyStore := sdk.NewKVStoreKey(params.ModuleName)
	supplyKeyStore := sdk.NewKVStoreKey(supply.ModuleName)
	ccmStoreKey := sdk.NewKVStoreKey(ccm.StoreKey)
	headerSyncStoreKey := sdk.NewKVStoreKey(types.StoreKey)
	lockProxyStoreKey := sdk.NewKVStoreKey(lockproxy.StoreKey)
	btcxStoreKey := sdk.NewKVStoreKey(btcx.StoreKey)
	ftStoreKey := sdk.NewKVStoreKey(ft.StoreKey)

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

	headersyncKeeper := NewKeeper(cdc, headerSyncStoreKey, paramKeeper.Subspace(types.DefaultParamspace))
	ccmKeeper := ccm.NewKeeper(cdc, ccmStoreKey, paramKeeper.Subspace(ccm.DefaultParamspace), headersyncKeeper, nil)
	lockproxyKeeper := lockproxy.NewKeeper(cdc, lockProxyStoreKey, paramKeeper.Subspace(lockproxy.DefaultParamspace), authKeeper, supplyKeeper, ccmKeeper)
	btcxKeeper := btcx.NewKeeper(cdc, btcxStoreKey, paramKeeper.Subspace(btcx.DefaultParamspace), authKeeper, bankKeeper, supplyKeeper, ccmKeeper)
	ftKeeper := ft.NewKeeper(cdc, ftStoreKey, paramKeeper.Subspace(ft.DefaultParamspace), authKeeper, bankKeeper, supplyKeeper, lockproxyKeeper, ccmKeeper)
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
func Test_headersync_PolyHeader(t *testing.T) {
	var header polytype.Header
	h0s, _ := hex.DecodeString(header0)
	source := polycommon.NewZeroCopySource(h0s)
	header.Deserialization(source)
	fmt.Printf("header is %s\n", header.String())
}

func Test_headersync_SyncGenesisHeader(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx

	h0s, _ := hex.DecodeString(header0)
	err := input.headerSyncKeeper.SyncGenesisHeader(ctx, h0s)
	if err != nil {
		t.Errorf("SyncGenesisHeader error:%v", err)
	}
	var chainId uint64 = 0
	var height uint32 = 0
	keyHeights := input.headerSyncKeeper.GetKeyHeights(ctx, chainId)
	fmt.Printf("keyHeights are %+v\n", keyHeights)

	height = 1
	keyHeight, err := input.headerSyncKeeper.FindKeyHeight(ctx, height, chainId)
	if err != nil {
		t.Errorf("FindKeyHeight error:%v", err)
	}
	fmt.Printf("keyHeight of chainId: %d for height: %d is %d\n", chainId, height, keyHeight)

	h1s, _ := hex.DecodeString(header1)
	h100s, _ := hex.DecodeString(header100)
	err = input.headerSyncKeeper.SyncBlockHeaders(ctx, [][]byte{h1s, h100s})
	if err != nil {
		t.Errorf("SyncBlockHeaders error:%v", err)
	}

	height = 90
	keyHeight, err = input.headerSyncKeeper.FindKeyHeight(ctx, height, chainId)
	if err != nil {
		t.Errorf("FindKeyHeight error:%v", err)
	}
	fmt.Printf("keyHeight of chainId: %d for height: %d is %d\n", chainId, height, keyHeight)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func moduleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
