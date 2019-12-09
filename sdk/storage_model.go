// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
)

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

func (info *PaymentInformation) String() string {
	return fmt.Sprintf(
		`{ "Receiver": %s, "Amount": %s, "Height": %s }`,
		info.Receiver,
		info.Amount,
		info.Height,
	)
}

type BillingDescription struct {
	Start    Height
	End      Height
	Payments []*PaymentInformation
}

func (desc *BillingDescription) String() string {
	return fmt.Sprintf(
		`
			"Start": %s,
			"End": %s,
			"Payments": %s,
		`,
		desc.Start,
		desc.End,
		desc.Payments,
	)
}

type ReplicatorInfo struct {
	Account                     *PublicAccount
	Start                       Height
	End                         Height
	Index                       int
	ActiveFilesWithoutDeposit   map[Hash]bool
}

func (info *ReplicatorInfo) String() string {
	return fmt.Sprintf(
		`
			"Account": %s,
			"Start": %s,
			"End": %s,
			"Index": %d,
			"ActiveFilesWithoutDeposit": %s,
		`,
		info.Account,
		info.Start,
		info.End,
		info.Index,
		info.ActiveFilesWithoutDeposit,
	)
}

type Drive struct {
	DriveAccount     *PublicAccount
	Start            Height
	State            DriveState
	OwnerAccount     *PublicAccount
	RootHash         *Hash
	Duration         Duration
	BillingPeriod    Duration
	BillingPrice     Amount
	DriveSize        StorageSize
	Replicas         uint16
	MinReplicators   uint16
	PercentApprovers uint8
	BillingHistory   []*BillingDescription
	Files            map[Hash]StorageSize
	Replicators      map[string]*ReplicatorInfo
	UploadPayments   []*PaymentInformation
}

func (drive *Drive) String() string {
	return fmt.Sprintf(
		`
			"DriveAccount": %s,
			"Start": %s,
			"State": %d,
			"OwnerAccount": %s,
			"RootHash": %s,
			"Duration": %d,
			"BillingPeriod": %d,
			"BillingPrice": %d,
			"DriveSize": %d,
			"Replicas": %d,
			"MinReplicators": %d,
			"PercentApprovers": %d,
			"BillingHistory": %s,
			"Files": %s,
			"Replicators": %s,
			"UploadPayments": %s,
		`,
		drive.DriveAccount,
		drive.Start,
		drive.State,
		drive.OwnerAccount,
		drive.RootHash,
		drive.Duration,
		drive.BillingPeriod,
		drive.BillingPrice,
		drive.DriveSize,
		drive.Replicas,
		drive.MinReplicators,
		drive.PercentApprovers,
		drive.BillingHistory,
		drive.Files,
		drive.Replicators,
		drive.UploadPayments,
	)
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