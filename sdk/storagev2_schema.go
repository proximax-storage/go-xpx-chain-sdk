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
			newArrayAttribute("blsPublicKey", ByteSize),
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
			newArrayAttribute("driveKey", ByteSize),
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
			newScalarAttribute("proversCount", ShortSize),
			newScalarAttribute("verificationOpinionsCount", ShortSize),
			newTableArrayAttribute("provers", schema{
				[]schemaAttribute{
					newArrayAttribute("prover", ByteSize),
				},
			}.schemaDefinition),
			newTableArrayAttribute("verificationOpinions", schema{
				[]schemaAttribute{
					newArrayAttribute("verifier", ByteSize),
					newArrayAttribute("blsSignature", ByteSize),
					newTableArrayAttribute("results", schema{
						[]schemaAttribute{
							newArrayAttribute("prover", ByteSize),
							newScalarAttribute("result", ByteSize),
						},
					}.schemaDefinition),
				},
			}.schemaDefinition),
		},
	}
}
