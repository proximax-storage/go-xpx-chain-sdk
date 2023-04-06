// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import "fmt"

type automaticExecutionsInfoDTO struct {
	AutomaticExecutionFileName          string    `json:"automaticExecutionFileName"`
	AutomaticExecutionsFunctionName     string    `json:"automaticExecutionsFunctionName"`
	AutomaticExecutionsNextBlockToCheck uint64DTO `json:"automaticExecutionsNextBlockToCheck"`
	AutomaticExecutionCallPayment       uint64DTO `json:"automaticExecutionCallPayment"`
	AutomaticDownloadCallPayment        uint64DTO `json:"automaticDownloadCallPayment"`
	AutomatedExecutionsNumber           uint32    `json:"automatedExecutionsNumber"`
	AutomaticExecutionsPrepaidSince     uint64DTO `json:"automaticExecutionsPrepaidSince"`
}

func (ref *automaticExecutionsInfoDTO) toStruct(networkType NetworkType) (*AutomaticExecutionsInfo, error) {

	return &AutomaticExecutionsInfo{
		AutomaticExecutionFileName:          ref.AutomaticExecutionFileName,
		AutomaticExecutionsFunctionName:     ref.AutomaticExecutionsFunctionName,
		AutomaticExecutionsNextBlockToCheck: ref.AutomaticExecutionsNextBlockToCheck.toStruct(),
		AutomaticExecutionCallPayment:       ref.AutomaticExecutionCallPayment.toStruct(),
		AutomaticDownloadCallPayment:        ref.AutomaticDownloadCallPayment.toStruct(),
		AutomatedExecutionsNumber:           ref.AutomatedExecutionsNumber,
		AutomaticExecutionsPrepaidSince:     ref.AutomaticExecutionsPrepaidSince.toStruct(),
	}, nil
}

type servicePaymentDTO struct {
	MosaicId uint64DTO `json:"mosaicId"`
	Amount   uint64DTO `json:"amount"`
}

func (ref *servicePaymentDTO) toStruct(NetworkType NetworkType) (*ServicePayment, error) {

	mosaicId, err := NewMosaicId(ref.MosaicId.toUint64())
	if err != nil {
		return nil, err
	}

	return &ServicePayment{
		MosaicId: mosaicId,
		Amount:   ref.Amount.toStruct(),
	}, nil
}

type servicePaymentDtos []*servicePaymentDTO

func (ref *servicePaymentDtos) toStruct(networkType NetworkType) ([]*ServicePayment, error) {
	var (
		dtos            = *ref
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
	CallId               hashDto            `json:"callId"`
	Caller               string             `json:"caller"`
	FileName             string             `json:"fileName"`
	FunctionName         string             `json:"functionName"`
	ActualArguments      string             `json:"actualArguments"`
	ExecutionCallPayment uint64DTO          `json:"executionCallPayment"`
	DownloadCallPayment  uint64DTO          `json:"downloadCallPayment"`
	ServicePayments      servicePaymentDtos `json:"servicePayments"`
	BlockHeight          uint64DTO          `json:"blockHeight"`
}

func (ref *contractCallDTO) toStruct(networkType NetworkType) (*ContractCall, error) {
	servicePayments, err := ref.ServicePayments.toStruct(networkType)
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
		CallId:               callId,
		Caller:               caller,
		FileName:             ref.FileName,
		FunctionName:         ref.FunctionName,
		ActualArguments:      []byte(ref.ActualArguments),
		ExecutionCallPayment: ref.ExecutionCallPayment.toStruct(),
		DownloadCallPayment:  ref.DownloadCallPayment.toStruct(),
		ServicePayments:      servicePayments,
		BlockHeight:          ref.BlockHeight.toStruct(),
	}, nil
}

type contractCallDTOs []*contractCallDTO

