package test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/gaia/app"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
	"testing"
)

func TestTOCosmosRoutine(t *testing.T) {
	acctAddr, _ := sdk.AccAddressFromBech32("cosmos1yq77njl39m92jzdsxqswwk2uq7y70jcraxxy7e")

	appCdc := app.MakeCodec()
	client, err := rpchttp.NewHTTP(NetConfig["gaia"].Ip, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	status, _ := client.Status()
	fmt.Println("currentHeight:", status.SyncInfo.LatestBlockHeight)

	bp := bank.NewQueryBalanceParams(acctAddr)
	raw, err := appCdc.MarshalJSON(bp)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.ABCIQueryWithOptions("/custom/bank/balances", raw, rpchttp.ABCIQueryOptions{Prove: true, Height: status.SyncInfo.LatestBlockHeight})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Response.Value, res.Response.Proof)

	p := auth.NewQueryAccountParams(acctAddr)
	raw, err = appCdc.MarshalJSON(p)
	if err != nil {
		t.Fatal(err)
	}
	res, err = client.ABCIQueryWithOptions("/store/acc/key", raw, rpchttp.ABCIQueryOptions{Prove: true, Height: status.SyncInfo.LatestBlockHeight - 1})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Response.Height, "hex:", res.Response.Proof)

	res, err = client.ABCIQueryWithOptions("/store/acc/key", raw, rpchttp.ABCIQueryOptions{Prove: true, Height: status.SyncInfo.LatestBlockHeight - 1000})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Response.Height, "hex:", res.Response.Proof)

	res, err = client.ABCIQueryWithOptions("/store/acc/key", raw, rpchttp.ABCIQueryOptions{Prove: true, Height: status.SyncInfo.LatestBlockHeight})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(res.Response.Height, "hex:", res.Response.Proof)

	for i := status.SyncInfo.LatestBlockHeight - 10; i < status.SyncInfo.LatestBlockHeight; i++ {
		res, err = client.ABCIQueryWithOptions("/store/acc/key", raw, rpchttp.ABCIQueryOptions{Prove: true, Height: int64(i)})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(res.Response.Height, "hex:", res.Response.Proof)
	}

	//
	//paramHash, _ := hex.DecodeString("44a62f471e82bf7de0c49aca023db0a8559cffd10ccaf744f2f31ce3aa152286")
	//
	//res, err = client.ABCIQueryWithOptions("/store/ccm/key", ccm.GetCrossChainTxKey(paramHash), rpchttp.ABCIQueryOptions{Prove: true, Height: status.SyncInfo.LatestBlockHeight - 10})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//fmt.Println(res.Response.Height, "hex:", res.Response.Proof)
	//

}
