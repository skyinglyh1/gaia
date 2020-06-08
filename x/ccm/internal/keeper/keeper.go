package keeper

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/gaia/x/ccm/internal/types"
	polycommon "github.com/cosmos/gaia/x/headersync/poly-utils/common"
	polytype "github.com/cosmos/gaia/x/headersync/poly-utils/core/types"
	"github.com/cosmos/gaia/x/headersync/poly-utils/merkle"
	ccmc "github.com/cosmos/gaia/x/headersync/poly-utils/native/service/cross_chain_manager/common"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	ttype "github.com/tendermint/tendermint/types"
)

type KeeperI interface {
	ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, height uint32, proofStr string, headerBs []byte) error
	CreateCrossChainTx(ctx sdk.Context, toChainId uint64, fromContractHash, toContractHash []byte, method string, args []byte) error
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
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) MountUnlockKeeperMap(ulKeeperMap map[string]types.UnlockKeeper) {
	k.ulKeeperMap = make(map[string]types.UnlockKeeper)
	for key, value := range ulKeeperMap {
		k.ulKeeperMap[key] = value
	}
	//k.ulKeeperMap = ulKeeperMap
}

func (k Keeper) IfContainToContract(ctx sdk.Context, keystore string, toContractAddr []byte, fromChainId uint64) *types.QueryContainToContractRes {
	k.Logger(ctx).Info(fmt.Sprintf("k.unkeeperMap is %+v ", k.ulKeeperMap))
	unlockKeeper, ok := k.ulKeeperMap[keystore]
	if !ok {
		return &types.QueryContainToContractRes{KeyStore: keystore, Info: "map doesnot contain current keystore"}
	}
	var res types.QueryContainToContractRes
	res.KeyStore = keystore
	res.Exist = unlockKeeper.ContainToContractAddr(ctx, toContractAddr, fromChainId)
	k.Logger(ctx).Info(fmt.Sprintf("key is %+v ", keystore))
	k.Logger(ctx).Info(fmt.Sprintf("IfContains %+v ", unlockKeeper.ContainToContractAddr(ctx, toContractAddr, fromChainId)))

	return &res
}

func (k Keeper) SetDenomCreator(ctx sdk.Context, denom string, creator sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(GetDenomToCreatorKey(denom), creator.Bytes())
}

func (k Keeper) GetDenomCreator(ctx sdk.Context, denom string) sdk.AccAddress {
	return ctx.KVStore(k.storeKey).Get(GetDenomToCreatorKey(denom))
}

func (k Keeper) CreateCrossChainTx(ctx sdk.Context, toChainId uint64, fromContractHash, toContractHash []byte, method string, args []byte) error {
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
	sink := polycommon.NewZeroCopySink(nil)
	txParam.Serialization(sink)

	store := ctx.KVStore(k.storeKey)

	txParamHash := tmhash.Sum(sink.Bytes())
	store.Set(GetCrossChainTxKey(txParamHash), sink.Bytes())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCrossChainTx,
			sdk.NewAttribute(types.AttributeKeyStatus, "1"),
			sdk.NewAttribute(types.AttributeCrossChainId, crossChainId.String()),
			sdk.NewAttribute(types.AttributeKeyTxParamHash, hex.EncodeToString(txParamHash)),
			sdk.NewAttribute(types.AttributeKeyMakeTxParam, hex.EncodeToString(sink.Bytes())),
		),
	})
	return nil
}

func (k Keeper) ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, height uint32, proofStr string, headerBs []byte) error {
	storedHeader, err := k.hsKeeper.GetHeaderByHeight(ctx, fromChainId, height)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ProcessCrossChainTx error:%s", err.Error()))
	}
	if storedHeader == nil {
		header := new(polytype.Header)
		if err := header.Deserialization(polycommon.NewZeroCopySource(headerBs)); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ProcessCrossChainTx error:%s", types.ErrDeserializeHeader(err).Error()))
		}
		if err := k.hsKeeper.ProcessHeader(ctx, header); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ProcessCrossChainTx error:%s", err.Error()))
		}
		storedHeader = header

	}

	proof, e := hex.DecodeString(proofStr)
	if e != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ProcessCrossChainTx, decode proof hex string to byte error:%s", e.Error()))
	}

	merkleValue, err := k.VerifyToCosmosTx(ctx, proof, fromChainId, storedHeader)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ProcessCrossChainTx, error:%s", err.Error()))
	}
	if merkleValue.MakeTxParam.ToChainID != types.CurrentChainCrossChainId {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ProcessCrossChainTx toChainId is not for this chain"))
	}
	// check if tocontractAddress is lockproxy module account, if yes, invoke lockproxy.unlock(), otherwise, invoke btcx.unlock
	k.Logger(ctx).Info(fmt.Sprintf("k.unkeeperMap is %+v ", k.ulKeeperMap))

	// TODO: invoke target module method
	for key, unlockKeeper := range k.ulKeeperMap {
		k.Logger(ctx).Info(fmt.Sprintf("key is %+v ", key))
		k.Logger(ctx).Info(fmt.Sprintf("IfContains %+v ", unlockKeeper.ContainToContractAddr(ctx, merkleValue.MakeTxParam.ToContractAddress, fromChainId)))

		if unlockKeeper.ContainToContractAddr(ctx, merkleValue.MakeTxParam.ToContractAddress, fromChainId) {
			if err := unlockKeeper.Unlock(ctx, fromChainId, merkleValue.MakeTxParam.FromContractAddress, merkleValue.MakeTxParam.ToContractAddress, merkleValue.MakeTxParam.Args); err != nil {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("ProcessCrossChainTx, unlock errror:%v", err))
			}
			return nil
		}
	}

	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("Cannot find any unlock keeper to perform 'unlock' method for toContractAddr:%x, fromChainId:%d", merkleValue.MakeTxParam.ToContractAddress, fromChainId))
}

func (k Keeper) VerifyToCosmosTx(ctx sdk.Context, proof []byte, fromChainId uint64, header *polytype.Header) (*ccmc.ToMerkleValue, error) {
	value, err := merkle.MerkleProve(proof, header.CrossStateRoot[:])
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("VerifyToCosmosTx, merkle.MerkleProve veify error:%s", err.Error()))
	}

	merkleValue := new(ccmc.ToMerkleValue)
	if err := merkleValue.Deserialization(polycommon.NewZeroCopySource(value)); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("VerifyToCosmosTx, ToMerkleValue Deserialization error:%s", err.Error()))
	}

	if err := k.checkDoneTx(ctx, fromChainId, merkleValue.MakeTxParam.CrossChainID); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("VerifyToCosmosTx, error:%s", err.Error()))
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

func (k Keeper) getCrossChainId(ctx sdk.Context) (sdk.Int, error) {
	store := ctx.KVStore(k.storeKey)
	idBs := store.Get(CrossChainIdKey)
	if idBs == nil {
		return sdk.NewInt(0), nil
	}
	var crossChainId sdk.Int
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(idBs, &crossChainId); err != nil {
		return sdk.NewInt(0), types.ErrUnmarshalSpecificTypeFail(crossChainId, err)
	}

	return crossChainId, nil
}
func (k Keeper) setCrossChainId(ctx sdk.Context, crossChainId sdk.Int) error {
	store := ctx.KVStore(k.storeKey)
	idBs, err := k.cdc.MarshalBinaryLengthPrefixed(crossChainId)
	if err != nil {
		return types.ErrMarshalSpecificTypeFail(crossChainId, err)
	}
	store.Set(CrossChainIdKey, idBs)
	return nil
}
