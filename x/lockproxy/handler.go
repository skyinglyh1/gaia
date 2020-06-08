package lockproxy

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gaia/x/lockproxy/internal/keeper"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
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
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateLockProxy(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateLockProxy) (*sdk.Result, error) {

	//err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)

	err := k.CreateLockProxy(ctx, msg.Creator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBindProxyHash(ctx sdk.Context, k keeper.Keeper, msg types.MsgBindProxyHash) (*sdk.Result, error) {
	if err := k.BindProxyHash(ctx, msg.Operator, msg.ToChainId, msg.ToChainProxyHash); err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBindAssetHash(ctx sdk.Context, k keeper.Keeper, msg types.MsgBindAssetHash) (*sdk.Result, error) {

	err := k.BindAssetHash(ctx, msg.Operator, msg.SourceAssetDenom, msg.ToChainId, msg.ToAssetHash, msg.InitialAmt)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("handleMsgBindAssetHash, %v", err))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgLock(ctx sdk.Context, k keeper.Keeper, msg types.MsgLock) (*sdk.Result, error) {

	err := k.Lock(ctx, msg.LockProxyHash, msg.FromAddress, msg.SourceAssetDenom, msg.ToChainId, msg.ToAddressBs, *msg.Value)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("handleMsgLock, %v", err))
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
