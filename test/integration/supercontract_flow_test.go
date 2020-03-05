// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"fmt"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSuperContractFlowTransaction(t *testing.T) {
	driveAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(driveAccount)

	replicatorAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(replicatorAccount)

	var storageSize uint64 = 10000
	var billingPrice uint64 = 50
	var billingPeriod = 10

	driveTx, err := client.NewPrepareDriveTransaction(
		sdk.NewDeadline(time.Hour),
		defaultAccount.PublicAccount,
		sdk.Duration(billingPeriod),
		sdk.Duration(billingPeriod),
		sdk.Amount(billingPrice),
		sdk.StorageSize(storageSize),
		1,
		1,
		1,
	)
	driveTx.ToAggregate(driveAccount.PublicAccount)
	assert.Nil(t, err)

	transferStorageToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		replicatorAccount.Address,
		[]*sdk.Mosaic{sdk.Storage(storageSize)},
		sdk.NewPlainMessage(""),
	)
	transferStorageToReplicator.ToAggregate(defaultAccount.PublicAccount)
	assert.Nil(t, err)

	transferXpxToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.Address,
		[]*sdk.Mosaic{sdk.Xpx(10000000)},
		sdk.NewPlainMessage(""),
	)
	transferXpxToReplicator.ToAggregate(defaultAccount.PublicAccount)
	assert.Nil(t, err)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{driveTx, transferStorageToReplicator, transferXpxToReplicator},
		)
	}, defaultAccount, driveAccount)
	assert.Nil(t, result.error)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewJoinToDriveTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount.PublicAccount,
		)
	}, replicatorAccount)
	assert.Nil(t, result.error)

	var fileSize uint64 = 147
	fileHash, err := sdk.StringToHash("AA2D2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093")
	assert.Nil(t, err)

	fsTx, err := client.NewDriveFileSystemTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.PublicAccount,
		&sdk.Hash{1},
		&sdk.Hash{},
		[]*sdk.Action{
			{
				FileHash: fileHash,
				FileSize: sdk.StorageSize(fileSize),
			},
		},
		[]*sdk.Action{},
	)
	fsTx.ToAggregate(defaultAccount.PublicAccount)
	assert.Nil(t, err)

	transferStreamingToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		replicatorAccount.Address,
		[]*sdk.Mosaic{sdk.Streaming(fileSize)},
		sdk.NewPlainMessage(""),
	)
	transferStreamingToReplicator.ToAggregate(defaultAccount.PublicAccount)
	assert.Nil(t, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{fsTx, transferStreamingToReplicator},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	superContract, err := client.NewAccount()
	assert.Nil(t, err)
	deploy, err := client.NewDeployTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.PublicAccount,
		defaultAccount.PublicAccount,
		fileHash,
		123,
	)
	deploy.ToAggregate(superContract.PublicAccount)
	assert.Nil(t, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{transferStreamingToReplicator, deploy},
		)
	}, defaultAccount, superContract)
	assert.Nil(t, result.error)

	contract, err := client.SuperContract.GetSuperContract(ctx, superContract.PublicAccount)
	assert.Nil(t, err)

	contracts, err := client.SuperContract.GetDriveSuperContracts(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	assert.Equal(t, contract, contracts[0])

	initiator, err := client.NewAccount()
	assert.Nil(t, err)
	transferSCToInitiator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		initiator.Address,
		[]*sdk.Mosaic{sdk.SuperContractMosaic(1000)},
		sdk.NewPlainMessage(""),
	)
	transferSCToInitiator.ToAggregate(defaultAccount.PublicAccount)

	assert.Nil(t, err)
	execute, err := client.NewStartExecuteTransaction(
		sdk.NewDeadline(time.Hour),
		superContract.PublicAccount,
		[]*sdk.Mosaic{sdk.SuperContractMosaic(1000)},
		"GoGoGo",
		[]int64{},
	)
	execute.ToAggregate(initiator.PublicAccount)
	assert.Nil(t, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{transferSCToInitiator, execute},
		)
	}, defaultAccount, initiator, replicatorAccount)
	assert.Nil(t, result.error)

	operationToken := result.Transaction.(*sdk.AggregateTransaction).InnerTransactions[1].GetAbstractTransaction().UniqueAggregateHash
	operation, err := client.SuperContract.GetOperation(ctx, operationToken)
	assert.Nil(t, err)

	operations, err := client.SuperContract.GetOperationsByAccount(ctx, initiator.PublicAccount)
	assert.Nil(t, err)
	assert.Equal(t, operation, operations[0])

	operationIdentify, err := client.NewOperationIdentifyTransaction(
		sdk.NewDeadline(time.Hour),
		operationToken,
	)
	operationIdentify.ToAggregate(superContract.PublicAccount)
	assert.Nil(t, err)

	scFileHash, err := sdk.StringToHash("BA2D2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093")
	assert.Nil(t, err)
	scFs, err := client.NewSuperContractFileSystemTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.PublicAccount,
		&sdk.Hash{2},
		&sdk.Hash{1},
		[]*sdk.Action{
			{
				FileHash: scFileHash,
				FileSize: sdk.StorageSize(5000),
			},
		},
		[]*sdk.Action{},
	)
	scFs.ToAggregate(superContract.PublicAccount)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{operationIdentify, transferSCToInitiator, scFs},
		)
	}, defaultAccount, replicatorAccount)

	endExecute, err := client.NewEndExecuteTransaction(
		sdk.NewDeadline(time.Hour),
		[]*sdk.Mosaic{sdk.SuperContractMosaic(1000)},
		operationToken,
		sdk.Success,
	)
	endExecute.ToAggregate(superContract.PublicAccount)
	assert.Nil(t, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{transferSCToInitiator, endExecute},
		)
	}, defaultAccount, replicatorAccount)
}
