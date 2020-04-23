package keeper

import (
	"encoding/binary"

	"github.com/cosmos/gaia/x/headersync/internal/types"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = types.ModuleName
)

// Keys for distribution store
// Items are stored with the following key: values
//
// - 0x00<proposalID_Bytes>: FeePol
//
// - 0x01: sdk.ConsAddress
//
// - 0x02<valAddr_Bytes>: ValidatorOutstandingRewards
//
// - 0x03<accAddr_Bytes>: sdk.AccAddress
//
// - 0x04<valAddr_Bytes><accAddr_Bytes>: DelegatorStartingInfo
//
// - 0x05<valAddr_Bytes><period_Bytes>: ValidatorHistoricalRewards
//
// - 0x06<valAddr_Bytes>: ValidatorCurrentRewards
//
// - 0x07<valAddr_Bytes>: ValidatorCurrentRewards
//
// - 0x08<valAddr_Bytes><height>: ValidatorSlashEvent
var (

	BlockHeaderPrefix = []byte{0x01}
	BlockHashPrefix = []byte{0x02}
	BlockCurrentHeightKey = []byte("currentHeight")
	ConsensusPeerPrefix = []byte{0x03}
	KeyHeightPrefix = []byte{0x04}

)



func GetBlockHeaderKey(chainId uint64, blockHash []byte) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(append(BlockHeaderPrefix, b...), blockHash...)
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

