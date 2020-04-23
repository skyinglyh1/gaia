package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
)

const (
	// query balance path
	QueryProxyHash     = "proxy_hash"
	QueryAssetHash     = "asset_hash"
	QueryCrossedAmount = "crossed_amount"
	QueryCrossedLimit  = "crossed_limit"
	QueryOperator      = "operator"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryProxyHash:
			return queryProxyHash(ctx, req, k)
		case QueryAssetHash:
			return queryAssetHash(ctx, req, k)
		case QueryCrossedAmount:
			return queryCrossedAmount(ctx, req, k)
		case QueryCrossedLimit:
			return queryCrossedLimit(ctx, req, k)
		case QueryOperator:
			return queryOperator(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}
func queryProxyHash(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryProxyHashParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	proxyHashBs := k.GetProxyHash(ctx, params.ChainId)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, proxyHashBs)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryAssetHash(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAssetHashParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	assetHashBs := k.GetAssetHash(ctx, params.SourceAssetDenom, params.ChainId)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, assetHashBs)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryCrossedAmount(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAssetHashParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	crossedAmount := k.GetCrossedAmount(ctx, params.SourceAssetDenom, params.ChainId)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, crossedAmount)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryCrossedLimit(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAssetHashParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	crossedLimit := k.GetCrossedLimit(ctx, params.SourceAssetDenom, params.ChainId)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, crossedLimit)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryOperator(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryAssetHashParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	operator := k.GetOperator(ctx)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, operator.Operator)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
