package sdk

import (
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
	utils "github.com/proximax-storage/go-xpx-utils"
)

// Modify Drive Transaction
type ModifyDriveTransaction struct {
	AbstractTransaction
	PriceDelta          Amount
	DurationDelta       Duration
	SizeDelta           SizeDelta
	ReplicasDelta       int8
	MinReplicatorsDelta int8
	MinApproversDelta   int8
}

func (tx *ModifyDriveTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ModifyDriveTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"PriceDelta": %d,
			"DurationDelta": %d,
			"SizeDelta": %d,
			"MinReplicatorsDelta": %d,
			"MinApproversDelta": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.PriceDelta,
		tx.DurationDelta,
		tx.SizeDelta,
		tx.MinReplicatorsDelta,
		tx.MinApproversDelta,
	)
}

func (tx *ModifyDriveTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	durationV := transactions.TransactionBufferCreateUint32Vector(builder, tx.DurationDelta.toArray())

	transactions.ModifyDriveTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	am := transactions.TransactionBufferCreateUint32Vector(builder, tx.PriceDelta.toArray())

	sizeDeltaV := transactions.TransactionBufferCreateUint32Vector(builder, tx.SizeDelta.toArray())
	transactions.ModifyDriveTransactionBufferAddPriceDelta(builder, am)
	transactions.ModifyDriveTransactionBufferAddDurationDelta(builder, durationV)

	transactions.ModifyDriveTransactionBufferAddSizeDelta(builder, sizeDeltaV)

	transactions.ModifyDriveTransactionBufferAddReplicasDelta(builder, tx.ReplicasDelta)
	transactions.ModifyDriveTransactionBufferAddMinReplicatorsDelta(builder, tx.MinReplicatorsDelta)
	transactions.ModifyDriveTransactionBufferAddMinApproversDelta(builder, tx.MinApproversDelta)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return modifyDriveTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *ModifyDriveTransaction) Size() int {
	return ModifyDriveHeaderSize
}

type modifyDriveTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		PriceDelta          uint64DTO `json:"priceDelta"`
		DurationDelta       uint64DTO `json:"durationDelta"`
		SizeDelta           uint64DTO `json:"sizeDelta"`
		ReplicasDelta       int8      `json:"replicasDelta"`
		MinReplicatorsDelta int8      `json:"minReplicatorsDelta"`
		MinApproversDelta   int8      `json:"minApproversDelta"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *modifyDriveTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &ModifyDriveTransaction{
		*atx,
		dto.Tx.PriceDelta.toStruct(),
		dto.Tx.DurationDelta.toStruct(),
		dto.Tx.SizeDelta.toStruct(),
		dto.Tx.ReplicasDelta,
		dto.Tx.MinReplicatorsDelta,
		dto.Tx.MinApproversDelta,
	}, nil
}

// Join Drive Transaction

type JoinToDriveTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
}

func (tx *JoinToDriveTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *JoinToDriveTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey,
	)
}

func (tx *JoinToDriveTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	b, err := utils.HexDecodeStringOdd(tx.DriveKey.PublicKey)
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

func (dto *joinToDriveTransactionDTO) toStruct() (Transaction, error) {
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

// Drive File System Transaction
type File struct {
	FileHash *Hash
}
type AddAction struct {
	File
	FileSize FileSize
}

type RemoveAction struct {
	File
}
type DriveFileSystemTransaction struct {
	AbstractTransaction
	RootHash      *Hash
	XorRootHash   *Hash
	AddActions    []*AddAction
	RemoveActions []*RemoveAction
}

func (tx *DriveFileSystemTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DriveFileSystemTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"RootHash": %d,
			"XorRootHash": %d,
			"AddActions": %s,
			"RemoveActions": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.RootHash,
		tx.XorRootHash,
		tx.AddActions,
		tx.RemoveActions,
	)
}

func addActionsToArrayToBuffer(builder *flatbuffers.Builder, addActions []*AddAction) (flatbuffers.UOffsetT, error) {
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

func removeActionsToArrayToBuffer(builder *flatbuffers.Builder, removeActions []*RemoveAction) (flatbuffers.UOffsetT, error) {
	msb := make([]flatbuffers.UOffsetT, len(removeActions))
	for i, m := range removeActions {

		rhV := transactions.TransactionBufferCreateByteVector(builder, m.FileHash[:])
		transactions.RemoveActionBufferStart(builder)
		transactions.RemoveActionBufferAddFileHash(builder, rhV)
		msb[i] = transactions.TransactionBufferEnd(builder)
	}
	return transactions.TransactionBufferCreateUOffsetVector(builder, msb), nil
}

func (tx *DriveFileSystemTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	rhV := transactions.TransactionBufferCreateByteVector(builder, tx.RootHash[:])
	xhV := transactions.TransactionBufferCreateByteVector(builder, tx.XorRootHash[:])

	addActionsV, err := addActionsToArrayToBuffer(builder, tx.AddActions)
	if err != nil {
		return nil, err
	}

	removeActionsV, err := removeActionsToArrayToBuffer(builder, tx.RemoveActions)
	if err != nil {
		return nil, err
	}

	transactions.DriveFileSystemTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DriveFileSystemTransactionBufferAddRootHash(builder, rhV)
	transactions.DriveFileSystemTransactionBufferAddXorRootHash(builder, xhV)

	transactions.DriveFileSystemTransactionBufferAddAddActionsCount(builder, uint8(len(tx.AddActions)))
	transactions.DriveFileSystemTransactionBufferAddRemoveActionsCount(builder, uint8(len(tx.RemoveActions)))

	transactions.DriveFileSystemTransactionBufferAddAddActions(builder, addActionsV)
	transactions.DriveFileSystemTransactionBufferAddRemoveActions(builder, removeActionsV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return driveFileSystemTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DriveFileSystemTransaction) Size() int {
	return DriveFileSystemHeaderSize + len(tx.AddActions) + len(tx.RemoveActions)
}

type driveFileSystemAddActionDTO struct {
	FileHash hashDto   `json:"fileHash"`
	FileSize uint64DTO `json:"fileSize"`
}

type driveFileSystemRemoveActionDTO struct {
	FileHash hashDto `json:"fileHash"`
}

type driveFileSystemTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		RootHash           hashDto                           `json:"rootHash"`
		XorRootHash        hashDto                           `json:"xorRootHash"`
		AddActionsCount    int8                              `json:"addActionsCount"`
		RemoveActionsCount int8                              `json:"removeActionsCount"`
		AddActions         []*driveFileSystemAddActionDTO    `json:"addActions"`
		RemoveActions      []*driveFileSystemRemoveActionDTO `json:"removeActions"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *driveFileSystemTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}
	rHash, err := dto.Tx.RootHash.Hash()
	if err != nil {
		return nil, err
	}
	xorRootHash, err := dto.Tx.RootHash.Hash()
	if err != nil {
		return nil, err
	}

	addActs, err := addActionsDTOArrayToStruct(dto.Tx.AddActions)
	if err != nil {
		return nil, err
	}

	removeActs, err := removeActionsDTOArrayToStruct(dto.Tx.RemoveActions)
	if err != nil {
		return nil, err
	}

	return &DriveFileSystemTransaction{
		*atx,
		rHash,
		xorRootHash,
		addActs,
		removeActs,
	}, nil
}

