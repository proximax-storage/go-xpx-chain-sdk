// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"net/url"
	"strconv"
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

func (active *ActiveDataModification) String() string {
	return fmt.Sprintf(
		`
			"Id": %s,
			"Owner": %s,
			"DownloadDataCdi": %s,
			"UploadSize": %d,
		`,
		active.Id.String(),
		active.Owner.String(),
		active.DownloadDataCdi.String(),
		active.UploadSize,
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
		"ActiveDataModifications": %+v,
		"CompletedDataModifications": %+v,
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

type BcDrivesPageOptions struct {
	BcDrivesPageFilters
	PaginationOrderingOptions
}

type BcDrivesPageFilters struct {
	Size            BcDrivesValue64 `url:""`
	UsedSize        BcDrivesValue64 `url:""`
	MetaFilesSize   BcDrivesValue64 `url:""`
	ReplicatorCount BcDrivesValue16 `url:""`
}

type BcDrivesValue64 struct {
	Value64           uint64
	BcDrivesValueType BcDrivesValueType
}

func (sV BcDrivesValue64) EncodeValues(key string, v *url.Values) error {
	if Size == sV.BcDrivesValueType {
		v.Add(Size.String(), strconv.FormatUint(sV.Value64, 10))
	} else if FromSize == sV.BcDrivesValueType {
		u := uint64DTO(uint64ToArray(sV.Value64))
		v.Add(FromSize.String(), u.toStruct().String())
	} else if ToSize == sV.BcDrivesValueType {
		u := uint64DTO(uint64ToArray(sV.Value64))
		v.Add(ToSize.String(), u.toStruct().String())
	} else if UsedSize == sV.BcDrivesValueType {
		v.Add(UsedSize.String(), strconv.FormatUint(sV.Value64, 10))
	} else if FromUsedSize == sV.BcDrivesValueType {
		u := uint64DTO(uint64ToArray(sV.Value64))
		v.Add(FromUsedSize.String(), u.toStruct().String())
	} else if ToUsedSize == sV.BcDrivesValueType {
		u := uint64DTO(uint64ToArray(sV.Value64))
		v.Add(ToUsedSize.String(), u.toStruct().String())
	} else if MetaFilesSize == sV.BcDrivesValueType {
		v.Add(MetaFilesSize.String(), strconv.FormatUint(sV.Value64, 10))
	} else if FromMetaFilesSize == sV.BcDrivesValueType {
		u := uint64DTO(uint64ToArray(sV.Value64))
		v.Add(FromMetaFilesSize.String(), u.toStruct().String())
	} else if ToMetaFilesSize == sV.BcDrivesValueType {
		u := uint64DTO(uint64ToArray(sV.Value64))
		v.Add(ToMetaFilesSize.String(), u.toStruct().String())
	}

	return nil
}

type BcDrivesValue16 struct {
	Value16           uint16
	BcDrivesValueType BcDrivesValueType
}

func (sV BcDrivesValue16) EncodeValues(key string, v *url.Values) error {
	if ReplicatorCount == sV.BcDrivesValueType {
		v.Add(ReplicatorCount.String(), strconv.FormatUint(uint64(sV.Value16), 10))
	} else if FromReplicatorCount == sV.BcDrivesValueType {
		v.Add(FromReplicatorCount.String(), strconv.FormatUint(uint64(sV.Value16), 10))
	} else if ToReplicatorCount == sV.BcDrivesValueType {
		v.Add(ToReplicatorCount.String(), strconv.FormatUint(uint64(sV.Value16), 10))
	}

	return nil
}

type BcDrivesValueType string

const (
	Size                BcDrivesValueType = "size"
	FromSize            BcDrivesValueType = "fromSize"
	ToSize              BcDrivesValueType = "toSize"
	UsedSize            BcDrivesValueType = "usedSize"
	FromUsedSize        BcDrivesValueType = "fromUsedSize"
	ToUsedSize          BcDrivesValueType = "toUsedSize"
	MetaFilesSize       BcDrivesValueType = "metaFilesSize"
	FromMetaFilesSize   BcDrivesValueType = "fromMetaFilesSize"
	ToMetaFilesSize     BcDrivesValueType = "toMetaFilesSize"
	ReplicatorCount     BcDrivesValueType = "replicatorCount"
	FromReplicatorCount BcDrivesValueType = "fromReplicatorCount"
	ToReplicatorCount   BcDrivesValueType = "toReplicatorCount"
)

func (vT BcDrivesValueType) String() string {
	return string(vT)
}

type DriveInfo struct {
	Drive                          *PublicAccount
	LastApprovedDataModificationId *Hash
	DataModificationIdIsValid      bool
	InitialDownloadWork            uint64
	Index                          int
}

func (info *DriveInfo) String() string {
	return fmt.Sprintf(
		`
			"Drive": %s, 
		    "LastApprovedDataModificationId": %s,
			"DataModificationIdIsValid": %t,
			"InitialDownloadWork": %d,
			"Index": %d
		`,
		info.Drive,
		info.LastApprovedDataModificationId,
		info.DataModificationIdIsValid,
		info.InitialDownloadWork,
		info.Index,
	)
}

type Replicator struct {
	ReplicatorAccount *PublicAccount
	Version           uint32
	Capacity          Amount
	BLSKey            string
	Drives            map[string]*DriveInfo
}

func (replicator *Replicator) String() string {
	return fmt.Sprintf(
		`
		ReplicatorAccount: %s, 
		Version: %d,
		Capacity: %d,
		BLSKey: %s,
		Drives: %+v,
		`,
		replicator.ReplicatorAccount,
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

type ReplicatorsPageOptions struct {
	ReplicatorsPageFilters
	PaginationOrderingOptions
}

type ReplicatorsPageFilters struct {
	Version  ReplicatorsValue32 `url:""`
	Capacity ReplicatorsValue64 `url:""`
}

type ReplicatorsValue64 struct {
	Value64              uint64
	ReplicatorsValueType ReplicatorsValueType
}

func (sV ReplicatorsValue64) EncodeValues(key string, v *url.Values) error {
	if Capacity == sV.ReplicatorsValueType {
		v.Add(Capacity.String(), strconv.FormatUint(sV.Value64, 10))
	} else if FromCapacity == sV.ReplicatorsValueType {
		u := uint64DTO(uint64ToArray(sV.Value64))
		v.Add(FromCapacity.String(), u.toStruct().String())
	} else if ToCapacity == sV.ReplicatorsValueType {
		u := uint64DTO(uint64ToArray(sV.Value64))
		v.Add(ToCapacity.String(), u.toStruct().String())
	}

	return nil
}

type ReplicatorsValue32 struct {
	Value32              uint32
	ReplicatorsValueType ReplicatorsValueType
}

func (sV ReplicatorsValue32) EncodeValues(key string, v *url.Values) error {
	if Version == sV.ReplicatorsValueType {
		v.Add(Version.String(), strconv.FormatUint(uint64(sV.Value32), 10))
	} else if FromVersion == sV.ReplicatorsValueType {
		v.Add(FromVersion.String(), strconv.FormatUint(uint64(sV.Value32), 10))
	} else if ToVersion == sV.ReplicatorsValueType {
		v.Add(ToVersion.String(), strconv.FormatUint(uint64(sV.Value32), 10))
	}

	return nil
}

type ReplicatorsValueType string

const (
	Version      ReplicatorsValueType = "version"
	FromVersion  ReplicatorsValueType = "fromVersion"
	ToVersion    ReplicatorsValueType = "toVersion"
	Capacity     ReplicatorsValueType = "capacity"
	FromCapacity ReplicatorsValueType = "fromCapacity"
	ToCapacity   ReplicatorsValueType = "toCapacity"
)

func (vT ReplicatorsValueType) String() string {
	return string(vT)
}

// Replicator Onboarding Transaction
type ReplicatorOnboardingTransaction struct {
	AbstractTransaction
	Capacity     Amount
	BlsPublicKey string
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
	DriveKey string
}
