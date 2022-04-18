// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package integration

import (
	"testing"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLiquidityProviderTransaction(t *testing.T) {
	mosaicId, err := sdk.NewMosaicId(XPXID)
	assert.Nil(t, err)

	requiredAmount := uint64(10000000)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			managerAccount.Address,
			[]*sdk.Mosaic{sdk.Xpx(requiredAmount)},
			sdk.NewPlainMessage("Test"),
		)
	}, defaultAccount)
	require.Nil(t, result.error, result.error)

	slashingAccount, err := client.NewAccount()
	require.Nil(t, err, err)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCreateLiquidityProviderTransaction(
			sdk.NewDeadline(time.Hour),
			mosaicId,
			sdk.Amount(requiredAmount/2),
			sdk.Amount(10),
			100,
			5,
			slashingAccount.PublicAccount,
			10,
			10,
		)
	}, managerAccount)
	assert.Nil(t, result.error, result.error)

	lps, err := client.LiquidityProvider.GetLiquidityProviders(ctx, nil)
	assert.Nil(t, err, err)
	require.NotNil(t, lps.LiquidityProviders)

	expectedLp := lps.LiquidityProviders[0]
	lp, err := client.LiquidityProvider.GetLiquidityProvider(ctx, expectedLp.ProviderKey)
	assert.Nil(t, err, err)
	assert.NotNil(t, lp)
	assert.EqualValues(t, expectedLp, lp)
}

func TestManualRateChangeTransactionTransaction(t *testing.T) {
	//TODO transaction passes successfully, but there is strange behavior with bool unmarshalling
	t.SkipNow()

	mosaicId, err := sdk.NewMosaicId(XPXID)
	assert.Nil(t, err)

	requiredAmount := uint64(100)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			managerAccount.Address,
			[]*sdk.Mosaic{sdk.Xpx(requiredAmount)},
			sdk.NewPlainMessage("Test"),
		)
	}, defaultAccount)
	require.Nil(t, result.error)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewManualRateChangeTransaction(
			sdk.NewDeadline(time.Hour),
			mosaicId,
			true,
			sdk.Amount(requiredAmount),
			true,
			sdk.Amount(300),
		)
	}, managerAccount)
	assert.Nil(t, result.error)
}
