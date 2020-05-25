package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/ft/internal/types"
	"github.com/cosmos/gaia/x/lockproxy"
)

func (k Keeper) CreateAndDelegateCoinToProxy(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin) sdk.Error {

	if k.ExistDenom(ctx, coin.Denom) {
		return sdk.ErrInternal(fmt.Sprintf("CreateCoins Error: denom:%s already exist", coin.Denom))
	}
	//k.SetOperator(ctx, denom, creator)
	k.ccmKeeper.SetDenomCreator(ctx, coin.Denom, creator)

	if err := k.supplyKeeper.MintCoins(ctx, lockproxy.ModuleName, sdk.NewCoins(coin)); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("CreateAndDelegateCoinToProxy error: %v", err))
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateAndDelegateCoinToProxy,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, coin.Denom),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, coin.Amount.String()),
		),
	})
	k.Logger(ctx).Info(fmt.Sprintf("creator:%s initialized coin: %s ", creator.String(), coin.String()))
	return nil
}
