// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func streamStartTransactionSchema() *schema {
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
			newArrayAttribute("expectedUploadSize", IntSize),
			newScalarAttribute("folderNameSize", ShortSize),
			newArrayAttribute("feedbackFeeAmount", IntSize),
			newArrayAttribute("folderName", ByteSize),
		},
	}
}

func streamFinishTransactionSchema() *schema {
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
			newArrayAttribute("streamId", ByteSize),
			newArrayAttribute("actualUploadSize", IntSize),
			newArrayAttribute("streamStructureCdi", ByteSize),
		},
	}
}

func streamPaymentTransactionSchema() *schema {
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
			newArrayAttribute("streamId", ByteSize),
			newArrayAttribute("additionalUploadSize", IntSize),
		},
	}
}
