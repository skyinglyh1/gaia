package types_test

import (
	"encoding/hex"
	"fmt"
	polycommon "github.com/cosmos/gaia/x/headersync/poly-utils/common"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
	"testing"
)

func TestTxArgs_Serialization(t *testing.T) {
	argsBs, _ := hex.DecodeString("06657468786363148af0a30541c146b214f5ac59c7999eacae0592d600407a10f35a0000000000000000000000000000000000000000000000000000")
	txArgs := new(types.TxArgs)
	err := txArgs.Deserialization(polycommon.NewZeroCopySource(argsBs), 32)
	if err != nil {
		t.Errorf("Deserialization Error: %v", err)
	}

	fmt.Printf("txargs.ToAssetHash is %x\n", txArgs.ToAssetHash)
	fmt.Printf("txargs.ToAddress is %x\n", txArgs.ToAddress)
	fmt.Printf("txargs.Amount is %s\n", txArgs.Amount.String())
}
