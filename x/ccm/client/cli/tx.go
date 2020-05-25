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
	"github.com/cosmos/gaia/x/ccm/internal/types"
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
		SendProcessCrossChainTxTxCmd(cdc),
	)...)
	return txCmd
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
