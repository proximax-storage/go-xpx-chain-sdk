package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testKeyRecordGroupEntryJson = `{
		"lockFundRecordGroup": {
			"identifier": "90936FF3536858CBEA8EE0EAAB99FE9EC4EF5EF1F66366569A",
			"records": [
				{
					"key": [
							120,
							0
						],
					"activeRecord": [
					{
						"id":[
							3646934825,
							3576016193
						],
						"amount":[
							10000005,
							0
						]
					}	
					],
					"inactiveRecords": [
						[
							{
								"id":[
									3646934825,
									3576016193
								],
								"amount":[
									10000005,
									0
								]
							}	
						]
					]
				},
				{		
					"key": [
							1233,
							0
						],
					"activeRecord": [
					],
					"inactiveRecords": [
					]
				}
			]
		}
	}`
	testHeightRecordGroupEntryJson = `{
		"lockFundRecordGroup": {
			"identifier": [
				120,
				0
			],
			"records": [
				{
					"key": "90936FF3536858CBEA8EE0EAAB99FE9EC4EF5EF1F66366569A",
					"activeRecord": [
					],
					"inactiveRecords": [
					]
				},
				{		
					"key": "90936FF3536858CBEA8EE0EAAB99FE9EC4EF5EF1F66366569B",
					"activeRecord": [
					{
						"id":[
							3646934825,
							3576016193
						],
						"amount":[
							10000005,
							0
						]
					}	
					],
					"inactiveRecords": [
						[
							{
								"id":[
									3646934825,
									3576016193
								],
								"amount":[
									10000005,
									0
								]
							}	
						]
					]
				}
			]
		}
	}`

	testKeyRecordGroupEntryJsonArr    = "[" + testKeyRecordGroupEntryJson + "]"
	testHeightRecordGroupEntryJsonArr = "[" + testHeightRecordGroupEntryJson + "]"
)

var (
	testKeyRecordKey      = "90936FF3536858CBEA8EE0EAAB99FE9EC4EF5EF1F66366569A"
	testPublicAccount, _  = NewAccountFromPublicKey(testKeyRecordKey, PublicTest)
	testPublicAccount2, _ = NewAccountFromPublicKey("90936FF3536858CBEA8EE0EAAB99FE9EC4EF5EF1F66366569B", PublicTest)

	testLockFundRecordEmpty = LockFundRecord{
		ActiveRecord:    []*Mosaic{},
		InactiveRecords: []*[]*Mosaic{},
	}
	testLockFundRecord = LockFundRecord{
		ActiveRecord: []*Mosaic{
			newMosaicPanic(newAssetIdPanic(uint64DTO{3646934825, 3576016193}), Amount(uint64DTO{10000005, 0}.toUint64())),
		},
		InactiveRecords: []*[]*Mosaic{
			{
				newMosaicPanic(newAssetIdPanic(uint64DTO{3646934825, 3576016193}), Amount(uint64DTO{10000005, 0}.toUint64())),
			},
		},
	}

	testKeyRecordGroupEntry = &LockFundKeyRecord{
		Identifier: testPublicAccount,
		Records: map[Height]*LockFundRecord{
			Height(120):  &testLockFundRecord,
			Height(1233): &testLockFundRecordEmpty,
		},
	}

	testHeightRecordGroupEntry = &LockFundHeightRecord{
		Identifier: Height(120),
		Records: map[string]*LockFundRecord{
			testPublicAccount.PublicKey:  &testLockFundRecordEmpty,
			testPublicAccount2.PublicKey: &testLockFundRecord,
		},
	}
)

func Test_LockFundService_GetHeightRecordGroupEntry(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(lockFundHeightRecordGroupRoute, Height(120)),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testHeightRecordGroupEntryJson,
	})
	lockFundClient := mock.getPublicTestClientUnsafe().LockFund

	defer mock.Close()

	record, err := lockFundClient.GetLockFundHeightRecords(ctx, Height(120))
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, testHeightRecordGroupEntry, record)
}

func Test_LockFundService_GetKeyRecordGroupEntry(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(lockFundKeyRecordGroupRoute, testPublicAccount.PublicKey),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testKeyRecordGroupEntryJson,
	})
	lockFundClient := mock.getPublicTestClientUnsafe().LockFund

	defer mock.Close()

	record, err := lockFundClient.GetLockFundKeyRecords(ctx, testPublicAccount)
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, testKeyRecordGroupEntry, record)
}
