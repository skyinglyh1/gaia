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
	return fmt.Sprintf(`
  Creator:        	 			%s
  TotalSupply:		 			%s
  RedeemScriptHash(AssetHash):  %x
  RedeemScipt: 					%x
`, msg.Creator.String(), msg.TotalSupply.String(), msg.RedeemScriptHash, msg.RedeemScipt)
}

type DenomCrossChainInfo struct {
	DenomInfo
	ToChainId   uint64
	ToAssetHash []byte
}

func (msg DenomCrossChainInfo) String() string {
	return msg.DenomInfo.String() + fmt.Sprintf(`
  ToChainId:       	 			%d
  ToAssetHash:		 			%x
`, msg.ToChainId, msg.ToAssetHash)
}
