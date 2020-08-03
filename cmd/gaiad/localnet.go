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
	"github.com/polynetwork/cosmos-poly-module/ccm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"io/ioutil"
	"os"
	"strconv"
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
		Use:   "config-genesis-file [[coin][,[coin]] [account_name]  [wallet_path] [chain_id]",
		Short: "config genesis account and chainId to genesis.json",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			coins, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			//flagHome := viper.GetString(flags.FlagHome)
			acctName := args[1]
			walletPath := args[2]
			chainId, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			acctInfo, err := importAcctFromWalletFile(cmd, acctName, walletPath, viper.GetString(flagClientHome))
			if err != nil {
				return err
			}
			return addGenesisAcct(cmd, coins, acctInfo, cdc, config, chainId)
		},
	}
	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	cmd.Flags().String(flagVestingAmt, "", "amount of coins for vesting accounts")
	cmd.Flags().Uint64(flagVestingStart, 0, "schedule start time (unix epoch) for vesting accounts")
	cmd.Flags().Uint64(flagVestingEnd, 0, "schedule end time (unix epoch) for vesting accounts")
	return cmd
}

func importAcctFromWalletFile(cmd *cobra.Command, acctName string, walletPath string, defaultClientHome string) (keys.Info, error) {

	//inBuf := bufio.NewReader(cmd.InOrStdin())
	//pwd, err := input.GetString(
	//	"Enter your account passphrase: ", inBuf)
	//if err != nil {
	//	return nil, fmt.Errorf("GetPassparase error: %v", err)
	//}
	//
	//bz, err := ioutil.ReadFile(walletPath)
	//if err != nil {
	//	return nil, err
	//}
	//privKey, _, err := mintkey.UnarmorDecryptPrivKey(string(bz), string(pwd))
	//if err != nil {
	//	return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	//}
	// restore key.Info from private key

	buf := bufio.NewReader(cmd.InOrStdin())
	kb, err := keys.NewKeyring(sdk.KeyringServiceName(), "file", defaultClientHome, buf)
	if err != nil {
		return nil, err
	}

	bz, err := ioutil.ReadFile(walletPath)
	if err != nil {
		return nil, err
	}

	passphrase, err := input.GetPassword("Enter passphrase to decrypt your key:", buf)
	if err != nil {
		return nil, err
	}
	if err := kb.ImportPrivKey(acctName, string(bz), passphrase); err != nil {
		return nil, fmt.Errorf("kb.ImportPrivKey error: %v", err)
	}
	acctInfo, err := kb.Get(acctName)
	if err != nil {
		return nil, fmt.Errorf("kb.Get(%s), error: %v", acctName, err)
	}
	return acctInfo, nil

	//if err := tmos.EnsureDir(readPath, 0700); err != nil {
	//	return nil, err
	//}
	//file := filepath.Join(readPath, fmt.Sprintf("%v.json", "key_seed"))
	//contents, err := tmos.ReadFile(file)
	//if err != nil {
	//	return nil, err
	//}
	//info := map[string]string{}
	//if err := json.Unmarshal(contents, &info); err != nil {
	//	return nil, err
	//}
	//seed, ok := info["secret"]
	//if !ok {
	//	return nil, fmt.Errorf("seed_key not created")
	//}
	//inBuf := bufio.NewReader(cmd.InOrStdin())
	//kb, err := keys.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), inBuf)
	//info, err := kb.CreateAccount("name", seed, "123454678", ckeys.DefaultKeyPass, hdPath, algo)
	//
	//return nil, nil
}

func generateNewAcct(cmd *cobra.Command, clientHome string) (keys.Info, error) {
	kb, err := ckeys.NewKeyBaseFromDir(clientHome)
	if err != nil {
		return nil, fmt.Errorf("NewKeyBaseFromDir error:%v", err)
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
		return nil, err
	}

	if keyPass == "" {
		keyPass = ckeys.DefaultKeyPass
	}

	_, secret, err := server.GenerateSaveCoinKey(kb, clientHome, keyPass, true)
	if err != nil {
		_ = os.RemoveAll(clientHome)
		return nil, err
	}

	info := map[string]string{"secret": secret}

	cliPrint, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	// save private key seed words
	if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), clientHome, cliPrint); err != nil {
		return nil, err
	}
	acctInfo, err := kb.Get(clientHome)
	if err != nil || acctInfo == nil {
		return nil, fmt.Errorf("keybase.Get(%s), error:%v, acountInfo is %+v\n", err, acctInfo)
	}
	return acctInfo, nil
}

func addGenesisAcct(cmd *cobra.Command, coins sdk.Coins, acctInfo keys.Info, cdc *codec.Codec, config *tmconfig.Config, chainId uint64) error {
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

	var ccmGenesis ccm.GenesisState
	if err := app.MakeCodec().UnmarshalJSON(appState[ccm.ModuleName], &ccmGenesis); err != nil {
		return fmt.Errorf("UnmarshalJSON, ccm genesis bytes to auth.GenesisState without accounts failed")
	}
	ccmGenesis.Params.ChainIdInPolyNet = chainId
	appCcmStateJson, err := codec.MarshalJSONIndent(cdc, ccmGenesis)
	if err != nil {
		return err
	}
	appState[ccm.ModuleName] = appCcmStateJson

	appStateJSON, err := cdc.MarshalJSON(appState)
	if err != nil {
		return err
	}

	// export app state
	genDoc.AppState = appStateJSON

	return genutil.ExportGenesisFile(genDoc, genFile)
}
