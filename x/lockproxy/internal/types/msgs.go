package types

import (
	"fmt"
	"github.com/cosmos/gaia/x/ccm"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/hex"
)

// Governance message types and routes
const (
	TypeMsgCreateLockProxy = "create_lock_proxy"

	TypeMsgBindProxyHash = "bind_proxy_hash"
	TypeMsgBindAssetHash = "bind_asset_hash"
	TypeMsgLock          = "lock"
)

// MsgSend - high level transaction of the coin module
type MsgCreateLockProxy struct {
	Creator sdk.AccAddress
}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgCreateLockProxy(creator sdk.AccAddress) MsgCreateLockProxy {
	return MsgCreateLockProxy{creator}
}

// Route Implements Msg.
func (msg MsgCreateLockProxy) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateLockProxy) Type() string { return TypeMsgCreateLockProxy }

// ValidateBasic Implements Msg.
func (msg MsgCreateLockProxy) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdk.ErrInvalidAddress(msg.Creator.String())
	}
	return nil
}

// GetSigners Implements Msg.
func (msg MsgCreateLockProxy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// GetSignBytes Implements Msg.
func (msg MsgCreateLockProxy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

type MsgBindProxyHash struct {
	Operator         sdk.AccAddress
	ToChainId        uint64
	ToChainProxyHash []byte
}

func NewMsgBindProxyHash(operator sdk.AccAddress, toChainId uint64, toChainProxyHash []byte) MsgBindProxyHash {
	return MsgBindProxyHash{operator, toChainId, toChainProxyHash}
}

//nolint
func (msg MsgBindProxyHash) Route() string { return RouterKey }
func (msg MsgBindProxyHash) Type() string  { return TypeMsgBindProxyHash }

// Implements Msg.
func (msg MsgBindProxyHash) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdk.ErrInvalidAddress(msg.Operator.String())
	}
	if msg.ToChainId <= 0 || msg.ToChainId == ccm.CurrentChainCrossChainId {
		return ErrInvalidChainId(DefaultCodespace, msg.ToChainId)
	}
	if len(msg.ToChainProxyHash) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrEmptyTargetHash(DefaultCodespace, hex.EncodeToString(msg.ToChainProxyHash))
	}

	return nil
}

func (msg MsgBindProxyHash) String() string {
	return fmt.Sprintf(`MsgBindProxyHash:
  Operator:       		%s(%x)
  ToChainId:			%d
  ToChainProxyHash:     %s
`, msg.Operator.String(), msg.Operator.Bytes(), msg.ToChainId, hex.EncodeToString(msg.ToChainProxyHash))
}

// Implements Msg.
func (msg MsgBindProxyHash) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgBindProxyHash) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}

type MsgBindAssetHash struct {
	Operator         sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToAssetHash      []byte
	InitialAmt       sdk.Int
}

func NewMsgBindAssetParam(operator sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAssetHash []byte, initialAmt sdk.Int) MsgBindAssetHash {
	return MsgBindAssetHash{operator, sourceAssetDenom, toChainId, toAssetHash, initialAmt}
}

//nolint
func (msg MsgBindAssetHash) Route() string { return RouterKey }
func (msg MsgBindAssetHash) Type() string  { return TypeMsgBindAssetHash }

// Implements Msg.
func (msg MsgBindAssetHash) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdk.ErrInvalidAddress(msg.Operator.String())
	}
	if msg.SourceAssetDenom == "" {
		return sdk.ErrInternal(fmt.Sprintf("SourceAssetDenom is empty"))
	}
	if msg.ToChainId <= 0 {
		return ErrInvalidChainId(DefaultCodespace, msg.ToChainId)
	}
	if len(msg.ToAssetHash) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrEmptyTargetHash(DefaultCodespace, hex.EncodeToString(msg.ToAssetHash))
	}
	if msg.InitialAmt.IsNegative() {
		return sdk.ErrInternal(fmt.Sprintf("bind asset param limit should be positive"))
	}
	return nil
}

func (msg MsgBindAssetHash) String() string {
	return fmt.Sprintf(`Bind Proxy Hash Message:
  Signer:         	%s
  SourceAssetDenom: %s
  ToChainId:  		%d
  ToAssetHash:      %s
  Limit: 			%s
`, msg.Operator.String(), msg.SourceAssetDenom, msg.ToChainId, hex.EncodeToString(msg.ToAssetHash), msg.InitialAmt.String())
}

// Implements Msg.
func (msg MsgBindAssetHash) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgBindAssetHash) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}

type MsgLock struct {
	LockProxyHash    []byte
	FromAddress      sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToAddressBs      []byte
	Value            *sdk.Int
}

func NewMsgLock(lockProxyHash []byte, fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddress []byte, value *sdk.Int) MsgLock {
	return MsgLock{lockProxyHash, fromAddress, sourceAssetDenom, toChainId, toAddress, value}
}

//nolint
func (msg MsgLock) Route() string { return RouterKey }
func (msg MsgLock) Type() string  { return TypeMsgLock }

// Implements Msg.
func (msg MsgLock) ValidateBasic() error {
	if len(msg.LockProxyHash) == 0 {
		return sdk.ErrInternal("passed in lockProxyHash is empty")
	}
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress(msg.FromAddress.String())
	}
	if msg.SourceAssetDenom == "" {
		return sdk.ErrInternal(fmt.Sprintf("SourceAssetDenom is empty"))
	}
	if msg.ToChainId <= 0 {
		return ErrInvalidChainId(DefaultCodespace, msg.ToChainId)
	}
	if len(msg.ToAddressBs) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrEmptyTargetHash(DefaultCodespace, hex.EncodeToString(msg.ToAddressBs))
	}
	if msg.Value.IsNegative() {
		return sdk.ErrInternal(fmt.Sprintf("bind asset param limit should be positive"))
	}
	return nil
}

func (msg MsgLock) String() string {
	return fmt.Sprintf(`MsgLock:
  LockProxyHash: 		%x
  FromAddress:          %s
  SourceAssetDenom:     %s
  ToChainId:            %d
  ToAddress:            %x
  Value:                %s
`, msg.LockProxyHash, msg.FromAddress.String(), msg.SourceAssetDenom, msg.ToChainId, msg.ToAddressBs, msg.Value.String())
}

// Implements Msg.
func (msg MsgLock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgLock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
