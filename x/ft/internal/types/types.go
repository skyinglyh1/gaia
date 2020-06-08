package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DenomInfo struct {
	Creator     sdk.AccAddress
	Denom       string
	AssetHash   string
	TotalSupply sdk.Int
}

func (msg DenomInfo) String() string {
	return fmt.Sprintf(`
  Denom:			 %s
  AssetHash:		 %s
  Creator:        	 %s
  TotalSupply:		 %s
`, msg.Denom, msg.AssetHash, msg.Creator.String(), msg.TotalSupply.String())
}

type DenomCrossChainInfo struct {
	DenomInfo
	ToChainId   uint64
	ToAssetHash string
}

func (msg DenomCrossChainInfo) String() string {
	return msg.DenomInfo.String() + fmt.Sprintf(`
  ToChainId:       	 %d
  ToAssetHash:		 %s
`, msg.ToChainId, msg.ToAssetHash)
}
