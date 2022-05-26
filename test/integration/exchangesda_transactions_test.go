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
						sdk.Storage(10000),
						sdk.Streaming(100000000),
					},
					defaultAccount.PublicAccount,
					sdk.Duration(10000000),
				},
				{
					sdk.SdaOffer{
						sdk.Streaming(100),
						sdk.Storage(10000),
					},
					defaultAccount.PublicAccount,
					sdk.Duration(10000000),
				},
			},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
}

func TestPlaceSdaExchangeOfferTransaction(t *testing.T) {
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewPlaceSdaExchangeOfferTransaction(
			sdk.NewDeadline(time.Hour),
			[]*sdk.PlaceSdaOffer{
				{
					sdk.SdaOffer{
						sdk.Storage(10000),
						sdk.Streaming(100000000),
					},
					defaultAccount.PublicAccount,
					sdk.Duration(1000),
				},
				{
					sdk.SdaOffer{
						sdk.Streaming(100),
						sdk.Storage(10000),
					},
					defaultAccount.PublicAccount,
					sdk.Duration(1000),
				},
			},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
}

func TestRemoveSdaExchangeOfferTransaction(t *testing.T) {
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
	}, defaultAccount)
	assert.Nil(t, result.error)
}
