package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	CurrentChainIdKey        = []byte("chainId")
	CurrentChainCrossChainId = uint64(10)
	MaxDenomsPerAccount      = 200
)

type ChainIdParam struct {
	ChainId uint64
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&ChainIdParam{})
}

func (p ChainIdParam) String() string {
	return fmt.Sprintf(
		`ChainId Params:
  Cross chain chain id : %d
`,
		p.ChainId,
	)
}

// Implements params.ParamSet
func (p *ChainIdParam) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{CurrentChainIdKey, &p.ChainId},
	}
}
