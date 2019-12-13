package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testAccountExchangeInfoJson = `{
  "exchange": {
    "owner": "ED7A848FDEB2321EE97CE8AF265588C54B4A58C72117247C7205EB061865055C",
    "ownerAddress": "909145399E6B95592041FCD38B6EE6ED2F20DDF5439BA8FD2D",
    "buyOffers": [],
    "sellOffers": [
      {
        "mosaicId": [
          519256100,
          642862634
        ],
        "amount": [
          997650,
          0
        ],
        "initialAmount": [
          1000000,
          0
        ],
        "initialCost": [
          500000,
          0
        ],
        "deadline": [
          10000023,
          0
        ],
        "price": 0.5
      }
    ],
    "expiredBuyOffers": [],
    "expiredSellOffers": []
  }
}`
	testOfferJson = `{
    "mosaicId": [
      519256100,
      642862634
    ],
    "amount": [
      997650,
      0
    ],
    "initialAmount": [
      1000000,
      0
    ],
    "initialCost": [
      500000,
      0
    ],
    "deadline": [
      10000023,
      0
    ],
    "price": 0.5,
    "owner": "ED7A848FDEB2321EE97CE8AF265588C54B4A58C72117247C7205EB061865055C",
    "type": 0
  }`

	testOfferJsonArr    = "[" + testOfferJson + ", " + testOfferJson + "]"
)

var testExchangeAccount, _ = NewAccountFromPublicKey("ED7A848FDEB2321EE97CE8AF265588C54B4A58C72117247C7205EB061865055C", PublicTest)

var (
	testExchangeMosaicId, _    = NewMosaicId(0x26514E2A1EF33824)
	offer = &OfferInfo{
		Owner: testExchangeAccount,
		Type: SellOffer,
		Mosaic: newMosaicPanic(testExchangeMosaicId, uint64DTO{ 997650, 0 }.toStruct()),
		PriceNumerator: uint64DTO{ 500000, 0 }.toStruct(),
		PriceDenominator: uint64DTO{ 1000000, 0 }.toStruct(),
		Deadline: uint64DTO{ 10000023, 0 }.toStruct(),
	}

	testUserExchangeInfo = &UserExchangeInfo{
		Owner:  testExchangeAccount,
		Offers: map[OfferType]map[MosaicId]*OfferInfo{
			SellOffer: map[MosaicId]*OfferInfo{
				*testExchangeMosaicId: offer,
			},
			BuyOffer: make(map[MosaicId]*OfferInfo),
		},
	}
)

func TestExchangeService_GetAccountExchangeInfo(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(exchangeRoute, testExchangeAccount.PublicKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testAccountExchangeInfoJson,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Exchange

	defer mock.Close()

	info, err := exchangeClient.GetAccountExchangeInfo(ctx, testExchangeAccount)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, testUserExchangeInfo, info)
}

func TestExchangeService_GetExchangeOfferByAssetId(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(offersByMosaicRoute, SellOffer.String(), testExchangeMosaicId.toHexString()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testOfferJsonArr,
	})
	exchangeClient := mock.getPublicTestClientUnsafe().Exchange

	defer mock.Close()

	offers, err := exchangeClient.GetExchangeOfferByAssetId(ctx, testExchangeMosaicId, SellOffer)
	assert.Nil(t, err)
	assert.NotNil(t, offers)
	assert.Equal(t, len(offers), 2)
	assert.Equal(t, []*OfferInfo{ offer, offer }, offers)
}