func addActionsDTOArrayToStruct(addAction []*driveFileSystemAddActionDTO) ([]*AddAction, error) {
	acts := make([]*AddAction, len(addAction))
	var err error = nil
	for i, m := range addAction {
		h, err := m.FileHash.Hash()
		if err != nil {
			return nil, err
		}

		s := m.FileSize.toUint64()

		acts[i] = &AddAction{
			File{
				FileHash: h,
			},
			baseInt64(s),
		}

	}

	return acts, err
}

func removeActionsDTOArrayToStruct(removeAction []*driveFileSystemRemoveActionDTO) ([]*RemoveAction, error) {
	removes := make([]*RemoveAction, len(removeAction))
	var err error = nil
	for i, m := range removeAction {
		h, err := m.FileHash.Hash()
		if err != nil {
			return nil, err
		}
		removes[i] = &RemoveAction{
			File{
				FileHash: h,
			},
		}

	}

	return removes, err
}

// Files Deposit Transaction
type FilesDepositTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
	Files    []*File
}

func (tx *FilesDepositTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *FilesDepositTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveKey": %d,
			"Files": %d,
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

func (tx *FilesDepositTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	b, err := utils.HexDecodeStringOdd(tx.DriveKey.PublicKey)
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

	transactions.FilesDepositTransactionBufferAddFilesCount(builder, uint8(len(tx.Files)))
	transactions.FilesDepositTransactionBufferAddFiles(builder, flsV)

	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return filesDepositTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *FilesDepositTransaction) Size() int {
	return FilesDepositHeaderSize + len(tx.Files)
}

type fileDTO struct {
	FileHash hashDto `json:"fileHash"`
}

type filesDepositTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey   string     `json:"driveKey"`
		FilesCount int8       `json:"filesCount"`
		Files      []*fileDTO `json:"files"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *filesDepositTransactionDTO) toStruct() (Transaction, error) {
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

func filesDTOArrayToStruct(files []*fileDTO) ([]*File, error) {
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

// End Drive Transaction

type EndDriveTransaction struct {
	AbstractTransaction
}

func (tx *EndDriveTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *EndDriveTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
		`,
		tx.AbstractTransaction.String(),
	)
}

func (tx *EndDriveTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	transactions.EndDriveTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
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
		PriceDelta uint64DTO `json:"priceDelta"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *endDriveTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &EndDriveTransaction{
		*atx,
	}, nil
}
