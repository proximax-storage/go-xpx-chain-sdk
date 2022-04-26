// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
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
			"Capacity": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Capacity.String(),
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
			"Size": %s,
			"VerificationFeeAmount": %s,
			"ReplicatorCount": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveSize.String(),
		tx.VerificationFeeAmount.String(),
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

func NewDataModificationTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	downloadDataCdi *Hash,
	uploadSize StorageSize,
	feedbackFeeAmount Amount,
	networkType NetworkType,
) (*DataModificationTransaction, error) {
	if driveKey == nil {
		return nil, ErrNilAccount
	}

	if downloadDataCdi == nil {
		return nil, ErrNilHash
	}

	if uploadSize <= 0 {
		return nil, errors.New("uploadSize should be positive")
	}

	if feedbackFeeAmount <= 0 {
		return nil, errors.New("feedbackFeeAmount should be positive")
	}

	tx := DataModificationTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DataModificationVersion,
			Deadline:    deadline,
			Type:        DataModification,
			NetworkType: networkType,
		},
		DriveKey:          driveKey,
		DownloadDataCdi:   downloadDataCdi,
		UploadSize:        uploadSize,
		FeedbackFeeAmount: feedbackFeeAmount,
	}

	return &tx, nil
}

func (tx *DataModificationTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DataModificationTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s
			"DownloadDataCdi": %s
			"UploadSize": %s
			"FeedbackFeeAmount": %s
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
		tx.DownloadDataCdi.String(),
		tx.UploadSize.String(),
		tx.FeedbackFeeAmount.String(),
	)
}

func (tx *DataModificationTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	downloadDataCdiV := hashToBuffer(builder, tx.DownloadDataCdi)
	uploadSizeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.UploadSize.toArray())
	feedbackFeeAmount := transactions.TransactionBufferCreateUint32Vector(builder, tx.FeedbackFeeAmount.toArray())

	transactions.DataModificationTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DataModificationTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.DataModificationTransactionBufferAddDownloadDataCdi(builder, downloadDataCdiV)
	transactions.DataModificationTransactionBufferAddUploadSize(builder, uploadSizeV)
	transactions.DataModificationTransactionBufferAddFeedbackFeeAmount(builder, feedbackFeeAmount)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return dataModificationTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DataModificationTransaction) Size() int {
	return DataModificationHeaderSize
}

type dataModificationTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey          string    `json:"driveKey"`
		DownloadDataCdi   hashDto   `json:"downloadDataCdi"`
		UploadSize        uint64DTO `json:"uploadSize"`
		FeedbackFeeAmount uint64DTO `json:"feedbackFeeAmount"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *dataModificationTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	downloadDataCdi, err := dto.Tx.DownloadDataCdi.Hash()
	if err != nil {
		return nil, err
	}

	return &DataModificationTransaction{
		*atx,
		driveKey,
		downloadDataCdi,
		dto.Tx.UploadSize.toStruct(),
		dto.Tx.FeedbackFeeAmount.toStruct(),
	}, nil
}

func NewDataModificationCancelTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	downloadDataCdi *Hash,
	networkType NetworkType,
) (*DataModificationCancelTransaction, error) {
	if driveKey == nil {
		return nil, ErrNilAccount
	}

	if downloadDataCdi == nil {
		return nil, ErrNilHash
	}

	tx := DataModificationCancelTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DataModificationCancelVersion,
			Deadline:    deadline,
			Type:        DataModificationCancel,
			NetworkType: networkType,
		},
		DriveKey:        driveKey,
		DownloadDataCdi: downloadDataCdi,
	}

	return &tx, nil
}

func (tx *DataModificationCancelTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DataModificationCancelTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s
			"Id": %s
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
		tx.DownloadDataCdi.String(),
	)
}

func (tx *DataModificationCancelTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	downloadDataCdiV := hashToBuffer(builder, tx.DownloadDataCdi)

	transactions.DataModificationCancelTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DataModificationCancelTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.DataModificationCancelTransactionBufferAddDownloadDataCdi(builder, downloadDataCdiV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return dataModificationCancelTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DataModificationCancelTransaction) Size() int {
	return DataModificationCancelHeaderSize
}

type dataModificationCancelTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey        string  `json:"driveKey"`
		DownloadDataCdi hashDto `json:"downloadDataCdi"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *dataModificationCancelTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	downloadDataCdi, err := dto.Tx.DownloadDataCdi.Hash()
	if err != nil {
		return nil, err
	}

	return &DataModificationCancelTransaction{
		*atx,
		driveKey,
		downloadDataCdi,
	}, nil
}

func NewStoragePaymentTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	storageUnits Amount,
	networkType NetworkType,
) (*StoragePaymentTransaction, error) {
	if driveKey == nil {
		return nil, ErrNilAccount
	}

	if storageUnits <= 0 {
		return nil, errors.New("storageUnits should be positive")
	}

	tx := StoragePaymentTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StoragePaymentVersion,
			Deadline:    deadline,
			Type:        StoragePayment,
			NetworkType: networkType,
		},
		DriveKey:     driveKey,
		StorageUnits: storageUnits,
	}

	return &tx, nil
}

func (tx *StoragePaymentTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *StoragePaymentTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s
			"StorageUnits": %s
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
		tx.StorageUnits.String(),
	)
}

func (tx *StoragePaymentTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	storageUnitsV := transactions.TransactionBufferCreateUint32Vector(builder, tx.StorageUnits.toArray())

	transactions.StoragePaymentTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.StoragePaymentTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.StoragePaymentTransactionBufferAddStorageUnits(builder, storageUnitsV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storagePaymentTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *StoragePaymentTransaction) Size() int {
	return StoragePaymentHeaderSize
}

type storagePaymentTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey     string    `json:"driveKey"`
		StorageUnits uint64DTO `json:"storageUnits"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *storagePaymentTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	return &StoragePaymentTransaction{
		*atx,
		driveKey,
		dto.Tx.StorageUnits.toStruct(),
	}, nil
}

func NewDownloadPaymentTransaction(
	deadline *Deadline,
	downloadChannelId *Hash,
	downloadSize StorageSize,
	feedbackFeeAmount Amount,
	networkType NetworkType,
) (*DownloadPaymentTransaction, error) {
	if downloadChannelId == nil {
		return nil, ErrNilHash
	}

	if downloadSize <= 0 {
		return nil, errors.New("downloadSize should be positive")
	}

	if feedbackFeeAmount <= 0 {
		return nil, errors.New("downloadSize should be positive")
	}

	tx := DownloadPaymentTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DownloadPaymentVersion,
			Deadline:    deadline,
			Type:        DownloadPayment,
			NetworkType: networkType,
		},
		DownloadChannelId: downloadChannelId,
		DownloadSize:      downloadSize,
		FeedbackFeeAmount: feedbackFeeAmount,
	}

	return &tx, nil
}

func (tx *DownloadPaymentTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DownloadPaymentTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DownloadChannelId": %s
			"StorageUnits": %s
			"FeedbackFeeAmount": %s
		`,
		tx.AbstractTransaction.String(),
		tx.DownloadChannelId.String(),
		tx.DownloadSize.String(),
		tx.FeedbackFeeAmount.String(),
	)
}

func (tx *DownloadPaymentTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	downloadChannelIdV := hashToBuffer(builder, tx.DownloadChannelId)
	downloadSizeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.DownloadSize.toArray())
	feedbackFeeAmountV := transactions.TransactionBufferCreateUint32Vector(builder, tx.FeedbackFeeAmount.toArray())

	transactions.DownloadPaymentTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DownloadPaymentTransactionBufferAddDownloadChannelId(builder, downloadChannelIdV)
	transactions.DownloadPaymentTransactionBufferAddDownloadSize(builder, downloadSizeV)
	transactions.DownloadPaymentTransactionBufferAddFeedbackFeeAmount(builder, feedbackFeeAmountV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return downloadPaymentTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DownloadPaymentTransaction) Size() int {
	return DownloadPaymentHeaderSize
}

type downloadPaymentTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DownloadChannelId hashDto   `json:"downloadChannelId"`
		DownloadSize      uint64DTO `json:"downloadSize"`
		FeedbackFeeAmount uint64DTO `json:"feedbackFeeAmount"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *downloadPaymentTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	downloadChannelId, err := dto.Tx.DownloadChannelId.Hash()
	if err != nil {
		return nil, err
	}

	return &DownloadPaymentTransaction{
		*atx,
		downloadChannelId,
		dto.Tx.DownloadSize.toStruct(),
		dto.Tx.FeedbackFeeAmount.toStruct(),
	}, nil
}

func NewDownloadTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	downloadSize StorageSize,
	feedbackFeeAmount Amount,
	listOfPublicKeys []*PublicAccount,
	networkType NetworkType,
) (*DownloadTransaction, error) {
	if driveKey == nil {
		return nil, ErrNilAccount
	}

	if downloadSize <= 0 {
		return nil, errors.New("downloadSize should be positive")
	}

	if feedbackFeeAmount <= 0 {
		return nil, errors.New("feedbackFeeAmount should be positive")
	}

	tx := DownloadTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DownloadVersion,
			Deadline:    deadline,
			Type:        Download,
			NetworkType: networkType,
		},
		DriveKey:          driveKey,
		DownloadSize:      downloadSize,
		FeedbackFeeAmount: feedbackFeeAmount,
		ListOfPublicKeys:  listOfPublicKeys,
	}

	return &tx, nil
}

func (tx *DownloadTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DownloadTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s,
			"DownloadSizeBytes": %s,
			"FeedbackFeeAmount": %s,
			"ListOfPublicKeys": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
		tx.DownloadSize.String(),
		tx.FeedbackFeeAmount.String(),
		tx.ListOfPublicKeys,
	)
}

func (tx *DownloadTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	downloadSizeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.DownloadSize.toArray())
	feedbackFeeAmountV := transactions.TransactionBufferCreateUint32Vector(builder, tx.FeedbackFeeAmount.toArray())
	keysV, err := keysToArrayToBuffer(builder, tx.ListOfPublicKeys)
	if err != nil {
		return nil, err
	}

	listOfPublicKeysSizeB := make([]byte, ListOfPublicKeysSize)
	binary.LittleEndian.PutUint16(listOfPublicKeysSizeB, uint16(len(tx.ListOfPublicKeys)))
	listOfPublicKeysSizeV := transactions.TransactionBufferCreateByteVector(builder, listOfPublicKeysSizeB)

	transactions.DownloadTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DownloadTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.DownloadTransactionBufferAddDownloadSize(builder, downloadSizeV)
	transactions.DownloadTransactionBufferAddFeedbackFeeAmount(builder, feedbackFeeAmountV)
	transactions.DownloadTransactionBufferAddListOfPublicKeysSize(builder, listOfPublicKeysSizeV)
	transactions.DownloadTransactionBufferAddListOfPublicKeys(builder, keysV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return downloadTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DownloadTransaction) Size() int {
	return DownloadHeaderSize + KeySize*len(tx.ListOfPublicKeys)
}

type downloadTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey             string    `json:"driveKey"`
		DownloadSize         uint64DTO `json:"downloadSize"`
		FeedbackFeeAmount    uint64DTO `json:"feedbackFeeAmount"`
		ListOfPublicKeysSize uint16    `json:"listOfPublicKeysSize"`
		ListOfPublicKeys     []string  `json:"listOfPublicKeys"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *downloadTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	keys := make([]*PublicAccount, len(dto.Tx.ListOfPublicKeys))
	for i, k := range dto.Tx.ListOfPublicKeys {
		key, err := NewAccountFromPublicKey(k, atx.NetworkType)
		if err != nil {
			return nil, err
		}

		keys[i] = key
	}

	return &DownloadTransaction{
		*atx,
		driveKey,
		dto.Tx.DownloadSize.toStruct(),
		dto.Tx.FeedbackFeeAmount.toStruct(),
		keys,
	}, nil
}

func NewFinishDownloadTransaction(
	deadline *Deadline,
	downloadChannelId *Hash,
	feedbackFeeAmount Amount,
	networkType NetworkType,
) (*FinishDownloadTransaction, error) {
	if downloadChannelId == nil {
		return nil, ErrNilHash
	}

	if feedbackFeeAmount <= 0 {
		return nil, errors.New("downloadSize should be positive")
	}

	tx := FinishDownloadTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     FinishDownloadVersion,
			Deadline:    deadline,
			Type:        FinishDownload,
			NetworkType: networkType,
		},
		DownloadChannelId: downloadChannelId,
		FeedbackFeeAmount: feedbackFeeAmount,
	}

	return &tx, nil
}

func (tx *FinishDownloadTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *FinishDownloadTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DownloadChannelId": %s
			"FeedbackFeeAmount": %s
		`,
		tx.AbstractTransaction.String(),
		tx.DownloadChannelId.String(),
		tx.FeedbackFeeAmount.String(),
	)
}

func (tx *FinishDownloadTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	downloadChannelIdV := hashToBuffer(builder, tx.DownloadChannelId)
	feedbackFeeAmountV := transactions.TransactionBufferCreateUint32Vector(builder, tx.FeedbackFeeAmount.toArray())

	transactions.FinishDownloadTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.FinishDownloadTransactionBufferAddDownloadChannelId(builder, downloadChannelIdV)
	transactions.FinishDownloadTransactionBufferAddFeedbackFeeAmount(builder, feedbackFeeAmountV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return finishDownloadTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *FinishDownloadTransaction) Size() int {
	return FinishDownloadHeaderSize
}

type finishDownloadTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DownloadChannelId hashDto   `json:"downloadChannelId"`
		FeedbackFeeAmount uint64DTO `json:"feedbackFeeAmount"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *finishDownloadTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	downloadChannelId, err := dto.Tx.DownloadChannelId.Hash()
	if err != nil {
		return nil, err
	}

	return &FinishDownloadTransaction{
		*atx,
		downloadChannelId,
		dto.Tx.FeedbackFeeAmount.toStruct(),
	}, nil
}

func NewVerificationPaymentTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	verificationFeeAmount Amount,
	networkType NetworkType,
) (*VerificationPaymentTransaction, error) {
	if driveKey == nil {
		return nil, ErrNilAccount
	}

	if verificationFeeAmount <= 0 {
		return nil, errors.New("verificationFeeAmount should be positive")
	}

	tx := VerificationPaymentTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     VerificationPaymentVersion,
			Deadline:    deadline,
			Type:        VerificationPayment,
			NetworkType: networkType,
		},
		DriveKey:              driveKey,
		VerificationFeeAmount: verificationFeeAmount,
	}

	return &tx, nil
}

func (tx *VerificationPaymentTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *VerificationPaymentTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s
			"CerificationFeeAmount": %s
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
		tx.VerificationFeeAmount.String(),
	)
}

func (tx *VerificationPaymentTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	verificationFeeAmountV := transactions.TransactionBufferCreateUint32Vector(builder, tx.VerificationFeeAmount.toArray())

	transactions.VerificationPaymentTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.VerificationPaymentTransactionBufferAddDriveKey(builder, driveKeyV)
	transactions.VerificationPaymentTransactionBufferAddVerificationFeeAmount(builder, verificationFeeAmountV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return verificationPaymentTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *VerificationPaymentTransaction) Size() int {
	return VerificationPaymentHeaderSize
}

type verificationPaymentTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey              string    `json:"driveKey"`
		VerificationFeeAmount uint64DTO `json:"verificationFeeAmount"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *verificationPaymentTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	return &VerificationPaymentTransaction{
		*atx,
		driveKey,
		dto.Tx.VerificationFeeAmount.toStruct(),
	}, nil
}

func NewEndDriveVerificationTransactionV2(
	deadline *Deadline,
	driveKey *PublicAccount,
	verificationTrigger *Hash,
	shardId uint16,
	keys []*PublicAccount,
	signatures []*Signature,
	opinions []uint8,
	networkType NetworkType,
) (*EndDriveVerificationTransactionV2, error) {

	if driveKey == nil {
		return nil, ErrNilAccount
	}

	if verificationTrigger == nil {
		return nil, ErrNilHash
	}

	tx := EndDriveVerificationTransactionV2{
		AbstractTransaction: AbstractTransaction{
			Version:     EndDriveVerificationV2Version,
			Deadline:    deadline,
			Type:        EndDriveVerificationV2,
			NetworkType: networkType,
		},
		DriveKey:            driveKey,
		VerificationTrigger: verificationTrigger,
		ShardId:             shardId,
		Keys:                keys,
		Signatures:          signatures,
		Opinions:            opinions,
	}

	return &tx, nil
}

func (tx *EndDriveVerificationTransactionV2) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *EndDriveVerificationTransactionV2) String() string {
	keys := ""
	for _, k := range tx.Keys {
		keys += k.String() + "\t"
	}

	signatures := ""
	for _, s := range tx.Signatures {
		signatures += s.String() + "\t"
	}

	return fmt.Sprintf(
		`
		"AbstractTransaction": %s,
		"DriveKey": %s,
		"VerificationTrigger": %s,
		"ShardId": %d,
		"Keys": %s,
		"Provers": %s,
		"Opinions": %+v,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
		tx.VerificationTrigger.String(),
		tx.ShardId,
		keys,
		signatures,
		tx.Opinions,
	)
}

func keysToArrayToBuffer(builder *flatbuffers.Builder, keys []*PublicAccount) (flatbuffers.UOffsetT, error) {
	psB := make([]flatbuffers.UOffsetT, len(keys))
	for i, p := range keys {
		b, err := hex.DecodeString(p.PublicKey)
		if err != nil {
			return 0, err
		}

		pkV := transactions.TransactionBufferCreateByteVector(builder, b)

		transactions.KeysBufferStart(builder)
		transactions.KeysBufferAddKey(builder, pkV)
		psB[i] = transactions.KeysBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, psB), nil
}

func signaturesToArrayToBuffer(builder *flatbuffers.Builder, signatures []*Signature) (flatbuffers.UOffsetT, error) {
	sB := make([]flatbuffers.UOffsetT, len(signatures))
	for i, s := range signatures {
		bsV := transactions.TransactionBufferCreateByteVector(builder, s[:])

		transactions.SignaturesBufferStart(builder)
		transactions.SignaturesBufferAddSignature(builder, bsV)
		sB[i] = transactions.SignaturesBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, sB), nil
}

