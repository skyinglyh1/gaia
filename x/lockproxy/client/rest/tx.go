package rest

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"encoding/hex"
	"errors"
	"github.com/cosmos/gaia/x/lockproxy/internal/types"
	"strconv"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/lockproxy/create_lock_proxy", createLockProxyRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/lockproxy/bind_proxy/{%s}/{%s}", ToChainId, ToLockProxyHash), bindProxyRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/lockproxy/bind_asset", ToChainId, ToLockProxyHash), bindAssetRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/lockproxy/lock", ToChainId, ToLockProxyHash), lockRequestHandlerFn(cliCtx)).Methods("POST")

}

type BaseReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
}

type BindAssetHashReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Denom       string       `json:"denom" yaml:"denom"`
	ToChainId   uint64       `json:"to_chain_id" yaml:"to_chain_id"`
	ToAssetHash []byte       `json:"to_asset_hash" yaml:"to_asset_hash"`
	InitialAmt  *big.Int     `json:"initial_amt" yaml:"initial_amt"`
}

type LockReq struct {
	BaseReq   rest.BaseReq `json:"base_req" yaml:"base_req"`
	LockProxy []byte       `json:"lock_proxy" yaml:"lock_proxy"`
	Denom     string       `json:"denom" yarml:"denom"`
	ToChainId uint64       `json:"to_chain_id" yaml:"to_chain_id"`
	ToAddress []byte       `json:"to_address" yaml:"to_address"`
	Amount    *big.Int     `json:"amount" yaml:"amount"`
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func createLockProxyRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req BaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgCreateLockProxy(cliCtx.GetFromAddress())
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func bindProxyRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		toProxyHashStr := vars[ToLockProxyHash]

		targetChainId, err := strconv.ParseUint(vars[ToChainId], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if toProxyHashStr[0:2] == "0x" {
			toProxyHashStr = toProxyHashStr[2:]
		}
		targetProxyHash, err := hex.DecodeString(toProxyHashStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("decode hex string 'toProxyHash' error:"+err.Error()).Error())
			return
		}

		var req BaseReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgBindProxyHash(cliCtx.GetFromAddress(), targetChainId, targetProxyHash)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func bindAssetRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BindAssetHashReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgBindAssetHash(cliCtx.GetFromAddress(), req.Denom, req.ToChainId, req.ToAssetHash, sdk.NewIntFromBigInt(req.InitialAmt))
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func lockRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LockReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgLock(cliCtx.GetFromAddress(), req.LockProxy, req.Denom, req.ToChainId, req.ToAddress, sdk.NewIntFromBigInt(req.Amount))
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
