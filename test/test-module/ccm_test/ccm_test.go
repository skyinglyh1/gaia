package ccm_test

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/test/test-module"
	"github.com/polynetwork/cosmos-poly-module/ccm"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"testing"
)

func Test_ccm_MsgProcessCrossChainTx(t *testing.T) {
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

	proof := "aa200ba64a2426582d280bd25504fcd1a6c3ef6e032ae3d7bf2288e713c242d433a9030000000000000020c8ecde6ea88c60452c5db3d40e50f72af1cc3374f3996e85af3bcb3d11c284de08090900000000000014b7f398711664de1dd685d9ba3eee3b6b830a7d83070000000000000014c330431496364497d7257839737b5e4596f5ac0606756e6c6f636b1d14e931a4f7020caaacf3ce942567625ebbc0a0ab350a00000000000000"
	header164599, _ := hex.DecodeString("000000000000000000000000398bc72374e0e0a89d4ec6a738de438e7607658bff7edffa02d9e2ac1b864bc4c272a1ca6504bd4029fde4db321d55a4d07b618f8b34ad6f21f88f168ec7564994e785bea32490f28cab67304b10d3318b01835da6abee1ecfdb2e2d7b5555a6287f5551c9425aa9d4e4c18dd53e8dd513abe8eaeb32da22c6ffc89f4afca69b0909ce5ef78202007fa61c13e95eeca6fd11017b226c6561646572223a332c227672665f76616c7565223a22424c364d7467636336776463386c377752466e3673414e694c7265444c52552f6471546d7645496c4e484a4c4a6b714f73424a6f79546a67584e57684b617a53414c4c4c435870644932486637684865326a41615a4b6f3d222c227672665f70726f6f66223a225447332b49775a4a4a57594565495333436f73316657446f5243726435793478413241662f684e36635763537a567135476a312f694469644178746571743776696d756d39576571526f554b48594f6d766d536835413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a3132303030302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000523120502482acb6564b19b90653f6e9c806292e8aa83f78e7a9382a24a6efe41c0c06f3923120502eb1baab602c5899282561cdaaa7aabbcdd0ccfcbc3e79793ac24acf90778f35a23120502679930a42aaf3c69798ca8a3f12e134c019405818d783d11748e039de851598823120502468dd1899ed2d1cc2b829882a165a0ecb6a745af0c72eb2982d66b4311b4ef73231205031e0779f5c5ccb2612352fe4a200f99d3e7758e70ba53f607c59ff22a30f678ff0542011b96e7a2f93a1325c05f57878df21e68c0911f2f8c3c6101b99d09e22a145aa0966602923be49be741cc70976b94d10dd5908c74adf1f75f653e5b15e6f748e13a42011bf63f54bfda2185a7ddb9675a5a0f8c6272459353b0b50679e7296f8a8c218dad099a8a06dac8c934d48a665f1913e39dec76d70cb226f17701b3285050b895b042011b7cc1597550b0fded2ed21430eb1ca0519681653c1d21772c9b92b9c7894ac8b36907f0533105ba760660797d225a7b17d20aa965a8ea36764469c88b53c7b8cf42011b29d86def125e4d8377ad6e7cea3c40a0b414f064bcd75d45b82bf391c1b549535c027e56489ca8599702fb71b7f635523110dac7cf3f8a8a0643e2339cb171b342011cea7fce5152da3aedc658cd4ef2f5db45d4421e9c3f906c020825b1c3bd2694a93f1603f3c744c270e888d608ed87ed2a43b972875cb0519fd40aabbf93517c37")
	processCrossChainTxs := []struct {
		FromChainId      uint64
		RelayChainHeight uint32
		Proof            string
		header           []byte
	}{
		{3, 164599, proof, header164599},
	}
	for _, processCrossChainTx := range processCrossChainTxs {
		msg := ccm.NewMsgProcessCrossChainTx(fromAddr, processCrossChainTx.FromChainId, processCrossChainTx.RelayChainHeight, processCrossChainTx.Proof, processCrossChainTx.header)
		if err := test.SendMsg(client, fromAddr, fromPriv, appCdc, msg, test.NetConfig["gaia"].ChainId); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}
