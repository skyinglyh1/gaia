package common

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/gaia/x/ccm/internal/keeper"
	"github.com/cosmos/gaia/x/ccm/internal/types"
)

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryContainToContractAddr(cliCtx context.CLIContext, queryRoute string, keystore string, toContractAddr []byte, fromChainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryContractToContractAddr),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryContainToContract(keystore, toContractAddr, fromChainId)),
	)
	return res, err
}
