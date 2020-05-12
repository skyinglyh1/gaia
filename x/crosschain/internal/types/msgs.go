package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/hex"
)

// Governance message types and routes
const (
	TypeMsgSyncGenesis = "sync_genesis"
	TypeMsgSyncHeaders = "sync_headers"

	TypeMsgBindProxyHash       = "bind_proxy_hash"
	TypeMsgBindAssetHash       = "bind_asset_hash"
	TypeMsgLock                = "lock"
	TypeMsgProcessCrossChainTx = "process_cross_chain_tx"
	TypeMsgCreateCoins         = "create_coins"
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

type MsgBindProxyParam struct {
	Signer        sdk.AccAddress
	TargetChainId uint64
	TargetHash    []byte
}

func NewMsgBindProxyParam(signer sdk.AccAddress, targetChainId uint64, targetHash []byte) MsgBindProxyParam {
	return MsgBindProxyParam{signer, targetChainId, targetHash}
}

//nolint
func (msg MsgBindProxyParam) Route() string { return RouterKey }
func (msg MsgBindProxyParam) Type() string  { return TypeMsgBindProxyHash }

// Implements Msg.
func (msg MsgBindProxyParam) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.TargetChainId <= 0 {
		return ErrInvalidChainId(DefaultCodespace, msg.TargetChainId)
	}
	if len(msg.TargetHash) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrEmptyTargetHash(DefaultCodespace, hex.EncodeToString(msg.TargetHash))
	}

	return nil
}

func (msg MsgBindProxyParam) String() string {
	return fmt.Sprintf(`Bind Proxy Hash Message:
  Signer:         %s
  TargetChainId:  %d
  TargetHash:     %s
`, msg.Signer.String(), msg.TargetChainId, hex.EncodeToString(msg.TargetHash))
}

// Implements Msg.
func (msg MsgBindProxyParam) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgBindProxyParam) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgBindAssetParam struct {
	Signer             sdk.AccAddress
	SourceAssetDenom   string
	TargetChainId      uint64
	TargetAssetHash    []byte
	Limit              sdk.Int
	IsTargetChainAsset bool
}

func NewMsgBindAssetParam(signer sdk.AccAddress, sourceAssetDenom string, targetChainId uint64, targetAssetHash []byte, limit sdk.Int, isTargetChainAsset bool) MsgBindAssetParam {
	return MsgBindAssetParam{signer, sourceAssetDenom, targetChainId, targetAssetHash, limit, isTargetChainAsset}
}

//nolint
func (msg MsgBindAssetParam) Route() string { return RouterKey }
func (msg MsgBindAssetParam) Type() string  { return TypeMsgBindAssetHash }

// Implements Msg.
func (msg MsgBindAssetParam) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.SourceAssetDenom == "" {
		return sdk.ErrInternal(fmt.Sprintf("SourceAssetDenom is empty"))
	}
	if msg.TargetChainId <= 0 {
		return ErrInvalidChainId(DefaultCodespace, msg.TargetChainId)
	}
	if len(msg.TargetAssetHash) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrEmptyTargetHash(DefaultCodespace, hex.EncodeToString(msg.TargetAssetHash))
	}
	if msg.Limit.IsNegative() {
		return sdk.ErrInternal(fmt.Sprintf("bind asset param limit should be positive"))
	}
	return nil
}

func (msg MsgBindAssetParam) String() string {
	return fmt.Sprintf(`Bind Proxy Hash Message:
  Signer:         %s
  SourceAssetDenom: %s
  TargetChainId:  %d
  TargetAssetHash:     %s
  Limit: %s
  IsTargetChainAsset: %t
`, msg.Signer.String(), msg.SourceAssetDenom, msg.TargetChainId, hex.EncodeToString(msg.TargetAssetHash), msg.Limit.String(), msg.IsTargetChainAsset)
}

// Implements Msg.
func (msg MsgBindAssetParam) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgBindAssetParam) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
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

func NewMsgProcessCrossChainTx(submitter sdk.AccAddress, fromChainId uint64, height uint32, proof string, header []byte) MsgProcessCrossChainTx {
	return MsgProcessCrossChainTx{submitter, fromChainId, height, proof, header}
}

//nolint
func (msg MsgProcessCrossChainTx) Route() string { return RouterKey }
func (msg MsgProcessCrossChainTx) Type() string  { return TypeMsgProcessCrossChainTx }

// Implements Msg.
func (msg MsgProcessCrossChainTx) ValidateBasic() sdk.Error {
	if msg.Submitter.Empty() {
		return sdk.ErrInvalidAddress(msg.Submitter.String())
	}
	if msg.FromChainId <= 0 {
		return sdk.ErrInternal(fmt.Sprintf("MsgCrossChaintx.FromChainId should be positive"))
	}
	if len(msg.Proof) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return sdk.ErrInternal(fmt.Sprintf("MsgCrossChaintx.Proof should not be empty"))
	}
	if len(msg.Header) == 0 {
		return sdk.ErrInternal(fmt.Sprintf("MsgCrossChainTx.Header should not be empty"))
	}
	return nil
}

