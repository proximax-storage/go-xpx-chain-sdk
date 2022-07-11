package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testMosaicRestrictionEntryAddressJson = `
	{
		"id": "randomid",
		"mosaicRestrictionEntry": {
			"version": 1,
			"compositeHash": "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810",
			"entryType": 0,
			"targetAddress": "905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3A",
			"restrictions": [
				{
					"key": [123, 0],
					"value": [123, 0]
				},
				{
					"key": [124, 0],
					"value": [124, 0]
				},
				{
					"key": [125, 0],
					"value": [125, 0]
				}
			]
		}
	}`
	testMosaicRestrictionEntryGlobalJson = `
	{
		"id": "string",
		"mosaicRestrictionEntry": {
			"version": 1,
			"compositeHash": "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810",
			"entryType": 1,
			"mosaicId": [123, 0],
			"targetAddress": "",
			"restrictions": [
				{
					"key": [123, 0],
					"restriction": {
						"referenceMosaicId": [124, 0],
						"restrictionValue": [125, 0],
						"restrictionType": 0
					}
				},
				{
					"key": [124, 0],
					"restriction": {
						"referenceMosaicId": [124, 0],
						"restrictionValue": [125, 0],
						"restrictionType": 0
					}
				},
				{
					"key": [125, 0],
					"restriction": {
						"referenceMosaicId": [124, 0],
						"restrictionValue": [125, 0],
						"restrictionType": 0
					}
				}
			]
		}
	}`
	testMosaicRestrictionsSearchResultJson = `
	{
		"data": [
			{
				"id": "string",
				"mosaicRestrictionEntry": {
					"version": 1,
					"compositeHash": "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810",
					"entryType": 0,
					"targetAddress": "",
					"restrictions": [
						{
							"key": [123, 0],
							"value": [123, 0]
						},
						{
							"key": [124, 0],
							"value": [124, 0]
						},
						{
							"key": [125, 0],
							"value": [125, 0]
						}
					]
				}
			},
			{
				"id": "string",
				"mosaicRestrictionEntry": {
					"version": 1,
					"compositeHash": "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23811",
					"entryType": 1,
					"mosaicId": [124, 0],
					"targetAddress": "",
					"restrictions": [
						{
							"key": [123, 0],
							"restriction": {
								"referenceMosaicId": [124, 0],
								"restrictionValue": [125, 0],
								"restrictionType": 0
							}
						},
						{
							"key": [124, 0],
							"restriction": {
								"referenceMosaicId": [124, 0],
								"restrictionValue": [125, 0],
								"restrictionType": 0
							}
						},
						{
							"key": [125, 0],
							"restriction": {
								"referenceMosaicId": [124, 0],
								"restrictionValue": [125, 0],
								"restrictionType": 0
							}
						}
					]
				}
			},
			{
				"id": "string3",
				"mosaicRestrictionEntry": {
					"version": 1,
					"compositeHash": "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23812",
					"entryType": 1,
					"mosaicId": [125, 0],
					"targetAddress": "",
					"restrictions": [
						{
							"key": [123, 0],
							"restriction": {
								"referenceMosaicId": [124, 0],
								"restrictionValue": [125, 0],
								"restrictionType": 0
							}
						}
					]
				}
			}
		],
		"pagination": {
			"totalEntries": 3,
			"pageNumber": 1,
			"pageSize": 20,
			"totalPages": 1
		}
	}
	`
)

