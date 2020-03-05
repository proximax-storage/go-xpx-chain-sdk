// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewPrepareDriveTransaction(
	deadline *Deadline,
	owner *PublicAccount,
	duration Duration,
	billingPeriod Duration,
	billingPrice Amount,
	driveSize StorageSize,
	replicas uint16,
	minReplicators uint16,
	percentApprovers uint8,
	networkType NetworkType,
) (*PrepareDriveTransaction, error) {

	if owner == nil {
		return nil, ErrNilAccount
	}

	if duration == 0 {
		return nil, errors.New("duration should be positive")
	}

	if billingPeriod == 0 {
		return nil, errors.New("billingPeriod should be positive")
	}

	if (duration % billingPeriod) != 0 {
		return nil, errors.New("billingPeriod should be multiples of duration")
	}

	if billingPrice == 0 {
		return nil, errors.New("billingPrice should be positive")
	}

	if driveSize == 0 {
		return nil, errors.New("driveSize should be positive")
	}

	if replicas == 0 {
		return nil, errors.New("replicas should be positive")
	}

	if minReplicators == 0 {
		return nil, errors.New("minReplicators should be positive")
	}

	if percentApprovers == 0 || percentApprovers > 100 {
		return nil, errors.New("percentApprovers should be in range 1-100")
	}

	mctx := PrepareDriveTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     PrepareDriveVersion,
			Deadline:    deadline,
			Type:        PrepareDrive,
			NetworkType: networkType,
		},
		Owner:            owner,
		Duration:         duration,
		BillingPeriod:    billingPeriod,
		BillingPrice:     billingPrice,
		DriveSize:        driveSize,
		Replicas:         replicas,
		MinReplicators:   minReplicators,
		PercentApprovers: percentApprovers,
	}

	return &mctx, nil
}

func (tx *PrepareDriveTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *PrepareDriveTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Owner": %s,
			"Duration": %d,
			"BillingPeriod": %d,
			"BillingPrice": %d,
			"DriveSize": %d,
			"Replicas": %d,
			"MinReplicators": %d,
			"PercentApprovers": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.Owner,
		tx.Duration,
		tx.BillingPeriod,
		tx.BillingPrice,
		tx.DriveSize,
		tx.Replicas,
		tx.MinReplicators,
		tx.PercentApprovers,
	)
}

func (tx *PrepareDriveTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	ownerB, err := hex.DecodeString(tx.Owner.PublicKey)
	if err != nil {
		return nil, err
	}

	ownerV := transactions.TransactionBufferCreateByteVector(builder, ownerB)
	durationV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Duration.toArray())
	billingPeriodV := transactions.TransactionBufferCreateUint32Vector(builder, tx.BillingPeriod.toArray())
	billingPriceV := transactions.TransactionBufferCreateUint32Vector(builder, tx.BillingPrice.toArray())
	driveSizeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.DriveSize.toArray())

	transactions.PrepareDriveTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.PrepareDriveTransactionBufferAddOwner(builder, ownerV)
	transactions.PrepareDriveTransactionBufferAddDuration(builder, durationV)
	transactions.PrepareDriveTransactionBufferAddBillingPeriod(builder, billingPeriodV)
	transactions.PrepareDriveTransactionBufferAddBillingPrice(builder, billingPriceV)
	transactions.PrepareDriveTransactionBufferAddDriveSize(builder, driveSizeV)

	transactions.PrepareDriveTransactionBufferAddReplicas(builder, tx.Replicas)
	transactions.PrepareDriveTransactionBufferAddMinReplicators(builder, tx.MinReplicators)
	transactions.PrepareDriveTransactionBufferAddPercentApprovers(builder, tx.PercentApprovers)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return prepareDriveTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *PrepareDriveTransaction) Size() int {
	return PrepareDriveHeaderSize
}

type prepareDriveTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Owner            string    `json:"owner"`
		Duration         uint64DTO `json:"duration"`
		BillingPeriod    uint64DTO `json:"billingPeriod"`
		BillingPrice     uint64DTO `json:"billingPrice"`
		DriveSize        uint64DTO `json:"driveSize"`
		Replicas         uint16    `json:"replicas"`
		MinReplicators   uint16    `json:"minReplicators"`
		PercentApprovers uint8     `json:"percentApprovers"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *prepareDriveTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	owner, err := NewAccountFromPublicKey(dto.Tx.Owner, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &PrepareDriveTransaction{
		*atx,
		owner,
		dto.Tx.Duration.toStruct(),
		dto.Tx.BillingPeriod.toStruct(),
		dto.Tx.BillingPrice.toStruct(),
		dto.Tx.DriveSize.toStruct(),
		dto.Tx.Replicas,
		dto.Tx.MinReplicators,
		dto.Tx.PercentApprovers,
	}, nil
}

func NewJoinToDriveTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	networkType NetworkType,
) (*JoinToDriveTransaction, error) {

	if driveKey == nil {
		return nil, ErrNilAccount
	}

	tx := JoinToDriveTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     JoinToDriveVersion,
			Deadline:    deadline,
			Type:        JoinToDrive,
			NetworkType: networkType,
		},
		DriveKey: driveKey,
	}

	return &tx, nil
}

func (tx *JoinToDriveTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *JoinToDriveTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey,
	)
}

func (tx *JoinToDriveTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	hV := transactions.TransactionBufferCreateByteVector(builder, b)

	transactions.JoinToDriveTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.JoinToDriveTransactionBufferAddDriveKey(builder, hV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return joinDriveTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *JoinToDriveTransaction) Size() int {
	return JoinToDriveHeaderSize
}

type joinToDriveTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey string `json:"driveKey"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *joinToDriveTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	acc, err := NewAccountFromPublicKey(dto.Tx.DriveKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &JoinToDriveTransaction{
		*atx,
		acc,
	}, nil
}

func NewDriveFileSystemTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	newRootHash *Hash,
	oldRootHash *Hash,
	addActions []*Action,
	removeActions []*Action,
	networkType NetworkType,
) (*DriveFileSystemTransaction, error) {

	if driveKey == nil {
		return nil, ErrNilAccount
	}

	if newRootHash == nil || oldRootHash == nil {
		return nil, errors.New("rootHash should not be nil")
	}

	if newRootHash.Equal(oldRootHash) {
		return nil, ErrNoChanges
	}

	tx := DriveFileSystemTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DriveFileSystemVersion,
			Deadline:    deadline,
			Type:        DriveFileSystem,
			NetworkType: networkType,
		},
		DriveKey:      driveKey,
		NewRootHash:   newRootHash,
		OldRootHash:   oldRootHash,
		AddActions:    addActions,
		RemoveActions: removeActions,
	}

	return &tx, nil
}

func (tx *DriveFileSystemTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DriveFileSystemTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s,
			"NewRootHash": %s,
			"OldRootHash": %s,
			"AddActions": %s,
			"RemoveActions": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey,
		tx.NewRootHash,
		tx.OldRootHash,
		tx.AddActions,
		tx.RemoveActions,
	)
}

func actionsToArrayToBuffer(builder *flatbuffers.Builder, addActions []*Action) (flatbuffers.UOffsetT, error) {
	msb := make([]flatbuffers.UOffsetT, len(addActions))
	for i, m := range addActions {

		rhV := transactions.TransactionBufferCreateByteVector(builder, m.FileHash[:])
		sizeDV := transactions.TransactionBufferCreateUint32Vector(builder, m.FileSize.toArray())
		transactions.AddActionBufferStart(builder)
		transactions.AddActionBufferAddFileHash(builder, rhV)
		transactions.AddActionBufferAddFileSize(builder, sizeDV)
		msb[i] = transactions.TransactionBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb), nil
}

