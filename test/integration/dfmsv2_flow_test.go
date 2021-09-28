// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

func TestDriveV2FlowTransaction(t *testing.T) {
	config, err := client.Network.GetNetworkConfig(ctx)
	assert.Nil(t, err)

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

	const replicatorCount uint16 = 2
	var replicators [replicatorCount]*sdk.Account
	var storageSize uint64 = 500
	var streamingSize uint64 = 100
	var verificationFee = 100

	for i := replicatorCount; i != 0; {
		i--
		replicatorAccount, err := client.NewAccount()
		assert.Nil(t, err)
		fmt.Printf("replicatorAccount%d: %s\n", i, replicatorAccount)
		replicators[i] = replicatorAccount
	}

	// add storage and xpx mosaic to the replicator account

	var transfers [replicatorCount]*sdk.TransferTransaction
	for j := replicatorCount; j != 0; {
		j--
		transferMosaicsToReplicator, err := client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			replicators[j].Address,
			[]*sdk.Mosaic{sdk.Storage(storageSize), sdk.Streaming(streamingSize), sdk.Xpx(10000)},
			sdk.NewPlainMessage(""),
		)
		transfers[j] = transferMosaicsToReplicator
		transfers[j].ToAggregate(defaultAccount.PublicAccount)
		assert.Nil(t, err)
	}

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{transfers[0], transfers[1]},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)

	// end region

	// replicator onboarding transaction

	var rpOnboard [replicatorCount]*sdk.ReplicatorOnboardingTransaction
	for i := replicatorCount; i != 0; {
		// generate random BLS Public Key
		b := make([]byte, 32)
		_, err = rand.Read(b)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		var ikm [32]byte
		copy(ikm[:], b[:])
		sk := sdk.GenerateKeyPairFromIKM(ikm)
		blsKey := sk.PublicKey
		i--
		replicatorOnboardingTx, err := client.NewReplicatorOnboardingTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.Amount(50),
			blsKey,
		)
		rpOnboard[i] = replicatorOnboardingTx
		rpOnboard[i].ToAggregate(replicators[i].PublicAccount)
		assert.Nil(t, err)
		fmt.Printf("rpOnboard%d: %s\n", i, rpOnboard[i])
	}

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{rpOnboard[0], rpOnboard[1]},
		)
	}, replicators[0], replicators[1])
	assert.Nil(t, result.error)

	// end region

	// prepare bc drive transaction

	var ppBcDrive [replicatorCount]*sdk.PrepareBcDriveTransaction
	for j := replicatorCount; j != 0; {
		j--
		prepareBcDriveTx, err := client.NewPrepareBcDriveTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.StorageSize(storageSize/10),
			sdk.Amount(verificationFee),
			replicatorCount,
		)
		ppBcDrive[j] = prepareBcDriveTx
		ppBcDrive[j].ToAggregate(replicators[j].PublicAccount)
		assert.Nil(t, err)
		fmt.Printf("ppBcDrive%d: %s\n", j, ppBcDrive[j])
	}

	agTx, err := client.NewCompleteAggregateTransaction(
		sdk.NewDeadline(time.Hour),
		[]sdk.Transaction{ppBcDrive[0], ppBcDrive[1]},
	)
	assert.Nil(t, err)
	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return agTx, nil
	}, replicators[0], replicators[1])
	assert.Nil(t, result.error)

	// end region

	// drive closure transaction

	var drClosure [replicatorCount]*sdk.DriveClosureTransaction
	for i := replicatorCount; i != 0; {
		i--
		var driveKey = agTx.InnerTransactions[i].GetAbstractTransaction().UniqueAggregateHash.String()
		fmt.Println("driveKey: ", driveKey)
		driveClosureTx, err := client.NewDriveClosureTransaction(
			sdk.NewDeadline(time.Hour),
			driveKey,
		)
		drClosure[i] = driveClosureTx
		drClosure[i].ToAggregate(replicators[i].PublicAccount)
		assert.Nil(t, err)
		fmt.Printf("drClosure%d: %s\n", i, drClosure[i])
	}

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{drClosure[0], drClosure[1]},
		)
	}, replicators[0], replicators[1])
	assert.Nil(t, result.error)

	// end

}
