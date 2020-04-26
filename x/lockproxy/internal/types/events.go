package types

// Minting module event types
const (
	EventTypeBindProxy           = "bind_proxy_hash"
	AttributeKeyToChainId        = "to_chain_id"
	AttributeKeyToChainProxyHash = "to_chain_proxy_hash"

	EventTypeBindAsset           = "bind_asset_hash"
	AttributeKeySourceAssetDenom = "source_asset_denom"
	AttributeKeyFromAssetHash    = "from_asset_hash"
	AttributeKeyToChainAssetHash = "to_chain_asset_hash"

	EventTypeLock           = "lock"
	AttributeKeyFromAddress = "from_address"
	AttributeKeyToAddress   = "to_address"
	AttributeKeyAmount      = "amount"

	EventTypeCreateCrossChainTx = "make_from_cosmos_proof"
	AttributeCrossChainId       = "cross_chainId"
	AttributeKeyTxParamHash     = "make_tx_param_hash"
	AttributeKeyMakeTxParam     = "make_tx_param"

	EventTypeVerifyToCosmosProof                        = "verify_to_cosmos_proof"
	AttributeKeyMerkleValueTxHash                       = "merkle_value.txhash"
	AttributeKeyMerkleValueMakeTxParamTxHash            = "merkle_value.make_tx_param.txhash"
	AttributeKeyMerkleValueMakeTxParamToContractAddress = "merkle_value.make_tx_param.to_contract_address"
	AttributeKeyFromChainId                             = "from_chain_id"
	AtttributeKeyStatus = "status"

	EventTypeUnlock              = "unlock"
	AttributeKeyFromContractHash = "from_contract_hash"
	AttributeKeyToAssetDenom     = "to_asset_denom"

	AttributeValueCategory = ModuleName
)
