package test

import "encoding/hex"

const (
	//ip = "tcp://172.168.3.93:26657"
	//validatorWallet = "./wallets/validator"
	//operatorWallet  = "./wallets/operator"

	ip              = "tcp://172.168.3.95:26657"
	validatorWallet = "./wallets/172.168.3.94_node0"
	operatorWallet  = "./wallets/operator"

	user0Wallet  = "./wallets/user0"
	operatorAddr = "cosmos1c0n2e6kuzp03pqm3av9q2v0fqn6ql3z5c5ddw7"
	user0Addr    = "cosmos1ayc6faczpj42eu7wjsjkwcj7h0q2p2e4vrlkzf"
	user1Addr    = "cosmos1mtgmggm73d4mqv5kcc7hvtplryflwhl998dk5q"
	user2Addr    = "cosmos1vmg4h3etfpy9a8fyru44uz87sw9dwmvfdpw358"
	operatorPwd  = "12345678"
	//ChainID         = "testing"
	ChainID = "cc-cosmos"
)

var (
	header0         = "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010ae3a2d1cba9ed56653edab871d93f8a96294debb6169a62681552dfd6d0fc70000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c8365b000000001dac2b7c00000000fd1a057b226c6561646572223a343239343936373239352c227672665f76616c7565223a22484a675171706769355248566745716354626e6443456c384d516837446172364e4e646f6f79553051666f67555634764d50675851524171384d6f38373853426a2b38577262676c2b36714d7258686b667a72375751343d222c227672665f70726f6f66223a22785864422b5451454c4c6a59734965305378596474572f442f39542f746e5854624e436667354e62364650596370382f55706a524c572f536a5558643552576b75646632646f4c5267727052474b76305566385a69413d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a343239343936373239352c226e65775f636861696e5f636f6e666967223a7b2276657273696f6e223a312c2276696577223a312c226e223a372c2263223a322c22626c6f636b5f6d73675f64656c6179223a31303030303030303030302c22686173685f6d73675f64656c6179223a31303030303030303030302c22706565725f68616e647368616b655f74696d656f7574223a31303030303030303030302c227065657273223a5b7b22696e646578223a312c226964223a2231323035303238313732393138353430623262353132656165313837326132613265336132386439383963363064393564616238383239616461376437646437303664363538227d2c7b22696e646578223a322c226964223a2231323035303338623861663632313065636664636263616232323535326566386438636634316336663836663963663961623533643836353734316366646238333366303662227d2c7b22696e646578223a332c226964223a2231323035303234383261636236353634623139623930363533663665396338303632393265386161383366373865376139333832613234613665666534316330633036663339227d2c7b22696e646578223a342c226964223a2231323035303236373939333061343261616633633639373938636138613366313265313334633031393430353831386437383364313137343865303339646538353135393838227d2c7b22696e646578223a352c226964223a2231323035303234363864643138393965643264316363326238323938383261313635613065636236613734356166306337326562323938326436366234333131623465663733227d2c7b22696e646578223a362c226964223a2231323035303265623162616162363032633538393932383235363163646161613761616262636464306363666362633365373937393361633234616366393037373866333561227d2c7b22696e646578223a372c226964223a2231323035303331653037373966356335636362323631323335326665346132303066393964336537373538653730626135336636303763353966663232613330663637386666227d5d2c22706f735f7461626c65223a5b362c342c332c352c362c312c322c352c342c372c342c322c332c332c372c362c352c342c362c352c312c342c332c312c322c352c322c322c362c312c342c352c342c372c322c332c342c312c352c372c342c312c322c322c352c362c342c342c322c372c332c362c362c352c312c372c332c312c362c312c332c332c322c342c342c312c352c362c352c312c322c362c372c352c362c332c342c372c372c332c322c372c312c352c362c352c322c332c362c322c362c312c372c372c372c312c372c342c332c332c332c322c312c372c355d2c226d61785f626c6f636b5f6368616e67655f76696577223a36303030307d7d9fe171f3fe643eb1c188400b828ba184816fc9ac0000"
	header1         = "000000000000000000000000f7259d9da6edb2672055c4f0efd8729f921ff4f2ea6cfe2c632bf9137a8eabbc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000d43e5bb4e452c5130a39a3fa2f4e738e84b8caba1ab8a525eb0c379224a0c48d6c6dba5e010000005a6108a580c36ac2fd0c017b226c6561646572223a322c227672665f76616c7565223a22424c48634b703946724866376b64383866685a3644724748314f735178726f795a6a66766165664d5546337673517a36764a654e2b3252657a524a515a396e686143554759645544745869533232355851584b773563413d222c227672665f70726f6f66223a223037366b5331617a4551714a6e61706774546e554e4b5131576649435755596a2f65554e693469714b46615a4c3345614b715338385855737241396267594152717a4763764c6635792f435a612f745653336e504a773d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d000000000000000000000000000000000000000005231205038b8af6210ecfdcbcab22552ef8d8cf41c6f86f9cf9ab53d865741cfdb833f06b231205028172918540b2b512eae1872a2a2e3a28d989c60d95dab8829ada7d7dd706d658231205031e0779f5c5ccb2612352fe4a200f99d3e7758e70ba53f607c59ff22a30f678ff23120502679930a42aaf3c69798ca8a3f12e134c019405818d783d11748e039de851598823120502eb1baab602c5899282561cdaaa7aabbcdd0ccfcbc3e79793ac24acf90778f35a0542011b86005fa58d8286db7873bf9f1f116b59757518b36568bd2fe3e4c52d80710bc026a25f8dd3b45aa609e1c0f9b01cf43f2b94d061c936862dcedfb5d3c125830f42011caf13504a2c253135307f440cfa7053d0c96268c20c882b19c85753a1e4cc72fb1344f3ef00535304d3ad908959d393c906548bb078c52f14c6fd60036193072242011b2574b5ee43fb9345e90c1e3c8269a49b4f8b45266ccd6e783ffb858a9766c96362df590aa2e89bc8c086ddb2a4c80dc43b9eae52cbb539f8ddccfa61e018293142011bd8cf6c36d04358ed8bc4055ae372a5302dc18a7b4e56959a1be01b3a20b831c94e04e5623518512cbb38d2b80d6e4c2bb3e246f0f2cd94251f0f2ba54475eb4142011b46e5c26aab0b23e0594f0769909b36c4c2f6a9ea6393a17ff680ea7a901e00f31bb8271d2c1d019486fe7d142f3ddc943d9bee3d71c890da5e66d0f20eb53b9c"
	header2         = "000000000000000000000000200964d7695fe6a46fcebae721c86e49a0ad27dfb957c692b03c3a259f557b5400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000627cc2d7499f71db41e137509bc52ecb22be84ac6fe9e2ce9d4fc69aca05b5708a6dba5e02000000f4a2f2179ec18853fd0c017b226c6561646572223a372c227672665f76616c7565223a22424c583176784e2f686353354c704f435a6370706e733471543348635030556948564e69586b70716f4a50546f3677436c597949664e4e6b54754566542b614459587a745a736b42356432494741334d4a5764635972593d222c227672665f70726f6f66223a2267754761493444715a52717762454362392f4a5976304f78622f6e34756873494a4b674d3347446b564146715968614b4a416e63546b6c6936536f554d46724c666d485150584950774e4a434c762f70672b525548513d3d222c226c6173745f636f6e6669675f626c6f636b5f6e756d223a302c226e65775f636861696e5f636f6e666967223a6e756c6c7d000000000000000000000000000000000000000005231205031e0779f5c5ccb2612352fe4a200f99d3e7758e70ba53f607c59ff22a30f678ff23120502679930a42aaf3c69798ca8a3f12e134c019405818d783d11748e039de8515988231205028172918540b2b512eae1872a2a2e3a28d989c60d95dab8829ada7d7dd706d65823120502482acb6564b19b90653f6e9c806292e8aa83f78e7a9382a24a6efe41c0c06f39231205038b8af6210ecfdcbcab22552ef8d8cf41c6f86f9cf9ab53d865741cfdb833f06b0542011c0b95961b716dafac5d0c1b5786433c91baf3b1fa48deb0d85825c3978521abb750426488ae3326da3a3e01f601c10b97e07123018bc8897284b5f9b2ddff343242011c373c58c6191f4112fffd19c78dea62dc99179346a41051c88fcab91c3e2a863c5fd1584f8164926a869af8826ebfd4294e015fe6199d5bfd86d594943875b50442011bd1489f4dc150dd22827f11b07b73c6e6ca71273eb44313969d0f73fa8d62a2701c1ea4b48f8b045704ca7e53f150908eecf4a16c4bcc1765ee419a715ed0f2a442011b50b28acba80bf88bab43c42566c66a4d05d9b39810017d05fd4b0dd2dce6df5a31cf1a51c8951c82eeb598ea598f72ad08bd83e7c1fbc28d3cc8a1e59aebb77842011c23a8532272456b79fdd446f09569a2568a2d7058c8f70a01c6a3d5596c2d89bc625e7bda1be1c3161e6c1fc6fb9cd483bff3cf4dc26e4f7bd4f247a3d49e529a"
	RedeemKey, _    = hex.DecodeString("c330431496364497d7257839737b5e4596f5ac06")
	RedeemScriptStr = "552102dec9a415b6384ec0a9331d0cdf02020f0f1e5731c327b86e2b5a92455a289748210365b1066bcfa21987c3e207b92e309b95ca6bee5f1133cf04d6ed4ed265eafdbc21031104e387cd1a103c27fdc8a52d5c68dec25ddfb2f574fbdca405edfd8c5187de21031fdb4b44a9f20883aff505009ebc18702774c105cb04b1eecebcb294d404b1cb210387cda955196cc2b2fc0adbbbac1776f8de77b563c6d2a06a77d96457dc3d0d1f2102dd7767b6a7cc83693343ba721e0f5f4c7b4b8d85eeb7aec20d227625ec0f59d321034ad129efdab75061e8d4def08f5911495af2dae6d3e9a4b6e7aeb5186fa432fc57ae"
	RedeemScript, _ = hex.DecodeString(RedeemScriptStr)
)

var (
	btcHashInBtcDev   = RedeemKey
	btcHahInEthDev, _ = hex.DecodeString("740C1a496A750a3C3F9A6Ca7e822C6BC776962eA")
	btcHahInOntDev, _ = hex.DecodeString("b7f398711664de1dd685d9ba3eee3b6b830a7d83")
)

var (
	proxyInOntHashDev, _ = hex.DecodeString("50478b75da76f14bb8358318b62897b97de043dd")
	proxyInEthHashDev, _ = hex.DecodeString("71CF3de5e27EcF7379a8EE74eF32C021dD068d8d")
)

var (
	oep4denInOntDev, _ = hex.DecodeString("99981b7485df558eb63f45ee19dcb0458b83ed25")
	oep4denInEthDev, _ = hex.DecodeString("AF6FB1B9a813295Bf9Ccd3139Ec013D5718069De")
)

var (
	oep4IndInOntDev, _ = hex.DecodeString("")
	oep4IndInEthDev, _ = hex.DecodeString("")
)
