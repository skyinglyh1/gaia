package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/version"
	mctype "github.com/ontio/multi-chain/core/types"
	"strings"

	"github.com/spf13/cobra"

	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/crosschain/client/common"
	"github.com/cosmos/gaia/x/crosschain/internal/types"
	"strconv"
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

	ccQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryHeader(queryRoute, cdc),
			GetCmdQueryCurrentHeight(queryRoute, cdc),

			GetCmdQueryProxyHash(queryRoute, cdc),
			GetCmdQueryAssetHash(queryRoute, cdc),
			GetCmdQueryCrossedAmount(queryRoute, cdc),
			GetCmdQueryCrossedLimit(queryRoute, cdc),
			GetCmdQueryOperator(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}



// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryHeader(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "header [chainId] [height]",
		Args:  cobra.ExactArgs(2),
		Short: "Query header of chainId of height",
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

			chainIdStr := args[0]
			heightStr := args[1]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			height, err := strconv.ParseUint(heightStr, 10, 32)
			if err != nil {
				return err
			}

			res, err := common.QueryHeader(cliCtx, queryRoute, uint64(chainId), uint32(height))
			if err != nil {
				return err
			}
			var header mctype.Header
			cdc.MustUnmarshalJSON(res, &header)
			fmt.Printf("header of height:%d is:\n %s\n", header.Height, MCHeader{header}.String())
			return nil
			//return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}

type MCHeader struct {
	mctype.Header
}

func (header MCHeader) String() string {
	blockHash := header.Hash()
	return fmt.Sprintf(`
	Version: 		 :%d
	ChainID          :%d
	PrevBlockHash    :%s
	TransactionsRoot :%s
	CrossStateRoot   :%s
	BlockRoot        :%s
	Timestamp        :%d
	Height           :%d
	ConsensusData    :%d
	ConsensusPayload :%s
	NextBookkeeper   :%s

	Bookkeepers []keypair.PublicKey
	SigData     [][]byte

	hash 			:%s
	
`, header.Version, header.ChainID, header.PrevBlockHash.ToHexString(), header.TransactionsRoot.ToHexString(), header.CrossStateRoot.ToHexString(),
		header.BlockRoot.ToHexString(), header.Timestamp, header.Height, header.ConsensusData, hex.EncodeToString(header.ConsensusPayload), header.NextBookkeeper.ToBase58(), blockHash.ToHexString())
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryCurrentHeight(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "height [chainId]",
		Args:  cobra.ExactArgs(1),
		Short: "Query block height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query crosschain height 0
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			chainIdStr := args[0]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}

			res, err := common.QueryCurrentHeaderHeight(cliCtx, queryRoute, chainId)
			fmt.Printf("cli.query.res = %v\n", res)
			if err != nil {
				return err
			}
			var height uint32
			cdc.MustUnmarshalJSON(res, &height)
			fmt.Printf("current synced header height of chainid:%d is: %d\n", chainId, height)
			return nil
			//return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}


// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryProxyHash(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proxy [chainId]",
		Short: "Query the proxy hash",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query proxy contract hash bond with self in another blockchain 
with chainId

Example:
$ %s query crosschain proxy 3
`,
				version.ClientName,
			),
		),
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
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query crosschain asset height 0
`,
				version.ClientName,
			),
		),
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
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the asset crossed amount in chainId chain corresponding with soureAssetDenom
Example:
$ %s query crosschain crossedamount stake 3
`,
				version.ClientName,
			),
		),
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
			var crossedAmount sdk.Int
			cdc.MustUnmarshalJSON(res, &crossedAmount)

			fmt.Printf("crossed_amount: %s\n", crossedAmount.String())
			return nil
		},
	}
}

func GetCmdQueryCrossedLimit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "crossedlimit [sourceassetdenom] [chainId]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the asset crossed limit in chainId chain corresponding with soureAssetDenom
Example:
$ %s query crosschain crossedlimit stake 3
`,
				version.ClientName,
			),
		),
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
			var crossedLimit sdk.Int
			cdc.MustUnmarshalJSON(res, &crossedLimit)
			fmt.Printf("crossed_limit: %s\n", crossedLimit.String())
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}

func GetCmdQueryOperator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the operator AccAddrss who have the right to invoke CreateCoins(),
BindProxyHash() and BindAssetHash(),
Example:
$ %s query crosschain operator
`,
				version.ClientName,
			),
		),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := common.QueryOperator(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			var operator sdk.AccAddress
			cdc.MustUnmarshalJSON(res, &operator)
			fmt.Printf("operator: %s\n", operator.String())
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}
