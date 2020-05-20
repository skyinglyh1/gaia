package test

import (
	"encoding/hex"
	"errors"
	"fmt"
	. "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/x/crosschain"
	"github.com/ontio/ontology/common"
	"github.com/tendermint/tendermint/crypto"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"io/ioutil"
	"testing"
	"time"
)

const (
	//ip = "tcp://172.168.3.93:26657"
	ip              = "tcp://172.168.3.93:26657"
	validatorWallet = "./wallets/validator"
	operatorWallet  = "./wallets/operator"
	user0Wallet     = "./wallets/user0"
	operatorPwd     = "12345678"
	ChainID         = "testing"
)

var (
	header0         = "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010ae3a2d1cba9ed56653edab871d93f8a96294debb6169a62681552dfd6d0fc70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c8365b000000001dac2b7c00000000fd1a057b226c6561646572223a343239343936373239352c227672665f76616c7565223a22484a675171706769355248566745716354626e6443456c384d516837446172364e4e646f6f79553051666f67555634764d50675851524171384d6f38373853426a2b38577262676c2b36714d7258686b667a72375751343d222c227672665f70726f6f66223a22785864422b5451454c4c6a59734965305378596474572f442f39542f746e5854624e436667354e62364650596370382f55706a524c572f536a5558643552576b75646632646f4c5267727052474b76305566385a69413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a343239343936373239352c226e65775f636861696e5f636f6e666967223a7b2276657273696f6e223a312c2276696577223a312c226e223a372c2263223a322c22626c6f636b5f6d73675f64656c6179223a31303030303030303030302c22686173685f6d73675f64656c6179223a31303030303030303030302c22706565725f68616e647368616b655f74696d656f7574223a31303030303030303030302c227065657273223a5b7b22696e646578223a312c226964223a2231323035303238313732393138353430623262353132656165313837326132613265336132386439383963363064393564616238383239616461376437646437303664363538227d2c7b22696e646578223a322c226964223a2231323035303338623861663632313065636664636263616232323535326566386438636634316336663836663963663961623533643836353734316366646238333366303662227d2c7b22696e646578223a332c226964223a2231323035303234383261636236353634623139623930363533663665396338303632393265386161383366373865376139333832613234613665666534316330633036663339227d2c7b22696e646578223a342c226964223a2231323035303236373939333061343261616633633639373938636138613366313265313334633031393430353831386437383364313137343865303339646538353135393838227d2c7b22696e646578223a352c226964223a2231323035303234363864643138393965643264316363326238323938383261313635613065636236613734356166306337326562323938326436366234333131623465663733227d2c7b22696e646578223a362c226964223a2231323035303265623162616162363032633538393932383235363163646161613761616262636464306363666362633365373937393361633234616366393037373866333561227d2c7b22696e646578223a372c226964223a2231323035303331653037373966356335636362323631323335326665346132303066393964336537373538653730626135336636303763353966663232613330663637386666227d5d2c22706f735f7461626c65223a5b362c342c332c352c362c312c322c352c342c372c342c322c332c332c372c362c352c342c362c352c312c342c332c312c322c352c322c322c362c312c342c352c342c372c322c332c342c312c352c372c342c312c322c322c352c362c342c342c322c372c332c362c362c352c312c372c332c312c362c312c332c332c322c342c342c312c352c362c352c312c322c362c372c352c362c332c342c372c372c332c322c372c312c352c362c352c322c332c362c322c362c312c372c372c372c312c372c342c332c332c332c322c312c372c355d2c226d61785f626c6f636b5f6368616e67655f76696577223a36303030307d7d9fe171f3fe643eb1c188400b828ba184816fc9ac0000"
	header1         = "0000000000000000000000002a8944a1f753de95e505607d82f70c8a2291c4cf3ef6277b47b282628882c6e400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a39f2294c54fe8e3575b3b3a531f6a91605ef153f8ed5344094c4d67ba489a3f1645a55e0100000036810287538edf27fd0c017b226c6561646572223a332c227672665f76616c7565223a224244714d70463966716a7754596d4e4c6c50795162464e505066546f374261786e677543775563697844713733636856386c6f6e525774474558596d77387671355279372f41505778434d573737433773594c33542b733d222c227672665f70726f6f66223a2246456c6453493638683352532f366f473951676e334835705136443244434c65367966614a485a7169745a6d696b50347a647843354a7943664554662f436e304f474c6a3966736b476a5566523668516370797568673d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000523120502ff7dfc705bb5ac68d2e89230c662999aeb18284131e9fce94ef9faf5b991753d23120503d101838807ec4079a46fe98d6bd9a0690abcbd8ce16e0fbc4520c7c7ef7885db2312050333c44837ddb94ad5fa0eb460b0f49154f90503619d4d2c8ee08330fb815841d22312050344017ccca820d90f0ebb461df46337b0923b0ae2bce583cee1a2624c9208e20823120503c6e681e5154efba64c75b0aa1ca54489bae7653077d16ddd9726f365be30622d0542051c5c15d9e88471206d349fddceb5a48b8b0c261ea9a6b306cf994e56bb4c7ac1846f5f97356268f930ffd37f542f449baaaf897a5fb8be8f6480656c0e6e596b7342051b288d849d6db42b4122de8f781c905b85ceabdcdc7e16a8a2418758d06ceea0895fd9eb9fdd713af03fa9c9593bf930e25dba59331bc27d63cf021ead892de06c42051ba1c6e68c6809c529cdcbc003ab411f000308828f1ffe7bc1d997b3a99eec7a311c27c81291680e4b29cd4bc7ddb2f5c20f615e18512652180d863fb255519a5842051bca95efb2dbff60a538e1726c249df04cd2cd9127f2c55becfefcf51961a5d9072a0d36bbe349befabd34d08ac03535491bfa1513c686e24a316fd7795f54ada042051c9eec1630be76366a2c374ca5ff349d61280a213739c81a872b46b3a8b0df3d2a74c3d03aa2218b3eba4e8d7a617ce49d27fed0d36acb4ce59f4944e0bb5ee3d1"
	header2         = "0000000000000000000000006eb34794b01ee25e05b5e41f35685bbf0f285ef0bb73905fddd72f6bc3c1ca1d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008c05a61780e1873482c24cc98a4fbb53f01c45ee96e2e3afec992097e52def133445a55e0200000097cc879256d3b949fd0c017b226c6561646572223a332c227672665f76616c7565223a22424d63306f455a56513546374c707a436c3962424a4779324b507a2b32776e37636f6c717741745732754846746f58497a3541705142574c35357a30444e5a386f58636f4c4a6d337a526b596b686f2b7a30486648396f3d222c227672665f70726f6f66223a222f7252754d4d592b31534867504c475470325a31703469645743645a7748684646496d3533613534796a79466e5651504b6635397158756d6f383430417859565375696175796c35564e477a6e753456333248356e673d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000623120502ff7dfc705bb5ac68d2e89230c662999aeb18284131e9fce94ef9faf5b991753d23120503c6e681e5154efba64c75b0aa1ca54489bae7653077d16ddd9726f365be30622d2312050312f10239513143c02984bc4ea9d854883fcdf49732dc727dfa747482fc807e642312050344017ccca820d90f0ebb461df46337b0923b0ae2bce583cee1a2624c9208e2082312050333c44837ddb94ad5fa0eb460b0f49154f90503619d4d2c8ee08330fb815841d223120503d101838807ec4079a46fe98d6bd9a0690abcbd8ce16e0fbc4520c7c7ef7885db0642051b15e124c12818c34003566fb5460e9563ee50b20a6ab90dd5511661035d9abbd84cd98115da932e2470854247498105a783546446039c6f6771cecff9c0a59a2d42051bc5fe5c235adcedc9975649804fcbb9fb1b1d68839cfa911389502280500be2e90c9b21f46366774a626f8aee5313c40379fcd71062ae794f155726dbf5cce9ba42051cf31e932fe19b7aee60f04a6786d8b85cd2bcc8571ad375479b7bf4e5dc6a9cc03abad6b70564edf58f9149dd558ac87c3b827bd9d93280784614ab922dfdeac842051c905500ca7a76dba83aa157bc8a725c952486c0bbcfe1c5502ce6e5f96f6ab56a64e0b0b088646f1b4eebb99f6fe708651dd39bc6bbba3145bad82ccd4e40be9842051bce7a5b564889ae0a093160a04a4d01e4a9251c5d84839a2f679b9d387ae94ed432a96b7aaccf7d4a9252d5aa32f76554b24b15ff99df892a2dbed47a6c3edf1942051c585e4851df6fa9626048d46db2a5abfb2b74817fc3163e9fed856ea25ca623ee5dd294e290d5b14599dc7452069f987ce54ddc1bb1a17db687aced9197a14e2f"
	RedeemKey, _    = hex.DecodeString("c330431496364497d7257839737b5e4596f5ac06")
	RedeemScript, _ = hex.DecodeString("552102dec9a415b6384ec0a9331d0cdf02020f0f1e5731c327b86e2b5a92455a289748210365b1066bcfa21987c3e207b92e309b95ca6bee5f1133cf04d6ed4ed265eafdbc21031104e387cd1a103c27fdc8a52d5c68dec25ddfb2f574fbdca405edfd8c5187de21031fdb4b44a9f20883aff505009ebc18702774c105cb04b1eecebcb294d404b1cb210387cda955196cc2b2fc0adbbbac1776f8de77b563c6d2a06a77d96457dc3d0d1f2102dd7767b6a7cc83693343ba721e0f5f4c7b4b8d85eeb7aec20d227625ec0f59d321034ad129efdab75061e8d4def08f5911495af2dae6d3e9a4b6e7aeb5186fa432fc57ae")
)

