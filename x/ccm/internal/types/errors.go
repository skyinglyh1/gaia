//nolint
package types

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const DefaultCodespace = ModuleName

var (
	ErrUnknownChainIdType           = sdkerrors.Register(DefaultCodespace, 1, "ErrUnknownChainIdType")
	ErrDeserializeHeaderFailType    = sdkerrors.Register(DefaultCodespace, 2, "ErrDeserializeHeaderFailType")
	ErrFindKeyHeightFailType        = sdkerrors.Register(DefaultCodespace, 3, "ErrFindKeyHeightFailType")
	ErrGetConsensusPeersFailType    = sdkerrors.Register(DefaultCodespace, 4, "ErrGetConsensusPeersFailType")
	ErrBookKeeperNumErrType         = sdkerrors.Register(DefaultCodespace, 5, "ErrBookKeeperNumErrType")
	ErrInvalidPublicKeyType         = sdkerrors.Register(DefaultCodespace, 6, "ErrInvalidPublicKeyType")
	ErrInvalidMultiSignatureType    = sdkerrors.Register(DefaultCodespace, 7, "ErrInvalidMultiSignatureType")
	ErrUnmarshalBlockInfoFailType   = sdkerrors.Register(DefaultCodespace, 8, "ErrUnmarshalBlockInfoFailType")
	ErrMarshalSpecificType          = sdkerrors.Register(DefaultCodespace, 9, "ErrMarshalSpecificType")
	ErrUnmarshalSpecificType        = sdkerrors.Register(DefaultCodespace, 10, "ErrUnmarshalSpecificType")
	ErrEmptyTargetHashType          = sdkerrors.Register(DefaultCodespace, 11, "ErrEmptyTargetHashType")
	ErrProposalHandlerNotExists     = sdkerrors.Register(DefaultCodespace, 12, "ErrProposalHandlerNotExists")
	ErrBelowCrossedLimitType        = sdkerrors.Register(DefaultCodespace, 13, "ErrBelowCrossedLimitType")
	ErrCrossedAmountOverflowType    = sdkerrors.Register(DefaultCodespace, 14, "ErrCrossedAmountOverflowType")
	ErrSupplyKeeperMintCoinFailType = sdkerrors.Register(DefaultCodespace, 15, "ErrSupplyKeeperMintCoinFailType")
	ErrSendCoinsToModuleFailType    = sdkerrors.Register(DefaultCodespace, 16, "ErrSendCoinsToModuleFailType")
	ErrSendCoinsFromModuleFailType  = sdkerrors.Register(DefaultCodespace, 17, "ErrSendCoinsFromModuleFailType")
	ErrCrossedAmountOverLimitType   = sdkerrors.Register(DefaultCodespace, 18, "ErrCrossedAmountOverLimitType")
	ErrGetCrossChainIdFailType      = sdkerrors.Register(DefaultCodespace, 19, "ErrGetCrossChainIdFailType")
	ErrCreateNegativeCoinsType      = sdkerrors.Register(DefaultCodespace, 20, "ErrCreateNegativeCoinsType")
)

func ErrDeserializeHeader(err error) error {
	return sdkerrors.Wrap(ErrDeserializeHeaderFailType, fmt.Sprintf("Header deserialization error:%s", err.Error()))
}

func ErrMarshalSpecificTypeFail(o interface{}, err error) error {
	return sdkerrors.Wrap(ErrUnmarshalSpecificType, fmt.Sprintf("marshal type: %s, error: %s", reflect.TypeOf(o).String(), err.Error()))
}

func ErrFindKeyHeight(height uint32, chainId uint64) error {
	return sdkerrors.Wrap(ErrFindKeyHeightFailType, fmt.Sprintf("findKeyHeight error: can not find key height with height:%d and chainId:%d", height, chainId))
}

func ErrGetConsensusPeers(height uint32, chainId uint64) error {
	return sdkerrors.Wrap(ErrGetConsensusPeersFailType, fmt.Sprintf("get consensus peers empty error: chainId: %d, height: %d", height, chainId))
}

func ErrBookKeeperNum(headerBookKeeperNum int, consensusNodeNum int) error {
	return sdkerrors.Wrap(ErrBookKeeperNumErrType, fmt.Sprintf("header Bookkeepers number:%d must more than 2/3 consensus node number:%d", headerBookKeeperNum, consensusNodeNum))
}

func ErrInvalidPublicKey(pubkey string) error {
	return sdkerrors.Wrap(ErrInvalidPublicKeyType, fmt.Sprintf("invalid pubkey error:%s", pubkey))
}

func ErrVerifyMultiSignatureFail(err error, height uint32) error {
	return sdkerrors.Wrap(ErrInvalidMultiSignatureType, fmt.Sprintf("verify multi signature error:%s, height:%d", err.Error(), height))
}

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

func ErrSupplyKeeperMintCoinsFail() error {
	return sdkerrors.Wrap(ErrSupplyKeeperMintCoinFailType, fmt.Sprintf("supplyKeeper mint coins failed"))
}

func ErrSendCoinsToModuleFail(amt sdk.Coins, fromAddr sdk.AccAddress, toAcct sdk.AccAddress) error {
	return sdkerrors.Wrap(ErrSendCoinsToModuleFailType, fmt.Sprintf("send coins:%s from account:%s to Module account:%s error", amt.String(), fromAddr.String(), toAcct.String()))
}

func ErrSendCoinsFromModuleFail(amt sdk.Coins, fromAddr sdk.AccAddress, toAcct sdk.AccAddress) error {
	return sdkerrors.Wrap(ErrSendCoinsFromModuleFailType, fmt.Sprintf("send coins:%s from Module account:%s to receiver account:%s error", amt.String(), fromAddr.String(), toAcct.String()))
}

func ErrCreateCrossChainTx(err error) error {
	return sdkerrors.Wrap(ErrGetCrossChainIdFailType, fmt.Sprintf("create cross chain tx error:%v", err))
}

func ErrCreateNegativeCoins(coins sdk.Coins) error {
	return sdkerrors.Wrap(ErrCreateNegativeCoinsType, fmt.Sprintf("create negative coins:%s", coins.String()))
}
