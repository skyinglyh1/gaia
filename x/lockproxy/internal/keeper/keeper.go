package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	ttype "github.com/tendermint/tendermint/types"

	"bytes"
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
	mcc "github.com/ontio/multi-chain/common"
	mctype "github.com/ontio/multi-chain/core/types"
	"github.com/ontio/multi-chain/merkle"
	ccmc "github.com/ontio/multi-chain/native/service/cross_chain_manager/common"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strconv"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// Keeper of the mint store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	paramSpace   params.Subspace
	supplyKeeper types.SupplyKeeper
	hsKeeper     types.SyncKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, supplyKeeper types.SupplyKeeper, hsKeeper types.SyncKeeper) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the lockproxy module account has not been set")
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper: supplyKeeper,
		hsKeeper:     hsKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// get the minter
func (k Keeper) GetOperator(ctx sdk.Context) (operator types.Operator) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.OperatorKey)
	if b == nil {
		operator = types.Operator{}
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &operator)
	return
}

// set the minter
func (k Keeper) SetOperator(ctx sdk.Context, operator types.Operator) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(operator)
	store.Set(types.OperatorKey, b)
}

//______________________________________________________________________

// GetParams returns the total set of minting parameters.
func (k Keeper) GetCoinsParam(ctx sdk.Context) (params types.CoinsParam) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetCoinsParam(ctx sdk.Context, coinsParam types.CoinsParam) {
	k.paramSpace.SetParamSet(ctx, &coinsParam)
}

func (k Keeper) BindProxyHash(ctx sdk.Context, targetChainId uint64, targetProxyHash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetBindProxyKey(targetChainId), targetProxyHash)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindProxy,
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(targetChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainProxyHash, hex.EncodeToString(targetProxyHash)),
		),
	})
}

func (k Keeper) GetProxyHash(ctx sdk.Context, toChainId uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(GetBindProxyKey(toChainId))
}

func (k Keeper) GetAssetHash(ctx sdk.Context, sourceAssetDenom string, toChainId uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := DenomToHash(sourceAssetDenom)
	return store.Get(GetBindAssetKey(sourceAssetHash.Bytes(), toChainId))
}

func (k Keeper) GetCrossedAmount(ctx sdk.Context, sourceAssetDenom string, toChainId uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := DenomToHash(sourceAssetDenom)
	crossedAmountBs := store.Get(GetCrossedAmountKey(sourceAssetHash.Bytes(), toChainId))
	crossedAmount := sdk.NewInt(0)
	if crossedAmountBs != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(crossedAmountBs, &crossedAmount)
	}
	return crossedAmount
}


func (k Keeper) GetCrossedLimit(ctx sdk.Context, sourceAssetDenom string, toChainId uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := DenomToHash(sourceAssetDenom)
	crossedLimitBs := store.Get(GetCrossedLimitKey(sourceAssetHash.Bytes(), toChainId))
	crossedLimit := sdk.NewInt(0)
	if crossedLimitBs != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(crossedLimitBs, &crossedLimit)
	}
	return crossedLimit
}

func (k Keeper) BindAssetHash(ctx sdk.Context, sourceAssetDenom string, targetChainId uint64, targetAssetHash []byte, limit sdk.Int, isTargetChainAsset bool) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := DenomToHash(sourceAssetDenom)
	// store the target asset hash
	store.Set(GetBindAssetKey(sourceAssetHash.Bytes(), targetChainId), targetAssetHash)
	// make sure the new limit is greater than the stored limit
	storedCrossedLimit := k.GetCrossedLimit(ctx, sourceAssetDenom, targetChainId)
	if limit.BigInt().Cmp(storedCrossedLimit.BigInt()) != 1 {
		return sdk.ErrInternal(fmt.Sprintf("new Limit:%s should be greater than stored Limit:%s", limit.String(), storedCrossedLimit.String()))
	}
	if isTargetChainAsset {
		increment := limit.Sub(storedCrossedLimit)
		storedCrossedAmount := k.GetCrossedAmount(ctx, sourceAssetDenom, targetChainId)

		newCrossedAmount := storedCrossedAmount.Add(increment)
		if newCrossedAmount.BigInt().Cmp(storedCrossedAmount.BigInt()) != 1 {
			return sdk.ErrInternal(fmt.Sprintf("new crossedAmount:%s is not greater than stored crossed amount:%s", newCrossedAmount.String(), storedCrossedAmount.String()))
		}
		store.Set(GetCrossedAmountKey(sourceAssetHash.Bytes(), targetChainId), k.cdc.MustMarshalBinaryLengthPrefixed(newCrossedAmount))
	}

	store.Set(GetCrossedLimitKey(sourceAssetHash.Bytes(), targetChainId), k.cdc.MustMarshalBinaryLengthPrefixed(limit))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindProxy,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(targetChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(targetAssetHash)),
		),
	})
	return nil
}

