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

func NewEndDriveVerificationTransactionV2(
	deadline *Deadline,
	driveKey *PublicAccount,
	verificationTrigger *Hash,
	shardId uint16,
	keys []*PublicAccount,
	signatures []string,
	opinions []uint8,
	networkType NetworkType,
) (*EndDriveVerificationTransactionV2, error) {
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
		keys += k.String()
	}

	signatures := ""
	for _, s := range tx.Signatures {
		signatures += s
	}

	return fmt.Sprintf(
		`
		"AbstractTransaction": %s,
		"Drive": %s,
		"VerificationTrigger": %s,
		"ShardId": %d,
		"Keys": %s,
		"Provers": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
		tx.VerificationTrigger.String(),
		tx.ShardId,
		keys,
		signatures,
	)
}

func keysToArrayToBuffer(builder *flatbuffers.Builder, provers []*PublicAccount) (flatbuffers.UOffsetT, error) {
	psB := make([]flatbuffers.UOffsetT, len(provers))
	for i, p := range provers {
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

func signaturesToArrayToBuffer(builder *flatbuffers.Builder, signatures []string) (flatbuffers.UOffsetT, error) {
	sB := make([]flatbuffers.UOffsetT, len(signatures))
	for i, s := range signatures {
		bsV := transactions.TransactionBufferCreateByteVector(builder, []byte(s))

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
		oB[byteNumber] |= o << i % 8

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

	b, err := hex.DecodeString(tx.DriveKey.PublicKey)
	if err != nil {
		return nil, err
	}

	dkV := transactions.TransactionBufferCreateByteVector(builder, b)
	vtV := transactions.TransactionBufferCreateByteVector(builder, tx.VerificationTrigger[:])

	pkV, err := keysToArrayToBuffer(builder, tx.Keys)
	if err != nil {
		return nil, err
	}

	sV, err := signaturesToArrayToBuffer(builder, tx.Signatures)
	if err != nil {
		return nil, err
	}

	voV := opinionsToArrayToBuffer(builder, len(tx.Keys), len(tx.Signatures), tx.Opinions)

	transactions.EndDriveVerificationTransactionV2BufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.EndDriveVerificationTransactionV2BufferAddDriveKey(builder, dkV)
	transactions.EndDriveVerificationTransactionV2BufferAddVerificationTrigger(builder, vtV)
	transactions.EndDriveVerificationTransactionV2BufferAddKeyCount(builder, uint8(len(tx.Keys)))
	transactions.EndDriveVerificationTransactionV2BufferAddJudgingKeyCount(builder, uint8(len(tx.Signatures)))
	transactions.EndDriveVerificationTransactionV2BufferAddKeys(builder, pkV)
	transactions.EndDriveVerificationTransactionV2BufferAddSignatures(builder, sV)
	transactions.EndDriveVerificationTransactionV2BufferAddOpinions(builder, voV)
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

func parseOpinions(opinions []uint8, keyCount, signaturesCount int) []uint8 {
	opinionsCount := keyCount * signaturesCount
	parsedOpinions := make([]uint8, opinionsCount)

	byteNumber := 0
	for i := 0; i < opinionsCount; i++ {
		parsedOpinions[i] = opinions[byteNumber] & uint8(i%8)

		if i != 0 && i%8 == 0 {
			byteNumber++
		}
	}

	return parsedOpinions
}

type endDriveVerificationTransactionV2DTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey            string   `json:"driveKey"`
		VerificationTrigger hashDto  `json:"verificationTrigger"`
		ShardId             uint16   `json:"shardId"`
		Keys                []string `json:"keys"`
		Signatures          []string `json:"signatures"`
		Opinions            []uint8  `json:"opinions"`
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

	parsedOpinions := parseOpinions(dto.Tx.Opinions, len(keys), len(dto.Tx.Signatures))

	return &EndDriveVerificationTransactionV2{
		*atx,
		driveAccount,
		verificationTrigger,
		dto.Tx.ShardId,
		keys,
		dto.Tx.Signatures,
		parsedOpinions,
	}, nil
}
