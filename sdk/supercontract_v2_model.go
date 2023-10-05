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
	AutomaticExecutionFileName          string
	AutomaticExecutionsFunctionName     string
	AutomaticExecutionsNextBlockToCheck Height
	AutomaticExecutionCallPayment       Amount
	AutomaticDownloadCallPayment        Amount
	AutomatedExecutionsNumber           uint32
	AutomaticExecutionsPrepaidSince     Height
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

type ContractCall struct {
	CallId               *Hash
	Caller               *PublicAccount
	FileName             string
	FunctionName         string
	ActualArguments      []byte
	ExecutionCallPayment Amount
	DownloadCallPayment  Amount
	ServicePayments      []*Mosaic
	BlockHeight          Height
}

func (contractCall *ContractCall) String() string {
	return fmt.Sprintf(
		`
			"CallId": %s,
			"Caller": %s,
			"FileName": %s,
			"FunctionName": %s,
			"ActualArguments": %d,
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
	StartBatchId uint64
	T            []byte
	R            []byte
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
	ExecutorKey        *PublicAccount
	NextBatchToApprove uint64
	PoEx               ProofOfExecution
}

func (executorInfo *ExecutorInfo) String() string {
	return fmt.Sprintf(
		`
			"ExecutorKey": %s,
			"NextBatchToApprove": %d,
			"PoEx": %+v,
		`,
		executorInfo.ExecutorKey,
		executorInfo.NextBatchToApprove,
		executorInfo.PoEx,
	)
}

type CompletedCall struct {
	CallId        *Hash
	Caller        *PublicAccount
	Status        uint16
	ExecutionWork Amount
	DownloadWork  Amount
}

func (completedCall *CompletedCall) String() string {
	return fmt.Sprintf(
		`
			"CallId": %s,
			"Caller": %s,
			"Status": %d,
			"ExecutionWork": %d,
			"DownloadWork": %d,
		`,
		completedCall.CallId,
		completedCall.Caller,
		completedCall.Status,
		completedCall.ExecutionWork,
		completedCall.DownloadWork,
	)
}

type Batch struct {
	BatchId                     uint64
	Success                     bool
	PoExVerificationInformation []byte
	CompletedCalls              []*CompletedCall
}

func (batch *Batch) String() string {
	return fmt.Sprintf(
		`
			"BatchId": %d,
			"Success": %t,
			"PoExVerificationInformation": %v,
			"CompletedCalls": %+v,
		`,
		batch.BatchId,
		batch.Success,
		batch.PoExVerificationInformation,
		batch.CompletedCalls,
	)
}

type SuperContractV2 struct {
	SuperContractKey                *PublicAccount
	DriveKey                        *PublicAccount
	ExecutionPaymentKey             *PublicAccount
	Assignee                        *PublicAccount
	Creator                         *PublicAccount
	DeploymentBaseModificationsInfo *Hash
	AutomaticExecutionsInfo         *AutomaticExecutionsInfo
	RequestedCalls                  []*ContractCall
	ExecutorsInfo                   []*ExecutorInfo
	Batches                         []*Batch
	ReleasedTransactions            []*Hash
}

func (superContractV2 *SuperContractV2) String() string {
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
		superContractV2.SuperContractKey,
		superContractV2.DriveKey,
		superContractV2.ExecutionPaymentKey,
		superContractV2.Assignee,
		superContractV2.Creator,
		superContractV2.DeploymentBaseModificationsInfo,
		superContractV2.AutomaticExecutionsInfo,
		superContractV2.RequestedCalls,
		superContractV2.ExecutorsInfo,
		superContractV2.Batches,
		superContractV2.ReleasedTransactions,
	)
}

type SuperContractsV2Page struct {
	SuperContractsV2 []*SuperContractV2
	Pagination       Pagination
}

type SuperContractsV2PageOptions struct {
	SuperContractsV2PageFilters
	PaginationOrderingOptions
}

type SuperContractsV2PageFilters struct {
	DriveKey string `url:"owner,omitempty"`

	Creator string `url:"owner,omitempty"`
}

// end of supercontract entry

// Automatic Executions Payment Transaction
type AutomaticExecutionsPaymentTransaction struct {
	AbstractTransaction
	ContractKey               *PublicAccount
	AutomaticExecutionsNumber uint32
}

// Manual Call Transaction
type ManualCallTransaction struct {
	AbstractTransaction
	ContractKey          *PublicAccount
	ExecutionCallPayment Amount
	DownloadCallPayment  Amount
	FileName             string
	FunctionName         string
	ActualArguments      []byte
	ServicePayments      []*Mosaic
}

// Deploy Contract Transaction
type DeployContractTransaction struct {
	AbstractTransaction
	DriveKey                       *PublicAccount
	ExecutionCallPayment           Amount
	DownloadCallPayment            Amount
	AutomaticExecutionCallPayment  Amount
	AutomaticDownloadCallPayment   Amount
	AutomaticExecutionsNumber      uint32
	Assignee                       *PublicAccount
	FileName                       string
	FunctionName                   string
	ActualArguments                []byte
	ServicePayments                []*Mosaic
	AutomaticExecutionFileName     string
	AutomaticExecutionFunctionName string
}

// Successful End Batch Execution Transaction
type RawProofsOfExecution struct {
	StartBatchId uint64
	T            []byte
	R            []byte
	F            []byte
	K            []byte
}

func (rawPoex *RawProofsOfExecution) String() string {
	return fmt.Sprintf(
		`
		"StartBatchId": %d,
		"T": %v,
		"R": %v,
		"F": %v,
			"K": %v,
			`,
		rawPoex.StartBatchId,
		rawPoex.T,
		rawPoex.R,
		rawPoex.F,
		rawPoex.K,
	)
}

type ExtendedCallDigest struct {
	CallId                  *Hash
	Manual                  bool
	Block                   Height
	Status                  uint16
	ReleasedTransactionHash *Hash
}

func (extendedCallDigest *ExtendedCallDigest) String() string {
	return fmt.Sprintf(
		`
			"CallId": %s,
			"Manual": %t,
			"Block": %d,
			"Status": %d,
			"ReleasedTransactionHash": %s,
		`,
		extendedCallDigest.CallId,
		extendedCallDigest.Manual,
		extendedCallDigest.Block,
		extendedCallDigest.Status,
		extendedCallDigest.ReleasedTransactionHash,
	)
}

type CallPayment struct {
	ExecutionPayment Amount
	DownloadPayment  Amount
}

func (callPayment *CallPayment) String() string {
	return fmt.Sprintf(
		`
			"ExecutionPayment": %d,
			"DownloadPayment": %d,
		`,
		callPayment.ExecutionPayment,
		callPayment.DownloadPayment,
	)
}

type EndBatchExecution struct {
	ContractKey                         *PublicAccount
	BatchId                             uint64
	AutomaticExecutionsNextBlockToCheck Height
	PublicKeys                          []*PublicAccount
	Signatures                          []*Signature
	ProofsOfExecutions                  []*RawProofsOfExecution
	CallPayments                        []*CallPayment
}

type SuccessfulEndBatchExecutionTransaction struct {
	AbstractTransaction
	EndBatchExecutionInfo                   EndBatchExecution
	StorageHash                             *Hash
	UsedSizedBytes                          uint64
	MetaFilesSizeBytes                      uint64
	ProofOfExecutionVerificationInformation []byte
	CallDigests                             []*ExtendedCallDigest
}

// Unsuccessful End Batch Execution Transaction
type ShortCallDigest struct {
	CallId *Hash
	Manual bool
	Block  Height
}

func (shortCallDigest *ShortCallDigest) String() string {
	return fmt.Sprintf(
		`
			"CallId": %s,
			"Manual": %t,
			"Block": %d,
		`,
		shortCallDigest.CallId,
		shortCallDigest.Manual,
		shortCallDigest.Block,
	)
}

type UnsuccessfulEndBatchExecutionTransaction struct {
	AbstractTransaction
	EndBatchExecution
	CallDigests []*ShortCallDigest
}
