// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewMosaicAddressRestrictionTransaction(deadline *Deadline, assetId AssetId, restrictionKey uint64, previousRestrictionValue uint64, newRestrictionValue uint64, targetAddress *Address, networkType NetworkType) (*MosaicAddressRestrictionTransaction, error) {

	tx := MosaicAddressRestrictionTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     MosaicAddressRestrictionVersion,
			Deadline:    deadline,
			Type:        AccountAddressRestriction,
			NetworkType: networkType,
		},
		MosaicId:                 assetId,
		RestrictionKey:           restrictionKey,
		PreviousRestrictionValue: previousRestrictionValue,
		NewRestrictionValue:      newRestrictionValue,
		TargetAddress:            targetAddress,
	}

	return &tx, nil
}

func NewMosaicGlobalRestrictionTransaction(deadline *Deadline,
	mosaicId AssetId,
	referenceMosaicId AssetId,
	restrictionKey uint64,
	previousRestrictionValue uint64,
	previousRestrictionType MosaicRestrictionType,
	newRestrictionValue uint64,
	newRestrictionType MosaicRestrictionType, networkType NetworkType) (*MosaicGlobalRestrictionTransaction, error) {

	tx := MosaicGlobalRestrictionTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     MosaicAddressRestrictionVersion,
			Deadline:    deadline,
			Type:        AccountAddressRestriction,
			NetworkType: networkType,
		},
		MosaicId:                 mosaicId,
		ReferenceMosaicId:        referenceMosaicId,
		RestrictionKey:           restrictionKey,
		PreviousRestrictionValue: previousRestrictionValue,
		PreviousRestrictionType:  previousRestrictionType,
		NewRestrictionValue:      newRestrictionValue,
		NewRestrictionType:       newRestrictionType,
	}

	return &tx, nil
}

func (tx *MosaicAddressRestrictionTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}
func (tx *MosaicGlobalRestrictionTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *MosaicGlobalRestrictionTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"MosaicId": %d,
			"ReferenceMosaicId": %d,
			"RestrictionKey": %d,
			"PreviousRestrictionValue": %d,
			"PreviousRestrictionType": %d,
			"NewRestrictionValue": %d,
			"NewRestrictionType": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.MosaicId,
		tx.ReferenceMosaicId,
		tx.RestrictionKey,
		tx.PreviousRestrictionValue,
		tx.PreviousRestrictionType,
		tx.NewRestrictionValue,
		tx.NewRestrictionType,
	)
}

func (tx *MosaicAddressRestrictionTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"MosaicId": %d,
			"RestrictionKey": %d,
			"PreviousRestrictionValue": %d,
			"NewRestrictionValue": %d,
			"Address": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.MosaicId,
		tx.RestrictionKey,
		tx.PreviousRestrictionValue,
		tx.NewRestrictionValue,
		tx.TargetAddress,
	)
}
func (tx *MosaicGlobalRestrictionTransaction) Size() int {
	return MosaicGlobalRestrictionHeaderSize
}
func (tx *MosaicAddressRestrictionTransaction) Size() int {
	return MosaicAddressRestrictionHeaderSize
}

func (tx *MosaicGlobalRestrictionTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	mosaicIdOffset := transactions.TransactionBufferCreateUint32Vector(builder, tx.MosaicId.toArray())
	referenceMosaicIdOffset := transactions.TransactionBufferCreateUint32Vector(builder, tx.ReferenceMosaicId.toArray())
	restrictionKeyOffset := transactions.TransactionBufferCreateUint32Vector(builder, uint64ToArray(tx.RestrictionKey))
	previousRestrictionValueOffset := transactions.TransactionBufferCreateUint32Vector(builder, uint64ToArray(tx.PreviousRestrictionValue))
	newRestrictionValueOffset := transactions.TransactionBufferCreateUint32Vector(builder, uint64ToArray(tx.NewRestrictionValue))

	transactions.MosaicGlobalRestrictionTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.MosaicGlobalRestrictionTransactionBufferAddMosaicId(builder, mosaicIdOffset)
	transactions.MosaicGlobalRestrictionTransactionBufferAddReferenceMosaicId(builder, referenceMosaicIdOffset)
	transactions.MosaicGlobalRestrictionTransactionBufferAddRestrictionKey(builder, restrictionKeyOffset)
	transactions.MosaicGlobalRestrictionTransactionBufferAddPreviousRestrictionValue(builder, previousRestrictionValueOffset)
	transactions.MosaicGlobalRestrictionTransactionBufferAddPreviousRestrictionType(builder, uint8(tx.PreviousRestrictionType))
	transactions.MosaicGlobalRestrictionTransactionBufferAddNewRestrictionValue(builder, newRestrictionValueOffset)
	transactions.MosaicGlobalRestrictionTransactionBufferAddNewRestrictionType(builder, uint8(tx.NewRestrictionType))
	t := transactions.MosaicGlobalRestrictionTransactionBufferEnd(builder)
	builder.Finish(t)

	return mosaicGlobalRestrictionTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *MosaicAddressRestrictionTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	a, err := tx.TargetAddress.Encode()
	if err != nil {
		return nil, err
	}
	mosaicIdOffset := transactions.TransactionBufferCreateUint32Vector(builder, tx.MosaicId.toArray())
	restrictionKeyOffset := transactions.TransactionBufferCreateUint32Vector(builder, uint64ToArray(tx.RestrictionKey))
	previousRestrictionValueOffset := transactions.TransactionBufferCreateUint32Vector(builder, uint64ToArray(tx.PreviousRestrictionValue))
	newRestrictionValueOffset := transactions.TransactionBufferCreateUint32Vector(builder, uint64ToArray(tx.NewRestrictionValue))
	addressOffset := transactions.TransactionBufferCreateByteVector(builder, a)

	transactions.MosaicGlobalRestrictionTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.MosaicAddressRestrictionTransactionBufferAddMosaicId(builder, mosaicIdOffset)
	transactions.MosaicAddressRestrictionTransactionBufferAddRestrictionKey(builder, restrictionKeyOffset)
	transactions.MosaicAddressRestrictionTransactionBufferAddPreviousRestrictionValue(builder, previousRestrictionValueOffset)
	transactions.MosaicAddressRestrictionTransactionBufferAddNewRestrictionValue(builder, newRestrictionValueOffset)
	transactions.MosaicAddressRestrictionTransactionBufferAddTargetAddress(builder, addressOffset)
	t := transactions.MosaicAddressRestrictionTransactionBufferEnd(builder)
	builder.Finish(t)

	return mosaicAddressRestrictionTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type MosaicAddressRestrictionTransactionDto struct {
	Tx struct {
		abstractTransactionDTO
		MosaicId                 mosaicIdDTO
		RestrictionKey           uint64DTO
		PreviousRestrictionValue uint64DTO
		NewRestrictionValue      uint64DTO
		TargetAddress            *Address
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

type MosaicGlobalRestrictionTransactionDto struct {
	Tx struct {
		abstractTransactionDTO
		MosaicId                 mosaicIdDTO
		ReferenceMosaicId        mosaicIdDTO
		RestrictionKey           uint64DTO
		PreviousRestrictionValue uint64DTO
		PreviousRestrictionType  uint8
		NewRestrictionValue      uint64DTO
		NewRestrictionType       uint8
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *MosaicAddressRestrictionTransactionDto) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}
	return &MosaicAddressRestrictionTransaction{
		*atx,
		mosaicId,
		dto.Tx.RestrictionKey.toUint64(),
		dto.Tx.PreviousRestrictionValue.toUint64(),
		dto.Tx.NewRestrictionValue.toUint64(),
		dto.Tx.TargetAddress,
	}, nil
}

func (dto *MosaicGlobalRestrictionTransactionDto) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	mosaicId, err := dto.Tx.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	referenceMosaicId, err := dto.Tx.ReferenceMosaicId.toStruct()
	if err != nil {
		return nil, err
	}
	return &MosaicGlobalRestrictionTransaction{
		*atx,
		mosaicId,
		referenceMosaicId,
		dto.Tx.RestrictionKey.toUint64(),
		dto.Tx.PreviousRestrictionValue.toUint64(),
		MosaicRestrictionType(dto.Tx.PreviousRestrictionType),
		dto.Tx.NewRestrictionValue.toUint64(),
		MosaicRestrictionType(dto.Tx.NewRestrictionType),
	}, nil
}
