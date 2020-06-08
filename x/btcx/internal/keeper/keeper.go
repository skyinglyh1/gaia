package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/gaia/x/btcx/exported"
	"github.com/cosmos/gaia/x/btcx/internal/types"
	polycommon "github.com/cosmos/gaia/x/headersync/poly-utils/common"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"strconv"
)

// Keeper of the mint store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	paramSpace   params.Subspace
	authKeeper   types.AccountKeeper
	bankKeeper   types.BankKeeper
	supplyKeeper types.SupplyKeeper
	ccmKeeper    types.CCMKeeper
	exported.UnlockKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, ak types.AccountKeeper, bk types.BankKeeper, supplyKeeper types.SupplyKeeper, ccmKeeper types.CCMKeeper) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the %s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		authKeeper:   ak,
		bankKeeper:   bk,
		supplyKeeper: supplyKeeper,
		ccmKeeper:    ccmKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) EnsureAccountExist(ctx sdk.Context, addr sdk.AccAddress) sdk.Error {
	acct := k.authKeeper.GetAccount(ctx, addr)
	if acct == nil {
		return sdk.ErrUnknownAddress(fmt.Sprintf("lockproxy: account %s does not exist", addr.String()))
	}
	return nil
}

func (k Keeper) CreateDenom(ctx sdk.Context, creator sdk.AccAddress, denom string, redeemScript string) sdk.Error {
	if reason, exist := k.ExistDenom(ctx, denom); exist {
		return sdk.ErrInternal(fmt.Sprintf("CreateCoins Error: denom:%s already exist, due to reason:%s", denom, reason))
	}
	//k.SetOperator(ctx, denom, creator)
	k.ccmKeeper.SetDenomCreator(ctx, denom, creator)

	redeemScriptBs, err := hex.DecodeString(redeemScript)
	if err != nil {
		panic("Invalid Redeem Script")
	}
	store := ctx.KVStore(k.storeKey)
	scriptHashKeyBs := btcutil.Hash160(redeemScriptBs)
	store.Set(GetDenomToScriptHashKey(denom), scriptHashKeyBs)
	store.Set(GetScriptHashToDenomKey(scriptHashKeyBs), []byte(denom))

	store.Set(GetScriptHashToRedeemScript(scriptHashKeyBs), redeemScriptBs)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateDenom,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, denom),
			sdk.NewAttribute(types.AttributeKeyRedeemKey, hex.EncodeToString(scriptHashKeyBs)),
			sdk.NewAttribute(types.AttributeKeyRedeemScript, hex.EncodeToString(redeemScriptBs)),
		),
	})
	k.Logger(ctx).Info(fmt.Sprintf("creator:%s initialized denom: %s ", creator.String(), denom))
	return nil
}

func (k Keeper) BindAssetHash(ctx sdk.Context, creator sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAssetHash []byte) sdk.Error {
	if !k.ValidCreator(ctx, sourceAssetDenom, creator) {
		//return sdk.ErrInternal(fmt.Sprintf("BindAssetHash, creator is not valid, expect:%s, got:%s", k.GetOperator(ctx, sourceAssetDenom).String(), creator.String()))
		return sdk.ErrInternal(fmt.Sprintf("BindAssetHash, creator is not valid, expect:%s, got:%s", k.ccmKeeper.GetDenomCreator(ctx, sourceAssetDenom).String(), creator.String()))
	}
	store := ctx.KVStore(k.storeKey)
	// for lock usage
	scriptHash := store.Get(GetDenomToScriptHashKey(sourceAssetDenom))
	// for unlock usage
	store.Set(GetScriptHashAndChainIdToAssetHashKey(scriptHash, toChainId), toAssetHash)
	// for unlock usage

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindAsset,
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeyFromAssetHash, hex.EncodeToString(sdk.AccAddress(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(toAssetHash)),
		),
	})
	return nil
}

