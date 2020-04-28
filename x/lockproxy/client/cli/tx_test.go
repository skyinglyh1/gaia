package cli_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/davecgh/go-spew/spew"
)


func Test_ParseCoins(t *testing.T) {

	args0 := "1000000000ont,1000000000000000000ong"
	//args1 := "1000000000stake,1000000000validatortoken"
	coins, err := sdk.ParseCoins(args0)
	if err != nil {
		t.Errorf("parsecoins error:%v", err)
	}
	spew.Printf("coins are %v\n", coins)
}

func Test_UnmarshalOperator(t *testing.T) {

	operator, err := sdk.AccAddressFromHex("b1ea3cf8713100ae1e84aecb5f80ecf741416be8")
	if err != nil {
		t.Errorf("could not unmarshal result to sdk.AccAddress:%v", err)
	}
	spew.Printf("opeartor are %s\n", operator.String())


}
func Test_GetHexAddressFromBench32(t *testing.T) {
	user2 := "cosmos1cwphz9u9qss84vk4g5sktfcxttwvm6qk3upd9z"
	user2Addr, err := sdk.AccAddressFromBech32(user2)
	if err != nil {
		t.Errorf("err = %v", err)
	}
	fmt.Printf("user2 Hex Address = %s\n", hex.EncodeToString(user2Addr.Bytes()))

	toAddress := make(sdk.AccAddress, len(user2Addr.Bytes()))
	copy(toAddress, user2Addr.Bytes())
	fmt.Printf("toAdddress = %s\n", hex.EncodeToString(toAddress))

}