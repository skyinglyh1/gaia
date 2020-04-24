package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/gaia/x/lockproxy/internal/keeper"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
)

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryProxyHash(cliCtx context.CLIContext, queryRoute string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryProxyHash),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryProxyHashParams(chainId)),
	)
	return res, err
}

func QueryAssetHash(cliCtx context.CLIContext, queryRoute string, sourceAssetDenom string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryAssetHash),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryAssetHashParams(sourceAssetDenom, chainId)),
	)
	return res, err
}

func QueryCrossedAmount(cliCtx context.CLIContext, queryRoute string, sourceAssetDenom string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryCrossedAmount),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryCrossedAmountParam(sourceAssetDenom, chainId)),
	)
	return res, err
}

func QueryCrossedLimit(cliCtx context.CLIContext, queryRoute string, sourceAssetDenom string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryCrossedLimit),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryCrossedLimitParam(sourceAssetDenom, chainId)),
	)
	return res, err
}

func QueryOperator(cliCtx context.CLIContext, queryRoute string) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryOperator),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryOperatorParam()),
	)
	return res, err
}
