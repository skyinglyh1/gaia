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
	"github.com/cosmos/gaia/x/btcx/internal/types"
	"math/big"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "btcx module send transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(
		SendCreateCoinTxCmd(cdc),
		SendBindAssetHashTxCmd(cdc),
		SendLockTxCmd(cdc),
	)...)
	return txCmd
}

func SendCreateCoinTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createcoin [coins_str]",
		Short: "Create coin by from, and from will become denom creator automaticall",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx btcx createcoin btccc redeemscript 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := args[0]
			redeemScript := args[1]

			msg := types.NewMsgCreateCoin(cliCtx.GetFromAddress(), denom, redeemScript)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendBindAssetHashTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bindassethash [source_asset_denom] [to_chainId] [to_asset_hash]",
		Short: "bind asset hash by the operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx btcx bindassethash btca 3 12341234
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetDenom := args[0]

			toChainId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			toAssetHashStr := args[2]
			if toAssetHashStr[0:2] == "0x" {
				toAssetHashStr = toAssetHashStr[2:]
			}
			toAssetHash, err := hex.DecodeString(toAssetHashStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'toAssetHash' error:%v", err)
			}
			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBindAssetParam(cliCtx.GetFromAddress(), sourceAssetDenom, toChainId, toAssetHash)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendLockTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock [source_asset_denom] [to_chain_id] [to_address] [amount]",
		Short: "lock amount of source_asset_denom and aim to release amount in to_chain_id chain to to_address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx crosschain lock btca 3 616f2a4a38396ff203ea01e6c070ae421bb8ce2d 123 
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
