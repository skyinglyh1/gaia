package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/crosschain/internal/types"
	mcc "github.com/ontio/multi-chain/common"
	"github.com/ontio/multi-chain/consensus/vbft/config"
	mcsig "github.com/ontio/multi-chain/core/signature"
	mctype "github.com/ontio/multi-chain/core/types"
	"sort"
	"strconv"
)

type HeaderSyncKeeper interface {

	SyncGenesisHeader(ctx sdk.Context, genesisHeader []byte) sdk.Error
	SyncBlockHeaders(ctx sdk.Context, headers [][]byte) sdk.Error
	ProcessHeader(ctx sdk.Context, header *mctype.Header) sdk.Error
	HeaderSyncViewKeeper
}


func (keeper CrossChainKeeper) SyncGenesisHeader(ctx sdk.Context, genesisHeaderBytes []byte) sdk.Error {
	genesisHeader := &mctype.Header{}

	source := mcc.NewZeroCopySource(genesisHeaderBytes)
	if err := genesisHeader.Deserialization(source); err != nil {
		return types.ErrDeserializeHeader(types.DefaultCodespace, err)
	}

	if err := keeper.SetBlockHeader(ctx, genesisHeader); err != nil {
		return err
	}
	if err := keeper.UpdateConsensusPeer(ctx, genesisHeader); err != nil {
		return err
	}

	return nil
}

func (keeper CrossChainKeeper) SyncBlockHeaders(ctx sdk.Context, headers [][]byte) sdk.Error {
	for _, headerBytes := range headers {
		header := &mctype.Header{}
		source := mcc.NewZeroCopySource(headerBytes)
		if err := header.Deserialization(source); err != nil {
			return types.ErrDeserializeHeader(types.DefaultCodespace, err)
		}
		h, err := keeper.GetHeaderByHeight(ctx, header.ChainID, header.Height)
		if err != nil {
			return sdk.ErrInternal(fmt.Sprintf("SyncBlockHeader chainId=%d, height=%d, err:%s", header.ChainID, header.Height, err.Error()))
		}

		if h == nil {
			if err := keeper.ProcessHeader(ctx, header); err != nil {
				return sdk.ErrInternal(fmt.Sprintf("SyncBlockHeader error:%s", err.Error()))
			}
		}
	}
	return nil
}

func (keeper CrossChainKeeper) ProcessHeader(ctx sdk.Context, header *mctype.Header) sdk.Error {
	if err := keeper.VerifyHeader(ctx, header); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("processHeader, %s", err.Error()))
	}
	if err := keeper.SetBlockHeader(ctx, header); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("processHeader, %s", err.Error()))
	}
	if err := keeper.UpdateConsensusPeer(ctx, header); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("processHeader, %s", err.Error()))
	}
	return nil
}

type HeaderSyncViewKeeper interface {
	GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*mctype.Header, sdk.Error)
	GetHeaderByHash(ctx sdk.Context, chainId uint64, hash mcc.Uint256) (*mctype.Header, sdk.Error)
	GetCurrentHeight(ctx sdk.Context, chainId uint64) (uint32, sdk.Error)
	GetConsensusPeers(ctx sdk.Context, chainId uint64, height uint32) (*types.ConsensusPeers, sdk.Error)
	GetKeyHeights(ctx sdk.Context, chainId uint64) *types.KeyHeights
}


