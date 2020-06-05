package btcx

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/btcx/internal/keeper"
	"github.com/cosmos/gaia/x/btcx/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateDenom:
			return handleMsgCreateDenom(ctx, k, msg)
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

func handleMsgCreateDenom(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateDenom) sdk.Result {

	err := k.CreateDenom(ctx, msg.Creator, msg.Denom, msg.RedeemScript)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("handleMsgCreateDenom, error, %v", err)).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// These functions assume everything has been authenticated,
// now we just perform action and save

func handleMsgBindAssetHash(ctx sdk.Context, k keeper.Keeper, msg types.MsgBindAssetHash) sdk.Result {

	if err := k.BindAssetHash(ctx, msg.Creator, msg.SourceAssetDenom, msg.ToChainId, msg.ToAssetHash); err != nil {
		return sdk.ErrInternal(fmt.Sprintf("handleMsgBindAssetHash, error, %v", err)).Result()
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
