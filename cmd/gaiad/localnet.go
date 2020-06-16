package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/input"
	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/gaia/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"os"
)

const (
	flagClientHome   = "home-client"
	flagVestingStart = "vesting-start-time"
	flagVestingEnd   = "vesting-end-time"
	flagVestingAmt   = "vesting-amount"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec,
	defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account  [coin][,[coin]]",
		Short: "Add genesis account to genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			//flagHome := viper.GetString(flags.FlagHome)
			clientHome := defaultClientHome

			//readPath := filepath.Join(clientHome)
			//if err := tmos.EnsureDir(readPath, 0700); err != nil {
			//	return err
			//}
			//file := filepath.Join(readPath, fmt.Sprintf("%v.json", "key_seed"))
			//contents, err := tmos.ReadFile(file)
			//if err != nil {
			//	return err
			//}
			//info := map[string]string{}
			//if err := json.Unmarshal(contents, &info); err != nil {
			//	return err
			//}
			//seed, ok := info["secret"]
			//if !ok {
			//	return fmt.Errorf("seed_key not created")
			//}
			//inBuf := bufio.NewReader(cmd.InOrStdin())
			//kb, err := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), inBuf)
			//info, err := kb.CreateAccount("name", seed, "123454678", ckeys.DefaultKeyPass, hdPath, algo)
			//

			kb, err := ckeys.NewKeyBaseFromDir(clientHome)
			if err != nil {
				return fmt.Errorf("NewKeyBaseFromDir error:%v", err)
			}

			buf := bufio.NewReader(cmd.InOrStdin())
			prompt := fmt.Sprintf(
				"Password for account '%s' (default %s):", clientHome, ckeys.DefaultKeyPass,
			)

			keyPass, err := input.GetPassword(prompt, buf)
			if err != nil && keyPass != "" {
				// An error was returned that either failed to read the password from
				// STDIN or the given password is not empty but failed to meet minimum
				// length requirements.
				return err
			}

			if keyPass == "" {
				keyPass = ckeys.DefaultKeyPass
			}

			_, secret, err := server.GenerateSaveCoinKey(kb, clientHome, keyPass, true)
			if err != nil {
				_ = os.RemoveAll(clientHome)
				return err
			}

			info := map[string]string{"secret": secret}

			cliPrint, err := json.Marshal(info)
			if err != nil {
				return err
			}

			// save private key seed words
			if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), clientHome, cliPrint); err != nil {
				return err
			}
			acctInfo, err := kb.Get(clientHome)
			if err != nil || acctInfo == nil {
				return fmt.Errorf("keybase.Get(%s), error:%v, acountInfo is %+v\n", err, acctInfo)
			}
			return addGenesisAcct(cmd, coins, acctInfo, cdc, config)
		},
	}
	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Uint64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Uint64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")
	return cmd
}

func addGenesisAcct(cmd *cobra.Command, coins sdk.Coins, acctInfo keys.Info, cdc *codec.Codec, config *tmconfig.Config) error {
	var accts exported.GenesisAccounts

	accts = append(accts, &auth.BaseAccount{
		Address:       acctInfo.GetAddress(),
		Coins:         coins,
		PubKey:        acctInfo.GetPubKey(),
		AccountNumber: 0,
		Sequence:      0,
	})

	// retrieve the app state
	genFile := config.GenesisFile()
	appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
	if err != nil {
		return err
	}
	var authGenesis auth.GenesisState
	if err := app.MakeCodec().UnmarshalJSON(appState[auth.ModuleName], &authGenesis); err != nil {
		return fmt.Errorf("UnmarshalJSON, auth genesis bytes to auth.GenesisState without accounts failed")
	}
	authGenesis.Accounts = accts

	appAuthStateJson, err := codec.MarshalJSONIndent(cdc, authGenesis)
	if err != nil {
		return err
	}

	appState[auth.ModuleName] = appAuthStateJson

	appStateJSON, err := cdc.MarshalJSON(appState)
	if err != nil {
		return err
	}

	// export app state
	genDoc.AppState = appStateJSON

	return genutil.ExportGenesisFile(genDoc, genFile)
}
