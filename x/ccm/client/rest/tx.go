package rest

import (
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/gaia/x/ccm/internal/types"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/ccm/processcrosschaintx", ProcessCrossChainTxRequestHandlerFn(cliCtx)).Methods("POST")

}

// SendReq defines the properties of a send request's body.
type ProcessCrossChainTxReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	FromChainId uint64       `json:"from_chain_id" yaml:"from_chain_id"`
	Height      uint32       `json:"height" yaml:"height"`
	Proof       string       `json:"proof" yaml:"proof"`
	Header      string       `json:"header" yaml:"header"`
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func ProcessCrossChainTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProcessCrossChainTxReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		header, err := hex.DecodeString(req.Header)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgProcessCrossChainTx(fromAddr, req.FromChainId, req.Height, req.Proof, header)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
