package rest

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/gaia/x/ccm/client/common"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		fmt.Sprintf("/ccm/if_contain_contract/{%s}/{%s}/{%s}", ModuleStoreKey, ToContract, FromChainId),
		queryIfContainContract(cliCtx, queryRoute),
	).Methods("GET")

}

func queryIfContainContract(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)

		toContract, err := hex.DecodeString(vars[ToContract])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		fromChainId, err := strconv.ParseUint(vars[FromChainId], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, ok := checkResponseQueryIfContainContractResponse(w, cliCtx, queryRoute, vars[ModuleStoreKey], toContract, fromChainId)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryIfContainContractResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, keyStore string, toContract []byte, fromChainId uint64) (res []byte, ok bool) {

	res, err := common.QueryContainToContractAddr(cliCtx, queryRoute, keyStore, toContract, fromChainId)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
