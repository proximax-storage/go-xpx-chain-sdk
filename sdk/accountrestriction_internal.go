// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/base32"
	"encoding/hex"
)

type AccountRestrictionDto struct {
	RestrictionFlags AccountRestrictionFlags `json:"restrictionFlags"`
	Values           []interface{}           `json:"values"`
}

type AccountRestrictionsDtoContainer struct {
	AccountRestrictions AccountRestrictionsDto `json:"accountRestrictions"`
}

type AccountRestrictionsDto struct {
	Version      uint32                  `json:"version"`
	Address      string                  `json:"address"`
	Restrictions []AccountRestrictionDto `json:"restrictions"`
}

func (ref *AccountRestrictionDto) toStruct(networkType NetworkType) (*AccountRestriction, error) {
	accountRestriction := AccountRestriction{}
	accountRestriction.RestrictionFlags = ref.RestrictionFlags
	values := make([]interface{}, len(ref.Values))
	if ref.RestrictionFlags&AccountRestrictionFlag_Address == AccountRestrictionFlag_Address {
		for i, value := range ref.Values {
			bytes, err := hex.DecodeString(value.(string))
			if err != nil {
				return nil, err
			}
			val := NewAddress(base32.StdEncoding.EncodeToString(bytes), networkType)
			values[i] = val
		}
	} else if ref.RestrictionFlags&AccountRestrictionFlag_MosaicId == AccountRestrictionFlag_MosaicId {
		for i, value := range ref.Values {
			val := value.([]interface{})
			mosaicIdDto := mosaicIdDTO{uint32(val[0].(float64)), uint32(val[1].(float64))}
			mosaicId, err := (&mosaicIdDto).toStruct()
			if err != nil {
				return nil, err
			}
			values[i] = mosaicId
		}
	} else if ref.RestrictionFlags&AccountRestrictionFlag_TransactionType == AccountRestrictionFlag_TransactionType {
		for i, value := range ref.Values {
			values[i] = EntityType(value.(float64))
		}
	}
	accountRestriction.Values = values

	return &accountRestriction, nil
}

func (ref *AccountRestrictionsDtoContainer) toStruct(networkType NetworkType) (*AccountRestrictions, error) {
	accountRestrictions := AccountRestrictions{}
	accountRestrictions.Version = ref.AccountRestrictions.Version
	bytes, err := hex.DecodeString(ref.AccountRestrictions.Address)
	if err != nil {
		return nil, err
	}
	accountRestrictions.Address = NewAddress(base32.StdEncoding.EncodeToString(bytes), networkType)
	restrictions := make([]AccountRestriction, len(ref.AccountRestrictions.Restrictions))
	accountRestrictions.Restrictions = restrictions
	for i, restriction := range ref.AccountRestrictions.Restrictions {
		rst, err := restriction.toStruct(networkType)
		if err != nil {
			return nil, err
		}
		accountRestrictions.Restrictions[i] = *rst
	}

	return &accountRestrictions, nil
}

type AccountRestrictionsPageDTO struct {
	AccountRestrictions []AccountRestrictionsDtoContainer `json:"data"`

	Pagination struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (t *AccountRestrictionsPageDTO) toStruct(networkType NetworkType) (*AccountRestrictionsPage, error) {
	page := &AccountRestrictionsPage{
		Restrictions: make([]AccountRestrictions, len(t.AccountRestrictions)),
		Pagination: Pagination{
			TotalEntries: t.Pagination.TotalEntries,
			PageNumber:   t.Pagination.PageNumber,
			PageSize:     t.Pagination.PageSize,
			TotalPages:   t.Pagination.TotalPages,
		},
	}
	for i, t := range t.AccountRestrictions {
		restrictions, err := t.toStruct(networkType)
		if err != nil {
			return nil, err
		}
		page.Restrictions[i] = *restrictions
	}

	return page, nil
}
