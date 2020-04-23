package keeper

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// coming from "github.com/ethereum/go-ethereum/common/math"
func Pad32Bytes(bigint *big.Int) ([]byte, error) {
	ret := make([]byte, 32)
	if bigint.Cmp(big.NewInt(0)) < 1 {
		return nil, fmt.Errorf("Pad32Bytes only support positive big.Int, but got:%s", bigint.String())
	}
	bigBs := bigint.Bytes()
	if len(bigBs) > 32 || (len(bigBs) == 32 && bigBs[31]&0x80 == 1) {
		return nil, fmt.Errorf("Pad32Bytes only support maximum 2**255-1 big.Int, but got:%s", bigint.String())
	}
	copy(ret[:len(bigBs)], bigBs)
	copy(ret[len(bigBs):], make([]byte, 32-len(bigBs)))
	return ToArrayReverse(ret), nil
}

func Unpad32Bytes(paddedBs []byte) (*big.Int, error) {
	paddedBs = ToArrayReverse(paddedBs)
	if len(paddedBs) != 32 {
		return nil, fmt.Errorf("Unpad32Bytes only support 32 bytes value, but got:%s", hex.EncodeToString(paddedBs))
	}
	nonZeroPos := 31
	for i := nonZeroPos; i >= 0; i-- {
		p := paddedBs[i]
		if p != 0x0 {
			nonZeroPos = i
			break
		}
	}
	if nonZeroPos == 31 && paddedBs[31]&0x80 == 1 {
		return nil, fmt.Errorf("Unpad32Bytes only support 32 bytes value, but got:%s", hex.EncodeToString(paddedBs))
	}

	return big.NewInt(0).SetBytes(paddedBs[:nonZeroPos+1]), nil

}

func ToArrayReverse(arr []byte) []byte {
	l := len(arr)
	x := make([]byte, 0)
	for i := l - 1; i >= 0; i-- {
		x = append(x, arr[i])
	}
	return x
}
