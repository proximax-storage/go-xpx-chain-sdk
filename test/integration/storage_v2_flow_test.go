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
	var replicatorsBlsKeys [replicatorCount]*sdk.KeyPair
	var storageSize uint64 = 500
	var verificationFee = 100

	owner, err := client.NewAccount()
	require.NoError(t, err, err)
	fmt.Printf("owner: %s\n", owner)

	for i := 0; i < len(replicators); i++ {
		replicators[i], err = client.NewAccount()
		require.NoError(t, err, err)
		fmt.Printf("replicatorAccount[%d]: %s\n", i, replicators[i])

		var ikm [32]byte
		_, err = rand.Read(ikm[:])
		require.NoError(t, err, err)

		replicatorsBlsKeys[i] = sdk.GenerateKeyPairFromIKM(ikm)
		fmt.Printf("replicatorsBlsKeys[%d]: %s\n", i, replicatorsBlsKeys[i])
	}

	// add storage and xpx mosaic to the drive account

	transferMosaicsToDrive, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		owner.Address,
		[]*sdk.Mosaic{sdk.Storage(storageSize / 10), sdk.Xpx(10000)},
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
	assert.NoError(t, result.error, result.error)

	// end region

	// add storage, streaming and xpx mosaic to the replicator accounts

	transfers := make([]sdk.Transaction, replicatorCount)
	for i := 0; i < len(replicators); i++ {
		transferMosaicsToReplicator, err := client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			replicators[i].Address,
			[]*sdk.Mosaic{sdk.Storage(storageSize), sdk.Xpx(10000)},
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
	assert.NoError(t, result.error, result.error)

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
	assert.NoError(t, result.error, result.error)

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
	assert.NoError(t, result.error, result.error)

	driveKey := strings.ToUpper(result.Transaction.GetAbstractTransaction().TransactionHash.String())
	driveAcc, err := sdk.NewAccountFromPublicKey(driveKey, client.NetworkType())
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
				driveAcc,
				block.BlockHash, // TODO get a real verificationTrigger
				1,
				keys,
				signatures,
				opinions,
			)
		}, defaultAccount)
		assert.NoError(t, result.error, result.error)
	})

	// drive closure transaction

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewDriveClosureTransaction(
			sdk.NewDeadline(time.Hour),
			driveKey,
		)
	}, owner)
	assert.NoError(t, result.error, result.error)

	// end region

	// replicator offboarding transaction

	for i := 0; i < len(replicators); i++ {
		result = sendTransaction(t, func() (sdk.Transaction, error) {
			return client.NewReplicatorOffboardingTransaction(sdk.NewDeadline(time.Hour))
		}, replicators[i])
		assert.Nil(t, result.error)
	}

	// end region

}
