package types

// Minting module event types
const (
	AttributeValueCategory = ModuleName

	EventTypeCreateLockProxy = "create_lock_proxy"
	AttributeKeyCreator      = "creator"
	AttributeKeyProxyHash    = "lock_proxy_hash"

	EventTypeBindProxy           = "bind_proxy_hash"
	AttributeKeyLockProxy        = "lock_proxy_hash"
	AttributeKeyToChainId        = "to_chain_id"
	AttributeKeyToChainProxyHash = "to_chain_proxy_hash"

	EventTypeBindAsset           = "bind_asset_hash"
	AttributeKeySourceAssetDenom = "source_asset_denom"
	AttributeKeySourceAssetHash  = "source_asset_hash"
	AttributeKeyToChainAssetHash = "to_chain_asset_hash"
	AttributeKeyInitialAmt       = "initial_amt"
	EventTypeLock                = "lock"
	AttributeKeyFromAddress      = "from_address"
	AttributeKeyToAddress        = "to_address"
	AttributeKeyAmount           = "amount"

	AttributeKeyFromChainId = "from_chain_id"

	EventTypeUnlock              = "unlock"
	AttributeKeyFromContractHash = "from_contract_hash"
	AttributeKeyToAssetDenom     = "to_asset_denom"
)
