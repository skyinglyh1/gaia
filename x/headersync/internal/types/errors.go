//nolint
package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
)

const (
	CodeUnknownChainId sdk.CodeType = iota
	CodeDeserializeHeaderFailType
	CodeFindKeyHeightFailType
	CodeGetConsensusPeersFailType
	CodeBookKeeperNumErrType
	CodeInvalidPublicKeyType
	CodeInvalidMultiSignatureType
	CodeUnmarshalBlockInfoFailType
	CodeMarshalSpecificType
	CodeUnmarshalSpecificType
	CodeEmptyTargetHashType
	CodeProposalHandlerNotExists
	CodeBelowCrossedLimit
	CodeCrossedAmountOverflow
	CodeSupplyKeeperMintCoinFailType
	CodeSendCoinsToModuleFailType
	CodeSendCoinsFromModuleFailType
	CodeCrossedAmountOverLimitType
	CodeGetCrossChainIdFailType
	CodeCreateNegativeCoinsType

	DefaultCodespace sdk.CodespaceType = ModuleName
)

func ErrDeserializeHeader(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeDeserializeHeaderFailType, fmt.Sprintf("Header deserialization error:%s", err.Error()))
}

func ErrMarshalSpecificTypeFail(codespace sdk.CodespaceType, o interface{}, err error) sdk.Error {
	return sdk.NewError(codespace, CodeMarshalSpecificType, fmt.Sprintf("marshal type: %s, error: %s", reflect.TypeOf(o).String(), err.Error()))
}

func ErrUnmarshalSpecificTypeFail(codespace sdk.CodespaceType, o interface{}, err error) sdk.Error {
	return sdk.NewError(codespace, CodeUnmarshalSpecificType, fmt.Sprintf("marshal type: %s, error: %s", reflect.TypeOf(o).String(), err.Error()))
}

func ErrFindKeyHeight(codespace sdk.CodespaceType, height uint32, chainId uint64) sdk.Error {
	return sdk.NewError(codespace, CodeFindKeyHeightFailType, fmt.Sprintf("findKeyHeight error: can not find key height with height:%d and chainId:%d", height, chainId))
}

func ErrGetConsensusPeers(codespace sdk.CodespaceType, height uint32, chainId uint64) sdk.Error {
	return sdk.NewError(codespace, CodeGetConsensusPeersFailType, fmt.Sprintf("get consensus peers empty error: chainId: %d, height: %d", height, chainId))
}

func ErrBookKeeperNum(codespace sdk.CodespaceType, headerBookKeeperNum int, consensusNodeNum int) sdk.Error {
	return sdk.NewError(codespace, CodeBookKeeperNumErrType, fmt.Sprintf("header Bookkeepers number:%d must more than 2/3 consensus node number:%d", headerBookKeeperNum, consensusNodeNum))
}

func ErrInvalidPublicKey(codespace sdk.CodespaceType, pubkey string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPublicKeyType, fmt.Sprintf("invalid pubkey error:%s", pubkey))
}

func ErrVerifyMultiSignatureFail(codespace sdk.CodespaceType, err error, height uint32) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMultiSignatureType, fmt.Sprintf("verify multi signature error:%s, height:%d", err.Error(), height))
}

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
func ErrSupplyKeeperMintCoinsFail(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSupplyKeeperMintCoinFailType, fmt.Sprintf("supplyKeeper mint coins failed "))
}

func ErrSendCoinsToModuleFail(codespace sdk.CodespaceType, amt sdk.Coins, fromAddr sdk.AccAddress, toAcct sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeSendCoinsToModuleFailType, fmt.Sprintf("send coins:%s from account:%s to Module account:%s error", amt.String(), fromAddr.String(), toAcct.String()))
}

func ErrSendCoinsFromModuleFail(codespace sdk.CodespaceType, amt sdk.Coins, fromAddr sdk.AccAddress, toAcct sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeSendCoinsFromModuleFailType, fmt.Sprintf("send coins:%s from Module account:%s to receiver account:%s error", amt.String(), fromAddr.String(), toAcct.String()))
}

func ErrCreateCrossChainTx(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeGetCrossChainIdFailType, fmt.Sprintf("create cross chain tx error:%v", err))
}

func ErrCreateNegativeCoins(codespace sdk.CodespaceType, coins sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeCreateNegativeCoinsType, fmt.Sprintf("create negative coins:%s", coins.String()))
}
