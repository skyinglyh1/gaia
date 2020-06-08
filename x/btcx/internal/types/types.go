package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DenomInfo struct {
	Creator          sdk.AccAddress
	TotalSupply      sdk.Int
	RedeemScipt      string
	RedeemScriptHash string
}

func (msg DenomInfo) String() string {
	return fmt.Sprintf(`
  Creator:        	 			%s
  TotalSupply:		 			%s
  RedeemScriptHash(AssetHash):  %s
  RedeemScipt: 					%s
`, msg.Creator.String(), msg.TotalSupply.String(), msg.RedeemScriptHash, msg.RedeemScipt)
}

type DenomCrossChainInfo struct {
	DenomInfo
	ToChainId   uint64
	ToAssetHash string
}

func (msg DenomCrossChainInfo) String() string {
	return msg.DenomInfo.String() + fmt.Sprintf(`
  ToChainId:       	 			%d
  ToAssetHash:		 			%s
`, msg.ToChainId, msg.ToAssetHash)
}