func (tx *DriveFileSystemTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	driveV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)
	rhV := transactions.TransactionBufferCreateByteVector(builder, tx.NewRootHash[:])

	xorRootHash := tx.NewRootHash.Xor(tx.OldRootHash)
	xhV := transactions.TransactionBufferCreateByteVector(builder, xorRootHash[:])

	addActionsV, err := actionsToArrayToBuffer(builder, tx.AddActions)
	if err != nil {
		return nil, err
	}

	removeActionsV, err := actionsToArrayToBuffer(builder, tx.RemoveActions)
	if err != nil {
		return nil, err
	}

	addActionsCountB := make([]byte, AddActionsSize)
	binary.LittleEndian.PutUint16(addActionsCountB, uint16(len(tx.AddActions)))
	addActionsCountV := transactions.TransactionBufferCreateByteVector(builder, addActionsCountB)

	removeActionsCountB := make([]byte, RemoveActionsSize)
	binary.LittleEndian.PutUint16(removeActionsCountB, uint16(len(tx.RemoveActions)))
	removeActionsCountV := transactions.TransactionBufferCreateByteVector(builder, removeActionsCountB)

	transactions.DriveFileSystemTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DriveFileSystemTransactionBufferAddDriveKey(builder, driveV)
	transactions.DriveFileSystemTransactionBufferAddRootHash(builder, rhV)
	transactions.DriveFileSystemTransactionBufferAddXorRootHash(builder, xhV)

	transactions.DriveFileSystemTransactionBufferAddAddActionsCount(builder, addActionsCountV)
	transactions.DriveFileSystemTransactionBufferAddRemoveActionsCount(builder, removeActionsCountV)

	transactions.DriveFileSystemTransactionBufferAddAddActions(builder, addActionsV)
	transactions.DriveFileSystemTransactionBufferAddRemoveActions(builder, removeActionsV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return driveFileSystemTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DriveFileSystemTransaction) Size() int {
	return DriveFileSystemHeaderSize + (len(tx.AddActions)+len(tx.RemoveActions))*(Hash256+StorageSizeSize)
}

type driveFileSystemAddActionDTO struct {
	FileHash hashDto   `json:"fileHash"`
	FileSize uint64DTO `json:"fileSize"`
}

type driveFileSystemTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey           string                         `json:"driveKey"`
		RootHash           hashDto                        `json:"rootHash"`
		XorRootHash        hashDto                        `json:"xorRootHash"`
		AddActionsCount    uint16                         `json:"addActionsCount"`
		RemoveActionsCount uint16                         `json:"removeActionsCount"`
		AddActions         []*driveFileSystemAddActionDTO `json:"addActions"`
		RemoveActions      []*driveFileSystemAddActionDTO `json:"removeActions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *driveFileSystemTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	rHash, err := dto.Tx.RootHash.Hash()
	if err != nil {
		return nil, err
	}

	xorRootHash, err := dto.Tx.XorRootHash.Hash()
	if err != nil {
		return nil, err
	}

	addActs, err := actionsDTOArrayToStruct(dto.Tx.AddActions)
	if err != nil {
		return nil, err
	}

	removeActs, err := actionsDTOArrayToStruct(dto.Tx.RemoveActions)
	if err != nil {
		return nil, err
	}

	return &DriveFileSystemTransaction{
		*atx,
		driveKey,
		rHash,
		xorRootHash.Xor(rHash),
		addActs,
		removeActs,
	}, nil
}

func actionsDTOArrayToStruct(actions []*driveFileSystemAddActionDTO) ([]*Action, error) {
	acts := make([]*Action, len(actions))
	var err error = nil
	for i, m := range actions {
		h, err := m.FileHash.Hash()
		if err != nil {
			return nil, err
		}

		s := m.FileSize.toUint64()

		acts[i] = &Action{
			FileHash: h,
			FileSize: StorageSize(s),
		}

	}

	return acts, err
}

func NewFilesDepositTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	files []*File,
	networkType NetworkType,
) (*FilesDepositTransaction, error) {

	if driveKey == nil {
		return nil, ErrNilAccount
	}

	if len(files) == 0 {
		return nil, ErrNoChanges
	}

	tx := FilesDepositTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     FilesDepositVersion,
			Deadline:    deadline,
			Type:        FilesDeposit,
			NetworkType: networkType,
		},
		DriveKey: driveKey,
		Files:    files,
	}

	return &tx, nil
}

func (tx *FilesDepositTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *FilesDepositTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s,
			"Files": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey,
		tx.Files,
	)
}

func fileToArrayToBuffer(builder *flatbuffers.Builder, addActions []*File) (flatbuffers.UOffsetT, error) {
	msb := make([]flatbuffers.UOffsetT, len(addActions))
	for i, m := range addActions {

		rhV := transactions.TransactionBufferCreateByteVector(builder, m.FileHash[:])
		transactions.FileBufferStart(builder)
		transactions.AddActionBufferAddFileHash(builder, rhV)
		msb[i] = transactions.TransactionBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb), nil
}

func (tx *FilesDepositTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	hV := transactions.TransactionBufferCreateByteVector(builder, b)

	flsV, err := fileToArrayToBuffer(builder, tx.Files)
	if err != nil {
		return nil, err
	}

	transactions.DriveFileSystemTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.FilesDepositTransactionBufferAddDriveKey(builder, hV)

	transactions.FilesDepositTransactionBufferAddFilesCount(builder, uint16(len(tx.Files)))
	transactions.FilesDepositTransactionBufferAddFiles(builder, flsV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return filesDepositTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *FilesDepositTransaction) Size() int {
	return FilesDepositHeaderSize + len(tx.Files)*Hash256
}

type fileDepositDTO struct {
	FileHash hashDto `json:"fileHash"`
}

type filesDepositTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey   string            `json:"driveKey"`
		FilesCount uint16            `json:"filesCount"`
		Files      []*fileDepositDTO `json:"files"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *filesDepositTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	fls, err := filesDTOArrayToStruct(dto.Tx.Files)
	if err != nil {
		return nil, err
	}

	acc, err := NewAccountFromPublicKey(dto.Tx.DriveKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &FilesDepositTransaction{
		*atx,
		acc,
		fls,
	}, nil
}

