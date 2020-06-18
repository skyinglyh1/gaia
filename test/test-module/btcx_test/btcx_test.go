package btcx_test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/gaia/test/test-module"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gaia/app"
	"github.com/polynetwork/cosmos-poly-module/btcx"
	"github.com/tendermint/tendermint/crypto"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

func setupBtcx() (crypto.PrivKey, types.AccAddress) {
	fromPriv, fromAddr, err := test.GetCosmosPrivateKey(test.NetConfig["gaia"].OperatorWallet, []byte(test.NetConfig["gaia"].OperatorPwd))
	if err != nil {
		panic(fmt.Sprintf("GetCosmosPrivateKey error:%v", err))
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))
	return fromPriv, fromAddr
}

func Test_btcx_MsgCreateCoin(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupBtcx()
	creator := fromAddr
	msg := btcx.NewMsgCreateDenom(creator, "btc3", test.RedeemScriptStr)
	if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_btcx_MsgBindAssetHash(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()
	btcDenom := "btc3"
	fromPriv, fromAddr := setupBtcx()
	msgBindAssetParams := []struct {
		denom       string
		toChainId   uint64
		toAssetHash []byte
	}{
		{btcDenom, 1, test.HashMap["gaia"].BtcHash[btcDenom][1]},
		//{btcDenom, 2, test.HashMap["gaia"].BtcHash[btcDenom][2]},
		//{btcDenom, 3, test.HashMap["gaia"].BtcHash[btcDenom][3]},
		//{btcDenom, 4, test.HashMap["gaia"].BtcHash[btcDenom][4]},
		//{btcDenom, 5, test.HashMap["gaia"].BtcHash[btcDenom][5]},
	}
	for _, msgBindAssetParam := range msgBindAssetParams {
		msg := btcx.NewMsgBindAssetHash(fromAddr, msgBindAssetParam.denom, msgBindAssetParam.toChainId, msgBindAssetParam.toAssetHash)
		if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}

func Test_btcx_Lock(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := test.GetCosmosPrivateKey(test.NetConfig["gaia"].OperatorWallet, []byte(test.NetConfig["gaia"].OperatorPwd))
	if err != nil {
		panic(fmt.Sprintf("GetCosmosPrivateKey error:%v", err))
	}
	fmt.Printf("sender address bench32 is %s\n", fromAddr.String())
	//toEthAddr, _ := hex.DecodeString("5cD3143f91a13Fe971043E1e4605C1c23b46bF44")
	//toOntAddr, _ := common.AddressFromBase58("AQf4Mzu1YJrhz9f3aRkkwSm9n3qhXGSh4p")
	toChainIdAddrs := []struct {
		Denom     string
		ToChainId uint64
		ToAddr    []byte
		Amount    sdk.Int
	}{
		{"node0token", 1, []byte("mpCNjy4QYAmw8eumHJRbVtt6bMDVQvPpFn"), sdk.NewInt(10000)},
		//{"btc", 2, toEthAddr, sdk.NewInt(11000)},
		//{"btc", 3, toOntAddr[:], sdk.NewInt(2)},
	}
	for _, toChainIdAddr := range toChainIdAddrs {
		msg := btcx.NewMsgLock(fromAddr, toChainIdAddr.Denom, toChainIdAddr.ToChainId, toChainIdAddr.ToAddr, toChainIdAddr.Amount)
		if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}
