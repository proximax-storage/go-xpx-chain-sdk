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
	ActiveDataModification *ActiveDataModification
	State                  DataModificationState
}

func (completed *CompletedDataModification) String() string {
	return fmt.Sprintf(
		`
			"ActiveDataModification": %s,
			"State:" %d,
		`,
		completed.ActiveDataModification.String(),
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
		drive.BcDriveAccount.String(),
		drive.OwnerAccount.String(),
		drive.RootHash.String(),
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
			"Index": %d
		`,
		info.LastApprovedDataModificationId.String(),
		info.DataModificationIdIsValid,
		info.InitialDownloadWork,
		info.Index,
	)
}

type Replicator struct {
	ReplicatorAccount *PublicAccount
	Version           int32
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
		replicator.ReplicatorAccount.String(),
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

type VerificationResult struct {
	Prover *PublicAccount
	Result bool
}

func (vr *VerificationResult) String() string {
	return fmt.Sprintf(
		`
		"Prover": %s,
		"Result": %t
		`,
		vr.Prover.String(),
		vr.Result,
	)
}

type VerificationResults []*VerificationResult

func (vrs VerificationResults) String() string {
	var str string
	for _, vr := range vrs {
		str += vr.String()
	}

	return str
}

func (vrs VerificationResults) Size() int {
	return len(vrs) * (KeySize + 1)
}

type VerificationOpinion struct {
	Verifier     *PublicAccount
	BlsSignature BLSSignature
	Results      VerificationResults
}

func (vo *VerificationOpinion) String() string {
	return fmt.Sprintf(
		`
		"Verifier": %s,
		"BlsSignature": %s,
		"Results": %s
		`,
		vo.Verifier.String(),
		vo.BlsSignature.HexString(),
		vo.Results.String(),
	)
}

func (vo *VerificationOpinion) Size() int {
	return KeySize + BlsSignatureSize + vo.Results.Size()
}

type EndDriveVerificationTransactionV2 struct {
	AbstractTransaction
	DriveKey             *PublicAccount
	VerificationTrigger  *Hash
	Provers              []*PublicAccount
	VerificationOpinions []*VerificationOpinion
}
