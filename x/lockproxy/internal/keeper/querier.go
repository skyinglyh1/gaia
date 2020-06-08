package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
)

const (
	QueryProxyByOperator = "query_proxy_by_operator"
	QueryProxyHash       = "proxy_hash"
	QueryAssetHash       = "asset_hash"
	QueryLockedAmt       = "locked_amount"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryProxyByOperator:
			return queryProxyByOperator(ctx, req, k)
		case QueryProxyHash:
			return queryProxyHash(ctx, req, k)
		case QueryAssetHash:
			return queryAssetHash(ctx, req, k)
		case QueryLockedAmt:
			return queryLockedAmount(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}

func queryProxyByOperator(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryProxyByOperator

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	proxyHash := k.GetLockProxyByOperator(ctx, params.Operator)
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, proxyHash)
	if e != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", e.Error()))
	}

	return bz, nil
}

func queryProxyHash(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryProxyHashParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	proxyHashBs := k.GetProxyHash(ctx, params.LockProxyHash, params.ChainId)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, proxyHashBs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", err.Error()))
	}

	return bz, nil
}

func queryAssetHash(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryAssetHashParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	assetHashBs := k.GetAssetHash(ctx, params.LockProxyHash, params.SourceAssetDenom, params.ChainId)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, assetHashBs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", err.Error()))
	}

	return bz, nil
}

func queryLockedAmount(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryLockedAmtParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	crossedAmount := k.GetLockedAmount(ctx, params.SourceAssetDenom)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, crossedAmount)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", err.Error()))
	}

	return bz, nil
}
