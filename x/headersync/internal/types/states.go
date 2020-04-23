package types

import (
	"github.com/ontio/multi-chain/common"
)

type Peer struct {
	Index      uint32
	PeerPubkey string
}

type KeyHeights struct {
	HeightList []uint32
}

type ConsensusPeers struct {
	ChainID uint64
	Height  uint32
	PeerMap map[string]*Peer
}

type Header struct {
	Version          uint32
	ChainID          uint64
	PrevBlockHash    common.Uint256
	TransactionsRoot common.Uint256
	CrossStateRoot   common.Uint256
	BlockRoot        common.Uint256
	Timestamp        uint32
	Height           uint32
	ConsensusData    uint64
	ConsensusPayload []byte
	NextBookkeeper   common.Address

	//Program *program.Program
	Bookkeepers [][]byte // refer to ontology-crypto\keypair\key.go SerializePublicKey() method
	SigData     [][]byte

	hash *common.Uint256
}
