package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/gaia/x/crosschain/internal/keeper"
	"github.com/cosmos/gaia/x/crosschain/internal/types"
)



// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryHeader(cliCtx context.CLIContext, queryRoute string, chainId uint64, height uint32) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryHeader),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryHeaderParams(chainId, height)),
	)
	return res, err
}

func QueryCurrentHeaderHeight(cliCtx context.CLIContext, queryRoute string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryCurrentHeight),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryCurrentHeightParams(chainId)),
	)

	return res, err
}


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