func Test_SendTxHeaderSync(t *testing.T) {
	//client := rpchttp.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	header1, _ := hex.DecodeString(header0)
	msg := crosschain.NewMsgSyncGenesisParam(fromAddr, header1)
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_CreateCoins(t *testing.T) {
	//client := rpchttp.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	creator := fromAddr
	coins, err := sdk.ParseCoins("1000000000ontc")
	//coins, err := sdk.ParseCoins("1000000000000btcc")
	if err != nil {
		t.Errorf("parse coins err:%v", err)
	}
	msg := crosschain.NewMsgCreateCoins(creator, coins)
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}

}

func Test_SetRedeemScript(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	msg := crosschain.NewMsgSetRedeemScript(fromAddr, "btcc", RedeemKey, RedeemScript)
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_BindProxyHash(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	ontProxy, _ := hex.DecodeString("50478b75da76f14bb8358318b62897b97de043dd")
	//ethProxy, _ := hex.DecodeString("71CF3de5e27EcF7379a8EE74eF32C021dD068d8d")
	proxyHashInOtherChain := []struct {
		ChainId   uint64
		ProxyHash []byte
	}{
		{3, ontProxy},
		//{2, ethProxy},
	}
	for _, proxyhash := range proxyHashInOtherChain {
		msg := crosschain.NewMsgBindProxyParam(fromAddr, proxyhash.ChainId, proxyhash.ProxyHash)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}

func Test_BindAssetHash(t *testing.T) {
	//client := rpchttp.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	coins, _ := sdk.ParseCoins("1000000000ontc")
	//coins, _ := sdk.ParseCoins("1000000000000btcc")

	//btcHashInBtc := RedeemKey
	//btcHashInEth, _ := hex.DecodeString("740C1a496A750a3C3F9A6Ca7e822C6BC776962eA")
	//btcHashInOnt, _ := hex.DecodeString("b7f398711664de1dd685d9ba3eee3b6b830a7d83")
	ontHashInOnt, _ := hex.DecodeString("0000000000000000000000000000000000000001")
	ontHashInEth, _ := hex.DecodeString("2516A471195f020f132af65b6502f2Bd355C553c")
	assetHashInOtherChain := []struct {
		Denom            string
		TargetChainId    uint64
		TargetChainAsset []byte
		limit            sdk.Int
	}{
		//{"btcc", 1, btcHashInBtc, sdk.NewInt(-1)},
		//{"btcc", 2, btcHashInEth, sdk.NewInt(-1)},
		//{"btcc", 3, btcHashInOnt, sdk.NewInt(-1)},
		{"ontc", 2, ontHashInEth, sdk.NewInt(-1)},
		{"ontc", 3, ontHashInOnt, sdk.NewInt(-1)},
	}
	for _, assetHashInfo := range assetHashInOtherChain {
		msg := crosschain.NewMsgBindAssetParam(fromAddr, assetHashInfo.Denom, assetHashInfo.TargetChainId, assetHashInfo.TargetChainAsset, coins.AmountOf(assetHashInfo.Denom))
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}

func Test_Lock(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(user0Wallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	//toEthAddr, _ := hex.DecodeString("5cD3143f91a13Fe971043E1e4605C1c23b46bF44")
	toOntAddr, _ := common.AddressFromBase58("AQf4Mzu1YJrhz9f3aRkkwSm9n3qhXGSh4p")
	toChainIdAddrs := []struct {
		Denom     string
		ToChainId uint64
		ToAddr    []byte
		Amount    sdk.Int
	}{
		//{"btcc", 1, []byte("mpCNjy4QYAmw8eumHJRbVtt6bMDVQvPpFn"), sdk.NewInt(10000)},
		//{"btcc", 2, toEthAddr, sdk.NewInt(11000)},
		//{"btcc", 3, toOntAddr[:], sdk.NewInt(234)},
		{"ontc", 3, toOntAddr[:], sdk.NewInt(20)},
		//{"ontc", 2, toEthAddr, sdk.NewInt(10)},
	}
	for _, toChainIdAddr := range toChainIdAddrs {
		msg := crosschain.NewMsgLock(fromAddr, toChainIdAddr.Denom, toChainIdAddr.ToChainId, toChainIdAddr.ToAddr, &toChainIdAddr.Amount)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}

}
func Test_CheckTxSuccess(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	txHash := "E41D0D509DD7E2FE928471E60A0731AE55CD90BFAE52412F3D590F25C10D2B8B"
	CheckTxSuccessful(client, txHash)
}

//cosmos1ayc6faczpj42eu7wjsjkwcj7h0q2p2e4vrlkzf

func CheckTxSuccessful(client rpchttp.Client, txHash string) {
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
		"cosmos1v8rpqa4valmgx5a8gnnstsecv8xg76sgzw4820",
		"cosmos17ud4tm64emwfrlgq0aafhguxajtc7w4gseapra",
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

//func Test_ProcessCrossChainTx(t *testing.T) {
//	client := rpchttp.NewHTTP(ip, "/websocket")
//	appCdc := app.MakeCodec()
//
//	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
//	if err != nil {
//		t.Errorf("err = %v ", err)
//	}
//	fmt.Printf("acct = %v\n", fromAddr.String())
//	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))
//
//	//msg := new(headersync.MsgSyncGenesisParam)
//	//msg.Syncer = addr
//	//msg.GenesisHeader, _ = hex.DecodeString(genesisHeaderStr)
//
//	msg := crosschain.NewMsgProcessCrossChainTx()
//
//	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
//		t.Errorf("sendMsg error:%v", err)
//	}
//}

func sendMsg(client *rpchttp.HTTP, fromAddr sdk.AccAddress, signerPriv crypto.PrivKey, appCdc *codec.Codec, msg sdk.Msg) error {
	// make sure the account exist in auth module
	bs, err := appCdc.MarshalJSON(auth.NewQueryAccountParams(fromAddr))
	if err != nil {
		return fmt.Errorf("marshaljson , auth.NewQueryAccountParams error:%v", err)
	}
	res, err := client.ABCIQueryWithOptions(fmt.Sprintf("custom/%s/%s", auth.QuerierRoute, auth.QueryAccount), bs, rpchttp.ABCIQueryOptions{Prove: true})
	if err != nil {
		return fmt.Errorf("ABCIQuery , error:%v", err)
	}
	if !res.Response.IsOK() {
		return fmt.Errorf("not resp.IsOK")
	}
	//get exported account
	var expAcct exported.Account
	if err := appCdc.UnmarshalJSON(res.Response.Value, &expAcct); err != nil {
		return fmt.Errorf("Get exported account error:%v", err)
	}
	sequenceNumber := expAcct.GetSequence()

	gasPrice, err := sdk.ParseDecCoins("0.0000000000001stake")
	if err != nil {
		return fmt.Errorf("ParseDecCoins to get gasprice error:%v", err)
	}
	var gas uint64 = 200000
	fee, err := CalcCosmosFees(gasPrice, gas)
	if err != nil {
		return fmt.Errorf("CalcCosmosFees to get gasprice error:%v", err)
	}
	toSign := auth.StdSignMsg{
		Sequence:      sequenceNumber,
		AccountNumber: expAcct.GetAccountNumber(),
		//ChainID:       "testing",
		ChainID: ChainID,
		Msgs:    []sdk.Msg{msg},
		Fee:     auth.NewStdFee(200000, fee),
	}
	sig, err := signerPriv.Sign(toSign.Bytes())
	if err != nil {
		return fmt.Errorf("failed to sign raw tx: (error: %v, raw tx: %x)", err, toSign.Bytes())
	}
	tx := auth.NewStdTx([]sdk.Msg{msg}, toSign.Fee, []auth.StdSignature{{signerPriv.PubKey(), sig}}, toSign.Memo)

	txEncoder := auth.DefaultTxEncoder(appCdc)
	rawTx, err := txEncoder(tx)
	if err != nil {
		return fmt.Errorf("failed to encode signed tx: %v", err)
	}

	broadRes, err := client.BroadcastTxSync(rawTx)
	if err != nil {
		return fmt.Errorf("failed to broadcast tx: (error: %v, raw tx: %x)", err, rawTx)

	}
	fmt.Printf("ResultBroadCastTxSync is %v\n", *broadRes)
	if broadRes.Code == 0 {
		fmt.Printf("hash is %x\n", broadRes.Hash)
		time.Sleep(6 * time.Second)
		CheckTxSuccessful(client, hex.EncodeToString(broadRes.Hash))
	}

	return nil
}

func GetCosmosPrivateKey(path string, pwd []byte) (crypto.PrivKey, types.AccAddress, error) {
	bz, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, types.AccAddress{}, err
	}

	privKey, err := mintkey.UnarmorDecryptPrivKey(string(bz), string(pwd))
	if err != nil {
		return nil, types.AccAddress{}, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	return privKey, types.AccAddress(privKey.PubKey().Address().Bytes()), nil
}

func CalcCosmosFees(gasPrice types.DecCoins, gas uint64) (types.Coins, error) {
	if gasPrice.IsZero() {
		return types.Coins{}, errors.New("gas price is zero")
	}
	if gas == 0 {
		return types.Coins{}, errors.New("gas is zero")
	}
	glDec := types.NewDec(int64(gas))
	fees := make(types.Coins, len(gasPrice))
	for i, gp := range gasPrice {
		fee := gp.Amount.Mul(glDec)
		fees[i] = types.NewCoin(gp.Denom, fee.Ceil().RoundInt())
	}
	return fees, nil
}
