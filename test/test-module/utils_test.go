package test

import (
	"encoding/hex"
	"fmt"
	. "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/gaia/app"
	"github.com/davecgh/go-spew/spew"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"testing"
)

func Test_Transfer_StakeCoin_From_Operator(t *testing.T) {
	client, err := rpchttp.NewHTTP(ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey("./wallets/172.168.3.94_node0", []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))
	_, operator, _ := GetCosmosPrivateKey("./wallets/operator", []byte(operatorPwd))
	fmt.Printf("operator = %v\n", operator.String())
	receivers := []string{
		//operator.String(),
		//user0Addr,
		"cosmos1nztkr7cvp6cvq4s9apyu4emayw0e3trl68gj3f",
	}
	amt, _ := sdk.ParseCoins("1000stake")
	for _, receiverAddrStr := range receivers {
		toAddr, _ := sdk.AccAddressFromBech32(receiverAddrStr)
		msg := bank.MsgSend{
			FromAddress: fromAddr,
			ToAddress:   toAddr,
			Amount:      amt,
		}
		if err := SendMsg(client, fromAddr, fromPriv, appCdc, msg, NetConfig["gaia"].ChainId); err != nil {
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

func Test_TestBytes(t *testing.T) {
	methodBs := make([]byte, 0)
	methodBs = returnBytes()
	fmt.Printf("methodBs = %v\n", methodBs)

}

func returnBytes() []byte {
	res := []byte{1, 2, 3}
	return res
}

func Test_CheckTxSuccessful(t *testing.T) {
	txHash := "35249D74B42585FE080BD9A49545919F257038BDC95DF8CB709E082CB71B3FF1"
	client, err := rpchttp.NewHTTP("tcp://13.251.218.38:30002", "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	cliCtx := NewCLIContext().WithCodec(app.MakeCodec()).WithClient(client).WithFrom("cosmos1vnnpptmw2vlm5h06ej3t23vsx6jaqgtcwexesm").WithTrustNode(true)
	output, err := utils.QueryTx(cliCtx, txHash)
	if err != nil {
		fmt.Printf("QueryTx err:%v", err)
	}
	if output.Code != 0 {
		fmt.Printf("Tx:%s failed \n result:%s\n", txHash, output.String())

	} else {
		fmt.Printf("Tx:%s Success \n result:%s\n", txHash, output.String())
	}
}

func Test_convertHex(t *testing.T) {
	strs := []string{"reth", "reth3", "btc2", "ether3"}
	for _, str := range strs {
		ha := hex.EncodeToString([]byte(str))
		fmt.Printf("hex(%s) = %s, len = %d\n", str, ha, len(ha))
	}
}
