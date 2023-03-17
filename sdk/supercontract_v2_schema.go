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
			newScalarAttribute("automaticExecutionsNumber", IntSize),
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
			newScalarAttribute("FileNameSize", ShortSize),
			newScalarAttribute("FunctionNameSize", ShortSize),
			newScalarAttribute("ActualArgumentsSize", ShortSize),
			newArrayAttribute("ExecutionCallPayment", IntSize),
			newArrayAttribute("DownloadCallPayment", IntSize),
			newScalarAttribute("ServicePaymentsCount", ByteSize),
			newArrayAttribute("FileName", IntSize),
			newArrayAttribute("FunctionName", IntSize),
			newArrayAttribute("ActualArguments", IntSize),
			newArrayAttribute("ServicePayments", IntSize),
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
			newScalarAttribute("FileNameSize", ShortSize),
			newScalarAttribute("FunctionNameSize", ShortSize),
			newScalarAttribute("ActualArgumentsSize", ShortSize),
			newArrayAttribute("ExecutionCallPayment", IntSize),
			newArrayAttribute("DownloadCallPayment", IntSize),
			newScalarAttribute("ServicePaymentsCount", ByteSize),
			newScalarAttribute("AutomaticExecutionFileNameSize", ShortSize),
			newScalarAttribute("AutomaticExecutionFunctionNameSize", ShortSize),
			newArrayAttribute("AutomaticExecutionCallPayment", IntSize),
			newArrayAttribute("AutomaticDownloadCallPayment", IntSize),
			newScalarAttribute("AutomaticExecutionsNumber", IntSize),
			newArrayAttribute("Assignee", ByteSize),
			newArrayAttribute("FileName", IntSize),
			newArrayAttribute("FunctionName", IntSize),
			newArrayAttribute("ActualArguments", IntSize),
			newArrayAttribute("ServicePayments", IntSize),
			newArrayAttribute("AutomaticExecutionFileName", IntSize),
			newArrayAttribute("AutomaticExecutionFunctionName", IntSize),
		},
	}
}