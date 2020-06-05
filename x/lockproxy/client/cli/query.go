package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/version"
	"strings"

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
	ccQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ccQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryProxyByOperator(queryRoute, cdc),
			GetCmdQueryProxyHash(queryRoute, cdc),
			GetCmdQueryAssetHash(queryRoute, cdc),
			GetCmdQueryLockedAmount(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryProxyByOperator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proxyhash [operator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query lockproxy hex string by the operator/creator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the created lock proxy contract hash created by the operator

Example:
$ %s query %s height 0
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			operatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := common.QueryProxyByOperator(cliCtx, queryRoute, operatorAddr)

			if err != nil {
				return err
			}
			var proxyHash []byte
			cdc.MustUnmarshalJSON(res, &proxyHash)
			fmt.Printf("creator:%s with lock proxy hash:%x \n", operatorAddr, proxyHash)
			return nil
		},
	}
}

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryProxyHash(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proxyhash [lockproxy] [chainId]",
		Short: "Query the proxy hash deployed in another chain with chainId",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query proxy contract hash bond with self in another blockchain 
with chainId

Example:
$ %s query %s proxyhash 3
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			lockProxyHash, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}
			chainIdStr := args[1]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryProxyHash(cliCtx, queryRoute, lockProxyHash, chainId)
			if err != nil {
				return err
			}
			var proxyHash []byte
			cdc.MustUnmarshalJSON(res, &proxyHash)
			fmt.Printf("toChain proxy_hash : %s\n", hex.EncodeToString(proxyHash))
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryAssetHash(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "assethash [lockproxy/operator] [sourceassetdenom] [chainId]",
		Short: "Query the asset hash in chainId chain corresponding with soureAssetDenom",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query %s assethash height 0
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			lockProxy, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				lockProxyBs, err1 := hex.DecodeString(args[0])
				if err1 != nil {
					return sdk.ErrInternal(fmt.Sprintf("lockproxy: %s or operator error: %s", err, err1))
				}
				lockProxy = append(lockProxy, lockProxyBs...)
			}

			sourceAssetdenom := args[1]

			chainIdStr := args[2]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryAssetHash(cliCtx, queryRoute, lockProxy, sourceAssetdenom, chainId)
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

func GetCmdQueryLockedAmount(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "lockedamt [denom]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the asset crossed amount in chainId chain corresponding with soureAssetDenom
Example:
$ %s query %s lockedamt stake
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetdenom := args[0]

			res, err := common.QueryLockedAmt(cliCtx, queryRoute, sourceAssetdenom)
			if err != nil {
				return err
			}
			var crossedLimit sdk.Int
			cdc.MustUnmarshalJSON(res, &crossedLimit)
			fmt.Printf("locked_amount for%s : %s\n", sourceAssetdenom, crossedLimit.String())
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}
