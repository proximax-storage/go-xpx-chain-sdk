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
	Status          OperationStatus
	Executor        *PublicAccount
	LockedMosaics   []*Mosaic
	Transactions    []*Transaction
}

type SuperContract struct {
	Account         *PublicAccount
	Drive           *Drive
	FileHash        *Hash
	FunctionsList   map[string]string
}

func (s *SuperContract) String() string {
	return fmt.Sprintf(
		`
			"Account": %s,
			"Drive": %s,
			"FileHash": %s,
			"FunctionsList": %+v,
		`,
		s.Account,
		s.Drive,
		s.FileHash,
		s.FunctionsList,
	)
}

type DeployTransaction struct {
	AbstractTransaction
	DriveAccount                *PublicAccount
	SuperContractAccount        *PublicAccount
	FileHash                    *Hash
	FunctionsList               map[string]string
}

type ExecuteTransaction struct {
	AbstractTransaction
	SuperContract       *PublicAccount
	LockMosaics         []*Mosaic
	Function            string
}

type StartOperationTransaction struct {
	AbstractTransaction
	OperationExecutor           *PublicAccount
	Mosaics                     []*Mosaic
}

// Must be aggregated in AOT
type EndOperationTransaction struct {
	AbstractTransaction
	UsedMosaics             []*Mosaic
	Status                  OperationStatus
}

type AggregateOperationTransaction struct {
	AbstractTransaction
	OperationHash       *Hash
	InnerTransactions   []Transaction
	Cosignatures        []*AggregateTransactionCosignature
}