func filesDTOArrayToStruct(files []*fileDepositDTO) ([]*File, error) {
	filesResult := make([]*File, len(files))
	var err error = nil
	for i, m := range files {
		h, err := m.FileHash.Hash()
		if err != nil {
			return nil, err
		}
		filesResult[i] = &File{
			FileHash: h,
		}

	}

	return filesResult, err
}

func NewEndDriveTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	networkType NetworkType,
) (*EndDriveTransaction, error) {

	if driveKey == nil {
		return nil, ErrNilAccount
	}

	tx := EndDriveTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     EndDriveVersion,
			Deadline:    deadline,
			Type:        EndDrive,
			NetworkType: networkType,
		},
		DriveKey: driveKey,
	}

	return &tx, nil
}

func (tx *EndDriveTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *EndDriveTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey,
	)
}

func (tx *EndDriveTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	hV := transactions.TransactionBufferCreateByteVector(builder, b)

	transactions.EndDriveTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.EndDriveTransactionBufferAddDriveKey(builder, hV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return endDriveTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *EndDriveTransaction) Size() int {
	return EndDriveHeaderSize
}

type endDriveTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey string `json:"driveKey"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *endDriveTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	return &EndDriveTransaction{
		*atx,
		driveKey,
	}, nil
}

func NewDriveFilesRewardTransaction(
	deadline *Deadline,
	infos []*UploadInfo,
	networkType NetworkType,
) (*DriveFilesRewardTransaction, error) {

	if len(infos) == 0 {
		return nil, ErrNoChanges
	}

	tx := DriveFilesRewardTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DriveFilesRewardVersion,
			Deadline:    deadline,
			Type:        DriveFilesReward,
			NetworkType: networkType,
		},
		UploadInfos: infos,
	}

	return &tx, nil
}

func (tx *DriveFilesRewardTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DriveFilesRewardTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"UploadInfos": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.UploadInfos,
	)
}

func uploadInfosToArrayToBuffer(builder *flatbuffers.Builder, infos []*UploadInfo) (flatbuffers.UOffsetT, error) {
	infosb := make([]flatbuffers.UOffsetT, len(infos))
	for j, info := range infos {
		rb, err := hex.DecodeString(info.Participant.PublicKey)
		if err != nil {
			return 0, err
		}

		rV := transactions.TransactionBufferCreateByteVector(builder, rb)
		uV := transactions.TransactionBufferCreateUint32Vector(builder, info.UploadedSize.toArray())

		transactions.UploadInfoBufferStart(builder)
		transactions.UploadInfoBufferAddReplicator(builder, rV)
		transactions.UploadInfoBufferAddUploaded(builder, uV)
		infosb[j] = transactions.UploadInfoBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, infosb), nil
}

func (tx *DriveFilesRewardTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	infosV, err := uploadInfosToArrayToBuffer(builder, tx.UploadInfos)
	if err != nil {
		return nil, err
	}

	transactions.DriveFilesRewardTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.DriveFilesRewardTransactionBufferAddUploadInfosCount(builder, uint16(len(tx.UploadInfos)))
	transactions.DriveFilesRewardTransactionBufferAddUploadInfos(builder, infosV)
	t := transactions.DriveFilesRewardTransactionBufferEnd(builder)
	builder.Finish(t)

	return driveFilesRewardTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DriveFilesRewardTransaction) Size() int {
	return TransactionHeaderSize + 2 + len(tx.UploadInfos)*(Hash256+StorageSizeSize)
}

type uploadInfoDTO struct {
	Participant string    `json:"participant"`
	Uploaded    uint64DTO `json:"uploaded"`
}

func (dto *uploadInfoDTO) toStruct(networkType NetworkType) (*UploadInfo, error) {
	acc, err := NewAccountFromPublicKey(dto.Participant, networkType)
	if err != nil {
		return nil, err
	}

	return &UploadInfo{
		acc,
		dto.Uploaded.toStruct(),
	}, nil
}

type driveFilesRewardTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		UploadInfos []*uploadInfoDTO `json:"uploadInfos"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *driveFilesRewardTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	uploadInfos := make([]*UploadInfo, len(dto.Tx.UploadInfos))

	for i, u := range dto.Tx.UploadInfos {
		info, err := u.toStruct(atx.NetworkType)
		if err != nil {
			return nil, err
		}

		uploadInfos[i] = info
	}

	return &DriveFilesRewardTransaction{
		*atx,
		uploadInfos,
	}, nil
}

