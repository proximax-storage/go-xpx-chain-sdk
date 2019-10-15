package sdk

func modifyDriveTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("priceDelta", IntSize),
			newScalarAttribute("durationDelta", IntSize),
			newScalarAttribute("sizeDelta", ByteSize),
			newScalarAttribute("replicasDelta", IntSize),
			newScalarAttribute("minReplicatorsDelta", IntSize),
			newScalarAttribute("minApproversDelta", IntSize),
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
			newScalarAttribute("rootHash", ByteSize),
			newScalarAttribute("xorRootHash", ByteSize),
			newScalarAttribute("addActionsCount", ByteSize),
			newScalarAttribute("removeActionsCount", ByteSize),
			newTableArrayAttribute("addActions", schema{
				[]schemaAttribute{
					newArrayAttribute("fileHash", ByteSize),
					newScalarAttribute("fileSize", IntSize),
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
			newScalarAttribute("driveKey", ByteSize),
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
			newScalarAttribute("driveKey", ByteSize),
			newScalarAttribute("filesCount", ByteSize),
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
		},
	}
}
