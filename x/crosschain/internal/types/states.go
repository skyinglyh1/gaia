package types

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
