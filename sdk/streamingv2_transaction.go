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

func NewStreamStartTransaction(
	deadline *Deadline,
	driveKey string,
	expectedUploadSize StorageSize,
	folderName string,
	feedbackFeeAmount Amount,
	networkType NetworkType,
) (*StreamStartTransaction, error) {

	tx := StreamStartTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     StreamStartVersion,
			Type:        StreamStart,
			NetworkType: networkType,
		},
		DriveKey:           driveKey,
		ExpectedUploadSize: expectedUploadSize,
		FolderName:         folderName,
		FeedbackFeeAmount:  feedbackFeeAmount,
	}

	return &tx, nil
}

func (tx *StreamStartTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *StreamStartTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s,
			"ExpectedUploadSize": %d,
			"FolderName": %s,
			"FeedbackFeeAmount", %v,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey,
		tx.ExpectedUploadSize,
		tx.FolderName,
		tx.FeedbackFeeAmount,
	)
}

func (tx *StreamStartTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	bytes := []byte(tx.FolderName)
	fp := transactions.TransactionBufferCreateByteVector(builder, bytes)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	expectedUploadSizeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.ExpectedUploadSize.toArray())
	feedbackFeeAmountV := transactions.TransactionBufferCreateUint32Vector(builder, tx.FeedbackFeeAmount.toArray())

	transactions.StreamStartTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.StreamStartTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.StreamStartTransactionBufferAddExpectedUploadSize(builder, expectedUploadSizeV)
	transactions.StreamStartTransactionBufferAddFolderNameSize(builder, uint16(tx.FolderNameSize()))

	transactions.StreamStartTransactionBufferAddFolderName(builder, fp)

	transactions.StreamStartTransactionBufferAddFeedbackFeeAmount(builder, feedbackFeeAmountV)

	t := transactions.StreamStartTransactionBufferEnd(builder)
	builder.Finish(t)

	return streamStartTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *StreamStartTransaction) Size() int {
	return StreamStartHeaderSize + tx.FolderNameSize()
}

func (tx *StreamStartTransaction) FolderNameSize() int {
	return len(tx.FolderName)
}

type streamStartTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey           string    `json:"driveKey"`
		ExpectedUploadSize uint64DTO `json:"expectedUploadSize"`
		FolderName         string    `json:"folderName"`
		FeedbackFeeAmount  uint64DTO `json:"feedbackFeeAmount"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *streamStartTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &StreamStartTransaction{
		*atx,
		dto.Tx.DriveKey,
		dto.Tx.ExpectedUploadSize.toStruct(),
		dto.Tx.FolderName,
		dto.Tx.FeedbackFeeAmount.toStruct(),
	}, nil
}

func NewStreamFinishTransaction(
	deadline *Deadline,
	driveKey string,
	streamId string,
	actualUploadSize StorageSize,
	streamStructureCdi string,
	networkType NetworkType,
) (*StreamFinishTransaction, error) {

	tx := StreamFinishTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     StreamFinishVersion,
			Type:        StreamFinish,
			NetworkType: networkType,
		},
		DriveKey:           driveKey,
		StreamId:           streamId,
		ActualUploadSize:   actualUploadSize,
		StreamStructureCdi: streamStructureCdi,
	}

	return &tx, nil
}

func (tx *StreamFinishTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *StreamFinishTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s,
			"StreamId": %s
			"ActualUploadSize": %d,
			"StreamStructureCdi": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey,
		tx.StreamId,
		tx.ActualUploadSize,
		tx.StreamStructureCdi,
	)
}

func (tx *StreamFinishTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey)
	if err != nil {
		return nil, err
	}

	streamIdB, err := hex.DecodeString(tx.StreamId)
	if err != nil {
		return nil, err
	}

	streamStructureCdi, err := hex.DecodeString(tx.StreamStructureCdi)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	streamIdV := transactions.TransactionBufferCreateByteVector(builder, streamIdB)
	actualUploadSizeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.ActualUploadSize.toArray())
	streamStructureCdiV := transactions.TransactionBufferCreateByteVector(builder, streamStructureCdi)

	transactions.StreamFinishTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.StreamFinishTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.StreamFinishTransactionBufferAddStreamId(builder, streamIdV)
	transactions.StreamFinishTransactionBufferAddActualUploadSize(builder, actualUploadSizeV)
	transactions.StreamFinishTransactionBufferAddStreamStructureCdi(builder, streamStructureCdiV)

	t := transactions.StreamFinishTransactionBufferEnd(builder)
	builder.Finish(t)

	return streamFinishTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *StreamFinishTransaction) Size() int {
	return StreamFinishHeaderSize
}

type streamFinishTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey           string    `json:"driveKey"`
		StreamId           string    `json:"streamId"`
		ActualUploadSize   uint64DTO `json:"actualUploadSize"`
		StreamStructureCdi string    `json:"streamStructureCdi"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *streamFinishTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &StreamFinishTransaction{
		*atx,
		dto.Tx.DriveKey,
		dto.Tx.StreamId,
		dto.Tx.ActualUploadSize.toStruct(),
		dto.Tx.StreamStructureCdi,
	}, nil
}
