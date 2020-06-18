package btcx_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/test/test-module"
	"github.com/polynetwork/cosmos-poly-module/btcx"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"testing"
)

func setupRestReq(restIp string, fromAddr sdk.AccAddress, signerPriv crypto.PrivKey, appCdc *codec.Codec) (*rest.BaseReq, *auth.BaseAccount, error) {
	_, body, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "GET", fmt.Sprintf("/auth/accounts/%s", fromAddr.String()), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("setupRestReq, SendRequest to get auth.account error:%v", err)
	}
	var resp rest.ResponseWithHeight
	appCdc.MustUnmarshalJSON(body, &resp)
	var acc auth.BaseAccount
	auth.ModuleCdc.MustUnmarshalJSON(resp.Result, &acc)
	accnum := acc.GetAccountNumber()
	sequence := acc.GetSequence()

	fees := sdk.NewCoins(sdk.NewInt64Coin("stake", 1))
	baseReq := rest.NewBaseReq(
		fromAddr.String(), "", test.ChainID, "200000", fmt.Sprintf("%f", 1.0), accnum, sequence, fees, nil, false,
	)
	return &baseReq, &acc, nil
}

func Test_btcx_GetDenomInfo(t *testing.T) {
	_, body, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "GET", fmt.Sprintf("/btcx/denom_info/%s", "btc2"), nil)
	if err != nil {
		t.Errorf("GetDenomInfo, SendRequest Error:%v", err)
	}
	cdc := app.MakeCodec()
	var resp rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &resp)
	var denomInfo btcx.DenomInfo
	cdc.MustUnmarshalJSON(resp.Result, &denomInfo)
	fmt.Printf("denomInfo is %s", denomInfo.String())
}

func Test_btcx_GetDenomInfoWithChainId(t *testing.T) {
	_, body, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "GET", fmt.Sprintf("/btcx/denom_cc_info/%s/%s", "btc3", "1"), nil)
	if err != nil {
		t.Errorf("GetDenomInfo, SendRequest Error:%v", err)
	}

	cdc := app.MakeCodec()
	var resp rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &resp)

	var denomInfo btcx.DenomCrossChainInfo
	cdc.MustUnmarshalJSON(resp.Result, &denomInfo)
	fmt.Printf("denomCrossChainInfoId is %s", denomInfo.String())
}

func Test_btcxC_CreateCoinThroughRest(t *testing.T) {
	fromPriKey, fromAddr := setupBtcx()
	cdc := app.MakeCodec()
	baseReq, acc, err := setupRestReq(test.NetConfig["gaia"].RestIp, fromAddr, fromPriKey, cdc)
	sr := btcx.CreateCoinReq{
		BaseReq:      *baseReq,
		Denom:        "btct1",
		RedeemScript: "1234",
	}
	req, err := cdc.MarshalJSON(sr)
	require.NoError(t, err)
	// generate msg to be signed
	_, payload, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "POST", fmt.Sprintf("/btcx/create_coin"), req)
	if err != nil {
		t.Errorf("CreateCoinThroughRest, post create_coin to get msg to sign Error:%v", err)
	}

	// sign and broadcast
	_, body, err := test.SignAndBroadcastGenTx(test.NetConfig["gaia"].RestIp, "op", test.NetConfig["gaia"].OperatorPwd, payload, *acc, flags.DefaultGasAdjustment, false, cdc)
	if err != nil {
		t.Errorf("CreateCoinThroughRest, SignAndBroadcastGenTx Error:%v", err)
	}

	var txResp sdk.TxResponse
	err = cdc.UnmarshalJSON([]byte(body), &txResp)
	require.NoError(t, err)
	fmt.Printf("txResp is %+v\n", txResp)
}
