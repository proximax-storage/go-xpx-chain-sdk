// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testAccountSdaExchangeInfoJson = `{
        "exchangesda": {
            "owner": "ED7A848FDEB2321EE97CE8AF265588C54B4A58C72117247C7205EB061865055C",
            "ownerAddress": "909145399E6B95592041FCD38B6EE6ED2F20DDF5439BA8FD2D",
            "sdaOfferBalances": [
                {
                    "mosaicIdGive": [
                        519256100,
                        642862634
                    ],
                    "mosaicIdGet": [
                        145530229,
                        1818060917
                    ],
                    "currentMosaicGiveAmount": [
                        997650,
                        0
                    ],
                    "currentMosaicGetAmount": [
                        997650,
                        0
                    ],
                    "initialMosaicGiveAmount": [
                        50000,
                        0
                    ],
                    "initialMosaicGetAmount": [
                        10000,
                        0
                    ],
                    "deadline": [
                        10000023,
                        0
                    ]
                }
            ],
            "expiredSdaOfferBalances": []
        }
    }`

	testSdaOfferBalanceJson = `{
        "owner": "ED7A848FDEB2321EE97CE8AF265588C54B4A58C72117247C7205EB061865055C",
        "mosaicIdGive": [
            519256100,
            642862634
        ],
        "mosaicIdGet": [
            145530229,
            1818060917
        ],
        "currentMosaicGiveAmount": [
            997650,
            0
        ],
        "currentMosaicGetAmount": [
            997650,
            0
        ],
        "initialMosaicGiveAmount": [
            50000,
            0
        ],
        "initialMosaicGetAmount": [
            10000,
            0
        ],
        "deadline": [
            10000023,
            0
        ]
    }`

	testSdaOfferBalanceJsonArr = "[" + testSdaOfferBalanceJson + ", " + testSdaOfferBalanceJson + "]"
)

var testSdaExchangeAccount, _ = NewAccountFromPublicKey("ED7A848FDEB2321EE97CE8AF265588C54B4A58C72117247C7205EB061865055C", PublicTest)

var (
	testSdaExchangeMosaicIdGive, _ = NewMosaicId(0x26514E2A1EF33824)
	testSdaExchangeMosaicIdGet, _  = NewMosaicId(0x6C5D687508AC9D75)

	testSdaOfferBalance = &SdaOfferBalance{
		Owner:             testSdaExchangeAccount,
		MosaicGive:        newMosaicPanic(testSdaExchangeMosaicIdGive, uint64DTO{997650, 0}.toStruct()),
		MosaicGet:         newMosaicPanic(testSdaExchangeMosaicIdGet, uint64DTO{997650, 0}.toStruct()),
		InitialAmountGive: uint64DTO{50000, 0}.toStruct(),
		InitialAmountGet:  uint64DTO{10000, 0}.toStruct(),
		Deadline:          uint64DTO{10000023, 0}.toStruct(),
	}

	testUserSdaExchangeInfo = &UserSdaExchangeInfo{
		Owner:            testSdaExchangeAccount,
		SdaOfferBalances: []*SdaOfferBalance{testSdaOfferBalance},
	}
)

func TestSdaExchangeService_GetAccountSdaExchangeInfo(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(exchangeSdaRoute, testSdaExchangeAccount.PublicKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testAccountSdaExchangeInfoJson,
	})

	sdaExchangeClient := mock.getPublicTestClientUnsafe().SdaExchange

	defer mock.Close()

	info, err := sdaExchangeClient.GetAccountSdaExchangeInfo(ctx, testSdaExchangeAccount)
	assert.Nil(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, testUserSdaExchangeInfo, info)
}

func TestSdaExchangeService_GetSdaExchangeOfferByAssetId(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(sdaOffersByMosaicRoute, "give", testSdaExchangeMosaicIdGive.toHexString()),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testSdaOfferBalanceJsonArr,
	})

	sdaExchangeClient := mock.getPublicTestClientUnsafe().SdaExchange

	defer mock.Close()

	offers, err := sdaExchangeClient.GetSdaExchangeOfferByAssetId(ctx, testSdaExchangeMosaicIdGive, "give")
	assert.Nil(t, err)
	assert.NotNil(t, offers)
	assert.Equal(t, len(offers), 2)
	assert.Equal(t, []*SdaOfferBalance{testSdaOfferBalance, testSdaOfferBalance}, offers)
}
