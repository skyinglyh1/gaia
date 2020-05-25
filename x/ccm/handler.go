package ccm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gaia/x/ccm/internal/keeper"
	"github.com/cosmos/gaia/x/ccm/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgProcessCrossChainTx:
			return handleMsgProcessCrossChainTx(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
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
