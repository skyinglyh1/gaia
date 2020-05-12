package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/cosmos/gaia/x/crosschain/internal/types"
	"strconv"
)

type Keeper interface {
	HeaderSyncKeeper
	LockProxyKeeper
	GetModuleAccount(ctx sdk.Context) exported.ModuleAccountI
	CreateCoins(ctx sdk.Context, creator sdk.AccAddress, coins sdk.Coins) sdk.Error
	SetRedeemScript(ctx sdk.Context, denom string, redeemKey []byte, redeemScript []byte)
	BindNoVMChainAssetHash(ctx sdk.Context, denom string, targetChainId uint64, targetAssetHash []byte, limit sdk.Int) sdk.Error
}

// Keeper of the mint store
type CrossChainKeeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	paramSpace   params.Subspace
	authKeeper   types.AccountKeeper
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
		authKeeper:   ak,
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

func (k CrossChainKeeper) CreateCoins(ctx sdk.Context, creator sdk.AccAddress, coins sdk.Coins) sdk.Error {
	if k.GetOperator(ctx).Operator.Empty() {
		k.SetOperator(ctx, types.Operator{creator})
	}
	if !coins.IsAllPositive() {
		return types.ErrCreateNegativeCoins(types.DefaultCodespace, coins)
	}

	var increments sdk.Coins
	storedSupplyCoins := k.supplyKeeper.GetSupply(ctx).GetTotal()
	for _, coin := range coins {
		oldCoinAmount := storedSupplyCoins.AmountOf(coin.Denom)
		if oldCoinAmount == sdk.ZeroInt() {
			storedSupplyCoins = append(storedSupplyCoins, coin)
		} else {
			increment := coin.Sub(sdk.NewCoin(coin.Denom, oldCoinAmount))
			increments = append(increments, increment)
		}
	}

	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, increments); err != nil {
		return types.ErrSupplyKeeperMintCoinsFail(types.DefaultCodespace)
	}
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("minted %s from %s module account", coins.String(), types.ModuleName))
	return nil
}

func (k CrossChainKeeper) SetRedeemScript(ctx sdk.Context, denom string, redeemKey []byte, redeemScript []byte) {
	store := ctx.KVStore(k.storeKey)
	//calculatedRedeemKey := btcutil.Hash160(redeemScriptBytes)
	store.Set(GetRedeemScriptKey(redeemKey), redeemScript)
	store.Set(GetDenomToHashKey(denom), redeemKey)
	store.Set(GetHashKeyToDenom(redeemKey), types.DenomToHash(denom))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetRedeemScript,
			sdk.NewAttribute(types.AttributeKeyRedeemKey, hex.EncodeToString(redeemKey)),
			sdk.NewAttribute(types.AttributeKeyRedeemScript, hex.EncodeToString(redeemScript)),
		),
	})
}

//BindAssetHash(ctx sdk.Context, sourceAssetDenom string, targetChainId uint64, targetAssetHash []byte, limit sdk.Int, isTargetChainAsset bool) sdk.Error
func (k CrossChainKeeper) BindNoVMChainAssetHash(ctx sdk.Context, denom string, targetChainId uint64, targetAssetHash []byte, limit sdk.Int) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	sourceAssetKey := store.Get(GetDenomToHashKey(denom))
	if sourceAssetKey == nil {
		return sdk.ErrInternal(fmt.Sprintf("there is no script key corresponded with denom:%s, please SetRedeemScript first", denom))
	}
	if err := k.BindAssetHash(ctx, string(sourceAssetKey), targetChainId, targetAssetHash, limit, true); err != nil {
		return err
	}

	store.Set(GetKeyToHashKey(sourceAssetKey, targetChainId), targetAssetHash)
	store.Set(GetContractToScriptKey(targetAssetHash, targetChainId), sourceAssetKey)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetNoVmChainAssetHash,
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(targetAssetHash)),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(targetChainId, 10)),
			sdk.NewAttribute(types.AttributeKeySourceRedeemKey, hex.EncodeToString(sourceAssetKey)),
		),
	})
	return nil
}