func opinionsToArrayToBuffer(builder *flatbuffers.Builder, keyCount, signaturesCount int, opinions []uint8) flatbuffers.UOffsetT {
	count := (signaturesCount*keyCount + 7) / 8
	oB := make([]byte, count)

	byteNumber := 0
	for i, o := range opinions {
		oB[byteNumber] |= o << (uint8(i) % 8)

		if i != 0 && i%8 == 0 {
			byteNumber++
		}
	}

	return transactions.TransactionBufferCreateByteVector(builder, oB)
}

func (tx *EndDriveVerificationTransactionV2) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	verificationTriggerV := hashToBuffer(builder, tx.VerificationTrigger)

	shardIdB := make([]byte, ShardIdSize)
	binary.LittleEndian.PutUint16(shardIdB, tx.ShardId)
	shardIdV := transactions.TransactionBufferCreateByteVector(builder, shardIdB)

	keysV, err := keysToArrayToBuffer(builder, tx.Keys)
	if err != nil {
		return nil, err
	}

	signaturesV, err := signaturesToArrayToBuffer(builder, tx.Signatures)
	if err != nil {
		return nil, err
	}

	opinionsV := opinionsToArrayToBuffer(builder, len(tx.Keys), len(tx.Signatures), tx.Opinions)

	transactions.EndDriveVerificationTransactionV2BufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.EndDriveVerificationTransactionV2BufferAddDriveKey(builder, driveKeyV)
	transactions.EndDriveVerificationTransactionV2BufferAddVerificationTrigger(builder, verificationTriggerV)
	transactions.EndDriveVerificationTransactionV2BufferAddShardId(builder, shardIdV)
	transactions.EndDriveVerificationTransactionV2BufferAddKeyCount(builder, uint8(len(tx.Keys)))
	transactions.EndDriveVerificationTransactionV2BufferAddJudgingKeyCount(builder, uint8(len(tx.Signatures)))
	transactions.EndDriveVerificationTransactionV2BufferAddKeys(builder, keysV)
	transactions.EndDriveVerificationTransactionV2BufferAddSignatures(builder, signaturesV)
	transactions.EndDriveVerificationTransactionV2BufferAddOpinions(builder, opinionsV)
	t := transactions.EndDriveVerificationTransactionV2BufferEnd(builder)
	builder.Finish(t)

	return endDriveVerificationV2TransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *EndDriveVerificationTransactionV2) Size() int {
	return EndDriveVerificationV2HeaderSize +
		len(tx.Keys)*KeySize +
		len(tx.Signatures)*SignatureSize +
		(len(tx.Signatures)*len(tx.Keys)+7)/8
}

