// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/hex"
	"errors"
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewReplicatorOnboardingTransaction(
	deadline *Deadline,
	capacity Amount,
	networkType NetworkType,
) (*ReplicatorOnboardingTransaction, error) {

	if capacity <= 0 {
		return nil, errors.New("capacity should be positive")
	}

	tx := ReplicatorOnboardingTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     ReplicatorOnboardingVersion,
			Type:        ReplicatorOnboarding,
			NetworkType: networkType,
		},
		Capacity: capacity,
	}

	return &tx, nil
}

func (tx *ReplicatorOnboardingTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ReplicatorOnboardingTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Capacity": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.Capacity,
	)
}

func (tx *ReplicatorOnboardingTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	capacityV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Capacity.toArray())

	transactions.ReplicatorOnboardingTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.ReplicatorOnboardingTransactionBufferAddCapacity(builder, capacityV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return replicatorOnboardingTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *ReplicatorOnboardingTransaction) Size() int {
	return ReplicatorOnboardingHeaderSize
}

type replicatorOnboardingTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Capacity uint64DTO `json:"capacity"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *replicatorOnboardingTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &ReplicatorOnboardingTransaction{
		*atx,
		dto.Tx.Capacity.toStruct(),
	}, nil
}

func NewPrepareBcDriveTransaction(
	deadline *Deadline,
	driveSize StorageSize,
	verificationFeeAmount Amount,
	replicatorCount uint16,
	networkType NetworkType,
) (*PrepareBcDriveTransaction, error) {

	if driveSize == 0 {
		return nil, errors.New("driveSize should be positive")
	}

	if verificationFeeAmount == 0 {
		return nil, errors.New("verificationFeeAmount should be positive")
	}

	if replicatorCount == 0 {
		return nil, errors.New("replicatorCount should be positive")
	}

	tx := PrepareBcDriveTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     PrepareBcDriveVersion,
			Deadline:    deadline,
			Type:        PrepareBcDrive,
			NetworkType: networkType,
		},
		DriveSize:             driveSize,
		VerificationFeeAmount: verificationFeeAmount,
		ReplicatorCount:       replicatorCount,
	}

	return &tx, nil
}

func (tx *PrepareBcDriveTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *PrepareBcDriveTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveSize": %d,
			"VerificationFeeAmount": %d,
			"ReplicatorCount": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveSize,
		tx.VerificationFeeAmount,
		tx.ReplicatorCount,
	)
}

func (tx *PrepareBcDriveTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveSizeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.DriveSize.toArray())
	verificationFeeAmountV := transactions.TransactionBufferCreateUint32Vector(builder, tx.VerificationFeeAmount.toArray())

	transactions.PrepareBcDriveTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.PrepareBcDriveTransactionBufferAddDriveSize(builder, driveSizeV)
	transactions.PrepareBcDriveTransactionBufferAddVerificationFeeAmount(builder, verificationFeeAmountV)

	transactions.PrepareBcDriveTransactionBufferAddReplicatorCount(builder, tx.ReplicatorCount)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return prepareBcDriveTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *PrepareBcDriveTransaction) Size() int {
	return PrepareBcDriveHeaderSize
}

type prepareBcDriveTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveSize             uint64DTO `json:"driveSize"`
		VerificationFeeAmount uint64DTO `json:"verificationFeeAmount"`
		ReplicatorCount       uint16    `json:"replicatorCount"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *prepareBcDriveTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &PrepareBcDriveTransaction{
		*atx,
		dto.Tx.DriveSize.toStruct(),
		dto.Tx.VerificationFeeAmount.toStruct(),
		dto.Tx.ReplicatorCount,
	}, nil
}

func NewDriveClosureTransaction(
	deadline *Deadline,
	drive string,
	networkType NetworkType,
) (*DriveClosureTransaction, error) {

	if len(drive) == 0 {
		return nil, ErrNilAccount
	}

	tx := DriveClosureTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DriveClosureVersion,
			Deadline:    deadline,
			Type:        DriveClosure,
			NetworkType: networkType,
		},
		Drive: drive,
	}

	return &tx, nil
}

func (tx *DriveClosureTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DriveClosureTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Drive": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Drive,
	)
}

func (tx *DriveClosureTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveB, err := hex.DecodeString(tx.Drive)
	if err != nil {
		return nil, err
	}

	driveV := transactions.TransactionBufferCreateByteVector(builder, driveB)

	transactions.DriveClosureTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DriveClosureTransactionBufferAddDrive(builder, driveV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return driveClosureTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DriveClosureTransaction) Size() int {
	return DriveClosureHeaderSize
}

type driveClosureTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Drive string `json:"drive"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *driveClosureTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &DriveClosureTransaction{
		*atx,
		dto.Tx.Drive,
	}, nil
}

func NewReplicatorOffboardingTransaction(
	deadline *Deadline,
	networkType NetworkType,
) (*ReplicatorOffboardingTransaction, error) {

	return &ReplicatorOffboardingTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     ReplicatorOffboardingVersion,
			Type:        ReplicatorOffboarding,
			NetworkType: networkType,
		},
	}, nil
}

func (tx *ReplicatorOffboardingTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ReplicatorOffboardingTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
		`,
		tx.AbstractTransaction.String(),
	)
}

func (tx *ReplicatorOffboardingTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.ReplicatorOffboardingTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return replicatorOffboardingTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *ReplicatorOffboardingTransaction) Size() int {
	return ReplicatorOffboardingHeaderSize
}

type replicatorOffboardingTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *replicatorOffboardingTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &ReplicatorOffboardingTransaction{
		*atx,
	}, nil
} 