package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/gaia/x/btcx/client/common"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		fmt.Sprintf("/btcx/denom_info/{%s}", Denom),
		queryDemonHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/btcx/denom_info/{%s}/{%s}", Denom, ChainId),
		queryDemonWithChainIdHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
}

func queryDemonHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)
		denom := vars[Denom]

		res, ok := checkResponseQueryDenomInfoResponse(w, cliCtx, queryRoute, denom)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDemonWithChainIdHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)
		denom := vars[Denom]
		chainId, ok := rest.ParseUint64OrReturnBadRequest(w, vars[ChainId])
		if !ok {
			return
		}
		res, ok := checkResponseQueryDenomInfoWithChainIdResponse(w, cliCtx, queryRoute, denom, chainId)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryDenomInfoResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, denom string) (res []byte, ok bool) {

	res, err := common.QueryDenomInfo(cliCtx, queryRoute, denom)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
func checkResponseQueryDenomInfoWithChainIdResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, denom string, chainId uint64) (res []byte, ok bool) {

	res, err := common.QueryDenomInfoWithId(cliCtx, queryRoute, denom, chainId)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
