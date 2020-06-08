package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gaia/x/ft/internal/types"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {

		case types.QueryDenom:
			return queryDenomInfo(ctx, req, k)
		case types.QueryDenomWithid:
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
	assetHashBs := k.GetDenomInfo(ctx, params.Denom)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, assetHashBs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", err.Error()))
	}

	return bz, nil
}

func queryDenomInfoWithId(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDenomInfoWithId

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("failed to parse params: %s", err))
	}
	crossedAmount := k.GetDenomInfoWithId(ctx, params.Denom, params.ChainId)

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, crossedAmount)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("could not marshal result to JSON: %s", err))
	}

	return bz, nil
}