func (k Keeper) CreateCoins(ctx sdk.Context, creator sdk.AccAddress, coins sdk.Coins) sdk.Error {
	if k.GetOperator(ctx).Operator.Empty() {
		k.SetOperator(ctx, types.Operator{creator})
	}
	// create the account if it doesn't yet exist
	//if !coins.IsZero() {
	//	return sdk.ErrInternal(fmt.Sprintf("only support create coins with initial zero supply"))
	//}
	zeroSupplyCoins := make([]sdk.Coin, 0)
	for _, coin := range coins {
		zeroSupplyCoins = append(zeroSupplyCoins, sdk.NewCoin(coin.Denom, sdk.NewInt(0)))
	}
	oldSupplyCoins := k.supplyKeeper.GetSupply(ctx).GetTotal()
	newSupplyCoinsToBeAdd := supply.NewSupply(oldSupplyCoins.Add(sdk.NewCoins(zeroSupplyCoins...)))

	k.supplyKeeper.SetSupply(ctx, newSupplyCoinsToBeAdd)

	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("supplyKeeper mint coins failed "))
	}
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("minted %s from %s module account", coins.String(), types.ModuleName))
	return nil
}


func (k Keeper) Lock(ctx sdk.Context, fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddressBs []byte, value sdk.Int) sdk.Error {
	// burn coin of sourceAssetDenom
	amt := sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, value))
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, fromAddress, types.ModuleName, amt); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("send coins:%s from account:%s to Module account:%s error:%v", amt.String(), fromAddress.String(), types.ModuleName, err))
	}
	//if err := k.supplyKeeper.BurnCoins(ctx, types.ModuleName, amt); err != nil {
	//	return sdk.ErrInternal(fmt.Sprintf("burn coins:%s from module account:%s error:%v", amt.String(), types.ModuleName, err))
	//}

	// make sure new crossed amount is strictly greater than old crossed amount and no less than the limit
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := DenomToHash(sourceAssetDenom)
	storedCrossedAmount := k.GetCrossedAmount(ctx, sourceAssetDenom, toChainId)
	storedCrossedLimit := k.GetCrossedLimit(ctx, sourceAssetDenom, toChainId)
	newCrossedAmount := storedCrossedAmount.Add(value)

	if newCrossedAmount.GTE(storedCrossedLimit) {
		return sdk.ErrInternal(fmt.Sprintf("new crossed amount:%s should be greater than crossed limit:%s ", newCrossedAmount.String(), storedCrossedLimit.String()))
	}

	// increase the new crossed amount by value
	store.Set(GetCrossedAmountKey(sourceAssetHash.Bytes(), toChainId), k.cdc.MustMarshalBinaryLengthPrefixed(newCrossedAmount))
	// get target chain proxy hash from storage
	toChainProxyHash := store.Get(GetBindProxyKey(toChainId))
	// get target asset hash from storage
	toChainAssetHash := store.Get(GetBindAssetKey(sourceAssetHash.Bytes(), toChainId))

	//  CreateCrossChainTx
	args := TxArgs{
		ToAssetHash: toChainAssetHash,
		ToAddress:   toAddressBs,
		Amount:      value.BigInt(),
	}
	sink := mcc.NewZeroCopySink(nil)
	if err := args.Serialization(sink); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("TxArgs Serialization error:%v", err))
	}
	if err := k.createCrossChainTx(ctx, toChainId, toChainProxyHash, "unlock", sink.Bytes()); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("create cross chain tx error:%v", err))
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AtttributeKeyStatus, strconv.FormatUint(1, 10)),
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainProxyHash, hex.EncodeToString(toChainProxyHash)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(toAddressBs)),
			sdk.NewAttribute(types.AttributeKeyFromAddress, fromAddress.String()),
			sdk.NewAttribute(types.AttributeKeyToAddress, hex.EncodeToString(toAddressBs)),
			sdk.NewAttribute(types.AttributeKeyAmount, value.String()),
		),
	})

	return nil
}

