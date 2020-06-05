package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

const (
	ModuleStoreKey = "module_store_key"
	ToContract     = "to_contract"
	FromChainId    = "from_chain_id"
)

// RegisterRoutes registers minting module REST handlers on the provided router.
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	registerQueryRoutes(cliCtx, r, queryRoute)
	registerTxRoutes(cliCtx, r)
}
