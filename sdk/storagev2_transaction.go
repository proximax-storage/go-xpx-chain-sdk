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

func NewDownloadTransaction(
	deadline *Deadline,
	downloadSize StorageSize,
	feedbackFeeAmount Amount,
	listOfPublicKeys []*Hash,
	networkType NetworkType,
) (*DownloadTransaction, error) {

	if downloadSize == 0 {
		return nil, errors.New("nothing to download")
	}

	if feedbackFeeAmount == 0 {
		return nil, errors.New("feedbackFeeAmount should be positive")
	}

	if len(listOfPublicKeys) == 0 {
		return nil, ErrNoChanges
	}

	tx := DownloadTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     DownloadVersion,
			Type:        Download,
			NetworkType: networkType,
		},
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
			"DownloadSize": %s,
			"FeedbackFeeAmount": %d,
			"ListOfPublicKeys": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DownloadSize,
		tx.FeedbackFeeAmount,
		tx.ListOfPublicKeys,
	)
}

func hashesToArrayToBuffer(builder *flatbuffers.Builder, hashes []*Hash) (flatbuffers.UOffsetT, error) {
	msb := make([]flatbuffers.UOffsetT, len(hashes))
	for i, m := range hashes {
		rhV := transactions.TransactionBufferCreateByteVector(builder, m[:])
		transactions.HashesBufferStart(builder)
		transactions.HashesBufferAddHashes(builder, rhV)
		msb[i] = transactions.TransactionBufferEnd(builder)
	}

	return transactions.TransactionBufferCreateUOffsetVector(builder, msb), nil
}

func (tx *DownloadTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	downloadSizeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.DownloadSize.toArray())
	feedbackFeeAmountV := transactions.TransactionBufferCreateUint32Vector(builder, tx.FeedbackFeeAmount.toArray())
	lpksV, err := hashesToArrayToBuffer(builder, tx.ListOfPublicKeys)
	if err != nil {
		return nil, err
	}

	transactions.DownloadTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DownloadTransactionBufferAddDownloadSize(builder, downloadSizeV)
	transactions.DownloadTransactionBufferAddFeedbackFeeAmount(builder, feedbackFeeAmountV)
	transactions.DownloadTransactionBufferAddPublicKeyCount(builder, uint16(len(tx.ListOfPublicKeys)))
	transactions.DownloadTransactionBufferAddListOfPublicKeys(builder, lpksV)

	t := transactions.DownloadTransactionBufferEnd(builder)
	builder.Finish(t)

	return downloadTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DownloadTransaction) Size() int {
	return DownloadHeaderSize + len(tx.ListOfPublicKeys)*Hash256
}

type downloadTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DownloadSize      uint64DTO  `json:"downloadSize"`
		FeedbackFeeAmount uint64DTO  `json:"feedbackFeeAmount"`
		ListOfPublicKeys  []*hashDto `json:"listOfPublicKeys"`
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

	publicKeys := make([]*Hash, len(dto.Tx.ListOfPublicKeys))

	for i, pub := range dto.Tx.ListOfPublicKeys {
		publicKey, err := pub.Hash()
		if err != nil {
			return nil, err
		}

		publicKeys[i] = publicKey
	}

	return &DownloadTransaction{
		*atx,
		dto.Tx.DownloadSize.toStruct(),
		dto.Tx.FeedbackFeeAmount.toStruct(),
		publicKeys,
	}, nil
}

func NewDownloadApprovalTransaction(
	deadline *Deadline,
	downloadChannelId *Hash,
	sequenceNumber uint16,
	responseToFinishDownloadTransaction uint8,
	publicKeys []*Hash,
	signatures []*Hash,
	presentOpinions []uint8,
	opinions []Opinions,
	networkType NetworkType,
) (*DownloadApprovalTransaction, error) {

	if downloadChannelId == nil {
		return nil, ErrNilHash
	}

	if len(opinions) == 0 {
		return nil, ErrNoChanges
	}

	tx := DownloadApprovalTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     DownloadApprovalVersion,
			Type:        DownloadApproval,
			NetworkType: networkType,
		},
		DownloadChannelId:                   downloadChannelId,
		SequenceNumber:                      sequenceNumber,
		ResponseToFinishDownloadTransaction: responseToFinishDownloadTransaction,
		PublicKeys:                          publicKeys,
		Signatures:                          signatures,
		PresentOpinions:                     presentOpinions,
		Opinions:                            opinions,
	}

	return &tx, nil
}