func (k Keeper) Lock(ctx sdk.Context, fromAddr sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddr []byte, amount sdk.Int) sdk.Error {
	// transfer back to btc
	store := ctx.KVStore(k.storeKey)
	redeemScriptHash := store.Get(GetDenomToScriptHashKey(sourceAssetDenom))
	if redeemScriptHash == nil {
		return sdk.ErrInternal(fmt.Sprintf("Invoke Lock of `btcx` module for denom: %s is illgeal", sourceAssetDenom))
	}
	sink := polycommon.NewZeroCopySink(nil)
	// contruct args bytes
	if toChainId == 1 {
		redeemScript := store.Get(GetScriptHashToRedeemScript(redeemScriptHash))
		toBtcArgs := types.ToBTCArgs{
			ToBtcAddress: toAddr,
			Amount:       amount.BigInt().Uint64(),
			RedeemScript: redeemScript,
		}
		if err := toBtcArgs.Serialization(sink); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("ToBTCArgs Serialization error:%v", err))
		}
	} else {
		args := types.BTCArgs{
			ToBtcAddress: toAddr,
			Amount:       amount.BigInt().Uint64(),
		}
		if err := args.Serialization(sink); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("BTCArgs Serialization error:%v", err))
		}
	}

	// burn coins from fromAddr
	if err := k.BurnCoins(ctx, fromAddr, sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, amount))); err != nil {
		return types.ErrMintCoinsFail(types.DefaultCodespace)
	}
	// get toAssetHash from storage
	toAssetHash := store.Get(GetScriptHashAndChainIdToAssetHashKey(redeemScriptHash, toChainId))
	// invoke cross_chain_manager module to construct cosmos proof
	if sdkErr := k.ccmKeeper.CreateCrossChainTx(ctx, toChainId, redeemScriptHash, toAssetHash, "unlock", sink.Bytes()); sdkErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Lock, CreateCrossChainTx error:%v", sdkErr))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeySourceAssetHash, hex.EncodeToString([]byte(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(toAssetHash)),
			sdk.NewAttribute(types.AttributeKeyFromAddress, fromAddr.String()),
			sdk.NewAttribute(types.AttributeKeyToAddress, hex.EncodeToString(toAddr)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) sdk.Error {

	var args types.BTCArgs
	if err := args.Deserialization(polycommon.NewZeroCopySource(argsBs)); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("unlock, Deserialize args error:%s", err))
	}
	store := ctx.KVStore(k.storeKey)
	toAssetHash := store.Get(GetScriptHashAndChainIdToAssetHashKey(toContractAddr, fromChainId))
	if !bytes.Equal(fromContractAddr, toAssetHash) {
		return sdk.ErrInternal(fmt.Sprintf("Unlock, fromContractaddr:%x is not the stored assetHash:%x", fromContractAddr, toAssetHash))
	}
	toDenom := string(store.Get(GetScriptHashToDenomKey(toContractAddr)))

	toAccAddr := sdk.AccAddress(args.ToBtcAddress)
	amount := sdk.NewIntFromBigInt(big.NewInt(0).SetUint64(args.Amount))
	if sdkErr := k.MintCoins(ctx, toAccAddr, sdk.NewCoins(sdk.NewCoin(toDenom, amount))); sdkErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Unlock, burnCoins from Addr:%s error:%v", toAccAddr.String(), sdkErr))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(types.AttributeKeyFromChainId, strconv.FormatUint(fromChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyFromContractHash, hex.EncodeToString(fromContractAddr)),
			sdk.NewAttribute(types.AttributeKeyToAssetDenom, toDenom),
			sdk.NewAttribute(types.AttributeKeyToAddress, toAccAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) GetDenomInfo(ctx sdk.Context, denom string) *types.DenomInfo {

	store := ctx.KVStore(k.storeKey)
	//operator := store.Get(GetDenomToOperatorKey(denom))
	operator := k.ccmKeeper.GetDenomCreator(ctx, denom)
	if len(operator) == 0 {
		return nil
	}
	denomInfo := new(types.DenomInfo)
	denomInfo.Creator = operator
	denomInfo.TotalSupply = k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
	denomInfo.RedeemScriptHash = store.Get(GetDenomToScriptHashKey(denom))
	denomInfo.RedeemScipt = store.Get(GetScriptHashToRedeemScript(denomInfo.RedeemScriptHash))
	return denomInfo
}

func (k Keeper) GetDenomCrossChainInfo(ctx sdk.Context, denom string, toChainId uint64) *types.DenomCrossChainInfo {
	denomInfo := new(types.DenomCrossChainInfo)
	denomInfo.DenomInfo = *k.GetDenomInfo(ctx, denom)
	denomInfo.ToChainId = toChainId
	denomInfo.ToAssetHash = ctx.KVStore(k.storeKey).Get(GetScriptHashAndChainIdToAssetHashKey(denomInfo.RedeemScriptHash, toChainId))
	return denomInfo
}

func (k Keeper) ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool {
	return ctx.KVStore(k.storeKey).Get((GetScriptHashAndChainIdToAssetHashKey(toContractAddr, fromChainId))) != nil
}

func (k Keeper) ValidCreator(ctx sdk.Context, denom string, creator sdk.AccAddress) bool {
	//store := ctx.KVStore(k.storeKey)
	//return bytes.Equal(store.Get(GetDenomToOperatorKey(denom)), creator.Bytes())
	return bytes.Equal(k.ccmKeeper.GetDenomCreator(ctx, denom), creator.Bytes())
}

func (k Keeper) ExistDenom(ctx sdk.Context, denom string) (string, bool) {
	storedSupplyCoins := k.supplyKeeper.GetSupply(ctx).GetTotal()
	//return storedSupplyCoins.AmountOf(denom) != sdk.ZeroInt() || len(k.GetOperator(ctx, denom)) != 0
	if len(k.ccmKeeper.GetDenomCreator(ctx, denom)) != 0 {
		return fmt.Sprintf("k.ccmKeeper.GetDenomCreator(ctx,%s) is %x", denom, k.ccmKeeper.GetDenomCreator(ctx, denom)), true
	}
	if !storedSupplyCoins.AmountOf(denom).Equal(sdk.ZeroInt()) {
		return fmt.Sprintf("supply.AmountOf(%s) is %s", denom, storedSupplyCoins.AmountOf(denom).String()), true
	}
	return "", false
}

//func (k Keeper) SetOperator(ctx sdk.Context, denom string, creator sdk.AccAddress) {
//	ctx.KVStore(k.storeKey).Set(GetDenomToOperatorKey(denom), creator.Bytes())
//}
//
//func (k Keeper) GetOperator(ctx sdk.Context, denom string) sdk.AccAddress {
//	return ctx.KVStore(k.storeKey).Get(GetDenomToOperatorKey(denom))
//
//}

// MintCoins creates new coins from thin air and adds it to the module account.
// Panics if the name maps to a non-minter module account or if the amount is invalid.
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
