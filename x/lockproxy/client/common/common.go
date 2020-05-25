package common

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/gaia/x/lockproxy/internal/keeper"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
)

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryProxyByOperator(cliCtx context.CLIContext, queryRoute string, operator sdk.AccAddress) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryProxyByOperator),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryProxyByOperator(operator)),
	)
	return res, err
}

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryProxyHash(cliCtx context.CLIContext, queryRoute string, lockProxyHash []byte, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryProxyHash),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryProxyHashParams(lockProxyHash, chainId)),
	)
	return res, err
}

func QueryAssetHash(cliCtx context.CLIContext, queryRoute string, lockProxyHash []byte, sourceAssetDenom string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryAssetHash),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryAssetHashParams(lockProxyHash, sourceAssetDenom, chainId)),
	)
	return res, err
}

func QueryLockedAmt(cliCtx context.CLIContext, queryRoute string, sourceAssetDenom string) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryLockedAmt),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryLockedAmtParam(sourceAssetDenom)),
	)
	return res, err
}
