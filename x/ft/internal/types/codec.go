package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {

	cdc.RegisterConcrete(MsgCreateDenom{}, ModuleName+"/MsgCreateDenom", nil)
	cdc.RegisterConcrete(MsgCreateCoinAndDelegateToProxy{}, ModuleName+"/MsgCreateCoinAndDelegateToProxy", nil)
	cdc.RegisterConcrete(MsgBindAssetHash{}, ModuleName+"/MsgBindAssetHash", nil)
	cdc.RegisterConcrete(MsgLock{}, ModuleName+"/MsgLock", nil)
	cdc.RegisterConcrete(MsgCreateCoins{}, ModuleName+"/MsgCreateCoins", nil)

}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
