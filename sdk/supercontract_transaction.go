// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"errors"
	"fmt"
)

func NewDeployTransaction(deadline *Deadline, drive, supercontract *PublicAccount, fileHash *Hash,
		functionsList []string, networkType NetworkType) (*DeployTransaction, error) {

	if drive == nil {
		return nil, ErrNilAccount
	}

	if supercontract == nil {
		return nil, ErrNilAccount
	}

	tx := DeployTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     DeployVersion,
			Deadline:    deadline,
			Type:        Deploy,
			NetworkType: networkType,
		},
		DriveAccount:           drive,
		SuperContractAccount:   supercontract,
		FileHash:               fileHash,
		FunctionsList:          functionsList,
	}

	return &tx, nil
}

func (tx *DeployTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *DeployTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"DriveAccount": %s,
			"SuperContractAccount": %s,
			"FileHash": %s,
			"FunctionsList": %+v,
		`,
		tx.AbstractTransaction.String(),
		tx.DriveAccount,
		tx.SuperContractAccount,
		tx.FileHash,
		tx.FunctionsList,
	)
}

func (tx *DeployTransaction) Size() int {
	return 0
}

func (tx *DeployTransaction) Bytes() ([]byte, error) {
	return nil, nil
}

func NewExecuteTransaction(deadline *Deadline, supercontract *PublicAccount, mosaics []*Mosaic, function string, networkType NetworkType) (*ExecuteTransaction, error) {
	if supercontract == nil {
		return nil, ErrNilAccount
	}
	if len(function) == 0 {
		return nil, errors.New("Function should be not empty")
	}
	if len(mosaics) == 0 {
		return nil, errors.New("Mosaics should be not empty")
	}

	tx := ExecuteTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     ExecuteVersion,
			Deadline:    deadline,
			Type:        Execute,
			NetworkType: networkType,
		},
		SuperContract:      supercontract,
		LockMosaics:        mosaics,
		Function:           function,
	}

	return &tx, nil
}

func (tx *ExecuteTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ExecuteTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"SuperContract": %s,
			"LockMosaics": %s,
			"Function": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.SuperContract,
		tx.LockMosaics,
		tx.Function,
	)
}

func (tx *ExecuteTransaction) Size() int {
	return 0
}

func (tx *ExecuteTransaction) Bytes() ([]byte, error) {
	return nil, nil
}

func NewAggregateOperationBoundedTransaction(deadline *Deadline, innerTxs []Transaction, hash *Hash, networkType NetworkType) (*AggregateOperationTransaction, error) {
	if innerTxs == nil {
		return nil, errors.New("innerTransactions must not be nil")
	}

	tx := AggregateOperationTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     AggregateOperationBondedVersion,
			Deadline:    deadline,
			Type:        AggregateOperationBonded,
			NetworkType: networkType,
		},
		OperationHash:        hash,
		InnerTransactions:    innerTxs,
	}

	return &tx, nil
}

func NewAggregateOperationCompleteTransaction(deadline *Deadline, innerTxs []Transaction, hash *Hash, networkType NetworkType) (*AggregateOperationTransaction, error) {
	if innerTxs == nil {
		return nil, errors.New("innerTransactions must not be nil")
	}

	tx := AggregateOperationTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     AggregateOperationCompletedVersion,
			Deadline:    deadline,
			Type:        AggregateOperationCompleted,
			NetworkType: networkType,
		},
		OperationHash:        hash,
		InnerTransactions:    innerTxs,
	}

	return &tx, nil
}

func (tx *AggregateOperationTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *AggregateOperationTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"OperationHash": %s,
			"InnerTransactions": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.OperationHash,
		tx.InnerTransactions,
	)
}

func (tx *AggregateOperationTransaction) Size() int {
	return 0
}

func (tx *AggregateOperationTransaction) Bytes() ([]byte, error) {
	return nil, nil
}

func NewEndOperationTransaction(deadline *Deadline, mosaics []*Mosaic, status OperationStatus, networkType NetworkType) (*EndOperationTransaction, error) {
	if status == Unknown {
		return nil, errors.New("Status should be not unknown")
	}
	if len(mosaics) == 0 {
		return nil, errors.New("Mosaics should be not empty")
	}

	tx := EndOperationTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     EndOperationVersion,
			Deadline:    deadline,
			Type:        EndOperation,
			NetworkType: networkType,
		},
		UsedMosaics:    mosaics,
		Status:         status,
	}

	return &tx, nil
}

func (tx *EndOperationTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *EndOperationTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"Status": %d,
			"UsedMosaics": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.Status,
		tx.UsedMosaics,
	)
}

func (tx *EndOperationTransaction) Size() int {
	return 0
}

func (tx *EndOperationTransaction) Bytes() ([]byte, error) {
	return nil, nil
}
