package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {

	cdc.RegisterConcrete(MsgBindAssetHash{}, ModuleName+"/MsgBindAssetHash", nil)
	cdc.RegisterConcrete(MsgLock{}, ModuleName+"/MsgLock", nil)
	cdc.RegisterConcrete(MsgProcessCrossChainTx{}, ModuleName+"/MsgProcessCrossChainTx", nil)
	cdc.RegisterConcrete(sdk.Int{}, "sdk/Int", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
