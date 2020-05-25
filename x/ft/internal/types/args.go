package types

import (
	"fmt"
	mcc "github.com/ontio/multi-chain/common"
	"math/big"
)

type TxArgs struct {
	ToAddress []byte
	Amount    *big.Int
}

func (this *TxArgs) Serialization(sink *mcc.ZeroCopySink, intLen int) error {
	sink.WriteVarBytes(this.ToAddress)
	paddedAmountBs, err := Pad32Bytes(this.Amount, intLen)
	if err != nil {
		return fmt.Errorf("TxArgs Serialization error:%v", err)
	}
	sink.WriteBytes(mcc.ToArrayReverse(paddedAmountBs))
	return nil
}

func (this *TxArgs) Deserialization(source *mcc.ZeroCopySource, intLen int) error {
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

	this.ToAddress = toAddress
	this.Amount = amount
	return nil
}
