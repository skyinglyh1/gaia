package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	polycommon "github.com/cosmos/gaia/x/headersync/poly-utils/common"
	selfexported "github.com/cosmos/gaia/x/lockproxy/exported"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	paramSpace   params.Subspace
	authKeeper   types.AccountKeeper
	supplyKeeper types.SupplyKeeper
	ccmKeeper    types.CrossChainManager
	selfexported.UnlockKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, ak types.AccountKeeper, supplyKeeper types.SupplyKeeper, ccmKeeper types.CrossChainManager) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the %s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		authKeeper:   ak,
		supplyKeeper: supplyKeeper,
		ccmKeeper:    ccmKeeper,
	}
}

func (k Keeper) GetModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) EnsureAccountExist(ctx sdk.Context, addr sdk.AccAddress) error {
	acct := k.authKeeper.GetAccount(ctx, addr)
	if acct == nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, fmt.Sprintf("lockproxy: account %s does not exist", addr.String()))
	}
	return nil
}

func (k Keeper) ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool {
	return ctx.KVStore(k.storeKey).Get((GetBindProxyKey(toContractAddr, fromChainId))) != nil
}

func (k Keeper) CreateLockProxy(ctx sdk.Context, creator sdk.AccAddress) error {
	if k.EnsureLockProxyExist(ctx, creator) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("CreateLockProxy Error: creator:%s already created lockproxy contract with hash:%x", creator.String(), creator.Bytes()))
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(GetOperatorToLockProxyKey(creator), creator.Bytes())
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateLockProxy,
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyProxyHash, hex.EncodeToString(creator.Bytes())),
		),
	})
	ctx.Logger().With("module", fmt.Sprintf("creator:%s initialized a lockproxy contract with hash: %x", creator.String(), creator.Bytes()))
	return nil
}

func (k Keeper) EnsureLockProxyExist(ctx sdk.Context, creator sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return bytes.Equal(store.Get(GetOperatorToLockProxyKey(creator)), creator)
}

func (k Keeper) GetLockProxyByOperator(ctx sdk.Context, operator sdk.AccAddress) []byte {
	store := ctx.KVStore(k.storeKey)
	proxyBytes := store.Get(GetOperatorToLockProxyKey(operator))
	if len(proxyBytes) == 0 || !bytes.Equal(operator.Bytes(), proxyBytes) {
		return nil
	}
	return proxyBytes
}

func (k Keeper) BindProxyHash(ctx sdk.Context, operator sdk.AccAddress, toChainId uint64, toProxyHash []byte) error {
	if !k.EnsureLockProxyExist(ctx, operator) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("BindProxyHash Error: operator:%s have NOT created lockproxy contract", operator.String()))
	}
	store := ctx.KVStore(k.storeKey)

	store.Set(GetBindProxyKey(operator, toChainId), toProxyHash)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindProxy,
			sdk.NewAttribute(types.AttributeKeyLockProxy, hex.EncodeToString(operator.Bytes())),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainProxyHash, hex.EncodeToString(toProxyHash)),
		),
	})
	return nil
}

func (k Keeper) GetProxyHash(ctx sdk.Context, operator sdk.AccAddress, toChainId uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(GetBindProxyKey(operator, toChainId))
}

func (k Keeper) BindAssetHash(ctx sdk.Context, operator sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAssetHash []byte, initialAmt sdk.Int) error {
	// ensure the operator has created the lockproxy contract
	if !k.EnsureLockProxyExist(ctx, operator) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("BindAssetHash,operator:%s have NOT created lockproxy contract", operator.String()))
	}
	// ensure the sourceAssetDenom has already been created with non-zero supply
	if !k.ExistDenom(ctx, sourceAssetDenom) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("BindAssetHash, sourceAssetDenom: %s NOT ", sourceAssetDenom))
	}
	//	ensure the passed-in initialAmt is equal to the balance of lockproxy module account
	moduleAcct := k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	if !moduleAcct.GetCoins().AmountOf(sourceAssetDenom).Equal(initialAmt) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("BindAssetHash, operator:%s, denom:%s, initialAmt incorrect, expect:%s, got:%s", operator.String(), sourceAssetDenom, moduleAcct.GetCoins().AmountOf(sourceAssetDenom).String(), initialAmt.String()))
	}

	store := ctx.KVStore(k.storeKey)
	// store the to asset hash based on the lockproxy contract (operator) and sourceAssetHash + toChainId
	store.Set(GetBindAssetHashKey(operator, []byte(sourceAssetDenom), toChainId), toAssetHash)
	// store the initial crossed amount
	store.Set(GetCrossedAmountKey([]byte(sourceAssetDenom)), k.cdc.MustMarshalBinaryLengthPrefixed(initialAmt))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindAsset,
			sdk.NewAttribute(types.AttributeKeyLockProxy, hex.EncodeToString(operator.Bytes())),
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeySourceAssetHash, hex.EncodeToString([]byte(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(toAssetHash)),
			sdk.NewAttribute(types.AttributeKeyInitialAmt, initialAmt.String()),
		),
	})
	return nil
}

func (k Keeper) ExistDenom(ctx sdk.Context, denom string) bool {
	storedSupplyCoins := k.supplyKeeper.GetSupply(ctx).GetTotal()
	return !storedSupplyCoins.AmountOf(denom).Equal(sdk.ZeroInt())
}

func (k Keeper) GetAssetHash(ctx sdk.Context, lockProxyHash []byte, sourceAssetDenom string, toChainId uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(GetBindAssetHashKey(lockProxyHash, []byte(sourceAssetDenom), toChainId))
}

