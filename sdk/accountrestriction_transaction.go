// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewAccountAddressRestrictionTransaction(deadline *Deadline, flags uint16, additions []*Address, deletions []*Address, networkType NetworkType) (*AccountAddressRestrictionTransaction, error) {

	tx := AccountAddressRestrictionTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     AccountAddressRestrictionVersion,
			Deadline:    deadline,
			Type:        AccountAddressRestriction,
			NetworkType: networkType,
		},
		RestrictionFlags:     flags,
		RestrictionAdditions: additions,
		RestrictionDeletions: deletions,
	}

	return &tx, nil
}

func NewAccountMosaicRestrictionTransaction(deadline *Deadline, flags uint16, additions []AssetId, deletions []AssetId, networkType NetworkType) (*AccountMosaicRestrictionTransaction, error) {

	tx := AccountMosaicRestrictionTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     AccountMosaicRestrictionVersion,
			Deadline:    deadline,
			Type:        AccountMosaicRestriction,
			NetworkType: networkType,
		},
		RestrictionFlags:     flags,
		RestrictionAdditions: additions,
		RestrictionDeletions: deletions,
	}

	return &tx, nil
}

func NewAccountOperationRestrictionTransaction(deadline *Deadline, flags uint16, additions []EntityType, deletions []EntityType, networkType NetworkType) (*AccountOperationRestrictionTransaction, error) {

	tx := AccountOperationRestrictionTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     AccountOperationRestrictionVersion,
			Deadline:    deadline,
			Type:        AccountOperationRestriction,
			NetworkType: networkType,
		},
		RestrictionFlags:     flags,
		RestrictionAdditions: additions,
		RestrictionDeletions: deletions,
	}

	return &tx, nil
}

