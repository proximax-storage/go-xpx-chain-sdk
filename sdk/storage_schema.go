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
