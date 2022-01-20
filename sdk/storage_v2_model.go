// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
)

type DataModificationState uint8

const (
	Succeeded DataModificationState = iota
	Cancelled
)

type ActiveDataModification struct {
	Id                 *Hash
	Owner              *PublicAccount
	DownloadDataCdi    *Hash
	ExpectedUploadSize StorageSize
	ActualUploadSize   StorageSize
	FolderName         string
	ReadyForApproval   bool
}

func (active *ActiveDataModification) String() string {
	return fmt.Sprintf(
		`
			"Id": %s,
			"Owner": %s,
			"DownloadDataCdi": %s,
			"ExpectedUploadSize": %d,
			"ActualUploadSize": %d,
			"FolderName": %s,
			"ReadyForApproval": %t, 
		`,
		active.Id.String(),
		active.Owner.String(),
		active.DownloadDataCdi.String(),
		active.ExpectedUploadSize,
		active.ActualUploadSize,
		active.FolderName,
		active.ReadyForApproval,
	)
}

type CompletedDataModification struct {
	*ActiveDataModification
	State DataModificationState
}

func (completed *CompletedDataModification) String() string {
	return fmt.Sprintf(
		`
			"ActiveDataModification": %+v,
			"State:" %d,
		`,
		completed.ActiveDataModification,
		completed.State,
	)
}

type ConfirmedUsedSize struct {
	Replicator *PublicAccount
	Size       StorageSize
}

func (confirmed *ConfirmedUsedSize) String() string {
	return fmt.Sprintf(
		`
			"Replicator": %s,
			"Size:" %d,
		`,
		confirmed.Replicator,
		confirmed.Size,
	)
}

type Shard struct {
	DownloadChannelId *Hash
	Replicators       []*PublicAccount
}

type Verification struct {
	VerificationTrigger *Hash
	Expiration          *Timestamp
	Expired             bool
	Shards              []*Shard
}

func (verification *Verification) String() string {
	return fmt.Sprintf(
		`
			"VerificationTrigger": %s,
			"Expiration:" %s,
			"VerificationOpinions:" %t,
			"Shards:" %+v,
		`,
		verification.VerificationTrigger,
		verification.Expiration.String(),
		verification.Expired,
		verification.Shards,
	)
}

type BcDrive struct {
	MultisigAccount            *PublicAccount
	Owner                      *PublicAccount
	RootHash                   *Hash
	Size                       StorageSize
	UsedSize                   StorageSize
	MetaFilesSize              StorageSize
	ReplicatorCount            uint16
	OwnerCumulativeUploadSize  StorageSize
	ActiveDataModifications    []*ActiveDataModification
	CompletedDataModifications []*CompletedDataModification
	ConfirmedUsedSizes         []*ConfirmedUsedSize
	Replicators                []*PublicAccount
	OffboardingReplicators     []*PublicAccount
	Verifications              []*Verification
	Shards                     []*Shard
}

func (drive *BcDrive) String() string {
	return fmt.Sprintf(
		`
		"MultisigAccount": %s,
		"Owner": %s,
		"RootHash": %s,
		"Size": %d,
		"UsedSize": %d,
		"MetaFilesSize": %d,
		"ReplicatorCount": %d,
		"OwnerCumulativeUploadSize": %d,
		"ActiveDataModifications": %+v,
		"CompletedDataModifications": %+v,
		"ConfirmedUsedSizes": %+v,
		"Replicators": %s,
		"OffboardingReplicators": %s,
		"Verifications": %+v,
		"Shards": %+v,
		`,
		drive.MultisigAccount,
		drive.Owner,
		drive.RootHash,
		drive.Size,
		drive.UsedSize,
		drive.MetaFilesSize,
		drive.ReplicatorCount,
		drive.OwnerCumulativeUploadSize,
		drive.ActiveDataModifications,
		drive.CompletedDataModifications,
		drive.ConfirmedUsedSizes,
		drive.Replicators,
		drive.OffboardingReplicators,
		drive.Verifications,
		drive.Shards,
	)
}

