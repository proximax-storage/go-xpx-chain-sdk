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
	exchangeAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(exchangeAccount)
	var exchangeAmount uint64 = 1000000

	config, err := client.Network.GetNetworkConfig(ctx)
	assert.Nil(t, err)

	config.NetworkConfig.Sections["plugin:catapult.plugins.exchange"].Fields["longOfferKey"].Value = exchangeAccount.PublicAccount.PublicKey

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewNetworkConfigTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.Duration(1),
			config.NetworkConfig,
			config.SupportedEntityVersions)
	}, nemesisAccount)
	assert.Nil(t, result.error)

	time.Sleep(time.Minute)

	driveAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(driveAccount)

	replicatorAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(replicatorAccount)

	var storageSize uint64 = 10000
	var billingPrice uint64 = 50

	driveTx, err := client.NewPrepareDriveTransaction(
		sdk.NewDeadline(time.Hour),
		defaultAccount.PublicAccount,
		sdk.Duration(1),
		sdk.Duration(1),
		sdk.Amount(billingPrice),
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

	transferXpxToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.Address,
		[]*sdk.Mosaic{sdk.Xpx(10000000)},
		sdk.NewPlainMessage(""),
	);
	transferXpxToReplicator.ToAggregate(defaultAccount.PublicAccount)
	assert.Nil(t, err)

	transferXpxToExchange, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		exchangeAccount.Address,
		[]*sdk.Mosaic{sdk.Storage(exchangeAmount)},
		sdk.NewPlainMessage(""),
	);
	transferXpxToExchange.ToAggregate(defaultAccount.PublicAccount)
	assert.Nil(t, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{driveTx, transferStorageToReplicator, transferXpxToReplicator, transferXpxToExchange},
		)
	}, defaultAccount, driveAccount)
	assert.Nil(t, result.error)

	if err := wsc.AddDriveStateHandlers(driveAccount.Address, func(info *sdk.DriveStateInfo) bool {
		if info.DriveKey != driveAccount.PublicAccount.PublicKey {
			return false
		}
		fmt.Printf("Got drive state: %v \n", info)
		return true
	}); err != nil {
		panic(err)
	}

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
		&sdk.Hash{1},
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
			&sdk.Hash{1},
			[]*sdk.AddAction{},
			[]*sdk.RemoveAction{
				{
					FileHash: fileHash,
				},
			},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	fmt.Println(defaultAccount)
	drives, err := client.Storage.GetAccountDrives(ctx, defaultAccount.PublicAccount, sdk.AllRoles)
	assert.Nil(t, err)
	fmt.Println(drives)

	drive, err := client.Storage.GetDrive(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	fmt.Println(drive)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewAddExchangeOfferTransaction(
			sdk.NewDeadline(time.Hour),
			[]*sdk.AddOffer{
				{
					sdk.Offer{
						sdk.SellOffer,
						sdk.Storage(exchangeAmount),
						sdk.Amount(exchangeAmount / 2),
					},
					sdk.Duration(10000000),
				},
			},
		)
	}, exchangeAccount)
	assert.Nil(t, result.error)

	exchangeInfo, err := client.Exchange.GetAccountExchangeInfo(ctx, exchangeAccount.PublicAccount)
	assert.Nil(t, err)
	fmt.Println(exchangeInfo)


	infos, err := client.Exchange.GetExchangeOfferByAssetId(ctx, sdk.StorageNamespaceId, sdk.SellOffer)
	assert.Nil(t, err)
	info := infos[0]
	confirmation, err := info.ConfirmOffer(sdk.Amount(billingPrice))
	assert.Nil(t, err)

	exchangeOfferTransaction, err := client.NewExchangeOfferTransaction(
		sdk.NewDeadline(time.Hour),
		[]*sdk.ExchangeConfirmation{
			confirmation,
		},
	)
	exchangeOfferTransaction.ToAggregate(driveAccount.PublicAccount)
	assert.Nil(t, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{exchangeOfferTransaction},
		)
	}, replicatorAccount)
	assert.Nil(t, result.error)

	drive, err = client.Storage.GetDrive(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	fmt.Println(drive)

	time.Sleep(time.Minute)


	endDriveTx, err := client.NewEndDriveTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.PublicAccount,
	)
	endDriveTx.ToAggregate(driveAccount.PublicAccount)
	assert.Nil(t, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{endDriveTx},
		)
	}, replicatorAccount)
	assert.Nil(t, result.error)

	drive, err = client.Storage.GetDrive(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	fmt.Println(drive)
}
