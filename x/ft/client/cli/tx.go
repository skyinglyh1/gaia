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
	"github.com/cosmos/gaia/x/ft/internal/types"
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
		SendCreateAndDelegateCoinToProxyTxCmd(cdc),

		SendCreateDenomTxCmd(cdc),
		SendBindAssetHashTxCmd(cdc),
		SendLockTxCmd(cdc),
		SendCreateCoinTxCmd(cdc),
	)...)
	return txCmd
}

func SendCreateAndDelegateCoinToProxyTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createdelegatedenom [creator] [coin] [lock_proxy]",
		Short: "Create coin by creator, and immediately delegate to the lock proxy module account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx ft createdelegatedenom [creator_address] [1000bch] 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			creator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			lockProxy, err := hex.DecodeString(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCoinAndDelegateToProxy(creator, coin, sdk.AccAddress(lockProxy))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
func SendCreateDenomTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createdenom [creator] [denom]",
		Short: "Create denom by creator, and from will become operator automaticall",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx ft createdenom [creator_address],ont 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			creator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			//for i, coin := range coins {
			//	coins[i] = sdk.NewCoin(coin.Denom, sdk.NewInt(0))
			//}
			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgCreateDenom(creator, args[1])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendBindAssetHashTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bindassethash [source_asset_denom] [target_chainId] [target_asset_hash] [initialAmount]",
		Short: "bind asset hash by the operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx crosschain bindassethash ont 3 00000000000000000001 100000 true 
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

			targetAssetHashStr := args[2]
			if targetAssetHashStr[0:2] == "0x" {
				targetAssetHashStr = targetAssetHashStr[2:]
			}
			toAssetHash, err := hex.DecodeString(targetAssetHashStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'targetProxyHash' error:%v", err)
			}

			msg := types.NewMsgBindAssetHash(cliCtx.GetFromAddress(), sourceAssetDenom, toChainId, toAssetHash)
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
			msg := types.NewMsgLock(cliCtx.GetFromAddress(), sourceAssetDenom, toChainId, toAddress, value)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
func SendCreateCoinTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createcoins [creator] [coins]",
		Short: "Create coins by creator, and from will receive all the coins",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx ft createdenom [creator_address],ont 
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			creator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			//for i, coin := range coins {
			//	coins[i] = sdk.NewCoin(coin.Denom, sdk.NewInt(0))
			//}
			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgCreateCoins(creator, args[1])
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
