package types

import (
	"fmt"
	polycommon "github.com/cosmos/gaia/x/headersync/poly-utils/common"
	"math/big"
)

type TxArgs struct {
	ToAssetHash []byte
	ToAddress   []byte
	Amount      *big.Int
}

func (this *TxArgs) Serialization(sink *polycommon.ZeroCopySink, intLen int) error {
	sink.WriteVarBytes(this.ToAssetHash)
	sink.WriteVarBytes(this.ToAddress)
	paddedAmountBs, err := Pad32Bytes(this.Amount, intLen)
	if err != nil {
		return fmt.Errorf("TxArgs Serialization error:%v", err)
	}
	sink.WriteBytes(polycommon.ToArrayReverse(paddedAmountBs))
	return nil
}

func (this *TxArgs) Deserialization(source *polycommon.ZeroCopySource, intLen int) error {
	txHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize txHash error")
	}
	toAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize ToAddress error")
	}
	paddedAmountBs, eof := source.NextBytes(uint64(intLen))
	if eof {
		return fmt.Errorf("TxArgs deserialize Amount error")
	}
	amount, err := Unpad32Bytes(paddedAmountBs, intLen)
	if err != nil {
		return fmt.Errorf("TxArgs Deserialization error:%v", err)
	}

	this.ToAssetHash = txHash
	this.ToAddress = toAddress
	this.Amount = amount
	return nil
}
