package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/gaia/x/ft/internal/types"
)

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryDenomInfo(cliCtx context.CLIContext, queryRoute string, denom string) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDenom),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryDenomInfo(denom)),
	)
	return res, err
}

func QueryDenomCrossChainInfo(cliCtx context.CLIContext, queryRoute string, denom string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDenomCrossChain),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryDenomCrossChainInfo(denom, chainId)),
	)

	return res, err
}
