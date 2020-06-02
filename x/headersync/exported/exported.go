package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	polycommon "github.com/cosmos/gaia/x/headersync/poly-utils/common"
	polytype "github.com/cosmos/gaia/x/headersync/poly-utils/core/types"
)

// DelegationI delegation bond for a delegated proof of stake system
type HeaderSyncKeeper interface {
	ProcessHeader(ctx sdk.Context, header *polytype.Header) sdk.Error
	GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*polytype.Header, sdk.Error)
	GetHeaderByHash(ctx sdk.Context, chainId uint64, hash polycommon.Uint256) (*polytype.Header, sdk.Error)
	GetCurrentHeight(ctx sdk.Context, chainId uint64) (uint32, sdk.Error)
}
