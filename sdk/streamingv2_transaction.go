// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/hex"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewStreamStartTransaction(
	deadline *Deadline,
	driveKey string,
	expectedUploadSize StorageSize,
	folder string,
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
		Folder:             folder,
		FeedbackFeeAmount:  feedbackFeeAmount,
	}

	return &tx, nil
}

func (tx *StreamStartTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *StreamStartTransaction) String() string {
	return ""
}

func (tx *StreamStartTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	bytes := []byte(tx.Folder)
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

	transactions.StreamStartTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.StreamStartTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.StreamStartTransactionBufferAddExpectedUploadSize(builder, uint64(tx.ExpectedUploadSize))
	transactions.StreamStartTransactionBufferAddFolderSize(builder, uint16(tx.FolderSize()))

	transactions.StreamStartTransactionBufferAddFolder(builder, fp)

	transactions.StreamStartTransactionBufferAddFeedbackFeeAmount(builder, uint64(tx.FeedbackFeeAmount))

	t := transactions.StreamStartTransactionBufferEnd(builder)
	builder.Finish(t)

	return streamStartTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *StreamStartTransaction) Size() int {
	return StreamStartHeaderSize + tx.FolderSize()
}

func (tx *StreamStartTransaction) FolderSize() int {
	return len(tx.Folder)
}

type streamStartTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey           string      `json:"driveKey"`
		ExpectedUploadSize StorageSize `json:"blsKey"`
		Folder             string      `json:"folder"`
		FeedbackFeeAmount  Amount      `json:"feedbackFeeAmount"`
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
		dto.Tx.ExpectedUploadSize,
		dto.Tx.Folder,
		dto.Tx.FeedbackFeeAmount,
	}, nil
}
