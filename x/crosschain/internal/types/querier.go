package types



// QueryBalanceParams defines the params for querying an account balance.
type QueryHeaderParams struct {
	ChainId uint64
	Height  uint32
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryHeaderParams(chainId uint64, height uint32) QueryHeaderParams {
	return QueryHeaderParams{ChainId: chainId, Height: height}
}

type QueryCurrentHeightParams struct {
	ChainId uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryCurrentHeightParams(chainId uint64) QueryCurrentHeightParams {
	return QueryCurrentHeightParams{ChainId: chainId}
}


// QueryBalanceParams defines the params for querying an account balance.
type QueryProxyHashParam struct {
	ChainId uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryProxyHashParams(chainId uint64) QueryProxyHashParam {
	return QueryProxyHashParam{ChainId: chainId}
}

type QueryAssetHashParam struct {
	SourceAssetDenom string
	ChainId          uint64
}

func NewQueryAssetHashParams(sourceAssetDenom string, chainId uint64) QueryAssetHashParam {
	return QueryAssetHashParam{SourceAssetDenom: sourceAssetDenom, ChainId: chainId}
}

type QueryCrossedAmountParam struct {
	SourceAssetDenom string
	ChainId          uint64
}

func NewQueryCrossedAmountParam(sourceAssetDenom string, chainId uint64) QueryCrossedAmountParam {
	return QueryCrossedAmountParam{SourceAssetDenom: sourceAssetDenom, ChainId: chainId}
}

type QueryCrossedLimitParam struct {
	SourceAssetDenom string
	ChainId          uint64
}

func NewQueryCrossedLimitParam(sourceAssetDenom string, chainId uint64) QueryCrossedLimitParam {
	return QueryCrossedLimitParam{SourceAssetDenom: sourceAssetDenom, ChainId: chainId}
}

type QueryOperatorParam struct{}

func NewQueryOperatorParam() QueryOperatorParam {
	return QueryOperatorParam{}
}
