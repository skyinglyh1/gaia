package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/version"
	polytype "github.com/cosmos/gaia/x/headersync/poly-utils/core/types"

	"github.com/spf13/cobra"

	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/gaia/x/headersync/client/common"
	"github.com/cosmos/gaia/x/headersync/internal/types"
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
func GetCmdQueryHeader(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "header [chainId] [height]",
		Args:  cobra.ExactArgs(2),
		Short: "Query header of chainId of height",
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
			var header polytype.Header
			cdc.MustUnmarshalJSON(res, &header)
			fmt.Printf("header of height:%d is:\n %s\n", header.Height, header.String())
			return nil
		},
	}
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryCurrentHeight(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "height [chainId]",
		Args:  cobra.ExactArgs(1),
		Short: "Query block height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query crosschain height 0
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			chainIdStr := args[0]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}

			res, err := common.QueryCurrentHeaderHeight(cliCtx, queryRoute, chainId)
			fmt.Printf("cli.query.res = %v\n", res)
			if err != nil {
				return err
			}
			var height uint32
			cdc.MustUnmarshalJSON(res, &height)
			fmt.Printf("current synced header height of chainid:%d is: %d\n", chainId, height)
			return nil
			//return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}

func GetCmdQueryKeyHeights(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "keyheights [chainId]",
		Args:  cobra.ExactArgs(1),
		Short: "Query block height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query crosschain height 0
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			chainIdStr := args[0]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}

			res, err := common.QueryKeyHeights(cliCtx, queryRoute, chainId)
			fmt.Printf("cli.query.res = %v\n", res)
			if err != nil {
				return err
			}
			var keyHeights types.KeyHeights
			cdc.MustUnmarshalJSON(res, &keyHeights)
			fmt.Printf("Key heights of chainid: %d is: %+v\n", chainId, keyHeights)
			return nil
			//return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}

func GetCmdQueryKeyHeight(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "keyheight [chainId] [height]",
		Args:  cobra.ExactArgs(2),
		Short: "Query block height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query crosschain height 0
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			chainId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			height, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			res, err := common.QueryKeyHeight(cliCtx, queryRoute, chainId, uint32(height))
			fmt.Printf("cli.query.res = %v\n", res)
			if err != nil {
				return err
			}
			var keyHeight uint32
			cdc.MustUnmarshalJSON(res, &keyHeight)
			fmt.Printf("Key height of chainid: %d is: %d\n", chainId, keyHeight)
			return nil
			//return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}
