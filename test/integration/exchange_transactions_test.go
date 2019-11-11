// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddExchangeOfferTransaction(t *testing.T) {
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewAddExchangeOfferTransaction(
			sdk.NewDeadline(time.Hour),
			[]*sdk.AddOffer{
				{
					sdk.Offer{
						sdk.SellOffer,
						sdk.Storage(10000000),
						sdk.Amount(10000000),
					},
					sdk.Duration(100),
				},
			},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
}

func TestRemoveExchangeOfferTransaction(t *testing.T) {
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewRemoveExchangeOfferTransaction(
			sdk.NewDeadline(time.Hour),
			[]*sdk.RemoveOffer{
				{
					sdk.SellOffer,
					sdk.StorageNamespaceId,
				},
			},
		)
	}, defaultAccount)
	assert.Nil(t, result.error)
}