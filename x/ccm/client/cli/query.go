package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/gaia/x/ccm/client/common"
	"github.com/cosmos/gaia/x/ccm/internal/types"
	"github.com/spf13/cobra"
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
func GetCmdQueryContainToContractAddr(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "ifcontainaddr [denom]",
		Args:  cobra.ExactArgs(3),
		Short: "Query denom info",
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