func (tx *AccountAddressRestrictionTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}
func (tx *AccountMosaicRestrictionTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}
func (tx *AccountOperationRestrictionTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AccountAddressRestrictionTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"RestrictionFlags": %d,
			"RestrictionAdditions": %T,
			"restrictionDeletions": %T,
		`,
		tx.AbstractTransaction.String(),
		tx.RestrictionFlags,
		tx.RestrictionAdditions,
		tx.RestrictionDeletions,
	)
}
func (tx *AccountMosaicRestrictionTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"RestrictionFlags": %d,
			"RestrictionAdditions": %T,
			"restrictionDeletions": %T,
		`,
		tx.AbstractTransaction.String(),
		tx.RestrictionFlags,
		tx.RestrictionAdditions,
		tx.RestrictionDeletions,
	)
}
func (tx *AccountOperationRestrictionTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"RestrictionFlags": %d,
			"RestrictionAdditions": %T,
			"restrictionDeletions": %T,
		`,
		tx.AbstractTransaction.String(),
		tx.RestrictionFlags,
		tx.RestrictionAdditions,
		tx.RestrictionDeletions,
	)
}

func (tx *AccountAddressRestrictionTransaction) Size() int {
	return AccountAddressRestrictionHeaderSize + (AddressSize * (len(tx.RestrictionAdditions) + len(tx.RestrictionDeletions)))
}
func (tx *AccountMosaicRestrictionTransaction) Size() int {
	return AccountMosaicRestrictionHeaderSize + (MosaicIdSize * (len(tx.RestrictionAdditions) + len(tx.RestrictionDeletions)))
}
func (tx *AccountOperationRestrictionTransaction) Size() int {
	return AccountOperationRestrictionHeaderSize + (HalfWordFlagsSize * (len(tx.RestrictionAdditions) + len(tx.RestrictionDeletions)))
}

func (tx *AccountAddressRestrictionTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	additionsVector := make([]flatbuffers.UOffsetT, len(tx.RestrictionAdditions))

	for i, mos := range tx.RestrictionAdditions {
		a, err := mos.Decode()
		if err != nil {
			return nil, err
		}
		addressOffset := transactions.TransactionBufferCreateByteVector(builder, a)
		transactions.AddressBufferStart(builder)
		transactions.AddressBufferAddAddress(builder, addressOffset)
		addressBufferOffset := transactions.AddressBufferEnd(builder)

		additionsVector[i] = addressBufferOffset
	}
	transactions.AccountAddressRestrictionTransactionBufferStartRestrictionAdditionsVector(builder, len(additionsVector))
	for i := len(additionsVector) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(additionsVector[i])
	}
	additionsVectorF := builder.EndVector(len(additionsVector))

	deletionsVector := make([]flatbuffers.UOffsetT, len(tx.RestrictionDeletions))

	for i, mos := range tx.RestrictionDeletions {
		a, err := mos.Decode()
		if err != nil {
			return nil, err
		}
		addressOffset := transactions.TransactionBufferCreateByteVector(builder, a)
		transactions.AddressBufferStart(builder)
		transactions.AddressBufferAddAddress(builder, addressOffset)
		addressBufferOffset := transactions.AddressBufferEnd(builder)

		deletionsVector[i] = addressBufferOffset
	}

	transactions.AccountAddressRestrictionTransactionBufferStartRestrictionDeletionsVector(builder, len(deletionsVector))
	for i := len(deletionsVector) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(deletionsVector[i])
	}
	deletionsVectorF := builder.EndVector(len(deletionsVector))

	transactions.AccountAddressRestrictionTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionFlags(builder, tx.RestrictionFlags)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionAdditionsCount(builder, uint8(len(tx.RestrictionAdditions)))
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionDeletionsCount(builder, uint8(len(tx.RestrictionDeletions)))
	transactions.AccountAddressRestrictionTransactionBufferAddAccountRestrictionTransactionBodyReserved1(builder, 0)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionAdditions(builder, additionsVectorF)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionDeletions(builder, deletionsVectorF)
	t := transactions.AccountAddressRestrictionTransactionBufferEnd(builder)
	builder.Finish(t)

	return accountAddressRestrictionTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AccountMosaicRestrictionTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	additionsVector := make([]flatbuffers.UOffsetT, len(tx.RestrictionAdditions))

	for i, mos := range tx.RestrictionAdditions {
		idOffset := transactions.TransactionBufferCreateUint32Vector(builder, mos.toArray())
		transactions.MosaicIdStart(builder)
		transactions.MosaicIdAddId(builder, idOffset)
		mosaicIdBufferOffset := transactions.MosaicIdEnd(builder)
		additionsVector[i] = mosaicIdBufferOffset
	}
	transactions.AccountMosaicRestrictionTransactionBufferStartRestrictionAdditionsVector(builder, len(additionsVector))
	for i := len(additionsVector) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(additionsVector[i])
	}
	additionsVectorF := builder.EndVector(len(additionsVector))

	deletionsVector := make([]flatbuffers.UOffsetT, len(tx.RestrictionDeletions))

	for i, mos := range tx.RestrictionDeletions {
		idOffset := transactions.TransactionBufferCreateUint32Vector(builder, mos.toArray())
		transactions.MosaicIdStart(builder)
		transactions.MosaicIdAddId(builder, idOffset)
		mosaicIdBufferOffset := transactions.MosaicIdEnd(builder)
		deletionsVector[i] = mosaicIdBufferOffset
	}

	transactions.AccountMosaicRestrictionTransactionBufferStartRestrictionDeletionsVector(builder, len(deletionsVector))
	for i := len(deletionsVector) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(deletionsVector[i])
	}
	deletionsVectorF := builder.EndVector(len(deletionsVector))

	transactions.AccountAddressRestrictionTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionFlags(builder, tx.RestrictionFlags)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionAdditionsCount(builder, uint8(len(tx.RestrictionAdditions)))
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionDeletionsCount(builder, uint8(len(tx.RestrictionDeletions)))
	transactions.AccountAddressRestrictionTransactionBufferAddAccountRestrictionTransactionBodyReserved1(builder, 0)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionAdditions(builder, additionsVectorF)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionDeletions(builder, deletionsVectorF)
	t := transactions.AccountAddressRestrictionTransactionBufferEnd(builder)
	builder.Finish(t)

	return accountMosaicRestrictionTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *AccountOperationRestrictionTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.AccountOperationRestrictionTransactionBufferStartRestrictionAdditionsVector(builder, len(tx.RestrictionAdditions))
	for i, _ := range tx.RestrictionAdditions {
		builder.PrependUint16(uint16(tx.RestrictionAdditions[i]))
	}
	additionsVectorF := builder.EndVector(len(tx.RestrictionAdditions))

	transactions.AccountOperationRestrictionTransactionBufferStartRestrictionDeletionsVector(builder, len(tx.RestrictionDeletions))
	for i, _ := range tx.RestrictionDeletions {
		builder.PrependUint16(uint16(tx.RestrictionDeletions[i]))
	}
	deletionsVectorF := builder.EndVector(len(tx.RestrictionDeletions))

	transactions.AccountAddressRestrictionTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionFlags(builder, tx.RestrictionFlags)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionAdditionsCount(builder, uint8(len(tx.RestrictionAdditions)))
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionDeletionsCount(builder, uint8(len(tx.RestrictionDeletions)))
	transactions.AccountAddressRestrictionTransactionBufferAddAccountRestrictionTransactionBodyReserved1(builder, 0)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionAdditions(builder, additionsVectorF)
	transactions.AccountAddressRestrictionTransactionBufferAddRestrictionDeletions(builder, deletionsVectorF)
	t := transactions.AccountAddressRestrictionTransactionBufferEnd(builder)
	builder.Finish(t)

	return accountOperationRestrictionTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type AccountAddressRestrictionTransactionDto struct {
	Tx struct {
		abstractTransactionDTO
		RestrictionFlags     uint16
		RestrictionAdditions []string
		RestrictionDeletions []string
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

type AccountMosaicRestrictionTransactionDto struct {
	Tx struct {
		abstractTransactionDTO
		RestrictionFlags     uint16
		RestrictionAdditions []*mosaicIdDTO
		RestrictionDeletions []*mosaicIdDTO
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

type AccountOperationRestrictionTransactionDto struct {
	Tx struct {
		abstractTransactionDTO
		RestrictionFlags     uint16
		RestrictionAdditions []*uint16
		RestrictionDeletions []*uint16
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *AccountAddressRestrictionTransactionDto) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	restrictionAdditions := make([]*Address, len(dto.Tx.RestrictionAdditions))

	for i, entry := range dto.Tx.RestrictionAdditions {
		restrictionAdditions[i] = NewAddress(entry, NotSupportedNet)
	}

	restrictionDeletions := make([]*Address, len(dto.Tx.RestrictionAdditions))

	for i, entry := range dto.Tx.RestrictionAdditions {
		restrictionDeletions[i] = NewAddress(entry, NotSupportedNet)
	}

	return &AccountAddressRestrictionTransaction{
		*atx,
		dto.Tx.RestrictionFlags,
		restrictionAdditions,
		restrictionDeletions,
	}, nil
}

func (dto *AccountMosaicRestrictionTransactionDto) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	restrictionAdditions := make([]AssetId, len(dto.Tx.RestrictionAdditions))

	for i, entry := range dto.Tx.RestrictionAdditions {

		assetId, err := entry.toStruct()
		if err != nil {
			return nil, err
		}
		restrictionAdditions[i] = assetId
	}

	restrictionDeletions := make([]AssetId, len(dto.Tx.RestrictionAdditions))

	for i, entry := range dto.Tx.RestrictionAdditions {

		assetId, err := entry.toStruct()
		if err != nil {
			return nil, err
		}
		restrictionDeletions[i] = assetId
	}

	return &AccountMosaicRestrictionTransaction{
		*atx,
		dto.Tx.RestrictionFlags,
		restrictionAdditions,
		restrictionDeletions,
	}, nil
}

func (dto *AccountOperationRestrictionTransactionDto) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	restrictionAdditions := make([]EntityType, len(dto.Tx.RestrictionAdditions))

	for i, entry := range dto.Tx.RestrictionAdditions {

		restrictionAdditions[i] = EntityType(*entry)
	}

	restrictionDeletions := make([]EntityType, len(dto.Tx.RestrictionAdditions))

	for i, entry := range dto.Tx.RestrictionAdditions {

		restrictionDeletions[i] = EntityType(*entry)
	}

	return &AccountOperationRestrictionTransaction{
		*atx,
		dto.Tx.RestrictionFlags,
		restrictionAdditions,
		restrictionDeletions,
	}, nil
}
