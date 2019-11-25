// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/str"
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

type BillingDescription struct {
	Start    Height
	End      Height
	Payments []*PaymentInformation
}

type ReplicatorInfo struct {
	Account             *PublicAccount
	Start               Height
	End                 Height
	Deposit             Amount
	Index               int
	FilesWithoutDeposit map[Hash]uint16
}

type FileActionType uint8

const (
	AddFile FileActionType = iota
	RemoveFile
)

type FileAction struct {
	Type 	FileActionType
	Height 	Height
}

type FileInfo struct {
	FileSize 	StorageSize
	Deposit  	Amount
	Payments 	[]*PaymentInformation
	Actions 	[]*FileAction
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
}

func (drive *Drive) String () string {
	return str.StructToString(
		"Drive",
		str.NewField("DriveKey", str.StringPattern, drive.DriveKey),
		str.NewField("Start", str.StringPattern, drive.Start),
		str.NewField("State", str.StringPattern, drive.State),
		str.NewField("Owner", str.StringPattern, drive.Owner),
		str.NewField("RootHash", str.StringPattern, drive.RootHash),
		str.NewField("Duration", str.StringPattern, drive.Duration),
		str.NewField("BillingPeriod", str.StringPattern, drive.BillingPeriod),
		str.NewField("BillingPrice", str.StringPattern, drive.BillingPrice),
		str.NewField("DriveSize", str.StringPattern, drive.DriveSize),
		str.NewField("Replicas", str.StringPattern, drive.Replicas),
		str.NewField("MinReplicators", str.StringPattern, drive.MinReplicators),
		str.NewField("PercentApprovers", str.StringPattern, drive.PercentApprovers),
		str.NewField("BillingHistory", str.StringPattern, drive.BillingHistory),
		str.NewField("Files", str.StringPattern, drive.Files),
		str.NewField("Replicators", str.StringPattern, drive.Replicators),
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

type AddAction struct {
	FileHash *Hash
	FileSize StorageSize
}

func (action *AddAction) String() string {
	return fmt.Sprintf(
		`
			"FileHash": %s,
			"FileSize": %s,
		`,
		action.FileHash,
		action.FileSize,
	)
}

type RemoveAction = File

type DriveFileSystemTransaction struct {
	AbstractTransaction
	DriveKey      *PublicAccount
	NewRootHash   *Hash
	OldRootHash   *Hash
	AddActions    []*AddAction
	RemoveActions []*RemoveAction
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

type DeletedFile struct {
	File
	UploadInfos []*UploadInfo
}

func (file *DeletedFile) String() string {
	return fmt.Sprintf(
		`
			"FileHash": %s,
			"UploadInfos": %s,
		`,
		file.FileHash,
		file.UploadInfos,
	)
}

// Delete Reward Transaction

type DeleteRewardTransaction struct {
	AbstractTransaction
	DeletedFiles []*DeletedFile
}
