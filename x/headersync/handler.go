package headersync

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/headersync/internal/keeper"
	"github.com/cosmos/gaia/x/headersync/internal/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSyncGenesisParam:
			return handleMsgGenesisHeader(ctx, k, msg)

		case types.MsgSyncHeadersParam:
			return handleMsgBlockHeaders(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized bank message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
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
