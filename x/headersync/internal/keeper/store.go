package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/headersync/internal/types"
	mcc "github.com/ontio/multi-chain/common"
	"github.com/ontio/multi-chain/consensus/vbft/config"
	mcsig "github.com/ontio/multi-chain/core/signature"
	mctype "github.com/ontio/multi-chain/core/types"
	"sort"
	"strconv"
)

func (keeper BaseKeeper) SetBlockHeader(ctx sdk.Context, blockHeader *mctype.Header) sdk.Error {

	store := ctx.KVStore(keeper.storeKey)
	blockHash := blockHeader.Hash()
	sink := mcc.NewZeroCopySink(nil)
	if err := blockHeader.Serialization(sink); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("GenesisHeader Serialization err:%v", err))
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
func (keeper BaseKeeper) GetCurrentHeight(ctx sdk.Context, chainId uint64) uint32 {
	store := ctx.KVStore(keeper.storeKey)
	heightBs := store.Get(GetBlockCurHeightKey(chainId))
	if heightBs == nil {
		return 0
	}
	var height uint32
	types.ModuleCdc.MustUnmarshalJSON(heightBs, &height)
	return height

}

func (keeper BaseKeeper) GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*mctype.Header, sdk.Error) {
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
		return nil, sdk.ErrInternal(fmt.Sprintf("Header Serialization err:%v", err))
	}
	return header, nil

}
func (keeper BaseKeeper) GetHeaderByHash(ctx sdk.Context, chainId uint64, hash mcc.Uint256) (*mctype.Header, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	headerBytes := store.Get(GetBlockHeaderKey(chainId, hash.ToArray()))
	if headerBytes == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("get block header error: chainid = %d, hash = %s", chainId, hex.EncodeToString(hash.ToArray())))
	}
	header := new(mctype.Header)
	source := mcc.NewZeroCopySource(headerBytes)
	if err := header.Deserialization(source); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Header Serialization err:%v", err))
	}
	return header, nil

}

func (keeper BaseKeeper) UpdateConsensusPeer(ctx sdk.Context, blockHeader *mctype.Header) sdk.Error {

	blkInfo := &vconfig.VbftBlockInfo{}
	if err := json.Unmarshal(blockHeader.ConsensusPayload, blkInfo); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("updateConsensusPeer, unmarshal blockInfo error: %s", err))
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
			return sdk.ErrInternal(fmt.Sprintf("updateConsensusPeer, set ConsensusPeer eerror: %s", err))
		}
	}

	return nil
}

func (keeper BaseKeeper) SetConsensusPeers(ctx sdk.Context, consensusPeers *types.ConsensusPeers) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)

	//store.Set(GetConsensusPeerKey(consensusPeers.ChainID, consensusPeers.Height), types.ModuleCdc.MustMarshalBinaryLengthPrefixed(consensusPeers))

	//bz, err := types.ModuleCdc.MarshalBinaryBare(consensusPeers)
	//bz, err := types.ModuleCdc.MarshalBinaryBare(consensusPeers)
	bz, err := types.ModuleCdc.MarshalJSON(consensusPeers)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("types.ModuleCdc.MarshalBinaryBare(ConsensusPeers) error:%v", err))
	}
	store.Set(GetConsensusPeerKey(consensusPeers.ChainID, consensusPeers.Height), bz)

	// update key heights
	keyHeights, err := keeper.GetKeyHeights(ctx, consensusPeers.ChainID)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("GetKeyHeights error:%v", err))
	}

	keyHeights.HeightList = append(keyHeights.HeightList, consensusPeers.Height)

	if err = keeper.SetKeyHeights(ctx, consensusPeers.ChainID, keyHeights); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("SetKeyHeights error:%v", err))
	}
	return nil
}

func (keeper BaseKeeper) GetConsensusPeers(ctx sdk.Context, chainId uint64, height uint32) (*types.ConsensusPeers, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)

	consensusPeerBytes := store.Get(GetConsensusPeerKey(chainId, height))
	if consensusPeerBytes == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("get consensus peers error: chainid = %d, height = %d", chainId, height))
	}
	consensusPeers := &types.ConsensusPeers{
		ChainID: chainId,
		Height:  height,
		PeerMap: make(map[string]*types.Peer),
	}
	if err := types.ModuleCdc.UnmarshalJSON(consensusPeerBytes, consensusPeers); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("types.ModuleCdc.UnmarshalJSON(consensusPeers) error:%v", err))
	}
	return consensusPeers, nil
}

func (keeper BaseKeeper) SetKeyHeights(ctx sdk.Context, chainId uint64, keyHeights *types.KeyHeights) sdk.Error {
	//first sort the list  (big -> small)
	sort.SliceStable(keyHeights.HeightList, func(i, j int) bool {
		return keyHeights.HeightList[i] > keyHeights.HeightList[j]
	})
	store := ctx.KVStore(keeper.storeKey)
	store.Set(GetKeyHeightsKey(chainId), types.ModuleCdc.MustMarshalBinaryLengthPrefixed(keyHeights))
	return nil
}

func (keeper BaseKeeper) GetKeyHeights(ctx sdk.Context, chainId uint64) (*types.KeyHeights, sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	keyHeightBytes := store.Get(GetKeyHeightsKey(chainId))

	keyHeights := &types.KeyHeights{
		HeightList: make([]uint32, 0),
	}
	if keyHeightBytes == nil {
		return keyHeights, nil
	}
	types.ModuleCdc.MustUnmarshalBinaryLengthPrefixed(keyHeightBytes, keyHeights)
	return keyHeights, nil
}

func (keeper BaseKeeper) VerifyHeader(ctx sdk.Context, header *mctype.Header) sdk.Error {
	height := header.Height
	keyHeight, err := keeper.findKeyHeight(ctx, height, header.ChainID)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("verifyHeader, findKeyHeight error:%v", err))
	}
	consensusPeer, err := keeper.GetConsensusPeers(ctx, header.ChainID, keyHeight)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("verifyHeader,%v", err))
	}
	if len(header.Bookkeepers)*3 < len(consensusPeer.PeerMap)*2 {
		return sdk.ErrInternal(fmt.Sprintf("verifyHeader, header Bookkeepers num %d must more than 2/3 consensus node num %d", len(header.Bookkeepers), len(consensusPeer.PeerMap)))
	}
	for _, bookkeeper := range header.Bookkeepers {
		pubkey := vconfig.PubkeyID(bookkeeper)
		_, present := consensusPeer.PeerMap[pubkey]
		if !present {
			return sdk.ErrInternal(fmt.Sprintf("verifyHeader, invalid pubkey error:%v", pubkey))
		}
	}
	hash := header.Hash()
	er := mcsig.VerifyMultiSignature(hash[:], header.Bookkeepers, len(header.Bookkeepers), header.SigData)
	if er != nil {
		return sdk.ErrInternal(fmt.Sprintf("verifyHeader, VerifyMultiSignature error:%s, heigh:%d", err, header.Height))
	}
	return nil

}

func (keeper BaseKeeper) findKeyHeight(ctx sdk.Context, height uint32, chainID uint64) (uint32, sdk.Error) {
	keyHeights, err := keeper.GetKeyHeights(ctx, chainID)
	if err != nil {
		return 0, sdk.ErrInternal(fmt.Sprintf("findKeyHeight, GetKeyHeights error: %v", err))
	}
	for _, v := range keyHeights.HeightList {
		if (height - v) > 0 {
			return v, nil
		}
	}
	return 0, sdk.ErrInternal(fmt.Sprintf("findKeyHeight, can not find key height with height %d", height))
}
