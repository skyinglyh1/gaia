package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegationI delegation bond for a delegated proof of stake system
type CCMKeeper interface {
	ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, height uint32, proofStr string, headerBs []byte) error
	CreateCrossChainTx(ctx sdk.Context, toChainId uint64, fromContractHash, toContractHash []byte, method string, args []byte) error
}
