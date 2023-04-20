// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func placeSdaExchangeOfferTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("sdaOfferCount", ByteSize),
			newTableArrayAttribute("offers", schema{
				[]schemaAttribute{
					newArrayAttribute("mosaicIdGive", IntSize),
					newArrayAttribute("mosaicAmountGive", IntSize),
					newArrayAttribute("mosaicIdGet", IntSize),
					newArrayAttribute("mosaicAmountGet", IntSize),
					newArrayAttribute("duration", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func removeSdaExchangeOfferTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("sdaOfferCount", ByteSize),
			newTableArrayAttribute("offers", schema{
				[]schemaAttribute{
					newArrayAttribute("mosaicIdGive", IntSize),
					newArrayAttribute("mosaicIdGet", IntSize),
				},
			}.schemaDefinition),
		},
	}
}
