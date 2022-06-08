// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import "fmt"

type UserSdaExchangeInfo struct {
	Owner            *PublicAccount
	SdaOfferBalances map[MosaicId]map[MosaicId]*SdaOfferBalance
}

func (info *UserSdaExchangeInfo) String() string {
	return fmt.Sprintf(
		`
			"Owner": %s,
			"SdaOfferBalances": %+v,
		`,
		info.Owner,
		info.SdaOfferBalances,
	)
}

type SdaOfferBalance struct {
	Owner             *PublicAccount
	MosaicGive        *Mosaic
	MosaicGet         *Mosaic
	InitialAmountGive Amount
	InitialAmountGet  Amount
	Deadline          Height
}

func (info *SdaOfferBalance) String() string {
	return fmt.Sprintf(
		`
			"Owner": %s,
			"MosaicGive": %s,
			"MosaicGet": %s,
			"InitialAmountGive": %s,
			"InitialAmountGet": %s,
			"Deadline": %d,
		`,
		info.Owner,
		info.MosaicGive,
		info.MosaicGet,
		info.InitialAmountGive,
		info.InitialAmountGet,
		info.Deadline,
	)
}

type SdaOffer struct {
	MosaicGive *Mosaic
	MosaicGet  *Mosaic
}

type PlaceSdaOffer struct {
	SdaOffer
	Duration Duration
}

func (offer *PlaceSdaOffer) String() string {
	return fmt.Sprintf(
		`
			"AssetIdGive": %s,
			"AmountGive": %s,
			"AssetIdGet": %s,
			"AmountGet": %s,
			"Duration": %s,
		`,
		offer.MosaicGive.AssetId,
		offer.MosaicGive.Amount,
		offer.MosaicGet.AssetId,
		offer.MosaicGet.Amount,
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
