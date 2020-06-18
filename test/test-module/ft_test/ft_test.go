package ft_test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/gaia/test/test-module"
	"github.com/polynetwork/cosmos-poly-module/ft"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gaia/app"
	"github.com/polynetwork/cosmos-poly-module/headersync/poly-utils/common"
	"github.com/tendermint/tendermint/crypto"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

func setupFt() (crypto.PrivKey, types.AccAddress) {
	fromPriv, fromAddr, err := test.GetCosmosPrivateKey(test.NetConfig["gaia"].OperatorWallet, []byte(test.NetConfig["gaia"].OperatorPwd))
	if err != nil {
		panic(fmt.Sprintf("GetCosmosPrivateKey error:%v", err))
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))
	return fromPriv, fromAddr
}

func Test_ft_MsgCreateDenom(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupFt()
	creator := fromAddr

	msg := ft.NewMsgCreateDenom(creator, "")
	if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_ft_MsgCreateAndDelegateCoinToProxy(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()
	fromPriv, fromAddr := setupFt()
	creator := fromAddr
	coins, err := sdk.ParseCoin("100000000000000000000000000ether3")
	if err != nil {
		t.Errorf("ParseCoin error:%v", err)
		return
	}
	msg := ft.NewMsgCreateCoinAndDelegateToProxy(creator, coins, test.HashMap["gaia"].ProxyHash[5])
	if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_ft_MsgBindAssetHash(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupFt()
	msgBindAssetParams := []struct {
		denom       string
		toChainId   uint64
		toAssetHash []byte
	}{
		{"oep41", 2, []byte{1, 2, 3}},
		{"oep41", 3, []byte{1, 2, 4}},
	}
	for _, msgBindAssetParam := range msgBindAssetParams {
		msg := ft.NewMsgBindAssetHash(fromAddr, msgBindAssetParam.denom, msgBindAssetParam.toChainId, msgBindAssetParam.toAssetHash)
		if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}

func Test_ft_Lock(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
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
		msg := ft.NewMsgLock(fromAddr, toChainIdAddr.Denom, toChainIdAddr.ToChainId, toChainIdAddr.ToAddr, toChainIdAddr.Amount)
		if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}

func Test_ft_CreateCoins(t *testing.T) {
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()

	fromPriv, fromAddr := setupFt()

	coinsToBeCreated := "100000000mst,100000000msu,100000000msv"
	msg := ft.NewMsgCreateCoins(fromAddr, coinsToBeCreated)
	if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}