var (
	testMosaicAddressRestrictionEntry = MosaicRestrictionEntry{
		Version:       1,
		CompositeHash: "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810",
		MosaicId:      &MosaicId{},
		EntryType:     0,
		Restrictions: []interface{}{
			&AddressMosaicRestriction{
				Key:   uint64DTO{123, 0}.toUint64(),
				Value: uint64DTO{123, 0}.toUint64(),
			},
			&AddressMosaicRestriction{
				Key:   uint64DTO{124, 0}.toUint64(),
				Value: uint64DTO{124, 0}.toUint64(),
			},
			&AddressMosaicRestriction{
				Key:   uint64DTO{125, 0}.toUint64(),
				Value: uint64DTO{125, 0}.toUint64(),
			},
		},
	}
	testMosaicGlobalRestrictionEntry = MosaicRestrictionEntry{
		Version:       1,
		EntryType:     1,
		MosaicId:      newMosaicIdPanic(uint64DTO{123, 0}.toUint64()),
		CompositeHash: "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810",
		Restrictions: []interface{}{
			&GlobalMosaicRestriction{
				Key: uint64DTO{123, 0}.toUint64(),
				Value: GlobalRestrictionValue{
					ReferenceMosaicId: newMosaicIdPanic(uint64DTO{124, 0}.toUint64()),
					RestrictionValue:  uint64DTO{125, 0}.toUint64(),
					RestrictionType:   0,
				},
			},
			&GlobalMosaicRestriction{
				Key: uint64DTO{124, 0}.toUint64(),
				Value: GlobalRestrictionValue{
					ReferenceMosaicId: newMosaicIdPanic(uint64DTO{124, 0}.toUint64()),
					RestrictionValue:  uint64DTO{125, 0}.toUint64(),
					RestrictionType:   0,
				},
			},
			&GlobalMosaicRestriction{
				Key: uint64DTO{125, 0}.toUint64(),
				Value: GlobalRestrictionValue{
					ReferenceMosaicId: newMosaicIdPanic(uint64DTO{124, 0}.toUint64()),
					RestrictionValue:  uint64DTO{125, 0}.toUint64(),
					RestrictionType:   0,
				},
			},
		},
	}

	testMosaicRestrictionsSearchResult = MosaicRestrictionsPage{
		Restrictions: []MosaicRestrictionEntry{
			{
				Version:       1,
				MosaicId:      &MosaicId{},
				EntryType:     0,
				CompositeHash: "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810",
				Restrictions: []interface{}{
					&AddressMosaicRestriction{
						Key:   uint64DTO{123, 0}.toUint64(),
						Value: uint64DTO{123, 0}.toUint64(),
					},
					&AddressMosaicRestriction{
						Key:   uint64DTO{124, 0}.toUint64(),
						Value: uint64DTO{124, 0}.toUint64(),
					},
					&AddressMosaicRestriction{
						Key:   uint64DTO{125, 0}.toUint64(),
						Value: uint64DTO{125, 0}.toUint64(),
					},
				},
			},
			{
				Version:       1,
				EntryType:     1,
				MosaicId:      newMosaicIdPanic(uint64DTO{124, 0}.toUint64()),
				CompositeHash: "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23811",
				Restrictions: []interface{}{
					&GlobalMosaicRestriction{
						Key: uint64DTO{123, 0}.toUint64(),
						Value: GlobalRestrictionValue{
							ReferenceMosaicId: newMosaicIdPanic(uint64DTO{124, 0}.toUint64()),
							RestrictionValue:  uint64DTO{125, 0}.toUint64(),
							RestrictionType:   0,
						},
					},
					&GlobalMosaicRestriction{
						Key: uint64DTO{124, 0}.toUint64(),
						Value: GlobalRestrictionValue{
							ReferenceMosaicId: newMosaicIdPanic(uint64DTO{124, 0}.toUint64()),
							RestrictionValue:  uint64DTO{125, 0}.toUint64(),
							RestrictionType:   0,
						},
					},
					&GlobalMosaicRestriction{
						Key: uint64DTO{125, 0}.toUint64(),
						Value: GlobalRestrictionValue{
							ReferenceMosaicId: newMosaicIdPanic(uint64DTO{124, 0}.toUint64()),
							RestrictionValue:  uint64DTO{125, 0}.toUint64(),
							RestrictionType:   0,
						},
					},
				},
			},
			{
				Version:       1,
				MosaicId:      newMosaicIdPanic(uint64DTO{125, 0}.toUint64()),
				EntryType:     1,
				CompositeHash: "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23812",
				Restrictions: []interface{}{
					&GlobalMosaicRestriction{
						Key: uint64DTO{123, 0}.toUint64(),
						Value: GlobalRestrictionValue{
							ReferenceMosaicId: newMosaicIdPanic(uint64DTO{124, 0}.toUint64()),
							RestrictionValue:  uint64DTO{125, 0}.toUint64(),
							RestrictionType:   0,
						},
					},
				},
			},
		},
		Pagination: Pagination{
			TotalEntries: 3,
			PageNumber:   1,
			PageSize:     20,
			TotalPages:   1,
		},
	}
)

func Test_MosaicRestrictionService_GetMosaicRestrictions_Address(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(mosaicRestrictionsRoute, "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810"),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testMosaicRestrictionEntryAddressJson,
	})
	accountRestrictionClient := mock.getPublicTestClientUnsafe().MosaicRestriction

	defer mock.Close()

	record, err := accountRestrictionClient.GetMosaicRestrictions(ctx, "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810")
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, testMosaicAddressRestrictionEntry, *record)
}

func Test_MosaicRestrictionService_GetMosaicRestrictions_Global(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(mosaicRestrictionsRoute, "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810"),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testMosaicRestrictionEntryGlobalJson,
	})
	accountRestrictionClient := mock.getPublicTestClientUnsafe().MosaicRestriction

	defer mock.Close()

	record, err := accountRestrictionClient.GetMosaicRestrictions(ctx, "C8FC3FB54FDDFBCE0E8C71224990124E4EEC5AD5D30E592EDFA9524669A23810")
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, testMosaicGlobalRestrictionEntry, *record)
}

func Test_MosaicRestrictionService_SearchAccountRestrictions(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(mosaicRestrictionsRoute, ""),
		AcceptedHttpMethods: []string{http.MethodPost},
		RespHttpCode:        200,
		RespBody:            testMosaicRestrictionsSearchResultJson,
	})
	mosaicRestrictionClient := mock.getPublicTestClientUnsafe().MosaicRestriction

	defer mock.Close()

	options := MosaicRestrictionsPageOptions{
		MosaicId:  nil,
		EntryType: nil,
		Address:   nil,
		PaginationOrderingOptions: PaginationOrderingOptions{
			PageSize:      20,
			PageNumber:    1,
			Offset:        "",
			SortField:     "",
			SortDirection: "",
		},
	}
	record, err := mosaicRestrictionClient.SearchMosaicRestrictions(ctx, &options)
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, testMosaicRestrictionsSearchResult, *record)
}
