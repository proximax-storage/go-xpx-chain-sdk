// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
)

type OfferType uint8

const (
	SellOffer OfferType = iota
	BuyOffer
	UnknownType
)

func(o OfferType) CounterOffer() OfferType {
	switch o {
	case SellOffer:
		return BuyOffer
	case BuyOffer:
		return SellOffer
	default:
		return UnknownType
	}
}

type UserExchangeInfo struct {
	Owner	*PublicAccount
	Offers	map[OfferType]map[MosaicId]OfferInfo
}

type OfferInfo struct {
	Type				OfferType
	Owner				*PublicAccount
	Mosaic				*Mosaic
	PriceNumerator		Amount
	PriceDenominator	Amount
	Dealine				Height
}

func(o *OfferInfo) Cost(amount Amount) (Amount, error) {
	if o.Mosaic.Amount < amount {
		return 0, errors.New("You can't get more mosaics when in offer")
	}

	switch o.Type {
	case SellOffer:
		// If user want to buy mosaic, we round the cost towards the seller(because we buy part of mosaics)
		return Amount(math.Ceil(float64(o.PriceNumerator * amount) / float64(o.PriceDenominator))), nil
	case BuyOffer:
		// If user want to sell mosaic, we round the cost towards the buyer(because we sell part of mosaics)
		return Amount(math.Floor(float64(o.PriceNumerator * amount) / float64(o.PriceDenominator))), nil
	default:
		return 0, errors.New("Unknown offer type")
	}
}

func(o *OfferInfo) ConfirmOffer(amount Amount) (*ExchangeConfirmation, error) {
	cost, err := o.Cost(amount)
	if err != nil {
		return nil, err
	}

	confirmation := &ExchangeConfirmation{
		Offer{
			Type:		o.Type,
			Mosaic:		newMosaicPanic(o.Mosaic.AssetId, amount ),
			Cost:		cost,
		},
		o.Owner,
	}

	return confirmation, nil
}

type Offer struct {
	Type 		OfferType
	Mosaic		*Mosaic
	Cost 		Amount
}

type AddOffer struct {
	Offer
	Duration	Duration
}

func (offer *AddOffer) String() string {
	return fmt.Sprintf(
		`
			"Type": %d,
			"AssetId": %s,
			"Amount": %s,
			"Cost": %s,
			"Duration": %s,
		`,
		offer.Type,
		offer.Mosaic.AssetId,
		offer.Mosaic.Amount,
		offer.Cost,
		offer.Duration,
	)
}

// Add Exchange Offer Transaction
type AddExchangeOfferTransaction struct {
	AbstractTransaction
	Offers				[]*AddOffer
}

type ExchangeConfirmation struct {
	Offer
	Owner				*PublicAccount
}

func (offer *ExchangeConfirmation) String() string {
	return fmt.Sprintf(
		`
			"Type": %d,
			"AssetId": %s,
			"Amount": %s,
			"Cost": %s,
			"Owner": %s,
		`,
		offer.Type,
		offer.Mosaic.AssetId,
		offer.Mosaic.Amount,
		offer.Cost,
		offer.Owner,
	)
}

// Exchange Transaction
type ExchangeOfferTransaction struct {
	AbstractTransaction
	Confirmations		[]*ExchangeConfirmation
}

type RemoveOffer struct {
	Type		OfferType
	AssetId		AssetId
}

func (offer *RemoveOffer) String() string {
	return fmt.Sprintf(
		`
			"Type": %d,
			"AssetId": %s,
		`,
		offer.Type,
		offer.AssetId,
	)
}

// Remove Exchange Offer Transaction
type RemoveExchangeOfferTransaction struct {
	AbstractTransaction
	Offers			[]*RemoveOffer
}