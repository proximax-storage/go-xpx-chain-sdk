package sdk

func storageDriveTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("actionType", IntSize),
			newArrayAttribute("action", ByteSize),
		},
	}
}

func storagePrepareDriveSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newArrayAttribute("duration", IntSize),
			newArrayAttribute("size", IntSize),
			newArrayAttribute("replicas", IntSize),
		},
	}
}

func storageDriveProlongationSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newArrayAttribute("duration", IntSize),
		},
	}
}

func storageFileDriveSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newArrayAttribute("hash", ByteSize),
			newArrayAttribute("parentHash", ByteSize),
			newArrayAttribute("name", ByteSize),
		},
	}
}
func storageFileSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newArrayAttribute("file", ByteSize),
		},
	}
}
func storageDirectorySchema() *schema {
	return &schema{
		[]schemaAttribute{
			newArrayAttribute("directory", ByteSize),
		},
	}
}
func storageFileOperationSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newArrayAttribute("source", ByteSize),
			newArrayAttribute("destination", ByteSize),
		},
	}
}
func storageFileHashSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newArrayAttribute("fileHash", ByteSize),
		},
	}
}
