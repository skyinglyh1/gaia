package types // noalias

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mcc "github.com/ontio/multi-chain/common"
	mctype "github.com/ontio/multi-chain/core/types"
)

// SupplyKeeper defines the expected supply keeper
type HeaderSyncKeeper interface {
	ProcessHeader(ctx sdk.Context, header *mctype.Header) sdk.Error
	GetHeaderByHeight(ctx sdk.Context, chainId uint64, height uint32) (*mctype.Header, sdk.Error)
	GetHeaderByHash(ctx sdk.Context, chainId uint64, hash mcc.Uint256) (*mctype.Header, sdk.Error)
	GetCurrentHeight(ctx sdk.Context, chainId uint64) (uint32, sdk.Error)
}

type SupplyI interface {
	SetTotal(total sdk.Coins) SupplyI
}

type UnlockKeeper interface {
	Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) sdk.Error
	ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool
}
