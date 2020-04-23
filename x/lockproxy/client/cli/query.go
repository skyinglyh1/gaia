package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/lockproxy/client/common"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
	"strconv"
)

// GetQueryCmd returns the cli query commands for the minting module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	mintingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the minting module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	mintingQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryProxyHash(queryRoute, cdc),
			GetCmdQueryAssetHash(queryRoute, cdc),
			GetCmdQueryCrossedAmount(queryRoute, cdc),
			GetCmdQueryCrossedLimit(queryRoute, cdc),
			GetCmdQueryOperator(queryRoute, cdc),
		)...,
	)

	return mintingQueryCmd
}

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryProxyHash(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proxy [chainId]",
		Short: "Query the proxy hash",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			chainIdStr := args[0]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryProxyHash(cliCtx, queryRoute, chainId)
			if err != nil {
				return err
			}
			var proxyHash []byte
			cdc.MustUnmarshalJSON(res, &proxyHash)
			fmt.Printf("proxy_hash: %s\n", hex.EncodeToString(proxyHash))
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryAssetHash(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "asset [sourceassetdenom] [chainId]",
		Short: "Query the asset hash in chainId chain corresponding with soureAssetDenom",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetdenom := args[0]

			chainIdStr := args[1]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryAssetHash(cliCtx, queryRoute, sourceAssetdenom, chainId)
			if err != nil {
				return err
			}
			var assetHash []byte
			cdc.MustUnmarshalJSON(res, &assetHash)
			fmt.Printf("asset_hash: %s\n", hex.EncodeToString(assetHash))
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}

// GetCmdQueryAnnualProvisions implements a command to return the current minting
func GetCmdQueryCrossedAmount(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "crossedamount [sourceassetdenom] [chainId]",
		Short: "Query the asset hash in chainId chain corresponding with soureAssetDenom",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetdenom := args[0]

			chainIdStr := args[1]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryCrossedAmount(cliCtx, queryRoute, sourceAssetdenom, chainId)
			if err != nil {
				return err
			}
			var assetHash []byte
			cdc.MustUnmarshalJSON(res, &assetHash)

			fmt.Printf("crossed_amount: %s\n", hex.EncodeToString(assetHash))
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}

func GetCmdQueryCrossedLimit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "crossedamount [sourceassetdenom] [chainId]",
		Short: "Query the asset hash in chainId chain corresponding with soureAssetDenom",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetdenom := args[0]

			chainIdStr := args[1]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryCrossedLimit(cliCtx, queryRoute, sourceAssetdenom, chainId)
			if err != nil {
				return err
			}
			var assetHash []byte
			cdc.MustUnmarshalJSON(res, &assetHash)
			fmt.Printf("crossed_limit: %s\n", hex.EncodeToString(assetHash))
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}

func GetCmdQueryOperator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "operator",
		Short: "Query the asset hash in chainId chain corresponding with soureAssetDenom",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := common.QueryOperator(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			var assetHash sdk.AccAddress
			cdc.MustUnmarshalJSON(res, &assetHash)
			fmt.Printf("operator: %s\n", hex.EncodeToString(assetHash))
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}
