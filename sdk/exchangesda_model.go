// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import "fmt"

type SdaOffer struct {
	MosaicGive *Mosaic
	MosaicGet  *Mosaic
}

type PlaceSdaOffer struct {
	SdaOffer
	Owner    *PublicAccount
	Duration Duration
}

func (offer *PlaceSdaOffer) String() string {
	return fmt.Sprintf(
		`
			"AssetIdGive": %s,
			"AmountGive": %s,
			"AssetIdGet": %s,
			"AmountGet": %s,
			"Owner": %s,
			"Duration": %s,
		`,
		offer.MosaicGive.AssetId,
		offer.MosaicGive.Amount,
		offer.MosaicGet.AssetId,
		offer.MosaicGet.Amount,
		offer.Owner,
		offer.Duration,
	)
}

// Place SDA-SDA Exchange Offer Transaction
type PlaceSdaExchangeOfferTransaction struct {
	AbstractTransaction
	Offers []*PlaceSdaOffer
}

type RemoveSdaOffer struct {
	AssetIdGive AssetId
	AssetIdGet  AssetId
}

func (offer *RemoveSdaOffer) String() string {
	return fmt.Sprintf(
		`
			"AssetIdGive": %s,
			"AssetIdGet": %s,
		`,
		offer.AssetIdGive,
		offer.AssetIdGet,
	)
}

// Remove SDA-SDA Exchange Offer Transaction
type RemoveSdaExchangeOfferTransaction struct {
	AbstractTransaction
	Offers []*RemoveSdaOffer
}