type endDriveVerificationTransactionV2DTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey            string         `json:"driveKey"`
		VerificationTrigger hashDto        `json:"verificationTrigger"`
		ShardId             uint16         `json:"shardId"`
		Keys                []string       `json:"publicKeys"`
		Signatures          []signatureDto `json:"signatures"`
		Opinions            []uint8        `json:"opinions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *endDriveVerificationTransactionV2DTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	driveAccount, err := NewAccountFromPublicKey(dto.Tx.DriveKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	verificationTrigger, err := dto.Tx.VerificationTrigger.Hash()
	if err != nil {
		return nil, fmt.Errorf("error parsing VerificationTrigger: %v", err)
	}

	keys := make([]*PublicAccount, len(dto.Tx.Keys))
	for i, k := range dto.Tx.Keys {
		keys[i], err = NewAccountFromPublicKey(k, atx.NetworkType)
		if err != nil {
			return nil, err
		}
	}

	signatures := make([]*Signature, len(dto.Tx.Signatures))
	for i, s := range dto.Tx.Signatures {
		signatures[i], err = s.Signature()
		if err != nil {
			return nil, err
		}
	}

	parsedOpinions := parseOpinions(dto.Tx.Opinions, uint8(len(keys)), uint8(len(dto.Tx.Signatures)))

	return &EndDriveVerificationTransactionV2{
		*atx,
		driveAccount,
		verificationTrigger,
		dto.Tx.ShardId,
		keys,
		signatures,
		parsedOpinions,
	}, nil
}

func NewDriveClosureTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	networkType NetworkType,
) (*DriveClosureTransaction, error) {

	if driveKey == nil {
		return nil, ErrNilAccount
	}

	tx := DriveClosureTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DriveClosureVersion,
			Deadline:    deadline,
			Type:        DriveClosure,
			NetworkType: networkType,
		},
		DriveKey: driveKey,
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
			"DriveKey": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
	)
}

func (tx *DriveClosureTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	driveV := transactions.TransactionBufferCreateByteVector(builder, driveB)

	transactions.DriveClosureTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DriveClosureTransactionBufferAddDriveKey(builder, driveV)
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
		DriveKey string `json:"driveKey"`
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

	driveKey, err := NewAccountFromPublicKey(dto.Tx.DriveKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &DriveClosureTransaction{
		*atx,
		driveKey,
	}, nil
}

func NewReplicatorOffboardingTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	networkType NetworkType,
) (*ReplicatorOffboardingTransaction, error) {

	if driveKey == nil {
		return nil, ErrNilAccount
	}

	return &ReplicatorOffboardingTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     ReplicatorOffboardingVersion,
			Type:        ReplicatorOffboarding,
			NetworkType: networkType,
		},
		DriveKey: driveKey,
	}, nil
}

func (tx *ReplicatorOffboardingTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ReplicatorOffboardingTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
	)
}

func (tx *ReplicatorOffboardingTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	driveV := transactions.TransactionBufferCreateByteVector(builder, driveB)

	transactions.ReplicatorOffboardingTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.ReplicatorOffboardingTransactionBufferAddDriveKey(builder, driveV)
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
		DriveKey string `json:"driveKey"`
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

	driveKey, err := NewAccountFromPublicKey(dto.Tx.DriveKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &ReplicatorOffboardingTransaction{
		*atx,
		driveKey,
	}, nil
}

func (tx *DataModificationApprovalTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DataModificationApprovalTransaction) String() string {
	//keys := ""
	//for _, k := range tx.PublicKeys {
	//	keys += k.String() + "\t"
	//}
	//
	//signatures := ""
	//for _, s := range tx.Signatures {
	//	signatures += s.String() + "\t"
	//}

	return fmt.Sprintf(
		`
			"DriveKey": %s, 
			"DataModificationId": %s, 
			"FileStructureCdi": %s, 
			"FileStructureSizeBytes": %d, 
			"MetaFilesSizeBytes": %d, 
			"UsedDriveSizeBytes": %d,
		`,
		tx.DriveKey.String(),
		tx.DataModificationId.String(),
		tx.FileStructureCdi.String(),
		tx.FileStructureSizeBytes,
		tx.MetaFilesSizeBytes,
		tx.UsedDriveSizeBytes,
		//tx.JudgingKeysCount,
		//tx.OverlappingKeysCount,
		//tx.JudgedKeysCount,
		//keys,
		//signatures,
		//tx.PresentOpinions,
		//tx.Opinions,
	)
}

func (tx *DataModificationApprovalTransaction) Bytes() ([]byte, error) {
	return nil, errors.New("cannot get bytes of DataModificationApprovalTransaction")
}

func (tx *DataModificationApprovalTransaction) Size() int {
	return DataModificationApprovalHeaderSize
}

type dataModificationApprovalTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey               string  `json:"driveKey"`
		DataModificationId     hashDto `json:"dataModificationId"`
		FileStructureCdi       hashDto `json:"fileStructureCdi"`
		FileStructureSizeBytes uint64  `json:"fileStructureSizeBytes"`
		MetaFilesSizeBytes     uint64  `json:"metaFilesSizeBytes"`
		UsedDriveSizeBytes     uint64  `json:"usedDriveSizeBytes"`
		//JudgingKeysCount     uint8          `json:"judgingKeysCount"`
		//OverlappingKeysCount uint8          `json:"overlappingKeysCount"`
		//JudgedKeysCount      uint8          `json:"judgedKeysCount"`
		//OpinionElementCount  uint16         `json:"opinionElementCount"`
		//PublicKeys           []string       `json:"publicKeys"`
		//Signatures           []signatureDto `json:"signatures"`
		//PresentOpinions      []uint8        `json:"presentOpinions"`
		//Opinions             []uint64DTO    `json:"opinions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *dataModificationApprovalTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	dataModificationId, err := dto.Tx.DataModificationId.Hash()
	if err != nil {
		return nil, fmt.Errorf("error parsing DownloadChannelId: %v", err)
	}

	fileStructureCdi, err := dto.Tx.FileStructureCdi.Hash()
	if err != nil {
		return nil, fmt.Errorf("error parsing ApprovalTrigger: %v", err)
	}

	//pKeys := make([]*PublicAccount, len(dto.Tx.PublicKeys))
	//for i, k := range dto.Tx.PublicKeys {
	//	pKeys[i], err = NewAccountFromPublicKey(k, atx.NetworkType)
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//
	//signatures := make([]*Signature, len(dto.Tx.Signatures))
	//for i, s := range dto.Tx.Signatures {
	//	signatures[i], err = s.Signature()
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//
	//totalJudgingKeysCount := dto.Tx.JudgingKeysCount + dto.Tx.OverlappingKeysCount
	//presentOpinionByteCount := uint8(len(dto.Tx.PresentOpinions))
	//overlappingKeysCount := (((presentOpinionByteCount * 8) - 7) / totalJudgingKeysCount) - dto.Tx.JudgedKeysCount
	//totalJudgedKeysCount := dto.Tx.JudgedKeysCount + overlappingKeysCount
	//
	//parsedPresentOpinions := parseOpinions(dto.Tx.PresentOpinions, totalJudgedKeysCount, totalJudgingKeysCount)
	//
	//opinions := make([]uint64, len(dto.Tx.Opinions))
	//for i, o := range dto.Tx.Opinions {
	//	opinions[i] = o.toUint64()
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	return &DataModificationApprovalTransaction{
		*atx,
		driveKey,
		dataModificationId,
		fileStructureCdi,
		dto.Tx.FileStructureSizeBytes,
		dto.Tx.MetaFilesSizeBytes,
		dto.Tx.UsedDriveSizeBytes,
		//dto.Tx.JudgingKeysCount,
		//dto.Tx.OverlappingKeysCount,
		//dto.Tx.JudgedKeysCount,
		//dto.Tx.OpinionElementCount,
		//pKeys,
		//signatures,
		//parsedPresentOpinions,
		//opinions,
	}, nil
}

func (tx *DataModificationSingleApprovalTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DataModificationSingleApprovalTransaction) String() string {
	keys := ""
	for _, k := range tx.PublicKeys {
		keys += k.String() + "\t"
	}

	return fmt.Sprintf(
		`
			"DriveKey": %s, 
			"DataModificationId": %s,
			"PublicKeys": %s,
			"Opinions": %+v,
		`,
		tx.DriveKey.String(),
		tx.DataModificationId.String(),
		keys,
		tx.Opinions,
	)
}

func (tx *DataModificationSingleApprovalTransaction) Bytes() ([]byte, error) {
	return nil, errors.New("cannot get bytes of DataModificationApprovalTransaction")
}

func (tx *DataModificationSingleApprovalTransaction) Size() int {
	return DataModificationSignleApprovalHeaderSize +
		len(tx.PublicKeys)*KeySize +
		len(tx.Opinions)*OpinionSizeSize
}

type dataModificationSingleApprovalTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey           string      `json:"driveKey"`
		DataModificationId hashDto     `json:"dataModificationId"`
		PublicKeysCount    uint8       `json:"publicKeysCount"`
		PublicKeys         []string    `json:"publicKeys"`
		Opinions           []uint64DTO `json:"opinions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *dataModificationSingleApprovalTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	dataModificationId, err := dto.Tx.DataModificationId.Hash()
	if err != nil {
		return nil, fmt.Errorf("error parsing DownloadChannelId: %v", err)
	}

	pKeys := make([]*PublicAccount, len(dto.Tx.PublicKeys))
	for i, k := range dto.Tx.PublicKeys {
		pKeys[i], err = NewAccountFromPublicKey(k, atx.NetworkType)
		if err != nil {
			return nil, err
		}
	}

	opinions := make([]uint64, len(dto.Tx.Opinions))
	for i, o := range dto.Tx.Opinions {
		opinions[i] = o.toUint64()
		if err != nil {
			return nil, err
		}
	}

	return &DataModificationSingleApprovalTransaction{
		*atx,
		driveKey,
		dataModificationId,
		dto.Tx.PublicKeysCount,
		pKeys,
		opinions,
	}, nil
}

func (tx *DownloadApprovalTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DownloadApprovalTransaction) String() string {
	keys := ""
	for _, k := range tx.PublicKeys {
		keys += k.String() + "\t"
	}

	signatures := ""
	for _, s := range tx.Signatures {
		signatures += s.String() + "\t"
	}

	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DownloadChannelId": %s,
			"ApprovalTrigger": %s,
			"PublicKeys": %s,
			"Signatures": %s,
			"PresentOpinions": %+v,
			"Opinions": %+v,
		`,
		tx.AbstractTransaction.String(),
		tx.DownloadChannelId.String(),
		tx.ApprovalTrigger.String(),
		keys,
		signatures,
		tx.PresentOpinions,
		tx.Opinions,
	)
}

