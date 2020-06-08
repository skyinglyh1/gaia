package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/ft/internal/types"
	polycommon "github.com/cosmos/gaia/x/headersync/poly-utils/common"
	"strconv"
)

func (k Keeper) CreateDenom(ctx sdk.Context, creator sdk.AccAddress, denom string) sdk.Error {
	if reason, exist := k.ExistDenom(ctx, denom); exist {
		return sdk.ErrInternal(fmt.Sprintf("CreateDenom Error: denom:%s already exist, due to reason:%s", denom, reason))
	}
	//k.SetOperator(ctx, denom, creator)
	k.ccmKeeper.SetDenomCreator(ctx, denom, creator)
	ctx.KVStore(k.storeKey).Set(GetIndependentCrossDenomKey(denom), []byte(denom))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCoin,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, denom),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
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
	if !bytes.Equal([]byte(sourceAssetDenom), store.Get(GetIndependentCrossDenomKey(sourceAssetDenom))) {
		return sdk.ErrInternal(fmt.Sprintf("BindAssetHash, denom:%s is not designed to be able to be bondAssetHash through this interface", sourceAssetDenom))

	}
	store.Set(GetBindAssetHashKey([]byte(sourceAssetDenom), toChainId), toAssetHash)

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

	sink := polycommon.NewZeroCopySink(nil)
	args := types.TxArgs{
		ToAddress: toAddr,
		Amount:    amount.BigInt(),
	}
	if err := args.Serialization(sink, 32); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("TxArgs Serialization error:%v", err))
	}

	// burn coins from fromAddr
	if err := k.BurnCoins(ctx, fromAddr, sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, amount))); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("ft_crossed_independently.Lock.BurnCoins error:%v", err))
	}
	// get toAssetHash from storage
	toAssetHash := store.Get(GetBindAssetHashKey([]byte(sourceAssetDenom), toChainId))
	// invoke cross_chain_manager module to construct cosmos proof
	if sdkErr := k.ccmKeeper.CreateCrossChainTx(ctx, toChainId, []byte(sourceAssetDenom), toAssetHash, "unlock", sink.Bytes()); sdkErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Lock, CreateCrossChainTx error:%v", sdkErr))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
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

	var args types.TxArgs
	if err := args.Deserialization(polycommon.NewZeroCopySource(argsBs), 32); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("unlock, Deserialize args error:%s", err))
	}

	store := ctx.KVStore(k.storeKey)
	denom := string(toContractAddr)
	storedFromAssetHash := store.Get(GetBindAssetHashKey([]byte(denom), fromChainId))
	if !bytes.Equal(fromContractAddr, storedFromAssetHash) {
		return sdk.ErrInternal(fmt.Sprintf("Unlock, fromContractaddr:%x is not the stored assetHash:%x", fromContractAddr, storedFromAssetHash))
	}

	toAccAddr := sdk.AccAddress(args.ToAddress)
	amount := sdk.NewIntFromBigInt(args.Amount)
	if sdkErr := k.MintCoins(ctx, toAccAddr, sdk.NewCoins(sdk.NewCoin(denom, amount))); sdkErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("Unlock, burnCoins from Addr:%s error:%v", toAccAddr.String(), sdkErr))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(types.AttributeKeyFromChainId, strconv.FormatUint(fromChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyFromContractHash, hex.EncodeToString(fromContractAddr)),
			sdk.NewAttribute(types.AttributeKeyToAssetDenom, denom),
			sdk.NewAttribute(types.AttributeKeyToAddress, toAccAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) GetDenomInfo(ctx sdk.Context, denom string) *types.DenomInfo {
	operator := k.ccmKeeper.GetDenomCreator(ctx, denom)
	if len(operator) == 0 {
		return nil
	}
	return &types.DenomInfo{
		Creator:     operator,
		Denom:       denom,
		AssetHash:   []byte(denom),
		TotalSupply: k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom),
	}
}

func (k Keeper) GetDenomCrossChainInfo(ctx sdk.Context, denom string, toChainId uint64) *types.DenomCrossChainInfo {
	return &types.DenomCrossChainInfo{
		DenomInfo:   k.GetDenomInfo(ctx, denom),
		ToChainId:   toChainId,
		ToAssetHash: ctx.KVStore(k.storeKey).Get(GetBindAssetHashKey([]byte(denom), toChainId)),
	}
}

func (k Keeper) ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool {
	return ctx.KVStore(k.storeKey).Get((GetBindAssetHashKey(toContractAddr, fromChainId))) != nil
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
