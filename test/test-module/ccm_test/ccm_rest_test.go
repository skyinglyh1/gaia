package ccm_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/gaia/app"
	"github.com/cosmos/gaia/test/test-module"
	"github.com/polynetwork/cosmos-poly-module/ccm"
	"testing"
)

func Test_lockproxy_GetCCMParameters(t *testing.T) {
	_, body, err := test.SendRequest(test.NetConfig["switcheo"].RestIp, "GET", fmt.Sprintf("/ccm/parameters"), nil)
	if err != nil {
		t.Errorf("Test_lockproxy_GetAssetHash, SendRequest Error:%v", err)
	}
	cdc := app.MakeCodec()
	var resp rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &resp)
	var params ccm.Params
	cdc.MustUnmarshalJSON(resp.Result, &params)
	fmt.Printf("ccm.Params is %+v\n", params)
}
