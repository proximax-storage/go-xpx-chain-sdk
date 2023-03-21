// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import "fmt"


// drive contract entry (will be implement in the future)
// type DriveContract struct {
// 	DriveContractKey 				*PublicAccount
// 	ContractKey 					*PublicAccount
// }

// func (driveContract *DriveContract) String() string {
// 	return fmt.Sprintf(
// 		`
// 			"DriveContractKey": %s,
// 			"ContractKey": %s,
// 		`,
// 		driveContract.DriveContractKey.String(),
// 		driveContract.ContractKey.String(),
// 	)
// }

// type DriveContractPage struct {
// 	DriveContracts 	[]*DriveContract
// 	Pagination 		Pagination
// }

// type DriveContractFileters struct {
// 	ContractKey string `url:"consumerKey, omitempty"`
// }

// type DriveContractPageOptions struct {
// 	BcDrivesPageFilters
// 	PaginationOrderingOptions
// }
// end of drive contract entry

// supercontract entry
type AutomaticExecutionsInfo struct {
	AutomaticExecutionFileName 			string
	AutomaticExecutionsFunctionName 	string
	AutomaticExecutionsNextBlockToCheck Height
	AutomaticExecutionCallPayment 		Amount
	AutomaticDownloadCallPayment 		Amount
	AutomatedExecutionsNumber 			uint32 `default:"0"`
	AutomaticExecutionsPrepaidSince 	Height
}

func (executionsInfo *AutomaticExecutionsInfo) String() string {
	return fmt.Sprintf(
		`
			"AutomaticExecutionFileName": %s,
			"AutomaticExecutionsFunctionName": %s,
			"AutomaticExecutionsNextBlockToCheck": %d,
			"AutomaticExecutionCallPayment": %d,
			"AutomaticDownloadCallPayment": %d,
			"AutomatedExecutionsNumber": %d,
			"AutomaticExecutionsPrepaidSince": %d,
		`,
		executionsInfo.AutomaticExecutionFileName,
		executionsInfo.AutomaticExecutionsFunctionName,
		executionsInfo.AutomaticExecutionsNextBlockToCheck,
		executionsInfo.AutomaticExecutionCallPayment,
		executionsInfo.AutomaticDownloadCallPayment,
		executionsInfo.AutomatedExecutionsNumber,
		executionsInfo.AutomaticExecutionsPrepaidSince,
	)
}

type ServicePayment struct {
	MosaicId 	*Mosaic
	Amount 		Amount
}

func (servicePayment *ServicePayment) String() string {
	return fmt.Sprintf(
		`
			"MosaicId": %d,
			"Amount": %d,
		`,
		servicePayment.MosaicId,
		servicePayment.Amount,
	)
}

type ContractCall struct {
	CallId 								*Hash
	Caller 								*PublicAccount
	FileName 							string
	FunctionName 						string
	ActualArguments 					string
	ExecutionCallPayment 				Amount
	DownloadCallPayment	 				Amount
	ServicePayments 					[]*ServicePayment
	BlockHeight 						Height
}

func (contractCall *ContractCall) String() string {
	return fmt.Sprintf(
		`
			"CallId": %s,
			"Caller": %s,
			"FileName": %s,
			"FunctionName": %s,
			"ActualArguments": %s,
			"ExecutionCallPayment": %d,
			"DownloadCallPayment": %d,
			"ServicePayments": %+v,
			"BlockHeight": %d,
		`,
		contractCall.CallId,
		contractCall.Caller,
		contractCall.FileName,
		contractCall.FunctionName,
		contractCall.ActualArguments,
		contractCall.ExecutionCallPayment,
		contractCall.DownloadCallPayment,
		contractCall.ServicePayments,
		contractCall.BlockHeight,
	)
}

type ProofOfExecution struct {
	StartBatchId 						uint64 `default:"0"`
	T  									[]byte
	R  									[]byte
}	

func (proofOfExecution *ProofOfExecution) String() string {
	return fmt.Sprintf(
		`
			"StartBatchId": %d,
			"T": %v,
			"R": %v,
		`,
		proofOfExecution.StartBatchId,
		proofOfExecution.T,
		proofOfExecution.R,
	)
}

type ExecutorInfo struct {
	NextBatchToApproave 				uint64 `default:"0"`
	PoEx 								ProofOfExecution
}

