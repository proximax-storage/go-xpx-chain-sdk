// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func addExchangeOfferTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("offersCount", ByteSize),
			newTableArrayAttribute("offers", schema{
				[]schemaAttribute{
					newArrayAttribute("mosaicId", IntSize),
					newArrayAttribute("mosaicAmount", IntSize),
					newArrayAttribute("cost", IntSize),
					newScalarAttribute("type", ByteSize),
					newArrayAttribute("duration", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func exchangeOfferTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("offersCount", ByteSize),
			newTableArrayAttribute("offers", schema{
				[]schemaAttribute{
					newArrayAttribute("mosaicId", IntSize),
					newArrayAttribute("mosaicAmount", IntSize),
					newArrayAttribute("cost", IntSize),
					newScalarAttribute("type", ByteSize),
					newArrayAttribute("owner", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func removeExchangeOfferTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("offersCount", ByteSize),
			newTableArrayAttribute("offers", schema{
				[]schemaAttribute{
					newArrayAttribute("mosaicId", IntSize),
					newScalarAttribute("type", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}
