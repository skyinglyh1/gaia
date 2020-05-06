package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/gaia/x/crosschain/internal/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

type testInput struct {
	cdc      *codec.Codec
	ctx      sdk.Context
	authKeeper auth.AccountKeeper
	paramKeeper       params.Keeper
	supplyKeeper supply.Keeper
	bankKeeper bank.Keeper
	ccKeeper Keeper
}

var header0 = "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c033644e70a2b4f8de4a15c4a0cd79315673b8346d033804807058f3ff4252900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c8365b000000001dac2b7c00000000fd1a057b226c6561646572223a343239343936373239352c227672665f76616c7565223a22484a675171706769355248566745716354626e6443456c384d516837446172364e4e646f6f79553051666f67555634764d50675851524171384d6f38373853426a2b38577262676c2b36714d7258686b667a72375751343d222c227672665f70726f6f66223a22785864422b5451454c4c6a59734965305378596474572f442f39542f746e5854624e436667354e62364650596370382f55706a524c572f536a5558643552576b75646632646f4c5267727052474b76305566385a69413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a343239343936373239352c226e65775f636861696e5f636f6e666967223a7b2276657273696f6e223a312c2276696577223a312c226e223a372c2263223a322c22626c6f636b5f6d73675f64656c6179223a31303030303030303030302c22686173685f6d73675f64656c6179223a31303030303030303030302c22706565725f68616e647368616b655f74696d656f7574223a31303030303030303030302c227065657273223a5b7b22696e646578223a312c226964223a2231323035303364313031383338383037656334303739613436666539386436626439613036393061626362643863653136653066626334353230633763376566373838356462227d2c7b22696e646578223a322c226964223a2231323035303361366231623065316336393737663434663336323332613566336236316236613835396234636535313437633439616363666139613432663438336631323034227d2c7b22696e646578223a332c226964223a2231323035303266663764666337303562623561633638643265383932333063363632393939616562313832383431333165396663653934656639666166356239393137353364227d2c7b22696e646578223a342c226964223a2231323035303334343031376363636138323064393066306562623436316466343633333762303932336230616532626365353833636565316132363234633932303865323038227d2c7b22696e646578223a352c226964223a2231323035303331326631303233393531333134336330323938346263346561396438353438383366636466343937333264633732376466613734373438326663383037653634227d2c7b22696e646578223a362c226964223a2231323035303333336334343833376464623934616435666130656234363062306634393135346639303530333631396434643263386565303833333066623831353834316432227d2c7b22696e646578223a372c226964223a2231323035303363366536383165353135346566626136346337356230616131636135343438396261653736353330373764313664646439373236663336356265333036323264227d5d2c22706f735f7461626c65223a5b362c352c342c332c372c322c372c372c352c352c322c322c322c322c362c352c322c342c312c332c342c312c342c332c332c322c342c352c372c312c342c332c342c352c332c352c352c342c322c312c342c332c312c352c352c352c322c362c342c332c312c362c322c322c312c332c332c322c332c372c372c362c342c342c362c372c372c362c322c362c372c372c312c332c342c312c352c362c322c372c342c342c362c352c312c332c352c372c352c332c312c362c312c322c362c362c312c372c362c362c372c332c372c312c315d2c226d61785f626c6f636b5f6368616e67655f76696577223a31303030307d7dfd9e5473b163f591a8829d83288809d97c20ab2a0000"
var header1 = "000000000000000000000000e48232a8468647e98bf2af215912a02d81bae8f94f0eca4e01de1a86ec6331110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020b31282c6e546fa19c94726399748489950401f2c381f8232550cece1789ea05b98685e010000003705355644d2d063fd0c017b226c6561646572223a332c227672665f76616c7565223a224244714d70463966716a7754596d4e4c6c50795162464e505066546f374261786e677543775563697844713733636856386c6f6e525774474558596d77387671355279372f41505778434d573737433773594c33542b733d222c227672665f70726f6f66223a226d7272714956696352442b534963354c636b6a42426e427577355047466d2b4c5759414a6e4b486e4437537471636f4434376343506c6f6c366473716a42364f446150445a7a774430417332586330335937724c75513d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000623120502ff7dfc705bb5ac68d2e89230c662999aeb18284131e9fce94ef9faf5b991753d23120503d101838807ec4079a46fe98d6bd9a0690abcbd8ce16e0fbc4520c7c7ef7885db2312050344017ccca820d90f0ebb461df46337b0923b0ae2bce583cee1a2624c9208e2082312050333c44837ddb94ad5fa0eb460b0f49154f90503619d4d2c8ee08330fb815841d22312050312f10239513143c02984bc4ea9d854883fcdf49732dc727dfa747482fc807e6423120503c6e681e5154efba64c75b0aa1ca54489bae7653077d16ddd9726f365be30622d0642051c0265da7ed153880b8c78e594314fba05624a72ff06457d82f8db51ce35782b1e20b82154e8ad55113f155a47e1ccb1b56c773567e6678682e8f0308942f8329542051c0d0099185ff11de774e55b8ef4601c4b8672a5136debcd8ed50edb25b240866b23fc4bf4602e3d331219296cab8089ee6f6e0513b00ddee0f30468f58fcd56c142051ccdfc928b266182ba0fbed9741099b67517b14594e0c7220b97021c3b7b2efa240ccd293cdb4a7e1472bfd96e1ace841a67ac6b221a78d49f80a26dc81940ec2542051ca73d1477fffff931ea1846515ee62dc486dca1d1029e39fd9db3ce9821d8c7c004eaa8c38221ce4bcc267361297c2d62e38d68ac456e0d38f2f11ed57ceff05b42051b505034287c85586950d395af3a307976005ddb9f15eaef9af0e29d2ed187dc322e5a69123a3e964a91ad66c3519218182df9a3056aa50d68bc2201471ff4580942051b7e1151398a8700ced185d7e5428d15d118cc0bd450671f0a4548bab8e6aac13432e6ee56b1daa9395e6d9b39b71fb4e5252d65119291ba942d914cebc62eb787"
var header100 = "0000000000000000000000005905e989bcc6d6bc9e07e70403244d140d2c7b7cc728813bae0ac41ad70910fc000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004e51e28cfe3e5fce4862ff758343eb4a8de8169bb1ebadb5f15a68f9e909b1643ba0685e64000000f9fca3bf095522b4fd0c017b226c6561646572223a352c227672665f76616c7565223a2242502b4a356265477956366a66332f4c525a344256554b55496a78337a54755268426450686a7763496b5239794f796a34456b5445624a7971352b506a686e374f2b383179563348353167753273644f787848457977773d222c227672665f70726f6f66223a22616f66594c465457494c79514a54795868327533334d5a7934764e4e7035567341733752412b445030514d6b53756e693662736265307671737444493052554d6b6f545545612f716c45737232563069584439514c773d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d0000000000000000000000000000000000000000052312050312f10239513143c02984bc4ea9d854883fcdf49732dc727dfa747482fc807e6423120503d101838807ec4079a46fe98d6bd9a0690abcbd8ce16e0fbc4520c7c7ef7885db2312050344017ccca820d90f0ebb461df46337b0923b0ae2bce583cee1a2624c9208e20823120503a6b1b0e1c6977f44f36232a5f3b61b6a859b4ce5147c49accfa9a42f483f12042312050333c44837ddb94ad5fa0eb460b0f49154f90503619d4d2c8ee08330fb815841d20542051cb3282836d80eb2d0e90b689b7ee7c4a728892821ad253a124967fcc42c98640424bc2e0aad81e6f524196b120740538dd4c58c070342b16820a7e2006a9640c942051cdc90a420b3fd2dbcdfc05f3f929847562baa4f64278bbe7f28314a012645eea00c46bad250bb990d3c04b61d27f0d8bf04294d11048dfc0795c5dab23ad3b1c842051ba7ee9c6d0995e02272171dd9cc1d5e3eed20a4a41977a7c416968323c3ffbf935fa19adc107247bcbee2ca7eb37d6f21eb8d7964ed8b6e714f34a26e3fa76f1142051cb6e1d5cf25897220ea40fc81b1f390a2b9bd0c5085a9d08889eb0abbe94ad40c0278802d73ea8458be483fd6962d7f7b360125d0b717c3ba040eafe70cac9e1842051ca96c8a95e1731945eefce36640b11235dd9960f2bdbe20177ceefb05a9447fbf4791d12c7aaf296950f81cbb934d3f221dce2c38515813a0e894072f1ceb9ab5"
var header101 = "000000000000000000000000a5a30d77c60c9357a9b6928a482f393ef92f976beed3427c79af792e4eb87f12000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000bd5cf85551494e3c976a876627fa44fe8d1c3a916713380904e6352f7e3b04859a0685e65000000df41ce77c8f0aa7cfd0c017b226c6561646572223a362c227672665f76616c7565223a224242442b786571675854544d782f5a5a413674492f587354625744597143646e2b4e504b657a6568434d7a4335307a786662755a79695764492b526e3752584c36377568454a616d5151782b4c6d58496843385530614d3d222c227672665f70726f6f66223a22737a39635479517747752f6542346a7039384d302f62463475346d424a7247364235453059456e70794b654237655575793867764d45756a64704f632f504f78704c546b6c476569746446426c356c615733776a67513d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d0000000000000000000000000000000000000000052312050333c44837ddb94ad5fa0eb460b0f49154f90503619d4d2c8ee08330fb815841d223120503d101838807ec4079a46fe98d6bd9a0690abcbd8ce16e0fbc4520c7c7ef7885db2312050344017ccca820d90f0ebb461df46337b0923b0ae2bce583cee1a2624c9208e2082312050312f10239513143c02984bc4ea9d854883fcdf49732dc727dfa747482fc807e6423120503c6e681e5154efba64c75b0aa1ca54489bae7653077d16ddd9726f365be30622d0542051b008e43ad265add5727e4b1edc1eecd84681d1e335113a2030d2d077acd372c0a7f7da08702c6457b10419f37bbbf7eb9637f0c5016374d1ad505b11814174e1242051bc8bfe9a56634040e1e23e090269bde5f9279ade10ce7a75dfca5076fc0b63a231cba9db79ce97b7cbae26e522d69fcf308b237af72036792b2269708b2744a9a42051b9e052a2b09bea339b4f1637f738e3b0ca880e83df5af487ca3d73bdb5b2a76d736a6e8368e57255868027198af4af2b3216f72c289efedfb1a81c8241c51782d42051c80421a109388561e312f7e917d2ce00de612dba37c243178bb51d5bf4a3d577e619b4572a5a5a99260b3ca694014ae13007edcd5a7ef2acc20c41b6bdf5aaa8e42051b51d114b06cd0e13d30ce43bd32bcb1532092b4fff460c87bee92e78b1117ac2e31d547d972be4fae27ed94d9391616a2b221df1bfc0c926572fac0c6c331e75d"
var (
	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		types.ModuleName:             {supply.Burner, supply.Minter},
	}
)
func setupTestInput() testInput {
	newDb := dbm.NewMemDB()

	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	types.RegisterCodec(cdc)

	//cdc.RegisterInterface((*interface{})(nil), nil)

	authKeyStore := sdk.NewKVStoreKey(auth.StoreKey)

	crosschainKeyStore := sdk.NewKVStoreKey(types.ModuleName)
	paramKeyStore := sdk.NewKVStoreKey(params.ModuleName)
	tKey := sdk.NewTransientStoreKey(params.TStoreKey)
	supplyKeyStore := sdk.NewKVStoreKey(supply.ModuleName)

	ms := store.NewCommitMultiStore(newDb)
	ms.MountStoreWithDB(authKeyStore, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(supplyKeyStore, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(paramKeyStore, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(crosschainKeyStore, sdk.StoreTypeIAVL, newDb)
	ms.MountStoreWithDB(tKey, sdk.StoreTypeIAVL, newDb)

	ms.LoadLatestVersion()

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.AccAddress([]byte("moduleAcc")).String()] = true

	paramKeeper := params.NewKeeper(cdc, paramKeyStore, tKey, params.DefaultCodespace)

	authKeeper := auth.NewAccountKeeper(cdc, authKeyStore, paramKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(authKeeper, paramKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, moduleAccountAddrs())
	supplyKeeper := supply.NewKeeper(cdc, supplyKeyStore, authKeeper, bankKeeper, maccPerms)

	ccKeeper := NewCrossChainKeeper(cdc, crosschainKeyStore, paramKeeper.Subspace(types.DefaultParamspace), authKeeper, supplyKeeper)


	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	initialTotalSupply := sdk.NewCoins(sdk.NewCoin("nativecoin", sdk.NewInt(123456)))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(initialTotalSupply))

	return testInput{cdc: cdc, ctx: ctx, authKeeper: authKeeper, ccKeeper:ccKeeper, supplyKeeper:supplyKeeper, bankKeeper:bankKeeper}
}


// ModuleAccountAddrs returns all the app's module account addresses.
func moduleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func Test_Keeper_CreateCoins(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx

	//lockKeeper.GetModuleAccount(ctx)
	fmt.Printf("lockKeeper.GetModuleAccount = %v", input.ccKeeper.GetModuleAccount(ctx).String())


	addr := sdk.AccAddress([]byte(types.ModuleName))
	acc := input.authKeeper.NewAccountWithAddress(ctx, addr)
	moduleAcc := input.ccKeeper.GetModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	lockProxyModuleAddress := moduleAcc.GetAddress()
	fmt.Printf("crosschainModuleAddress  = %s\n", hex.EncodeToString(lockProxyModuleAddress.Bytes()))


	creator := acc.GetAddress()
	coinsStr := "1000000000ont,1000000000000000000ong"
	coins, err := sdk.ParseCoins(coinsStr)
	if err != nil {
		t.Errorf("parsecoins error:%v", err)
	}
	assert.Equal(t, addr, creator)

	total0 := input.supplyKeeper.GetSupply(ctx).GetTotal()


	err = input.ccKeeper.CreateCoins(ctx, creator, coins)
	if err != nil {
		t.Errorf("CreateCoins error:%v", err)
	}

	operator := input.ccKeeper.GetOperator(ctx)
	if !operator.Operator.Equals(acc.GetAddress()) {
		t.Errorf("storage operator not equal addr address, expect:%s, got:%s", acc.GetAddress().String(), operator.Operator.String())
	}


	total1 := input.supplyKeeper.GetSupply(ctx).GetTotal()

	assert.Equal(t, total1.IsEqual(total0.Add(coins)), true)

	balanceCoins := input.bankKeeper.GetCoins(ctx, input.supplyKeeper.GetModuleAddress(types.ModuleName))
	assert.Equal(t, balanceCoins.IsEqual(coins), true)
}


func Test_Keeper_BindProxyAndAssetHash(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx

	// generate base account with initial coins
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	creatorAcct := auth.NewBaseAccountWithAddress(addr)
	creatorAcct.Coins = sdk.NewCoins(sdk.NewCoin("ont", sdk.NewInt(1000000000)))
	creatorAcct.PubKey = pubKey
	creatorAcct.AccountNumber = uint64(0)
	input.authKeeper.SetAccount(ctx, &creatorAcct)


	// invoke GetModuleAccount() to initialize the lockproxy module account
	moduleAcc := input.ccKeeper.GetModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// create some coin in order to initialize the lockproxy operator
	coinsStr := "100somecoin"
	coins, err := sdk.ParseCoins(coinsStr)
	if err != nil {
		t.Errorf("parsecoins error:%v", err)
	}
	err = input.ccKeeper.CreateCoins(ctx, creatorAcct.GetAddress(), coins)
	if err != nil {
		t.Errorf("CreateCoins error:%v", err)
	}

	// operator should be the creator who first creates the some coin through CreateCoins() method
	operator := input.ccKeeper.GetOperator(ctx)
	assert.Equal(t, operator.Operator.String(), creatorAcct.Address.String())


	proxyHash := []byte{ 1, 2, 3, 4}
	input.ccKeeper.BindProxyHash(ctx, 3, proxyHash)
	storedProxyHash := input.ccKeeper.GetProxyHash(ctx, 3)
	assert.Equal(t, proxyHash, storedProxyHash)


	assetHash := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	crossedLimit := sdk.NewInt(1000000000)
	err = input.ccKeeper.BindAssetHash(ctx, "ont", 3,  assetHash, crossedLimit, true)
	assert.Nil(t, err)

	storedAssetHash := input.ccKeeper.GetAssetHash(ctx, "ont", 3)
	assert.Equal(t, assetHash, storedAssetHash)

	storedCrossedAmount := input.ccKeeper.GetCrossedAmount(ctx, "ont", 3)
	assert.Equal(t, storedCrossedAmount.String(), sdk.NewInt(1000000000).String())

	storedCrossedLimit := input.ccKeeper.GetCrossedLimit(ctx, "ont", 3)
	assert.Equal(t, storedCrossedLimit.String(), crossedLimit.String())




	assetHash1 := []byte{4, 3, 2, 1}
	sourceAssetDenom := coins[0].Denom
	crossedLimit1 := coins[0].Amount
	err = input.ccKeeper.BindAssetHash(ctx, sourceAssetDenom, 3,  assetHash1, crossedLimit1, false)
	assert.Nil(t, err)

	storedAssetHash1 := input.ccKeeper.GetAssetHash(ctx, sourceAssetDenom, 3)
	assert.Equal(t, assetHash1, storedAssetHash1)

	storedCrossedAmount1 := input.ccKeeper.GetCrossedAmount(ctx, sourceAssetDenom, 3)
	assert.Equal(t, storedCrossedAmount1.String(), sdk.NewInt(0).String())

	storedCrossedLimit1 := input.ccKeeper.GetCrossedLimit(ctx, sourceAssetDenom, 3)
	assert.Equal(t, storedCrossedLimit1.String(), crossedLimit1.String())

}


func Test_Keeper_Lock(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx

	// generate base account with initial coins
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	creatorAcct := auth.NewBaseAccountWithAddress(addr)
	coinsupply := sdk.NewCoins(sdk.NewCoin("ont", sdk.NewInt(1000000000)))
	creatorAcct.Coins = coinsupply
	creatorAcct.PubKey = pubKey
	creatorAcct.AccountNumber = uint64(0)
	input.authKeeper.SetAccount(ctx, &creatorAcct)
	input.supplyKeeper.SetSupply(ctx, supply.NewSupply(coinsupply))

	// invoke GetModuleAccount() to initialize the lockproxy module account
	moduleAcc := input.ccKeeper.GetModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// create some coin in order to initialize the lockproxy operator
	coinsStr := "100somecoin"
	coins, err := sdk.ParseCoins(coinsStr)
	if err != nil {
		t.Errorf("parsecoins error:%v", err)
	}
	err = input.ccKeeper.CreateCoins(ctx, creatorAcct.GetAddress(), coins)
	if err != nil {
		t.Errorf("CreateCoins error:%v", err)
	}

	// operator should be the creator who first creates the some coin through CreateCoins() method
	operator := input.ccKeeper.GetOperator(ctx)
	assert.Equal(t, operator.Operator.String(), creatorAcct.Address.String())


	proxyHash := []byte{ 1, 2, 3, 4}
	input.ccKeeper.BindProxyHash(ctx, 3, proxyHash)
	storedProxyHash := input.ccKeeper.GetProxyHash(ctx, 3)
	assert.Equal(t, proxyHash, storedProxyHash)


	assetHash1 := []byte{4, 3, 2, 1}
	err = input.ccKeeper.BindAssetHash(ctx, "ont", 3,  assetHash1, sdk.NewInt(1000000000), false)
	assert.Nil(t, err)


	toAddressBs := []byte{1, 1, 1, 1}
	amount := sdk.NewInt(10)
	err = input.ccKeeper.Lock(ctx, creatorAcct.Address, "ont", 3, toAddressBs, amount)
	assert.Nil(t, err)
}




func Test_Keeper_ProcessCrossChainTx(t *testing.T) {
	input := setupTestInput()
	ctx := input.ctx

	// generate base account with initial coins
	privKey := secp256k1.GenPrivKey()
	pubKey := privKey.PubKey()
	addr := sdk.AccAddress(pubKey.Address())
	creatorAcct := auth.NewBaseAccountWithAddress(addr)
	coinsupply := sdk.NewCoins(sdk.NewCoin("ont", sdk.NewInt(1000000000)))
	creatorAcct.Coins = coinsupply
	creatorAcct.PubKey = pubKey
	creatorAcct.AccountNumber = uint64(0)
	input.authKeeper.SetAccount(ctx, &creatorAcct)
	input.supplyKeeper.SetSupply(ctx, supply.NewSupply(coinsupply))

	// invoke GetModuleAccount() to initialize the lockproxy module account
	moduleAcc := input.ccKeeper.GetModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// create some coin in order to initialize the lockproxy operator
	coinsStr := "100somecoin"
	coins, err := sdk.ParseCoins(coinsStr)
	if err != nil {
		t.Errorf("parsecoins error:%v", err)
	}
	err = input.ccKeeper.CreateCoins(ctx, creatorAcct.GetAddress(), coins)
	if err != nil {
		t.Errorf("CreateCoins error:%v", err)
	}

	// operator should be the creator who first creates the some coin through CreateCoins() method
	operator := input.ccKeeper.GetOperator(ctx)
	assert.Equal(t, operator.Operator.String(), creatorAcct.Address.String())


	proxyHash := []byte{ 1, 2, 3, 4}
	input.ccKeeper.BindProxyHash(ctx, 3, proxyHash)
	storedProxyHash := input.ccKeeper.GetProxyHash(ctx, 3)
	assert.Equal(t, proxyHash, storedProxyHash)


	assetHash1 := []byte{4, 3, 2, 1}
	err = input.ccKeeper.BindAssetHash(ctx, "ont", 3,  assetHash1, sdk.NewInt(1000000000), false)
	assert.Nil(t, err)






	var header0 = "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c033644e70a2b4f8de4a15c4a0cd79315673b8346d033804807058f3ff4252900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c8365b000000001dac2b7c00000000fd1a057b226c6561646572223a343239343936373239352c227672665f76616c7565223a22484a675171706769355248566745716354626e6443456c384d516837446172364e4e646f6f79553051666f67555634764d50675851524171384d6f38373853426a2b38577262676c2b36714d7258686b667a72375751343d222c227672665f70726f6f66223a22785864422b5451454c4c6a59734965305378596474572f442f39542f746e5854624e436667354e62364650596370382f55706a524c572f536a5558643552576b75646632646f4c5267727052474b76305566385a69413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a343239343936373239352c226e65775f636861696e5f636f6e666967223a7b2276657273696f6e223a312c2276696577223a312c226e223a372c2263223a322c22626c6f636b5f6d73675f64656c6179223a31303030303030303030302c22686173685f6d73675f64656c6179223a31303030303030303030302c22706565725f68616e647368616b655f74696d656f7574223a31303030303030303030302c227065657273223a5b7b22696e646578223a312c226964223a2231323035303364313031383338383037656334303739613436666539386436626439613036393061626362643863653136653066626334353230633763376566373838356462227d2c7b22696e646578223a322c226964223a2231323035303361366231623065316336393737663434663336323332613566336236316236613835396234636535313437633439616363666139613432663438336631323034227d2c7b22696e646578223a332c226964223a2231323035303266663764666337303562623561633638643265383932333063363632393939616562313832383431333165396663653934656639666166356239393137353364227d2c7b22696e646578223a342c226964223a2231323035303334343031376363636138323064393066306562623436316466343633333762303932336230616532626365353833636565316132363234633932303865323038227d2c7b22696e646578223a352c226964223a2231323035303331326631303233393531333134336330323938346263346561396438353438383366636466343937333264633732376466613734373438326663383037653634227d2c7b22696e646578223a362c226964223a2231323035303333336334343833376464623934616435666130656234363062306634393135346639303530333631396434643263386565303833333066623831353834316432227d2c7b22696e646578223a372c226964223a2231323035303363366536383165353135346566626136346337356230616131636135343438396261653736353330373764313664646439373236663336356265333036323264227d5d2c22706f735f7461626c65223a5b362c352c342c332c372c322c372c372c352c352c322c322c322c322c362c352c322c342c312c332c342c312c342c332c332c322c342c352c372c312c342c332c342c352c332c352c352c342c322c312c342c332c312c352c352c352c322c362c342c332c312c362c322c322c312c332c332c322c332c372c372c362c342c342c362c372c372c362c322c362c372c372c312c332c342c312c352c362c322c372c342c342c362c352c312c332c352c372c352c332c312c362c312c322c362c362c312c372c362c362c372c332c372c312c315d2c226d61785f626c6f636b5f6368616e67655f76696577223a31303030307d7dfd9e5473b163f591a8829d83288809d97c20ab2a0000"
	mcSerializedHeader0Bs, _ := hex.DecodeString(header0)
	err = input.ccKeeper.SyncGenesisHeader(ctx, mcSerializedHeader0Bs)
	if err != nil {
		t.Errorf("SyncGenesisHeader error:%v", err)
	}


	header3307Bs, _ := hex.DecodeString("00000000000000000000000059f1daa0274cd31173dc2086cb10dd902497d54220e2447e699dd277c0f3df02000000000000000000000000000000000000000000000000000000000000000048374a23e74b942e9132fd413596cf8429957b79fe0aa130830a22e304f3d16eed109f40659334c5d1d6b0784db2380557c642e39f89821dddace55a9da549bfe5ff695eeb0c0000012f04bf42e36758fd0c017b226c6561646572223a362c227672665f76616c7565223a22424a73674144332f58725375565a426968635171637a6e774a384d2b766e656367623259324b44716b317141687568476b68574466344753735639444556324a4161626756464c7345776d365a61386a557a3542676a773d222c227672665f70726f6f66223a225a69476937315636624262636574493538354b4537377951346d6d42583174307a575a526d517866655475684f6a49646a47764c6876646243486736616e544a5a2b4d4f624a5679784778336775325264776c4d4b773d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d0000000000000000000000000000000000000000052312050333c44837ddb94ad5fa0eb460b0f49154f90503619d4d2c8ee08330fb815841d223120503d101838807ec4079a46fe98d6bd9a0690abcbd8ce16e0fbc4520c7c7ef7885db2312050344017ccca820d90f0ebb461df46337b0923b0ae2bce583cee1a2624c9208e20823120503c6e681e5154efba64c75b0aa1ca54489bae7653077d16ddd9726f365be30622d2312050312f10239513143c02984bc4ea9d854883fcdf49732dc727dfa747482fc807e640542051b2b302b7c3508a7dd547083ecdfd53efe02f473e8c4b35c3c46030a076a637a4e3fcb2e646c2cdb5ae4470574f73a896b26ff0dd6262ecc69b4bc3a93772d5c6a42051bbebf90f32376d4f2cccdd74d879dd22fe540f303e7e20d2a89c1074607e2aab70db254e3882adfc2d3174c255e6d469e58e001cf48e0a8d28a05441cf18a550b42051b9a4e60dff925475899e2b7736f15b2bb20775a4dcc5f2c3808f63f965b570e1336921585df8fea5ddc6713612907c0383af42a2156105546214ea1c621c08beb42051bf3f45dacfac533f3e2bb36357578da4f47562a9bc30b0591363e56e6dc5579822ca6dc612755944478c1b73ecbd6d13805eb4fa2053ca1eed9cdcc3ce4c9066442051b977ebf2bc5eeb253f90408a0a66e926a7561dc7262ead0b1b455d9cd4c096db63daf20950f7b30a51ada8c6112815bc9a37b8193a3266fbe33b71324871fab4d")
	fromChainId := uint64(3)
	height := uint32(3307)
	proof3306 := ""
	err = input.ccKeeper.ProcessCrossChainTx(ctx, fromChainId, height, proof3306, header3307Bs)
	assert.Nil(t, err)



}