package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	ttype "github.com/tendermint/tendermint/types"

	"bytes"
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/gaia/x/crosschain/internal/types"
	mcc "github.com/ontio/multi-chain/common"
	mctype "github.com/ontio/multi-chain/core/types"
	"github.com/ontio/multi-chain/merkle"
	ccmc "github.com/ontio/multi-chain/native/service/cross_chain_manager/common"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"strconv"
)

type LockProxyKeeper interface {

	SetOperator(ctx sdk.Context, operator types.Operator)
	BindProxyHash(ctx sdk.Context, targetChainId uint64, targetProxyHash []byte)
	BindAssetHash(ctx sdk.Context, sourceAssetDenom string, targetChainId uint64, targetAssetHash []byte, limit sdk.Int, isTargetChainAsset bool) sdk.Error

	CreateCoins(ctx sdk.Context, creator sdk.AccAddress, coins sdk.Coins) sdk.Error
	Lock(ctx sdk.Context, fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddressBs []byte, value sdk.Int) sdk.Error
	ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, height uint32, proofStr string, headerBs []byte) sdk.Error
	LockProxyViewKeeper
}
type LockProxyViewKeeper interface {
	GetOperator(ctx sdk.Context) (operator types.Operator)
	GetProxyHash(ctx sdk.Context, toChainId uint64) []byte
	GetAssetHash(ctx sdk.Context, sourceAssetDenom string, toChainId uint64) []byte
	GetCrossedAmount(ctx sdk.Context, sourceAssetDenom string, toChainId uint64) sdk.Int
	GetCrossedLimit(ctx sdk.Context, sourceAssetDenom string, toChainId uint64) sdk.Int
}

// Logger returns a module-specific logger.
func (k CrossChainKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// get the minter
func (k CrossChainKeeper) GetOperator(ctx sdk.Context) (operator types.Operator) {
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
func (k CrossChainKeeper) SetOperator(ctx sdk.Context, operator types.Operator) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(operator)
	store.Set(types.OperatorKey, b)
}



func (k CrossChainKeeper) BindProxyHash(ctx sdk.Context, targetChainId uint64, targetProxyHash []byte) {
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

func (k CrossChainKeeper) GetProxyHash(ctx sdk.Context, toChainId uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(GetBindProxyKey(toChainId))
}

func (k CrossChainKeeper) GetAssetHash(ctx sdk.Context, sourceAssetDenom string, toChainId uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := types.DenomToHash(sourceAssetDenom)
	return store.Get(GetBindAssetKey(sourceAssetHash.Bytes(), toChainId))
}

func (k CrossChainKeeper) GetCrossedAmount(ctx sdk.Context, sourceAssetDenom string, toChainId uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := types.DenomToHash(sourceAssetDenom)
	crossedAmountBs := store.Get(GetCrossedAmountKey(sourceAssetHash.Bytes(), toChainId))
	crossedAmount := sdk.NewInt(0)
	if crossedAmountBs != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(crossedAmountBs, &crossedAmount)
	}
	return crossedAmount
}


func (k CrossChainKeeper) GetCrossedLimit(ctx sdk.Context, sourceAssetDenom string, toChainId uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := types.DenomToHash(sourceAssetDenom)
	crossedLimitBs := store.Get(GetCrossedLimitKey(sourceAssetHash.Bytes(), toChainId))
	crossedLimit := sdk.NewInt(0)
	if crossedLimitBs != nil {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(crossedLimitBs, &crossedLimit)
	}
	return crossedLimit
}

func (k CrossChainKeeper) BindAssetHash(ctx sdk.Context, sourceAssetDenom string, targetChainId uint64, targetAssetHash []byte, limit sdk.Int, isTargetChainAsset bool) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := types.DenomToHash(sourceAssetDenom)
	// store the target asset hash
	store.Set(GetBindAssetKey(sourceAssetHash.Bytes(), targetChainId), targetAssetHash)
	// make sure the new limit is greater than the stored limit
	storedCrossedLimit := k.GetCrossedLimit(ctx, sourceAssetDenom, targetChainId)
	if limit.BigInt().Cmp(storedCrossedLimit.BigInt()) != 1 {
		return types.ErrBelowCrossedLimit(types.DefaultCodespace, limit, storedCrossedLimit)
	}
	if isTargetChainAsset {
		increment := limit.Sub(storedCrossedLimit)
		storedCrossedAmount := k.GetCrossedAmount(ctx, sourceAssetDenom, targetChainId)

		newCrossedAmount := storedCrossedAmount.Add(increment)
		if newCrossedAmount.BigInt().Cmp(storedCrossedAmount.BigInt()) != 1 {
			return types.ErrCrossedAmountOverflow(types.DefaultCodespace, newCrossedAmount, storedCrossedLimit)
		}
		store.Set(GetCrossedAmountKey(sourceAssetHash.Bytes(), targetChainId), k.cdc.MustMarshalBinaryLengthPrefixed(newCrossedAmount))
	}

	store.Set(GetCrossedLimitKey(sourceAssetHash.Bytes(), targetChainId), k.cdc.MustMarshalBinaryLengthPrefixed(limit))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindAsset,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeyFromAssetHash, hex.EncodeToString(sourceAssetHash)),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(targetChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(targetAssetHash)),
		),
	})
	return nil
}

