// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
)

type automaticExecutionsInfoDTO struct {
	AutomaticExecutionFileName 			string		`json:"automaticExecutionFileName"`
	AutomaticExecutionsFunctionName 	string		`json:"automaticExecutionsFunctionName"`
	AutomaticExecutionsNextBlockToCheck uint64DTO	`json:"automaticExecutionsNextBlockToCheck"`
	AutomaticExecutionCallPayment 		uint64DTO	`json:"automaticExecutionCallPayment"`
	AutomaticDownloadCallPayment 		uint64DTO	`json:"automaticDownloadCallPayment"`
	AutomatedExecutionsNumber 			uint32		`json:"automatedExecutionsNumber"`
	AutomaticExecutionsPrepaidSince 	uint64DTO	`json:"automaticExecutionsPrepaidSince"`
}

func (ref *automaticExecutionsInfoDTO) toStruct(networkType NetworkType) (*AutomaticExecutionsInfo, error) {

	return &AutomaticExecutionsInfo{
		AutomaticExecutionFileName: 			ref.AutomaticExecutionFileName,
		AutomaticExecutionsFunctionName:		ref.AutomaticExecutionsFunctionName,
		AutomaticExecutionsNextBlockToCheck:	ref.AutomaticExecutionsNextBlockToCheck.toStruct(),
		AutomaticExecutionCallPayment:			ref.AutomaticExecutionCallPayment.toStruct(),
		AutomaticDownloadCallPayment:			ref.AutomaticDownloadCallPayment.toStruct(),
		AutomatedExecutionsNumber:				ref.AutomatedExecutionsNumber,
		AutomaticExecutionsPrepaidSince:		ref.AutomaticExecutionsPrepaidSince.toStruct(),
	}, nil
}

type servicePaymentDTO struct {
	MosaicId 	uint64DTO 	`json:"mosaicId"`
	Amount 		uint64DTO 	`json:"amount"`	
}

func (ref *servicePaymentDTO) toStruct(NetworkType NetworkType) (*ServicePayment, error) {

	mosaicId, err := NewMosaicId(ref.MosaicId.toUint64())
	if err != nil {
		return nil, err
	}

	return &ServicePayment{
		MosaicId: mosaicId,
		Amount: ref.Amount.toStruct(),
	}, nil
}

type servicePaymentDtos []*servicePaymentDTO