func (k Keeper) createCrossChainTx(ctx sdk.Context, toChainId uint64, toContractHash []byte, method string, args []byte) sdk.Error {
	crossChainId, err := k.getCrossChainId(ctx)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("get cross chain id from storage error:%v", err))
	}
	k.setCrossChainId(ctx, crossChainId.Add(sdk.NewInt(1)))

	var ttx ttype.Tx
	copy(ttx, ctx.TxBytes())
	txHash := ttx.Hash()
	crossChainIdBs := crossChainId.BigInt().Bytes()
	txParam := ccmc.MakeTxParam{
		TxHash:              txHash,
		CrossChainID:        crossChainIdBs,
		FromContractAddress: k.supplyKeeper.GetModuleAddress(types.ModuleName),
		ToChainID:           toChainId,
		ToContractAddress:   toContractHash,
		Method:              method,
		Args:                args,
	}
	sink := mcc.NewZeroCopySink(nil)
	txParam.Serialization(sink)

	store := ctx.KVStore(k.storeKey)

	txParamHash := tmhash.Sum(sink.Bytes())
	store.Set(GetCrossChainTxKey(txParamHash), sink.Bytes())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCrossChainTx,
			sdk.NewAttribute(types.AttributeCrossChainId, crossChainId.String()),
			sdk.NewAttribute(types.AttributeKeyTxParamHash, hex.EncodeToString(txParamHash)),
			sdk.NewAttribute(types.AttributeKeyMakeTxParam, hex.EncodeToString(sink.Bytes())),
		),
	})
	return nil
}

//
func (k Keeper) getCrossChainId(ctx sdk.Context) (sdk.Int, error) {
	store := ctx.KVStore(k.storeKey)
	idBs := store.Get(CrossChainIdKey)

	if idBs == nil {
		return sdk.NewInt(0), nil
	}
	crossChainId := new(sdk.Int)
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(idBs, crossChainId); err != nil {
		return sdk.NewInt(0), err
	}

	return *crossChainId, nil
}
func (k Keeper) setCrossChainId(ctx sdk.Context, crossChainId sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(CrossChainIdKey, k.cdc.MustMarshalBinaryLengthPrefixed(crossChainId))
}

func (k Keeper) ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, height uint32, proofStr string, headerBs []byte) sdk.Error {
	storedHeader, sdkErr := k.hsKeeper.GetHeaderByHeight(ctx, fromChainId, height)
	if sdkErr != nil {
		return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx error:%v", sdkErr))
	}
	if storedHeader == nil {
		header := new(mctype.Header)
		if err := header.Deserialization(mcc.NewZeroCopySource(headerBs)); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx error:%v", err))
		}
		if err := k.hsKeeper.ProcessHeader(ctx, header); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx error:%v", err))
		}
		storedHeader = header

	}

	proof, err := hex.DecodeString(proofStr)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx decode proof string to byte error:%v", err))
	}

	merkleValue, err := k.VerifyToCosmosTx(ctx, proof, fromChainId, storedHeader)

	if merkleValue.MakeTxParam.ToChainID != types.CurrentChainCrossChainId {
		return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx toChainId is not for this chain"))
	}

	if !bytes.Equal(merkleValue.MakeTxParam.ToContractAddress, k.supplyKeeper.GetModuleAddress(types.ModuleName)) {
		return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx, merkleValue.MakeTxParam.ToContractAddress:%s is not the lockproxy module account address:%s",
			hex.EncodeToString(merkleValue.MakeTxParam.ToContractAddress),
			hex.EncodeToString(k.supplyKeeper.GetModuleAddress(types.ModuleName).Bytes())))
	}

	if err := k.unlock(ctx, fromChainId, merkleValue.MakeTxParam.FromContractAddress, merkleValue.MakeTxParam.Args); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx, unlock errror:%v", err))
	}

	return nil
}

