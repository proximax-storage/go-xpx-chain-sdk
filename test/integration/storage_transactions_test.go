// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"fmt"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

func TestDriveFlowTransaction(t *testing.T) {
	driveAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(driveAccount)

	replicatorAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(replicatorAccount)

	var storageSize uint64 = 10000

	driveTx, err := client.NewPrepareDriveTransaction(
		sdk.NewDeadline(time.Hour),
		defaultAccount.PublicAccount,
		sdk.Duration(100),
		sdk.Duration(50),
		sdk.Amount(50),
		sdk.StorageSize(storageSize),
		1,
		1,
		1,
	);
	driveTx.ToAggregate(driveAccount.PublicAccount)
	assert.Nil(t, err)

	transferStorageToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		replicatorAccount.Address,
		[]*sdk.Mosaic{sdk.Storage(storageSize)},
		sdk.NewPlainMessage(""),
	);
	transferStorageToReplicator.ToAggregate(defaultAccount.PublicAccount)
	assert.Nil(t, err)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{driveTx, transferStorageToReplicator},
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

	var fileSize uint64 = 50
	fileHash, err := sdk.StringToHash("AA2d2427E105A9B60DF634553849135DF629F1408A018D02B07A70CAFFB43093")
	assert.Nil(t, err)

	fsTx, err := client.NewDriveFileSystemTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.PublicAccount,
		&sdk.Hash{},
		&sdk.Hash{},
		[]*sdk.AddAction{
			{
				FileHash: fileHash,
				FileSize: sdk.StorageSize(fileSize),
			},
		},
		[]*sdk.RemoveAction{},
	)
	fsTx.ToAggregate(defaultAccount.PublicAccount)
	assert.Nil(t,err)

	transferStreamingToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		replicatorAccount.Address,
		[]*sdk.Mosaic{sdk.Streaming(fileSize)},
		sdk.NewPlainMessage(""),
	);
	transferStreamingToReplicator.ToAggregate(defaultAccount.PublicAccount)
	assert.Nil(t, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{fsTx, transferStreamingToReplicator},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewFilesDepositTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount.PublicAccount,
			[]*sdk.File{
				{
					FileHash: fileHash,
				},
			},
		)
	}, replicatorAccount)
	assert.Nil(t, result.error)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewDriveFileSystemTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount.PublicAccount,
			&sdk.Hash{},
			&sdk.Hash{},
			[]*sdk.AddAction{},
			[]*sdk.RemoveAction{
				{
					FileHash: fileHash,
				},
			},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	drives, err := client.Storage.GetAccountDrives(ctx, defaultAccount.PublicAccount, sdk.AllRoles)
	assert.Nil(t, err)
	fmt.Println(drives)
}
