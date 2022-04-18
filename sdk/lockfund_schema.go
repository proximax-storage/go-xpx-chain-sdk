// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func lockFundTransferTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("duration", IntSize),
			newScalarAttribute("action", ByteSize),
			newScalarAttribute("mosaicsCount", ByteSize),
			newTableArrayAttribute("mosaics", schema{
				[]schemaAttribute{
					newArrayAttribute("id", IntSize),
					newArrayAttribute("amount", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func lockFundCancelUnlockTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("targetHeight", IntSize),
		},
	}
}
