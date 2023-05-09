// Copyright 2023 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

// parse function for deploy contract txn and manual call txn
func parseData(
	builder *flatbuffers.Builder,
	executionCallPayment *Amount,
	downloadCallPayment *Amount,
	fileName string,
	functionName string,
	servicePayments []*Mosaic,
) (flatbuffers.UOffsetT, flatbuffers.UOffsetT, flatbuffers.UOffsetT, flatbuffers.UOffsetT, flatbuffers.UOffsetT) {

	executionCall := transactions.TransactionBufferCreateUint32Vector(builder, executionCallPayment.toArray())
	downloadCall := transactions.TransactionBufferCreateUint32Vector(builder, downloadCallPayment.toArray())

	fileBytes := []byte(fileName)
	file := transactions.TransactionBufferCreateByteVector(builder, fileBytes)
	functionBytes := []byte(functionName)
	function := transactions.TransactionBufferCreateByteVector(builder, functionBytes)

	mb := make([]flatbuffers.UOffsetT, len(servicePayments))
	for i, it := range servicePayments {
		id := transactions.TransactionBufferCreateUint32Vector(builder, it.AssetId.toArray())
		am := transactions.TransactionBufferCreateUint32Vector(builder, it.Amount.toArray())
		transactions.MosaicBufferStart(builder)
		transactions.MosaicBufferAddId(builder, id)
		transactions.MosaicBufferAddAmount(builder, am)
		mb[i] = transactions.MosaicBufferEnd(builder)
	}
	mV := transactions.TransactionBufferCreateUOffsetVector(builder, mb)
	return executionCall, downloadCall, file, function, mV
}

// automatic executions payment transaction
func NewAutomaticExecutionsPaymentTransaction(
	deadline *Deadline,
	contractKey *PublicAccount,
	automaticExecutionsNumber uint32,
	networkType NetworkType,
) (*AutomaticExecutionsPaymentTransaction, error) {

	tx := AutomaticExecutionsPaymentTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     AutomaticExecutionsPaymentVersion,
			Type:        AutomaticExecutionsPayment,
			NetworkType: networkType,
		},
		ContractKey:               contractKey,
		AutomaticExecutionsNumber: automaticExecutionsNumber,
	}

	return &tx, nil
}

func (tx *AutomaticExecutionsPaymentTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AutomaticExecutionsPaymentTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"ContractKey": %s,
			"AutomaticExecutionsNumber": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.ContractKey.String(),
		tx.AutomaticExecutionsNumber,
	)
}

func (tx *AutomaticExecutionsPaymentTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	contractKeyB, err := hex.DecodeString(tx.ContractKey.PublicKey)
	if err != nil {
		return nil, err
	}

	contractKeyV := transactions.TransactionBufferCreateByteVector(builder, contractKeyB)

	transactions.AutomaticExecutionsPaymentTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.AutomaticExecutionsPaymentTransactionBufferAddContractKey(builder, contractKeyV)
	transactions.AutomaticExecutionsPaymentTransactionBufferAddAutomaticExecutionsNumber(builder, tx.AutomaticExecutionsNumber)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return automaticExecutionsPaymentTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AutomaticExecutionsPaymentTransaction) Size() int {
	return AutomaticExecutionsPaymentHeaderSize
}

type automaticExecutionsPaymentTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ContractKey               string `json:"contractKey"`
		AutomaticExecutionsNumber uint32 `json:"automaticExecutionsNumber"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *automaticExecutionsPaymentTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	contractKey, err := NewAccountFromPublicKey(dto.Tx.ContractKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &AutomaticExecutionsPaymentTransaction{
		*atx,
		contractKey,
		dto.Tx.AutomaticExecutionsNumber,
	}, nil
}

// manual call transaction
func NewManualCallTransaction(
	deadline *Deadline,
	contractKey *PublicAccount,
	executionCallPayment Amount,
	downloadCallPayment Amount,
	fileName string,
	functionName string,
	actualArguments []byte,
	servicePayments []*Mosaic,
	networkType NetworkType,
) (*ManualCallTransaction, error) {
	tx := ManualCallTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     AutomaticExecutionsPaymentVersion,
			Type:        AutomaticExecutionsPayment,
			NetworkType: networkType,
		},
		ContractKey:          contractKey,
		ExecutionCallPayment: executionCallPayment,
		DownloadCallPayment:  downloadCallPayment,
		FileName:             fileName,
		FunctionName:         functionName,
		ActualArguments:      actualArguments,
		ServicePayments:      servicePayments,
	}
	return &tx, nil
}

func (tx *ManualCallTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ManualCallTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	contractKeyB, err := hex.DecodeString(tx.ContractKey.PublicKey)
	if err != nil {
		return nil, err
	}

	contractKeyV := transactions.TransactionBufferCreateByteVector(builder, contractKeyB)
	executionCallPayment, downloadCallPayment, fileName, functionName, mV := parseData(builder,
		&tx.ExecutionCallPayment, &tx.DownloadCallPayment, tx.FileName, tx.FunctionName, tx.ServicePayments)
	actualArgV := transactions.TransactionBufferCreateByteVector(builder, tx.ActualArguments)

	fileNameSizeBuf := make([]byte, FileNameSize)
	binary.LittleEndian.PutUint16(fileNameSizeBuf, uint16(len(tx.FileName)))
	fileNameSizeV := transactions.TransactionBufferCreateByteVector(builder, fileNameSizeBuf)

	functionNameSizeBuf := make([]byte, FunctionNameSize)
	binary.LittleEndian.PutUint16(functionNameSizeBuf, uint16(len(tx.FunctionName)))
	functionNameSizeV := transactions.TransactionBufferCreateByteVector(builder, functionNameSizeBuf)

	actualArgumentsSizeBuf := make([]byte, ActualArgumentsSize)
	binary.LittleEndian.PutUint16(actualArgumentsSizeBuf, uint16(len(tx.ActualArguments)))
	actualArgumentsSizeV := transactions.TransactionBufferCreateByteVector(builder, actualArgumentsSizeBuf)

	transactions.ManualCallTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.ManualCallTransactionBufferAddContractKey(builder, contractKeyV)
	transactions.ManualCallTransactionBufferAddExecutionCallPayment(builder, executionCallPayment)
	transactions.ManualCallTransactionBufferAddDownloadCallPayment(builder, downloadCallPayment)
	transactions.ManualCallTransactionBufferAddFileNameSize(builder, fileNameSizeV)
	transactions.ManualCallTransactionBufferAddFileName(builder, fileName)
	transactions.ManualCallTransactionBufferAddFunctionNameSize(builder, functionNameSizeV)
	transactions.ManualCallTransactionBufferAddFunctionName(builder, functionName)
	transactions.ManualCallTransactionBufferAddActualArgumentsSize(builder, actualArgumentsSizeV)
	transactions.ManualCallTransactionBufferAddActualArguments(builder, actualArgV)
	transactions.ManualCallTransactionBufferAddServicePaymentsCount(builder, byte(len(tx.ServicePayments)))
	transactions.ManualCallTransactionBufferAddServicePayments(builder, mV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return manualCallTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *ManualCallTransaction) Size() int {
	return ManualCallHeaderSize +
		len([]byte(tx.FileName)) +
		len([]byte(tx.FunctionName)) +
		len(tx.ActualArguments) +
		len(tx.ServicePayments)*(MosaicIdSize+AmountSize)
}

type manualCallTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ContractKey          string       `json:"contractKey"`
		FileNameSize         uint16       `json:"fileNameSize"`
		FunctionNameSize     uint16       `json:"functionNameSize"`
		ActualArgumentsSize  uint16       `json:"actualArgumentsSize"`
		ExecutionCallPayment uint64DTO    `json:"executionCallPayment"`
		DownloadCallPayment  uint64DTO    `json:"downloadCallPayment"`
		ServicePaymentsCount uint8        `json:"servicePaymentsCount"`
		FileName             string       `json:"fileName"`
		FunctionName         string       `json:"functionName"`
		ActualArguments      string       `json:"actualArguments"`
		ServicePayments      []*mosaicDTO `json:"servicePayments"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *manualCallTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	contractKey, err := NewAccountFromPublicKey(dto.Tx.ContractKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	mosaics := make([]*Mosaic, len(dto.Tx.ServicePayments))

	for i, mosaic := range dto.Tx.ServicePayments {
		msc, err := mosaic.toStruct()
		if err != nil {
			return nil, err
		}

		mosaics[i] = msc
	}

	return &ManualCallTransaction{
		*atx,
		contractKey,
		dto.Tx.ExecutionCallPayment.toStruct(),
		dto.Tx.DownloadCallPayment.toStruct(),
		dto.Tx.FileName,
		dto.Tx.FunctionName,
		[]byte(dto.Tx.ActualArguments),
		mosaics,
	}, nil
}