func (msg MsgProcessCrossChainTx) String() string {
	return fmt.Sprintf(`Bind Proxy Hash Message:
  Submitter:         %s
  FromChainId: %d
  Height:  %d
  Proof:     %s
  Header: %s
`, msg.Submitter.String(), msg.FromChainId, msg.Height, msg.Proof, hex.EncodeToString(msg.Header))
}

// Implements Msg.
func (msg MsgProcessCrossChainTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgProcessCrossChainTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Submitter}
}

type MsgCreateCoins struct {
	Creator sdk.AccAddress
	Coins   sdk.Coins
}

func NewMsgCreateCoins(creator sdk.AccAddress, coins sdk.Coins) MsgCreateCoins {
	return MsgCreateCoins{Creator: creator, Coins: coins}
}

//nolint
func (msg MsgCreateCoins) Route() string { return RouterKey }
func (msg MsgCreateCoins) Type() string  { return TypeMsgCreateCoins }

// Implements Msg.
func (msg MsgCreateCoins) ValidateBasic() sdk.Error {
	if msg.Coins.Empty() {
		return sdk.ErrInvalidAddress(msg.Creator.String())
	}
	//if !msg.Coins.IsZero() {
	//	return sdk.ErrInternal(fmt.Sprintf("Coins is Not Zero"))
	//}
	return nil
}

func (msg MsgCreateCoins) String() string {
	return fmt.Sprintf(`Create Coins Message:
  Creator:         %s
  Coins: 		   %s
`, msg.Creator.String(), msg.Coins.String())
}

// Implements Msg.
func (msg MsgCreateCoins) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgCreateCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

type MsgSetRedeemScript struct {
	Operator     sdk.AccAddress
	Denom        string
	RedeemKey    []byte
	RedeemScript []byte
}

func NewMsgSetRedeemScript(operator sdk.AccAddress, denom string, redeemKey, redeemScript []byte) MsgSetRedeemScript {
	return MsgSetRedeemScript{Operator: operator, Denom: denom, RedeemKey: redeemKey, RedeemScript: redeemScript}
}

//nolint
func (msg MsgSetRedeemScript) Route() string { return RouterKey }
func (msg MsgSetRedeemScript) Type() string  { return TypeMsgCreateCoins }

// Implements Msg.
func (msg MsgSetRedeemScript) ValidateBasic() sdk.Error {
	if len(msg.RedeemKey) == 0 {
		return sdk.ErrInternal(fmt.Sprintf("empty redeem key"))
	}
	if len(msg.RedeemScript) == 0 {
		return sdk.ErrInternal(fmt.Sprintf("empty redeem script"))
	}
	return nil
}

func (msg MsgSetRedeemScript) String() string {
	return fmt.Sprintf(`Create Coins Message:
  Operator:         %s
  RedeemKey: 		%x
  RedeemScript: 	%x
`, msg.Operator.String(), msg.RedeemKey, msg.RedeemScript)
}

// Implements Msg.
func (msg MsgSetRedeemScript) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgSetRedeemScript) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}

type MsgBindNoVMChainAssetHash struct {
	Signer          sdk.AccAddress
	Denom           string
	TargetChainId   uint64
	TargetAssetHash []byte
	Limit           sdk.Int
}

func NewMsgBindNoVMChainAssetHash(signer sdk.AccAddress, denom string, targetChainId uint64, targetAssetHash []byte, limit sdk.Int) MsgBindNoVMChainAssetHash {
	return MsgBindNoVMChainAssetHash{signer, denom, targetChainId, targetAssetHash, limit}
}

//nolint
func (msg MsgBindNoVMChainAssetHash) Route() string { return RouterKey }
func (msg MsgBindNoVMChainAssetHash) Type() string  { return TypeMsgBindAssetHash }

// Implements Msg.
func (msg MsgBindNoVMChainAssetHash) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.Denom == "" {
		return sdk.ErrInternal(fmt.Sprintf("SourceAssetDenom is empty"))
	}
	if msg.TargetChainId <= 0 {
		return ErrInvalidChainId(DefaultCodespace, msg.TargetChainId)
	}
	if len(msg.TargetAssetHash) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrEmptyTargetHash(DefaultCodespace, hex.EncodeToString(msg.TargetAssetHash))
	}
	if msg.Limit.IsNegative() {
		return sdk.ErrInternal(fmt.Sprintf("bind asset param limit should be positive"))
	}
	return nil
}

func (msg MsgBindNoVMChainAssetHash) String() string {
	return fmt.Sprintf(`BindNoVMChainAssetHash message:
  Signer:         %s
  Denom: %s
  TargetChainId:  %d
  TargetAssetHash:     %s
  Limit: %s
`, msg.Signer.String(), msg.Denom, msg.TargetChainId, hex.EncodeToString(msg.TargetAssetHash), msg.Limit.String())
}

// Implements Msg.
func (msg MsgBindNoVMChainAssetHash) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgBindNoVMChainAssetHash) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
