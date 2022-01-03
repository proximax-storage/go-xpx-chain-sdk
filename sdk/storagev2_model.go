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
	ActiveDataModification []*ActiveDataModification
	State                  DataModificationState
}

func (completed *CompletedDataModification) String() string {
	return fmt.Sprintf(
		`
			"ActiveDataModification": %s,
			"State:" %d,
		`,
		completed.ActiveDataModification,
		completed.State,
	)
}

type ConfirmedUsedSize struct {
	Replicator *Hash
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

type VerificationState uint8

const (
	PendingVerification VerificationState = iota
	CanceledVerification
	FinishedVerification
)

type VerificationOpinion struct {
	Prover *Hash
	Result uint8
}

func (verificationOpinion *VerificationOpinion) String() string {
	return fmt.Sprintf(
		`
			"Prover": %s,
			"Result:" %d,
		`,
		verificationOpinion.Prover,
		verificationOpinion.Result,
	)
}

type Verification struct {
	VerificationTrigger  *Hash
	State                VerificationState
	VerificationOpinions []*VerificationOpinion
}

func (verification *Verification) String() string {
	return fmt.Sprintf(
		`
			"VerificationTrigger": %s,
			"State:" %d,
			"VerificationOpinions:" %+v,
		`,
		verification.VerificationTrigger,
		verification.State,
		verification.VerificationOpinions,
	)
}

type BcDrive struct {
	BcDriveAccount             *PublicAccount
	OwnerAccount               *PublicAccount
	RootHash                   *Hash
	DriveSize                  StorageSize
	UsedSize                   StorageSize
	MetaFilesSize              StorageSize
	ReplicatorCount            uint16
	OwnerCumulativeUploadSize  StorageSize
	ActiveDataModifications    []*ActiveDataModification
	CompletedDataModifications []*CompletedDataModification
	ConfirmedUsedSizes         []*ConfirmedUsedSize
	Replicators                []*Hash
	Verifications              []*Verification
}

func (drive *BcDrive) String() string {
	return fmt.Sprintf(
		`
		"BcDriveAccount": %s,
		"OwnerAccount": %s,
		"RootHash": %s,
		"DriveSize": %d,
		"UsedSize": %d,
		"MetaFilesSize": %d,
		"ReplicatorCount": %d,
		"OwnerCumulativeUploadSize": %d,
		"ActiveDataModifications": %+v,
		"CompletedDataModifications": %+v,
		"ConfirmedUsedSizes": %+v,
		"Replicators": %s,
		"Verifications": %+v,
		`,
		drive.BcDriveAccount,
		drive.OwnerAccount,
		drive.RootHash,
		drive.DriveSize,
		drive.UsedSize,
		drive.MetaFilesSize,
		drive.ReplicatorCount,
		drive.OwnerCumulativeUploadSize,
		drive.ActiveDataModifications,
		drive.CompletedDataModifications,
		drive.ConfirmedUsedSizes,
		drive.Replicators,
		drive.Verifications,
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
	Drive                          *PublicAccount
	LastApprovedDataModificationId *Hash
	DataModificationIdIsValid      bool
	InitialDownloadWork            StorageSize
}

func (info *DriveInfo) String() string {
	return fmt.Sprintf(
		`
			"Drive": %s, 
		    "LastApprovedDataModificationId": %s,
			"DataModificationIdIsValid": %t,
			"InitialDownloadWork": %d,
		`,
		info.Drive,
		info.LastApprovedDataModificationId,
		info.DataModificationIdIsValid,
		info.InitialDownloadWork,
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

type Payment struct {
	Replicator *Hash
	Payment    Amount
}

func (payment *Payment) String() string {
	return fmt.Sprintf(
		`
			"Replicator": %s,
			"Payment:" %d,
		`,
		payment.Replicator,
		payment.Payment,
	)
}

type DownloadChannel struct {
	Id                    *Hash
	Consumer              *Hash
	DownloadSize          StorageSize
	DownloadApprovalCount uint16
	ListOfPublicKeys      []*Hash
	CumulativePayments    []*Payment
}

func (downloadChannel *DownloadChannel) String() string {
	return fmt.Sprintf(
		`
			"Id": %s,
			"Consumer": %s,
			"DownloadSize": %d,
			"DownloadApprovalCount": %d,
			"ListOfPublicKeys": %s,
			"CumulativePayments": %+v,
		`,
		downloadChannel.Id,
		downloadChannel.Consumer,
		downloadChannel.DownloadSize,
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

// Drive Closure Transaction
type DriveClosureTransaction struct {
	AbstractTransaction
	Drive string
}

// Download Transaction
type DownloadTransaction struct {
	AbstractTransaction
	DownloadSize      StorageSize
	FeedbackFeeAmount Amount
	ListOfPublicKeys  []*Hash
}

// Download Approval Transaction
type DownloadApprovalTransaction struct {
	AbstractTransaction
	DownloadChannelId                   *Hash
	SequenceNumber                      uint16
	ResponseToFinishDownloadTransaction uint8
	PublicKeys                          []*Hash
	Signatures                          []*Hash
	PresentOpinions                     []*uint8
	Opinions                            []*Opinions
}

// Download Approval Transaction
type DownloadPaymentTransaction struct {
	AbstractTransaction
	DownloadChannelId *Hash
	DownloadSize      StorageSize
	FeedbackFeeAmount Amount
}

// Finish Download Transaction
type FinishDownloadTransaction struct {
	AbstractTransaction
	DownloadChannelId *Hash
	FeedbackFeeAmount Amount
}
