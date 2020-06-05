package cli

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/gaia/x/ccm/client/common"
	"github.com/cosmos/gaia/x/ccm/internal/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
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
			GetCmdQueryContainToContractAddr(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryContainToContractAddr(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ifContainContract [module_store_key] [to_contract_addr] [from_chain_id]",
		Args:  cobra.ExactArgs(3),
		Short: "Query if module_store_key module should be targeted to execute `unlock` logic based on ToMerkleValue.MakeTxParam.ToContractAddress and ToMerkleValue.FromChainId",
		Long: strings.TrimSpace(
			fmt.Sprintf(`

Example:
$ %s query %s ifContainContract btcx c330431496364497d7257839737b5e4596f5ac06 2
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			keystore := args[0]
			toContractAddr, _ := hex.DecodeString(args[1])
			fromChainId, _ := strconv.ParseInt(args[2], 10, 64)

			resBs, err := common.QueryContainToContractAddr(cliCtx, queryRoute, keystore, toContractAddr, uint64(fromChainId))
			if err != nil {
				return err
			}
			var res types.QueryContainToContractRes
			cdc.MustUnmarshalJSON(resBs, &res)
			fmt.Printf("QueryContainToContractAddr res is:\n %s\n", res.String())
			return nil
		},
	}
}
