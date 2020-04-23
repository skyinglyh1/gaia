package types

import (
	"crypto/elliptic"
	"github.com/cosmos/cosmos-sdk/codec"
	mctype "github.com/ontio/multi-chain/core/types"
	ontec "github.com/ontio/ontology-crypto/ec"
	ontkeypair "github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-crypto/signature"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSyncGenesisParam{}, "gaia/MsgSyncGenesisParam", nil)
	cdc.RegisterConcrete(MsgSyncHeadersParam{}, "gaia/MsgSyncHeadersParam", nil)
	cdc.RegisterConcrete(ConsensusPeers{}, "gaia/ConsensusPeers", nil)
	cdc.RegisterConcrete(mctype.Header{}, "gaiad/mcHeader", nil)
	cdc.RegisterInterface((*ontkeypair.PublicKey)(nil), nil)
	cdc.RegisterConcrete((*ontec.PublicKey)(nil), "ontecpub", nil)
	cdc.RegisterConcrete((*signature.DSASignature)(nil), "ont.DSASignature", nil)
	cdc.RegisterInterface((*elliptic.Curve)(nil), nil)
}

// module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
