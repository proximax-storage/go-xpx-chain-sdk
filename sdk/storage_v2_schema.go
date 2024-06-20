// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func replicatorOnboardingTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("capacity", IntSize),
			newArrayAttribute("nodeBootKey", ByteSize),
			newArrayAttribute("message", ByteSize),
			newArrayAttribute("messageSignature", ByteSize),
		},
	}
}

func prepareBcDriveTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveSize", IntSize),
			newArrayAttribute("verificationFeeAmount", IntSize),
			newScalarAttribute("replicatorCount", ShortSize),
		},
	}
}

func dataModificationTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveKey", ByteSize),
			newArrayAttribute("downloadDataCdi", ByteSize),
			newArrayAttribute("uploadSize", IntSize),
			newArrayAttribute("feedbackFeeAmount", IntSize),
		},
	}
}

func dataModificationCancelTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveKey", ByteSize),
			newArrayAttribute("downloadDataCdi", ByteSize),
		},
	}
}

func storagePaymentTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveKey", ByteSize),
			newArrayAttribute("storageUnits", IntSize),
		},
	}
}

func downloadPaymentTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveKey", ByteSize),
			newArrayAttribute("downloadSize", IntSize),
			newArrayAttribute("feedbackFeeAmount", IntSize),
		},
	}
}

func downloadTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveKey", ByteSize),
			newArrayAttribute("downloadSize", IntSize),
			newArrayAttribute("feedbackFeeAmount", IntSize),
			newArrayAttribute("listOfPublicKeysSize", ByteSize),
			newTableArrayAttribute("listOfPublicKeys", schema{
				[]schemaAttribute{
					newArrayAttribute("key", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func finishDownloadTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("downloadChannelId", ByteSize),
			newArrayAttribute("feedbackFeeAmount", IntSize),
		},
	}
}

func verificationPaymentTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveKey", ByteSize),
			newArrayAttribute("VerificationFeeAmount", IntSize),
		},
	}
}

func endDriveVerificationV2TransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveKey", ByteSize),
			newArrayAttribute("verificationTrigger", ByteSize),
			newArrayAttribute("shardId", ByteSize),
			newScalarAttribute("keyCount", ByteSize),
			newScalarAttribute("judgingKeyCount", ByteSize),
			newTableArrayAttribute("keys", schema{
				[]schemaAttribute{
					newArrayAttribute("key", ByteSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("signatures", schema{
				[]schemaAttribute{
					newArrayAttribute("signature", ByteSize),
				},
			}.schemaDefinition),
			newScalarAttribute("opinions", ByteSize),
		},
	}
}

func driveClosureTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveKey", ByteSize),
		},
	}
}

func replicatorOffboardingTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("driveKey", ByteSize),
		},
	}
}

func downloadApprovalTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("downloadChannelId", ByteSize),
			newArrayAttribute("approvalTrigger", ByteSize),
			newScalarAttribute("sequenceNumber", ByteSize),
			newScalarAttribute("responseToFinishDownloadTransaction", ByteSize),
			newScalarAttribute("judgingCount", ByteSize),
			newScalarAttribute("overlappingCount", ByteSize),
			newScalarAttribute("judgedCount", ByteSize),
			newScalarAttribute("opinionElementCount", ByteSize),
			newTableArrayAttribute("publicKeys", schema{
				[]schemaAttribute{
					newArrayAttribute("key", ByteSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("signatures", schema{
				[]schemaAttribute{
					newArrayAttribute("signature", ByteSize),
				},
			}.schemaDefinition),
			newArrayAttribute("presentOpinions", ByteSize),
			newTableArrayAttribute("opinions", schema{
				[]schemaAttribute{
					newArrayAttribute("Opinion", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func replicatorsCleanupTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("replicatorCount", ByteSize),
			newTableArrayAttribute("replicatorKeys", schema{
				[]schemaAttribute{
					newArrayAttribute("key", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func replicatorTreeRebuildTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("replicatorCount", ByteSize),
			newTableArrayAttribute("replicatorKeys", schema{
				[]schemaAttribute{
					newArrayAttribute("key", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}
