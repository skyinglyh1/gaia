package ccm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gaia/x/ccm/internal/keeper"
	"github.com/cosmos/gaia/x/ccm/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgProcessCrossChainTx:
			return handleMsgProcessCrossChainTx(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized staking message type: %T", msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgProcessCrossChainTx(ctx sdk.Context, k keeper.Keeper, msg types.MsgProcessCrossChainTx) (*sdk.Result, error) {

	err := k.ProcessCrossChainTx(ctx, msg.FromChainId, msg.Height, msg.Proof, msg.Header)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("handleMsgCrossChainTx for unlock error, %v", err))
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
