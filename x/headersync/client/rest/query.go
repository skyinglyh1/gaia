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
	r.HandleFunc("/headersync/header/{chainId}/{height}", QueryHeaderRequestHandlerFn(cliCtx, queryRoute)).Methods("GET")
	r.HandleFunc("/headersync/currentheight/{chainId}", QueryCurrentHeaderHeightRequestHandlerFn(cliCtx, queryRoute)).Methods("GET")
}

func QueryHeaderRequestHandlerFn(cliCtx context.CLIContext,queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// query for rewards from a particular delegator
		vars := mux.Vars(r)
		chainIdStr := vars["chainId"]
		heightStr := vars["height"]

		chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		height, err := strconv.ParseUint(heightStr, 10, 32)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		res, ok := checkResponseQueryHeaderResponse(w, cliCtx, queryRoute, chainId, uint32(height))
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}


func QueryCurrentHeaderHeightRequestHandlerFn(cliCtx context.CLIContext,queryRoute string) http.HandlerFunc {
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

		res, ok := checkResponseQueryCurrentHeaderHeightResponse(w, cliCtx, queryRoute, chainId)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryHeaderResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, chainId uint64, height uint32,
) (res []byte, ok bool) {

	res, err := common.QueryHeader(cliCtx, queryRoute, chainId, height)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
func checkResponseQueryCurrentHeaderHeightResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, chainId uint64) (res []byte, ok bool) {

	res, err := common.QueryCurrentHeaderHeight(cliCtx, queryRoute, chainId)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}