func presentOpinionsToArrayToBuffer(builder *flatbuffers.Builder, totalJudgedKeysCount, totalJudgingKeysCount int, presentOpinions []uint8) flatbuffers.UOffsetT {
	count := (totalJudgingKeysCount*totalJudgedKeysCount + 7) / 8
	poB := make([]byte, count)

	byteNumber := 0
	for i, o := range presentOpinions {
		poB[byteNumber] |= o << i % 8

		if i != 0 && i%8 == 0 {
			byteNumber++
		}
	}

	return transactions.TransactionBufferCreateByteVector(builder, poB)
}

func opinionsToJaggedArrayToBuffer(builder *flatbuffers.Builder, opinions []*Opinion) (flatbuffers.UOffsetT, error) {
	oB := make([]flatbuffers.UOffsetT, len(opinions))
	for i, op := range opinions {
		oV := transactions.TransactionBufferCreateUint32Vector(builder, op.Opinion[i].toArray())
		transactions.OpinionsBufferStart(builder)
		transactions.OpinionsBufferAddOpinion(builder, oV)
		oB[i] = transactions.OpinionsBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, oB), nil
}

func (tx *DownloadApprovalTransaction) Bytes() ([]byte, error) {
	return nil, errors.New("cannot get bytes of DownloadApprovalTransaction")
}

func (tx *DownloadApprovalTransaction) Size() int {
	return DownloadApprovalHeaderSize +
		len(tx.PublicKeys)*KeySize +
		len(tx.Signatures)*SignatureSize +
		len(tx.PresentOpinions) +
		len(tx.Opinions)*OpinionSizeSize
}

type downloadApprovalTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DownloadChannelId    hashDto        `json:"downloadChannelId"`
		ApprovalTrigger      hashDto        `json:"approvalTrigger"`
		JudgingKeysCount     uint8          `json:"judgingKeysCount"`
		OverlappingKeysCount uint8          `json:"overlappingKeysCount"`
		JudgedKeysCount      uint8          `json:"judgedKeysCount"`
		OpinionElementCount  uint16         `json:"opinionElementCount"`
		PublicKeys           []string       `json:"publicKeys"`
		Signatures           []signatureDto `json:"signatures"`
		PresentOpinions      []uint8        `json:"presentOpinions"`
		Opinions             []uint64DTO    `json:"opinions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *downloadApprovalTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	downloadChannelId, err := dto.Tx.DownloadChannelId.Hash()
	if err != nil {
		return nil, fmt.Errorf("error parsing DownloadChannelId: %v", err)
	}

	approvalTrigger, err := dto.Tx.ApprovalTrigger.Hash()
	if err != nil {
		return nil, fmt.Errorf("error parsing ApprovalTrigger: %v", err)
	}

	pKeys := make([]*PublicAccount, len(dto.Tx.PublicKeys))
	for i, k := range dto.Tx.PublicKeys {
		pKeys[i], err = NewAccountFromPublicKey(k, atx.NetworkType)
		if err != nil {
			return nil, err
		}
	}

	signatures := make([]*Signature, len(dto.Tx.Signatures))
	for i, s := range dto.Tx.Signatures {
		signatures[i], err = s.Signature()
		if err != nil {
			return nil, err
		}
	}

	totalJudgingKeysCount := dto.Tx.JudgingKeysCount + dto.Tx.OverlappingKeysCount
	presentOpinionByteCount := uint8(len(dto.Tx.PresentOpinions))
	overlappingKeysCount := (((presentOpinionByteCount * 8) - 7) / totalJudgingKeysCount) - dto.Tx.JudgedKeysCount
	totalJudgedKeysCount := dto.Tx.JudgedKeysCount + overlappingKeysCount

	parsedPresentOpinions := parseOpinions(dto.Tx.PresentOpinions, totalJudgedKeysCount, totalJudgingKeysCount)

	opinions := make([]uint64, len(dto.Tx.Opinions))
	for i, o := range dto.Tx.Opinions {
		opinions[i] = o.toUint64()
		if err != nil {
			return nil, err
		}
	}

	return &DownloadApprovalTransaction{
		*atx,
		downloadChannelId,
		approvalTrigger,
		dto.Tx.JudgingKeysCount,
		dto.Tx.OverlappingKeysCount,
		dto.Tx.JudgedKeysCount,
		dto.Tx.OpinionElementCount,
		pKeys,
		signatures,
		parsedPresentOpinions,
		opinions,
	}, nil
}

func parseOpinions(opinions []uint8, keyCount, signaturesCount uint8) []uint8 {
	opinionsCount := keyCount * signaturesCount
	parsedOpinions := make([]uint8, opinionsCount)

	byteNumber := 0
	for i := uint8(0); i < opinionsCount; i++ {
		parsedOpinions[i] = opinions[byteNumber] & (i % 8)

		if i != 0 && i%8 == 0 {
			byteNumber++
		}
	}

	return parsedOpinions
}
