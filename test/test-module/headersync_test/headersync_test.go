package headersync_test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/test/test-module"
	"github.com/polynetwork/cosmos-poly-module/headersync"
	polycommon "github.com/polynetwork/cosmos-poly-module/headersync/poly-utils/common"
	polytype "github.com/polynetwork/cosmos-poly-module/headersync/poly-utils/core/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"testing"
)

func Test_Deserialize_GenesisHeader(t *testing.T) {
	genesisHeader := &polytype.Header{}
	heade0Bs, _ := hex.DecodeString(test.Header0)
	source := polycommon.NewZeroCopySource(heade0Bs)
	if err := genesisHeader.Deserialization(source); err != nil {
		t.Errorf("Deserialize...... err:%+v", err)
	}
}

func Test_headersync_MsgSyncGenesisHeader(t *testing.T) {
	//client := rpchttp.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := test.GetCosmosPrivateKey(test.NetConfig["gaia"].OperatorWallet, []byte(test.NetConfig["gaia"].OperatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	heade0Bs, _ := hex.DecodeString(test.Header0)
	msg := headersync.NewMsgSyncGenesisParam(fromAddr, heade0Bs)
	if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_headersync_MsgSyncBlockHeaders(t *testing.T) {
	//client := rpchttp.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
	client, err := rpchttp.NewHTTP(test.NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := test.GetCosmosPrivateKey(test.NetConfig["gaia"].OperatorWallet, []byte(test.NetConfig["gaia"].OperatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	heade1Bs, _ := hex.DecodeString(test.Header1)
	heade2Bs, _ := hex.DecodeString(test.Header2)
	msg := headersync.NewMsgSyncHeadersParam(fromAddr, [][]byte{heade1Bs, heade2Bs})
	if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}
