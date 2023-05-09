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

func NewDeployTransaction(deadline *Deadline, drive, owner *PublicAccount, fileHash *Hash, vmVersion uint64, networkType NetworkType) (*DeployTransaction, error) {
	if drive == nil {
		return nil, ErrNilAccount
	}

	if owner == nil {
		return nil, ErrNilAccount
	}

	tx := DeployTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DeployVersion,
			Deadline:    deadline,
			Type:        Deploy,
			NetworkType: networkType,
		},
		DriveAccount: drive,
		Owner:        owner,
		FileHash:     fileHash,
		VMVersion:    vmVersion,
	}

	return &tx, nil
}

func (tx *DeployTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DeployTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveAccount": %s,
			"Owner": %s,
			"FileHash": %s,
			"VMVersion": %+d,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveAccount,
		tx.Owner,
		tx.FileHash,
		tx.VMVersion,
	)
}

func (tx *DeployTransaction) Size() int {
	return DeployHeaderSize
}

func (tx *DeployTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	hV := hashToBuffer(builder, tx.FileHash)
	dB, err := hex.DecodeString(tx.DriveAccount.PublicKey)
	if err != nil {
		return nil, err
	}

	dV := transactions.TransactionBufferCreateByteVector(builder, dB)

	ownerBytes, err := hex.DecodeString(tx.Owner.PublicKey)
	if err != nil {
		return nil, err
	}

	ownerVector := transactions.TransactionBufferCreateByteVector(builder, ownerBytes)
	vV := transactions.TransactionBufferCreateUint32Vector(builder, baseInt64(tx.VMVersion).toArray())

	transactions.DeployTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.DeployTransactionBufferAddDriveKey(builder, dV)
	transactions.DeployTransactionBufferAddOwner(builder, ownerVector)
	transactions.DeployTransactionBufferAddFileHash(builder, hV)
	transactions.DeployTransactionBufferAddVmVersion(builder, vV)
	t := transactions.DeployTransactionBufferEnd(builder)
	builder.Finish(t)

	return deployTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type deployTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Drive     string    `json:"drive"`
		Owner     string    `json:"owner"`
		FileHash  hashDto   `json:"fileHash"`
		VMVersion uint64DTO `json:"vmVersion"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *deployTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	drive, err := NewAccountFromPublicKey(dto.Tx.Drive, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	owner, err := NewAccountFromPublicKey(dto.Tx.Owner, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	hash, err := dto.Tx.FileHash.Hash()
	if err != nil {
		return nil, err
	}

	return &DeployTransaction{
		*atx,
		drive,
		owner,
		hash,
		dto.Tx.VMVersion.toUint64(),
	}, nil
}

func NewStartExecuteTransaction(deadline *Deadline, supercontract *PublicAccount, mosaics []*Mosaic, function string, functionParameters []int64, networkType NetworkType) (*StartExecuteTransaction, error) {
	if supercontract == nil {
		return nil, ErrNilAccount
	}
	if len(function) == 0 {
		return nil, errors.New("Function should be not empty")
	}
	if len(mosaics) == 0 {
		return nil, errors.New("Mosaics should be not empty")
	}

	tx := StartExecuteTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     StartExecuteVersion,
			Deadline:    deadline,
			Type:        StartExecute,
			NetworkType: networkType,
		},
		SuperContract:      supercontract,
		LockMosaics:        mosaics,
		Function:           function,
		FunctionParameters: functionParameters,
	}

	return &tx, nil
}

func (tx *StartExecuteTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *StartExecuteTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"SuperContract": %s,
			"LockMosaics": %s,
			"Function": %s,
			"functionParameters": %+v,
		`,
		tx.AbstractTransaction.String(),
		tx.SuperContract,
		tx.LockMosaics,
		tx.Function,
		tx.FunctionParameters,
	)
}

func (tx *StartExecuteTransaction) Size() int {
	return StartExecuteHeaderSize + len(tx.Function) + len(tx.FunctionParameters)*8 + len(tx.LockMosaics)*(AmountSize+MosaicIdSize)
}

