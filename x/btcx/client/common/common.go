package common

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/gaia/x/btcx/internal/keeper"
	"github.com/cosmos/gaia/x/btcx/internal/types"
)

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryDenomInfo(cliCtx context.CLIContext, queryRoute string, denom string) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryDenomInfo),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryDenomInfo(denom)),
	)
	return res, err
}

func QueryDenomInfoWithId(cliCtx context.CLIContext, queryRoute string, denom string, toChainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryDenomInfoWithid),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryDenomInfoWithId(denom, toChainId)),
	)
	return res, err
}
