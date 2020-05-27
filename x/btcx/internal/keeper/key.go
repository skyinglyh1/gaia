package keeper

import (
	"encoding/binary"

	"github.com/cosmos/gaia/x/btcx/internal/types"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = types.ModuleName
)

var (
	DenomToOperatorPrefix          = []byte{0x01}
	ChainIdToAssetHashPrefix       = []byte{0x02}
	DenomToScriptHashPrefix        = []byte{0x03}
	ScriptHashToRedeemScriptPrefix = []byte{0x04}

	ScriptHashToDenomPrefix = []byte{0x05}
	DenomToHashPrefix       = []byte{0x06}
)

func GetDenomToOperatorKey(denom string) []byte {
	return append(DenomToOperatorPrefix, []byte(denom)...)
}

func GetScriptHashAndChainIdToAssetHashKey(scriptHash []byte, chainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(append(ChainIdToAssetHashPrefix, scriptHash...), b...)
}

func GetScriptHashToDenomKey(scriptHash []byte) []byte {
	return append(ScriptHashToDenomPrefix, scriptHash...)
}

func GetScriptHashToRedeemScript(scriptHashKeyBs []byte) []byte {
	return append(ScriptHashToRedeemScriptPrefix, scriptHashKeyBs...)
}

func GetDenomToScriptHashKey(denom string) []byte {
	return append(DenomToScriptHashPrefix, []byte(denom)...)
}

func GetDenomToAssetHashKey(denom string) []byte {
	return append(DenomToHashPrefix, []byte(denom)...)
}
