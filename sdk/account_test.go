// Copyright 2018 ProximaX Limited. All rights reserved. // Use of this source code is governed by the Apache 2.0 // license that can be found in the LICENSE file.
package sdk

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/proximax-storage/go-xpx-utils/tests"
	"github.com/stretchr/testify/assert"
	"testing"
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
)

var (
	nemTestAddress1 = "SAONSOGFZZHNEIBRYXHDTDTBR2YSAXKTITRFHG2Y"
	nemTestAddress2 = "SBJ5D7TFIJWPY56JBEX32MUWI5RU6KVKZYITQ2HA"
	publicKey1      = "27F6BEF9A7F75E33AE2EB2EBA10EF1D6BEA4D30EBD5E39AF8EE06E96E11AE2A9"
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

func TestAccountService_Transactions(t *testing.T) {
	mockServer.AddRouter(&mock.Router{
		Path:     fmt.Sprintf("/account/%s/transactions", publicKey1),
		RespBody: "[" + transactionJson + "]",
	})

	transactions, err := accountClient.Transactions(
		context.Background(),
		&PublicAccount{
			&Address{MijinTest, nemTestAddress2},
			publicKey1,
		},
		&AccountTransactionsOption{},
	)

	assert.Nilf(t, err, "AccountService.Transactions returned error: %s", err)

	for _, tx := range transactions {
		tests.ValidateStringers(t, transaction, tx)
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
