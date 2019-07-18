// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func accountLinkTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("remoteAccountKey", ByteSize),
			newScalarAttribute("linkAction", ByteSize),
		},
	}
}

func accountPropertyTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("propertyType", ByteSize),
			newScalarAttribute("modificationCount", ByteSize),
			newTableArrayAttribute("modifications", schema{
				[]schemaAttribute{
					newScalarAttribute("modificationType", ByteSize),
					newArrayAttribute("value", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func aliasTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("actionType", ByteSize),
			newArrayAttribute("namespaceId", IntSize),
			newArrayAttribute("aliasId", ByteSize),
		},
	}
}

func aggregateTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("transactionsSize", IntSize),
			newArrayAttribute("transactions", ByteSize),
		},
	}
}

func catapultConfigTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("applyHeightDelta", IntSize),
			newScalarAttribute("blockChainConfigSize", IntSize),
			newScalarAttribute("supportedEntityVersionsSize", IntSize),
			newArrayAttribute("blockChainConfig", ByteSize),
			newArrayAttribute("supportedEntityVersionsSize", ByteSize),
		},
	}
}

func catapultUpdateTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("upgradePeriod", IntSize),
			newArrayAttribute("newCatapultVersion", IntSize),
		},
	}
}

func modifyMetadataTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("metadataType", ByteSize),
			newArrayAttribute("metadataId", ByteSize),
			newTableArrayAttribute("modifications", schema{
				[]schemaAttribute{
					newScalarAttribute("size", IntSize),
					newScalarAttribute("modificationType", ByteSize),
					newScalarAttribute("keySize", ByteSize),
					newArrayAttribute("valueSize", ByteSize),
					newArrayAttribute("key", ByteSize),
					newArrayAttribute("value", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func mosaicDefinitionTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("mosaicNonce", IntSize),
			newArrayAttribute("mosaicId", IntSize),
			newScalarAttribute("numOptionalProperties", ByteSize),
			newScalarAttribute("flags", ByteSize),
			newScalarAttribute("divisibility", ByteSize),
			newTableArrayAttribute("modifications", schema{
				[]schemaAttribute{
					newScalarAttribute("mosaicPropertyId", ByteSize),
					newArrayAttribute("value", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func mosaicSupplyChangeTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("mosaicId", IntSize),
			newScalarAttribute("direction", ByteSize),
			newArrayAttribute("delta", IntSize),
		},
	}
}

func transferTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("recipient", ByteSize),
			newScalarAttribute("messageSize", ShortSize),
			newScalarAttribute("numMosaics", ByteSize),
			newTableAttribute("message", schema{
				[]schemaAttribute{
					newScalarAttribute("type", ByteSize),
					newArrayAttribute("payload", ByteSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("mosaics", schema{
				[]schemaAttribute{
					newArrayAttribute("id", IntSize),
					newArrayAttribute("amount", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func modifyMultisigAccountTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("minRemovalDelta", ByteSize),
			newScalarAttribute("minApprovalDelta", ByteSize),
			newScalarAttribute("numModifications", ByteSize),
			newTableArrayAttribute("modification", schema{
				[]schemaAttribute{
					newScalarAttribute("type", ByteSize),
					newArrayAttribute("cosignatoryPublicKey", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func modifyContractTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("durationDelta", IntSize),
			newArrayAttribute("hash", ByteSize),
			newScalarAttribute("numCustomers", ByteSize),
			newScalarAttribute("numExecutors", ByteSize),
			newScalarAttribute("numVerifiers", ByteSize),
			newTableArrayAttribute("customers", schema{
				[]schemaAttribute{
					newScalarAttribute("type", ByteSize),
					newArrayAttribute("cosignatoryPublicKey", ByteSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("executors", schema{
				[]schemaAttribute{
					newScalarAttribute("type", ByteSize),
					newArrayAttribute("cosignatoryPublicKey", ByteSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("verifiers", schema{
				[]schemaAttribute{
					newScalarAttribute("type", ByteSize),
					newArrayAttribute("cosignatoryPublicKey", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func registerNamespaceTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("namespaceType", ByteSize),
			newArrayAttribute("durationParentId", IntSize),
			newArrayAttribute("namespaceId", IntSize),
			newScalarAttribute("namespaceNameSize", ByteSize),
			newArrayAttribute("name", ByteSize),
		},
	}
}

func lockFundsTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("mosaicId", IntSize),
			newArrayAttribute("mosaicAmount", IntSize),
			newArrayAttribute("duration", IntSize),
			newArrayAttribute("hash", ByteSize),
		},
	}
}

func secretLockTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("mosaicId", IntSize),
			newArrayAttribute("mosaicAmount", IntSize),
			newArrayAttribute("duration", IntSize),
			newScalarAttribute("hashAlgorithm", ByteSize),
			newArrayAttribute("secret", ByteSize),
			newArrayAttribute("recipient", ByteSize),
		},
	}
}

func secretProofTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("hashAlgorithm", ByteSize),
			newArrayAttribute("secret", ByteSize),
			newArrayAttribute("recipient", ByteSize),
			newScalarAttribute("proofSize", ShortSize),
			newArrayAttribute("proof", ByteSize),
		},
	}
}
