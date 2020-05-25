package keeper

import (
	"encoding/binary"

	"github.com/cosmos/gaia/x/ccm/internal/types"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = types.ModuleName
)

var (
	CrossChainTxDetailPrefix = []byte{0x01}
	CrossChainDoneTxPrefix   = []byte{0x02}
	DenomToCreatorPrefix     = []byte{0x03}

	CrossChainIdKey = []byte("crosschainid")
)

func GetCrossChainTxKey(crossChainTxSum []byte) []byte {
	return append(CrossChainTxDetailPrefix, crossChainTxSum...)
}
func GetDoneTxKey(fromChainId uint64, crossChainid []byte) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, fromChainId)
	return append(append(CrossChainDoneTxPrefix, b...), crossChainid...)
}

func GetDenomToCreatorKey(denom string) []byte {
	return append(DenomToCreatorPrefix, []byte(denom)...)
}