// deploy contract transaction
func NewDeployContractTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	executionCallPayment Amount,
	downloadCallPayment Amount,
	automaticExecutionCallPayment Amount,
	automaticDownloadCallPayment Amount,
	automaticExecutionsNumber uint32,
	assignee *PublicAccount,
	fileName string,
	functionName string,
	actualArguments []byte,
	servicePayments []*Mosaic,
	automaticExecutionFileName string,
	automaticExecutionFunctionName string,
	networkType NetworkType,
) (*DeployContractTransaction, error) {

	tx := DeployContractTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     DeployContractVersion,
			Type:        DeployContract,
			NetworkType: networkType,
		},
		DriveKey:                       driveKey,
		ExecutionCallPayment:           executionCallPayment,
		DownloadCallPayment:            downloadCallPayment,
		AutomaticExecutionCallPayment:  automaticExecutionCallPayment,
		AutomaticDownloadCallPayment:   automaticDownloadCallPayment,
		AutomaticExecutionsNumber:      automaticExecutionsNumber,
		Assignee:                       assignee,
		FileName:                       fileName,
		FunctionName:                   functionName,
		ActualArguments:                actualArguments,
		ServicePayments:                servicePayments,
		AutomaticExecutionFileName:     automaticExecutionFileName,
		AutomaticExecutionFunctionName: automaticExecutionFunctionName,
	}

	return &tx, nil
}

