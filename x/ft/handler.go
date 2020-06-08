package ft

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gaia/x/ft/internal/keeper"
	"github.com/cosmos/gaia/x/ft/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
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
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateAndDelegateCoinToProxy(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateAndDelegateCoinToProxy) (*sdk.Result, error) {

	//err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)

	err := k.CreateCoinAndDelegateToProxy(ctx, msg.Creator, msg.Coin, msg.LockProxyHash)
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

// Handle MsgMultiSend.
func handleMsgCreateDenom(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateDenom) (*sdk.Result, error) {
	err := k.CreateDenom(ctx, msg.Creator, msg.Denom)
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

func handleMsgBindAssetHash(ctx sdk.Context, k keeper.Keeper, msg types.MsgBindAssetHash) (*sdk.Result, error) {

	err := k.BindAssetHash(ctx, msg.Creator, msg.SourceAssetDenom, msg.ToChainId, msg.ToAssetHash)
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

func handleMsgLock(ctx sdk.Context, k keeper.Keeper, msg types.MsgLock) (*sdk.Result, error) {

	err := k.Lock(ctx, msg.FromAddress, msg.SourceAssetDenom, msg.ToChainId, msg.ToAddressBs, *msg.Value)
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

func handleMsgCreateCoins(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateCoins) (*sdk.Result, error) {
	coins, err := sdk.ParseCoins(msg.Coins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("handleMsgCreateCoins, parseCoins error:%v", err))
	}
	sdkErr := k.CreateCoins(ctx, msg.Creator, coins)
	if sdkErr != nil {
		return nil, sdkErr
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
