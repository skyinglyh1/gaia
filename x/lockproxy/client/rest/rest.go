package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

const (
	Operator        = "operator"
	LockProxyHash   = "lock_proxy_hash"
	ToChainId       = "to_chain_id"
	AssetDenom      = "asset_denom"
	ToLockProxyHash = "to_lock_proxy_hash"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r)
}
