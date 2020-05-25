package rest

import (
	"encoding/hex"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/gaia/x/lockproxy/client/common"
	"strconv"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		"/lockproxy/proxyhash/{lockproxy}/{chainId}",
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
		lockproxyStr := vars["lockproxy"]
		lockproxy, err := hex.DecodeString(lockproxyStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		chainIdStr := vars["chainId"]
		chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, ok := checkResponseQueryProxyHashResponse(w, cliCtx, queryRoute, lockproxy, chainId)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryProxyHashResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, lockproxy []byte, chainId uint64) (res []byte, ok bool) {

	res, err := common.QueryProxyHash(cliCtx, queryRoute, lockproxy, chainId)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