func (k Keeper) GetLockedAmount(ctx sdk.Context, sourceAssetDenom string) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	crossedAmountBs := store.Get(GetCrossedAmountKey([]byte(sourceAssetDenom)))
	crossedAmount := sdk.NewInt(0)
	if crossedAmountBs != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(crossedAmountBs, &crossedAmount)
	}
	return crossedAmount
}
func (k Keeper) setLockededAmt(ctx sdk.Context, sourceAssetHash []byte, lockedAmt sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetCrossedAmountKey(sourceAssetHash), k.cdc.MustMarshalBinaryLengthPrefixed(lockedAmt))
}

func (k Keeper) Lock(ctx sdk.Context, lockProxyHash []byte, fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddressBs []byte, value sdk.Int) error {
	// send coin of sourceAssetDenom from fromAddress to module account address
	amt := sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, value))
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, fromAddress, types.ModuleName, amt); err != nil {
		return types.ErrSendCoinsToModuleFail(amt, fromAddress, k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetAddress())
	}
	store := ctx.KVStore(k.storeKey)

	sourceAssetHash := []byte(sourceAssetDenom)
	toChainAssetHash := store.Get(GetBindAssetHashKey(lockProxyHash, sourceAssetHash, toChainId))

	// get target asset hash from storage
	sink := polycommon.NewZeroCopySink(nil)
	args := types.TxArgs{
		ToAssetHash: toChainAssetHash,
		ToAddress:   toAddressBs,
		Amount:      value.BigInt(),
	}
	if err := args.Serialization(sink, 32); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("TxArgs Serialization error:%v", err))
	}
	// get target chain proxy hash from storage
	toChainProxyHash := store.Get(GetBindProxyKey(lockProxyHash, toChainId))
	fromContractHash := lockProxyHash
	if err := k.ccmKeeper.CreateCrossChainTx(ctx, toChainId, fromContractHash, toChainProxyHash, "unlock", sink.Bytes()); err != nil {
		return types.ErrCreateCrossChainTx(err)
	}
	if amt.AmountOf(sourceAssetDenom).IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("the coin being crossed has negative value, coin:%s", amt.String()))
	}
	k.setLockededAmt(ctx, sourceAssetHash, k.GetLockedAmount(ctx, sourceAssetDenom).Add(amt.AmountOf(sourceAssetDenom)))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyLockProxy, hex.EncodeToString(fromContractHash)),
			sdk.NewAttribute(types.AttributeKeyToChainProxyHash, hex.EncodeToString(toChainProxyHash)),
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeySourceAssetHash, hex.EncodeToString([]byte(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(toChainAssetHash)),
			sdk.NewAttribute(types.AttributeKeyFromAddress, fromAddress.String()),
			sdk.NewAttribute(types.AttributeKeyToAddress, hex.EncodeToString(toAddressBs)),
			sdk.NewAttribute(types.AttributeKeyAmount, value.String()),
		),
	})

	return nil
}

func (k Keeper) Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) error {

	proxyHash := k.GetProxyHash(ctx, toContractAddr, fromChainId)
	if len(proxyHash) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("the proxyHash is empty with chainId = %d", fromChainId))
	}
	if !bytes.Equal(proxyHash, fromContractAddr) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("stored toProxyHash correlated with lockproxyHash:%x is not equal to fromContractAddress-expect:%x, got:%x", toContractAddr, proxyHash, fromContractAddr))
	}
	args := new(types.TxArgs)
	if err := args.Deserialization(polycommon.NewZeroCopySource(argsBs), 32); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("unlock, Deserialization args error:%s", err))
	}
	toAssetHash := args.ToAssetHash
	toAddress := args.ToAddress
	amount := args.Amount

	// to asset hash should be the hex format string of source asset denom name, NOT Module account address
	toAssetDenom := string(toAssetHash)
	if len(k.GetAssetHash(ctx, toContractAddr, toAssetDenom, fromChainId)) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("toAssetHash:%x of denom:%s doesnot belong to the current lock proxy hash:%x", toAssetHash, toAssetDenom, toContractAddr))
	}

	// mint coin of sourceAssetDenom
	amt := sdk.NewCoins(sdk.NewCoin(toAssetDenom, sdk.NewIntFromBigInt(amount)))
	toAcctAddress := make(sdk.AccAddress, len(toAddress))
	copy(toAcctAddress, toAddress)

	if err := k.EnsureAccountExist(ctx, toAddress); err != nil {
		return err
	}
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddress, amt); err != nil {
		return types.ErrSendCoinsFromModuleFail(amt, k.GetModuleAccount(ctx).GetAddress(), toAddress)
	}
	newCrossedAmt := k.GetLockedAmount(ctx, toAssetDenom).Sub(sdk.NewIntFromBigInt(amount))
	if newCrossedAmt.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("new crossed amount is negative, storedCrossedAmt:%s, amount:%s", k.GetLockedAmount(ctx, toAssetDenom).String(), sdk.NewIntFromBigInt(amount).String()))
	}
	k.setLockededAmt(ctx, toAssetHash, newCrossedAmt)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(types.AttributeKeyFromChainId, strconv.FormatUint(fromChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyFromContractHash, hex.EncodeToString(fromContractAddr)),
			sdk.NewAttribute(types.AttributeKeyToAssetDenom, toAssetDenom),
			sdk.NewAttribute(types.AttributeKeyToAddress, toAcctAddress.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}
