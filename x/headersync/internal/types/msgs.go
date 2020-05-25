package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Governance message types and routes
const (
	TypeMsgSyncGenesis = "sync_genesis"
	TypeMsgSyncHeaders = "sync_headers"
)

// MsgSend - high level transaction of the coin module
type MsgSyncGenesisParam struct {
	Syncer        sdk.AccAddress
	GenesisHeader []byte
}

var _ sdk.Msg = MsgSyncGenesisParam{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSyncGenesisParam(syncer sdk.AccAddress, genesisHeader []byte) MsgSyncGenesisParam {
	return MsgSyncGenesisParam{Syncer: syncer, GenesisHeader: genesisHeader}
}

// Route Implements Msg.
func (msg MsgSyncGenesisParam) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSyncGenesisParam) Type() string { return TypeMsgSyncGenesis }

// ValidateBasic Implements Msg.
func (msg MsgSyncGenesisParam) ValidateBasic() sdk.Error {
	if msg.Syncer.Empty() {
		return sdk.ErrInvalidAddress(msg.Syncer.String())
	}
	if len(msg.GenesisHeader) == 0 {
		return sdk.ErrInvalidAddress("missing GenesisHeader bytes")
	}
	return nil
}

// GetSigners Implements Msg.
func (msg MsgSyncGenesisParam) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Syncer}
}

// GetSignBytes Implements Msg.
func (msg MsgSyncGenesisParam) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// MsgMultiSend - high level transaction of the coin module
type MsgSyncHeadersParam struct {
	Syncer  sdk.AccAddress
	Headers [][]byte
}

var _ sdk.Msg = MsgSyncHeadersParam{}

// NewMsgMultiSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSyncHeadersParam(syncer sdk.AccAddress, headers [][]byte) MsgSyncHeadersParam {
	return MsgSyncHeadersParam{Syncer: syncer, Headers: headers}
}

// Route Implements Msg
func (msg MsgSyncHeadersParam) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgSyncHeadersParam) Type() string { return TypeMsgSyncHeaders }

// ValidateBasic Implements Msg.
func (msg MsgSyncHeadersParam) ValidateBasic() sdk.Error {
	if msg.Syncer.Empty() {
		return sdk.ErrInvalidAddress(msg.Syncer.String())
	}
	if len(msg.Headers) == 0 {
		return sdk.ErrInvalidAddress("missing BlockHeaders bytes")
	}
	if len(msg.Headers[0]) == 0 {
		return sdk.ErrInvalidAddress("missing BlockHeader bytes")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSyncHeadersParam) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSyncHeadersParam) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Syncer}
}
