package types

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// coming from "github.com/ethereum/go-ethereum/common/math"
func Pad32Bytes(bigint *big.Int, intLen int) ([]byte, error) {
	ret := make([]byte, intLen)
	if bigint.Cmp(big.NewInt(0)) < 1 {
		return nil, fmt.Errorf("Pad32Bytes only support positive big.Int, but got:%s", bigint.String())
	}
	bigBs := bigint.Bytes()
	if len(bigBs) > intLen || (len(bigBs) == intLen && bigBs[intLen-1]&0x80 == 1) {
		return nil, fmt.Errorf("Pad32Bytes only support maximum 2**255-1 big.Int, but got:%s", bigint.String())
	}
	copy(ret[:len(bigBs)], bigBs)
	copy(ret[len(bigBs):], make([]byte, intLen-len(bigBs)))
	return ToArrayReverse(ret), nil
}

func Unpad32Bytes(paddedBs []byte, intLen int) (*big.Int, error) {
	//paddedBs = ToArrayReverse(paddedBs)
	if len(paddedBs) != intLen {
		return nil, fmt.Errorf("Unpad32Bytes only support 32 bytes value, but got:%s", hex.EncodeToString(paddedBs))
	}
	nonZeroPos := intLen - 1
	for i := nonZeroPos; i >= 0; i-- {
		p := paddedBs[i]
		if p != 0x0 {
			nonZeroPos = i
			break
		}
	}
	if nonZeroPos == intLen-1 && paddedBs[intLen-1]&0x80 == 1 {
		return nil, fmt.Errorf("Unpad32Bytes only support 32 bytes value, but got:%s", hex.EncodeToString(paddedBs))
	}

	return big.NewInt(0).SetBytes(ToArrayReverse(paddedBs[:nonZeroPos+1])), nil
}

func ToArrayReverse(arr []byte) []byte {
	l := len(arr)
	x := make([]byte, 0)
	for i := l - 1; i >= 0; i-- {
		x = append(x, arr[i])
	}
	return x
}
