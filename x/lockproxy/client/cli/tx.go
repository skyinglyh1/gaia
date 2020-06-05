package cli

import (
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"

	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
	"math/big"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s module send transaction subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(

		SendCreateLockProxyTxCmd(cdc),
		SendBindProxyHashTxCmd(cdc),
		SendBindAssetHashTxCmd(cdc),
		SendLockTxCmd(cdc),
	)...)
	return txCmd
}

func SendCreateLockProxyTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createlockproxy [creator]",
		Short: "Create lockproxy contract by creator, needs creator's signature",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx %s createlockproxy [creator_address] 
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			msg := types.NewMsgCreateLockProxy(cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendBindProxyHashTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bindproxyhash [to_chain_id] [to_chain_proxy_hash]",
		Short: "bindproxyhash by the operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx %s bindproxyhash 3 11223344556677889900 
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			toChainIdStr := args[0]
			toChainProxyHashStr := args[1]

			targetChainId, err := strconv.ParseUint(toChainIdStr, 10, 64)
			if err != nil {
				return err
			}
			if toChainProxyHashStr[0:2] == "0x" {
				toChainProxyHashStr = toChainProxyHashStr[2:]
			}
			toChainProxyHash, err := hex.DecodeString(toChainProxyHashStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'targetProxyHash' error:%v", err)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBindProxyHash(cliCtx.GetFromAddress(), targetChainId, toChainProxyHash)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendBindAssetHashTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bindassethash [source_asset_denom] [to_chainId] [to_asset_hash] [initial_balance_of_lockproxy_module_acct_for_asset_denom]",
		Short: "bind asset hash by the operator, ",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx %s bindassethash ont 3 00000000000000000001 100000
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetDenom := args[0]

			toChainIdStr := args[1]
			toChainId, err := strconv.ParseUint(toChainIdStr, 10, 64)
			if err != nil {
				return err
			}

			toAssetHashStr := args[2]
			if toAssetHashStr[0:2] == "0x" {
				toAssetHashStr = toAssetHashStr[2:]
			}
			toAssetHash, err := hex.DecodeString(toAssetHashStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'targetProxyHash' error:%v", err)
			}

			limitBigInt, ok := big.NewInt(0).SetString(args[3], 10)
			if !ok {
				return fmt.Errorf("read limit as big int from args[3] failed")
			}
			initialAmt := sdk.NewIntFromBigInt(limitBigInt)
			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBindAssetHash(cliCtx.GetFromAddress(), sourceAssetDenom, toChainId, toAssetHash, initialAmt)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendLockTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock [lock_proxy_hash] [source_asset_denom] [to_chain_id] [to_address] [amount]",
		Short: "lock amount of source_asset_denom and aim to release amount in to_chain_id chain to to_address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx %s lock 12341234 ont 3 616f2a4a38396ff203ea01e6c070ae421bb8ce2d 123 
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			lockProxyHash, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			sourceAssetDenom := args[1]

			toChainIdStr := args[2]
			toChainId, err := strconv.ParseUint(toChainIdStr, 10, 64)
			if err != nil {
				return err
			}

			toAddressStr := args[3]
			toAddress, err := hex.DecodeString(toAddressStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'toAddress' error:%v", err)
			}

			valueBigInt, ok := big.NewInt(0).SetString(args[4], 10)
			if !ok {
				return fmt.Errorf("read value as big int from args[3] failed")
			}
			value := sdk.NewIntFromBigInt(valueBigInt)

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgLock(lockProxyHash, cliCtx.GetFromAddress(), sourceAssetDenom, toChainId, toAddress, value)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
