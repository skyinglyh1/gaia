package keeper

import (
	"encoding/binary"

	"github.com/cosmos/gaia/x/lockproxy/internal/types"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = types.ModuleName
)

var (
	BindProxyPrefix     = []byte{0x01}
	BindAssetPrefix     = []byte{0x02}
	CrossedLimitPrefix  = []byte{0x03}
	CrossedAmountPrefix = []byte{0x04}

	CrossChainIdKey          = []byte("crosschainid")
	CrossChainTxDetailPrefix = []byte{0x05}
	CrossChainDoneTxPrefix   = []byte{0x06}
)

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

func GetCrossChainTxKey(crossChainTxSum []byte) []byte {
	return append(CrossChainTxDetailPrefix, crossChainTxSum...)
}
func GetDoneTxKey(fromChainId uint64, crossChainid []byte) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, fromChainId)
	return append(append(CrossChainDoneTxPrefix, b...), crossChainid...)
}
