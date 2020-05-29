package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryProxyByOperator struct {
	Operator sdk.AccAddress
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryProxyByOperator(operator sdk.AccAddress) QueryProxyByOperator {
	return QueryProxyByOperator{Operator: operator}
}

// QueryBalanceParams defines the params for querying an account balance.
type QueryProxyHashParam struct {
	LockProxyHash []byte
	ChainId       uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryProxyHashParams(lockProxyHash []byte, chainId uint64) QueryProxyHashParam {
	return QueryProxyHashParam{lockProxyHash, chainId}
}

type QueryAssetHashParam struct {
	LockProxyHash    []byte
	SourceAssetDenom string
	ChainId          uint64
}

func NewQueryAssetHashParams(lockProxyHash []byte, sourceAssetDenom string, chainId uint64) QueryAssetHashParam {
	return QueryAssetHashParam{LockProxyHash: lockProxyHash, SourceAssetDenom: sourceAssetDenom, ChainId: chainId}
}

type QueryLockedAmtParam struct {
	SourceAssetDenom string
}

func NewQueryLockedAmtParam(sourceAssetDenom string) QueryLockedAmtParam {
	return QueryLockedAmtParam{SourceAssetDenom: sourceAssetDenom}
}
