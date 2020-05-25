package lockproxy

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/lockproxy/internal/keeper"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateLockProxy:
			return handleMsgCreateLockProxy(ctx, k, msg)
		case types.MsgBindProxyHash:
			return handleMsgBindProxyHash(ctx, k, msg)
		case types.MsgBindAssetHash:
			return handleMsgBindAssetHash(ctx, k, msg)
		case types.MsgLock:
			return handleMsgLock(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateLockProxy(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateLockProxy) sdk.Result {

	//err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)

	err := k.CreateLockProxy(ctx, msg.Creator)
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

func handleMsgBindProxyHash(ctx sdk.Context, k keeper.Keeper, msg types.MsgBindProxyHash) sdk.Result {
	if err := k.BindProxyHash(ctx, msg.Operator, msg.ToChainId, msg.ToChainProxyHash); err != nil {
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

	err := k.BindAssetHash(ctx, msg.Operator, msg.SourceAssetDenom, msg.TargetChainId, msg.TargetAssetHash, msg.InitialAmt)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("handleMsgBindAssetHash, %v", err)).Result()
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

	err := k.Lock(ctx, msg.LockProxyHash, msg.FromAddress, msg.SourceAssetDenom, msg.ToChainId, msg.ToAddressBs, *msg.Value)
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
