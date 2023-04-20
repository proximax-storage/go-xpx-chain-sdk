// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"crypto/rand"
	"encoding/hex"
	math "math/rand"
	"sort"
	"testing"
	"time"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Asset struct {
	MosaicId    *sdk.MosaicId
	NamespaceId *sdk.NamespaceId
}

func createLinkedMosaicIdAndNamespace(t *testing.T, owner *sdk.Account) Asset {
	t.Helper()

	name := make([]byte, 5)

	_, err := rand.Read(name)
	assert.NoError(t, err, err)
	nameHex := hex.EncodeToString(name)

	namespaceId, err := sdk.NewNamespaceIdFromName(nameHex)
	assert.NoError(t, err, err)

	registerTx, err := client.NewRegisterRootNamespaceTransaction(
		sdk.NewDeadline(time.Hour),
		nameHex,
		sdk.Duration(0),
	)
	assert.NoError(t, err, err)
	registerTx.ToAggregate(owner.PublicAccount)

	r := math.New(math.NewSource(time.Now().UTC().UnixNano()))
	nonce := r.Uint32()
	supply := sdk.Amount(10000000000000)

	mosaicId, err := sdk.NewMosaicIdFromNonceAndOwner(nonce, owner.PublicAccount.PublicKey)
	assert.NoError(t, err, err)
	mosaicDefinitionTx, err := client.NewMosaicDefinitionTransaction(
		sdk.NewDeadline(time.Hour),
		nonce,
		owner.PublicAccount.PublicKey,
		sdk.NewMosaicProperties(true, true, 4, sdk.Duration(0)),
	)
	assert.NoError(t, err, err)
	mosaicDefinitionTx.ToAggregate(owner.PublicAccount)
	mosaicSupplyTx, err := client.NewMosaicSupplyChangeTransaction(
		sdk.NewDeadline(time.Hour),
		mosaicId,
		sdk.Increase,
		supply,
	)
	assert.NoError(t, err, err)
	mosaicSupplyTx.ToAggregate(owner.PublicAccount)

	aliasTx, err := client.NewMosaicAliasTransaction(
		sdk.NewDeadline(time.Hour),
		mosaicId,
		namespaceId,
		sdk.AliasLink,
	)
	assert.Nil(t, err)
	aliasTx.ToAggregate(owner.PublicAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewCompleteAggregateTransaction(
			sdk.NewDeadline(time.Hour),
			[]sdk.Transaction{registerTx, mosaicDefinitionTx, mosaicSupplyTx, aliasTx},
		)
	}, owner)
	assert.Nil(t, result.error)

	var assets Asset
	assets.MosaicId = mosaicId
	assets.NamespaceId = namespaceId

	return assets
}

func placeSdaExchangeOfferTransaction(t *testing.T, owner *sdk.Account, assets []Asset, duration sdk.Duration) {
	t.Helper()

	// add SdaOne and SdaTwo to the account owner
	a := []*sdk.MosaicId{assets[0].MosaicId, assets[1].MosaicId}
	sort.Slice(a, func(i, j int) bool {
		return a[i].Id() < a[j].Id()
	})
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewTransferTransaction(
			sdk.NewDeadline(time.Hour),
			owner.Address,
			[]*sdk.Mosaic{{AssetId: a[0], Amount: 2000}, {AssetId: a[1], Amount: 2000}},
			sdk.NewPlainMessage(""),
		)
	}, nemesisAccount)
	require.NoError(t, result.error, result.error)
	// end region

	result = sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewPlaceSdaExchangeOfferTransaction(
			sdk.NewDeadline(time.Hour),
			[]*sdk.PlaceSdaOffer{
				{
					SdaOffer: sdk.SdaOffer{
						MosaicGive: &sdk.Mosaic{AssetId: assets[0].MosaicId, Amount: 1000},
						MosaicGet:  &sdk.Mosaic{AssetId: assets[1].MosaicId, Amount: 2000},
					},
					Duration: duration,
				},
				{
					SdaOffer: sdk.SdaOffer{
						MosaicGive: &sdk.Mosaic{AssetId: assets[1].MosaicId, Amount: 200},
						MosaicGet:  &sdk.Mosaic{AssetId: assets[0].MosaicId, Amount: 100},
					},
					Duration: duration,
				},
			},
		)
	}, owner)
	assert.Nil(t, result.error)
}

func removeSdaExchangeOfferTransaction(t *testing.T, owner *sdk.Account, assets []Asset) {
	t.Helper()

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewRemoveSdaExchangeOfferTransaction(
			sdk.NewDeadline(time.Hour),
			[]*sdk.RemoveSdaOffer{
				{
					AssetIdGive: assets[0].NamespaceId,
					AssetIdGet:  assets[1].NamespaceId,
				},
				{
					AssetIdGive: assets[1].NamespaceId,
					AssetIdGet:  assets[0].NamespaceId,
				},
			},
		)
	}, owner)
	assert.Nil(t, result.error)
}

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

	// Create 2 eternal Sirius Digital Assets
	var assets []Asset
	SdaOne := createLinkedMosaicIdAndNamespace(t, nemesisAccount)
	SdaTwo := createLinkedMosaicIdAndNamespace(t, nemesisAccount)
	assets = append(assets, SdaOne, SdaTwo)

	placeSdaExchangeOfferTransaction(t, defaultAccount, assets, sdk.Duration(10000000))
	removeSdaExchangeOfferTransaction(t, defaultAccount, assets)
}

func TestPlaceAndRemoveSdaExchangeOfferTransaction(t *testing.T) {
	owner, err := client.NewAccountFromPrivateKey("E8230C5CD4EB49F6AC7FD191373E9B6322C5258CE6EE03119944A312CE6226F6")
	require.NoError(t, err, err)

	// Create 2 eternal Sirius Digital Assets
	var assets []Asset
	SdaOne := createLinkedMosaicIdAndNamespace(t, nemesisAccount)
	SdaTwo := createLinkedMosaicIdAndNamespace(t, nemesisAccount)
	assets = append(assets, SdaOne, SdaTwo)
	placeSdaExchangeOfferTransaction(t, owner, assets, sdk.Duration(1000))
	removeSdaExchangeOfferTransaction(t, owner, assets)
}
