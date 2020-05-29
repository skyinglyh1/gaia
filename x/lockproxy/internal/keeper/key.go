package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gaia/x/lockproxy/internal/types"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = types.ModuleName
)

var (
	OperatorToLockProxyKey = []byte{0x01}
	BindProxyPrefix        = []byte{0x02}
	BindAssetPrefix        = []byte{0x03}
	CrossedAmountPrefix    = []byte{0x04}
)

func GetOperatorToLockProxyKey(operator sdk.AccAddress) []byte {
	return append(OperatorToLockProxyKey, operator...)
}

func GetBindProxyKey(proxyHash []byte, toChainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, toChainId)
	return append(append(BindProxyPrefix, proxyHash...), b...)
}

func GetBindAssetHashKey(lockProxyHash []byte, sourceAssetHash []byte, targetChainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, targetChainId)
	return append(append(append(BindAssetPrefix, lockProxyHash...), sourceAssetHash...), b...)
}

func GetCrossedAmountKey(sourceAssetHash []byte) []byte {
	return append(CrossedAmountPrefix, sourceAssetHash...)
}
