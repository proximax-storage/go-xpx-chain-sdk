// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"time"
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
	Id          uint32
	Replicators []*PublicAccount
}

type Verification struct {
	VerificationTrigger *Hash
	Expiration          *Timestamp
	Duration            time.Duration
	Shards              []*Shard
}

func (verification *Verification) String() string {
	return fmt.Sprintf(
		`
			"VerificationTrigger": %s,
			"Expiration:" %s,
			"Duration:" %s,
			"Shards:" %+v,
		`,
		verification.VerificationTrigger,
		verification.Expiration.String(),
		verification.Duration.String(),
		verification.Shards,
	)
}

type DownloadShard struct {
	DownloadChannelId *Hash
}

func (ds *DownloadShard) String() string {
	return fmt.Sprintf(`"DownloadChannelId": %s`, ds.DownloadChannelId.String())
}

type UploadInfoStorageV2 struct {
	Key        *PublicAccount
	UploadSize uint64
}

func (uis *UploadInfoStorageV2) String() string {
	return fmt.Sprintf(
		`
		"Key": %s,
		"UploadSize": %d
		`,
		uis.Key.String(),
		uis.UploadSize,
	)
}

type DataModificationShard struct {
	Replicator             *PublicAccount
	ActualShardReplicators []*UploadInfoStorageV2
	FormerShardReplicators []*UploadInfoStorageV2
	OwnerUpload            uint64
}

func (uis *DataModificationShard) String() string {
	actualShardReplicators := ""
	for _, asr := range uis.ActualShardReplicators {
		actualShardReplicators += asr.String() + " "
	}

	formerShardReplicators := ""
	for _, fsr := range uis.FormerShardReplicators {
		formerShardReplicators += fsr.String() + " "
	}

	return fmt.Sprintf(
		`
		"Replicator": %s,
		"ActualShardReplicators": %s,
		"FormerShardReplicators": %s,
		"OwnerUpload": %d
		`,
		uis.Replicator.String(),
		actualShardReplicators,
		formerShardReplicators,
		uis.OwnerUpload,
	)
}

type BcDrive struct {
	MultisigAccount            *PublicAccount
	Owner                      *PublicAccount
	RootHash                   *Hash
	Size                       StorageSize
	UsedSizeBytes              StorageSize
	MetaFilesSizeBytes         StorageSize
	ReplicatorCount            uint16
	ActiveDataModifications    []*ActiveDataModification
	CompletedDataModifications []*CompletedDataModification
	ConfirmedUsedSizes         []*ConfirmedUsedSize
	Replicators                []*PublicAccount
	OffboardingReplicators     []*PublicAccount
	Verification               *Verification
	DownloadShards             []*DownloadShard
	DataModificationShards     []*DataModificationShard
}

func (drive *BcDrive) String() string {
	return fmt.Sprintf(
		`
		"MultisigAccount": %s,
		"Owner": %s,
		"RootHash": %s,
		"Size": %d,
		"UsedSizeBytes": %d,
		"MetaFilesSizeBytes": %d,
		"ReplicatorCount": %d,
		"ActiveDataModifications": %+v,
		"CompletedDataModifications": %+v,
		"ConfirmedUsedSizes": %+v,
		"Replicators": %s,
		"OffboardingReplicators": %s,
		"Verification": %+v,
		"DownloadShards": %+v,
		"DataModificationShards": %+v,
		`,
		drive.MultisigAccount,
		drive.Owner,
		drive.RootHash,
		drive.Size,
		drive.UsedSizeBytes,
		drive.MetaFilesSizeBytes,
		drive.ReplicatorCount,
		drive.ActiveDataModifications,
		drive.CompletedDataModifications,
		drive.ConfirmedUsedSizes,
		drive.Replicators,
		drive.OffboardingReplicators,
		drive.Verification,
		drive.DownloadShards,
		drive.DataModificationShards,
	)
}

type BcDrivesPage struct {
	BcDrives   []*BcDrive
	Pagination Pagination
}

type BcDrivesPageOptions struct {
	BcDrivesPageFilters
	PaginationOrderingOptions
}

type BcDrivesPageFilters struct {
	Owner string `url:"owner,omitempty"`

	Size     uint64 `url:"size,omitempty"`
	ToSize   uint64 `url:"toSize,omitempty"`
	FromSize uint64 `url:"toSize,omitempty"`

	UsedSize     uint64 `url:"usedSize,omitempty"`
	ToUsedSize   uint64 `url:"toUsedSize,omitempty"`
	FromUsedSize uint64 `url:"toUsedSize,omitempty"`

	MetaFilesSize     uint64 `url:"metaFilesSize,omitempty"`
	ToMetaFilesSize   uint64 `url:"toMetaFilesSize,omitempty"`
	FromMetaFilesSize uint64 `url:"toMetaFilesSize,omitempty"`

	ReplicatorCount     uint64 `url:"replicatorCount,omitempty"`
	ToReplicatorCount   uint64 `url:"toReplicatorCount,omitempty"`
	FromReplicatorCount uint64 `url:"toReplicatorCount,omitempty"`
}

type DriveInfo struct {
	DriveKey                            *PublicAccount
	LastApprovedDataModificationId      *Hash
	InitialDownloadWork                 StorageSize
	LastCompletedCumulativeDownloadWork StorageSize
}

