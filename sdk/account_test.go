// Copyright 2018 ProximaX Limited. All rights reserved. // Use of this source code is governed by the Apache 2.0 // license that can be found in the LICENSE file.
package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/proximax-storage/go-xpx-utils/tests"
	"github.com/stretchr/testify/assert"
)

var (
	account = &AccountInfo{
		Address:         &Address{MijinTest, "SAONSOGFZZHNEIBRYXHDTDTBR2YSAXKTITRFHG2Y"},
		AddressHeight:   uint64DTO{1, 0}.toStruct(),
		PublicKey:       "F3824119C9F8B9E81007CAA0EDD44F098458F14503D7C8D7C24F60AF11266E57",
		PublicKeyHeight: uint64DTO{0, 0}.toStruct(),
		AccountType:     MainAccount,
		LinkedAccount: &PublicAccount{
			&Address{Type: MijinTest, Address: "SDYVPENRSMSGU24XSSCQPHKKWYUNKYFDLAVTUMMS"},
			"F2D7845487664F4417232C93771C337FA34B78BE053EF22C4EAFB2005BD65006",
		},
		Mosaics: []*Mosaic{
			newMosaicPanic(newMosaicIdPanic(uint64DTO{298950589, 1817567325}.toUint64()), uint64DTO{3863990592, 95248}.toStruct()),
		},
		Reputation: 0.9,
	}

	accountProperties = &AccountProperties{
		Address:            &Address{MijinTest, "SAONSOGFZZHNEIBRYXHDTDTBR2YSAXKTITRFHG2Y"},
		AllowedAddresses:   []*Address{{MijinTest, "SAONSOGFZZHNEIBRYXHDTDTBR2YSAXKTITRFHG2Y"}},
		AllowedMosaicId:    []*MosaicId{newMosaicIdPanic(uint64DTO{1486560344, 659392627}.toUint64())},
		AllowedEntityTypes: []EntityType{LinkAccount},
		BlockedAddresses:   []*Address{{MijinTest, "SAONSOGFZZHNEIBRYXHDTDTBR2YSAXKTITRFHG2Y"}},
		BlockedMosaicId:    []*MosaicId{newMosaicIdPanic(uint64DTO{1486560344, 659392627}.toUint64())},
		BlockedEntityTypes: []EntityType{LinkAccount},
	}

	accountClient = mockServer.getPublicTestClientUnsafe().Account

	harvester1 = &Harvester{
		Key:                    "119EAB9545B31613D88557F8E783DBD1D01790783C112742B597F14E28A8A50E",
		Owner:                  "1DBCFA374315B059FDA6B08A981737CECB73912D4689069CD71850DCC3AA3031",
		Address:                &Address{Public, "XDWAZJDTYD65Y456F6BR2WG2MMEVIT2ZYC6W4UT4"},
		DisabledHeight:         uint64DTO{0, 0}.toStruct(),
		LastSigningBlockHeight: uint64DTO{6001108, 0}.toStruct(),
		EffectiveBalance:       uint64DTO{1354348165, 205}.toStruct(),
		CanHarvest:             true,
		Activity:               0.000006784999772748008,
		Greed:                  0.1,
	}

	harvester2 = &Harvester{
		Key:                    "1837E8A42E75E974C52BF470DF39A2A07FB164867CC1B126FD79A60D37EE7544",
		Owner:                  "2F88213BFE21E22B3A082E5D89091066F892D5A92E1A0949B32209D77506E289",
		Address:                &Address{Public, "XCJ2A2GT7JBTHLVNAX4PRQUDBUSADFRI6QPJK2ME"},
		DisabledHeight:         uint64DTO{0, 0}.toStruct(),
		LastSigningBlockHeight: uint64DTO{6124616, 0}.toStruct(),
		EffectiveBalance:       uint64DTO{1686419070, 2979}.toStruct(),
		CanHarvest:             true,
		Activity:               -0.000003215000227251993,
		Greed:                  0.1,
	}

	harvesters = &HarvestersPage{
		Harvesters: []*Harvester{
			harvester1,
			harvester2,
		},
		Pagination: &Pagination{
			TotalEntries: 75,
			PageNumber:   1,
			PageSize:     2,
			TotalPages:   38,
		},
	}
)

