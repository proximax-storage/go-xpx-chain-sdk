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
	blsPublicKey string,
	networkType NetworkType,
) (*ReplicatorOnboardingTransaction, error) {

	if capacity <= 0 {
		return nil, errors.New("capacity should be positive")
	}

	if len(blsPublicKey) == 0 {
		return nil, ErrNilAccount
	}

	tx := ReplicatorOnboardingTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     ReplicatorOnboardingVersion,
			Type:        ReplicatorOnboarding,
			NetworkType: networkType,
		},
		Capacity:     capacity,
		BlsPublicKey: blsPublicKey,
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
			"BlsPublicKey": %+v,
		`,
		tx.AbstractTransaction.String(),
		tx.Capacity,
		tx.BlsPublicKey,
	)
}

func (tx *ReplicatorOnboardingTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	blsKeyB, err := hex.DecodeString(tx.BlsPublicKey)
	if err != nil {
		return nil, err
	}

	capacityV := transactions.TransactionBufferCreateUint32Vector(builder, tx.Capacity.toArray())
	blsKeyV := transactions.TransactionBufferCreateByteVector(builder, blsKeyB)

	transactions.ReplicatorOnboardingTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.ReplicatorOnboardingTransactionBufferAddCapacity(builder, capacityV)
	transactions.ReplicatorOnboardingTransactionBufferAddBlsPublicKey(builder, blsKeyV)

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
		Capacity     uint64DTO `json:"capacity"`
		BlsPublicKey string    `json:"blsKey"`
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
		dto.Tx.BlsPublicKey,
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
	driveKey string,
	networkType NetworkType,
) (*DriveClosureTransaction, error) {

	if len(driveKey) == 0 {
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
		tx.DriveKey,
	)
}

func (tx *DriveClosureTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	driveKeyB, err := hex.DecodeString(tx.DriveKey)
	if err != nil {
		return nil, err
	}

	driveKeyV := transactions.TransactionBufferCreateByteVector(builder, driveKeyB)

	transactions.DriveClosureTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DriveClosureTransactionBufferAddDriveKey(builder, driveKeyV)
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

	return &DriveClosureTransaction{
		*atx,
		dto.Tx.DriveKey,
	}, nil
}

func NewEndDriveVerificationTransactionV2(
	deadline *Deadline,
	driveKey *PublicAccount,
	verificationTrigger *Hash,
	provers []*PublicAccount,
	verificationOpinions []*VerificationOpinion,
	networkType NetworkType,
) (*EndDriveVerificationTransactionV2, error) {
	tx := EndDriveVerificationTransactionV2{
		AbstractTransaction: AbstractTransaction{
			Version:     EndDriveVerificationV2Version,
			Deadline:    deadline,
			Type:        EndDriveVerificationV2,
			NetworkType: networkType,
		},
		DriveKey:             driveKey,
		VerificationTrigger:  verificationTrigger,
		Provers:              provers,
		VerificationOpinions: verificationOpinions,
	}

	return &tx, nil
}

func (tx *EndDriveVerificationTransactionV2) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *EndDriveVerificationTransactionV2) String() string {
	provers := ""
	for _, p := range tx.Provers {
		provers += p.String()
	}

	verificationOpinions := ""
	for _, vo := range tx.VerificationOpinions {
		verificationOpinions += vo.String()
	}

	return fmt.Sprintf(
		`
		"AbstractTransaction": %s,
		"Drive": %s,
		"VerificationTrigger": %s,
		"Provers": %s,
		"VerificationOpinions": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveKey.String(),
		tx.VerificationTrigger.String(),
		provers,
		verificationOpinions,
	)
}

func resultsToArrayToBuffer(builder *flatbuffers.Builder, results VerificationResults) (flatbuffers.UOffsetT, error) {
	resultsB := make([]flatbuffers.UOffsetT, len(results))
	for i, r := range results {
		b, err := hex.DecodeString(r.Prover.PublicKey)
		if err != nil {
			return 0, err
		}

		pkV := transactions.TransactionBufferCreateByteVector(builder, b)

		rB := byte(0)
		if r.Result {
			rB = byte(1)
		}

		transactions.ResultBufferStart(builder)
		transactions.ResultBufferAddProver(builder, pkV)
		transactions.ResultBufferAddResult(builder, rB)
		resultsB[i] = transactions.ResultBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, resultsB), nil
}

func proversToArrayToBuffer(builder *flatbuffers.Builder, provers []*PublicAccount) (flatbuffers.UOffsetT, error) {
	psB := make([]flatbuffers.UOffsetT, len(provers))
	for i, p := range provers {
		b, err := hex.DecodeString(p.PublicKey)
		if err != nil {
			return 0, err
		}

		pkV := transactions.TransactionBufferCreateByteVector(builder, b)

		transactions.ProversBufferStart(builder)
		transactions.ProversBufferAddProver(builder, pkV)
		psB[i] = transactions.ProversBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, psB), nil
}