func (tx *DeployContractTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DeployContractTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	assigneeB, err := hex.DecodeString(tx.Assignee.PublicKey)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	assigneeV := transactions.TransactionBufferCreateByteVector(builder, assigneeB)

	executionCallPayment, downloadCallPayment, fileName, functionName, mV := parseData(builder,
		&tx.ExecutionCallPayment, &tx.DownloadCallPayment, tx.FileName, tx.FunctionName, tx.ServicePayments)

	actualArgV := transactions.TransactionBufferCreateByteVector(builder, tx.ActualArguments)

	automaticExecutionFileName := transactions.TransactionBufferCreateByteVector(builder, []byte(tx.AutomaticExecutionFileName))
	automaticExecutionFunctionName := transactions.TransactionBufferCreateByteVector(builder, []byte(tx.AutomaticExecutionFunctionName))

	automaticExecutionCallPayment := transactions.TransactionBufferCreateUint32Vector(builder, tx.AutomaticExecutionCallPayment.toArray())
	automaticDownloadCallPayment := transactions.TransactionBufferCreateUint32Vector(builder, tx.AutomaticDownloadCallPayment.toArray())

	fileNameSizeBuf := make([]byte, FileNameSize)
	binary.LittleEndian.PutUint16(fileNameSizeBuf, uint16(len(tx.FileName)))
	fileNameSizeV := transactions.TransactionBufferCreateByteVector(builder, fileNameSizeBuf)

	functionNameSizeBuf := make([]byte, FunctionNameSize)
	binary.LittleEndian.PutUint16(functionNameSizeBuf, uint16(len(tx.FunctionName)))
	functionNameSizeV := transactions.TransactionBufferCreateByteVector(builder, functionNameSizeBuf)

	actualArgumentsSizeBuf := make([]byte, ActualArgumentsSize)
	binary.LittleEndian.PutUint16(actualArgumentsSizeBuf, uint16(len(tx.ActualArguments)))
	actualArgumentsSizeV := transactions.TransactionBufferCreateByteVector(builder, actualArgumentsSizeBuf)

	automaticExecutionsFileSizeBuf := make([]byte, AutomaticExecutionsFileNameSize)
	binary.LittleEndian.PutUint16(automaticExecutionsFileSizeBuf, uint16(len(tx.AutomaticExecutionFileName)))
	automaticExecutionsFileSizeV := transactions.TransactionBufferCreateByteVector(builder, automaticExecutionsFileSizeBuf)

	automaticExecutionsFunctionNameSizeBuf := make([]byte, AutomaticExecutionsFunctionNameSize)
	binary.LittleEndian.PutUint16(automaticExecutionsFunctionNameSizeBuf, uint16(len(tx.AutomaticExecutionFunctionName)))
	automaticExecutionsFunctionNameSizeV := transactions.TransactionBufferCreateByteVector(builder, automaticExecutionsFunctionNameSizeBuf)

	transactions.DeployContractTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DeployContractTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.DeployContractTransactionBufferAddExecutionCallPayment(builder, executionCallPayment)
	transactions.DeployContractTransactionBufferAddDownloadCallPayment(builder, downloadCallPayment)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionCallPayment(builder, automaticExecutionCallPayment)
	transactions.DeployContractTransactionBufferAddAutomaticDownloadCallPayment(builder, automaticDownloadCallPayment)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionsNumber(builder, tx.AutomaticExecutionsNumber)
	transactions.DeployContractTransactionBufferAddAssignee(builder, assigneeV)
	transactions.DeployContractTransactionBufferAddFileNameSize(builder, fileNameSizeV)
	transactions.DeployContractTransactionBufferAddFileName(builder, fileName)
	transactions.DeployContractTransactionBufferAddFunctionNameSize(builder, functionNameSizeV)
	transactions.DeployContractTransactionBufferAddFunctionName(builder, functionName)
	transactions.DeployContractTransactionBufferAddActualArgumentsSize(builder, actualArgumentsSizeV)
	transactions.DeployContractTransactionBufferAddActualArguments(builder, actualArgV)
	transactions.DeployContractTransactionBufferAddServicePaymentsCount(builder, uint8(len(tx.ServicePayments)))
	transactions.DeployContractTransactionBufferAddServicePayments(builder, mV)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionFileNameSize(builder, automaticExecutionsFileSizeV)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionFileName(builder, automaticExecutionFileName)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionFunctionNameSize(builder, automaticExecutionsFunctionNameSizeV)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionFunctionName(builder, automaticExecutionFunctionName)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return deployContractTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DeployContractTransaction) Size() int {
	return DeployContractHeaderSize +
		len(tx.FileName) +
		len(tx.FunctionName) +
		len(tx.AutomaticExecutionFileName) +
		len(tx.AutomaticExecutionFunctionName) +
		len(tx.ActualArguments) +
		len(tx.ServicePayments)*(MosaicIdSize+AmountSize)
}

type deployContractTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey                           string       `json:"driveKey"`
		FileNameSize                       uint16       `json:"fileNameSize"`
		FunctionNameSize                   uint16       `json:"functionNameSize"`
		ActualArgumentsSize                uint16       `json:"actualArgumentsSize"`
		ExecutionCallPayment               uint64DTO    `json:"executionCallPayment"`
		DownloadCallPayment                uint64DTO    `json:"downloadCallPayment"`
		ServicePaymentsCount               uint8        `json:"servicePaymentsCount"`
		AutomaticExecutionFileNameSize     uint16       `json:"automaticExecutionFileNameSize"`
		AutomaticExecutionFunctionNameSize uint16       `json:"automaticExecutionFunctionNameSize"`
		AutomaticExecutionCallPayment      uint64DTO    `json:"automaticExecutionCallPayment"`
		AutomaticDownloadCallPayment       uint64DTO    `json:"automaticDownloadCallPayment"`
		AutomaticExecutionsNumber          uint32       `json:"automaticExecutionsNumber"`
		Assignee                           string       `json:"assignee"`
		FileName                           string       `json:"fileName"`
		FunctionName                       string       `json:"functionName"`
		ActualArguments                    string       `json:"actualArguments"`
		ServicePayments                    []*mosaicDTO `json:"servicePayments"`
		AutomaticExecutionFileName         string       `json:"automaticExecutionFileName"`
		AutomaticExecutionFunctionName     string       `json:"automaticExecutionFunctionName"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *deployContractTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	driveKey, err := NewAccountFromPublicKey(dto.Tx.DriveKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	assignee, err := NewAccountFromPublicKey(dto.Tx.Assignee, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	mosaics := make([]*Mosaic, len(dto.Tx.ServicePayments))

	for i, mosaic := range dto.Tx.ServicePayments {
		msc, err := mosaic.toStruct()
		if err != nil {
			return nil, err
		}

		mosaics[i] = msc
	}

	return &DeployContractTransaction{
		*atx,
		driveKey,
		dto.Tx.ExecutionCallPayment.toStruct(),
		dto.Tx.DownloadCallPayment.toStruct(),
		dto.Tx.AutomaticExecutionCallPayment.toStruct(),
		dto.Tx.AutomaticDownloadCallPayment.toStruct(),
		dto.Tx.AutomaticExecutionsNumber,
		assignee,
		dto.Tx.FileName,
		dto.Tx.FunctionName,
		[]byte(dto.Tx.ActualArguments),
		mosaics,
		dto.Tx.AutomaticExecutionFileName,
		dto.Tx.AutomaticExecutionFunctionName,
	}, nil
}

// SuccessfulEndBatchExecutionTransaction bytes() is not sent by the client
func (tx *SuccessfulEndBatchExecutionTransaction) Bytes() ([]byte, error) {
	return nil, nil
}

func (tx *SuccessfulEndBatchExecutionTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *SuccessfulEndBatchExecutionTransaction) Size() int {
	return SuccessfulEndBatchExecutionHeaderSize
}

type successfulEndBatchExecutionTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		EndBatchExecution                       *endBatchExecutionDTO  `json:"endBatchExecution"`
		StorageHash                             hashDto                `json:"storageHash"`
		UsedSizedBytes                          uint64DTO              `json:"usedSizedBytes"`
		MetaFilesSizeBytes                      uint64DTO              `json:"metaFilesSizeBytes"`
		ProofOfExecutionVerificationInformation string                 `json:"proofOfExecutionVerificationInformation"`
		CallDigests                             extendedCallDigestDTOs `json:"callDigests"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *successfulEndBatchExecutionTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	endBatch, err := dto.Tx.EndBatchExecution.toStruct(atx.NetworkType)
	if err != nil {
		return nil, err
	}

	storageHash, err := dto.Tx.StorageHash.Hash()
	if err != nil {
		return nil, err
	}

	callDigests, err := dto.Tx.CallDigests.toStruct(atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &SuccessfulEndBatchExecutionTransaction{
		*atx,
		*endBatch,
		storageHash,
		dto.Tx.UsedSizedBytes.toUint64(),
		dto.Tx.MetaFilesSizeBytes.toUint64(),
		[]byte(dto.Tx.ProofOfExecutionVerificationInformation),
		callDigests,
	}, nil
}

// UnsuccessfulEndBatchExecutionTransaction bytes() is not sent by the client
func (tx *UnsuccessfulEndBatchExecutionTransaction) Bytes() ([]byte, error) {
	return nil, nil
}

func (tx *UnsuccessfulEndBatchExecutionTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *UnsuccessfulEndBatchExecutionTransaction) Size() int {
	return UnsuccessfulEndBatchExecutionHeaderSize
}

type unsuccessfulEndBatchExecutionTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		EndBatchExecution *endBatchExecutionDTO `json:"endBatchExecution"`
		CallDigests       shortCallDigestDTOs   `json:"callDigests"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *unsuccessfulEndBatchExecutionTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	endBatch, err := dto.Tx.EndBatchExecution.toStruct(atx.NetworkType)
	if err != nil {
		return nil, err
	}

	callDigests, err := dto.Tx.CallDigests.toStruct(atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &UnsuccessfulEndBatchExecutionTransaction{
		*atx,
		*endBatch,
		callDigests,
	}, nil
}
