package lockproxy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - minter state
type GenesisState struct {
	Operator Operator `json:"operator" yaml:"operator"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(operator Operator) GenesisState {
	return GenesisState{
		Operator: operator,
	}
}

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetOperator(ctx, data.Operator)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	operator := keeper.GetOperator(ctx)
	return NewGenesisState(operator)
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := ValidateOperator(data.Operator)
	if err != nil {
		return err
	}
	return nil
}
