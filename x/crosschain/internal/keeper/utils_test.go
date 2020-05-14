package keeper

import (
	"github.com/cosmos/gaia/x/crosschain/internal/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func Test_Pad32Bytes(t *testing.T) {
	nums := []string{
		"1", "ff", "ffff",
	}
	for _, num := range nums {
		bigInt, ok := big.NewInt(0).SetString(num, 16)
		if !ok {
			t.Errorf("newint.setstring err")
		}
		res, err := types.Pad32Bytes(bigInt, 32)
		assert.Nil(t, err)
		//fmt.Printf("Pad32Bytes(0x%s) = %s, len = %d\n", num, hex.EncodeToString(res), len(res))

		numBigInt, err := types.Unpad32Bytes(res, 32)
		assert.Nil(t, err)
		//fmt.Printf("num       = %s\n", bigInt.String())
		//fmt.Printf("numBigInt = %s\n\n", numBigInt.String())
		assert.Equal(t, bigInt, numBigInt)
	}
}
