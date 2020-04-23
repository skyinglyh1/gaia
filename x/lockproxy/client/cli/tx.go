package cli

import (
	"github.com/spf13/cobra"

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
		Short:                      "lock proxy module transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(
		SendCreateCoinsTxCmd(cdc),
		SendBindProxyHashTxCmd(cdc),
		SendBindAssetHashTxCmd(cdc),
		SendLockTxCmd(cdc),
		SendProcessCrossChainTxTxCmd(cdc),
	)...)
	return txCmd
}

func SendCreateCoinsTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createcoins [coins_str]",
		Short: "Create coins by operator",
		Args:  cobra.ExactArgs(1),
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
		Short: "Create and sign a bindProxyHash tx",
		Args:  cobra.ExactArgs(2),
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
		Short: "Sync one block header",
		Args:  cobra.ExactArgs(5),
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
			msg := types.NewMsgBindAssetParam(cliCtx.GetFromAddress(), sourceAssetDenom, targetChainId, targetAssetHash, &limit, isTargetChainAsset)
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
		Use:   "lock [block_header_hex_string]",
		Short: "Sync one block header",
		Args:  cobra.ExactArgs(4),
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
			toAddress, err := sdk.AccAddressFromBech32(toAddressStr)
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
		Short: "Sync one block header",
		Args:  cobra.ExactArgs(4),
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