func (ref *contractCallDTOs) toStruct(networkType NetworkType) ([]*ContractCall, error) {
	var (
		dtos          = *ref
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
	StartBatchId uint64DTO `json:"startBatchId"`
	T            string    `json:"T"`
	R            string    `json:"R"`
}

func (ref *proofOfExecutionDTO) toStruct(networkType NetworkType) (*ProofOfExecution, error) {
	return &ProofOfExecution{
		StartBatchId: ref.StartBatchId.toUint64(),
		T:            []byte(ref.T),
		R:            []byte(ref.R),
	}, nil
}

type executorInfoDTO struct {
	ExecutorKey        string              `json:"executorKey"`
	NextBatchToApprove uint64DTO           `json:"nextBatchToApprove"`
	ProofOfExecution   proofOfExecutionDTO `json:"proofOfExecution"`
}

func (ref *executorInfoDTO) toStruct(networkType NetworkType) (*ExecutorInfo, error) {
	poex, err := ref.ProofOfExecution.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	executorKey, err := NewAccountFromPublicKey(ref.ExecutorKey, networkType)
	if err != nil {
		return nil, err
	}

	return &ExecutorInfo{
		ExecutorKey:        executorKey,
		NextBatchToApprove: ref.NextBatchToApprove.toUint64(),
		PoEx:               *poex,
	}, nil
}

type executorInfoDTOs []*executorInfoDTO

func (ref *executorInfoDTOs) toStruct(networkType NetworkType) ([]*ExecutorInfo, error) {
	var (
		dtos          = *ref
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
	CallId        hashDto   `json:"callId"`
	Caller        string    `json:"caller"`
	Status        uint16    `json:"status"`
	ExecutionWork uint64DTO `json:"executionWork"`
	DownloadWork  uint64DTO `json:"downloadWork"`
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
		CallId:        callId,
		Caller:        caller,
		Status:        ref.Status,
		ExecutionWork: ref.ExecutionWork.toStruct(),
		DownloadWork:  ref.DownloadWork.toStruct(),
	}, nil
}

type completedCallDTOs []*completedCallDTO

func (ref *completedCallDTOs) toStruct(networkType NetworkType) ([]*CompletedCall, error) {
	var (
		dtos           = *ref
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
	BatchId                     uint64DTO         `json:"batchId"`
	Success                     bool              `json:"success"`
	PoExVerificationInformation string            `json:"poExVerificationInformation"`
	CompletedCalls              completedCallDTOs `json:"completedCalls"`
}

func (ref *batchDTO) toStruct(networkType NetworkType) (*Batch, error) {
	completedCalls, err := ref.CompletedCalls.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &Batch{
		BatchId:                     ref.BatchId.toUint64(),
		Success:                     ref.Success,
		PoExVerificationInformation: []byte(ref.PoExVerificationInformation),
		CompletedCalls:              completedCalls,
	}, nil
}

type batchDTOs []*batchDTO

func (ref *batchDTOs) toStruct(networkType NetworkType) ([]*Batch, error) {
	var (
		dtos    = *ref
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

type releasedTransactionDTOs struct {
	ReleasedTransactions []string `json:"releasedTransactionHashs"`
}

func (ref *releasedTransactionDTOs) toStruct(networkType NetworkType) ([]string, error) {
	return ref.ReleasedTransactions, nil
}

type superContractV2DTO struct {
	SuperContractV2 struct {
		SuperContractKey                string                      `json:"superContractKey"`
		DriveKey                        string                      `json:"driveKey"`
		ExecutionPaymentKey             string                      `json:"executionPaymentKey"`
		Assignee                        string                      `json:"assignee"`
		Creator                         string                      `json:"creator"`
		DeploymentBaseModificationsInfo hashDto                     `json:"deploymentBaseModificationsInfo"`
		AutomaticExecutionsInfos        *automaticExecutionsInfoDTO `json:"automaticExecutionsInfo"`
		ContractCalls                   contractCallDTOs            `json:"requestCalls"`
		ExecutorInfos                   executorInfoDTOs            `json:"executorsInfo"`
		Batches                         batchDTOs                   `json:"batches"`
		ReleasedTransactions            []string                    `json:"releasedTransactions"`
	}
}

func (ref *superContractV2DTO) toStruct(networkType NetworkType) (*SuperContractV2, error) {
	superContractKey, err := NewAccountFromPublicKey(ref.SuperContractV2.SuperContractKey, networkType)
	if err != nil {
		return nil, err
	}

	driveKey, err := NewAccountFromPublicKey(ref.SuperContractV2.DriveKey, networkType)
	if err != nil {
		return nil, err
	}

	executionPaymentKey, err := NewAccountFromPublicKey(ref.SuperContractV2.ExecutionPaymentKey, networkType)
	if err != nil {
		return nil, err
	}

	assignee, err := NewAccountFromPublicKey(ref.SuperContractV2.Assignee, networkType)
	if err != nil {
		return nil, err
	}

	creator, err := NewAccountFromPublicKey(ref.SuperContractV2.Creator, networkType)
	if err != nil {
		return nil, err
	}

	deploymentBaseModificationsInfo, err := ref.SuperContractV2.DeploymentBaseModificationsInfo.Hash()
	if err != nil {
		return nil, err
	}

	automaticExecutionsInfo := &AutomaticExecutionsInfo{}
	if ref.SuperContractV2.AutomaticExecutionsInfos != nil {
		automaticExecutionsInfo, err = ref.SuperContractV2.AutomaticExecutionsInfos.toStruct(networkType)
		if err != nil {
			return nil, fmt.Errorf("sdk.superContractDto.toStruct SuperContractV2.AutomaticExecutionsInfos.toStruct: %v", err)
		}
	}

	requestedCalls, err := ref.SuperContractV2.ContractCalls.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.SuperContractV2.toStruct SuperContractV2.ContractCalls.toStruct: %v", err)
	}
	executorsInfo, err := ref.SuperContractV2.ExecutorInfos.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.SuperContractV2.toStruct SuperContractV2.ExecutorInfos.toStruct: %v", err)
	}
	batches, err := ref.SuperContractV2.Batches.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.SuperContractV2.toStruct SuperContractV2.Batches.toStruct: %v", err)
	}

	var infos = make([]*Hash, 0, len(ref.SuperContractV2.ReleasedTransactions))

	for _, iter := range ref.SuperContractV2.ReleasedTransactions {
		info, err := StringToHash(iter)
		if err != nil {
			return nil, err
		}

		infos = append(infos, info)
	}

	return &SuperContractV2{
		SuperContractKey:                superContractKey,
		DriveKey:                        driveKey,
		ExecutionPaymentKey:             executionPaymentKey,
		Assignee:                        assignee,
		Creator:                         creator,
		DeploymentBaseModificationsInfo: deploymentBaseModificationsInfo,
		AutomaticExecutionsInfo:         automaticExecutionsInfo,
		RequestedCalls:                  requestedCalls,
		ExecutorsInfo:                   executorsInfo,
		Batches:                         batches,
		ReleasedTransactions:            infos,
	}, nil
}

type superContractV2PageDTO struct {
	SuperContractsV2 []superContractV2DTO `json:"data"`

	Pagination struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (t *superContractV2PageDTO) toStruct(networkType NetworkType) (*SuperContractsV2Page, error) {
	page := &SuperContractsV2Page{
		SuperContractsV2: make([]*SuperContractV2, len(t.SuperContractsV2)),
		Pagination: Pagination{
			TotalEntries: t.Pagination.TotalEntries,
			PageNumber:   t.Pagination.PageNumber,
			PageSize:     t.Pagination.PageSize,
			TotalPages:   t.Pagination.TotalPages,
		},
	}

	var err error
	for i, t := range t.SuperContractsV2 {
		page.SuperContractsV2[i], err = t.toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}

	return page, nil
}

type rawProofsOfExecutionDTO struct {
	StartBatchId uint64DTO `json:"startBatchId"`
	T            string    `json:"T"`
	R            string    `json:"R"`
	F            string    `json:"F"`
	K            string    `json:"K"`
}

func (ref *rawProofsOfExecutionDTO) toStruct(networkType NetworkType) (*RawProofsOfExecution, error) {
	return &RawProofsOfExecution{
		StartBatchId: ref.StartBatchId.toUint64(),
		T:            []byte(ref.T),
		R:            []byte(ref.R),
		F:            []byte(ref.F),
		K:            []byte(ref.K),
	}, nil

}

type rawProofsOfExecutionDTOs []*rawProofsOfExecutionDTO

func (ref *rawProofsOfExecutionDTOs) toStruct(networkType NetworkType) ([]*RawProofsOfExecution, error) {
	var (
		dtos  = *ref
		poexs = make([]*RawProofsOfExecution, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		poexs = append(poexs, info)
	}

	return poexs, nil
}

type extendedCallDigestDTO struct {
	CallId                  hashDto   `json:"callId"`
	Manual                  bool      `json:"manual"`
	Block                   uint64DTO `json:"block"`
	Status                  uint16    `json:"status"`
	ReleasedTransactionHash hashDto   `json:"releasedTransactionHash"`
}

func (ref *extendedCallDigestDTO) toStruct(networkType NetworkType) (*ExtendedCallDigest, error) {
	callId, err := ref.CallId.Hash()
	if err != nil {
		return nil, err
	}

	releasedTransactionHash, err := ref.ReleasedTransactionHash.Hash()
	if err != nil {
		return nil, err
	}

	return &ExtendedCallDigest{
		CallId:                  callId,
		Manual:                  ref.Manual,
		Block:                   ref.Block.toStruct(),
		Status:                  ref.Status,
		ReleasedTransactionHash: releasedTransactionHash,
	}, nil
}

type extendedCallDigestDTOs []extendedCallDigestDTO

func (ref *extendedCallDigestDTOs) toStruct(networkType NetworkType) ([]*ExtendedCallDigest, error) {
	var (
		dtos  = *ref
		extendedCallDigests = make([]*ExtendedCallDigest, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		extendedCallDigests = append(extendedCallDigests, info)
	}

	return extendedCallDigests, nil
}

type callPaymentDTO struct {
	ExecutionPayment uint64DTO `json:"executionPayment"`
	DownloadPayment  uint64DTO `json:"downloadPayment"`
}

func (ref *callPaymentDTO) toStruct(networkType NetworkType) (*CallPayment, error) {
	return &CallPayment{
		ExecutionPayment: ref.ExecutionPayment.toStruct(),
		DownloadPayment:  ref.DownloadPayment.toStruct(),
	}, nil
}

type callPaymentDTOs []*callPaymentDTO

func (ref *callPaymentDTOs) toStruct(networkType NetworkType) ([]*CallPayment, error) {
	var (
		dtos         = *ref
		callPayments = make([]*CallPayment, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		callPayments = append(callPayments, info)
	}

	return callPayments, nil
}

type shortCallDigestDTO struct {
	CallId hashDto   `json:"callId"`
	Manual bool      `json:"manual"`
	Block  uint64DTO `json:"block"`
}

func (ref *shortCallDigestDTO) toStruct(networkType NetworkType) (*ShortCallDigest, error) {
	callId, err := ref.CallId.Hash()
	if err != nil {
		return nil, err
	}

	return &ShortCallDigest{
		CallId: callId,
		Manual: ref.Manual,
		Block:  ref.Block.toStruct(),
	}, nil
}

type shortCallDigestDTOs []shortCallDigestDTO

func (ref *shortCallDigestDTOs) toStruct(networkType NetworkType) ([]*ShortCallDigest, error) {
	var (
		dtos  = *ref
		shortCallDigests = make([]*ShortCallDigest, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		shortCallDigests = append(shortCallDigests, info)
	}

	return shortCallDigests, nil
}


type endBatchExecutionDTO struct {
	ContractKey                         string                   `json:"contractKey"`
	BatchId                             uint64DTO                `json:"batchId"`
	AutomaticExecutionsNextBlockToCheck uint64DTO                `json:"automaticExecutionsNextBlockToCheck"`
	PublicKeys                          []string                 `json:"publicKeys"`
	Signatures                          []signatureDto           `json:"signatures"`
	ProofsOfExecutions                  rawProofsOfExecutionDTOs `json:"proofsOfExecutions"`
	CallPayments                        callPaymentDTOs          `json:"callPayments"`
}

func (ref *endBatchExecutionDTO) toStruct(networkType NetworkType) (*EndBatchExecution, error) {
	contractKey, err := NewAccountFromPublicKey(ref.ContractKey, networkType)
	if err != nil {
		return nil, err
	}

	keys := make([]*PublicAccount, len(ref.PublicKeys))
	for i, k := range ref.PublicKeys {
		keys[i], err = NewAccountFromPublicKey(k, networkType)
		if err != nil {
			return nil, err
		}
	}

	signatures := make([]*Signature, len(ref.Signatures))
	for i, s := range ref.Signatures {
		signatures[i], err = s.Signature()
		if err != nil {
			return nil, err
		}
	}

	poexs, err := ref.ProofsOfExecutions.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.toStruct ProofsOfExecutions.toStruct: %v", err)
	}

	callPayments, err := ref.CallPayments.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.toStruct CallPayments.toStruct: %v", err)
	}

	return &EndBatchExecution{
		ContractKey:                         contractKey,
		BatchId:                             ref.BatchId.toUint64(),
		AutomaticExecutionsNextBlockToCheck: ref.AutomaticExecutionsNextBlockToCheck.toStruct(),
		PublicKeys:                          keys,
		Signatures:                          signatures,
		ProofsOfExecutions:                  poexs,
		CallPayments:                        callPayments,
	}, nil
}
