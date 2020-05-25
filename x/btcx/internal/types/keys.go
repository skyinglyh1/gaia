package types

// the one key to use for the keeper store

// nolint
const (
	// module name
	ModuleName = "btcx"

	// default paramspace for params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// Query endpoints supported by the minting querier
	QueryParameters = "parameters"
)