func (executorInfo *ExecutorInfo) String() string {
	return fmt.Sprintf(
		`
			"NextBatchToApproave": %d,
			"PoEx": %+v,
		`,
		executorInfo.NextBatchToApproave,
		executorInfo.PoEx,
	)
}

type Batch struct {
	Success bool
	PoExVerificationInformation			[]byte

}

func (batch *Batch) String() string {
	return fmt.Sprintf(
		`
			"Success": %t,
			"PoExVerificationInformation": %v,
		`,
		batch.Success,
		batch.PoExVerificationInformation,
	)
}

type SuperContractEntry struct {
	SuperContractKey 					*PublicAccount
	DriveKey 							*PublicAccount
	ExecutionPaymentKey					*PublicAccount
	Assignee 							*PublicAccount
	Creator 							*PublicAccount
	DeploymentBaseModificationsInfo 	*Hash
	AutomaticExecutionsInfo 			AutomaticExecutionsInfo
	RequestedCalls 						[]*ContractCall
	ExecutorsInfo 						map[PublicAccount]ExecutorInfo
	Batches 							map[uint64]Batch
	ReleasedTransactions 				[]*Hash
}

func (superContractEntry *SuperContractEntry) String() string {
	return fmt.Sprintf(
		`
			"SuperContractKey": %s,
			"DriveKey": %s,
			"ExecutionPaymentKey": %s,
			"Assignee": %s,
			"Creator": %s,
			"DeploymentBaseModificationsInfo": %s,
			"AutomaticExecutionsInfo": %+v,
			"RequestedCalls": %+v,
			"ExecutorsInfo": %+v,
			"Batches": %+v,
			"ReleasedTransactions": %v,
		`,
		superContractEntry.SuperContractKey.String(),
		superContractEntry.DriveKey.String(),
		superContractEntry.ExecutionPaymentKey.String(),
		superContractEntry.Assignee.String(),
		superContractEntry.Creator.String(),
		superContractEntry.DeploymentBaseModificationsInfo.String(),
		superContractEntry.AutomaticExecutionsInfo,
		superContractEntry.RequestedCalls,
		superContractEntry.ExecutorsInfo,
		superContractEntry.Batches,
		superContractEntry.ReleasedTransactions,
	)
}

type SuperContractEntriesPage struct {
	SuperContractEntries []*SuperContractEntry
	Pagination Pagination
}

type SuperContractEntriesPageOption struct {
	SuperContractEntriesPageFilters
	PaginationOrderingOptions
}

type SuperContractEntriesPageFilters struct {
	DriveKey string `url:"owner,omitempty"`

	Creator string `url:"owner,omitempty"`     
}
// end of supercontract entry

// Automatic Executions Payment Transaction
type AutomaticExecutionsPaymentTransaction struct {
	AbstractTransaction
	ContractKey							*PublicAccount
	AutomaticExecutionsNumber 			uint32
}

// Manual Call Transaction
type ManualCallTransaction struct {
	AbstractTransaction
	ContractKey							*PublicAccount
	FileNameSize						uint16
	FunctionNameSize					uint16
	ActualArgumentsSize					uint16
	ExecutionCallPayment				Amount
	DownloadCallPayment					Amount
	ServicePaymentsCount 				uint8
	FileName							string
	FunctionName						string
	ActualArguments						string
	ServicePayments						[]*Mosaic
}

// Deploy Contract Transaction
type DeployContractTransaction struct {
	AbstractTransaction
	DriveKey							*PublicAccount
	FileNameSize						uint16
	FunctionNameSize					uint16
	ActualArgumentsSize					uint16
	ExecutionCallPayment				Amount
	DownloadCallPayment					Amount
	ServicePaymentsCount 				uint8
	AutomaticExecutionFileNameSize		uint16
	AutomaticExecutionFunctionNameSize 	uint16
	AutomaticExecutionCallPayment 		Amount
	AutomaticDownloadCallPayment 		Amount
	AutomaticExecutionsNumber 			uint32
	Assignee							*PublicAccount
	FileName							string
	FunctionName						string
	ActualArguments						string
	ServicePayments						[]*Mosaic
	AutomaticExecutionFileName			string
	AutomaticExecutionFunctionName		string
}