func (tx *StartExecuteTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	sB, err := hex.DecodeString(tx.SuperContract.PublicKey)
	if err != nil {
		return nil, err
	}

	sV := transactions.TransactionBufferCreateByteVector(builder, sB)

	pB := make([]byte, len(tx.FunctionParameters)*8)
	for i, b := range tx.FunctionParameters {
		binary.LittleEndian.PutUint64(pB[8*i:8*(i+1)], uint64(b))
	}
	pV := transactions.TransactionBufferCreateByteVector(builder, pB)

	functionV := transactions.TransactionBufferCreateByteVector(builder, []byte(tx.Function))
	mb := make([]flatbuffers.UOffsetT, len(tx.LockMosaics))
	for i, mos := range tx.LockMosaics {
		id := transactions.TransactionBufferCreateUint32Vector(builder, mos.AssetId.toArray())
		am := transactions.TransactionBufferCreateUint32Vector(builder, mos.Amount.toArray())
		transactions.MosaicBufferStart(builder)
		transactions.MosaicBufferAddId(builder, id)
		transactions.MosaicBufferAddAmount(builder, am)
		mb[i] = transactions.MosaicBufferEnd(builder)
	}
	mV := transactions.TransactionBufferCreateUOffsetVector(builder, mb)

	dataSizeB := make([]byte, 2)
	binary.LittleEndian.PutUint16(dataSizeB, uint16(len(pB)))
	dataSizeV := transactions.TransactionBufferCreateByteVector(builder, dataSizeB)

	transactions.StartExecuteTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.StartExecuteTransactionBufferAddMosaicsCount(builder, uint8(len(tx.LockMosaics)))
	transactions.StartExecuteTransactionBufferAddFunctionSize(builder, uint8(len(tx.Function)))
	transactions.StartExecuteTransactionBufferAddDataSize(builder, dataSizeV)
	transactions.StartExecuteTransactionBufferAddSuperContract(builder, sV)
	transactions.StartExecuteTransactionBufferAddData(builder, pV)
	transactions.StartExecuteTransactionBufferAddFunction(builder, functionV)
	transactions.StartExecuteTransactionBufferAddMosaics(builder, mV)
	t := transactions.StartExecuteTransactionBufferEnd(builder)
	builder.Finish(t)

	return startExecuteTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type startExecuteTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		SuperContract string       `json:"superContract"`
		Function      string       `json:"function"`
		Data          string       `json:"data"`
		Mosaics       []*mosaicDTO `json:"mosaics"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *startExecuteTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	sc, err := NewAccountFromPublicKey(dto.Tx.SuperContract, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	b, err := hex.DecodeString(dto.Tx.Data)
	if err != nil {
		return nil, err
	}

	pB := make([]int64, len(b)/8)
	for i := 0; i < len(pB); i++ {
		pB[i] = int64(binary.LittleEndian.Uint64(b[8*i : 8*(i+1)]))
	}

	return &StartExecuteTransaction{
		*atx,
		sc,
		dto.Tx.Function,
		mosaics,
		pB,
	}, nil
}

func NewEndExecuteTransaction(deadline *Deadline, mosaics []*Mosaic, token *Hash, status OperationStatus, networkType NetworkType) (*EndExecuteTransaction, error) {
	if token == nil {
		return nil, ErrNilHash
	}
	if len(mosaics) == 0 {
		return nil, errors.New("Mosaics should be not empty")
	}

	tx := EndExecuteTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     EndExecuteVersion,
			Deadline:    deadline,
			Type:        EndExecute,
			NetworkType: networkType,
		},
		OperationToken: token,
		UsedMosaics:    mosaics,
		Status:         status,
	}

	return &tx, nil
}

func NewOperationIdentifyTransaction(deadline *Deadline, operationKey *Hash, networkType NetworkType) (*OperationIdentifyTransaction, error) {
	if operationKey == nil {
		return nil, ErrNilHash
	}

	tx := OperationIdentifyTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     OperationIdentifyVersion,
			Deadline:    deadline,
			Type:        OperationIdentify,
			NetworkType: networkType,
		},
		OperationHash: operationKey,
	}

	return &tx, nil
}

func (tx *OperationIdentifyTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *OperationIdentifyTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"OperationHash": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.OperationHash,
	)
}

func (tx *OperationIdentifyTransaction) Size() int {
	return OperationIdentifyHeaderSize
}

func (tx *OperationIdentifyTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	hV := hashToBuffer(builder, tx.OperationHash)

	transactions.OperationIdentifyTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.OperationIdentifyTransactionBufferAddOperationToken(builder, hV)
	t := transactions.OperationIdentifyTransactionBufferEnd(builder)
	builder.Finish(t)

	return operationIdentifyTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type operationIdentifyTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Token hashDto `json:"operationToken"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *operationIdentifyTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	hash, err := dto.Tx.Token.Hash()
	if err != nil {
		return nil, err
	}

	return &OperationIdentifyTransaction{
		*atx,
		hash,
	}, nil
}

func NewEndOperationTransaction(deadline *Deadline, mosaics []*Mosaic, token *Hash, status OperationStatus, networkType NetworkType) (*EndOperationTransaction, error) {
	if status == Unknown {
		return nil, errors.New("Status should be not unknown")
	}
	if len(mosaics) == 0 {
		return nil, errors.New("Mosaics should be not empty")
	}

	tx := EndOperationTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     EndOperationVersion,
			Deadline:    deadline,
			Type:        EndOperation,
			NetworkType: networkType,
		},
		UsedMosaics:    mosaics,
		OperationToken: token,
		Status:         status,
	}

	return &tx, nil
}

