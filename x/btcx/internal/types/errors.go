//nolint
package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const DefaultCodespace = ModuleName

var (
	ErrUnknownChainIdType         = sdkerrors.Register(DefaultCodespace, 1, "ErrUnknownChainIdType")
	ErrEmptyTargetHashType        = sdkerrors.Register(DefaultCodespace, 2, "ErrEmptyTargetHashType")
	ErrBelowCrossedLimitType      = sdkerrors.Register(DefaultCodespace, 3, "ErrBelowCrossedLimitType")
	ErrCrossedAmountOverflowType  = sdkerrors.Register(DefaultCodespace, 4, "ErrCrossedAmountOverflowType")
	ErrMintCoinFailType           = sdkerrors.Register(DefaultCodespace, 5, "ErrMintCoinFailType")
	ErrCrossedAmountOverLimitType = sdkerrors.Register(DefaultCodespace, 6, "ErrCrossedAmountOverLimitType")
)

func ErrInvalidChainId(chainId uint64) error {
	return sdkerrors.Wrap(ErrUnknownChainIdType, fmt.Sprintf("unknown chainId with id %d", chainId))
}

func ErrEmptyTargetHash(targetHashStr string) error {
	return sdkerrors.Wrap(ErrEmptyTargetHashType, fmt.Sprintf("empty target asset hash %s", targetHashStr))
}

func ErrBelowCrossedLimit(limit sdk.Int, storedCrossedLimit sdk.Int) error {
	return sdkerrors.Wrap(ErrBelowCrossedLimitType, fmt.Sprintf("new Limit:%s should be greater than stored Limit:%s", limit.String(), storedCrossedLimit.String()))
}

func ErrCrossedAmountOverLimit(newCrossedAmount sdk.Int, crossedLimit sdk.Int) error {
	return sdkerrors.Wrap(ErrCrossedAmountOverLimitType, fmt.Sprintf("new crossedAmount:%s is over the crossedLimit:%s", newCrossedAmount.String(), crossedLimit.String()))
}

func ErrCrossedAmountOverflow(newCrossedAmount sdk.Int, storedCrossedAmount sdk.Int) error {
	return sdkerrors.Wrap(ErrCrossedAmountOverflowType, fmt.Sprintf("new crossedAmount:%s is not greater than stored crossed amount:%s", newCrossedAmount.String(), storedCrossedAmount.String()))
}

func ErrMintCoinsFail() error {
	return sdkerrors.Wrap(ErrMintCoinFailType, fmt.Sprintf(" mint coins failed "))
}
