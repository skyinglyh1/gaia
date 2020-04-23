package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/headersync/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"fmt"
)

const (
	// query balance path
	QueryHeader = "header"
	QueryCurrentHeight = "current_height"

)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryHeader:
			return queryHeader(ctx, req, k)
		case QueryCurrentHeight:
			return queryCurrentHeight(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown bank query endpoint")
		}
	}
}

// queryBalance fetch an account's balance for the supplied height.
// Height and account address are passed as first and second path components respectively.
func queryHeader(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryHeaderParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	header, err := k.GetHeaderByHeight(ctx, params.ChainId, params.Height)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("queryHeader, %v", err))
	}
	bz, er := codec.MarshalJSONIndent(types.ModuleCdc, header)
	if er != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
func queryCurrentHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryCurrentHeightParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	height := k.GetCurrentHeight(ctx, params.ChainId)
	fmt.Printf("got height = %d\n", height)
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, height)
	fmt.Printf("internal.keeper.querier.go.bz = %v\n", bz)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
