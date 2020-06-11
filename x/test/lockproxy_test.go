package test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/gaia/x/lockproxy"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/x/headersync/poly-utils/common"
	"github.com/tendermint/tendermint/crypto"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

func setupLockProxy() (crypto.PrivKey, types.AccAddress) {
	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		panic(fmt.Sprintf("GetCosmosPrivateKey error:%v", err))
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))
	return fromPriv, fromAddr
}

func Test_lockproxy_MsgCreateLockProxy(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupLockProxy()
	creator := fromAddr

	msg := lockproxy.NewMsgCreateLockProxy(creator)
	fmt.Printf("creator: %s created lockproxy hash :%x\n", creator.String(), creator.Bytes())
	//creator: cosmos1vnnpptmw2vlm5h06ej3t23vsx6jaqgtcwexesm created lockproxy hash :64e610af6e533fba5dfacca2b5459036a5d02178
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}

}

func Test_lockproxy_MsgBindProxyHash(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupLockProxy()
	creator := fromAddr

	bindProxyHashs := []struct {
		toChainId   uint64
		toProxyHash []byte
	}{
		{3, proxyInOntHashDev},
		{2, proxyInEthHashDev},
	}
	for _, bindProxyHash := range bindProxyHashs {
		msg := lockproxy.NewMsgBindProxyHash(creator, bindProxyHash.toChainId, bindProxyHash.toProxyHash)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}

func Test_lockproxy_MsgBindAssetHash(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupLockProxy()
	msgBindAssetParams := []struct {
		denom       string
		toChainId   uint64
		toAssetHash []byte
		initialAmt  sdk.Int
	}{
		{"peo", 2, oep4denInEthDev, sdk.NewInt(10000000000000)},
		{"peo", 3, oep4denInOntDev, sdk.NewInt(10000000000000)},
	}
	for _, msgBindAssetParam := range msgBindAssetParams {
		msg := lockproxy.NewMsgBindAssetHash(fromAddr, msgBindAssetParam.denom, msgBindAssetParam.toChainId, msgBindAssetParam.toAssetHash, msgBindAssetParam.initialAmt)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}

func Test_lockproxy_Lock(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupLockProxy()

	//toEthAddr, _ := hex.DecodeString("5cD3143f91a13Fe971043E1e4605C1c23b46bF44")
	toOntAddr, _ := common.AddressFromBase58("AQf4Mzu1YJrhz9f3aRkkwSm9n3qhXGSh4p")
	toChainIdAddrs := []struct {
		Denom     string
		ToChainId uint64
		ToAddr    []byte
		Amount    sdk.Int
	}{
		//{"peo", 2, toEthAddr, sdk.NewInt(11000)},
		{"peo", 3, toOntAddr[:], sdk.NewInt(1)},
	}
	for _, toChainIdAddr := range toChainIdAddrs {
		msg := lockproxy.NewMsgLock(fromAddr, fromAddr.Bytes(), toChainIdAddr.Denom, toChainIdAddr.ToChainId, toChainIdAddr.ToAddr, toChainIdAddr.Amount)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}
