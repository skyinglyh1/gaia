package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/ft/internal/types"
)

const (
	QueryDenomInfo       = "denom_info"
	QueryDenomInfoWithid = "denom_info_id"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {

		case QueryDenomInfo:
			return queryDenomInfo(ctx, req, k)
		case QueryDenomInfoWithid:
			return queryDenomInfoWithId(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}

func queryDenomInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryDenomInfo

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	assetHashBs := k.GetDenomInfo(ctx, params.Denom)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, assetHashBs)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryDenomInfoWithId(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryDenomInfoWithId

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	crossedAmount := k.GetDenomInfoWithId(ctx, params.Denom, params.ChainId)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, crossedAmount)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
