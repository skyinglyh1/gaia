package ft_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/gaia/test/test-module"
	"github.com/polynetwork/cosmos-poly-module/ft"
	"testing"

	"github.com/cosmos/gaia/app"
)

func Test_ft_GetDenomInfo(t *testing.T) {
	_, body, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "GET", fmt.Sprintf("/ft/denom_info/%s", "oep41"), nil)
	if err != nil {
		t.Errorf("GetDenomInfo, SendRequest Error:%v", err)
	}
	cdc := app.MakeCodec()
	var resp rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &resp)
	var denomInfo ft.DenomInfo
	cdc.MustUnmarshalJSON(resp.Result, &denomInfo)
	fmt.Printf("denomInfo is %s\n", denomInfo.String())
}

func Test_ft_GetDenomInfoWithChainId(t *testing.T) {
	_, body, err := test.SendRequest(test.NetConfig["gaia"].RestIp, "GET", fmt.Sprintf("/ft/denom_cc_info/%s/%s", "oep41", "2"), nil)
	if err != nil {
		t.Errorf("GetDenomInfo, SendRequest Error:%v", err)
	}

	cdc := app.MakeCodec()
	var resp rest.ResponseWithHeight
	cdc.MustUnmarshalJSON(body, &resp)

	var denomInfo ft.DenomCrossChainInfo
	cdc.MustUnmarshalJSON(resp.Result, &denomInfo)
	fmt.Printf("denomCrossChainInfoId is %s\n", denomInfo.String())
}
