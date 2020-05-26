package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/gaia/x/ccm/internal/types"
	mcc "github.com/ontio/multi-chain/common"
	mctype "github.com/ontio/multi-chain/core/types"
	"github.com/ontio/multi-chain/merkle"
	ccmc "github.com/ontio/multi-chain/native/service/cross_chain_manager/common"
	"github.com/tendermint/tendermint/crypto/tmhash"
	ttype "github.com/tendermint/tendermint/types"
	"strconv"
)

type KeeperI interface {
	ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, height uint32, proofStr string, headerBs []byte) sdk.Error
	CreateCrossChainTx(ctx sdk.Context, toChainId uint64, fromContractHash, toContractHash []byte, method string, args []byte) sdk.Error
	SetDenomCreator(ctx sdk.Context, denom string, creator sdk.AccAddress)
	GetDenomCreator(ctx sdk.Context, denom string) sdk.AccAddress
}

// Keeper of the mint store
type Keeper struct {
	cdc         *codec.Codec
	storeKey    sdk.StoreKey
	paramSpace  params.Subspace
	hsKeeper    types.HeaderSyncKeeper
	ulKeeperMap map[string]types.UnlockKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, hsk types.HeaderSyncKeeper, ulKeeperMap map[string]types.UnlockKeeper) Keeper {

	return Keeper{
		cdc:         cdc,
		storeKey:    key,
		paramSpace:  paramSpace.WithKeyTable(types.ParamKeyTable()),
		hsKeeper:    hsk,
		ulKeeperMap: ulKeeperMap,
	}
}
func (k Keeper) MountUnlockKeeperMap(ulKeeperMap map[string]types.UnlockKeeper) {
	k.ulKeeperMap = ulKeeperMap
}

func (k Keeper) SetDenomCreator(ctx sdk.Context, denom string, creator sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(GetDenomToCreatorKey(denom), creator.Bytes())
}

func (k Keeper) GetDenomCreator(ctx sdk.Context, denom string) sdk.AccAddress {
	return ctx.KVStore(k.storeKey).Get(GetDenomToCreatorKey(denom))

}

func (k Keeper) CreateCrossChainTx(ctx sdk.Context, toChainId uint64, fromContractHash, toContractHash []byte, method string, args []byte) sdk.Error {
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
		FromContractAddress: fromContractHash,
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

func (k Keeper) ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, height uint32, proofStr string, headerBs []byte) sdk.Error {
	storedHeader, err := k.hsKeeper.GetHeaderByHeight(ctx, fromChainId, height)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx error:%s", err.Error()))
	}
	if storedHeader == nil {
		header := new(mctype.Header)
		if err := header.Deserialization(mcc.NewZeroCopySource(headerBs)); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx error:%s", types.ErrDeserializeHeader(types.DefaultCodespace, err).Error()))
		}
		if err := k.hsKeeper.ProcessHeader(ctx, header); err != nil {
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
	// check if tocontractAddress is lockproxy module account, if yes, invoke lockproxy.unlock(), otherwise, invoke btcx.unlock
	// TODO: invoke target module method
	for _, unlockKeeper := range k.ulKeeperMap {
		if unlockKeeper.ContainToContractAddr(ctx, merkleValue.MakeTxParam.ToContractAddress, fromChainId) {
			if err := unlockKeeper.Unlock(ctx, fromChainId, merkleValue.MakeTxParam.FromContractAddress, merkleValue.MakeTxParam.ToContractAddress, merkleValue.MakeTxParam.Args); err != nil {
				return sdk.ErrInternal(fmt.Sprintf("ProcessCrossChainTx, unlock errror:%v", err))
			}
			break
		}
	}

	return nil
}

func (k Keeper) VerifyToCosmosTx(ctx sdk.Context, proof []byte, fromChainId uint64, header *mctype.Header) (*ccmc.ToMerkleValue, sdk.Error) {
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

func (k Keeper) getCrossChainId(ctx sdk.Context) (sdk.Int, sdk.Error) {
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
func (k Keeper) setCrossChainId(ctx sdk.Context, crossChainId sdk.Int) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	idBs, err := k.cdc.MarshalBinaryLengthPrefixed(crossChainId)
	if err != nil {
		return types.ErrMarshalSpecificTypeFail(types.DefaultCodespace, crossChainId, err)
	}
	store.Set(CrossChainIdKey, idBs)
	return nil
}