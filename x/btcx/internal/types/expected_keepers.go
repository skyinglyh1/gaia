package types // noalias

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account

	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	GetAllAccounts(ctx sdk.Context) []authexported.Account
	SetAccount(ctx sdk.Context, acc authexported.Account)

	IterateAccounts(ctx sdk.Context, process func(authexported.Account) bool)
}

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, exported.ModuleAccountI)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	SetSupply(ctx sdk.Context, supply exported.SupplyI)
	GetSupply(ctx sdk.Context) (supply exported.SupplyI)
}

type SupplyI interface {
	SetTotal(total sdk.Coins) SupplyI
}

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	DelegateCoins(ctx sdk.Context, fromAdd, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	UndelegateCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error

	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
}

type CCMKeeper interface {
	ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, height uint32, proofStr string, headerBs []byte) sdk.Error
	CreateCrossChainTx(ctx sdk.Context, toChainId uint64, fromContractHash, toContractHash []byte, method string, args []byte) sdk.Error
	SetDenomCreator(ctx sdk.Context, denom string, creator sdk.AccAddress)
	GetDenomCreator(ctx sdk.Context, denom string) sdk.AccAddress
}