func NewStartDriveVerificationTransaction(
	deadline *Deadline,
	driveKey *PublicAccount,
	networkType NetworkType,
) (*StartDriveVerificationTransaction, error) {

	if driveKey == nil {
		return nil, ErrNilAccount
	}

	tx := StartDriveVerificationTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StartDriveVerificationVersion,
			Deadline:    deadline,
			Type:        StartDriveVerification,
			NetworkType: networkType,
		},
		DriveKey: driveKey,
	}

	return &tx, nil
}

func (tx *StartDriveVerificationTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *StartDriveVerificationTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey,
	)
}

func (tx *StartDriveVerificationTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	hV := transactions.TransactionBufferCreateByteVector(builder, b)

	transactions.StartDriveVerificationTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.StartDriveVerificationTransactionBufferAddDriveKey(builder, hV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return startDriveVerificationTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *StartDriveVerificationTransaction) Size() int {
	return StartDriveVerificationHeaderSize
}

type startDriveVerificationTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey string `json:"driveKey"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *startDriveVerificationTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	return &StartDriveVerificationTransaction{
		*atx,
		driveKey,
	}, nil
}

func NewEndDriveVerificationTransaction(
	deadline *Deadline,
	failures []*FailureVerification,
	networkType NetworkType,
) (*EndDriveVerificationTransaction, error) {

	tx := EndDriveVerificationTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     EndDriveVerificationVersion,
			Deadline:    deadline,
			Type:        EndDriveVerification,
			NetworkType: networkType,
		},
		Failures: failures,
	}

	return &tx, nil
}

func (tx *EndDriveVerificationTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *EndDriveVerificationTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Failures": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Failures,
	)
}

func failureVerificationsToArrayToBuffer(builder *flatbuffers.Builder, failures []*FailureVerification) (flatbuffers.UOffsetT, error) {
	failuresb := make([]flatbuffers.UOffsetT, len(failures))
	for i, f := range failures {
		rb, err := hex.DecodeString(f.Replicator.PublicKey)
		if err != nil {
			return 0, err
		}

		blockHashesb := make([]flatbuffers.UOffsetT, len(f.BlochHashes))

		for i, block := range f.BlochHashes {
			hV := hashToBuffer(builder, block)
			transactions.BlockHashBufferStart(builder)
			transactions.BlockHashBufferAddBlockHashe(builder, hV)
			blockHashesb[i] = transactions.BlockHashBufferEnd(builder)
		}

		rV := transactions.TransactionBufferCreateByteVector(builder, rb)
		hV := transactions.TransactionBufferCreateUOffsetVector(builder, blockHashesb)

		transactions.VerificationFailureBufferStart(builder)
		transactions.VerificationFailureBufferAddSize(builder, uint32(f.Size()))
		transactions.VerificationFailureBufferAddReplicator(builder, rV)
		transactions.VerificationFailureBufferAddBlockHashes(builder, hV)
		failuresb[i] = transactions.VerificationFailureBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, failuresb), nil
}

func (tx *EndDriveVerificationTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	failuresV, err := failureVerificationsToArrayToBuffer(builder, tx.Failures)
	if err != nil {
		return nil, err
	}

	transactions.EndDriveVerificationTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.EndDriveVerificationTransactionBufferAddFailures(builder, failuresV)
	t := transactions.EndDriveVerificationTransactionBufferEnd(builder)
	builder.Finish(t)

	return endDriveVerificationTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *EndDriveVerificationTransaction) Size() int {
	size := 0
	for _, f := range tx.Failures {
		size += f.Size()
	}

	return TransactionHeaderSize + size
}

type failureVerificationDTO struct {
	Replicator  string    `json:"replicator"`
	BlockHashes []hashDto `json:"blockHashes"`
}

func (dto *failureVerificationDTO) toStruct(networkType NetworkType) (*FailureVerification, error) {
	acc, err := NewAccountFromPublicKey(dto.Replicator, networkType)
	if err != nil {
		return nil, err
	}

	hashes := make([]*Hash, len(dto.BlockHashes))

	for i, h := range dto.BlockHashes {
		hash, err := h.Hash()
		if err != nil {
			return nil, err
		}
		hashes[i] = hash
	}

	return &FailureVerification{
		acc,
		hashes,
	}, nil
}

type endDriveVerificationTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Failures []*failureVerificationDTO `json:"verificationFailures"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *endDriveVerificationTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	failures := make([]*FailureVerification, len(dto.Tx.Failures))

	for i, f := range dto.Tx.Failures {
		failure, err := f.toStruct(atx.NetworkType)
		if err != nil {
			return nil, err
		}

		failures[i] = failure
	}

	return &EndDriveVerificationTransaction{
		*atx,
		failures,
	}, nil
}

