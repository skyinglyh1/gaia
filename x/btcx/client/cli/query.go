package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/gaia/x/btcx/client/common"
	"github.com/cosmos/gaia/x/btcx/internal/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for the minting module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	ccQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the crossChain module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return ccQueryCmd
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryDenomInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "denomInfo [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query denom info",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query block header for a specific height
already synced from another blockchain, normally, relayer-chain (with chainId=0), into current chain

Example:
$ %s query crosschain header 0 1
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := args[0]

			res, err := common.QueryDenomInfo(cliCtx, queryRoute, denom)
			if err != nil {
				return err
			}
			var denomInfo types.DenomInfo
			cdc.MustUnmarshalJSON(res, &denomInfo)
			fmt.Printf("denomInfo of denom:%s is:\n %s\n", denom, denomInfo.String())
			return nil
		},
	}
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryDenomInfoWithChainId(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "denomInfoId [denom] [chainId]",
		Args:  cobra.ExactArgs(2),
		Short: "Query denom info correlated with chainId",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query btcx denomInfoId btca 2
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := args[0]
			toChainId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := common.QueryDenomInfoWithId(cliCtx, queryRoute, denom, toChainId)
			if err != nil {
				return err
			}

			var denomInfoWithId types.DenomInfoWithId
			cdc.MustUnmarshalJSON(res, &denomInfoWithId)
			fmt.Printf("denomInfo of denom:%s for chainId:%d is:\n %s\n", denom, toChainId, denomInfoWithId.String())
			return nil
			//return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}
