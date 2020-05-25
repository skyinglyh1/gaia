package types

import (
	"fmt"
	"github.com/cosmos/gaia/x/ccm"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/hex"
)

// Governance message types and routes
const (
	TypeMsgCreateAndDelegateCoinToProxy = "create_delegate_to_proxy"
	TypeMsgCreateDenom                  = "create_denom"

	TypeMsgBindProxyHash       = "bind_proxy_hash"
	TypeMsgBindAssetHash       = "bind_asset_hash"
	TypeMsgLock                = "lock"
	TypeMsgProcessCrossChainTx = "process_cross_chain_tx"
)

// MsgSend - high level transaction of the coin module
type MsgCreateAndDelegateCoinToProxy struct {
	Creator sdk.AccAddress
	Coin    sdk.Coin
}

var _ sdk.Msg = MsgCreateAndDelegateCoinToProxy{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgCreateAndDelegateCoinToProxy(creator sdk.AccAddress, coin sdk.Coin) MsgCreateAndDelegateCoinToProxy {
	return MsgCreateAndDelegateCoinToProxy{Creator: creator, Coin: coin}
}

// Route Implements Msg.
func (msg MsgCreateAndDelegateCoinToProxy) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateAndDelegateCoinToProxy) Type() string { return TypeMsgCreateAndDelegateCoinToProxy }

// ValidateBasic Implements Msg.
func (msg MsgCreateAndDelegateCoinToProxy) ValidateBasic() sdk.Error {
	if msg.Creator.Empty() {
		return sdk.ErrInvalidAddress(msg.Creator.String())
	}
	if !msg.Coin.IsValid() {
		return sdk.ErrInternal(fmt.Sprintf("Invalid coin:%s", msg.Coin.String()))
	}
	return nil
}

// GetSigners Implements Msg.
func (msg MsgCreateAndDelegateCoinToProxy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// GetSignBytes Implements Msg.
func (msg MsgCreateAndDelegateCoinToProxy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

type MsgCreateDenom struct {
	Creator sdk.AccAddress
	Denom   string
}

func NewMsgCreateDenom(creator sdk.AccAddress, denom string) MsgCreateDenom {
	return MsgCreateDenom{Creator: creator, Denom: denom}
}

//nolint
func (msg MsgCreateDenom) Route() string { return RouterKey }
func (msg MsgCreateDenom) Type() string  { return TypeMsgCreateDenom }

// Implements Msg.
func (msg MsgCreateDenom) ValidateBasic() sdk.Error {
	if !msg.Creator.Empty() {
		return sdk.ErrInternal(fmt.Sprintf("MsgCreateDenom.Creator is empty"))
	}
	if _, err := sdk.ParseCoin(msg.Denom + "100"); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("MsgCreateDenom.Denom:%s is invalid", msg.Denom))
	}
	return nil
}

func (msg MsgCreateDenom) String() string {
	return fmt.Sprintf(`Create Coins Message:
  Creator:         %s
  Denom: 		   %s
`, msg.Creator.String(), msg.Denom)
}

// Implements Msg.
func (msg MsgCreateDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgCreateDenom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

type MsgBindAssetHash struct {
	Creator          sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToAssetHash      []byte
}

func NewMsgBindAssetHash(creator sdk.AccAddress, sourceAssetDenom string, targetChainId uint64, targetAssetHash []byte) MsgBindAssetHash {
	return MsgBindAssetHash{creator, sourceAssetDenom, targetChainId, targetAssetHash}
}

//nolint
func (msg MsgBindAssetHash) Route() string { return RouterKey }
func (msg MsgBindAssetHash) Type() string  { return TypeMsgBindAssetHash }

// Implements Msg.
func (msg MsgBindAssetHash) ValidateBasic() sdk.Error {
	if msg.Creator.Empty() {
		return sdk.ErrInvalidAddress(msg.Creator.String())
	}
	if msg.SourceAssetDenom == "" {
		return sdk.ErrInternal(fmt.Sprintf("SourceAssetDenom is empty"))
	}
	if msg.ToChainId <= 0 || msg.ToChainId == ccm.CurrentChainCrossChainId {
		return ErrInvalidChainId(DefaultCodespace, msg.ToChainId)
	}
	if len(msg.ToAssetHash) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrEmptyTargetHash(DefaultCodespace, hex.EncodeToString(msg.ToAssetHash))
	}
	return nil
}

func (msg MsgBindAssetHash) String() string {
	return fmt.Sprintf(`MsgBindAssetHash:
  DenomCreator:         %s
  SourceAssetDenom: 	%s
  TargetChainId:  		%d
  TargetAssetHash:      %x
`, msg.Creator.String(), msg.SourceAssetDenom, msg.ToChainId, msg.ToAssetHash)
}

// Implements Msg.
func (msg MsgBindAssetHash) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgBindAssetHash) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

type MsgLock struct {
	FromAddress      sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToAddressBs      []byte
	Value            *sdk.Int
}

func NewMsgLock(fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddress []byte, value *sdk.Int) MsgLock {
	return MsgLock{fromAddress, sourceAssetDenom, toChainId, toAddress, value}
}

//nolint
func (msg MsgLock) Route() string { return RouterKey }
func (msg MsgLock) Type() string  { return TypeMsgLock }

// Implements Msg.
func (msg MsgLock) ValidateBasic() sdk.Error {
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
	return fmt.Sprintf(`Bind Proxy Hash Message:
  FromAddress:         %s
  SourceAssetDenom: %s
  ToChainId:  %d
  ToAddress:     %s
  Value: %s
`, msg.FromAddress.String(), msg.SourceAssetDenom, msg.ToChainId, hex.EncodeToString(msg.ToAddressBs), msg.Value.String())
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

type MsgProcessCrossChainTx struct {
	Submitter   sdk.AccAddress
	FromChainId uint64
	Height      uint32
	Proof       string
	Header      []byte
}
