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
	Id              *Hash
	Owner           *PublicAccount
	DownloadDataCdi *Hash
	UploadSize      StorageSize
}

func (desc *ActiveDataModification) String() string {
	return fmt.Sprintf(
		`
			"Id": %s,
			"Owner": %s,
			"DownloadDataCdi": %s,
			"UploadSize": %d,
		`,
		desc.Id,
		desc.Owner,
		desc.DownloadDataCdi,
		desc.UploadSize,
	)
}

type CompletedDataModification struct {
	ActiveDataModification *ActiveDataModification
	State                  DataModificationState
}

func (desc *CompletedDataModification) String() string {
	return fmt.Sprintf(
		`
			"ActiveDataModification:" %s,
			"State:" %d,
		`,
		desc.ActiveDataModification.String(),
		desc.State,
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
	ActiveDataModifications    []*ActiveDataModification
	CompletedDataModifications []*CompletedDataModification
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
		"ActiveDataModifications": %s,
		"CompletedDataModifications": %s,
		`,
		drive.BcDriveAccount,
		drive.OwnerAccount,
		drive.RootHash,
		drive.DriveSize,
		drive.UsedSize,
		drive.MetaFilesSize,
		drive.ReplicatorCount,
		drive.ActiveDataModifications,
		drive.CompletedDataModifications,
	)
}

type BcDrivesPage struct {
	BcDrives   []*BcDrive
	Pagination Pagination
}

type DriveInfo struct {
	LastApprovedDataModificationId *Hash
	DataModificationIdIsValid      bool
	InitialDownloadWork            uint64
	Index                          int
}

func (info *DriveInfo) String() string {
	return fmt.Sprintf(
		`
			"LastApprovedDataModificationId": %s,
			"DataModificationIdIsValid": %t,
			"InitialDownloadWork": %d,
		`,
		info.LastApprovedDataModificationId,
		info.DataModificationIdIsValid,
		info.InitialDownloadWork,
	)
}

type Replicator struct {
	ReplicatorKey *PublicAccount
	Version       int32
	Capacity      Amount
	BLSKey        BLSPublicKey
	Drives        map[string]*DriveInfo
}

func (replicator *Replicator) String() string {
	return fmt.Sprintf(
		`
		ReplicatorKey: %s, 
		Version: %d,
		Capacity: %d,
		BLSKey: %s,
		Drives: %s,
		`,
		replicator.ReplicatorKey,
		replicator.Version,
		replicator.Capacity,
		replicator.BLSKey,
		replicator.Drives,
	)
}

type ReplicatorsPage struct {
	Replicators []*Replicator
	Pagination  Pagination
}

// Replicator Offboarding Transaction
type ReplicatorOffboardingTransaction struct {
	AbstractTransaction
}