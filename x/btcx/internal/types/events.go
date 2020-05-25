package types

// Minting module event types
const (
	AttributeValueCategory = ModuleName

	AttributeKeyToChainId = "to_chain_id"

	EventTypeBindAsset  = "bind_asset_hash"
	AttributeKeyCreator = "creator"

	AttributeKeyFromAssetHash    = "from_asset_hash"
	AttributeKeyToChainAssetHash = "to_chain_asset_hash"

	EventTypeLock                = "lock"
	AttributeKeySourceAssetHash  = "source_asset_hash"
	AttributeKeySourceAssetDenom = "source_asset_denom"
	AttributeKeyFromAddress      = "from_address"
	AttributeKeyToAddress        = "to_address"
	AttributeKeyAmount           = "amount"

	AttributeKeyFromChainId = "from_chain_id"
	AtttributeKeyStatus     = "status"

	EventTypeUnlock              = "unlock"
	AttributeKeyFromContractHash = "from_contract_hash"
	AttributeKeyToAssetDenom     = "to_asset_denom"

	EventTypeCreateCoin      = "create_coin"
	AttributeKeyRedeemKey    = "redeem_key"
	AttributeKeyRedeemScript = "redeem_script"
)
