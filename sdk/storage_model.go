package sdk

import (
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

type StorageAction interface {
	String() string
	Size() int
	generateBytes() ([]byte, error)
}

type StorageFile struct {
	Hash       *Hash
	ParentHash *Hash
	Name       string
}

func (s *StorageFile) String() string {
	return fmt.Sprintf("Hash: %s,ParentHash: %s, Name: %s", s.Hash, s.ParentHash, s.Name)
}

func (s *StorageFile) Size() int {
	return FileSize
}

func (s *StorageFile) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	hV := transactions.TransactionBufferCreateByteVector(builder, s.Hash[:])
	pV := transactions.TransactionBufferCreateByteVector(builder, s.ParentHash[:])
	n := builder.CreateString(s.Name)

	transactions.DriveFileBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, s.Size())
	transactions.DriveFileBufferAddHash(builder, hV)
	transactions.DriveFileBufferAddParentHash(builder, pV)
	transactions.DriveFileBufferAddNameSize(builder, byte(len(s.Name)))
	transactions.DriveFileBufferAddName(builder, n)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storageFileDriveSchema().serialize(builder.FinishedBytes()), nil
}

type driveTransactionDTO struct {
	Hash       hashDto `json:"hash"`
	ParentHash hashDto `json:"parentHash"`
	Name       string  `json:"name"`
}

type StorageDrivePrepareAction struct {
	Duration  Duration
	DriveSize DriveSize
	Replicas  Replicas
}

func (s *StorageDrivePrepareAction) String() string {
	return fmt.Sprintf("Duration: %s, DriveSize: %s, Replicas: %s", s.Duration, s.DriveSize, s.Replicas)
}

func (s *StorageDrivePrepareAction) Size() int {
	return PrepareDriveSize
}

func (s *StorageDrivePrepareAction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	dV := transactions.TransactionBufferCreateUint32Vector(builder, s.Duration.toArray())
	dsV := transactions.TransactionBufferCreateUint32Vector(builder, s.DriveSize.toArray())
	rV := transactions.TransactionBufferCreateUint32Vector(builder, s.Replicas.toArray())

	transactions.StoragePrepareDriveBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, s.Size())
	transactions.StoragePrepareDriveBufferAddDuration(builder, dV)
	transactions.StoragePrepareDriveBufferAddDriveSize(builder, dsV)
	transactions.StoragePrepareDriveBufferAddReplicas(builder, rV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storagePrepareDriveSchema().serialize(builder.FinishedBytes()), nil
}

type storagePrepareDriveTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ActionInfo struct {
			Duration  uint64DTO `json:"duration"`
			DriveSize uint64DTO `json:"size"`
			Replicas  uint64DTO `json:"replicas"`
		} `json:"actionInfo"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *storagePrepareDriveTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &StorageTransaction{
		AbstractTransaction: *atx,
		Action: &StorageDrivePrepareAction{
			Duration:  dto.Tx.ActionInfo.Duration.toStruct(),
			DriveSize: dto.Tx.ActionInfo.Replicas.toStruct(),
			Replicas:  dto.Tx.ActionInfo.Replicas.toStruct(),
		},
	}, nil
}

type StorageDriveProlongationAction struct {
	Duration Duration
}

func (s *StorageDriveProlongationAction) String() string {
	return fmt.Sprintf("Duration: %s", s.Duration)
}

func (s *StorageDriveProlongationAction) Size() int {
	return DriveProlongationSize
}

func (s *StorageDriveProlongationAction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	dV := transactions.TransactionBufferCreateUint32Vector(builder, s.Duration.toArray())

	transactions.StorageDriveProlongationBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, s.Size())
	transactions.StorageDriveProlongationBufferAddDuration(builder, dV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storageDriveProlongationSchema().serialize(builder.FinishedBytes()), nil
}

type storageDriveProlongationTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ActionInfo struct {
			Duration uint64DTO `json:"duration"`
		} `json:"action"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *storageDriveProlongationTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &StorageTransaction{
		AbstractTransaction: *atx,
		Action: &StorageDriveProlongationAction{
			Duration: dto.Tx.ActionInfo.Duration.toStruct(),
		},
	}, nil
}

type StorageFileHashAction struct {
	FileHash *Hash
}

func (s *StorageFileHashAction) String() string {
	return fmt.Sprintf("FileHash: %s", s.FileHash)
}

func (s *StorageFileHashAction) Size() int {
	return FileHashSize
}

func (s *StorageFileHashAction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	hV := transactions.TransactionBufferCreateByteVector(builder, s.FileHash[:])
	transactions.StorageFileHashBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, s.Size())
	transactions.StorageFileHashBufferAddFileHash(builder, hV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storageFileHashSchema().serialize(builder.FinishedBytes()), nil
}

type storageFileHashTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ActionInfo struct {
			FileHash hashDto `json:"fileHash"`
		} `json:"action"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *storageFileHashTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}
	fileHash, err := dto.Tx.ActionInfo.FileHash.Hash()
	if err != nil {
		return nil, err
	}
	return &StorageTransaction{
		AbstractTransaction: *atx,
		Action: &StorageFileHashAction{
			FileHash: fileHash,
		},
	}, nil
}

type StorageFileAction struct {
	File *StorageFile
}

func (s *StorageFileAction) String() string {
	return fmt.Sprintf("File: %s", s.File.String())
}

func (s *StorageFileAction) Size() int {
	return FileSize
}

func (s *StorageFileAction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	d, err := s.File.generateBytes()
	if err != nil {
		return nil, err
	}
	dV := transactions.TransactionBufferCreateByteVector(builder, d)

	transactions.StorageFileBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, s.Size())
	transactions.StorageFileBufferAddDriveFile(builder, dV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storageFileSchema().serialize(builder.FinishedBytes()), nil
}

type StorageDirectoryAction struct {
	Directory *StorageFile
}

func (s *StorageDirectoryAction) String() string {
	return fmt.Sprintf("Directory: %s", s.Directory.String())
}

func (s *StorageDirectoryAction) Size() int {
	return FileSize
}

func (s *StorageDirectoryAction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	d, err := s.Directory.generateBytes()
	if err != nil {
		return nil, err
	}
	dV := transactions.TransactionBufferCreateByteVector(builder, d)

	transactions.StorageFileBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, s.Size())
	transactions.StorageFileBufferAddDriveFile(builder, dV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storageDirectorySchema().serialize(builder.FinishedBytes()), nil
}

type storageDirectoryTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ActionInfo struct {
			Directory driveTransactionDTO `json:"directory"`
		} `json:"action"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *storageDirectoryTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}
	fileHash, err := dto.Tx.ActionInfo.Directory.Hash.Hash()
	if err != nil {
		return nil, err
	}
	fileParentHash, err := dto.Tx.ActionInfo.Directory.ParentHash.Hash()
	if err != nil {
		return nil, err
	}
	return &StorageTransaction{
		AbstractTransaction: *atx,
		Action: &StorageDirectoryAction{
			Directory: &StorageFile{
				Hash:       fileHash,
				ParentHash: fileParentHash,
				Name:       dto.Tx.ActionInfo.Directory.Name,
			},
		},
	}, nil
}

type storageDriveVerificationTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *storageDriveVerificationTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &StorageTransaction{
		AbstractTransaction: *atx,
	}, nil
}

type storageFileTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ActionInfo struct {
			Directory driveTransactionDTO `json:"file"`
		} `json:"action"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *storageFileTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}
	fileHash, err := dto.Tx.ActionInfo.Directory.Hash.Hash()
	if err != nil {
		return nil, err
	}
	fileParentHash, err := dto.Tx.ActionInfo.Directory.ParentHash.Hash()
	if err != nil {
		return nil, err
	}
	return &StorageTransaction{
		AbstractTransaction: *atx,
		Action: &StorageFileAction{
			File: &StorageFile{
				Hash:       fileHash,
				ParentHash: fileParentHash,
				Name:       dto.Tx.ActionInfo.Directory.Name,
			},
		},
	}, nil
}

type StorageOperationFileAction struct {
	Source      *StorageFile
	Destination *StorageFile
}

func (s *StorageOperationFileAction) String() string {
	return fmt.Sprintf("Source: %s, Destination: %s", s.Source.String(), s.Destination.String())
}

func (s *StorageOperationFileAction) Size() int {
	return FileSize
}

func (s *StorageOperationFileAction) generateBytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	source, err := s.Source.generateBytes()
	if err != nil {
		return nil, err
	}
	destination, err := s.Destination.generateBytes()
	if err != nil {
		return nil, err
	}
	sourceV := transactions.TransactionBufferCreateByteVector(builder, source)
	destinationV := transactions.TransactionBufferCreateByteVector(builder, destination)

	transactions.StorageFileOperationBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, s.Size())
	transactions.StorageFileOperationBufferAddSource(builder, sourceV)
	transactions.StorageFileOperationBufferAddDestination(builder, destinationV)
	t := transactions.TransactionBufferEnd(builder)
	builder.Finish(t)

	return storageFileOperationSchema().serialize(builder.FinishedBytes()), nil
}

type storageFileOperationTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ActionInfo struct {
			Source      driveTransactionDTO `json:"source"`
			Destination driveTransactionDTO `json:"destination"`
		} `json:"action"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *storageFileOperationTransactionDTO) toStruct() (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}
	sFileHash, err := dto.Tx.ActionInfo.Source.Hash.Hash()
	if err != nil {
		return nil, err
	}
	sFileParentHash, err := dto.Tx.ActionInfo.Source.ParentHash.Hash()
	if err != nil {
		return nil, err
	}
	dFileHash, err := dto.Tx.ActionInfo.Destination.Hash.Hash()
	if err != nil {
		return nil, err
	}
	dFileParentHash, err := dto.Tx.ActionInfo.Destination.ParentHash.Hash()
	if err != nil {
		return nil, err
	}
	return &StorageTransaction{
		AbstractTransaction: *atx,
		Action: &StorageOperationFileAction{
			Source: &StorageFile{
				Hash:       sFileHash,
				ParentHash: sFileParentHash,
				Name:       dto.Tx.ActionInfo.Source.Name,
			},
			Destination: &StorageFile{
				Hash:       dFileHash,
				ParentHash: dFileParentHash,
				Name:       dto.Tx.ActionInfo.Destination.Name,
			},
		},
	}, nil
}
