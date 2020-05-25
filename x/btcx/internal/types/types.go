package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DenomInfo struct {
	Creator          sdk.AccAddress
	TotalSupply      sdk.Int
	RedeemScipt      []byte
	RedeemScriptHash []byte
}

func (msg DenomInfo) String() string {
	return fmt.Sprintf(`DenomInfo:
  Creator:        	 %s
  TotalSupply:		 %s
  RedeemScipt: 		 %x
  RedeemScriptHash:  %x
`, msg.Creator.String(), msg.TotalSupply.String(), msg.RedeemScipt, msg.RedeemScriptHash)
}

type DenomInfoWithId struct {
	DenomInfo
	ToChainId   uint64
	ToAssetHash []byte
}

func (msg DenomInfoWithId) String() string {
	return msg.DenomInfo.String() + fmt.Sprintf(`
  ToChainId:       	 %d
  ToAssetHash:		 %x
`, msg.ToChainId, msg.ToAssetHash)
}
