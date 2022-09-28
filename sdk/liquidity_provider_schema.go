// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func createLiquidityProviderTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("providerMosaicId", IntSize),
			newArrayAttribute("currencyDeposit", IntSize),
			newArrayAttribute("initialMosaicsMinting", IntSize),
			newScalarAttribute("slashingPeriod", IntSize),
			newScalarAttribute("windowSize", ShortSize),
			newArrayAttribute("slashingAccount", ByteSize),
			newScalarAttribute("alpha", IntSize),
			newScalarAttribute("beta", IntSize),
		},
	}
}

func manualRateChangeTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("providerMosaicId", IntSize),
			newScalarAttribute("currencyBalanceIncrease", ByteSize),
			newArrayAttribute("currencyBalanceChange", IntSize),
			newScalarAttribute("mosaicBalanceIncrease", ByteSize),
			newArrayAttribute("mosaicBalanceChange", IntSize),
		},
	}
}
