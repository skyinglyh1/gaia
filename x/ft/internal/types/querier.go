package types

const (
	QueryDenom           = "denom_info"
	QueryDenomCrossChain = "denom_cc_info"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryDenomInfo struct {
	Denom string
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryDenomInfo(denom string) QueryDenomInfo {
	return QueryDenomInfo{denom}
}

type QueryDenomCrossChainInfo struct {
	Denom   string
	ChainId uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryDenomCrossChainInfo(denom string, toChainId uint64) QueryDenomCrossChainInfo {
	return QueryDenomCrossChainInfo{denom, toChainId}
}
