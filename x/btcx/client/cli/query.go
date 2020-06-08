package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/gaia/x/btcx/client/common"
	"github.com/cosmos/gaia/x/btcx/internal/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// GetQueryCmd returns the cli query commands for the minting module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	ccQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ccQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryDenomInfo(queryRoute, cdc),
			GetCmdQueryDenomInfoWithChainId(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryDenomInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "denominfo [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query denom info",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a specific denom or coin info incluing the coin creator,  coin total supply, the 
redeem script and redeem script hash

Example:
$ %s query btcx denomInfo btcx
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
		Use:   "denomccinfo [denom] [chainId]",
		Args:  cobra.ExactArgs(2),
		Short: "Query denom info correlated with a specific chainId",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a specific denom or coin info correlated with a specific chainId incluing the coin creator,  coin total supply, the 
redeem script and redeem script hash, toChainId and the corresponding toAssetHash in hex format

Example:
$ %s query btcx denomInfoId btcx 2
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

			res, err := common.QueryDenomCrossChainInfo(cliCtx, queryRoute, denom, toChainId)
			if err != nil {
				return err
			}

			var denomCCInfo types.DenomCrossChainInfo
			cdc.MustUnmarshalJSON(res, &denomCCInfo)
			fmt.Printf("denom cross chain Info of denom:%s for chainId:%d is:\n %s\n", denom, toChainId, denomCCInfo.String())
			return nil
			//return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}
