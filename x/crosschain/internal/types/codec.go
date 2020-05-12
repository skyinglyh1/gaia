package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mctype "github.com/ontio/multi-chain/core/types"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func RegisterCodecForLockProxy(cdc *codec.Codec) {

	cdc.RegisterConcrete(MsgCreateCoins{}, "crosschain/MsgCreateCoins", nil)
	cdc.RegisterConcrete(MsgBindProxyParam{}, "crosschain/MsgBindProxyParam", nil)
	cdc.RegisterConcrete(MsgBindAssetParam{}, "crosschain/MsgBindAssetParam", nil)
	cdc.RegisterConcrete(MsgLock{}, "crosschain/MsgLock", nil)
	cdc.RegisterConcrete(MsgProcessCrossChainTx{}, "crosschain/MsgProcessCrossChainTx", nil)
	cdc.RegisterConcrete(sdk.Int{}, "sdk/Int", nil)
	cdc.RegisterConcrete(MsgBindNoVMChainAssetHash{}, "crosschain/MsgBindNoVMChainAssetHash", nil)
	cdc.RegisterConcrete(MsgSetRedeemScript{}, "crosschain/MsgSetRedeemScript", nil)
}

func RegisterCodecForHeaderSync(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSyncGenesisParam{}, "crosschain/MsgSyncGenesisParam", nil)
	cdc.RegisterConcrete(MsgSyncHeadersParam{}, "crosschain/MsgSyncHeadersParam", nil)
	cdc.RegisterConcrete(ConsensusPeers{}, "crosschain/ConsensusPeers", nil)
	cdc.RegisterConcrete(mctype.Header{}, "crosschain/mcHeader", nil)
}

func RegisterCodec(cdc *codec.Codec) {
	RegisterCodecForHeaderSync(cdc)
	RegisterCodecForLockProxy(cdc)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
