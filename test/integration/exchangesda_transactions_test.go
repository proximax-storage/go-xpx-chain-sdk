// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
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

func TestPlaceSdaExchangeOfferTransaction_LongOfferKey(t *testing.T) {
	configDelta := 5
	config, err := client.Network.GetNetworkConfig(ctx)
	assert.Nil(t, err)

	config.NetworkConfig.Sections["plugin:catapult.plugins.exchangesda"].Fields["longOfferKey"].Value = defaultAccount.PublicAccount.PublicKey

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewNetworkConfigTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.Duration(configDelta),
			config.NetworkConfig,
			config.SupportedEntityVersions)
	}, nemesisAccount)
	assert.Nil(t, result.error)

	waitForBlocksCount(t, configDelta)

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewPlaceSdaExchangeOfferTransaction(
			sdk.NewDeadline(time.Hour),
			[]*sdk.PlaceSdaOffer{
				{
					sdk.SdaOffer{
						sdk.Storage(1000),
						sdk.Streaming(2000),
					},
					sdk.Duration(10000000),
				},
				{
					sdk.SdaOffer{
						sdk.Streaming(2000),
						sdk.Storage(1000),
					},
					sdk.Duration(10000000),
				},
			},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
}

func TestPlaceSdaExchangeOfferTransaction(t *testing.T) {
	owner, err := client.NewAccountFromPrivateKey("E8230C5CD4EB49F6AC7FD191373E9B6322C5258CE6EE03119944A312CE6226F6")
	require.NoError(t, err, err)

	// add storage and streaming mosaic to the drive owner
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			owner.Address,
			[]*sdk.Mosaic{sdk.Storage(1000), sdk.Streaming(2000)},
			sdk.NewPlainMessage(""),
		)
	}, defaultAccount)
	require.NoError(t, result.error, result.error)
	// end region

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewPlaceSdaExchangeOfferTransaction(
			sdk.NewDeadline(time.Hour),
			[]*sdk.PlaceSdaOffer{
				{
					sdk.SdaOffer{
						sdk.Storage(100),
						sdk.Streaming(200),
					},
					sdk.Duration(1000),
				},
				{
					sdk.SdaOffer{
						sdk.Streaming(200),
						sdk.Storage(100),
					},
					sdk.Duration(1000),
				},
			},
		)
	}, owner)
	assert.Nil(t, result.error)
}

func TestRemoveSdaExchangeOfferTransaction(t *testing.T) {
	owner, err := client.NewAccountFromPrivateKey("E8230C5CD4EB49F6AC7FD191373E9B6322C5258CE6EE03119944A312CE6226F6")
	require.NoError(t, err, err)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewRemoveSdaExchangeOfferTransaction(
			sdk.NewDeadline(time.Hour),
			[]*sdk.RemoveSdaOffer{
				{
					sdk.StorageNamespaceId,
					sdk.StreamingNamespaceId,
				},
				{
					sdk.StreamingNamespaceId,
					sdk.StorageNamespaceId,
				},
			},
		)
	}, owner)
	assert.Nil(t, result.error)
}
