package test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/x/btcx"
	"github.com/tendermint/tendermint/crypto"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

func setupBtcx() (crypto.PrivKey, types.AccAddress) {
	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		panic(fmt.Sprintf("GetCosmosPrivateKey error:%v", err))
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))
	return fromPriv, fromAddr
}

func Test_btcx_MsgCreateCoin(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupBtcx()
	creator := fromAddr

	msg := btcx.NewMsgCreateCoin(creator, "btc", RedeemScriptStr)
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_btcx_MsgBindAssetHash(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupBtcx()
	msgBindAssetParams := []struct {
		denom       string
		toChainId   uint64
		toAssetHash []byte
	}{
		{"btc", 1, btcHashInBtcDev},
		{"btc", 2, btcHahInEthDev},
		{"btc", 3, btcHahInOntDev},
	}
	for _, msgBindAssetParam := range msgBindAssetParams {
		msg := btcx.NewMsgBindAssetParam(fromAddr, msgBindAssetParam.denom, msgBindAssetParam.toChainId, msgBindAssetParam.toAssetHash)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}

func Test_btcx_Lock(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(user0Wallet, []byte(operatorPwd))
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
		{"btc", 1, []byte("mpCNjy4QYAmw8eumHJRbVtt6bMDVQvPpFn"), sdk.NewInt(10000)},
		//{"btc", 2, toEthAddr, sdk.NewInt(11000)},
		//{"btc", 3, toOntAddr[:], sdk.NewInt(2)},
	}
	for _, toChainIdAddr := range toChainIdAddrs {
		msg := btcx.NewMsgLock(fromAddr, toChainIdAddr.Denom, toChainIdAddr.ToChainId, toChainIdAddr.ToAddr, &toChainIdAddr.Amount)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}
