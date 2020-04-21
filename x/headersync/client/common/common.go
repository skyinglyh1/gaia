package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/gaia/x/headersync/internal/types"
	"github.com/cosmos/gaia/x/headersync/internal/keeper"
)

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryHeader(cliCtx context.CLIContext, queryRoute string, chainId uint64, height uint32) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryHeader),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryHeaderParams(chainId, height)),
	)
	return res, err
}
