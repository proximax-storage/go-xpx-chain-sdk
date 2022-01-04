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
			newArrayAttribute("drive", ByteSize),
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
			newArrayAttribute("downloadSize", IntSize),
			newArrayAttribute("feedbackFeeAmount", IntSize),
			newScalarAttribute("publicKeyCount", ShortSize),
			newTableArrayAttribute("listOfPublicKeys", schema{
				[]schemaAttribute{
					newArrayAttribute("hashes", ByteSize),
				},
			}.schemaDefinition),
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
			newScalarAttribute("sequenceNumber", ShortSize),
			newScalarAttribute("responseToFinishDownloadTransaction", ByteSize),
			newScalarAttribute("judgingKeysCount", ByteSize),
			newScalarAttribute("overlappingKeysCount", ByteSize),
			newScalarAttribute("judgedKeysCount", ByteSize),
			newScalarAttribute("opinionElementCount", ByteSize),
			newTableArrayAttribute("publicKeys", schema{
				[]schemaAttribute{
					newArrayAttribute("hashes", ByteSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("signatures", schema{
				[]schemaAttribute{
					newArrayAttribute("hashes", ByteSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("presentOpinions", schema{
				[]schemaAttribute{
					newArrayAttribute("present", ByteSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("opinions", schema{
				[]schemaAttribute{
					newArrayAttribute("opinion", ByteSize),
				},
			}.schemaDefinition),
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
			newArrayAttribute("downloadChannelId", ByteSize),
			newArrayAttribute("downloadSize", IntSize),
			newArrayAttribute("feedbackFeeAmount", IntSize),
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
