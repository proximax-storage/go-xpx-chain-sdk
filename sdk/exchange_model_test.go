// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var exchangeAccount, _ = NewAccountFromPublicKey("ED7A848FDEB2321EE97CE8AF265588C54B4A58C72117247C7205EB061865055C", PublicTest)
var exchangeMosaicId, _    = NewMosaicId(0x26514E2A1EF33824)

func TestCounterOffer(t *testing.T) {
	assert.Equal(t, SellOffer.CounterOffer(), BuyOffer)
	assert.Equal(t, BuyOffer.CounterOffer(), SellOffer)
}

func TestCost_Cost_Overflow(t *testing.T) {
	offer := &OfferInfo{
		Owner: exchangeAccount,
		Type: SellOffer,
		Mosaic: newMosaicPanic(exchangeMosaicId, Amount(9000000000000000)),
		Deadline: Duration(10000023),
	}

	offer.PriceNumerator = 9000000000000000
	offer.PriceDenominator = 9000000000000000
	cost, err := offer.Cost(Amount(1000000000))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(1000000000))
}

func TestCost_SellOffer(t *testing.T) {
	offer := &OfferInfo{
		Owner: exchangeAccount,
		Type: SellOffer,
		Mosaic: newMosaicPanic(exchangeMosaicId, Amount(100)),
		Deadline: Duration(10000023),
	}

	offer.PriceNumerator = 1
	offer.PriceDenominator = 3
	cost, err := offer.Cost(Amount(1))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(1))

	cost, err = offer.Cost(Amount(2))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(1))

	cost, err = offer.Cost(Amount(3))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(1))

	cost, err = offer.Cost(Amount(4))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(2))

	cost, err = offer.Cost(Amount(12))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(4))
}

func TestCost_BuyOffer(t *testing.T) {
	offer := &OfferInfo{
		Owner: exchangeAccount,
		Type: BuyOffer,
		Mosaic: newMosaicPanic(exchangeMosaicId, Amount(100)),
		Deadline: Duration(10000023),
	}

	offer.PriceNumerator = 1
	offer.PriceDenominator = 3
	cost, err := offer.Cost(Amount(1))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(0))

	cost, err = offer.Cost(Amount(2))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(0))

	cost, err = offer.Cost(Amount(3))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(1))

	cost, err = offer.Cost(Amount(4))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(1))

	cost, err = offer.Cost(Amount(13))
	assert.Nil(t, err)
	assert.Equal(t, cost, Amount(4))
}

func TestCost_Not_Enough_Amount(t *testing.T) {
	offer := &OfferInfo{
		Owner: exchangeAccount,
		Type: BuyOffer,
		Mosaic: newMosaicPanic(exchangeMosaicId, Amount(100)),
		Deadline: Duration(10000023),
		PriceNumerator: Amount(123),
		PriceDenominator: Amount(321),
	}

	_, err := offer.Cost(Amount(101))
	assert.NotNil(t, err)
}

func TestCost_Unknown_Type(t *testing.T) {
	offer := &OfferInfo{
		Owner: exchangeAccount,
		Type: 3,
		Mosaic: newMosaicPanic(exchangeMosaicId, uint64DTO{ 100, 0 }.toStruct()),
		Deadline: uint64DTO{ 10000023, 0 }.toStruct(),
	}

	_, err := offer.Cost(Amount(1))
	assert.NotNil(t, err)
}

func TestConfirmOffer(t *testing.T) {
	offer := &OfferInfo{
		Owner: exchangeAccount,
		Type: SellOffer,
		Mosaic: newMosaicPanic(exchangeMosaicId, Amount(100)),
		Deadline: Duration(10000023),
		PriceNumerator: Amount(1),
		PriceDenominator: Amount(3),
	}

	confirmation, err := offer.ConfirmOffer(99)
	assert.Nil(t, err)
	assert.Equal(t, confirmation.Type, SellOffer)
	assert.Equal(t, confirmation.Owner, exchangeAccount)
	assert.Equal(t, confirmation.Mosaic, newMosaicPanic(exchangeMosaicId, Amount(99)))
	assert.Equal(t, confirmation.Cost, Amount(33))
}