func (k Keeper) VerifyToCosmosTx(ctx sdk.Context, proof []byte, fromChainId uint64, header *mctype.Header) (*ccmc.ToMerkleValue, sdk.Error) {
	value, err := merkle.MerkleProve(proof, header.CrossStateRoot[:])
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("VerifyToCosmosTx, merkle.MerkleProve veify error:%v", err))
	}

	merkleValue := new(ccmc.ToMerkleValue)
	if err := merkleValue.Deserialization(mcc.NewZeroCopySource(value)); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("VerifyToCosmosTx, ToMerkleValue Deserialization error:%v", err))
	}

	if err := k.checkDoneTx(ctx, fromChainId, merkleValue.MakeTxParam.CrossChainID); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("VerifyToCosmosTx, error:%v", err))
	}

	k.putDoneTx(ctx, fromChainId, merkleValue.MakeTxParam.CrossChainID)
	//if err != nil {
	//	return nil, sdk.ErrInternal(fmt.Sprintf("VerifyToCosmosTx, error:%v", err))
	//}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeVerifyToCosmosProof,
			sdk.NewAttribute(types.AttributeKeyMerkleValueTxHash, hex.EncodeToString(merkleValue.TxHash)),
			sdk.NewAttribute(types.AttributeKeyMerkleValueMakeTxParamTxHash, hex.EncodeToString(merkleValue.MakeTxParam.TxHash)),
			sdk.NewAttribute(types.AttributeKeyFromChainId, strconv.FormatUint(fromChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyMerkleValueMakeTxParamToContractAddress, hex.EncodeToString(merkleValue.MakeTxParam.ToContractAddress)),
		),
	})
	return merkleValue, nil

}

func (k Keeper) checkDoneTx(ctx sdk.Context, fromChainId uint64, crossChainId []byte) error {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(GetDoneTxKey(fromChainId, crossChainId))
	if value != nil {
		return fmt.Errorf("checkDoneTx, tx already done")
	}
	return nil
}
func (k Keeper) putDoneTx(ctx sdk.Context, fromChainId uint64, crossChainId []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetDoneTxKey(fromChainId, crossChainId), crossChainId)
}

func (k Keeper) unlock(ctx sdk.Context, fromChainId uint64, fromContractAddress []byte, argsBs []byte) sdk.Error {
	args := new(TxArgs)
	if err := args.Deserialization(mcc.NewZeroCopySource(argsBs)); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("unlock, error:%v", err))
	}

	proxyHash := k.GetProxyHash(ctx, fromChainId)
	if len(proxyHash) == 0 {
		return sdk.ErrInternal(fmt.Sprintf("the proxyHash is empty with chainId = %d", fromChainId))
	}
	if !bytes.Equal(proxyHash, fromContractAddress) {
		return sdk.ErrInternal(fmt.Sprintf("stored proxyHash is not equal to fromContractAddress, expect:%s, got:%s", hex.EncodeToString(proxyHash), hex.EncodeToString(fromContractAddress)))
	}
	// to asset hash should be the hex format string of source asset denom name, NOT Module account address
	toAssetDenom := HashToDenom(args.ToAssetHash)

	// mint coin of sourceAssetDenom
	amt := sdk.NewCoins(sdk.NewCoin(toAssetDenom, sdk.NewIntFromBigInt(args.Amount)))
	//if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, amt); err != nil {
	//	return sdk.ErrInternal(fmt.Sprintf("mint coins:%s to module account:%s error:%v", amt.String(), types.ModuleName, err))
	//}
	var toAddress sdk.AccAddress
	copy(toAddress, args.ToAddress)
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddress, amt); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("send coins:%s from module account:%s to receiver account:%s error:%v", amt.String(), types.ModuleName, toAddress.String(), err))
	}

	// update crossedAmount value
	crossedAmount := k.GetCrossedAmount(ctx, toAssetDenom, fromChainId)

	newCrossedAmount := crossedAmount.Sub(sdk.NewIntFromBigInt(args.Amount))
	if newCrossedAmount.GTE(crossedAmount) {
		return sdk.ErrInternal(fmt.Sprintf("new crossed amount:%s should be less than old crossed amount:%s", newCrossedAmount.String(), crossedAmount.String()))
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(GetCrossedAmountKey(DenomToHash(toAssetDenom), fromChainId), k.cdc.MustMarshalBinaryLengthPrefixed(newCrossedAmount))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(types.AttributeKeyFromChainId, strconv.FormatUint(fromChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyFromContractHash, hex.EncodeToString(fromContractAddress)),
			sdk.NewAttribute(types.AttributeKeyToAssetDenom, toAssetDenom),
			sdk.NewAttribute(types.AttributeKeyToAddress, toAddress.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, args.Amount.String()),
		),
	})
	return nil
}
