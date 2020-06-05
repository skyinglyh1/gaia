package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gaia/x/btcx/internal/types"
)

const (
	QueryDenomInfo       = "denom_info"
	QueryDenomInfoWithid = "denom_info_id"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryDenomInfo:
			return queryDenomInfo(ctx, req, k)
		case QueryDenomInfoWithid:
			return queryDenomInfoWithId(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}

func queryDenomInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDenomInfo

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	denomInfo := k.GetDenomInfo(ctx, params.Denom)

	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, denomInfo)
	if e != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", e))
	}

	return bz, nil
}

func queryDenomInfoWithId(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDenomInfoWithId

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	denomInfo := k.GetDenomInfoWithId(ctx, params.Denom, params.ChainId)

	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, denomInfo)
	if e != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", e))
	}

	return bz, nil
}
