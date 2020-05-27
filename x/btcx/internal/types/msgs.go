package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/hex"
)

// Governance message types and routes
const (
	TypeMsgBindAssetHash = "bind_asset_hash"
	TypeMsgLock          = "lock"
	TypeMsgCreateCoin    = "create_coin"
)

type MsgCreateCoin struct {
	Creator      sdk.AccAddress
	Denom        string
	RedeemScript string
}

func NewMsgCreateCoin(creator sdk.AccAddress, denom string, redeemScript string) MsgCreateCoin {
	return MsgCreateCoin{Creator: creator, Denom: denom, RedeemScript: redeemScript}
}

//nolint
func (msg MsgCreateCoin) Route() string { return RouterKey }
func (msg MsgCreateCoin) Type() string  { return TypeMsgCreateCoin }

// Implements Msg.
func (msg MsgCreateCoin) ValidateBasic() sdk.Error {
	if msg.Creator.Empty() {
		return sdk.ErrInternal(fmt.Sprintf("MsgCreateCoin.Creator is empty"))
	}
	if _, err := sdk.ParseCoins("10" + msg.Denom); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("MsgCreateCoin.Denom is illegal, err:%v", err))
	}
	if _, err := hex.DecodeString(msg.RedeemScript); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("MsgCreateCoin.RedeemScript is not hex string format, err:%v", err))
	}
	return nil
}

func (msg MsgCreateCoin) String() string {
	return fmt.Sprintf(`MsgCreateCoin:
  Creator:         %s
  Denom: 		   %s
  RedeemScript:    %s
`, msg.Creator.String(), msg.Denom, msg.RedeemScript)
}

// Implements Msg.
func (msg MsgCreateCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgCreateCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

type MsgBindAssetHash struct {
	Creator          sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToAssetHash      []byte
}

func NewMsgBindAssetParam(creator sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAssetHash []byte) MsgBindAssetHash {
	return MsgBindAssetHash{creator, sourceAssetDenom, toChainId, toAssetHash}
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
	if msg.ToChainId <= 0 {
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
  Creator:          %s
  SourceAssetDenom: %s
  TargetChainId:    %d
  TargetAssetHash:  %x
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