func (ref *servicePaymentDtos) toStruct(networkType NetworkType) ([]*ServicePayment, error) {
	var (
		dtos                    = *ref
		servicePayments = make([]*ServicePayment, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		servicePayments = append(servicePayments, info)
	}

	return servicePayments, nil
}

type contractCallDTO struct {
	CallId					hashDto 	`json:"callId"`
	Caller 					string		`json:"caller"`
	FileName				string 		`json:"fileName"`
	FunctionName			string 		`json:"functionName"`
	ActualArguments			string 		`json:"actualArguments"`
	ExecutionCallPayment 	uint64DTO 	`json:"executionCallPayment"`
	DownloadCallPayment 	uint64DTO 	`json:"downloadCallPayment"`
	servicePaymentDtos 
	BlockHeight				uint64DTO	`json:"blockHeight"`
}

func (ref *contractCallDTO) toStruct(networkType NetworkType) (*ContractCall, error) {
	servicePayments, err := ref.servicePaymentDtos.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	callId, err := ref.CallId.Hash()
	if err != nil {
		return nil, err
	}

	caller, err := NewAccountFromPublicKey(ref.Caller, networkType)
	if err != nil {
		return nil, err
	}

	return &ContractCall{
		CallId: callId,
		Caller: caller,
		FileName: ref.FileName,
		FunctionName: ref.FunctionName,
		ActualArguments: ref.ActualArguments,
		ExecutionCallPayment: ref.ExecutionCallPayment.toStruct(),
		DownloadCallPayment: ref.DownloadCallPayment.toStruct(),
		ServicePayments: servicePayments,
		BlockHeight: ref.BlockHeight.toStruct(),
	}, nil
}

type contractCallDTOs []*contractCallDTO

func (ref *contractCallDTOs) toStruct(networkType NetworkType) ([]*ContractCall, error) {
	var (
		dtos                    = *ref
		contractCalls = make([]*ContractCall, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		contractCalls = append(contractCalls, info)
	}

	return contractCalls, nil
}

type proofOfExecutionDTO struct {
	StartBatchId	uint64		`json:"startBatchId"`
	T				string		`json:"t"`
	R				string		`json:"r"`
}

func (ref *proofOfExecutionDTO) toStruct(networkType NetworkType) (*ProofOfExecution, error) {
	return &ProofOfExecution{
		StartBatchId: ref.StartBatchId,
		T: []byte(ref.T),
		R: []byte(ref.T),
	}, nil
} 

type executorInfoDTO struct {
	ExecutorKey				string	`json:"executorKey"`
	NextBatchToApproave 	uint64 	`json:"nextBatchToApproave"`
	proofOfExecutionDTO
}

func (ref *executorInfoDTO) toStruct(networkType NetworkType) (*ExecutorInfo, error) {
	poex, err := ref.proofOfExecutionDTO.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	executorKey, err := NewAccountFromPublicKey(ref.ExecutorKey, networkType)
	if err != nil {
		return nil, err
	}

	return &ExecutorInfo{
		ExecutorKey: executorKey,
		NextBatchToApproave: ref.NextBatchToApproave,
		PoEx: *poex,
	}, nil
}

type executorInfoDTOs []*executorInfoDTO

func (ref *executorInfoDTOs) toStruct(networkType NetworkType) ([]*ExecutorInfo, error) {
	var (
		dtos                    = *ref
		executorInfos = make([]*ExecutorInfo, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		executorInfos = append(executorInfos, info)
	}

	return executorInfos, nil
}

type completedCallDTO struct {
	CallId					hashDto 	`json:"callId"`
	Caller 					string		`json:"caller"`
	Status					uint16 		`json:"status"`
	ExecutionWork 			uint64DTO 	`json:"executionWork"`
	DownloadWork 			uint64DTO 	`json:"downloadWork"`
}

func (ref *completedCallDTO) toStruct(networkType NetworkType) (*CompletedCall, error) {
	callId, err := ref.CallId.Hash()
	if err != nil {
		return nil, err
	}

	caller, err := NewAccountFromPublicKey(ref.Caller, networkType)
	if err != nil {
		return nil, err
	}

	return &CompletedCall{
		CallId: callId,
		Caller: caller,
		Status: ref.Status,
		ExecutionWork: ref.ExecutionWork.toStruct(),
		DownloadWork: ref.DownloadWork.toStruct(),
	}, nil
}

type completedCallDTOs []*completedCallDTO

func (ref *completedCallDTOs) toStruct(networkType NetworkType) ([]*CompletedCall, error) {
	var (
		dtos                    = *ref
		completedCalls = make([]*CompletedCall, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		completedCalls = append(completedCalls, info)
	}

	return completedCalls, nil
}

type batchDTO struct {
	BatchId 						uint64		`json:"batchId"`	
	Success 						bool		`json:"success"`		
	PoExVerificationInformation		string		`json:"poExVerificationInformation"`
	completedCallDTOs
}

func (ref *batchDTO) toStruct(networkType NetworkType) (*Batch, error) {
	completedCalls, err := ref.completedCallDTOs.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &Batch{
		BatchId: ref.BatchId,
		Success: ref.Success,
		PoExVerificationInformation: []byte(ref.PoExVerificationInformation),
		CompletedCalls: completedCalls,
	}, nil
}

type batchDTOs []*batchDTO

func (ref *batchDTOs) toStruct(networkType NetworkType) ([]*Batch, error) {
	var (
		dtos                    = *ref
		batches = make([]*Batch, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		batches = append(batches, info)
	}

	return batches, nil
}

type releasedTransactionDTO struct {
	ReleasedTransactionHash		hashDto		`json:"releasedTransactionHash"`	
}

func (ref *releasedTransactionDTO) toStruct(networkType NetworkType) (*ReleasedTransaction, error) {
	releasedTransactionHash, err := ref.ReleasedTransactionHash.Hash()
	if err != nil {
		return nil, err
	}

	return &ReleasedTransaction{
		ReleasedTransactionHash: releasedTransactionHash,
	}, nil
}

type releasedTransactionDTOs []*releasedTransactionDTO

func (ref *releasedTransactionDTOs) toStruct(networkType NetworkType) ([]*ReleasedTransaction, error) {
	var (
		dtos                    = *ref
		releasedTransactions = make([]*ReleasedTransaction, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		releasedTransactions = append(releasedTransactions, info)
	}

	return releasedTransactions, nil
}

type superContractV2DTO struct {
	SuperContractKey 					string
	DriveKey 							string
	ExecutionPaymentKey					string
	Assignee 							string
	Creator 							string
	DeploymentBaseModificationsInfo 	hashDto
	automaticExecutionsInfoDTO 			
	contractCallDTOs 						
	executorInfoDTOs 						
	batchDTOs 							
	releasedTransactionDTOs
}

func (ref *superContractV2DTO) toStruct(networkType NetworkType) (*SuperContractV2, error) {
	superContractKey, err := NewAccountFromPublicKey(ref.SuperContractKey, networkType)
	if err != nil {
		return nil, err
	}

	driveKey, err := NewAccountFromPublicKey(ref.DriveKey, networkType)
	if err != nil {
		return nil, err
	}

	executionPaymentKey, err := NewAccountFromPublicKey(ref.ExecutionPaymentKey, networkType)
	if err != nil {
		return nil, err
	}

	assignee, err := NewAccountFromPublicKey(ref.Assignee, networkType)
	if err != nil {
		return nil, err
	}

	creator, err := NewAccountFromPublicKey(ref.Creator, networkType)
	if err != nil {
		return nil, err
	}

	deploymentBaseModificationsInfo, err := ref.DeploymentBaseModificationsInfo.Hash()
	if err != nil {
		return nil, err
	}

	automaticExecutionsInfo, err := ref.automaticExecutionsInfoDTO.toStruct(networkType)

	requestedCalls, err := ref.contractCallDTOs.toStruct(networkType)
	executorsInfo, err := ref.executorInfoDTOs.toStruct(networkType)
	batches, err := ref.batchDTOs.toStruct(networkType)
	releasedTransaction, err := ref.releasedTransactionDTOs.toStruct(networkType)
	

	return &SuperContractV2{
		SuperContractKey: superContractKey,
		DriveKey: driveKey,
		ExecutionPaymentKey: executionPaymentKey,
		Assignee: assignee,
		Creator: creator,
		DeploymentBaseModificationsInfo: deploymentBaseModificationsInfo,
		AutomaticExecutionsInfo: automaticExecutionsInfo,
		RequestedCalls: requestedCalls,
		ExecutorsInfo: executorsInfo,
		Batches: batches,
		ReleasedTransactions: releasedTransaction,
	}, nil
}
