package types

// QueryBalanceParams defines the params for querying an account balance.
type QueryDenomInfo struct {
	Denom string
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryDenomInfo(denom string) QueryDenomInfo {
	return QueryDenomInfo{denom}
}

type QueryDenomInfoWithId struct {
	Denom   string
	ChainId uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryDenomInfoWithId(denom string, toChainId uint64) QueryDenomInfoWithId {
	return QueryDenomInfoWithId{denom, toChainId}
}
