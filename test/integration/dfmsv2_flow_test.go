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

	driveAccount, err := client.NewAccountFromPrivateKey("EDFB348D4AAA333E6D73D9CAD1EA18FE3FE079CC3373E9E4E75A4FBD7D3476E0")
	assert.Nil(t, err)
	fmt.Printf("driveAccount: %s\n", driveAccount)

	for i := 0; i < len(replicators); i++ {
		replicatorAccount, err := client.NewAccount()
		assert.Nil(t, err)
		fmt.Printf("replicatorAccount%d: %s\n", i, replicatorAccount)
		replicators[i] = replicatorAccount
	}

	// add storage and xpx mosaic to the drive account

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			driveAccount.Address,
			[]*sdk.Mosaic{sdk.Storage(storageSize / 10), sdk.Xpx(10000)},
			sdk.NewPlainMessage(""),
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	// end region

	// add storage, streaming and xpx mosaic to the replicator accounts

	transfers := make([]sdk.Transaction, replicatorCount)
	for i := 0; i < len(replicators); i++ {
		transferMosaicsToReplicator, err := client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			replicators[i].Address,
			[]*sdk.Mosaic{sdk.Storage(storageSize), sdk.Streaming(streamingSize), sdk.Xpx(10000)},
			sdk.NewPlainMessage(""),
		)
		transferMosaicsToReplicator.ToAggregate(defaultAccount.PublicAccount)
		assert.NoError(t, err, err)
		transfers[i] = transferMosaicsToReplicator
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
	for i := 0; i < len(replicators); i++ {
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
	assert.NoError(t, err, err)
	fmt.Printf("ppBcDrive: %s\n", prepareBcDriveTx)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return prepareBcDriveTx, nil
	}, driveAccount)
	assert.Nil(t, result.error)

	// end region

	// drive closure transaction

	driveKey := strings.ToUpper(prepareBcDriveTx.GetAbstractTransaction().TransactionHash.String())
	fmt.Println("driveKey: ", driveKey)
	driveClosureTx, err := client.NewDriveClosureTransaction(
		sdk.NewDeadline(time.Hour),
		driveKey,
	)
	assert.NoError(t, err, err)
	fmt.Printf("drClosure: %s\n", driveClosureTx)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return driveClosureTx, nil
	}, driveAccount)
	assert.Nil(t, result.error)

	// end region

	// replicator offboarding transaction

	for i := 0; i < len(replicators); i++ {
		replicatorOffboardingTx, err := client.NewReplicatorOffboardingTransaction(sdk.NewDeadline(time.Hour))
		assert.Nil(t, err)
		fmt.Printf("rpOffboard%d: %s\n", i, replicatorOffboardingTx)

		result = sendTransaction(t, func() (sdk.Transaction, error) {
			return replicatorOffboardingTx, nil
		}, replicators[i])
		assert.Nil(t, result.error)
	}

	// end region

}
