package lockproxy_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/test/test-module"
	"testing"
)

func Test_lockproxy_GetLockProxyHashByOperator(t *testing.T) {
	operator := "cosmos1whx6nu3eaztwzqplkt6z5kuz0dhm8rsn34rx8w"
	_, body, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "GET", fmt.Sprintf("/lockproxy/proxy_hash_by_operator/%s", operator), nil)
	//operator := "cosmos1vnnpptmw2vlm5h06ej3t23vsx6jaqgtcwexesm"
	//_, body, err := SendRequest("172.168.3.94:9091", "GET", fmt.Sprintf("/lockproxy/proxy_hash_by_operator/%s", operator), nil)
	if err != nil {
		t.Errorf("Test_lockproxy_GetProxyHash, SendRequest Error:%v", err)
	}
	fmt.Printf("body is %s\n", string(body))
	cdc := app.MakeCodec()
	var resp rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &resp)
	var denomInfo []byte
	cdc.MustUnmarshalJSON(resp.Result, &denomInfo)
	fmt.Printf("proxyHash is %x\n", denomInfo)
}
func Test_lockproxy_GetProxyHash(t *testing.T) {
	//lockProxyHash := "db8afcccebc026c6cae1d541b25f80a83b065c8a"
	//toChainId := "2"
	lockProxyHash := "64e610af6e533fba5dfacca2b5459036a5d02178"
	toChainId := "2"
	_, body, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "GET", fmt.Sprintf("/lockproxy/proxy_hash/%s/%s", lockProxyHash, toChainId), nil)
	if err != nil {
		t.Errorf("Test_lockproxy_GetProxyHash, SendRequest Error:%v", err)
	}
	cdc := app.MakeCodec()
	var resp rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &resp)
	var denomInfo []byte
	cdc.MustUnmarshalJSON(resp.Result, &denomInfo)
	fmt.Printf("proxyHash is %x\n", denomInfo)
}

func Test_lockproxy_GetAssetHash(t *testing.T) {
	//lockProxyHash := "db8afcccebc026c6cae1d541b25f80a83b065c8a"
	//assetDenom := "reth3"
	//toChainId := "2"
	lockProxyHash := " 75cda9f239e896e1003fb2f42a5b827b6fb38e13"
	assetDenom := "ontcc"
	toChainId := "3"
	_, body, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "GET", fmt.Sprintf("/lockproxy/asset_hash/%s/%s/%s", lockProxyHash, assetDenom, toChainId), nil)
	if err != nil {
		t.Errorf("Test_lockproxy_GetAssetHash, SendRequest Error:%v", err)
	}
	cdc := app.MakeCodec()
	var resp rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &resp)
	var denomInfo []byte
	cdc.MustUnmarshalJSON(resp.Result, &denomInfo)
	fmt.Printf("assetHash is %x\n", denomInfo)
}

func Test_lockproxy_GetLockedAmt(t *testing.T) {
	_, body, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "GET", fmt.Sprintf("/lockproxy/locked_amount/%s", "ether3"), nil)
	if err != nil {
		t.Errorf("Test_lockproxy_GetAssetHash, SendRequest Error:%v", err)
	}
	cdc := app.MakeCodec()
	var resp rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &resp)
	var amt sdk.Int
	cdc.MustUnmarshalJSON(resp.Result, &amt)
	fmt.Printf("locked amount of denom: %s is %+v\n", "reth3", amt.String())
}
