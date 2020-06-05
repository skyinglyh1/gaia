package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegationI delegation bond for a delegated proof of stake system
type UnlockKeeper interface {
	Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) error
	ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool
}
