package sdk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
)

const (
	testAccountRestrictionEntryJson = `{
		"accountRestrictions": {
			"version": 1,
			"address": "905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3X",
			"restrictions": [
				{
					"restrictionFlags": 1,
					"values": [
						"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3A",
						"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3B",
						"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3C"
					]
				},
				{
					"restrictionFlags": 2,
					"values": [
						[1231231, 0],
						[1231232, 0],
						[1231233, 0]
					]
				},
				{
					"restrictionFlags": 4,
					"values": [
						123,
						124,
						125
					]
				}
			]
		}
	}`
	testAccountRestrictionsSearchResultJson = `
	{
		"data": [
			{
				"accountRestrictions": {
					"version": 1,
					"address": "905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3X",
					"restrictions": [
						{
							"restrictionFlags": 1,
							"values": [
								"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3A",
								"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3B",
								"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3C"
							]
						},
						{
							"restrictionFlags": 2,
							"values": [
								[1231231, 0],
								[1231232, 0],
								[1231233, 0]
							]
						},
						{
							"restrictionFlags": 4,
							"values": [
								123,
								124,
								125
							]
						}
					]
				}
			},
			{
				"accountRestrictions": {
					"version": 1,
					"address": "905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3Y",
					"restrictions": [
						{
							"restrictionFlags": 1,
							"values": [
								"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3A",
								"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3B",
								"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3C"
							]
						},
						{
							"restrictionFlags": 2,
							"values": [
								[1231231, 0],
								[1231232, 0],
								[1231233, 0]
							]
						},
						{
							"restrictionFlags": 4,
							"values": [
								123,
								124,
								125
							]
						}
					]
				}
			},
			{
				"accountRestrictions": {
					"version": 1,
					"address": "905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3Z",
					"restrictions": [
						{
							"restrictionFlags": 1,
							"values": [
								"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3A",
								"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3B",
								"905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3C"
							]
						},
						{
							"restrictionFlags": 2,
							"values": [
								[1231231, 0],
								[1231232, 0],
								[1231233, 0]
							]
						},
						{
							"restrictionFlags": 4,
							"values": [
								123,
								124,
								125
							]
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
	testAccountRestrictionEntryJsonArr = "[" + testAccountRestrictionEntryJson + "]"
)

func newAddressFromRawForceType(address string, networkType NetworkType) *Address {
	addressResult := newAddressFromRaw(address)
	addressResult.Type = networkType
	return addressResult
}

var (
	testAccountRestrictionEntry = AccountRestrictions{
		Version: 1,
		Address: NewAddress("905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3X", PublicTest),
		Restrictions: []AccountRestriction{
			AccountRestriction{
				RestrictionFlags: AccountRestrictionFlag_Address,
				Values: []interface{}{
					newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR2", NetworkType(168)),
					newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR3", NetworkType(168)),
					newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR4", NetworkType(168)),
				},
			},
			AccountRestriction{
				RestrictionFlags: AccountRestrictionFlag_MosaicId,
				Values: []interface{}{
					newMosaicIdPanic(uint64DTO{1231231, 0}.toUint64()),
					newMosaicIdPanic(uint64DTO{1231232, 0}.toUint64()),
					newMosaicIdPanic(uint64DTO{1231233, 0}.toUint64()),
				},
			},
			AccountRestriction{
				RestrictionFlags: AccountRestrictionFlag_TransactionType,
				Values: []interface{}{
					EntityType(123),
					EntityType(124),
					EntityType(125),
				},
			},
		},
	}
	testAccountRestrictionsSearchResult = &AccountRestrictionsPage{
		Restrictions: []AccountRestrictions{
			{
				Version: 1,
				Address: NewAddress("905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3X", PublicTest),
				Restrictions: []AccountRestriction{
					AccountRestriction{
						RestrictionFlags: AccountRestrictionFlag_Address,
						Values: []interface{}{
							newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR2", NetworkType(168)),
							newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR3", NetworkType(168)),
							newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR4", NetworkType(168)),
						},
					},
					AccountRestriction{
						RestrictionFlags: AccountRestrictionFlag_MosaicId,
						Values: []interface{}{
							newMosaicIdPanic(uint64DTO{1231231, 0}.toUint64()),
							newMosaicIdPanic(uint64DTO{1231232, 0}.toUint64()),
							newMosaicIdPanic(uint64DTO{1231233, 0}.toUint64()),
						},
					},
					AccountRestriction{
						RestrictionFlags: AccountRestrictionFlag_TransactionType,
						Values: []interface{}{
							EntityType(123),
							EntityType(124),
							EntityType(125),
						},
					},
				},
			},
			{
				Version: 1,
				Address: NewAddress("905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3Y", PublicTest),
				Restrictions: []AccountRestriction{
					AccountRestriction{
						RestrictionFlags: AccountRestrictionFlag_Address,
						Values: []interface{}{
							newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR2", NetworkType(168)),
							newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR3", NetworkType(168)),
							newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR4", NetworkType(168)),
						},
					},
					AccountRestriction{
						RestrictionFlags: AccountRestrictionFlag_MosaicId,
						Values: []interface{}{
							newMosaicIdPanic(uint64DTO{1231231, 0}.toUint64()),
							newMosaicIdPanic(uint64DTO{1231232, 0}.toUint64()),
							newMosaicIdPanic(uint64DTO{1231233, 0}.toUint64()),
						},
					},
					AccountRestriction{
						RestrictionFlags: AccountRestrictionFlag_TransactionType,
						Values: []interface{}{
							EntityType(123),
							EntityType(124),
							EntityType(125),
						},
					},
				},
			},
			{
				Version: 1,
				Address: NewAddress("905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3Z", PublicTest),
				Restrictions: []AccountRestriction{
					AccountRestriction{
						RestrictionFlags: AccountRestrictionFlag_Address,
						Values: []interface{}{
							newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR2", NetworkType(168)),
							newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR3", NetworkType(168)),
							newAddressFromRawForceType("SBN5BDMFV4ZCJJRMF3NLABGP6RBSE4PGMKZTHOR4", NetworkType(168)),
						},
					},
					AccountRestriction{
						RestrictionFlags: AccountRestrictionFlag_MosaicId,
						Values: []interface{}{
							newMosaicIdPanic(uint64DTO{1231231, 0}.toUint64()),
							newMosaicIdPanic(uint64DTO{1231232, 0}.toUint64()),
							newMosaicIdPanic(uint64DTO{1231233, 0}.toUint64()),
						},
					},
					AccountRestriction{
						RestrictionFlags: AccountRestrictionFlag_TransactionType,
						Values: []interface{}{
							EntityType(123),
							EntityType(124),
							EntityType(125),
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

func Test_AccountRestrictionService_GetAccountRestrictions(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(accountRestrictionsRoute, "905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3X"),
		AcceptedHttpMethods: []string{http.MethodGet},
		RespHttpCode:        200,
		RespBody:            testAccountRestrictionEntryJson,
	})
	accountRestrictionClient := mock.getPublicTestClientUnsafe().AccountRestriction

	defer mock.Close()

	record, err := accountRestrictionClient.GetAccountRestrictions(ctx, NewAddress("905BD08D85AF3224A62C2EDAB004CFF4432271E662B333BA3X", PublicTest))
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, testAccountRestrictionEntry, *record)
}

func Test_AccountRestrictionService_SearchAccountRestrictions(t *testing.T) {
	mock := newSdkMockWithRouter(&mock.Router{
		Path:                fmt.Sprintf(accountRestrictionsRoute, ""),
		AcceptedHttpMethods: []string{http.MethodPost},
		RespHttpCode:        200,
		RespBody:            testAccountRestrictionsSearchResultJson,
	})
	accountRestrictionClient := mock.getPublicTestClientUnsafe().AccountRestriction

	defer mock.Close()

	options := AccountRestrictionsPageOptions{
		Address: nil,
		PaginationOrderingOptions: PaginationOrderingOptions{
			PageSize:      20,
			PageNumber:    1,
			Offset:        "",
			SortField:     "",
			SortDirection: "",
		},
	}
	record, err := accountRestrictionClient.SearchAccountRestrictions(ctx, &options)
	assert.Nil(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, testAccountRestrictionsSearchResult, record)
}