func (k CrossChainKeeper) CreateCoins(ctx sdk.Context, creator sdk.AccAddress, coins sdk.Coins) sdk.Error {
	if k.GetOperator(ctx).Operator.Empty() {
		k.SetOperator(ctx, types.Operator{creator})
	}
	if !coins.IsAllPositive() {
		return types.ErrCreateNegativeCoins(types.DefaultCodespace, coins)
	}
	//zeroSupplyCoins := make([]sdk.Coin, 0)
	//for _, coin := range coins {
	//	zeroSupplyCoins = append(zeroSupplyCoins, sdk.NewCoin(coin.Denom, sdk.NewInt(0)))
	//}

	storedSupplyCoins := k.supplyKeeper.GetSupply(ctx).GetTotal()
	for _, coin := range coins {
		for i, oldCoin := range storedSupplyCoins {
			if coin.Denom == oldCoin.Denom {
				storedSupplyCoins[i].Amount = coin.Amount
				break
			}
		}
	}
	newCoinsSupply := supply.NewSupply(storedSupplyCoins)
	k.supplyKeeper.SetSupply(ctx, newCoinsSupply)
	// TODO: should only mint the increment coins to module account
	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return types.ErrSupplyKeeperMintCoinsFail(types.DefaultCodespace)
	}
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("minted %s from %s module account", coins.String(), types.ModuleName))
	return nil
}


func (k CrossChainKeeper) Lock(ctx sdk.Context, fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddressBs []byte, value sdk.Int) sdk.Error {
	// burn coin of sourceAssetDenom
	amt := sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, value))
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, fromAddress, types.ModuleName, amt); err != nil {
		return types.ErrSendCoinsToModuleFail(types.DefaultCodespace, amt, fromAddress,  k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetAddress())
	}
	// make sure new crossed amount is strictly greater than old crossed amount and no less than the limit
	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := types.DenomToHash(sourceAssetDenom)
	storedCrossedAmount := k.GetCrossedAmount(ctx, sourceAssetDenom, toChainId)
	storedCrossedLimit := k.GetCrossedLimit(ctx, sourceAssetDenom, toChainId)
	newCrossedAmount := storedCrossedAmount.Add(value)

	if newCrossedAmount.GTE(storedCrossedLimit) {
		return types.ErrCrossedAmountOverLimit(types.DefaultCodespace, newCrossedAmount, storedCrossedLimit)
	}
	// increase the new crossed amount by value
	store.Set(GetCrossedAmountKey(sourceAssetHash.Bytes(), toChainId), k.cdc.MustMarshalBinaryLengthPrefixed(newCrossedAmount))
	// get target chain proxy hash from storage
	toChainProxyHash := store.Get(GetBindProxyKey(toChainId))
	// get target asset hash from storage
	toChainAssetHash := store.Get(GetBindAssetKey(sourceAssetHash.Bytes(), toChainId))

	//  CreateCrossChainTx
	args := types.TxArgs{
		ToAssetHash: toChainAssetHash,
		ToAddress:   toAddressBs,
		Amount:      value.BigInt(),
	}
	sink := mcc.NewZeroCopySink(nil)
	if err := args.Serialization(sink); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("TxArgs Serialization error:%v", err))
	}
	if err := k.createCrossChainTx(ctx, toChainId, toChainProxyHash, "unlock", sink.Bytes()); err != nil {
		return types.ErrCreateCrossChainTx(types.DefaultCodespace, err)
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

func (k CrossChainKeeper) createCrossChainTx(ctx sdk.Context, toChainId uint64, toContractHash []byte, method string, args []byte) sdk.Error {
	crossChainId, err := k.getCrossChainId(ctx)
	if err != nil {
		return err
	}
	if err := k.setCrossChainId(ctx, crossChainId.Add(sdk.NewInt(1))); err != nil {
		return err
	}

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
func (k CrossChainKeeper) getCrossChainId(ctx sdk.Context) (sdk.Int, sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	idBs := store.Get(CrossChainIdKey)
	if idBs == nil {
		return sdk.NewInt(0), nil
	}
	var crossChainId sdk.Int
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(idBs, &crossChainId); err != nil {
		return sdk.NewInt(0), types.ErrUnmarshalSpecificTypeFail(types.DefaultCodespace, crossChainId, err)
	}

	return crossChainId, nil
}
func (k CrossChainKeeper) setCrossChainId(ctx sdk.Context, crossChainId sdk.Int) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	idBs, err := k.cdc.MarshalBinaryLengthPrefixed(crossChainId)
	if err != nil {
		return types.ErrMarshalSpecificTypeFail(types.DefaultCodespace, crossChainId, err)
	}
	store.Set(CrossChainIdKey, idBs)
	return nil
}

func (k CrossChainKeeper) ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, height uint32, proofStr string, headerBs []byte) sdk.Error {
	storedHeader, err := k.GetHeaderByHeight(ctx, fromChainId, height)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx error:%s", err.Error()))
	}
	if storedHeader == nil {
		header := new(mctype.Header)
		if err := header.Deserialization(mcc.NewZeroCopySource(headerBs)); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx error:%s", types.ErrDeserializeHeader(types.DefaultCodespace, err).Error()))
		}
		if err := k.ProcessHeader(ctx, header); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx error:%s", err.Error()))
		}
		storedHeader = header

	}

	proof, e := hex.DecodeString(proofStr)
	if e != nil {
		return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx, decode proof hex string to byte error:%s", e.Error()))
	}

	merkleValue, err := k.VerifyToCosmosTx(ctx, proof, fromChainId, storedHeader)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx, error:%s", err.Error()))
	}
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

