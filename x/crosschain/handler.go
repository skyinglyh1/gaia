package crosschain

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/crosschain/internal/keeper"
	"github.com/cosmos/gaia/x/crosschain/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSyncGenesisParam:
			return handleMsgGenesisHeader(ctx, k, msg)
		case types.MsgSyncHeadersParam:
			return handleMsgBlockHeaders(ctx, k, msg)

		case types.MsgCreateCoins:
			return handleMsgCreateCoins(ctx, k, msg)
		case types.MsgBindProxyParam:
			return handleMsgBindProxyParam(ctx, k, msg)
		case types.MsgBindAssetParam:
			return handleMsgBindAssetParam(ctx, k, msg)
		case types.MsgLock:
			return handleMsgLock(ctx, k, msg)
		case types.MsgProcessCrossChainTx:
			return handleMsgProcessCrossChainTx(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgGenesisHeader(ctx sdk.Context, k keeper.Keeper, msg types.MsgSyncGenesisParam) sdk.Result {

	//err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)

	err := k.SyncGenesisHeader(ctx, msg.GenesisHeader)
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
func handleMsgBlockHeaders(ctx sdk.Context, k keeper.Keeper, msg types.MsgSyncHeadersParam) sdk.Result {
	err := k.SyncBlockHeaders(ctx, msg.Headers)
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

func handleMsgCreateCoins(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateCoins) sdk.Result {
	if !k.GetOperator(ctx).Operator.Empty() && !k.GetOperator(ctx).Operator.Equals(msg.Creator) {
		return sdk.ErrInternal(fmt.Sprintf("only operator can bind proxy hash, expected:%s, got:%s", k.GetOperator(ctx).Operator.String(), msg.Creator.String())).Result()
	}
	err := k.CreateCoins(ctx, msg.Creator, msg.Coins)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("handleMsgCrossChainTx for unlock error, %v", err)).Result()
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

func handleMsgBindProxyParam(ctx sdk.Context, k keeper.Keeper, msg types.MsgBindProxyParam) sdk.Result {

	if !k.GetOperator(ctx).Operator.Equals(msg.Signer) {
		return sdk.ErrInternal(fmt.Sprintf("only operator can bind proxy hash, expected:%s, got:%s", k.GetOperator(ctx).Operator.String(), msg.Signer.String())).Result()
	}
	k.BindProxyHash(ctx, msg.TargetChainId, msg.TargetHash)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBindAssetParam(ctx sdk.Context, k keeper.Keeper, msg types.MsgBindAssetParam) sdk.Result {

	if !k.GetOperator(ctx).Operator.Equals(msg.Signer) {
		return sdk.ErrInternal(fmt.Sprintf("only operator can bind proxy hash, expected:%s, got:%s", k.GetOperator(ctx).Operator.String(), msg.Signer.String())).Result()
	}
	err := k.BindAssetHash(ctx, msg.SourceAssetDenom, msg.TargetChainId, msg.TargetAssetHash, msg.Limit, msg.IsTargetChainAsset)
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

func handleMsgProcessCrossChainTx(ctx sdk.Context, k keeper.Keeper, msg types.MsgProcessCrossChainTx) sdk.Result {

	err := k.ProcessCrossChainTx(ctx, msg.FromChainId, msg.Height, msg.Proof, msg.Header)
	if err != nil {
		return sdk.ErrInternal(fmt.Sprintf("handleMsgCrossChainTx for unlock error, %v", err)).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
