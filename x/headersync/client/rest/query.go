package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/gaia/x/headersync/client/common"
	"strconv"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		"/headersync/header/{chainId}/{height}",
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
		chainIdStr := vars["chainId"]
		chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		heightStr := vars["height"]
		height, err := strconv.ParseUint(heightStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, ok := checkResponseQueryProxyHashResponse(w, cliCtx, queryRoute, chainId, uint32(height))
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryProxyHashResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, chainId uint64, height uint32) (res []byte, ok bool) {

	res, err := common.QueryHeader(cliCtx, queryRoute, chainId, height)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
