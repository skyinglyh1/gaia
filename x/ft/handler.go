package ft

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/ft/internal/keeper"
	"github.com/cosmos/gaia/x/ft/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateAndDelegateCoinToProxy:
			return handleMsgCreateAndDelegateCoinToProxy(ctx, k, msg)

		case types.MsgCreateDenom:
			return handleMsgCreateDenom(ctx, k, msg)

		case types.MsgBindAssetHash:
			return handleMsgBindAssetHash(ctx, k, msg)
		case types.MsgLock:
			return handleMsgLock(ctx, k, msg)
		case types.MsgCreateCoins:
			return handleMsgCreateCoins(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateAndDelegateCoinToProxy(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateAndDelegateCoinToProxy) sdk.Result {

	//err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)

	err := k.CreateAndDelegateCoinToProxy(ctx, msg.Creator, msg.Coin)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// Handle MsgMultiSend.
func handleMsgCreateDenom(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateDenom) sdk.Result {
	err := k.CreateDenom(ctx, msg.Creator, msg.Denom)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBindAssetHash(ctx sdk.Context, k keeper.Keeper, msg types.MsgBindAssetHash) sdk.Result {

	err := k.BindAssetHash(ctx, msg.Creator, msg.SourceAssetDenom, msg.ToChainId, msg.ToAssetHash)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgLock(ctx sdk.Context, k keeper.Keeper, msg types.MsgLock) sdk.Result {

	err := k.Lock(ctx, msg.FromAddress, msg.SourceAssetDenom, msg.ToChainId, msg.ToAddressBs, *msg.Value)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("handleMsgLock, %v", err)).Result()
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgCreateCoins(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateCoins) sdk.Result {
	coins, err := sdk.ParseCoins(msg.Coins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("handleMsgCreateCoins, parseCoins error:%v", err)).Result()
	}
	sdkErr := k.CreateCoins(ctx, msg.Creator, coins)
	if sdkErr != nil {
		return sdkErr.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
