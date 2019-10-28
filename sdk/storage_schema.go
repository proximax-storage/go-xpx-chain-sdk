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
			newScalarAttribute("duration", IntSize),
			newScalarAttribute("billingPeriod", IntSize),
			newScalarAttribute("billingPrice", IntSize),
			newScalarAttribute("driveSize", ByteSize),
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
			newScalarAttribute("addActionsCount", ShortSize),
			newScalarAttribute("removeActionsCount", ShortSize),
			newTableArrayAttribute("addActions", schema{
				[]schemaAttribute{
					newArrayAttribute("fileHash", ByteSize),
					newArrayAttribute("fileSize", IntSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("removeActions", schema{
				[]schemaAttribute{
					newArrayAttribute("fileHash", ByteSize),
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
