package test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/gaia/x/ft"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/x/headersync/poly-utils/common"
	"github.com/tendermint/tendermint/crypto"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

func setupFt() (crypto.PrivKey, types.AccAddress) {
	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		panic(fmt.Sprintf("GetCosmosPrivateKey error:%v", err))
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))
	return fromPriv, fromAddr
}

func Test_ft_MsgCreateDenom(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupFt()
	creator := fromAddr

	msg := ft.NewMsgCreateDenom(creator, "oepindependent")
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_ft_MsgCreateAndDelegateCoinToProxy(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupFt()
	creator := fromAddr
	coins, err := sdk.ParseCoin("10000000000000peo")
	if err != nil {
		t.Errorf("ParseCoin error:%v", err)
		return
	}
	msg := ft.NewMsgCreateAndDelegateCoinToProxy(creator, coins)
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_ft_MsgBindAssetHash(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupFt()
	msgBindAssetParams := []struct {
		denom       string
		toChainId   uint64
		toAssetHash []byte
	}{
		{"oepindependent", 2, oep4IndInEthDev},
		{"oepindependent", 3, oep4IndInOntDev},
	}
	for _, msgBindAssetParam := range msgBindAssetParams {
		msg := ft.NewMsgBindAssetHash(fromAddr, msgBindAssetParam.denom, msgBindAssetParam.toChainId, msgBindAssetParam.toAssetHash)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}

func Test_ft_Lock(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupFt()

	toEthAddr, _ := hex.DecodeString("5cD3143f91a13Fe971043E1e4605C1c23b46bF44")
	toOntAddr, _ := common.AddressFromBase58("AQf4Mzu1YJrhz9f3aRkkwSm9n3qhXGSh4p")
	toChainIdAddrs := []struct {
		Denom     string
		ToChainId uint64
		ToAddr    []byte
		Amount    sdk.Int
	}{
		{"btc", 1, []byte("mpCNjy4QYAmw8eumHJRbVtt6bMDVQvPpFn"), sdk.NewInt(10000)},
		{"btc", 2, toEthAddr, sdk.NewInt(11000)},
		{"btc", 3, toOntAddr[:], sdk.NewInt(234)},
	}
	for _, toChainIdAddr := range toChainIdAddrs {
		msg := ft.NewMsgLock(fromAddr, toChainIdAddr.Denom, toChainIdAddr.ToChainId, toChainIdAddr.ToAddr, &toChainIdAddr.Amount)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}

func Test_ft_CreateCoins(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupFt()

	coinsToBeCreated := "100000000MST,100000000MSU,100000000MSV"
	msg := ft.NewMsgCreateCoins(fromAddr, coinsToBeCreated)
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}