func (k CrossChainKeeper) VerifyToCosmosTx(ctx sdk.Context, proof []byte, fromChainId uint64, header *mctype.Header) (*ccmc.ToMerkleValue, sdk.Error) {
	value, err := merkle.MerkleProve(proof, header.CrossStateRoot[:])
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("VerifyToCosmosTx, merkle.MerkleProve veify error:%s", err.Error()))
	}

	merkleValue := new(ccmc.ToMerkleValue)
	if err := merkleValue.Deserialization(mcc.NewZeroCopySource(value)); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("VerifyToCosmosTx, ToMerkleValue Deserialization error:%s", err.Error()))
	}

	if err := k.checkDoneTx(ctx, fromChainId, merkleValue.MakeTxParam.CrossChainID); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("VerifyToCosmosTx, error:%s", err.Error()))
	}

	k.putDoneTx(ctx, fromChainId, merkleValue.MakeTxParam.CrossChainID)

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

func (k CrossChainKeeper) checkDoneTx(ctx sdk.Context, fromChainId uint64, crossChainId []byte) error {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(GetDoneTxKey(fromChainId, crossChainId))
	if value != nil {
		return fmt.Errorf("checkDoneTx, tx already done")
	}
	return nil
}
func (k CrossChainKeeper) putDoneTx(ctx sdk.Context, fromChainId uint64, crossChainId []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetDoneTxKey(fromChainId, crossChainId), crossChainId)
}

func (k CrossChainKeeper) unlock(ctx sdk.Context, fromChainId uint64, fromContractAddress []byte, argsBs []byte) sdk.Error {
	args := new(types.TxArgs)
	if err := args.Deserialization(mcc.NewZeroCopySource(argsBs)); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("unlock, error:%s", err))
	}

	proxyHash := k.GetProxyHash(ctx, fromChainId)
	if len(proxyHash) == 0 {
		return sdk.ErrInternal(fmt.Sprintf("the proxyHash is empty with chainId = %d", fromChainId))
	}
	if !bytes.Equal(proxyHash, fromContractAddress) {
		return sdk.ErrInternal(fmt.Sprintf("stored proxyHash is not equal to fromContractAddress, expect:%x, got:%x", proxyHash, fromContractAddress))
	}
	// to asset hash should be the hex format string of source asset denom name, NOT Module account address
	toAssetDenom := types.HashToDenom(args.ToAssetHash)

	// mint coin of sourceAssetDenom
	amt := sdk.NewCoins(sdk.NewCoin(toAssetDenom, sdk.NewIntFromBigInt(args.Amount)))
	//if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, amt); err != nil {
	//	return sdk.ErrInternal(fmt.Sprintf("mint coins:%s to module account:%s error:%v", amt.String(), types.ModuleName, err))
	//}
	toAddress := make(sdk.AccAddress, len(args.ToAddress))
	copy(toAddress, args.ToAddress)

	if err := k.EnsureAccountExist(ctx, toAddress); err != nil {
		return err
	}
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddress, amt); err != nil {
		return types.ErrSendCoinsFromModuleFail(types.DefaultCodespace, amt, k.GetModuleAccount(ctx).GetAddress(), toAddress)
	}

	// update crossedAmount value
	crossedAmount := k.GetCrossedAmount(ctx, toAssetDenom, fromChainId)

	newCrossedAmount := crossedAmount.Sub(sdk.NewIntFromBigInt(args.Amount))
	if newCrossedAmount.GTE(crossedAmount) {
		return sdk.ErrInternal(fmt.Sprintf("new crossed amount:%s should be less than old crossed amount:%s", newCrossedAmount.String(), crossedAmount.String()))
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(GetCrossedAmountKey(types.DenomToHash(toAssetDenom), fromChainId), k.cdc.MustMarshalBinaryLengthPrefixed(newCrossedAmount))
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