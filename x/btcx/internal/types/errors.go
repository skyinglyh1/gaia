//nolint
package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeUnknownChainId sdk.CodeType = iota
	CodeEmptyTargetHashType
	CodeBelowCrossedLimit
	CodeCrossedAmountOverflow
	CodeMintCoinFailType
	CodeCrossedAmountOverLimitType
	DefaultCodespace sdk.CodespaceType = ModuleName
)

func ErrInvalidChainId(codespace sdk.CodespaceType, chainId uint64) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownChainId, fmt.Sprintf("unknown chainId with id %d", chainId))
}

func ErrEmptyTargetHash(codespace sdk.CodespaceType, targetHashStr string) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyTargetHashType, fmt.Sprintf("empty target asset hash %s", targetHashStr))
}
func ErrBelowCrossedLimit(codespace sdk.CodespaceType, limit sdk.Int, storedCrossedLimit sdk.Int) sdk.Error {
	return sdk.NewError(codespace, CodeBelowCrossedLimit, fmt.Sprintf("new Limit:%s should be greater than stored Limit:%s", limit.String(), storedCrossedLimit.String()))
}

func ErrCrossedAmountOverLimit(codespace sdk.CodespaceType, newCrossedAmount sdk.Int, crossedLimit sdk.Int) sdk.Error {
	return sdk.NewError(codespace, CodeCrossedAmountOverLimitType, fmt.Sprintf("new crossedAmount:%s is over the crossedLimit:%s", newCrossedAmount.String(), crossedLimit.String()))
}

func ErrCrossedAmountOverflow(codespace sdk.CodespaceType, newCrossedAmount sdk.Int, storedCrossedAmount sdk.Int) sdk.Error {
	return sdk.NewError(codespace, CodeCrossedAmountOverflow, fmt.Sprintf("new crossedAmount:%s is not greater than stored crossed amount:%s", newCrossedAmount.String(), storedCrossedAmount.String()))
}
func ErrMintCoinsFail(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMintCoinFailType, fmt.Sprintf(" mint coins failed "))
}
