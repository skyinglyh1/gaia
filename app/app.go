package app

import (
	"github.com/cosmos/gaia/x/btcx"
	"github.com/cosmos/gaia/x/ccm"
	"github.com/cosmos/gaia/x/ft"
	"github.com/cosmos/gaia/x/headersync"
	"github.com/cosmos/gaia/x/lockproxy"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"io"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	//hs "github.com/cosmos/gaia/x/headersync"
	//lp "github.com/cosmos/gaia/x/lockproxy"

	"os"
)

const appName = "GaiaApp"

var (
	//// default home directories for gaiacli
	//DefaultCLIHome = os.ExpandEnv("/home/skyinglyh/Go_Workspace/src/github.com/cosmos/gaia-build/.gaiacli")
	//// default home directories for gaiad
	//DefaultNodeHome = os.ExpandEnv("/home/skyinglyh/Go_Workspace/src/github.com/cosmos/gaia-build/.gaiad")

	//DefaultCLIHome = os.ExpandEnv("$HOME/.gaiacli")
	//DefaultNodeHome = os.ExpandEnv("$HOME/.gaiad")

	DefaultCLIHome  = os.ExpandEnv("./.gaiacli")
	DefaultNodeHome = os.ExpandEnv("./.gaiad")

	// The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		headersync.AppModuleBasic{}, ccm.AppModuleBasic{},
		lockproxy.AppModuleBasic{}, btcx.AppModuleBasic{}, ft.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		btcx.ModuleName:           {supply.Burner, supply.Minter},
		ft.ModuleName:             {supply.Burner, supply.Minter},
		lockproxy.ModuleName:      {supply.Minter},
	}
)

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc
}

// Extended ABCI application
type GaiaApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper

	//ccKeeper crosschain.Keeper

	ccmKeeper        ccm.Keeper
	headersyncKeeper headersync.Keeper
	lockproxyKeeper  lockproxy.Keeper
	btcxKeeper       btcx.Keeper
	ftKeeper         ft.Keeper

	// the module manager
	mm *module.Manager
}

// NewGaiaApp returns a reference to an initialized GaiaApp.
func NewGaiaApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *GaiaApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey,
		headersync.StoreKey, ccm.StoreKey,
		btcx.StoreKey, ft.StoreKey, lockproxy.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	app := &GaiaApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace)
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)

	//crosschainSubspace := app.paramsKeeper.Subspace(crosschain.DefaultParamspace)
	headerSyncSubspace := app.paramsKeeper.Subspace(headersync.DefaultParamspace)
	ccmSubspace := app.paramsKeeper.Subspace(ccm.DefaultParamspace)
	btcxSubspace := app.paramsKeeper.Subspace(btcx.DefaultParamspace)
	ftSubspace := app.paramsKeeper.Subspace(ft.DefaultParamspace)
	lockproxySubspace := app.paramsKeeper.Subspace(lockproxy.DefaultParamspace)

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, bankSubspace, bank.DefaultCodespace, app.ModuleAccountAddrs())
	app.supplyKeeper = supply.NewKeeper(app.cdc, keys[supply.StoreKey], app.accountKeeper, app.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		app.cdc, keys[staking.StoreKey], tkeys[staking.TStoreKey],
		app.supplyKeeper, stakingSubspace, staking.DefaultCodespace,
	)
	app.mintKeeper = mint.NewKeeper(app.cdc, keys[mint.StoreKey], mintSubspace, &stakingKeeper, app.supplyKeeper, auth.FeeCollectorName)
	app.distrKeeper = distr.NewKeeper(app.cdc, keys[distr.StoreKey], distrSubspace, &stakingKeeper,
		app.supplyKeeper, distr.DefaultCodespace, auth.FeeCollectorName, app.ModuleAccountAddrs())
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, keys[slashing.StoreKey], &stakingKeeper, slashingSubspace, slashing.DefaultCodespace,
	)
	app.crisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))
	app.govKeeper = gov.NewKeeper(
		app.cdc, keys[gov.StoreKey], app.paramsKeeper, govSubspace,
		app.supplyKeeper, &stakingKeeper, gov.DefaultCodespace, govRouter,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	//app.ccKeeper = crosschain.NewKeeper(app.cdc, keys[crosschain.StoreKey], crosschainSubspace, app.accountKeeper, app.supplyKeeper)
	app.headersyncKeeper = headersync.NewKeeper(app.cdc, keys[headersync.StoreKey], headerSyncSubspace)
	app.ccmKeeper = ccm.NewKeeper(app.cdc, keys[ccm.StoreKey], ccmSubspace, app.headersyncKeeper, nil)
	app.lockproxyKeeper = lockproxy.NewKeeper(app.cdc, keys[lockproxy.StoreKey], lockproxySubspace, app.accountKeeper, app.supplyKeeper, app.ccmKeeper)
	app.btcxKeeper = btcx.NewKeeper(app.cdc, keys[btcx.StoreKey], btcxSubspace, app.accountKeeper, app.bankKeeper, app.supplyKeeper, app.ccmKeeper)
	app.ftKeeper = ft.NewKeeper(app.cdc, keys[ft.StoreKey], ftSubspace, app.accountKeeper, app.bankKeeper, app.supplyKeeper, app.lockproxyKeeper, app.ccmKeeper)

	app.ccmKeeper.MountUnlockKeeperMap(map[string]ccm.UnlockKeeper{
		btcx.StoreKey:      app.btcxKeeper,
		ft.StoreKey:        app.ftKeeper,
		lockproxy.StoreKey: app.lockproxyKeeper,
	})

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.accountKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.distrKeeper, app.accountKeeper, app.supplyKeeper),

		//crosschain.NewAppModule(app.ccKeeper, app.supplyKeeper),

		headersync.NewAppModule(app.headersyncKeeper),
		lockproxy.NewAppModule(app.lockproxyKeeper),
		ccm.NewAppModule(app.ccmKeeper),
		btcx.NewAppModule(app.btcxKeeper),
		ft.NewAppModule(app.ftKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)

	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		genaccounts.ModuleName, distr.ModuleName, staking.ModuleName,
		auth.ModuleName, bank.ModuleName, slashing.ModuleName, gov.ModuleName,
		mint.ModuleName, supply.ModuleName, crisis.ModuleName, genutil.ModuleName,
		//crosschain.ModuleName,
		headersync.ModuleName, lockproxy.ModuleName, btcx.ModuleName, ft.ModuleName, ccm.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	return app
}

// application updates every begin block
func (app *GaiaApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *GaiaApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *GaiaApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *GaiaApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *GaiaApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
