package keeper

import (
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
		res, err := Pad32Bytes(bigInt)
		assert.Nil(t, err)
		//fmt.Printf("Pad32Bytes(0x%s) = %s, len = %d\n", num, hex.EncodeToString(res), len(res))

		numBigInt, err := Unpad32Bytes(res)
		assert.Nil(t, err)
		//fmt.Printf("num       = %s\n", bigInt.String())
		//fmt.Printf("numBigInt = %s\n\n", numBigInt.String())
		assert.Equal(t, bigInt, numBigInt)
	}
}