func (tx *EndOperationTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *EndOperationTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Status": %d,
			"OperationToken": %s,
			"UsedMosaics": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Status,
		tx.OperationToken,
		tx.UsedMosaics,
	)
}

func (tx *EndOperationTransaction) Size() int {
	return EndOperationHeaderSize + len(tx.UsedMosaics)*(AmountSize+MosaicIdSize)
}

func (tx *EndOperationTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	mb := make([]flatbuffers.UOffsetT, len(tx.UsedMosaics))
	for i, mos := range tx.UsedMosaics {
		id := transactions.TransactionBufferCreateUint32Vector(builder, mos.AssetId.toArray())
		am := transactions.TransactionBufferCreateUint32Vector(builder, mos.Amount.toArray())
		transactions.MosaicBufferStart(builder)
		transactions.MosaicBufferAddId(builder, id)
		transactions.MosaicBufferAddAmount(builder, am)
		mb[i] = transactions.MosaicBufferEnd(builder)
	}
	mV := transactions.TransactionBufferCreateUOffsetVector(builder, mb)
	hV := hashToBuffer(builder, tx.OperationToken)

	transactions.EndOperationTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.EndOperationTransactionBufferAddMosaicsCount(builder, uint8(len(tx.UsedMosaics)))
	transactions.EndOperationTransactionBufferAddOperationToken(builder, hV)
	transactions.EndOperationTransactionBufferAddStatus(builder, uint16(tx.Status))
	transactions.EndOperationTransactionBufferAddMosaics(builder, mV)
	t := transactions.EndOperationTransactionBufferEnd(builder)
	builder.Finish(t)

	return endOperationTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type endOperationTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		Mosaics []*mosaicDTO    `json:"mosaics"`
		Token   hashDto         `json:"operationToken"`
		Status  OperationStatus `json:"result"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *endOperationTransactionDTO) toStruct(*Hash) (Transaction, error) {
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

	hash, err := dto.Tx.Token.Hash()
	if err != nil {
		return nil, err
	}

	return &EndOperationTransaction{
		*atx,
		mosaics,
		hash,
		dto.Tx.Status,
	}, nil
}

func NewSuperContractFileSystemTransaction(
	deadline *Deadline,
	driveKey string,
	newRootHash *Hash,
	oldRootHash *Hash,
	addActions []*Action,
	removeActions []*Action,
	networkType NetworkType,
) (*SuperContractFileSystemTransaction, error) {
	tx, err := NewDriveFileSystemTransaction(deadline, driveKey, newRootHash, oldRootHash, addActions, removeActions, networkType)
	if err != nil {
		return nil, err
	}
	tx.Type = SuperContractFileSystem
	tx.Version = SuperContractFileSystemVersion

	return tx, nil
}

func NewDeactivateTransaction(deadline *Deadline, sc string, driveKey string, networkType NetworkType) (*DeactivateTransaction, error) {
	if len(sc) != 64 {
		return nil, errors.New("wrong super contract key")
	}
	if len(driveKey) != 64 {
		return nil, errors.New("wrong drive key")
	}

	tx := DeactivateTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DeactivateVersion,
			Deadline:    deadline,
			Type:        Deactivate,
			NetworkType: networkType,
		},
		SuperContract: sc,
		DriveKey:      driveKey,
	}

	return &tx, nil
}

func (tx *DeactivateTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DeactivateTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"SuperContract": %s,
			"DriveKey": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.SuperContract,
		tx.DriveKey,
	)
}

func (tx *DeactivateTransaction) Size() int {
	return DeactivateHeaderSize
}

func (tx *DeactivateTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	scB, err := hex.DecodeString(tx.SuperContract)
	if err != nil {
		return nil, err
	}
	scV := transactions.TransactionBufferCreateByteVector(builder, scB)

	driveB, err := hex.DecodeString(tx.DriveKey)
	if err != nil {
		return nil, err
	}
	driveV := transactions.TransactionBufferCreateByteVector(builder, driveB)

	transactions.DeactivateTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.DeactivateTransactionBufferAddSuperContract(builder, scV)
	transactions.DeactivateTransactionBufferAddDriveKey(builder, driveV)
	t := transactions.DeactivateTransactionBufferEnd(builder)
	builder.Finish(t)

	return deactivateTransactionSchema().serialize(builder.FinishedBytes()), nil
}

type deactivateTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		SuperContract string `json:"superContract"`
		Drive         string `json:"driveKey"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *deactivateTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	return &DeactivateTransaction{
		*atx,
		dto.Tx.SuperContract,
		dto.Tx.Drive,
	}, nil
}
