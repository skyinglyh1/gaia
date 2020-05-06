package types

import (
	"fmt"
	mcc "github.com/ontio/multi-chain/common"
	"math/big"
)

type TxArgs struct {
	ToAssetHash []byte
	ToAddress   []byte
	Amount      *big.Int
}

func (this *TxArgs) Serialization(sink *mcc.ZeroCopySink) error {
	sink.WriteVarBytes(this.ToAssetHash)
	sink.WriteVarBytes(this.ToAddress)
	paddedAmountBs, err := Pad32Bytes(this.Amount)
	if err != nil {
		return fmt.Errorf("TxArgs Serialization error:%v", err)
	}
	sink.WriteBytes(mcc.ToArrayReverse(paddedAmountBs))
	return nil
}

func (this *TxArgs) Deserialization(source *mcc.ZeroCopySource) error {
	txHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize txHash error")
	}
	toAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize ToAddress error")
	}
	paddedAmountBs, eof := source.NextBytes(32)
	if eof {
		return fmt.Errorf("TxArgs deserialize Amount error")
	}
	amount, err := Unpad32Bytes(paddedAmountBs)
	if err != nil {
		return fmt.Errorf("TxArgs Deserialization error:%v", err)
	}

	this.ToAssetHash = txHash
	this.ToAddress = toAddress
	this.Amount = amount
	return nil
}
