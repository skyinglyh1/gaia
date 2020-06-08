package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/gaia/x/btcx/internal/types"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/btcx/create_coin", createCoinRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/btcx/bind_asset_hash", bindAssetHashRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/btcx/lock", lockRequestHandlerFn(cliCtx)).Methods("POST")

}

// SendReq defines the properties of a send request's body.
type CreateCoinReq struct {
	BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
	Denom        string       `json:"denom" yaml:"denom"`
	RedeemScript string       `json:"redeem_script" yaml:"redeem_script"`
}

type BindAssetHashReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Denom       string       `json:"denom" yaml:"denom"`
	ToChainId   uint64       `json:"redeem_script" yaml:"redeem_script"`
	ToAssetHash []byte       `json:"to_asset_hash" yaml:"to_asset_hash"`
}

type LockReq struct {
	BaseReq          rest.BaseReq `json:"base_req" yaml:"base_req"`
	SourceAssetDenom string       `json:"source_asset_denom" yaml:"source_asset_denom"`
	ToChainId        uint64       `json:"to_Chain_id" yaml:"to_chain_id"`
	ToAddressBs      []byte       `json:"to_address_bs" yaml:"to_address_bs"`
	Value            *sdk.Int     `json:"value" yaml:"value"`
}

func createCoinRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateCoinReq
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
		msg := types.NewMsgCreateDenom(fromAddr, req.Denom, req.RedeemScript)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func bindAssetHashRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BindAssetHashReq
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
		msg := types.NewMsgBindAssetHash(fromAddr, req.Denom, req.ToChainId, req.ToAssetHash)
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
		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgLock(fromAddr, req.SourceAssetDenom, req.ToChainId, req.ToAddressBs, req.Value)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