func NewStartFileDownloadTransaction(
	deadline *Deadline,
	drive *PublicAccount,
	files []*DownloadFile,
	networkType NetworkType,
) (*StartFileDownloadTransaction, error) {

	if drive == nil {
		return nil, ErrNilAccount
	}

	if len(files) == 0 {
		return nil, ErrNoChanges
	}

	tx := StartFileDownloadTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StartFileDownloadVersion,
			Deadline:    deadline,
			Type:        StartFileDownload,
			NetworkType: networkType,
		},
		Drive: drive,
		Files: files,
	}

	return &tx, nil
}

func (tx *StartFileDownloadTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *StartFileDownloadTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Drive": %s,
			"Files": %+v,
		`,
		tx.AbstractTransaction.String(),
		tx.Drive,
		tx.Files,
	)
}

func downloadFilesToArrayToBuffer(builder *flatbuffers.Builder, files []*DownloadFile) (flatbuffers.UOffsetT, error) {
	filesb := make([]flatbuffers.UOffsetT, len(files))
	for i, f := range files {
		hV := hashToBuffer(builder, f.FileHash)
		sizeV := transactions.TransactionBufferCreateUint32Vector(builder, f.FileSize.toArray())

		transactions.AddActionBufferStart(builder)
		transactions.AddActionBufferAddFileSize(builder, sizeV)
		transactions.AddActionBufferAddFileHash(builder, hV)
		filesb[i] = transactions.AddActionBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, filesb), nil
}

func (tx *StartFileDownloadTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(tx.Drive.PublicKey)
	if err != nil {
		return nil, err
	}

	driveV := transactions.TransactionBufferCreateByteVector(builder, b)

	filesV, err := downloadFilesToArrayToBuffer(builder, tx.Files)
	if err != nil {
		return nil, err
	}

	transactions.StartFileDownloadTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.StartFileDownloadTransactionBufferAddDriveKey(builder, driveV)
	transactions.StartFileDownloadTransactionBufferAddFileCount(builder, uint16(len(tx.Files)))
	transactions.StartFileDownloadTransactionBufferAddFiles(builder, filesV)
	t := transactions.StartFileDownloadTransactionBufferEnd(builder)
	builder.Finish(t)

	return startFileDownloadTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *StartFileDownloadTransaction) Size() int {
	return StartFileDownloadHeaderSize + len(tx.Files)*(StorageSizeSize+Hash256)
}

type downloadFileDTO struct {
	Hash hashDto   `json:"fileHash"`
	Size uint64DTO `json:"fileSize"`
}

func (dto *downloadFileDTO) toStruct() (*DownloadFile, error) {
	hash, err := dto.Hash.Hash()
	if err != nil {
		return nil, err
	}

	return &DownloadFile{
		hash,
		dto.Size.toStruct(),
	}, nil
}

type startFileDownloadTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey string             `json:"driveKey"`
		Files    []*downloadFileDTO `json:"files"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *startFileDownloadTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	acc, err := NewAccountFromPublicKey(dto.Tx.DriveKey, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	files := make([]*DownloadFile, len(dto.Tx.Files))

	for i, f := range dto.Tx.Files {
		file, err := f.toStruct()
		if err != nil {
			return nil, err
		}

		files[i] = file
	}

	return &StartFileDownloadTransaction{
		*atx,
		acc,
		files,
	}, nil
}

func NewEndFileDownloadTransaction(
	deadline *Deadline,
	recipient *PublicAccount,
	operationToken *Hash,
	files []*DownloadFile,
	networkType NetworkType,
) (*EndFileDownloadTransaction, error) {

	if recipient == nil {
		return nil, ErrNilAccount
	}

	if operationToken == nil {
		return nil, ErrNilHash
	}

	if len(files) == 0 {
		return nil, ErrNoChanges
	}

	tx := EndFileDownloadTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     EndFileDownloadVersion,
			Deadline:    deadline,
			Type:        EndFileDownload,
			NetworkType: networkType,
		},
		Recipient:      recipient,
		OperationToken: operationToken,
		Files:          files,
	}

	return &tx, nil
}