const (
	accountInfoJson = `{  
   "meta":{  

   },
   "account":{  
      "address":"901CD938C5CE4ED22031C5CE398E618EB1205D5344E2539B58",
      "addressHeight":[  
         1,
         0
      ],
      "publicKey":"F3824119C9F8B9E81007CAA0EDD44F098458F14503D7C8D7C24F60AF11266E57",
      "publicKeyHeight":[  
         0,
         0
      ],
	  "accountType": 1,
      "linkedAccountKey": "F2D7845487664F4417232C93771C337FA34B78BE053EF22C4EAFB2005BD65006",
      "mosaics":[  
         {  
            "id":[  
               298950589,
               1817567325
            ],
            "amount":[  
               3863990592,
               95248
            ]
         }
      ]
   }
}
`
	accountPropertiesJson = `{
  "accountProperties": {
    "address": "901CD938C5CE4ED22031C5CE398E618EB1205D5344E2539B58",
    "properties": [
      {
        "propertyType": 1,
        "values": [
          "901CD938C5CE4ED22031C5CE398E618EB1205D5344E2539B58"
        ]
      },
      {
        "propertyType": 2,
        "values": [
          [
            1486560344,
            659392627
          ]
        ]
      },
      {
        "propertyType": 4,
        "values": [
          16716
        ]
      },
      {
        "propertyType": 129,
        "values": [
          "901CD938C5CE4ED22031C5CE398E618EB1205D5344E2539B58"
        ]
      },
      {
        "propertyType": 130,
        "values": [
          [
            1486560344,
            659392627
          ]
        ]
      },
      {
        "propertyType": 132,
        "values": [
          16716
        ]
      }
    ]
  }
}
`
	accountNameJson = `{  
    "address": "901CD938C5CE4ED22031C5CE398E618EB1205D5344E2539B58",
    "names": [
      "alias1",
      "alias2"
    ]
},
{  
    "address": "9053D1FE65426CFC77C9092FBD329647634F2AAACE113868E0",
    "names": [
      "alias3",
      "alias4"
    ]
}
`
	accountHarvestingJson = `[{
    "harvester": {
        "key": "119EAB9545B31613D88557F8E783DBD1D01790783C112742B597F14E28A8A50E",
        "owner": "1DBCFA374315B059FDA6B08A981737CECB73912D4689069CD71850DCC3AA3031",
        "address": "B8EC0CA473C0FDDC73BE2F831D58DA6309544F59C0BD6E527C",
        "disabledHeight": [0, 0],
        "lastSigningBlockHeight": [
            6001108, 0
        ],
        "effectiveBalance": [
            1354348165, 205
        ],
        "canHarvest": true,
        "activity": 0.000006784999772748008,
        "greed": 0.1
    }
}]
`
	harvestersPageJson = `{
    "data": [{
        "harvester": {
            "key": "119EAB9545B31613D88557F8E783DBD1D01790783C112742B597F14E28A8A50E",
            "owner": "1DBCFA374315B059FDA6B08A981737CECB73912D4689069CD71850DCC3AA3031",
            "address": "B8EC0CA473C0FDDC73BE2F831D58DA6309544F59C0BD6E527C",
            "disabledHeight": [
                0, 0
            ],
            "lastSigningBlockHeight": [
                6001108, 0
            ],
            "effectiveBalance": [
                1354348165,
                205
            ],
            "canHarvest": true,
            "activity": 0.000006784999772748008,
            "greed": 0.1
        },
        "meta": {
            "id": "631362904443c6c8d7366bba"
        }
    }, {
        "harvester": {
            "key": "1837E8A42E75E974C52BF470DF39A2A07FB164867CC1B126FD79A60D37EE7544",
            "owner": "2F88213BFE21E22B3A082E5D89091066F892D5A92E1A0949B32209D77506E289",
            "address": "B893A068D3FA4333AEAD05F8F8C2830D24019628F41E956984",
            "disabledHeight": [
                0, 0
            ],
            "lastSigningBlockHeight": [
                6124616, 0
            ],
            "effectiveBalance": [
                1686419070,
                2979
            ],
            "canHarvest": true,
            "activity": -0.000003215000227251993,
            "greed": 0.1
        },
        "meta": {
            "id": "632d81d69bc8787925d992a1"
        }
    }],
    "pagination": {
        "totalEntries": 75,
        "pageNumber": 1,
        "pageSize": 2,
        "totalPages": 38
    }
}
`
)

