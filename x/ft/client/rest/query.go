package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/gaia/x/ft/client/common"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		"/corsschain/proxyhash/{chainId}",
		queryProxyHashHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

}

func queryProxyHashHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)
		denom := vars["denom"]
		res, ok := checkResponseQueryProxyHashResponse(w, cliCtx, queryRoute, denom)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryProxyHashResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, denom string) (res []byte, ok bool) {

	res, err := common.QueryDenomInfo(cliCtx, queryRoute, denom)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
