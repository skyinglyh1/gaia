package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		"/corsschain/proxyhash/{chainId}",
		queryProxyHashHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

}

func queryProxyHashHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// do nothing
	}
}

func checkResponseQueryProxyHashResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, chainId uint64) (res []byte, ok bool) {

	// do nothing

	return res, true
}
