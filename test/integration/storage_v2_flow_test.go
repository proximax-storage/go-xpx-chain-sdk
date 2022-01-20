// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"crypto/rand"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

func TestDriveV2FlowTransaction(t *testing.T) {
	const replicatorCount uint16 = 2
	var replicators [replicatorCount]*sdk.Account
	var storageSize uint64 = 500
	var verificationFee = 100

	owner, err := client.NewAccount()
	require.NoError(t, err, err)
	fmt.Printf("owner: %s\n", owner)

	for i := 0; i < len(replicators); i++ {
		replicators[i], err = client.NewAccount()
		require.NoError(t, err, err)
		fmt.Printf("replicatorAccount[%d]: %s\n", i, replicators[i])
	}

	// add storage and xpx mosaic to the drive owner
	transferMosaicsToDrive, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		owner.Address,
		[]*sdk.Mosaic{sdk.Storage(storageSize), sdk.Streaming(storageSize * 2), sdk.Xpx(10000)},
		sdk.NewPlainMessage(""),
	)
	assert.NoError(t, err, err)
	transferMosaicsToDrive.ToAggregate(defaultAccount.PublicAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{transferMosaicsToDrive},
		)
	}, defaultAccount)
	require.NoError(t, result.error, result.error)

	// end region

	// add storage, streaming and xpx mosaic to the replicator accounts
	transfers := make([]sdk.Transaction, replicatorCount)
	for i := 0; i < len(replicators); i++ {
		transferMosaicsToReplicator, err := client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			replicators[i].Address,
			[]*sdk.Mosaic{sdk.Storage(storageSize), sdk.Streaming(storageSize * 2), sdk.Xpx(10000)},
			sdk.NewPlainMessage(""),
		)
		assert.NoError(t, err, err)

		transferMosaicsToReplicator.ToAggregate(defaultAccount.PublicAccount)
		transfers[i] = transferMosaicsToReplicator
	}

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			transfers,
		)
	}, defaultAccount)
	require.NoError(t, result.error, result.error)

	// end region

	// replicator onboarding transaction

	rpOnboards := make([]sdk.Transaction, replicatorCount)
	for i := 0; i < len(replicators); i++ {
		replicatorOnboardingTx, err := client.NewReplicatorOnboardingTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.Amount(storageSize),
		)
		assert.NoError(t, err, err)
		replicatorOnboardingTx.ToAggregate(replicators[i].PublicAccount)
		rpOnboards[i] = replicatorOnboardingTx
		fmt.Printf("rpOnboard%d: %s\n", i, rpOnboards[i])
	}

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			rpOnboards,
		)
	}, replicators[0], replicators[1:]...)
	require.NoError(t, result.error, result.error)

	// end region

	// prepare bc drive transaction
	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewPrepareBcDriveTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.StorageSize(storageSize/10),
			sdk.Amount(verificationFee),
			replicatorCount,
		)
	}, owner)
	require.NoError(t, result.error, result.error)

	driveKey := strings.ToUpper(result.Transaction.GetAbstractTransaction().TransactionHash.String())
	driveAccount, err := sdk.NewAccountFromPublicKey(driveKey, client.NetworkType())
	assert.NoError(t, err, err)
	fmt.Printf("Drive Account: %s", driveAccount.String())

	// end region

	// Data Modification

	downloadDataCdi := &sdk.Hash{}
	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewDataModificationTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount,
			downloadDataCdi,
			10,
			10,
		)
	}, owner)
	assert.NoError(t, result.error, result.error)

	// end region

	// Data Modification Cancel

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewDataModificationCancelTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount,
			downloadDataCdi,
		)
	}, owner)
	assert.NoError(t, result.error, result.error)

	// end region

	// Storage Payment

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStoragePaymentTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount,
			10,
		)
	}, owner)
	assert.NoError(t, result.error, result.error)

	// end region

	// Download

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewDownloadTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount,
			sdk.StorageSize(storageSize),
			100,
			[]*sdk.PublicAccount{},
		)
	}, owner)
	assert.NoError(t, result.error, result.error)

	// end region

	// Download Payment

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewDownloadPaymentTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount,
			10,
			10,
		)
	}, owner)
	assert.NoError(t, result.error, result.error)

	// end region

	// Finish Download

	downloadChannelId := &sdk.Hash{1} // TODO add real downloadChannelId
	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewFinishDownloadTransaction(
			sdk.NewDeadline(time.Hour),
			downloadChannelId,
			100,
		)
	}, owner)
	assert.NoError(t, result.error, result.error)

	// end region

	// VerificationPayment

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewVerificationPaymentTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount,
			100,
		)
	}, owner)
	assert.NoError(t, result.error, result.error)

	// end region

	t.Run("EndDriveVerificationV2", func(t *testing.T) {
		//t.SkipNow()

		// prepare same opinions
		opinions := make([]uint8, len(replicators)*len(replicators))
		for i, _ := range opinions {
			opinions[i] = 1
		}

		keys := make([]*sdk.PublicAccount, len(replicators))
		for i, r := range replicators {
			keys[i] = r.PublicAccount
		}

		signatures := make([]string, len(replicators))
		for i, _ := range replicators {
			var s [64]byte
			_, err = rand.Read(s[:])
			signatures[i] = string(s[:])
		}

		currHeight, err := client.Blockchain.GetBlockchainHeight(ctx)
		require.NoError(t, err, err)

		block, err := client.Blockchain.GetBlockByHeight(ctx, currHeight)
		require.NoError(t, err, err)

		result = sendTransaction(t, func() (sdk.Transaction, error) {
			return client.NewEndDriveVerificationTransactionV2(
				sdk.NewDeadline(time.Hour),
				driveAccount,
				block.BlockHash, // TODO get a real verificationTrigger
				1,
				keys,
				signatures,
				opinions,
			)
		}, defaultAccount)
		assert.NoError(t, result.error, result.error)
	})

	// replicator offboarding transaction

	for i := 0; i < len(replicators); i++ {
		result = sendTransaction(t, func() (sdk.Transaction, error) {
			return client.NewReplicatorOffboardingTransaction(
				sdk.NewDeadline(time.Hour),
				driveAccount,
			)
		}, replicators[i])
		assert.NoError(t, result.error, result.error)
	}

	// end region

	// drive closure transaction

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewDriveClosureTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount,
		)
	}, owner)
	require.NoError(t, result.error, result.error)

	// end region
}
