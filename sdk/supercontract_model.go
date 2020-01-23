// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
)

type OperationStatus uint8

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
	Account   *PublicAccount
	Drive     *Drive
	FileHash  *Hash
	VMVersion uint64
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

type DeployTransaction struct {
	AbstractTransaction
	DriveAccount         *PublicAccount
	SuperContractAccount *PublicAccount
	FileHash             *Hash
	VMVersion            uint64
}

type ExecuteTransaction struct {
	AbstractTransaction
	SuperContract      *PublicAccount
	LockMosaics        []*Mosaic
	Function           string
	FunctionParameters []int64
}

type StartOperationTransaction struct {
	AbstractTransaction
	OperationExecutor *PublicAccount
	Mosaics           []*Mosaic
}

type OperationIdentifyTransaction struct {
	AbstractTransaction
	OperationHash     *Hash
}

// Must be aggregated in AOT
type EndOperationTransaction struct {
	AbstractTransaction
	UsedMosaics []*Mosaic
	Status      OperationStatus
}
