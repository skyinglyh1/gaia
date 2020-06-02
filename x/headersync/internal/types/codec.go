package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	polytype "github.com/cosmos/gaia/x/headersync/poly-utils/core/types"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSyncGenesisParam{}, "crosschain/MsgSyncGenesisParam", nil)
	cdc.RegisterConcrete(MsgSyncHeadersParam{}, "crosschain/MsgSyncHeadersParam", nil)
	cdc.RegisterConcrete(ConsensusPeers{}, "crosschain/ConsensusPeers", nil)
	cdc.RegisterConcrete(polytype.Header{}, "crosschain/PolyHeader", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