var (
	nemTestAddress1  = "SAONSOGFZZHNEIBRYXHDTDTBR2YSAXKTITRFHG2Y"
	nemTestAddress2  = "SBJ5D7TFIJWPY56JBEX32MUWI5RU6KVKZYITQ2HA"
	harvesterAddress = "XDWAZJDTYD65Y456F6BR2WG2MMEVIT2ZYC6W4UT4"
)

var (
	accountNames = []*AccountName{
		{
			Address: newAddressFromRaw(nemTestAddress1),
			Names:   []string{"alias1", "alias2"},
		},
		{
			Address: newAddressFromRaw(nemTestAddress2),
			Names:   []string{"alias3", "alias4"},
		},
	}
)

func TestAccountService_GetAccountProperties(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf(accountPropertiesRoute, nemTestAddress1),
		RespBody: accountPropertiesJson,
	})

	accP, err := accountClient.GetAccountProperties(context.Background(), &Address{MijinTest, nemTestAddress1})

	assert.Nilf(t, err, "AccountService.GetAccountProperties returned error: %s", err)

	tests.ValidateStringers(t, accountProperties, accP)
}

func TestAccountService_GetAccountsProperties(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     accountsPropertiesRoute,
		RespBody: "[" + accountPropertiesJson + "]",
	})

	accountsProperties, err := accountClient.GetAccountsProperties(
		context.Background(),
		&Address{MijinTest, nemTestAddress1},
	)

	assert.Nilf(t, err, "AccountService.GetAccountsProperties returned error: %s", err)

	for _, accP := range accountsProperties {
		tests.ValidateStringers(t, accountProperties, accP)
	}
}

func TestAccountService_GetAccountInfo(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf("/account/%s", nemTestAddress1),
		RespBody: accountInfoJson,
	})

	acc, err := accountClient.GetAccountInfo(context.Background(), &Address{MijinTest, nemTestAddress1})

	assert.Nilf(t, err, "AccountService.GetAccountInfo returned error: %s", err)

	tests.ValidateStringers(t, account, acc)
}

func TestAccountService_GetAccountsInfo(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     "/account",
		RespBody: "[" + accountInfoJson + "]",
	})

	accounts, err := accountClient.GetAccountsInfo(
		context.Background(),
		&Address{MijinTest, nemTestAddress1},
	)

	assert.Nilf(t, err, "AccountService.GetAccountsInfo returned error: %s", err)

	for _, acc := range accounts {
		tests.ValidateStringers(t, account, acc)
	}
}

func TestAccountService_GetAccountsNames(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     accountNamesRoute,
		RespBody: "[" + accountNameJson + "]",
	})

	t.Run("return list of names as expect", func(t *testing.T) {

		names, err := accountClient.GetAccountNames(
			context.Background(),
			&Address{MijinTest, accountNames[0].Address.Address},
			&Address{MijinTest, accountNames[1].Address.Address},
		)

		assert.Nilf(t, err, "AccountService.GetAccountNames returned error: %s", err)

		for i, accNames := range names {
			tests.ValidateStringers(t, accountNames[i], accNames)
		}
	})
	t.Run("return error for empty accounts arguments as expect", func(t *testing.T) {

		_, err := accountClient.GetAccountNames(
			context.Background(),
		)

		assert.EqualError(t, err, ErrEmptyAddressesIds.Error())

	})
}

func newAddressFromRaw(addressString string) (address *Address) {
	address, err := NewAddressFromRaw(addressString)
	if err != nil {
		return nil
	}
	return address
}

func TestAccountService_GetAccountHarvestingByAddress(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf("/account/%s/harvesting", harvesterAddress),
		RespBody: accountHarvestingJson,
	})

	h, err := accountClient.GetAccountHarvesting(context.Background(), NewAddress(harvesterAddress, Public))

	assert.Nilf(t, err, "AccountService.GetAccountHarvesting returned error: %s", err)

	tests.ValidateStringers(t, harvester1, h)
}

func TestAccountService_GetHarvesters(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     "/harvesters",
		RespBody: harvestersPageJson,
	})

	h, err := accountClient.GetHarvesters(
		context.Background(),
		&PaginationOrderingOptions{
			PageSize:   2,
			PageNumber: 1,
		},
	)

	assert.Nilf(t, err, "AccountService.GetHarvesters returned error: %s", err)

	tests.ValidateStringers(t, harvesters, h)
}
