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
	"github.com/cosmos/gaia/x/crosschain/internal/types"
	"math/big"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "crosschain module send transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(
		SendSyncGenesisTxCmd(cdc),
		SendSyncHeaderTxCmd(cdc),

		SendCreateCoinsTxCmd(cdc),
		SendBindProxyHashTxCmd(cdc),
		SendBindAssetHashTxCmd(cdc),
		SendLockTxCmd(cdc),
		SendProcessCrossChainTxTxCmd(cdc),

		SendSetRedeemScriptTxCmd(cdc),
		SendBindNoVMTxCmd(cdc),
	)...)
	return txCmd
}

func SendSyncGenesisTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "syncgenesis [genesis_header_hexstring]",
		Short: "Create and sign a syncgenesis tx",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			genesisHeaderBytes, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}
			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgSyncGenesisParam(cliCtx.GetFromAddress(), genesisHeaderBytes)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendSyncHeaderTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "syncheader [block_header_hex_string]",
		Short: "Sync one block header",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			headerBytes, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgSyncHeadersParam(cliCtx.GetFromAddress(), [][]byte{headerBytes})
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendCreateCoinsTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createcoins [coins_str]",
		Short: "Create coins by from, and from will become operator automaticall",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx crosschain createcoins 1000000000ont,1000000000000000000ong 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}
			//for i, coin := range coins {
			//	coins[i] = sdk.NewCoin(coin.Denom, sdk.NewInt(0))
			//}
			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgCreateCoins(cliCtx.GetFromAddress(), coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
func SendBindProxyHashTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bindproxyhash [target_chain_id] [target_proxy_hash]",
		Short: "bindproxyhash by the operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx crosschain bindproxyhash 3 11223344556677889900 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			targetChainIdStr := args[0]
			targetProxyHashStr := args[1]

			targetChainId, err := strconv.ParseUint(targetChainIdStr, 10, 64)
			if err != nil {
				return err
			}
			if targetProxyHashStr[0:2] == "0x" {
				targetProxyHashStr = targetProxyHashStr[2:]
			}
			targetProxyHash, err := hex.DecodeString(targetProxyHashStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'targetProxyHash' error:%v", err)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBindProxyParam(cliCtx.GetFromAddress(), targetChainId, targetProxyHash)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendBindAssetHashTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bindassethash [source_asset_denom] [target_chainId] [target_asset_hash] [limit] [is_target_chain_asset]",
		Short: "bind asset hash by the operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx crosschain bindassethash ont 3 00000000000000000001 100000 true 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetDenom := args[0]

			targetChainIdStr := args[1]
			targetChainId, err := strconv.ParseUint(targetChainIdStr, 10, 64)
			if err != nil {
				return err
			}

			targetAssetHashStr := args[2]
			if targetAssetHashStr[0:2] == "0x" {
				targetAssetHashStr = targetAssetHashStr[2:]
			}
			targetAssetHash, err := hex.DecodeString(targetAssetHashStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'targetProxyHash' error:%v", err)
			}

			limitBigInt, ok := big.NewInt(0).SetString(args[3], 10)
			if !ok {
				return fmt.Errorf("read limit as big int from args[3] failed")
			}
			limit := sdk.NewIntFromBigInt(limitBigInt)

			isTargetChainAsset, err := strconv.ParseBool(args[4])
			if err != nil {
				return fmt.Errorf("read istargetChainAsset parameter from args[4] failed")
			}
			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBindAssetParam(cliCtx.GetFromAddress(), sourceAssetDenom, targetChainId, targetAssetHash, limit, isTargetChainAsset)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

//coins, err := sdk.ParseCoins(args[1])
//if err != nil {
//return err
//}

func SendLockTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock [source_asset_denom] [to_chain_id] [to_address] [amount]",
		Short: "lock amount of source_asset_denom and aim to release amount in to_chain_id chain to to_address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx crosschain lock ont 3 616f2a4a38396ff203ea01e6c070ae421bb8ce2d 123 
`,
				version.ClientName,
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

			toAddressStr := args[2]
			toAddress, err := hex.DecodeString(toAddressStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'toAddress' error:%v", err)
			}

			valueBigInt, ok := big.NewInt(0).SetString(args[3], 10)
			if !ok {
				return fmt.Errorf("read value as big int from args[3] failed")
			}
			value := sdk.NewIntFromBigInt(valueBigInt)

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgLock(cliCtx.GetFromAddress(), sourceAssetDenom, toChainId, toAddress, &value)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendProcessCrossChainTxTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "processcrosschaintx [from_chainId] [height] [proof] [header]",
		Short: "process cross chain tx targeting at current cosmos-type chain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx crosschain processcrosschaintx 3 1000 'proof_hex_str_at_height_1000' 'header_1000_can_be_empty_if_header_already_synced'
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			fromChainIdStr := args[0]
			fromChainId, err := strconv.ParseUint(fromChainIdStr, 10, 64)
			if err != nil {
				return err
			}

			heightStr := args[1]
			height, err := strconv.ParseUint(heightStr, 10, 32)
			if err != nil {
				return err
			}

			proof := args[2]
			header, err := hex.DecodeString(args[3])
			if err != nil {
				return fmt.Errorf("decode hex string 'header' error:%v", err)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgProcessCrossChainTx(cliCtx.GetFromAddress(), fromChainId, uint32(height), proof, header)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendSetRedeemScriptTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setredeem [denom] [redeem_key] [redeem_script]",
		Short: "set redeem script indexed by redeem key",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx crosschain createcoins 1000000000ont,1000000000000000000ong 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			denom := args[0]
			redeemKey, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}
			redeemScript, err := hex.DecodeString(args[2])
			if err != nil {
				return err
			}
			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgSetRedeemScript(cliCtx.GetFromAddress(), denom, redeemKey, redeemScript)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendBindNoVMTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bindnovm [source_denom] [target_chainId] [target_asset_hash] [limit]",
		Short: "bind asset hash by the operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx crosschain bindassethash hex(btc) 3 00000000000000000001 100000
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceDenom := args[0]

			targetChainIdStr := args[1]
			targetChainId, err := strconv.ParseUint(targetChainIdStr, 10, 64)
			if err != nil {
				return err
			}

			targetAssetHashStr := args[2]
			if targetAssetHashStr[0:2] == "0x" {
				targetAssetHashStr = targetAssetHashStr[2:]
			}
			targetAssetHash, err := hex.DecodeString(targetAssetHashStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'targetProxyHash' error:%v", err)
			}

			limitBigInt, ok := big.NewInt(0).SetString(args[3], 10)
			if !ok {
				return fmt.Errorf("read limit as big int from args[3] failed")
			}
			limit := sdk.NewIntFromBigInt(limitBigInt)

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBindNoVMChainAssetHash(cliCtx.GetFromAddress(), sourceDenom, targetChainId, targetAssetHash, limit)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
