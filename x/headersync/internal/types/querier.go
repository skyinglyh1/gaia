package types

// QueryBalanceParams defines the params for querying an account balance.
type QueryHeaderParams struct {
	ChainId uint64
	Height uint32
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryHeaderParams(chainId uint64, height uint32) QueryHeaderParams{
	return QueryHeaderParams{ChainId: chainId, Height:height}
}
type QueryHeaderHeightParams struct {
	ChainId uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryHeaderHeightParams(chainId uint64) QueryHeaderParams{
	return QueryHeaderParams{ChainId: chainId}
}
