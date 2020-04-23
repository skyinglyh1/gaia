package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/gaia/x/headersync/client/common"
	"github.com/cosmos/gaia/x/headersync/internal/types"
	mctype "github.com/ontio/multi-chain/core/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	hsQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the header sync module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	hsQueryCmd.AddCommand(
		GetCmdQueryHeader(queryRoute, cdc),
	)

	return hsQueryCmd
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryHeader(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "header [chainId] [height]",
		Args:  cobra.ExactArgs(2),
		Short: "Query header",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query distribution outstanding (un-withdrawn) rewards
for a validator and all their delegations.

Example:
$ %s query distr validator-outstanding-rewards cosmosvaloper1lwjmdnks33xwnmfayc64ycprww49n33mtm92ne
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)


			chainIdStr := args[0]
			heightStr := args[1]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			height, err := strconv.ParseUint(heightStr, 10, 32)
			if err != nil {
				return err
			}

			res, err := common.QueryHeader(cliCtx, queryRoute, uint64(chainId), uint32(height))
			if err != nil {
				return err
			}
			var header mctype.Header
			cdc.MustUnmarshalJSON(res, &header)
			return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}

type MCHeader struct {
	mctype.Header
}

func (header MCHeader) String() string {
	return fmt.Sprintf("header of height: %d is :%s\n", header.Height, header)
}

