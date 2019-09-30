package sdk

import (
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

type StorageTransaction struct {
	AbstractTransaction
	ActionType DriveActionType
	Action     StorageAction
}

// returns a StorageTransaction
func NewStorageDrivePrepareTransaction(deadline *Deadline, duration Duration, driveSize DriveSize, replicas Replicas, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StoragePrepareDrive,
		Action: &StorageDrivePrepareAction{
			Replicas:  replicas,
			DriveSize: driveSize,
			Duration:  duration,
		},
	}, nil
}

// returns a StorageTransaction
func NewStorageDriveProlongationTransaction(deadline *Deadline, duration Duration, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageDriveProlongation,
		Action: &StorageDriveProlongationAction{
			Duration: duration,
		},
	}, nil
}

func NewStorageFileDepositTransaction(deadline *Deadline, fileHash *Hash, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageFileDeposit,
		Action: &StorageFileHashAction{
			FileHash: fileHash,
		},
	}, nil
}

func NewStorageDriveDepositTransaction(deadline *Deadline, directoryHash *Hash, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageDriveDeposit,
		Action: &StorageFileHashAction{
			FileHash: directoryHash,
		},
	}, nil
}

func NewStorageFileDepositReturnTransaction(deadline *Deadline, fileHash *Hash, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageFileDepositReturn,
		Action: &StorageFileHashAction{
			FileHash: fileHash,
		},
	}, nil
}

func NewStorageDriveDepositReturnTransaction(deadline *Deadline, directoryHash *Hash, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageDriveDepositReturn,
		Action: &StorageFileHashAction{
			FileHash: directoryHash,
		},
	}, nil
}

func NewStorageFilePaymentTransaction(deadline *Deadline, fileHash *Hash, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageFilePayment,
		Action: &StorageFileHashAction{
			FileHash: fileHash,
		},
	}, nil
}

func NewStorageDrivePaymentTransaction(deadline *Deadline, directoryHash *Hash, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageDrivePayment,
		Action: &StorageFileHashAction{
			FileHash: directoryHash,
		},
	}, nil
}

func NewStorageCreateDirectoryTransaction(deadline *Deadline, directory *StorageFile, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageCreateDirectory,
		Action: &StorageFileAction{
			File: directory,
		},
	}, nil
}

func NewStorageRemoveDirectoryTransaction(deadline *Deadline, directory *StorageFile, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageRemoveDirectory,
		Action: &StorageFileAction{
			File: directory,
		},
	}, nil
}

func NewStorageUploadFileTransaction(deadline *Deadline, file *StorageFile, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageUploadFile,
		Action: &StorageFileAction{
			File: file,
		},
	}, nil
}

func NewStorageDownloadFileTransaction(deadline *Deadline, file *StorageFile, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageDownloadFile,
		Action: &StorageFileAction{
			File: file,
		},
	}, nil
}
func NewStorageDeleteFileTransaction(deadline *Deadline, file *StorageFile, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageDeleteFile,
		Action: &StorageFileAction{
			File: file,
		},
	}, nil
}

func NewStorageMoveFileTransaction(deadline *Deadline, source *StorageFile, destination *StorageFile, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageMoveFile,
		Action: &StorageOperationFileAction{
			Source:      source,
			Destination: destination,
		},
	}, nil
}

func NewStorageCopyFileTransaction(deadline *Deadline, source *StorageFile, destination *StorageFile, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageCopyFile,
		Action: &StorageOperationFileAction{
			Source:      source,
			Destination: destination,
		},
	}, nil
}
func NewStorageDriveVerificationTransaction(deadline *Deadline, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StorageDrive,
			NetworkType: networkType,
		},
		ActionType: StorageDriveVerification,
		Action:     nil,
	}, nil
}
func (tx *StorageTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *StorageTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"AbstractType": %d,
			"Action: %s"
		`,
		tx.AbstractTransaction.String(),
		tx.ActionType,
		tx.Action.String(),
	)
}

func (tx *StorageTransaction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}
	ai, err := tx.Action.generateBytes()
	if err != nil {
		return nil, err
	}

	dV := transactions.TransactionBufferCreateByteVector(builder, ai)

	transactions.StorageDriveTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.StorageDriveTransactionBufferAddActionType(builder, uint8(tx.ActionType))
	transactions.StorageDriveTransactionBufferAddAction(builder, dV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storageDriveTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *StorageTransaction) Size() int {
	return StorageTransactionHeaderSize + tx.Action.Size()
}
