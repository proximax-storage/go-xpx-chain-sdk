// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func mosaicGlobalRestrictionTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("mosaicId", IntSize),
			newArrayAttribute("referenceMosaicId", IntSize),
			newArrayAttribute("restrictionKey", IntSize),
			newArrayAttribute("previousRestrictionValue", IntSize),
			newArrayAttribute("newRestrictionValue", IntSize),
			newScalarAttribute("previousRestrictionType", ByteSize),
			newScalarAttribute("newRestrictionType", ByteSize),
		},
	}
}

func mosaicAddressRestrictionTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("mosaicId", IntSize),
			newArrayAttribute("restrictionKey", IntSize),
			newArrayAttribute("previousRestrictionValue", IntSize),
			newArrayAttribute("newRestrictionValue", IntSize),
			newArrayAttribute("targetAddress", ByteSize),
		},
	}
}
