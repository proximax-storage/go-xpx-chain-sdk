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
	soInfo, err := client.Resolve.GetMosaicInfoByAssetId(ctx, sdk.StorageNamespaceId)
	require.NoError(t, err, err)

	smInfo, err := client.Resolve.GetMosaicInfoByAssetId(ctx, sdk.StreamingNamespaceId)
	require.NoError(t, err, err)

	scInfo, err := client.Resolve.GetMosaicInfoByAssetId(ctx, sdk.SuperContractNamespaceId)
	require.NoError(t, err, err)

	mosaics := []*sdk.MosaicId{soInfo.MosaicId, smInfo.MosaicId, scInfo.MosaicId}
	for _, mosaic := range mosaics {
		lps, err := client.LiquidityProvider.GetLiquidityProviders(ctx, nil)
		assert.Nil(t, err, err)

		currencyDeposit := sdk.XpxRelative(100000)
		result := sendTransaction(t, func() (sdk.Transaction, error) {
			return client.NewTransferTransaction(
				sdk.NewDeadline(time.Hour),
				managerAccount.Address,
				[]*sdk.Mosaic{currencyDeposit},
				sdk.NewPlainMessage("Test"),
			)
		}, defaultAccount)
		require.Nil(t, result.error, result.error)

		slashingAccount, err := client.NewAccountFromPublicKey("0000000000000000000000000000000000000000000000000000000000000000")
		require.Nil(t, err, err)

		result = sendTransaction(t, func() (sdk.Transaction, error) {
			return client.NewCreateLiquidityProviderTransaction(
				sdk.NewDeadline(time.Hour),
				mosaic,
				sdk.XpxRelative(100000).Amount,
				sdk.XpxRelative(100000).Amount,
				500,
				5,
				slashingAccount,
				500,
				500,
			)
		}, managerAccount)
		assert.Nil(t, result.error, result.error)

		lpsAfter, err := client.LiquidityProvider.GetLiquidityProviders(ctx, &sdk.LiquidityProviderPageOptions{Owner: managerAccount.PublicAccount.PublicKey})
		assert.Nil(t, err, err)
		assert.Equal(t, lps.Pagination.TotalEntries+uint64(len(mosaics)), lpsAfter.Pagination.TotalEntries)
	}
}

func TestManualRateChangeTransactionTransaction(t *testing.T) {
	soInfo, err := client.Resolve.GetMosaicInfoByAssetId(ctx, sdk.StorageNamespaceId)
	require.NoError(t, err, err)

	smInfo, err := client.Resolve.GetMosaicInfoByAssetId(ctx, sdk.StreamingNamespaceId)
	require.NoError(t, err, err)

	scInfo, err := client.Resolve.GetMosaicInfoByAssetId(ctx, sdk.SuperContractNamespaceId)
	require.NoError(t, err, err)

	mosaicAmount := uint64(1000000)
	currencyAmount := uint64(2000000)
	mosaics := []*sdk.MosaicId{soInfo.MosaicId, smInfo.MosaicId, scInfo.MosaicId}
	for _, mosaic := range mosaics {
		result := sendTransaction(t, func() (sdk.Transaction, error) {
			return client.NewTransferTransaction(
				sdk.NewDeadline(time.Hour),
				managerAccount.Address,
				[]*sdk.Mosaic{sdk.XpxRelative(currencyAmount)},
				sdk.NewPlainMessage("Test"),
			)
		}, defaultAccount)
		require.Nil(t, result.error)

		result = sendTransaction(t, func() (sdk.Transaction, error) {
			return client.NewManualRateChangeTransaction(
				sdk.NewDeadline(time.Hour),
				mosaic,
				true,
				sdk.Amount(mosaicAmount),
				true,
				sdk.Amount(currencyAmount),
			)
		}, managerAccount)
		assert.Nil(t, result.error)
	}
}
