package types

// Minting module event types
const (
	AttributeValueCategory = ModuleName

	EventTypeCreateCrossChainTx = "make_from_cosmos_proof"
	AttributeCrossChainId       = "cross_chainId"
	AttributeKeyTxParamHash     = "make_tx_param_hash"
	AttributeKeyMakeTxParam     = "make_tx_param"

	EventTypeVerifyToCosmosProof                        = "verify_to_cosmos_proof"
	AttributeKeyMerkleValueTxHash                       = "merkle_value.txhash"
	AttributeKeyMerkleValueMakeTxParamTxHash            = "merkle_value.make_tx_param.txhash"
	AttributeKeyMerkleValueMakeTxParamToContractAddress = "merkle_value.make_tx_param.to_contract_address"
	AttributeKeyFromChainId                             = "from_chain_id"
)
