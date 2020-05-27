package keeper

import (
	"encoding/binary"

	"github.com/cosmos/gaia/x/ft/internal/types"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = types.ModuleName
)

var (
	BindAssetHashPrefix = []byte{0x01}
	BlockHashPrefix     = []byte{0x02}
	ConsensusPeerPrefix = []byte{0x03}
	KeyHeightPrefix     = []byte{0x04}

	BindProxyPrefix             = []byte{0x05}
	BindAssetPrefix             = []byte{0x06}
	CrossedLimitPrefix          = []byte{0x07}
	CrossedAmountPrefix         = []byte{0x08}
	LockedAmountPrefix          = []byte{0x09}
	CrossChainTxDetailPrefix    = []byte{0x09}
	CrossChainDoneTxPrefix      = []byte{0xa}
	RedeemKeyScriptPrefix       = []byte{0xb}
	RedeemToHashPrefix          = []byte{0xc}
	ContractToRedeemPrefix      = []byte{0xd}
	DenomToHashPrefix           = []byte{0xe}
	HashToDenomPrefix           = []byte{0xf}
	IndependentCrossDenomPrefix = []byte{0x10}
	BlockCurrentHeightKey       = []byte("currentHeight")
	CrossChainIdKey             = []byte("crosschainid")
)

func GetBindAssetHashKey(sourceDenomHash []byte, chainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(append(BindAssetHashPrefix, sourceDenomHash...), b...)
}

func GetBlockHashKey(chainId uint64, height uint32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	bh := make([]byte, 4)
	binary.LittleEndian.PutUint32(bh, height)
	return append(append(BlockHashPrefix, b...), bh...)
}
func GetBlockCurHeightKey(chainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(BlockCurrentHeightKey, b...)
}

func GetConsensusPeerKey(chainId uint64, height uint32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	bh := make([]byte, 4)
	binary.LittleEndian.PutUint32(bh, height)
	return append(append(ConsensusPeerPrefix, b...), bh...)
}

func GetKeyHeightsKey(chainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(append(KeyHeightPrefix, b...), b...)
}

func GetBindProxyKey(targetChainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, targetChainId)
	return append(BindProxyPrefix, b...)
}

func GetBindAssetKey(sourceAssetHash []byte, targetChainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, targetChainId)
	return append(append(BindAssetPrefix, sourceAssetHash...), b...)
}

func GetCrossedLimitKey(sourceAssetHash []byte, targetChainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, targetChainId)
	return append(append(CrossedLimitPrefix, sourceAssetHash...), b...)
}

func GetCrossedAmountKey(sourceAssetHash []byte, targetChainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, targetChainId)
	return append(append(CrossedAmountPrefix, sourceAssetHash...), b...)
}

func GetIndependentCrossDenomKey(denom string) []byte {
	return append(IndependentCrossDenomPrefix, []byte(denom)...)
}
