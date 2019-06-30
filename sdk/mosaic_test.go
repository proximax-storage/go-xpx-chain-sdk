// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/proximax-storage/go-xpx-utils/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	mosaicClient = mockServer.getPublicTestClientUnsafe().Mosaic
)

const (
	testMosaicPathID = "6C55E05D11D19FBD"

	testMosaicInfoJson = `
{
  "mosaic": {
    "mosaicId": [
		298950589,
		1817567325
    ],
    "supply": [
      3403414400,
      2095475
    ],
    "height": [
      1,
      0
    ],
    "owner": "321DE652C4D3362FC2DDF7800F6582F4A10CFEA134B81F8AB6E4BE78BBA4D18E",
	"revision": 1,
    "properties": [
	  {
    	"value": [
        	2,
        	0
      	],
    	"id": 0
	  },
	  {
    	"value": [
        	6,
        	0
      	],
    	"id": 1
	  },
	  {
    	"value": [
        	1,
        	0
      	],
    	"id": 2
	  }
    ]
  }
}`

	testMosaicNamesJson = `[
   {
      "mosaicId":[
         519256100,
         642862634
      ],
      "names":[
         "cat.storage"
      ]
   },
   {
      "mosaicId":[
         481110499,
         231112638
      ],
      "names":[
         "cat.currency"
      ]
   }
]`
)

var (
	mosaicCorr = &MosaicInfo{
		MosaicId: NewMosaicIdNoCheck(uint64DTO{298950589, 1817567325}.toUint64()),
		Supply:   uint64DTO{3403414400, 2095475}.toStruct(),
		Height:   uint64DTO{1, 0}.toStruct(),
		Owner: &PublicAccount{
			Address: &Address{
				Type:    mosaicClient.client.config.NetworkType,
				Address: "VBFBW6TUGLEWQIBCMTBMXXQORZKUP3WTVX36ZFE7",
			},

			PublicKey: "321DE652C4D3362FC2DDF7800F6582F4A10CFEA134B81F8AB6E4BE78BBA4D18E",
		},
		Revision: 1,
		Properties: NewMosaicProperties(
			false,
			true,
			false,
			6,
			uint64DTO{1, 0}.toStruct(),
		),
	}

	mosaicNames = []*MosaicName{
		{
			NewMosaicIdNoCheck(0x26514E2A1EF33824),
			[]string{"cat.storage"},
		},
		{
			NewMosaicIdNoCheck(0x0DC67FBE1CAD29E3),
			[]string{"cat.currency"},
		},
	}
)

func TestMosaicService_GetMosaic(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf(mosaicRoute, testMosaicPathID),
		RespBody: testMosaicInfoJson,
	})

	mscInfo, err := mosaicClient.GetMosaicInfo(ctx, mosaicCorr.MosaicId)

	assert.Nilf(t, err, "MosaicService.GetMosaic returned error: %s", err)
	tests.ValidateStringers(t, mosaicCorr, mscInfo)
}

func TestMosaicService_GetMosaics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockServer.AddRouter(&mock.Router{
			Path:     mosaicsRoute,
			RespBody: "[" + testMosaicInfoJson + "]",
			ReqJsonBodyStruct: struct {
				MosaicIds []string `json:"mosaicIds"`
			}{},
		})

		mscInfoArr, err := mosaicClient.GetMosaicInfos(ctx, []*MosaicId{mosaicCorr.MosaicId})

		assert.Nilf(t, err, "MosaicService.GetMosaics returned error: %s", err)

		for _, mscInfo := range mscInfoArr {
			tests.ValidateStringers(t, mosaicCorr, mscInfo)
		}
	})

	t.Run("empty url params", func(t *testing.T) {
		_, err := mosaicClient.GetMosaicInfos(ctx, []*MosaicId{})

		assert.NotNil(t, err, "MosaicService.GetMosaics returned error: %s", err)
	})
}

func TestMosaicService_GetMosaicsNames(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     mosaicNamesRoute,
		RespBody: testMosaicNamesJson,
		ReqJsonBodyStruct: struct {
			MosaicIds []string `json:"mosaicIds"`
		}{},
	})

	mscNameArr, err := mosaicClient.GetMosaicsNames(ctx, mosaicNames[0].MosaicId, mosaicNames[1].MosaicId)

	assert.Nilf(t, err, "MosaicService.GetMosaics returned error: %s", err)

	for i, mscName := range mscNameArr {
		tests.ValidateStringers(t, mosaicNames[i], mscName)
	}
}