type BcDrivesPage struct {
	BcDrives   []*BcDrive
	Pagination Pagination
}

type BcDrivesPageOptions struct {
	PaginationOrderingOptions
}

type DriveInfo struct {
	DriveKey                            *PublicAccount
	LastApprovedDataModificationId      *Hash
	DataModificationIdIsValid           bool
	InitialDownloadWork                 StorageSize
	LastCompletedCumulativeDownloadWork StorageSize
}

func (info *DriveInfo) String() string {
	return fmt.Sprintf(
		`
			"DriveKey": %s, 
		    "LastApprovedDataModificationId": %s,
			"DataModificationIdIsValid": %t,
			"InitialDownloadWork": %d,
			"LastCompletedCumulativeDownloadWork": %d,
		`,
		info.DriveKey,
		info.LastApprovedDataModificationId,
		info.DataModificationIdIsValid,
		info.InitialDownloadWork,
		info.LastCompletedCumulativeDownloadWork,
	)
}

type Replicator struct {
	ReplicatorAccount *PublicAccount
	Version           uint32
	Capacity          Amount
	Drives            []*DriveInfo
}

func (replicator *Replicator) String() string {
	return fmt.Sprintf(
		`
		ReplicatorAccount: %s, 
		Version: %d,
		Capacity: %d,
		Drives: %+v,
		`,
		replicator.ReplicatorAccount,
		replicator.Version,
		replicator.Capacity,
		replicator.Drives,
	)
}

type ReplicatorsPage struct {
	Replicators []*Replicator
	Pagination  Pagination
}

type ReplicatorsPageOptions struct {
	PaginationOrderingOptions
}

// Replicator Onboarding Transaction
type ReplicatorOnboardingTransaction struct {
	AbstractTransaction
	Capacity Amount
}

// Prepare Bc Drive Transaction
type PrepareBcDriveTransaction struct {
	AbstractTransaction
	DriveSize             StorageSize
	VerificationFeeAmount Amount
	ReplicatorCount       uint16
}

// Data Modification Transaction
type DataModificationTransaction struct {
	AbstractTransaction
	DriveKey          *PublicAccount
	DownloadDataCdi   *Hash
	UploadSize        StorageSize
	FeedbackFeeAmount Amount
}

// Data Modification Cancel Transaction
type DataModificationCancelTransaction struct {
	AbstractTransaction
	DriveKey        *PublicAccount
	DownloadDataCdi *Hash
}

// Storage Payment Transaction
type StoragePaymentTransaction struct {
	AbstractTransaction
	DriveKey     *PublicAccount
	StorageUnits Amount
}

// Download Payment Transaction
type DownloadPaymentTransaction struct {
	AbstractTransaction
	DriveKey          *PublicAccount
	DownloadSize      StorageSize
	FeedbackFeeAmount Amount
}

// Download  Transaction
type DownloadTransaction struct {
	AbstractTransaction
	DriveKey          *PublicAccount
	DownloadSize      StorageSize
	FeedbackFeeAmount Amount
	ListOfPublicKeys  []*PublicAccount
}

// Finish Download Transaction
type FinishDownloadTransaction struct {
	AbstractTransaction
	DownloadChannelId *Hash
	FeedbackFeeAmount Amount
}

// Verification Payment Transaction
type VerificationPaymentTransaction struct {
	AbstractTransaction
	DriveKey              *PublicAccount
	VerificationFeeAmount Amount
}

// End Drive Verification Transaction
type EndDriveVerificationTransactionV2 struct {
	AbstractTransaction
	DriveKey            *PublicAccount
	VerificationTrigger *Hash
	ShardId             uint16
	Keys                []*PublicAccount
	Signatures          []string
	Opinions            []uint8
}

// Drive Closure Transaction
type DriveClosureTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
}

// Replicator Offboarding Transaction
type ReplicatorOffboardingTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
}
