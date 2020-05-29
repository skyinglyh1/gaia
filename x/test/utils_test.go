package test

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/x/crosschain"
	"github.com/davecgh/go-spew/spew"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"reflect"
	"testing"
)

func Test_Transfer_StakeCoin_From_Operator(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(validatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))
	receivers := []string{
		user0Addr,
		"cosmos1v8rpqa4valmgx5a8gnnstsecv8xg76sgzw4820",
	}
	amt, _ := sdk.ParseCoins("1000000stake")
	for _, receiverAddrStr := range receivers {
		toAddr, _ := sdk.AccAddressFromBech32(receiverAddrStr)
		msg := bank.MsgSend{
			FromAddress: fromAddr,
			ToAddress:   toAddr,
			Amount:      amt,
		}
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}

func Test_ParseCoins(t *testing.T) {

	args0 := "1000000000ont,1000000000000000000ong"
	//args1 := "1000000000stake,1000000000validatortoken"
	coins, err := sdk.ParseCoins(args0)
	if err != nil {
		t.Errorf("parsecoins error:%v", err)
	}

	spew.Printf("coins are %v\n", coins)
}

func Test_UnmarshalOperator(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("cosmos1ayc6faczpj42eu7wjsjkwcj7h0q2p2e4vrlkzf")
	fmt.Printf("addr in hex is %x\n", addr.Bytes())

	addr1, err := sdk.AccAddressFromHex(hex.EncodeToString(addr))
	if err != nil {
		t.Errorf("could not unmarshal result to sdk.AccAddress:%v", err)
	}
	fmt.Printf("opeartor are %s\n", addr1.String())
	fromPriv, fromAddr, err := GetCosmosPrivateKey(user0Wallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("acct in Hash format is = %x\n", fromAddr.Bytes())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

}

func Test_DenomToHash(t *testing.T) {
	denom := "ontc"
	fmt.Printf("denom:%s, sourceasset hash is %x\n", denom, crosschain.DenomToHash(denom).Bytes())
}

func Test_TypeOf(t *testing.T) {
	var x crosschain.MsgProcessCrossChainTx
	fmt.Printf("xtype is %s\n", reflect.TypeOf(x).String())
	var x1 *crosschain.MsgProcessCrossChainTx
	fmt.Printf("xtype is %s\n", reflect.TypeOf(x1).String())
	//fromContractHash := make([]byte, 0)
	var fromContractHash []byte
	sourceAssetHash := []byte{1, 2, 3, 5}
	copy(fromContractHash, sourceAssetHash)
	//fromContractHash = append(fromContractHash, sourceAssetHash...)
	var from sdk.AccAddress
	copy(from[:], sourceAssetHash[:])
	fmt.Printf("fromContractHash = %v\n", fromContractHash)
	fmt.Printf("from = %v\n", from)
}

func Test_TestBytes(t *testing.T) {
	methodBs := make([]byte, 0)
	methodBs = returnBytes()
	fmt.Printf("methodBs = %v\n", methodBs)

}

func returnBytes() []byte {
	res := []byte{1, 2, 3}
	return res
}
