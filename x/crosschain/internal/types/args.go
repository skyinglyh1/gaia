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

func (this *TxArgs) Serialization(sink *mcc.ZeroCopySink, intLen int) error {
	sink.WriteVarBytes(this.ToAssetHash)
	sink.WriteVarBytes(this.ToAddress)
	paddedAmountBs, err := Pad32Bytes(this.Amount, intLen)
	if err != nil {
		return fmt.Errorf("TxArgs Serialization error:%v", err)
	}
	sink.WriteBytes(mcc.ToArrayReverse(paddedAmountBs))
	return nil
}

func (this *TxArgs) Deserialization(source *mcc.ZeroCopySource, intLen int) error {
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

type ToBTCArgs struct {
	ToBtcAddress []byte
	Amount       uint64
	RedeemScript []byte
}

func (this *ToBTCArgs) Serialization(sink *mcc.ZeroCopySink, intLen int) error {
	sink.WriteVarBytes(this.ToBtcAddress)
	paddedAmountBs, err := Pad32Bytes(new(big.Int).SetUint64(this.Amount), intLen)
	if err != nil {
		return fmt.Errorf("ToBTCArgs Serialization error:%v", err)
	}
	sink.WriteBytes(mcc.ToArrayReverse(paddedAmountBs))
	sink.WriteVarBytes(this.RedeemScript)
	return nil
}

func (this *ToBTCArgs) Deserialization(source *mcc.ZeroCopySource, intLen int) error {
	toBtcAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize toBtcAddress error")
	}

	paddedAmountBs, eof := source.NextBytes(uint64(intLen))
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize Amount error")
	}
	amount, err := Unpad32Bytes(paddedAmountBs, intLen)
	if err != nil {
		return fmt.Errorf("ToBTCArgs Deserialization error:%v", err)
	}
	redeemScript, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize redeemScript error")
	}

	this.ToBtcAddress = toBtcAddress
	this.Amount = amount.Uint64()
	this.RedeemScript = redeemScript
	return nil
}

type BTCArgs struct {
	ToBtcAddress []byte
	Amount       uint64
}

func (this *BTCArgs) Serialization(sink *mcc.ZeroCopySink, intLen int) error {
	sink.WriteVarBytes(this.ToBtcAddress)
	paddedAmountBs, err := Pad32Bytes(new(big.Int).SetUint64(this.Amount), intLen)
	if err != nil {
		return fmt.Errorf("ToBTCArgs Serialization error:%v", err)
	}
	sink.WriteBytes(mcc.ToArrayReverse(paddedAmountBs))
	return nil
}

func (this *BTCArgs) Deserialization(source *mcc.ZeroCopySource, intLen int) error {
	toBtcAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize toBtcAddress error")
	}

	paddedAmountBs, eof := source.NextBytes(uint64(intLen))
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize Amount error")
	}
	amount, err := Unpad32Bytes(paddedAmountBs, intLen)
	if err != nil {
		return fmt.Errorf("ToBTCArgs Deserialization error:%v", err)
	}
	this.ToBtcAddress = toBtcAddress
	this.Amount = amount.Uint64()
	return nil
}
