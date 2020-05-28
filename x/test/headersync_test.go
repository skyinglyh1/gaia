package test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/x/headersync"
	mcc "github.com/ontio/multi-chain/common"
	mctype "github.com/ontio/multi-chain/core/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"testing"
)

func Test_Deserialize_GenesisHeader(t *testing.T) {
	genesisHeader := &mctype.Header{}
	heade0Bs, _ := hex.DecodeString(header0)
	source := mcc.NewZeroCopySource(heade0Bs)
	if err := genesisHeader.Deserialization(source); err != nil {
		t.Errorf("Deserialize...... err:%+v", err)
	}
}

func Test_headersync_MsgSyncGenesisHeader(t *testing.T) {
	//client := rpchttp.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	heade0Bs, _ := hex.DecodeString(header0)
	msg := headersync.NewMsgSyncGenesisParam(fromAddr, heade0Bs)
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_headersync_MsgSyncBlockHeaders(t *testing.T) {
	//client := rpchttp.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	heade1Bs, _ := hex.DecodeString(header1)
	heade2Bs, _ := hex.DecodeString(header2)
	msg := headersync.NewMsgSyncHeadersParam(fromAddr, [][]byte{heade1Bs, heade2Bs})
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}
