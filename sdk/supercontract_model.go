// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
)

type OperationStatus uint16

const (
	Unknown OperationStatus = iota
	Started
	Success
	Failure
)

type Operation struct {
	Status        OperationStatus
	Executor      *PublicAccount
	LockedMosaics []*Mosaic
	Transactions  []*Transaction
}

type SuperContract struct {
	Account     *PublicAccount
	Drive       *PublicAccount
	FileHash    *Hash
	VMVersion   uint64
	Start       Height
	End         Height
}

func (s *SuperContract) String() string {
	return fmt.Sprintf(
		`
			"Account": %s,
			"Drive": %s,
			"FileHash": %s,
			"VMVersion": %d,
		`,
		s.Account,
		s.Drive,
		s.FileHash,
		s.VMVersion,
	)
}

type StartOperationTransaction struct {
	AbstractTransaction
	OperationExecutors  []*PublicAccount
	Mosaics             []*Mosaic
	Duration            Duration
}

type OperationIdentifyTransaction struct {
	AbstractTransaction
	OperationHash     *Hash
}

// Must be aggregated in AOT
type EndOperationTransaction struct {
	AbstractTransaction
	UsedMosaics         []*Mosaic
	OperationToken      *Hash
	Status              OperationStatus
}

type DeployTransaction struct {
	AbstractTransaction
	DriveAccount            *PublicAccount
	Owner                   *PublicAccount
	FileHash                *Hash
	VMVersion               uint64
}

type StartExecuteTransaction struct {
	AbstractTransaction
	SuperContract      *PublicAccount
	Function           string
	LockMosaics        []*Mosaic
	FunctionParameters []int64
}

type EndExecuteTransaction = EndOperationTransaction