package sdk

import (
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

type StorageTransaction struct {
	AbstractTransaction
	Action StorageAction
}

// returns a StorageTransaction
func NewStorageDrivePrepareTransaction(deadline *Deadline, duration Duration, driveSize DriveSize, replicas Replicas, networkType NetworkType) (*StorageTransaction, error) {
	return &StorageTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StorageVersion,
			Deadline:    deadline,
			Type:        StoragePrepareDrive,
			NetworkType: networkType,
		},
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
			Type:        StorageDriveProlongation,
			NetworkType: networkType,
		},
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
			Type:        StorageFileDeposit,
			NetworkType: networkType,
		},
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
			Type:        StorageDriveDeposit,
			NetworkType: networkType,
		},
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
			Type:        StorageFileDepositReturn,
			NetworkType: networkType,
		},
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
			Type:        StorageDriveDepositReturn,
			NetworkType: networkType,
		},
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
			Type:        StorageFilePayment,
			NetworkType: networkType,
		},
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
			Type:        StorageDrivePayment,
			NetworkType: networkType,
		},
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
			Type:        StorageCreateDirectory,
			NetworkType: networkType,
		},
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
			Type:        StorageRemoveDirectory,
			NetworkType: networkType,
		},
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
			Type:        StorageUploadFile,
			NetworkType: networkType,
		},
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
			Type:        StorageDownloadFile,
			NetworkType: networkType,
		},
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
			Type:        StorageDeleteFile,
			NetworkType: networkType,
		},
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
			Type:        StorageMoveFile,
			NetworkType: networkType,
		},
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
			Type:        StorageCopyFile,
			NetworkType: networkType,
		},
		Action: &StorageOperationFileAction{
			Source:      source,
			Destination: destination,
		},
	}, nil
}

func (tx *StorageTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *StorageTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Action: %s"
		`,
		tx.AbstractTransaction.String(),
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
	transactions.StorageDriveTransactionBufferAddAction(builder, dV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storageDriveTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *StorageTransaction) Size() int {
	return TransactionHeaderSize + tx.Action.Size()
}
