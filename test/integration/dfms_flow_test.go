// Copyright 2019 ProximaX Limited. All rights reserved.
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

func TestDriveFlowTransaction(t *testing.T) {
	exchangeAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Println(exchangeAccount)
	var exchangeAmount uint64 = 1000000

	config, err := client.Network.GetNetworkConfig(ctx)
	assert.Nil(t, err)

	config.NetworkConfig.Sections["plugin:catapult.plugins.exchange"].Fields["longOfferKey"].Value = exchangeAccount.PublicAccount.PublicKey

	configDelta := 2
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewNetworkConfigTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.Duration(configDelta),
			config.NetworkConfig,
			config.SupportedEntityVersions)
	}, nemesisAccount)
	assert.Nil(t, result.error)

	waitForBlocksCount(t, configDelta)

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
	assert.Nil(t, err)
	driveTx.ToAggregate(driveAccount.PublicAccount)

	transferStorageToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		replicatorAccount.Address,
		[]*sdk.Mosaic{sdk.Storage(storageSize)},
		sdk.NewPlainMessage(""),
	)
	assert.Nil(t, err)
	transferStorageToReplicator.ToAggregate(defaultAccount.PublicAccount)

	transferXpxToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.Address,
		[]*sdk.Mosaic{sdk.Xpx(10000000)},
		sdk.NewPlainMessage(""),
	)
	assert.Nil(t, err)
	transferXpxToReplicator.ToAggregate(defaultAccount.PublicAccount)

	transferXpxToExchange, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		exchangeAccount.Address,
		[]*sdk.Mosaic{sdk.Storage(exchangeAmount)},
		sdk.NewPlainMessage(""),
	)
	assert.Nil(t, err)
	transferXpxToExchange.ToAggregate(defaultAccount.PublicAccount)

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
	fsTx.ToAggregate(defaultAccount.PublicAccount)

	transferStreamingToReplicator, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		replicatorAccount.Address,
		[]*sdk.Mosaic{sdk.Streaming(fileSize)},
		sdk.NewPlainMessage(""),
	)
	assert.Nil(t, err)
	transferStreamingToReplicator.ToAggregate(defaultAccount.PublicAccount)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{fsTx, transferStreamingToReplicator},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	drives, err := client.Storage.GetAccountDrives(ctx, defaultAccount.PublicAccount, sdk.AllDriveRoles)
	assert.Nil(t, err)
	fmt.Println(drives)

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
	assert.Nil(t, err)
	exchangeOfferTransaction.ToAggregate(driveAccount.PublicAccount)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{exchangeOfferTransaction},
		)
	}, replicatorAccount)
	assert.Nil(t, result.error)

	verificationStatus, err := client.Storage.GetVerificationStatus(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	assert.False(t, verificationStatus.Active)
	assert.True(t, verificationStatus.Available)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStartDriveVerificationTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount.PublicAccount,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	verificationStatus, err = client.Storage.GetVerificationStatus(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	assert.True(t, verificationStatus.Active)
	assert.False(t, verificationStatus.Available)

	verificationTx, err := client.NewEndDriveVerificationTransaction(
		sdk.NewDeadline(time.Hour),
		[]*sdk.FailureVerification{},
	)
	assert.Nil(t, err)
	verificationTx.ToAggregate(driveAccount.PublicAccount)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{verificationTx},
		)
	}, replicatorAccount)
	assert.Nil(t, result.error)

	startFileDownloadTx, err := client.NewStartFileDownloadTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.PublicAccount,
		[]*sdk.DownloadFile{
			{
				FileHash: fileHash,
				FileSize: sdk.StorageSize(fileSize),
			},
		},
	)
	assert.Nil(t, err)
	startFileDownloadTx.ToAggregate(defaultAccount.PublicAccount)

	agTx, err := client.NewCompleteAggregateTransaction(
		sdk.NewDeadline(time.Hour),
		[]sdk.Transaction{startFileDownloadTx},
	)
	assert.Nil(t, err)
	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return agTx, nil
	}, defaultAccount)
	assert.Nil(t, result.error)

	uniqueHash, err := sdk.UniqueAggregateHash(agTx, startFileDownloadTx, client.GenerationHash())
	assert.Nil(t, err)
	donwloadInfo, err := client.Storage.GetDownloadInfo(ctx, uniqueHash)
	assert.Nil(t, err)

	donwloadInfos, err := client.Storage.GetAccountDownloadInfos(ctx, defaultAccount.PublicAccount)
	assert.Nil(t, err)
	assert.Equal(t, donwloadInfo, donwloadInfos[0])

	donwloadInfos, err = client.Storage.GetDriveDownloadInfos(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	assert.Equal(t, donwloadInfo, donwloadInfos[0])

	endFileDownloadTx, err := client.NewEndFileDownloadTransaction(
		sdk.NewDeadline(time.Hour),
		defaultAccount.PublicAccount,
		uniqueHash,
		[]*sdk.DownloadFile{
			{
				FileHash: fileHash,
				FileSize: sdk.StorageSize(fileSize),
			},
		},
	)
	assert.Nil(t, err)
	endFileDownloadTx.ToAggregate(driveAccount.PublicAccount)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{endFileDownloadTx},
		)
	}, replicatorAccount)
	assert.Nil(t, result.error)

	drive, err := client.Storage.GetDrive(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	fmt.Println(drive)

	waitForBlocksCount(t, billingPeriod)

	verificationStatus, err = client.Storage.GetVerificationStatus(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	assert.False(t, verificationStatus.Active)
	assert.False(t, verificationStatus.Available)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewDriveFileSystemTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount.PublicAccount.PublicKey,
			&sdk.Hash{},
			&sdk.Hash{1},
			[]*sdk.Action{},
			[]*sdk.Action{
				{
					FileHash: fileHash,
					FileSize: sdk.StorageSize(fileSize),
				},
			},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	drive, err = client.Storage.GetDrive(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)

	endDriveTx, err := client.NewEndDriveTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.PublicAccount,
	)
	assert.Nil(t, err)
	endDriveTx.ToAggregate(driveAccount.PublicAccount)

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

	driveFilesRewardTx, err := client.NewDriveFilesRewardTransaction(
		sdk.NewDeadline(time.Hour),
		[]*sdk.UploadInfo{
			{
				Participant:  replicatorAccount.PublicAccount,
				UploadedSize: 100,
			},
		},
	)
	assert.Nil(t, err)
	driveFilesRewardTx.ToAggregate(driveAccount.PublicAccount)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{driveFilesRewardTx},
		)
	}, replicatorAccount)
	assert.Nil(t, result.error)

	drive, err = client.Storage.GetDrive(ctx, driveAccount.PublicAccount)
	assert.Nil(t, err)
	fmt.Println(drive)
}
