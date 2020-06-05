package keeper

import (
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/ccm/internal/types"
)

const (
	QueryContractToContractAddr = "denom_info_id"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryContractToContractAddr:
			return queryContractToContractAddr(ctx, req, k)

		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown minting query endpoint: %s", path[0]))
		}
	}
}

func queryContractToContractAddr(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryContainToContract

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}
	resInfo := k.IfContainToContract(ctx, params.KeyStore, params.ToContractAddr, params.FromChainId)

	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, resInfo)
	if e != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", e.Error()))
	}

	return bz, nil
}