func (tx *DownloadApprovalTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DownloadApprovalTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DownloadChannelId": %s,
			"SequenceNumber": %d,
			"ResponseToFinishDownloadTransaction": %d,
			"PublicKeys": %s,
			"Signatures": %s,
			"PresentOpinions": %s,
			"Opinions": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.DownloadChannelId,
		tx.SequenceNumber,
		tx.ResponseToFinishDownloadTransaction,
		tx.PublicKeys,
		tx.Signatures,
		tx.PresentOpinions,
		tx.Opinions,
	)
}

func (tx *DownloadApprovalTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	downloadChannelIdV := transactions.TransactionBufferCreateByteVector(builder, tx.DownloadChannelId[:])
	pksV, err := hashesToArrayToBuffer(builder, tx.PublicKeys)
	if err != nil {
		return nil, err
	}
	sgsV, err := hashesToArrayToBuffer(builder, tx.Signatures)
	if err != nil {
		return nil, err
	}

	posB := make([]flatbuffers.UOffsetT, len(tx.PresentOpinions))
	for i, m := range tx.PresentOpinions {
		transactions.PresentOpinionsBufferStart(builder)
		transactions.PresentOpinionsBufferAddPresent(builder, m)
		transactions.PresentOpinionsBufferEnd(builder)
		posB[i] = transactions.TransactionBufferEnd(builder)
	}
	posV, err := transactions.TransactionBufferCreateUOffsetVector(builder, posB), nil
	if err != nil {
		return nil, err
	}

	opsB := make([]flatbuffers.UOffsetT, len(tx.Opinions))
	for i, m := range tx.Opinions {
		rhV := transactions.TransactionBufferCreateUint32Vector(builder, m.toArray())
		transactions.OpinionsBufferStart(builder)
		transactions.OpinionsBufferAddOpinion(builder, rhV)
		transactions.OpinionsBufferEnd(builder)
		opsB[i] = transactions.TransactionBufferEnd(builder)
	}
	opsV, err := transactions.TransactionBufferCreateUOffsetVector(builder, opsB), nil
	if err != nil {
		return nil, err
	}

	transactions.DownloadApprovalTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DownloadApprovalTransactionBufferAddDownloadChannelId(builder, downloadChannelIdV)
	transactions.DownloadApprovalTransactionBufferAddSequenceNumber(builder, tx.SequenceNumber)
	transactions.DownloadApprovalTransactionBufferAddResponseToFinishDownloadTransaction(builder, tx.ResponseToFinishDownloadTransaction)

	totalKeysCount := len(tx.PublicKeys)        // JudgingKeysCount + OverlappingKeysCount + JudgedKeysCount
	totalJudgingKeysCount := len(tx.Signatures) // JudgingKeysCount + OverlappingKeysCount
	judgedKeysCount := totalKeysCount - totalJudgingKeysCount
	totalJudgedKeysCount := (len(tx.PresentOpinions)*8 - 7) / totalJudgingKeysCount // OverlappingKeysCount + JudgedKeysCount
	judgingKeysCount := totalKeysCount - totalJudgedKeysCount
	overlappingKeysCount := totalKeysCount - judgedKeysCount - judgingKeysCount

	transactions.DownloadApprovalTransactionBufferAddJudgingKeysCount(builder, uint8(judgingKeysCount))
	transactions.DownloadApprovalTransactionBufferAddJudgedKeysCount(builder, uint8(judgedKeysCount))
	transactions.DownloadApprovalTransactionBufferAddOverlappingKeysCount(builder, uint8(overlappingKeysCount))
	transactions.DownloadApprovalTransactionBufferAddOpinionElementCount(builder, uint8(len(tx.Opinions)))
	transactions.DownloadApprovalTransactionBufferAddPublicKeys(builder, pksV)
	transactions.DownloadApprovalTransactionBufferAddSignatures(builder, sgsV)
	transactions.DownloadApprovalTransactionBufferAddPresentOpinions(builder, posV)
	transactions.DownloadApprovalTransactionBufferAddOpinions(builder, opsV)

	t := transactions.DownloadApprovalTransactionBufferEnd(builder)
	builder.Finish(t)

	return downloadApprovalTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DownloadApprovalTransaction) Size() int {
	return DownloadApprovalHeaderSize + len(tx.PublicKeys)*Hash256 + len(tx.Signatures)*Hash256 + len(tx.PresentOpinions) + len(tx.Opinions)*StorageSizeSize
}

type downloadApprovalTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DownloadChannelId                   hashDto     `json:downloadChannelId`
		SequenceNumber                      uint16      `json:sequenceNumber`
		ResponseToFinishDownloadTransaction uint8       `json:responseToFinishDownloadTransaction`
		PublicKeys                          []*hashDto  `json:publicKeys`
		Signatures                          []*hashDto  `json:signatures`
		PresentOpinions                     []uint8     `json:presentOpinions`
		Opinions                            []uint64DTO `json:opinions`
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
		return nil, err
	}

	publicKeys := make([]*Hash, len(dto.Tx.PublicKeys))
	for i, pub := range dto.Tx.PublicKeys {
		publicKey, err := pub.Hash()
		if err != nil {
			return nil, err
		}
		publicKeys[i] = publicKey
	}

	signatures := make([]*Hash, len(dto.Tx.Signatures))
	for i, sign := range dto.Tx.PublicKeys {
		signature, err := sign.Hash()
		if err != nil {
			return nil, err
		}
		signatures[i] = signature
	}

	presentOpinions := make([]uint8, len(dto.Tx.PresentOpinions))
	for i, present := range dto.Tx.PresentOpinions {
		presentOpinions[i] = present
	}

	opinions := make([]Opinions, len(dto.Tx.Opinions))
	for i, opinion := range dto.Tx.Opinions {
		opinions[i] = opinion.toStruct()
	}

	return &DownloadApprovalTransaction{
		*atx,
		downloadChannelId,
		dto.Tx.SequenceNumber,
		dto.Tx.ResponseToFinishDownloadTransaction,
		publicKeys,
		signatures,
		presentOpinions,
		opinions,
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

	if downloadSize == 0 {
		return nil, errors.New("downloadSize should be positive")
	}

	if feedbackFeeAmount == 0 {
		return nil, errors.New("feedbackFeeAmount should be positive")
	}

	tx := DownloadPaymentTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     DownloadPaymentVersion,
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
		"DownloadChannelId": %s,
		"DownloadSize": %d,
		"FeedbackFeeAmount": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.DownloadChannelId,
		tx.DownloadSize,
		tx.FeedbackFeeAmount,
	)
}

func (tx *DownloadPaymentTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	downloadChannelIdV := transactions.TransactionBufferCreateByteVector(builder, tx.DownloadChannelId[:])
	downloadSizeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.DownloadSize.toArray())
	feedbackFeeAmountV := transactions.TransactionBufferCreateUint32Vector(builder, tx.FeedbackFeeAmount.toArray())

	transactions.DownloadPaymentTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.DownloadPaymentTransactionBufferAddDownloadChannelId(builder, downloadChannelIdV)
	transactions.DownloadPaymentTransactionBufferAddDownloadSize(builder, downloadSizeV)
	transactions.DownloadPaymentTransactionBufferAddFeedbackFeeAmount(builder, feedbackFeeAmountV)

	t := transactions.DownloadPaymentTransactionBufferEnd(builder)
	builder.Finish(t)

	return downloadPaymentTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *DownloadPaymentTransaction) Size() int {
	return DownloadPaymentHeaderSize
}

type downloadPaymentTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DownloadChannelId *hashDto  `json:downloadChannelId`
		DownloadSize      uint64DTO `json:downloadSize`
		FeedbackFeeAmount uint64DTO `json:feedbackFeeAmount`
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

func NewFinishDownloadTransaction(
	deadline *Deadline,
	downloadChannelId *Hash,
	feedbackFeeAmount Amount,
	networkType NetworkType,
) (*FinishDownloadTransaction, error) {

	if downloadChannelId == nil {
		return nil, ErrNilHash
	}

	if feedbackFeeAmount == 0 {
		return nil, errors.New("feedbackFeeAmount should be positive")
	}

	tx := FinishDownloadTransaction{
		AbstractTransaction: AbstractTransaction{
			Deadline:    deadline,
			Version:     FinishDownloadVersion,
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
			"DownloadChannelId": %s,
			"FeedbackFeeAmount": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.DownloadChannelId,
		tx.FeedbackFeeAmount,
	)
}

func (tx *FinishDownloadTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	downloadChannelIdV := transactions.TransactionBufferCreateByteVector(builder, tx.DownloadChannelId[:])
	feedbackFeeAmountV := transactions.TransactionBufferCreateUint32Vector(builder, tx.FeedbackFeeAmount.toArray())

	transactions.FinishDownloadTransactionBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.FinishDownloadTransactionBufferAddDownloadChannelId(builder, downloadChannelIdV)
	transactions.FinishDownloadTransactionBufferAddFeedbackFeeAmount(builder, feedbackFeeAmountV)

	t := transactions.FinishDownloadTransactionBufferEnd(builder)
	builder.Finish(t)

	return finishDownloadTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func (tx *FinishDownloadTransaction) Size() int {
	return FinishDownloadHeaderSize
}

type finishDownloadTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		DownloadChannelId *hashDto  `json:downloadChannelId`
		FeedbackFeeAmount uint64DTO `json:feedbackFeeAmount`
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
