package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/headersync/internal/types"
)

const (
	QueryHeader        = "header"
	QueryCurrentHeight = "current_height"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryHeader:
			return queryHeader(ctx, req, k)
		case QueryCurrentHeight:
			return queryCurrentHeight(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}

func queryHeader(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryHeaderParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	header, err := k.GetHeaderByHeight(ctx, params.ChainId, params.Height)
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("queryHeader, %v", err))
	}
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, header)
	if e != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", e.Error()))
	}

	return bz, nil
}
func queryCurrentHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryCurrentHeightParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	height, err := k.GetCurrentHeight(ctx, params.ChainId)
	if err != nil {
		return nil, err
	}
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, height)
	fmt.Printf("internal.keeper.querier.go.bz = %v\n", bz)
	if e != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", e.Error()))
	}

	return bz, nil
}
