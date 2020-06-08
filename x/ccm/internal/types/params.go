package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	CurrentChainIdKey        = []byte("chainId")
	CurrentChainCrossChainId = uint64(11)
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
		params.NewParamSetPair(CurrentChainIdKey, &p.ChainId, validateChainId),
	}
}

func validateChainId(i interface{}) error {
	return nil
}
