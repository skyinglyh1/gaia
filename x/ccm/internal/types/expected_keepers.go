package types // noalias

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	polycommon "github.com/cosmos/gaia/x/headersync/poly-utils/common"
	polytype "github.com/cosmos/gaia/x/headersync/poly-utils/core/types"
)

// SupplyKeeper defines the expected supply keeper
type HeaderSyncKeeper interface {
	ProcessHeader(ctx sdk.Context, header *polytype.Header) error
	GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*polytype.Header, error)
	GetHeaderByHash(ctx sdk.Context, chainId uint64, hash polycommon.Uint256) (*polytype.Header, error)
	GetCurrentHeight(ctx sdk.Context, chainId uint64) (uint32, error)
}

type SupplyI interface {
	SetTotal(total sdk.Coins) SupplyI
}

type UnlockKeeper interface {
	Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) error
	ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool
}
