package test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/x/crosschain"
	"github.com/tendermint/tendermint/crypto"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"io/ioutil"
	"testing"
)

const (
	//ip = "tcp://172.168.3.93:26657"
	ip              = "tcp://172.168.3.93:26657"
	validatorWallet = "./wallets/validator"
	operatorWallet  = "./wallets/operator"
	operatorPwd     = "12345678"
	ChainID         = "testing"
)

var (
	header0         = "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c033644e70a2b4f8de4a15c4a0cd79315673b8346d033804807058f3ff4252900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c8365b000000001dac2b7c00000000fd1a057b226c6561646572223a343239343936373239352c227672665f76616c7565223a22484a675171706769355248566745716354626e6443456c384d516837446172364e4e646f6f79553051666f67555634764d50675851524171384d6f38373853426a2b38577262676c2b36714d7258686b667a72375751343d222c227672665f70726f6f66223a22785864422b5451454c4c6a59734965305378596474572f442f39542f746e5854624e436667354e62364650596370382f55706a524c572f536a5558643552576b75646632646f4c5267727052474b76305566385a69413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a343239343936373239352c226e65775f636861696e5f636f6e666967223a7b2276657273696f6e223a312c2276696577223a312c226e223a372c2263223a322c22626c6f636b5f6d73675f64656c6179223a31303030303030303030302c22686173685f6d73675f64656c6179223a31303030303030303030302c22706565725f68616e647368616b655f74696d656f7574223a31303030303030303030302c227065657273223a5b7b22696e646578223a312c226964223a2231323035303364313031383338383037656334303739613436666539386436626439613036393061626362643863653136653066626334353230633763376566373838356462227d2c7b22696e646578223a322c226964223a2231323035303361366231623065316336393737663434663336323332613566336236316236613835396234636535313437633439616363666139613432663438336631323034227d2c7b22696e646578223a332c226964223a2231323035303266663764666337303562623561633638643265383932333063363632393939616562313832383431333165396663653934656639666166356239393137353364227d2c7b22696e646578223a342c226964223a2231323035303334343031376363636138323064393066306562623436316466343633333762303932336230616532626365353833636565316132363234633932303865323038227d2c7b22696e646578223a352c226964223a2231323035303331326631303233393531333134336330323938346263346561396438353438383366636466343937333264633732376466613734373438326663383037653634227d2c7b22696e646578223a362c226964223a2231323035303333336334343833376464623934616435666130656234363062306634393135346639303530333631396434643263386565303833333066623831353834316432227d2c7b22696e646578223a372c226964223a2231323035303363366536383165353135346566626136346337356230616131636135343438396261653736353330373764313664646439373236663336356265333036323264227d5d2c22706f735f7461626c65223a5b362c352c342c332c372c322c372c372c352c352c322c322c322c322c362c352c322c342c312c332c342c312c342c332c332c322c342c352c372c312c342c332c342c352c332c352c352c342c322c312c342c332c312c352c352c352c322c362c342c332c312c362c322c322c312c332c332c322c332c372c372c362c342c342c362c372c372c362c322c362c372c372c312c332c342c312c352c362c322c372c342c342c362c352c312c332c352c372c352c332c312c362c312c322c362c362c312c372c362c362c372c332c372c312c315d2c226d61785f626c6f636b5f6368616e67655f76696577223a31303030307d7dfd9e5473b163f591a8829d83288809d97c20ab2a0000"
	header1         = "0000000000000000000000002a8944a1f753de95e505607d82f70c8a2291c4cf3ef6277b47b282628882c6e400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a39f2294c54fe8e3575b3b3a531f6a91605ef153f8ed5344094c4d67ba489a3f1645a55e0100000036810287538edf27fd0c017b226c6561646572223a332c227672665f76616c7565223a224244714d70463966716a7754596d4e4c6c50795162464e505066546f374261786e677543775563697844713733636856386c6f6e525774474558596d77387671355279372f41505778434d573737433773594c33542b733d222c227672665f70726f6f66223a2246456c6453493638683352532f366f473951676e334835705136443244434c65367966614a485a7169745a6d696b50347a647843354a7943664554662f436e304f474c6a3966736b476a5566523668516370797568673d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000523120502ff7dfc705bb5ac68d2e89230c662999aeb18284131e9fce94ef9faf5b991753d23120503d101838807ec4079a46fe98d6bd9a0690abcbd8ce16e0fbc4520c7c7ef7885db2312050333c44837ddb94ad5fa0eb460b0f49154f90503619d4d2c8ee08330fb815841d22312050344017ccca820d90f0ebb461df46337b0923b0ae2bce583cee1a2624c9208e20823120503c6e681e5154efba64c75b0aa1ca54489bae7653077d16ddd9726f365be30622d0542051c5c15d9e88471206d349fddceb5a48b8b0c261ea9a6b306cf994e56bb4c7ac1846f5f97356268f930ffd37f542f449baaaf897a5fb8be8f6480656c0e6e596b7342051b288d849d6db42b4122de8f781c905b85ceabdcdc7e16a8a2418758d06ceea0895fd9eb9fdd713af03fa9c9593bf930e25dba59331bc27d63cf021ead892de06c42051ba1c6e68c6809c529cdcbc003ab411f000308828f1ffe7bc1d997b3a99eec7a311c27c81291680e4b29cd4bc7ddb2f5c20f615e18512652180d863fb255519a5842051bca95efb2dbff60a538e1726c249df04cd2cd9127f2c55becfefcf51961a5d9072a0d36bbe349befabd34d08ac03535491bfa1513c686e24a316fd7795f54ada042051c9eec1630be76366a2c374ca5ff349d61280a213739c81a872b46b3a8b0df3d2a74c3d03aa2218b3eba4e8d7a617ce49d27fed0d36acb4ce59f4944e0bb5ee3d1"
	header2         = "0000000000000000000000006eb34794b01ee25e05b5e41f35685bbf0f285ef0bb73905fddd72f6bc3c1ca1d000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008c05a61780e1873482c24cc98a4fbb53f01c45ee96e2e3afec992097e52def133445a55e0200000097cc879256d3b949fd0c017b226c6561646572223a332c227672665f76616c7565223a22424d63306f455a56513546374c707a436c3962424a4779324b507a2b32776e37636f6c717741745732754846746f58497a3541705142574c35357a30444e5a386f58636f4c4a6d337a526b596b686f2b7a30486648396f3d222c227672665f70726f6f66223a222f7252754d4d592b31534867504c475470325a31703469645743645a7748684646496d3533613534796a79466e5651504b6635397158756d6f383430417859565375696175796c35564e477a6e753456333248356e673d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d00000000000000000000000000000000000000000623120502ff7dfc705bb5ac68d2e89230c662999aeb18284131e9fce94ef9faf5b991753d23120503c6e681e5154efba64c75b0aa1ca54489bae7653077d16ddd9726f365be30622d2312050312f10239513143c02984bc4ea9d854883fcdf49732dc727dfa747482fc807e642312050344017ccca820d90f0ebb461df46337b0923b0ae2bce583cee1a2624c9208e2082312050333c44837ddb94ad5fa0eb460b0f49154f90503619d4d2c8ee08330fb815841d223120503d101838807ec4079a46fe98d6bd9a0690abcbd8ce16e0fbc4520c7c7ef7885db0642051b15e124c12818c34003566fb5460e9563ee50b20a6ab90dd5511661035d9abbd84cd98115da932e2470854247498105a783546446039c6f6771cecff9c0a59a2d42051bc5fe5c235adcedc9975649804fcbb9fb1b1d68839cfa911389502280500be2e90c9b21f46366774a626f8aee5313c40379fcd71062ae794f155726dbf5cce9ba42051cf31e932fe19b7aee60f04a6786d8b85cd2bcc8571ad375479b7bf4e5dc6a9cc03abad6b70564edf58f9149dd558ac87c3b827bd9d93280784614ab922dfdeac842051c905500ca7a76dba83aa157bc8a725c952486c0bbcfe1c5502ce6e5f96f6ab56a64e0b0b088646f1b4eebb99f6fe708651dd39bc6bbba3145bad82ccd4e40be9842051bce7a5b564889ae0a093160a04a4d01e4a9251c5d84839a2f679b9d387ae94ed432a96b7aaccf7d4a9252d5aa32f76554b24b15ff99df892a2dbed47a6c3edf1942051c585e4851df6fa9626048d46db2a5abfb2b74817fc3163e9fed856ea25ca623ee5dd294e290d5b14599dc7452069f987ce54ddc1bb1a17db687aced9197a14e2f"
	RedeemKey, _    = hex.DecodeString("c330431496364497d7257839737b5e4596f5ac06")
	RedeemScript, _ = hex.DecodeString("552102dec9a415b6384ec0a9331d0cdf02020f0f1e5731c327b86e2b5a92455a289748210365b1066bcfa21987c3e207b92e309b95ca6bee5f1133cf04d6ed4ed265eafdbc21031104e387cd1a103c27fdc8a52d5c68dec25ddfb2f574fbdca405edfd8c5187de21031fdb4b44a9f20883aff505009ebc18702774c105cb04b1eecebcb294d404b1cb210387cda955196cc2b2fc0adbbbac1776f8de77b563c6d2a06a77d96457dc3d0d1f2102dd7767b6a7cc83693343ba721e0f5f4c7b4b8d85eeb7aec20d227625ec0f59d321034ad129efdab75061e8d4def08f5911495af2dae6d3e9a4b6e7aeb5186fa432fc57ae")
)

