// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/stretchr/testify/assert"
)

func TestSuperContractFlowTransaction(t *testing.T) {
	driveAccount, err := client.NewAccountFromVersion(1)
	assert.Nil(t, err)
	fmt.Println(driveAccount)

	replicatorAccount, err := client.NewAccountFromVersion(1)
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
	assert.Nil(t, err)
	driveTx.ToAggregate(driveAccount)

	transferStorageToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		replicatorAccount.Address,
		[]*sdk.Mosaic{sdk.Storage(storageSize)},
		sdk.NewPlainMessage(""),
	)
	assert.Nil(t, err)
	transferStorageToReplicator.ToAggregate(defaultAccount)

	transferXpxToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.Address,
		[]*sdk.Mosaic{sdk.Xpx(10000000)},
		sdk.NewPlainMessage(""),
	)
	assert.Nil(t, err)
	transferXpxToReplicator.ToAggregate(defaultAccount)

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
		driveAccount.PublicAccount.PublicKey,
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
	assert.Nil(t, err)
	fsTx.ToAggregate(defaultAccount)

	transferStreamingToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		replicatorAccount.Address,
		[]*sdk.Mosaic{sdk.Streaming(fileSize)},
		sdk.NewPlainMessage(""),
	)
	assert.Nil(t, err)
	transferStreamingToReplicator.ToAggregate(defaultAccount)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{fsTx, transferStreamingToReplicator},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	superContract, err := client.NewAccountFromVersion(1)
	assert.Nil(t, err)
	deploy, err := client.NewDeployTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.PublicAccount,
		defaultAccount.PublicAccount,
		fileHash,
		123,
	)
	assert.Nil(t, err)
	deploy.ToAggregate(superContract)

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

	initiator, err := client.NewAccountFromVersion(1)
	assert.Nil(t, err)
	transferSCToInitiator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		initiator.Address,
		[]*sdk.Mosaic{sdk.SuperContractMosaic(1000)},
		sdk.NewPlainMessage(""),
	)
	assert.Nil(t, err)
	transferSCToInitiator.ToAggregate(defaultAccount)

	assert.Nil(t, err)
	execute, err := client.NewStartExecuteTransaction(
		sdk.NewDeadline(time.Hour),
		superContract.PublicAccount,
		[]*sdk.Mosaic{sdk.SuperContractMosaic(1000)},
		"GoGoGo",
		[]int64{},
	)
	assert.Nil(t, err)
	execute.ToAggregate(initiator)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{transferSCToInitiator, execute},
		)
	}, defaultAccount, initiator, replicatorAccount)
	assert.Nil(t, result.error)

	operationToken := result.Transaction.(*sdk.AggregateTransactionV1).InnerTransactions[1].GetAbstractTransaction().UniqueAggregateHash
	operation, err := client.SuperContract.GetOperation(ctx, operationToken)
	assert.Nil(t, err)

	operations, err := client.SuperContract.GetOperationsByAccount(ctx, initiator.PublicAccount)
	assert.Nil(t, err)
	assert.Equal(t, operation, operations[0])

	operationIdentify, err := client.NewOperationIdentifyTransaction(
		sdk.NewDeadline(time.Hour),
		operationToken,
	)
	assert.Nil(t, err)
	operationIdentify.ToAggregate(superContract)

	scFileHash, err := sdk.StringToHash("BA2D2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093")
	assert.Nil(t, err)
	scFs, err := client.NewSuperContractFileSystemTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.PublicAccount.PublicKey,
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
	assert.Nil(t, err)
	scFs.ToAggregate(superContract)

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
	assert.Nil(t, err)
	endExecute.ToAggregate(superContract)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{transferSCToInitiator, endExecute},
		)
	}, defaultAccount, replicatorAccount)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewDeactivateTransaction(
			sdk.NewDeadline(time.Hour),
			superContract.PublicAccount.PublicKey,
			driveAccount.PublicAccount.PublicKey,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	sc, err := client.SuperContract.GetSuperContract(ctx, superContract.PublicAccount)
	assert.Nil(t, err)
	assert.Equal(t, sdk.SuperContractDeactivatedByParticipant, sc.State)
}
