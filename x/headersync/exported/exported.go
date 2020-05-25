package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mcc "github.com/ontio/multi-chain/common"
	mctype "github.com/ontio/multi-chain/core/types"
)

// DelegationI delegation bond for a delegated proof of stake system
type HeaderSyncKeeper interface {
	ProcessHeader(ctx sdk.Context, header *mctype.Header) sdk.Error
	GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*mctype.Header, sdk.Error)
	GetHeaderByHash(ctx sdk.Context, chainId uint64, hash mcc.Uint256) (*mctype.Header, sdk.Error)
	GetCurrentHeight(ctx sdk.Context, chainId uint64) (uint32, sdk.Error)
}