func Test_SendTxHeaderSync(t *testing.T) {
	//client := rpchttp.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey("operatorWallet", []byte(operatorPwd))
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
	coins, err := sdk.ParseCoins("1000000000ont,1000000000000000000ong,1000000000000btc")
	if err != nil {
		t.Errorf("parse coins err:%v", err)
	}
	msg := crosschain.NewMsgCreateCoins(creator, coins)
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

	ontProxy, _ := hex.DecodeString("")
	ethProxy, _ := hex.DecodeString("")
	proxyHashInOtherChain := []struct {
		ChainId   uint64
		ProxyHash []byte
	}{
		{3, ontProxy},
		{2, ethProxy},
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

	coins, _ := sdk.ParseCoins("1000000000ont,1000000000000000000ong")
	ontHash, _ := hex.DecodeString("0000000000000000000000000000000000000001")
	ongHash, _ := hex.DecodeString("0000000000000000000000000000000000000002")
	assetHashInOnt := [][]byte{
		ontHash, ongHash,
	}
	for i, coin := range coins {
		msg := crosschain.NewMsgBindAssetParam(fromAddr, coin.Denom, 3, assetHashInOnt[i], coins.AmountOf(coin.Denom), true)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
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

	msg := crosschain.NewMsgSetRedeemScript(fromAddr, "btc", RedeemKey, RedeemScript)
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
	}
}

func Test_BindNoVMChainAssetHash(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	redeemKey := RedeemKey
	limit := sdk.NewInt(1000000000000)
	//btcInOntHash, _ := hex.DecodeString("")
	//btcInEthHash, _ := hex.DecodeString("")
	btcInBtcHash := redeemKey
	btcAssetHashInNonBtcChain := []struct {
		ChainId           uint64
		AssetContractHash []byte
	}{
		//{3, btcInOntHash},
		//{2, btcInEthHash},
		{1, btcInBtcHash},
	}
	for _, btcAssetHash := range btcAssetHashInNonBtcChain {
		msg := crosschain.NewMsgBindNoVMChainAssetHash(fromAddr, "btc", btcAssetHash.ChainId, btcAssetHash.AssetContractHash, limit)
		if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
			t.Errorf("sendMsg error:%v", err)
		}
	}
}

func Test_Lock(t *testing.T) {
	client := rpchttp.NewHTTP(ip, "/websocket")
	appCdc := app.MakeCodec()

	fromPriv, fromAddr, err := GetCosmosPrivateKey(operatorWallet, []byte(operatorPwd))
	if err != nil {
		t.Errorf("err = %v ", err)
	}
	fmt.Printf("acct = %v\n", fromAddr.String())
	fmt.Printf("priv = %v\n", hex.EncodeToString(fromPriv.Bytes()))

	sourceAssetDenom := ""
	var toChainId uint64 = 1
	toAddress, _ := hex.DecodeString("")
	value := sdk.NewInt(100)
	msg := crosschain.NewMsgLock(fromAddr, sourceAssetDenom, toChainId, toAddress, &value)
	if err := sendMsg(client, fromAddr, fromPriv, appCdc, msg); err != nil {
		t.Errorf("sendMsg error:%v", err)
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
	fmt.Printf("ResultBroadCastTxSync is %v", broadRes)
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