func (tx *EndFileDownloadTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *EndFileDownloadTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Recipient": %s,
			"OperationToken": %s,
			"Files": %+v,
		`,
		tx.AbstractTransaction.String(),
		tx.Recipient,
		tx.OperationToken,
		tx.Files,
	)
}

func (tx *EndFileDownloadTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	hashV := hashToBuffer(builder, tx.OperationToken)

	b, err := hex.DecodeString(tx.Recipient.PublicKey)
	if err != nil {
		return nil, err
	}

	recipientV := transactions.TransactionBufferCreateByteVector(builder, b)

	filesV, err := downloadFilesToArrayToBuffer(builder, tx.Files)
	if err != nil {
		return nil, err
	}

	transactions.EndFileDownloadTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.EndFileDownloadTransactionBufferAddRecipient(builder, recipientV)
	transactions.EndFileDownloadTransactionBufferAddOperationToken(builder, hashV)
	transactions.EndFileDownloadTransactionBufferAddFileCount(builder, uint16(len(tx.Files)))
	transactions.EndFileDownloadTransactionBufferAddFiles(builder, filesV)
	t := transactions.EndFileDownloadTransactionBufferEnd(builder)
	builder.Finish(t)

	return endFileDownloadTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *EndFileDownloadTransaction) Size() int {
	return EndFileDownloadHeaderSize + len(tx.Files)*(Hash256+StorageSizeSize)
}

type endFileDownloadTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Recipient     string             `json:"fileRecipient"`
		OperationHash hashDto            `json:"operationToken"`
		Files         []*downloadFileDTO `json:"files"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *endFileDownloadTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	acc, err := NewAccountFromPublicKey(dto.Tx.Recipient, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	hash, err := dto.Tx.OperationHash.Hash()
	if err != nil {
		return nil, err
	}

	files := make([]*DownloadFile, len(dto.Tx.Files))

	for i, f := range dto.Tx.Files {
		file, err := f.toStruct()
		if err != nil {
			return nil, err
		}

		files[i] = file
	}

	return &EndFileDownloadTransaction{
		*atx,
		acc,
		hash,
		files,
	}, nil
}
