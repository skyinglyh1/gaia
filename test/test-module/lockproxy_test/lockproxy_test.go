package lockproxy_test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/gaia/test/test-module"
	"github.com/polynetwork/cosmos-poly-module/lockproxy"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gaia/app"
	"github.com/tendermint/tendermint/crypto"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

func setupLockProxy() (crypto.PrivKey, types.AccAddress) {
	fromPriv, fromAddr, err := test.GetCosmosPrivateKey(test.NetConfig["gaia"].OperatorWallet, []byte(test.NetConfig["gaia"].OperatorPwd))
	if err != nil {
		panic(fmt.Sprintf("GetCosmosPrivateKey error:%v", err))
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))
	return fromPriv, fromAddr
}

func Test_lockproxy_MsgCreateLockProxy(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupLockProxy()
	creator := fromAddr

	msg := lockproxy.NewMsgCreateLockProxy(creator)
	fmt.Printf("creator: %s created lockproxy hash: %x\n", creator.String(), creator.Bytes())
	//creator: cosmos1vnnpptmw2vlm5h06ej3t23vsx6jaqgtcwexesm created lockproxy hash :64e610af6e533fba5dfacca2b5459036a5d02178
	if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}

}

func Test_lockproxy_MsgBindProxyHash(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupLockProxy()
	creator := fromAddr

	bindProxyHashs := []struct {
		toChainId   uint64
		toProxyHash []byte
	}{
		//{3, test.HashMap["gaia"].ProxyHash[3]},
		{2, test.HashMap["gaia"].ProxyHash[2]},
	}
	for _, bindProxyHash := range bindProxyHashs {
		msg := lockproxy.NewMsgBindProxyHash(creator, bindProxyHash.toChainId, bindProxyHash.toProxyHash)
		if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}

func Test_lockproxy_MsgBindAssetHash(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupLockProxy()
	iniAmt, _ := new(big.Int).SetString("100000000000000000000000000", 10)
	msgBindAssetParams := []struct {
		denom       string
		toChainId   uint64
		toAssetHash []byte
		initialAmt  sdk.Int
	}{
		{"ether3", 2, test.HashMap["gaia"].EthHash["ether3"][2], sdk.NewIntFromBigInt(iniAmt)},
		//{"peo", 3, test.HashMap["gaia"].Oep4Hash["oep4"][3], sdk.NewInt(10000000000000)},
	}
	for _, msgBindAssetParam := range msgBindAssetParams {
		msg := lockproxy.NewMsgBindAssetHash(fromAddr, msgBindAssetParam.denom, msgBindAssetParam.toChainId, msgBindAssetParam.toAssetHash, msgBindAssetParam.initialAmt)
		if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}

func Test_lockproxy_Lock(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupLockProxy()

	toEthAddr, _ := hex.DecodeString("344cFc3B8635f72F14200aAf2168d9f75df86FD3")
	//toOntAddr, _ := common.AddressFromBase58("AQf4Mzu1YJrhz9f3aRkkwSm9n3qhXGSh4p")
	amt, _ := new(big.Int).SetString("100000000000000000", 10)
	toChainIdAddrs := []struct {
		Denom     string
		ToChainId uint64
		ToAddr    []byte
		Amount    sdk.Int
	}{
		{"ether3", 2, toEthAddr, sdk.NewIntFromBigInt(amt)},
		//{"peo", 3, toOntAddr[:], sdk.NewInt(1)},
	}
	for _, toChainIdAddr := range toChainIdAddrs {
		msg := lockproxy.NewMsgLock(fromAddr, fromAddr.Bytes(), toChainIdAddr.Denom, toChainIdAddr.ToChainId, toChainIdAddr.ToAddr, toChainIdAddr.Amount)
		if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}