func (keeper CrossChainKeeper) SetBlockHeader(ctx sdk.Context, blockHeader *mctype.Header) sdk.Error {

	store := ctx.KVStore(keeper.storeKey)
	blockHash := blockHeader.Hash()
	sink := mcc.NewZeroCopySink(nil)
	if err := blockHeader.Serialization(sink); err != nil {
		return types.ErrDeserializeHeader(types.DefaultCodespace, err)
	}
	store.Set(GetBlockHeaderKey(blockHeader.ChainID, blockHash.ToArray()), sink.Bytes())
	store.Set(GetBlockHashKey(blockHeader.ChainID, blockHeader.Height), types.ModuleCdc.MustMarshalJSON(blockHash))
	store.Set(GetBlockCurHeightKey(blockHeader.ChainID), types.ModuleCdc.MustMarshalJSON(blockHeader.Height))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSyncHeader,
			sdk.NewAttribute(types.AttributeKeyChainId, strconv.FormatUint(blockHeader.ChainID, 10)),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatUint(uint64(blockHeader.Height), 10)),
			sdk.NewAttribute(types.AttributeKeyBlockHash, hex.EncodeToString(blockHash[:])),
			sdk.NewAttribute(types.AttributeKeyNativeChainHeight, strconv.FormatUint(uint64(ctx.BlockHeight()), 10)),
		),
	})
	return nil
}
func (keeper CrossChainKeeper) GetCurrentHeight(ctx sdk.Context, chainId uint64) (uint32, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	heightBs := store.Get(GetBlockCurHeightKey(chainId))
	if heightBs == nil {
		return 0, nil
	}
	var height uint32
	if err := types.ModuleCdc.UnmarshalJSON(heightBs, &height); err != nil {
		return 0, types.ErrUnmarshalSpecificTypeFail(types.DefaultCodespace, height, err)
	}
	return height, nil

}

func (keeper CrossChainKeeper) GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*mctype.Header, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	hashBytes := store.Get(GetBlockHashKey(chainId, height))
	if hashBytes == nil {
		return nil, nil
	}
	blockHash := new(mcc.Uint256)
	types.ModuleCdc.MustUnmarshalJSON(hashBytes, blockHash)
	headerBytes := store.Get(GetBlockHeaderKey(chainId, blockHash.ToArray()))
	if headerBytes == nil {
		return nil, nil
	}
	header := new(mctype.Header)
	source := mcc.NewZeroCopySource(headerBytes)
	if err := header.Deserialization(source); err != nil {
		return nil, types.ErrDeserializeHeader(types.DefaultCodespace, err)
	}
	return header, nil

}
func (keeper CrossChainKeeper) GetHeaderByHash(ctx sdk.Context, chainId uint64, hash mcc.Uint256) (*mctype.Header, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	headerBytes := store.Get(GetBlockHeaderKey(chainId, hash.ToArray()))
	if headerBytes == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("get block header error: chainid = %d, hash = %s", chainId, hex.EncodeToString(hash.ToArray())))
	}
	header := new(mctype.Header)
	source := mcc.NewZeroCopySource(headerBytes)
	if err := header.Deserialization(source); err != nil {
		return nil, types.ErrDeserializeHeader(types.DefaultCodespace, err)
	}
	return header, nil

}

func (keeper CrossChainKeeper) UpdateConsensusPeer(ctx sdk.Context, blockHeader *mctype.Header) sdk.Error {

	blkInfo := &vconfig.VbftBlockInfo{}
	if err := json.Unmarshal(blockHeader.ConsensusPayload, blkInfo); err != nil {
		return types.ErrUnmarshalSpecificTypeFail(types.DefaultCodespace, blkInfo, err)
	}
	if blkInfo.NewChainConfig != nil {
		consensusPeers := &types.ConsensusPeers{
			ChainID: blockHeader.ChainID,
			Height:  blockHeader.Height,
			PeerMap: make(map[string]*types.Peer),
		}
		for _, p := range blkInfo.NewChainConfig.Peers {
			consensusPeers.PeerMap[p.ID] = &types.Peer{Index: p.Index, PeerPubkey: p.ID}
		}
		if err := keeper.SetConsensusPeers(ctx, consensusPeers); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("updateConsensusPeer, set ConsensusPeer error: %s", err.Error()))
		}
	}

	return nil
}

func (keeper CrossChainKeeper) SetConsensusPeers(ctx sdk.Context, consensusPeers *types.ConsensusPeers) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)

	bz, err := types.ModuleCdc.MarshalJSON(consensusPeers)
	if err != nil {
		return types.ErrMarshalSpecificTypeFail(types.DefaultCodespace, consensusPeers, err)
	}
	store.Set(GetConsensusPeerKey(consensusPeers.ChainID, consensusPeers.Height), bz)

	// update key heights
	keyHeights := keeper.GetKeyHeights(ctx, consensusPeers.ChainID)

	keyHeights.HeightList = append(keyHeights.HeightList, consensusPeers.Height)

	if e := keeper.SetKeyHeights(ctx, consensusPeers.ChainID, keyHeights); e != nil {
		return sdk.ErrInternal(fmt.Sprintf("SetConsensusPeers, set KeyHeights error: %s", e.Error()))
	}
	return nil
}

