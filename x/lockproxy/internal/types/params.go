package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyMintDenom             = []byte("MintDenom")
	KeyInflationRateChange   = []byte("InflationRateChange")
	KeyInflationMax          = []byte("InflationMax")
	KeyInflationMin          = []byte("InflationMin")
	KeyGoalBonded            = []byte("GoalBonded")
	KeyBlocksPerYear         = []byte("BlocksPerYear")
	KeyCoins                 = []byte("lockproxycoins")
	CurrentChainCrossChainId = uint64(6)
)

type Operator struct {
	Operator sdk.AccAddress
}

type CoinsParam struct {
	Coins sdk.Coins
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&CoinsParam{})
}

func NewCoinsParam(denom string) CoinsParam {
	denomCoin := sdk.Coin{denom, sdk.NewInt(0)}
	return CoinsParam{
		Coins: sdk.NewCoins(denomCoin),
	}
}

// default minting module parameters
func DefaultCoins() CoinsParam {
	return CoinsParam{
		Coins: sdk.Coins{
			sdk.Coin{
				Denom:  "MySimpleToken1",
				Amount: sdk.NewInt(10),
			},
			sdk.Coin{
				Denom:  "MySimpleToken2",
				Amount: sdk.NewInt(20),
			},
		},
	}

}

// validate params
func ValidateOperator(operator Operator) error {

	//if operator.Operator.Empty() {
	//	return fmt.Errorf("lockproxy parameter Operator can't be an empty string")
	//}
	return nil
}

func (p CoinsParam) String() string {
	return fmt.Sprintf(
		`Minting Params:
  Cross chain currently supported Coins : %s
`,
		p.Coins.String(),
	)
}

// Implements params.ParamSet
func (p *CoinsParam) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyCoins, &p.Coins},
	}
}
