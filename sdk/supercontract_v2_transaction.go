// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/hex"
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

// automatic executions payment transaction
func NewAutomaticExecutionsPaymentTransaction(
	deadline 					*Deadline,
	contractKey 				*PublicAccount,
	automaticExecutionsNumber 	uint32,
	networkType 				NetworkType,
) (*AutomaticExecutionsPaymentTransaction, error) {

	tx := AutomaticExecutionsPaymentTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     AutomaticExecutionsPaymentVersion,
			Type:        AutomaticExecutionsPayment,
			NetworkType: networkType,
		},
		ContractKey: contractKey,
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
		ContractKey 				string `json:"contractKey"`
		AutomaticExecutionsNumber 	uint32 `json:"automaticExecutionsNumber"`
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
	deadline 				*Deadline,
	contractKey 			*PublicAccount,
	fileNameSize 			uint16,
	functionNameSize 		uint16,
	actualArgumentsSize 	uint16,
	executionCallPayment 	Amount,
	downloadCallPayment 	Amount,
	servicePaymentsCount 	uint8,
	fileName 				string,
	functionName 			string,
	actualArguments 		string,
	servicePayments 		[]*MosaicId,
	networkType 			NetworkType,
) (*ManualCallTransaction, error) {
	tx := ManualCallTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     AutomaticExecutionsPaymentVersion,
			Type:        AutomaticExecutionsPayment,
			NetworkType: networkType,
		},
		ContractKey: 			contractKey,
		FileNameSize: 			fileNameSize,
		FunctionNameSize: 		functionNameSize,
		ActualArgumentsSize: 	actualArgumentsSize,
		ExecutionCallPayment: 	executionCallPayment,
		DownloadCallPayment: 	downloadCallPayment,
		ServicePaymentsCount: 	servicePaymentsCount,
		FileName: 				fileName,
		FunctionName: 			functionName,
		ActualArguments: 		actualArguments,
		ServicePayments: 		servicePayments,
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
	executionCallPayment := transactions.TransactionBufferCreateUint32Vector(builder, tx.ExecutionCallPayment.toArray())
	downloadCallPayment := transactions.TransactionBufferCreateUint32Vector(builder, tx.DownloadCallPayment.toArray())

	fileBytes := []byte(tx.FileName)
	fileName := transactions.TransactionBufferCreateByteVector(builder, fileBytes)
	functionBytes := []byte(tx.FunctionName)
	functionName := transactions.TransactionBufferCreateByteVector(builder, functionBytes)
	argumentBytes := []byte(tx.ActualArguments)
	actualArguments := transactions.TransactionBufferCreateByteVector(builder, argumentBytes)


	mb := make([]flatbuffers.UOffsetT, len(tx.ServicePayments))
	for i, it := range tx.ServicePayments {
		mos := transactions.TransactionBufferCreateUint32Vector(builder, it.toArray())
		transactions.MosaicBufferStart(builder)
		transactions.MosaicBufferAddId(builder, mos)
		mb[i] = transactions.MosaicBufferEnd(builder)
	}
	mV := transactions.TransactionBufferCreateUOffsetVector(builder, mb)

	transactions.ManualCallTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.ManualCallTransactionBufferAddContractKey(builder, contractKeyV)
	transactions.ManualCallTransactionBufferAddFileNameSize(builder, tx.FileNameSize)
	transactions.ManualCallTransactionBufferAddFunctionNameSize(builder, tx.FunctionNameSize)
	transactions.ManualCallTransactionBufferAddActualArgumentsSize(builder, tx.ActualArgumentsSize)
	transactions.ManualCallTransactionBufferAddExecutionCallPayment(builder, executionCallPayment)
	transactions.ManualCallTransactionBufferAddDownloadCallPayment(builder, downloadCallPayment)
	transactions.ManualCallTransactionBufferAddServicePaymentsCount(builder, tx.ServicePaymentsCount)
	transactions.ManualCallTransactionBufferAddFileName(builder, fileName)
	transactions.ManualCallTransactionBufferAddFunctionName(builder,functionName)
	transactions.ManualCallTransactionBufferAddActualArguments(builder, actualArguments) 
	transactions.ManualCallTransactionBufferAddServicePayments(builder, mV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return manualCallTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *ManualCallTransaction) Size() int {
	return ManualCallHeaderSize
}

type manualCallTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ContractKey 				string 			`json:"contractKey"`
		FileNameSize 				uint16 			`json:"fileNameSize"`
		FunctionNameSize 			uint16 			`json:"functionNameSize"`
		ActualArgumentsSize 		uint16 			`json:"actualArgumentsSize"`
		ExecutionCallPayment 		uint64DTO 		`json:"executionCallPayment"`
		DownloadCallPayment 		uint64DTO 		`json:"downloadCallPayment"`
		ServicePaymentsCount 		uint8 			`json:"servicePaymentsCount"`
		FileName 					string 			`json:"fileName"`
		FunctionName 				string 			`json:"functionName"`
		ActualArguments 			string 			`json:"actualArguments"`
		ServicePayments 			[]*mosaicIdDTO	`json:"servicePayments"`
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

	mosaics := make([]*MosaicId, len(dto.Tx.ServicePayments))

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
		dto.Tx.FileNameSize,
		dto.Tx.FunctionNameSize,
		dto.Tx.ActualArgumentsSize,
		dto.Tx.ExecutionCallPayment.toStruct(),
		dto.Tx.DownloadCallPayment.toStruct(),
		dto.Tx.ServicePaymentsCount,
		dto.Tx.FileName,
		dto.Tx.FunctionName,
		dto.Tx.ActualArguments,
		mosaics,
	}, nil
}

// deploy contract transaction
func NewDeployContractTransaction(
	deadline 							*Deadline,
	driveKey 							*PublicAccount,
	fileNameSize 						uint16,
	functionNameSize 					uint16,
	actualArgumentsSize 				uint16,
	executionCallPayment 				Amount,
	downloadCallPayment 				Amount,
	servicePaymentsCount 				uint8,
	automaticExecutionFileNameSize		uint16,
	automaticExecutionFunctionNameSize 	uint16,
	automaticExecutionCallPayment 		Amount,
	automaticDownloadCallPayment 		Amount,
	automaticExecutionsNumber 			uint32,
	assignee							*PublicAccount,
	fileName 							string,
	functionName 						string,
	actualArguments 					string,
	servicePayments 					[]*MosaicId,
	automaticExecutionFileName			string,
	automaticExecutionFunctionName		string,
	networkType 			NetworkType,
) (*DeployContractTransaction, error) {

	tx := DeployContractTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     DeployContractVersion,
			Type:        DeployContract,
			NetworkType: networkType,
		},
		DriveKey: 							driveKey,
		FileNameSize: 						fileNameSize,
		FunctionNameSize: 					functionNameSize,
		ActualArgumentsSize: 				actualArgumentsSize,
		ExecutionCallPayment: 				executionCallPayment,
		DownloadCallPayment: 				downloadCallPayment,
		ServicePaymentsCount: 				servicePaymentsCount,
		AutomaticExecutionFileNameSize:		automaticExecutionFileNameSize,
		AutomaticExecutionFunctionNameSize: automaticExecutionFunctionNameSize,
		AutomaticExecutionCallPayment: 		automaticExecutionCallPayment,
		AutomaticDownloadCallPayment:		automaticDownloadCallPayment,
		AutomaticExecutionsNumber: 			automaticExecutionsNumber,
		Assignee: 							assignee,
		FileName: 							fileName,
		FunctionName: 						functionName,
		ActualArguments: 					actualArguments,
		ServicePayments: 					servicePayments,
		AutomaticExecutionFileName: 		automaticExecutionFileName,
		AutomaticExecutionFunctionName: 	automaticExecutionFunctionName,
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

	executionCallPayment := transactions.TransactionBufferCreateUint32Vector(builder, tx.ExecutionCallPayment.toArray())
	downloadCallPayment := transactions.TransactionBufferCreateUint32Vector(builder, tx.DownloadCallPayment.toArray())

	fileBytes := []byte(tx.FileName)
	fileName := transactions.TransactionBufferCreateByteVector(builder, fileBytes)
	functionBytes := []byte(tx.FunctionName)
	functionName := transactions.TransactionBufferCreateByteVector(builder, functionBytes)
	argumentBytes := []byte(tx.ActualArguments)
	actualArguments := transactions.TransactionBufferCreateByteVector(builder, argumentBytes)
	automaticExecutionFileBytes := []byte(tx.FileName)
	automaticExecutionFileName := transactions.TransactionBufferCreateByteVector(builder, automaticExecutionFileBytes)
	automaticExecutionFunction := []byte(tx.FileName)
	automaticExecutionFunctionName := transactions.TransactionBufferCreateByteVector(builder, automaticExecutionFunction)

	automaticExecutionCallPayment := transactions.TransactionBufferCreateUint32Vector(builder, tx.AutomaticExecutionCallPayment.toArray())
	automaticDownloadCallPayment := transactions.TransactionBufferCreateUint32Vector(builder, tx.AutomaticDownloadCallPayment.toArray())
	mb := make([]flatbuffers.UOffsetT, len(tx.ServicePayments))
	for i, it := range tx.ServicePayments {
		mos := transactions.TransactionBufferCreateUint32Vector(builder, it.toArray())
		transactions.MosaicBufferStart(builder)
		transactions.MosaicBufferAddId(builder, mos)
		mb[i] = transactions.MosaicBufferEnd(builder)
	}
	mV := transactions.TransactionBufferCreateUOffsetVector(builder, mb)

	transactions.DeployContractTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DeployContractTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.DeployContractTransactionBufferAddFileNameSize(builder, tx.FileNameSize)
	transactions.DeployContractTransactionBufferAddFunctionNameSize(builder, tx.FunctionNameSize)
	transactions.DeployContractTransactionBufferAddActualArgumentsSize(builder, tx.ActualArgumentsSize)
	transactions.DeployContractTransactionBufferAddExecutionCallPayment(builder, executionCallPayment)
	transactions.DeployContractTransactionBufferAddDownloadCallPayment(builder, downloadCallPayment)
	transactions.DeployContractTransactionBufferAddServicePaymentsCount(builder, tx.ServicePaymentsCount)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionFileNameSize(builder, tx.AutomaticExecutionFileNameSize)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionFunctionNameSize(builder, tx.AutomaticExecutionFunctionNameSize)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionCallPayment(builder, automaticExecutionCallPayment)
	transactions.DeployContractTransactionBufferAddAutomaticDownloadCallPayment(builder, automaticDownloadCallPayment)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionsNumber(builder, tx.AutomaticExecutionsNumber)
	transactions.DeployContractTransactionBufferAddAssignee(builder, assigneeV)
	transactions.DeployContractTransactionBufferAddFileName(builder, fileName)
	transactions.DeployContractTransactionBufferAddFunctionName(builder,functionName)
	transactions.DeployContractTransactionBufferAddActualArguments(builder, actualArguments) 
	transactions.DeployContractTransactionBufferAddServicePayments(builder, mV)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionFileName(builder, automaticExecutionFileName)
	transactions.DeployContractTransactionBufferAddAutomaticExecutionFunctionName(builder, automaticExecutionFunctionName)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return deployContractTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DeployContractTransaction) Size() int {
	return DeployContractHeaderSize
}

type deployContractTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey 							string 			`json:"driveKey"`
		FileNameSize 						uint16 			`json:"fileNameSize"`
		FunctionNameSize 					uint16 			`json:"functionNameSize"`
		ActualArgumentsSize 				uint16 			`json:"actualArgumentsSize"`
		ExecutionCallPayment 				uint64DTO 		`json:"executionCallPayment"`
		DownloadCallPayment 				uint64DTO 		`json:"downloadCallPayment"`
		ServicePaymentsCount 				uint8 			`json:"servicePaymentsCount"`
		AutomaticExecutionFileNameSize 		uint16 			`json:"automaticExecutionFileNameSize"`
		AutomaticExecutionFunctionNameSize 	uint16 			`json:"automaticExecutionFunctionNameSize"`
		AutomaticExecutionCallPayment 		uint64DTO 		`json:"automaticExecutionCallPayment"`
		AutomaticDownloadCallPayment 		uint64DTO 		`json:"automaticDownloadCallPayment"`
		AutomaticExecutionsNumber 			uint32 			`json:"automaticExecutionsNumber"`
		Assignee 							string 			`json:"assignee"`
		FileName 							string 			`json:"fileName"`
		FunctionName 						string 			`json:"functionName"`
		ActualArguments 					string 			`json:"actualArguments"`
		ServicePayments 					[]*mosaicIdDTO	`json:"servicePayments"`
		AutomaticExecutionFileName 			string 			`json:"automaticExecutionFileName"`
		AutomaticExecutionFunctionName 		string 			`json:"automaticExecutionFunctionName"`
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

	mosaics := make([]*MosaicId, len(dto.Tx.ServicePayments))

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
		dto.Tx.FileNameSize,
		dto.Tx.FunctionNameSize,
		dto.Tx.ActualArgumentsSize,
		dto.Tx.ExecutionCallPayment.toStruct(),
		dto.Tx.DownloadCallPayment.toStruct(),
		dto.Tx.ServicePaymentsCount,
		dto.Tx.AutomaticExecutionFileNameSize,
		dto.Tx.AutomaticExecutionFunctionNameSize,
		dto.Tx.AutomaticExecutionCallPayment.toStruct(),
		dto.Tx.AutomaticDownloadCallPayment.toStruct(),
		dto.Tx.AutomaticExecutionsNumber,
		assignee,
		dto.Tx.FileName,
		dto.Tx.FunctionName,
		dto.Tx.ActualArguments,
		mosaics,
		dto.Tx.AutomaticExecutionFileName,
		dto.Tx.AutomaticExecutionFunctionName,
	}, nil
}