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

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

func TestDriveV2FlowTransaction(t *testing.T) {
	const replicatorCount uint16 = 2
	var replicators [replicatorCount]*sdk.Account
	var storageSize uint64 = 500
	var streamingSize uint64 = 100
	var verificationFee = 100

	driveAccount, err := client.NewAccount()
	assert.Nil(t, err)
	fmt.Printf("driveAccount: %s\n", driveAccount)

	for i := replicatorCount; i != 0; {
		i--
		replicatorAccount, err := client.NewAccount()
		assert.Nil(t, err)
		fmt.Printf("replicatorAccount%d: %s\n", i, replicatorAccount)
		replicators[i] = replicatorAccount
	}

	// add storage and xpx mosaic to the drive account

	transferMosaicsToDrive, err := client.NewTransferTransaction(
		sdk.NewDeadline(time.Hour),
		driveAccount.Address,
		[]*sdk.Mosaic{sdk.Storage(storageSize / 10), sdk.Xpx(10000)},
		sdk.NewPlainMessage(""),
	)
	transferMosaicsToDrive.ToAggregate(defaultAccount.PublicAccount)
	assert.NoError(t, err, err)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{transferMosaicsToDrive},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	// end region

	// add storage, streaming and xpx mosaic to the replicator accounts

	transfers := make([]sdk.Transaction, replicatorCount)
	for j := replicatorCount; j != 0; {
		j--
		transferMosaicsToReplicator, err := client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			replicators[j].Address,
			[]*sdk.Mosaic{sdk.Storage(storageSize), sdk.Streaming(streamingSize), sdk.Xpx(10000)},
			sdk.NewPlainMessage(""),
		)
		transferMosaicsToReplicator.ToAggregate(defaultAccount.PublicAccount)
		assert.NoError(t, err, err)
		transfers[j] = transferMosaicsToReplicator
	}

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			transfers,
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	// end region

	// replicator onboarding transaction

	rpOnboards := make([]sdk.Transaction, replicatorCount)
	for i := replicatorCount; i != 0; {
		// generate random BLS Public Key
		b := make([]byte, 32)
		_, err := rand.Read(b)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		var ikm [32]byte
		copy(ikm[:], b[:])
		sk := sdk.GenerateKeyPairFromIKM(ikm)
		blsKey := sk.PublicKey.HexString()
		i--
		replicatorOnboardingTx, err := client.NewReplicatorOnboardingTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.Amount(storageSize),
			blsKey,
		)
		replicatorOnboardingTx.ToAggregate(replicators[i].PublicAccount)
		assert.NoError(t, err, err)
		rpOnboards[i] = replicatorOnboardingTx
		fmt.Printf("rpOnboard%d: %s\n", i, rpOnboards[i])
	}

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			rpOnboards,
		)
	}, replicators[0], replicators[1:]...)
	assert.Nil(t, result.error)

	// end region

	// prepare bc drive transaction

	prepareBcDriveTx, err := client.NewPrepareBcDriveTransaction(
		sdk.NewDeadline(time.Hour),
		sdk.StorageSize(storageSize/10),
		sdk.Amount(verificationFee),
		replicatorCount,
	)
	prepareBcDriveTx.ToAggregate(driveAccount.PublicAccount)
	assert.NoError(t, err, err)
	fmt.Printf("ppBcDrive: %s\n", prepareBcDriveTx)

	agTx, err := client.NewCompleteAggregateTransaction(
		sdk.NewDeadline(time.Hour),
		[]sdk.Transaction{prepareBcDriveTx},
	)
	assert.Nil(t, err)
	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return agTx, nil
	}, driveAccount)
	assert.Nil(t, result.error)

	// end region

	// drive closure transaction

	drClosures := make([]sdk.Transaction, len(agTx.InnerTransactions))
	for i, agTxIn := range agTx.InnerTransactions {
		var driveKey = strings.ToUpper(agTxIn.GetAbstractTransaction().UniqueAggregateHash.String())
		fmt.Println("driveKey: ", driveKey)
		driveClosureTx, err := client.NewDriveClosureTransaction(
			sdk.NewDeadline(time.Hour),
			driveKey,
		)
		driveClosureTx.ToAggregate(driveAccount.PublicAccount)
		assert.NoError(t, err, err)
		drClosures[i] = driveClosureTx
		fmt.Printf("drClosure%d: %s\n", i, drClosures[i])
	}
	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			drClosures,
		)
	}, driveAccount)
	assert.Nil(t, result.error)

	// end

}
