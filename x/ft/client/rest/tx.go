package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"encoding/hex"
	"errors"
	"github.com/cosmos/gaia/x/ft/internal/types"
	"strconv"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/crosschain/bindproxy/{targetChainId}/{targetProxyHash}", BindProxyRequestHandlerFn(cliCtx)).Methods("POST")

}

// SendReq defines the properties of a send request's body.
type SendReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func BindProxyRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		targetChainIdStr := vars["targetChainId"]
		targetProxyHashStr := vars["targetProxyHash"]
		denom := vars["denom"]

		targetChainId, err := strconv.ParseUint(targetChainIdStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if targetProxyHashStr[0:2] == "0x" {
			targetProxyHashStr = targetProxyHashStr[2:]
		}
		targetProxyHash, err := hex.DecodeString(targetProxyHashStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("decode hex string 'targetProxyHash' error:"+err.Error()).Error())
			return
		}

		var req SendReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgBindAssetHash(cliCtx.GetFromAddress(), denom, targetChainId, targetProxyHash)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
