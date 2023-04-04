// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

func automaticExecutionsPaymentTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("contractKey", ByteSize),
			newArrayAttribute("automaticExecutionsNumber", IntSize),
		},
	}
}

func manualCallTransactionSchema() *schema {
	return &schema{
		[]schemaAttribute{
			newScalarAttribute("size", IntSize),
			newArrayAttribute("signature", ByteSize),
			newArrayAttribute("signer", ByteSize),
			newScalarAttribute("version", IntSize),
			newScalarAttribute("type", ShortSize),
			newArrayAttribute("maxFee", IntSize),
			newArrayAttribute("deadline", IntSize),
			newArrayAttribute("contractKey", ByteSize),
			newScalarAttribute("fileNameSize", ShortSize),
			newScalarAttribute("functionNameSize", ShortSize),
			newScalarAttribute("actualArgumentsSize", ShortSize),
			newArrayAttribute("executionCallPayment", IntSize),
			newArrayAttribute("downloadCallPayment", IntSize),
			newScalarAttribute("servicePaymentsCount", ByteSize),
			newArrayAttribute("fileName", ByteSize),
			newArrayAttribute("functionName", ByteSize),
			newArrayAttribute("actualArguments", ByteSize),
			newArrayAttribute("servicePayments", IntSize),
		},
	}
}

func deployContractTransactionSchema() *schema {
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
			newScalarAttribute("fileNameSize", ShortSize),
			newScalarAttribute("functionNameSize", ShortSize),
			newScalarAttribute("actualArgumentsSize", ShortSize),
			newArrayAttribute("executionCallPayment", IntSize),
			newArrayAttribute("downloadCallPayment", IntSize),
			newScalarAttribute("servicePaymentsCount", ByteSize),
			newScalarAttribute("automaticExecutionFileNameSize", ShortSize),
			newScalarAttribute("automaticExecutionFunctionNameSize", ShortSize),
			newArrayAttribute("automaticExecutionCallPayment", IntSize),
			newArrayAttribute("automaticDownloadCallPayment", IntSize),
			newScalarAttribute("automaticExecutionsNumber", IntSize),
			newArrayAttribute("assignee", ByteSize),
			newArrayAttribute("fileName", ByteSize),
			newArrayAttribute("functionName", ByteSize),
			newArrayAttribute("actualArguments", ByteSize),
			newArrayAttribute("servicePayments", IntSize),
			newArrayAttribute("automaticExecutionFileName", ByteSize),
			newArrayAttribute("automaticExecutionFunctionName", ByteSize),
		},
	}
}