func (keeper CrossChainKeeper) GetConsensusPeers(ctx sdk.Context, chainId uint64, height uint32) (*types.ConsensusPeers, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)

	consensusPeerBytes := store.Get(GetConsensusPeerKey(chainId, height))
	if consensusPeerBytes == nil {
		return nil, types.ErrGetConsensusPeers(types.DefaultCodespace, height, chainId)
	}
	consensusPeers := new(types.ConsensusPeers)
	if err := types.ModuleCdc.UnmarshalJSON(consensusPeerBytes, consensusPeers); err != nil {
		return nil, types.ErrUnmarshalSpecificTypeFail(types.DefaultCodespace, consensusPeers, err)
	}
	return consensusPeers, nil
}

func (keeper CrossChainKeeper) SetKeyHeights(ctx sdk.Context, chainId uint64, keyHeights *types.KeyHeights) sdk.Error {
	//first sort the list  (big -> small)
	sort.SliceStable(keyHeights.HeightList, func(i, j int) bool {
		return keyHeights.HeightList[i] > keyHeights.HeightList[j]
	})
	store := ctx.KVStore(keeper.storeKey)
	bz, err := types.ModuleCdc.MarshalBinaryLengthPrefixed(keyHeights)
	if err != nil {
		//return types.ErrMarshalKeyHeightsFail(types.DefaultCodespace, err)
		return types.ErrMarshalSpecificTypeFail(types.DefaultCodespace, keyHeights, err)
	}
	store.Set(GetKeyHeightsKey(chainId), bz)
	return nil
}

func (keeper CrossChainKeeper) GetKeyHeights(ctx sdk.Context, chainId uint64) *types.KeyHeights {
	store := ctx.KVStore(keeper.storeKey)
	keyHeightBytes := store.Get(GetKeyHeightsKey(chainId))
	keyHeights := new(types.KeyHeights)
	if keyHeightBytes == nil {
		return keyHeights
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(keyHeightBytes, keyHeights)
	return keyHeights
}

func (keeper CrossChainKeeper) VerifyHeader(ctx sdk.Context, header *mctype.Header) sdk.Error {
	height := header.Height
	keyHeight, err := keeper.findKeyHeight(ctx, height, header.ChainID)
	if err != nil {
		return err
	}
	consensusPeer, err := keeper.GetConsensusPeers(ctx, header.ChainID, keyHeight)
	if err != nil {
		return err
	}
	if len(header.Bookkeepers)*3 < len(consensusPeer.PeerMap)*2 {
		return types.ErrBookKeeperNum(types.DefaultCodespace, len(header.Bookkeepers), len(consensusPeer.PeerMap))
	}
	for _, bookkeeper := range header.Bookkeepers {
		pubkey := vconfig.PubkeyID(bookkeeper)
		_, present := consensusPeer.PeerMap[pubkey]
		if !present {
			return types.ErrInvalidPublicKey(types.DefaultCodespace, pubkey)
		}
	}
	hash := header.Hash()
	if e := mcsig.VerifyMultiSignature(hash[:], header.Bookkeepers, len(header.Bookkeepers), header.SigData); e != nil {
		return types.ErrVerifyMultiSignatureFail(types.DefaultCodespace, err, header.Height)
	}
	return nil
}

func (keeper CrossChainKeeper) findKeyHeight(ctx sdk.Context, height uint32, chainId uint64) (uint32, sdk.Error) {
	keyHeights := keeper.GetKeyHeights(ctx, chainId)
	for _, v := range keyHeights.HeightList {
		if (height - v) > 0 {
			return v, nil
		}
	}
	return 0, types.ErrFindKeyHeight(types.DefaultCodespace, height, chainId)
}
