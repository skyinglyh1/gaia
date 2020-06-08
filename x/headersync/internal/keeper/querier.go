package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gaia/x/headersync/internal/types"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryHeader:
			return queryHeader(ctx, req, k)
		case types.QueryCurrentHeight:
			return queryCurrentHeight(ctx, req, k)
		case types.QueryKeyHeights:
			return queryKeyHeights(ctx, req, k)
		case types.QueryKeyHeight:
			return queryKeyHeight(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}

func queryHeader(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryHeaderParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	header, err := k.GetHeaderByHeight(ctx, params.ChainId, params.Height)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("queryHeader, %v", err))
	}
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, header)
	if e != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", e.Error()))
	}

	return bz, nil
}
func queryCurrentHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryCurrentHeightParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	height, err := k.GetCurrentHeight(ctx, params.ChainId)
	if err != nil {
		return nil, err
	}
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, height)
	fmt.Printf("internal.keeper.querier.go.bz = %v\n", bz)
	if e != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", e.Error()))
	}

	return bz, nil
}

func queryKeyHeights(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryKeyHeightsParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	heights := k.GetKeyHeights(ctx, params.ChainId)
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, heights)
	if e != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", e.Error()))
	}

	return bz, nil
}

func queryKeyHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryKeyHeightParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	height, err := k.FindKeyHeight(ctx, params.Height, params.ChainId)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("queryKeyHeight, FindKeyHeight error:%s", err))
	}
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, height)
	if e != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", e.Error()))
	}

	return bz, nil
}
