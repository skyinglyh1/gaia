package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/version"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/gaia/x/ft/client/common"
	"github.com/cosmos/gaia/x/ft/internal/types"
	"strconv"
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

	ccQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryDenomInfo(queryRoute, cdc),
			GetCmdQueryDenomInfoWithId(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryDenomInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "denominfo [sourceassetdenom]",
		Short: "Query the asset hash in chainId chain corresponding with soureAssetDenom",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query crosschain asset height 0
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetdenom := args[0]

			res, err := common.QueryDenomInfo(cliCtx, queryRoute, sourceAssetdenom)
			if err != nil {
				return err
			}
			var denomInfo types.DenomInfo
			cdc.MustUnmarshalJSON(res, &denomInfo)
			fmt.Printf("denomInfo of denom:%s is:\n %s\n", sourceAssetdenom, denomInfo.String())
			return nil
		},
	}
}

func GetCmdQueryDenomInfoWithId(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "denominfowithid [sourceassetdenom] [chainid]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the asset crossed limit in chainId chain corresponding with soureAssetDenom
Example:
$ %s query crosschain crossedlimit stake 3
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetdenom := args[0]

			chainId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryDenomInfoWithId(cliCtx, queryRoute, sourceAssetdenom, chainId)
			if err != nil {
				return err
			}
			var denomInfo types.DenomInfoWithId
			cdc.MustUnmarshalJSON(res, &denomInfo)
			fmt.Printf("denomInfo in detail of denom:%s is:\n %s\n", sourceAssetdenom, denomInfo.String())
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}
