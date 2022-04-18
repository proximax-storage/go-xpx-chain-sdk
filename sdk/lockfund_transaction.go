// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewLockFundTransferTransaction(deadline *Deadline, duration Duration, action LockFundAction, mosaics []*Mosaic, networkType NetworkType) (*LockFundTransferTransaction, error) {

	tx := LockFundTransferTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     LockFundTransferVersion,
			Deadline:    deadline,
			Type:        LockFundTransfer,
			NetworkType: networkType,
		},
		Duration: duration,
		Mosaics:  mosaics,
		Action:   action,
	}

	return &tx, nil
}

func (tx *LockFundTransferTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *LockFundTransferTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Duration": %d,
			"Action": %d,
			"Mosaics": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Duration,
		tx.Action,
		tx.Mosaics,
	)
}

func (tx *LockFundTransferTransaction) Size() int {
	return LockFundTransferHeaderSize + ((MosaicIdSize + AmountSize) * len(tx.Mosaics))
}

func (tx *LockFundTransferTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	mb := make([]flatbuffers.UOffsetT, len(tx.Mosaics))
	for i, mos := range tx.Mosaics {
		id := transactions.TransactionBufferCreateUint32Vector(builder, mos.AssetId.toArray())
		am := transactions.TransactionBufferCreateUint32Vector(builder, mos.Amount.toArray())
		transactions.MosaicBufferStart(builder)
		transactions.MosaicBufferAddId(builder, id)
		transactions.MosaicBufferAddAmount(builder, am)
		mb[i] = transactions.MosaicBufferEnd(builder)
	}
	mV := transactions.TransactionBufferCreateUOffsetVector(builder, mb)

	duration := transactions.TransactionBufferCreateUint32Vector(builder, tx.Duration.toArray())
	transactions.LockFundTransferTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.LockFundTransferTransactionBufferAddAction(builder, uint8(tx.Action))
	transactions.LockFundTransferTransactionBufferAddDuration(builder, duration)
	transactions.LockFundTransferTransactionBufferAddMosaicsCount(builder, uint8(len(tx.Mosaics)))
	transactions.LockFundTransferTransactionBufferAddMosaics(builder, mV)
	t := transactions.LockFundTransferTransactionBufferEnd(builder)
	builder.Finish(t)

	return lockFundTransferTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type lockFundTransferTransactionDto struct {
	Tx struct {
		abstractTransactionDTO
		Duration Duration       `json:"duration"`
		Action   LockFundAction `json:"action"`
		Mosaics  []*mosaicDTO   `json:"mosaics"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *lockFundTransferTransactionDto) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	mosaics := make([]*Mosaic, len(dto.Tx.Mosaics))

	for i, mosaic := range dto.Tx.Mosaics {
		msc, err := mosaic.toStruct()
		if err != nil {
			return nil, err
		}

		mosaics[i] = msc
	}

	return &LockFundTransferTransaction{
		*atx,
		dto.Tx.Duration,
		dto.Tx.Action,
		mosaics,
	}, nil
}

func NewLockFundCancelUnlockTransaction(deadline *Deadline, targetHeight Height, networkType NetworkType) (*LockFundCancelUnlockTransaction, error) {

	tx := LockFundCancelUnlockTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     LockFundTransferVersion,
			Deadline:    deadline,
			Type:        LockFundTransfer,
			NetworkType: networkType,
		},
		TargetHeight: targetHeight,
	}

	return &tx, nil
}

func (tx *LockFundCancelUnlockTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *LockFundCancelUnlockTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"TargetHeight": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.TargetHeight,
	)
}

func (tx *LockFundCancelUnlockTransaction) Size() int {
	return LockFundCancelUnlockHeaderSize
}

func (tx *LockFundCancelUnlockTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}
	targetHeight := transactions.TransactionBufferCreateUint32Vector(builder, tx.TargetHeight.toArray())
	transactions.LockFundCancelUnlockTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.LockFundCancelUnlockTransactionBufferAddTargetHeight(builder, targetHeight)
	t := transactions.LockFundCancelUnlockTransactionBufferEnd(builder)
	builder.Finish(t)

	return lockFundCancelUnlockTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type lockFundCancelUnlockTransactionDto struct {
	Tx struct {
		abstractTransactionDTO
		TargetHeight Height `json:"targetHeight"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *lockFundCancelUnlockTransactionDto) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &LockFundCancelUnlockTransaction{
		*atx,
		dto.Tx.TargetHeight,
	}, nil
}
