// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func prepareDriveTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("owner", ByteSize),
			newArrayAttribute("duration", IntSize),
			newArrayAttribute("billingPeriod", IntSize),
			newArrayAttribute("billingPrice", IntSize),
			newArrayAttribute("driveSize", IntSize),
			newScalarAttribute("replicas", ShortSize),
			newScalarAttribute("minReplicators", ShortSize),
			newScalarAttribute("percentApprovers", ByteSize),
		},
	}
}

func driveFileSystemTransactionSchema() *schema {
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
			newArrayAttribute("rootHash", ByteSize),
			newArrayAttribute("xorRootHash", ByteSize),
			newArrayAttribute("addActionsCount", ByteSize),
			newArrayAttribute("removeActionsCount", ByteSize),
			newTableArrayAttribute("addActions", schema{
				[]schemaAttribute{
					newArrayAttribute("fileHash", ByteSize),
					newArrayAttribute("fileSize", IntSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("removeActions", schema{
				[]schemaAttribute{
					newArrayAttribute("fileHash", ByteSize),
					newArrayAttribute("fileSize", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func joinDriveTransactionSchema() *schema {
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

func filesDepositTransactionSchema() *schema {
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
			newScalarAttribute("filesCount", ShortSize),
			newTableArrayAttribute("files", schema{
				[]schemaAttribute{
					newArrayAttribute("fileHash", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func endDriveTransactionSchema() *schema {
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

func startDriveVerificationTransactionSchema() *schema {
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

func endDriveVerificationTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newTableArrayAttribute("failures", schema{
				[]schemaAttribute{
					newScalarAttribute("size", IntSize),
					newArrayAttribute("replicator", ByteSize),
					newTableArrayAttribute("blockHashes", schema{
						[]schemaAttribute{
							newArrayAttribute("blockHash", ByteSize),
						},
					}.schemaDefinition),
				},
			}.schemaDefinition),
		},
	}
}

func startFileDownloadTransactionSchema() *schema {
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
			newScalarAttribute("filesCount", ShortSize),
			newTableArrayAttribute("files", schema{
				[]schemaAttribute{
					newArrayAttribute("hash", ByteSize),
					newArrayAttribute("size", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func endFileDownloadTransactionSchema() *schema {
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
			newArrayAttribute("operationHash", ByteSize),
			newScalarAttribute("filesCount", ShortSize),
			newTableArrayAttribute("files", schema{
				[]schemaAttribute{
					newArrayAttribute("hash", ByteSize),
					newArrayAttribute("size", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func driveFilesRewardTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("uploadInfosCount", ShortSize),
			newTableArrayAttribute("uploadInfos", schema{
				[]schemaAttribute{
					newArrayAttribute("replicator", ByteSize),
					newArrayAttribute("uploaded", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func deployTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("drive", ByteSize),
			newArrayAttribute("owner", ByteSize),
			newArrayAttribute("fileHash", ByteSize),
			newArrayAttribute("vmVersion", IntSize),
		},
	}
}

func startExecuteTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("superContract", ByteSize),
			newScalarAttribute("functionSize", ByteSize),
			newScalarAttribute("mosaicsCount", ByteSize),
			newArrayAttribute("dataSize", ByteSize),
			newArrayAttribute("function", ByteSize),
			newTableArrayAttribute("mosaics", schema{
				[]schemaAttribute{
					newArrayAttribute("id", IntSize),
					newArrayAttribute("amount", IntSize),
				},
			}.schemaDefinition),
			newArrayAttribute("data", ByteSize),
		},
	}
}

func operationIdentifyTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("operationToken", ByteSize),
		},
	}
}

func endOperationTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("mosaicsCount", ByteSize),
			newArrayAttribute("operationToken", ByteSize),
			newScalarAttribute("status", ShortSize),
			newTableArrayAttribute("mosaics", schema{
				[]schemaAttribute{
					newArrayAttribute("id", IntSize),
					newArrayAttribute("amount", IntSize),
				},
			}.schemaDefinition),
		},
	}
}
