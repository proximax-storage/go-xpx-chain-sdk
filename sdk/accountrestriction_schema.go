// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func accountAddressRestrictionTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("restrictionFlags", ShortSize),
			newScalarAttribute("restrictionAdditionsCount", ByteSize),
			newScalarAttribute("restrictionDeletionsCount", ByteSize),
			newScalarAttribute("AccountRestrictionTransactionBody_Reserved1", IntSize),
			newTableArrayAttribute("restrictions", schema{
				[]schemaAttribute{
					newArrayAttribute("address", ByteSize),
				},
			}.schemaDefinition),
		},
	}
}

func accountMosaicRestrictionTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("restrictionFlags", ShortSize),
			newScalarAttribute("RestrictionAdditionsCount", ByteSize),
			newScalarAttribute("RestrictionDeletionsCount", ByteSize),
			newScalarAttribute("AccountRestrictionTransactionBody_Reserved1", IntSize),
			newTableArrayAttribute("restrictions", schema{
				[]schemaAttribute{
					newArrayAttribute("mosaicId", IntSize),
				},
			}.schemaDefinition),
		},
	}
}

func accountOperationRestrictionTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newScalarAttribute("restrictionFlags", ShortSize),
			newScalarAttribute("RestrictionAdditionsCount", ByteSize),
			newScalarAttribute("RestrictionDeletionsCount", ByteSize),
			newScalarAttribute("AccountRestrictionTransactionBody_Reserved1", IntSize),
			newArrayAttribute("restrictions", ShortSize),
		},
	}
}
