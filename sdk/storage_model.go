// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import "fmt"

type DriveState uint8

const (
	NotStarted DriveState = iota
	Pending
	InProgress
	Finished
)

type PaymentInformation struct {
	Receiver *PublicAccount
	Amount   Amount
	Height   Height
}

type BillingDescription struct {
	Start    Height
	End      Height
	Payments []*PaymentInformation
}

type ReplicatorInfo struct {
	Account                     *PublicAccount
	Start                       Height
	End                         Height
	Index                       int
	ActiveFilesWithoutDeposit   map[Hash]bool
}

type FileInfo struct {
	FileSize 	StorageSize
}

type Drive struct {
	DriveKey         *PublicAccount
	Start            Height
	State            DriveState
	Owner            *PublicAccount
	RootHash         *Hash
	Duration         Duration
	BillingPeriod    Duration
	BillingPrice     Amount
	DriveSize        StorageSize
	Replicas         uint16
	MinReplicators   uint16
	PercentApprovers uint8
	BillingHistory   []*BillingDescription
	Files            map[Hash]*FileInfo
	Replicators      map[string]*ReplicatorInfo
	UploadPayments   []*PaymentInformation
}

// Prepare Drive Transaction
type PrepareDriveTransaction struct {
	AbstractTransaction
	Owner            *PublicAccount
	Duration         Duration
	BillingPeriod    Duration
	BillingPrice     Amount
	DriveSize        StorageSize
	Replicas         uint16
	MinReplicators   uint16
	PercentApprovers uint8
}

// Join Drive Transaction

type JoinToDriveTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
}

type File struct {
	FileHash *Hash
}

func (file *File) String() string {
	return fmt.Sprintf(
		`
			"FileHash": %s,
		`,
		file.FileHash,
	)
}

type Action struct {
	FileHash *Hash
	FileSize StorageSize
}

func (action *Action) String() string {
	return fmt.Sprintf(
		`
			"FileHash": %s,
			"FileSize": %s,
		`,
		action.FileHash,
		action.FileSize,
	)
}

type DriveFileSystemTransaction struct {
	AbstractTransaction
	DriveKey      *PublicAccount
	NewRootHash   *Hash
	OldRootHash   *Hash
	AddActions    []*Action
	RemoveActions []*Action
}

// Files Deposit Transaction
type FilesDepositTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
	Files    []*File
}

// End Drive Transaction

type EndDriveTransaction struct {
	AbstractTransaction
	DriveKey *PublicAccount
}

type UploadInfo struct {
	Participant     *PublicAccount
	UploadedSize    Amount
}

func (info *UploadInfo) String() string {
	return fmt.Sprintf(
		`
			"Participant": %s,
			"UploadedSize": %s,
		`,
		info.Participant,
		info.UploadedSize,
	)
}

// Drive Files Reward Transaction

type DriveFilesRewardTransaction struct {
	AbstractTransaction
	UploadInfos []*UploadInfo
}

// Start Drive Verification Transaction

type StartDriveVerificationTransaction struct {
	AbstractTransaction
	DriveKey    *PublicAccount
}

type FailureVerification struct {
	Replicator  *PublicAccount
	BlochHash   *Hash
}

func (fail *FailureVerification) String() string {
	return fmt.Sprintf(
		`
			"Replicator": %s,
			"BlochHash": %s,
		`,
		fail.Replicator,
		fail.BlochHash,
	)
}

// End Drive Verification Transaction

type EndDriveVerificationTransaction struct {
	AbstractTransaction
	Failures []*FailureVerification
}

type VerificationStatus struct {
	Active      bool
	Available   bool
}