func verificationOpinionsToArrayToBuffer(builder *flatbuffers.Builder, verificationOpinions []*VerificationOpinion) (flatbuffers.UOffsetT, error) {
	voB := make([]flatbuffers.UOffsetT, len(verificationOpinions))
	for i, vo := range verificationOpinions {
		b, err := hex.DecodeString(vo.Verifier.PublicKey)
		if err != nil {
			return 0, err
		}

		pkV := transactions.TransactionBufferCreateByteVector(builder, b)
		bsV := transactions.TransactionBufferCreateByteVector(builder, []byte(vo.BlsSignature))
		rsV, err := resultsToArrayToBuffer(builder, vo.Results)
		if err != nil {
			return 0, err
		}

		transactions.VerificationOpinionBufferStart(builder)
		transactions.VerificationOpinionBufferAddVerifier(builder, pkV)
		transactions.VerificationOpinionBufferAddBlsSignature(builder, bsV)
		transactions.VerificationOpinionBufferAddResults(builder, rsV)
		voB[i] = transactions.VerificationOpinionBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, voB), nil
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

	psV, err := proversToArrayToBuffer(builder, tx.Provers)
	if err != nil {
		return nil, err
	}

	voV, err := verificationOpinionsToArrayToBuffer(builder, tx.VerificationOpinions)
	if err != nil {
		return nil, err
	}

	transactions.EndDriveVerificationTransactionV2BufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)
	transactions.EndDriveVerificationTransactionV2BufferAddDriveKey(builder, dkV)
	transactions.EndDriveVerificationTransactionV2BufferAddVerificationTrigger(builder, vtV)
	transactions.EndDriveVerificationTransactionV2BufferAddProversCount(builder, uint16(len(tx.Provers)))
	transactions.EndDriveVerificationTransactionV2BufferAddVerificationOpinionsCount(builder, uint16(len(tx.VerificationOpinions)))
	transactions.EndDriveVerificationTransactionV2BufferAddProvers(builder, psV)
	transactions.EndDriveVerificationTransactionV2BufferAddVerificationOpinions(builder, voV)
	t := transactions.EndDriveVerificationTransactionV2BufferEnd(builder)
	builder.Finish(t)

	return endDriveVerificationV2TransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *EndDriveVerificationTransactionV2) Size() int {
	size := EndDriveVerificationV2HeaderSize + len(tx.Provers)*KeySize
	for _, vo := range tx.VerificationOpinions {
		size += vo.Size()
	}

	return size
}

type verificationOpinionDTO struct {
	Verifier     string          `json:"verifier"`
	BlsSignature string          `json:"blsSignature"`
	Results      map[string]bool `json:"results"`
}

func verificationOpinionDTOArrayToStruct(verificationOpinionDTOs []*verificationOpinionDTO, networkType NetworkType) ([]*VerificationOpinion, error) {
	verificationOpinions := make([]*VerificationOpinion, 0, len(verificationOpinionDTOs))

	for _, dto := range verificationOpinionDTOs {
		verifier, err := NewAccountFromPublicKey(dto.Verifier, networkType)
		if err != nil {
			return nil, err
		}

		results := make(VerificationResults, 0, len(dto.Results))
		for prover, res := range dto.Results {
			acc, err := NewAccountFromPublicKey(prover, networkType)
			if err != nil {
				return nil, err
			}

			results = append(results, &VerificationResult{acc, res})
		}

		verificationOpinions = append(verificationOpinions, &VerificationOpinion{
			verifier,
			BLSSignature(dto.BlsSignature),
			results,
		})
	}

	return verificationOpinions, nil
}

type endDriveVerificationTransactionV2DTO struct {
	Tx struct {
		abstractTransactionDTO
		DriveKey             string                    `json:"driveKey"`
		VerificationTrigger  hashDto                   `json:"verificationTrigger"`
		Provers              []string                  `json:"provers"`
		VerificationOpinions []*verificationOpinionDTO `json:"verificationOpinions"`
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

	provers := make([]*PublicAccount, len(dto.Tx.Provers))
	for i, p := range dto.Tx.Provers {
		provers[i], err = NewAccountFromPublicKey(p, atx.NetworkType)
		if err != nil {
			return nil, err
		}
	}

	vos, err := verificationOpinionDTOArrayToStruct(dto.Tx.VerificationOpinions, atx.NetworkType)
	if err != nil {
		return nil, fmt.Errorf("error parsing VerificationOpinions: %v", err)
	}

	return &EndDriveVerificationTransactionV2{
		*atx,
		driveAccount,
		verificationTrigger,
		provers,
		vos,
	}, nil
}