func (info *DriveInfo) String() string {
	return fmt.Sprintf(
		`
			"DriveKey": %s, 
		    "LastApprovedDataModificationId": %s,
			"InitialDownloadWork": %d,
			"LastCompletedCumulativeDownloadWork": %d,
		`,
		info.DriveKey,
		info.LastApprovedDataModificationId,
		info.InitialDownloadWork,
		info.LastCompletedCumulativeDownloadWork,
	)
}

type Replicator struct {
	Account          *PublicAccount
	Version          uint32
	Drives           []*DriveInfo // TODO make map
	DownloadChannels []*Hash
}

func (replicator *Replicator) String() string {
	return fmt.Sprintf(
		`
		Account: %s, 
		Version: %d,
		Drives: %+v,
		DownloadChannels: %+v,
		`,
		replicator.Account,
		replicator.Version,
		replicator.Drives,
		replicator.DownloadChannels,
	)
}

type ReplicatorsPage struct {
	Replicators []*Replicator
	Pagination  Pagination
}

type ReplicatorsPageOptions struct {
	ReplicatorsPageFilters
	PaginationOrderingOptions
}

type ReplicatorsPageFilters struct {
	Version     uint32 `url:"version,omitempty"`
	ToVersion   uint32 `url:"toVersion,omitempty"`
	FromVersion uint32 `url:"fromVersion,omitempty"`

	Capacity     uint64 `url:"capacity,omitempty"`
	ToCapacity   uint64 `url:"toCapacity,omitempty"`
	FromCapacity uint64 `url:"fromCapacity,omitempty"`
}

type CumulativePayment struct {
	Replicator *PublicAccount
	Payment    Amount
}

func (payment *CumulativePayment) String() string {
	return fmt.Sprintf(
		`
			"Replicator": %s,
			"CumulativePayment:" %d,
		`,
		payment.Replicator,
		payment.Payment,
	)
}

type DownloadChannel struct {
	Id                    *Hash
	Consumer              *PublicAccount
	Drive                 *PublicAccount
	DownloadSizeMegabytes StorageSize
	DownloadApprovalCount uint16
	Finished              bool
	ListOfPublicKeys      []*PublicAccount
	ShardReplicators      []*PublicAccount
	CumulativePayments    []*CumulativePayment
}

func (downloadChannel *DownloadChannel) String() string {
	return fmt.Sprintf(
		`
			"Id": %s,
			"Consumer": %s,
			"Drive": %s,
			"DownloadSizeMegabytes": %d,
			"DownloadApprovalCount": %d,
			"ListOfPublicKeys": %s,
			"CumulativePayments": %+v,
		`,
		downloadChannel.Id,
		downloadChannel.Consumer,
		downloadChannel.Drive,
		downloadChannel.DownloadSizeMegabytes,
		downloadChannel.DownloadApprovalCount,
		downloadChannel.ListOfPublicKeys,
		downloadChannel.CumulativePayments,
	)
}

type DownloadChannelsPage struct {
	DownloadChannels []*DownloadChannel
	Pagination       Pagination
}

type DownloadChannelsPageOptions struct {
	DownloadChannelsFilters
	PaginationOrderingOptions
}

type DownloadChannelsFilters struct {
	Consumer string `url:"consumerKey,omitempty"`

	DownloadSize     uint64 `url:"downloadSize,omitempty"`
	ToDownloadSize   uint32 `url:"toDownloadSize,omitempty"`
	FromDownloadSize uint64 `url:"fromDownloadSize,omitempty"`

	DownloadApprovalCount     uint64 `url:"downloadApprovalCount,omitempty"`
	ToDownloadApprovalCount   uint64 `url:"toDownloadApprovalCount,omitempty"`
	FromDownloadApprovalCount uint64 `url:"fromDownloadApprovalCount,omitempty"`
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
	DownloadChannelId *Hash
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

// TODO add VerificationOpinions type with ability to covert to uint8

// End Drive Verification Transaction
type EndDriveVerificationTransactionV2 struct {
	AbstractTransaction
	DriveKey            *PublicAccount
	VerificationTrigger *Hash
	ShardId             uint16
	Keys                []*PublicAccount
	Signatures          []*Signature
	Opinions            uint8
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

type Opinion struct {
	Opinion []OpinionSize
}

// Data Modification Approval Transaction
type DataModificationApprovalTransaction struct {
	AbstractTransaction
	DriveKey               *PublicAccount
	DataModificationId     *Hash
	FileStructureCdi       *Hash
	FileStructureSizeBytes uint64
	MetaFilesSizeBytes     uint64
	UsedDriveSizeBytes     uint64
	JudgingKeysCount       uint8
	OverlappingKeysCount   uint8
	JudgedKeysCount        uint8
	OpinionElementCount    uint16
	PublicKeys             []*PublicAccount
	Signatures             []*Signature
	PresentOpinions        []uint8
	Opinions               []uint64
}

type DataModificationSingleApprovalTransaction struct {
	AbstractTransaction
	DriveKey           *PublicAccount
	DataModificationId *Hash
	PublicKeysCount    uint8
	PublicKeys         []*PublicAccount
	Opinions           []uint64
}

// Download Approval Transaction
type DownloadApprovalTransaction struct {
	AbstractTransaction
	DownloadChannelId    *Hash
	ApprovalTrigger      *Hash
	JudgingKeysCount     uint8
	OverlappingKeysCount uint8
	JudgedKeysCount      uint8
	OpinionElementCount  uint16
	PublicKeys           []*PublicAccount
	Signatures           []*Signature
	PresentOpinions      []uint8
	Opinions             []uint64
}
