// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type MosaicRestrictionDto struct {
	Key   AccountRestrictionFlags
	Value interface{}
}

type MosaicRestrictionsDtoContainer struct {
	MosaicRestrictionEntry MosaicRestrictionEntryDto `json:"mosaicRestrictionEntry"`
}

type GlobalRestrictionValueDto struct {
	ReferenceMosaicId mosaicIdDTO           `json:"referenceMosaicId"`
	RestrictionValue  uint64DTO             `json:"restrictionValue"`
	RestrictionType   MosaicRestrictionType `json:"restrictionType"`
}

func (ref *GlobalRestrictionValueDto) toStruct() (*GlobalRestrictionValue, error) {
	mosaicRestriction := GlobalRestrictionValue{}
	mosaicId, err := ref.ReferenceMosaicId.toStruct()
	if err != nil {
		return nil, err
	}
	mosaicRestriction.ReferenceMosaicId = mosaicId
	mosaicRestriction.RestrictionValue = ref.RestrictionValue.toUint64()
	mosaicRestriction.RestrictionType = ref.RestrictionType
	return &mosaicRestriction, nil
}

type AddressMosaicRestrictionDto struct {
	Key   uint64DTO `json:"key"`
	Value uint64DTO `json:"value"`
}

func (ref *AddressMosaicRestrictionDto) toStruct() *AddressMosaicRestriction {
	mosaicRestriction := AddressMosaicRestriction{}
	mosaicRestriction.Key = ref.Key.toUint64()
	mosaicRestriction.Value = ref.Value.toUint64()
	return &mosaicRestriction
}

type GlobalMosaicRestrictionDto struct {
	Key   uint64DTO                 `key`
	Value GlobalRestrictionValueDto `json:"restriction"`
}

func (ref *GlobalMosaicRestrictionDto) toStruct() (*GlobalMosaicRestriction, error) {
	mosaicRestriction := GlobalMosaicRestriction{}
	mosaicRestriction.Key = ref.Key.toUint64()
	value, err := ref.Value.toStruct()
	if err != nil {
		return nil, err
	}
	mosaicRestriction.Value = *value
	return &mosaicRestriction, nil
}

type MosaicRestrictionEntryDto struct {
	Version       uint32                     `json:"version"`
	CompositeHash string                     `json:"compositeHash"`
	EntryType     MosaicRestrictionEntryType `json:"entryType"`
	MosaicId      mosaicIdDTO                `json:"mosaicId"`
	Address       *Address                   `json:"address"`
	Restrictions  []map[string]interface{}   `json:"restrictions"`
}
type MosaicRestrictionEntryDtoContainer struct {
	MosaicRestrictionEntry MosaicRestrictionEntryDto `json:"mosaicRestrictionEntry"`
}

func (ref *MosaicRestrictionEntryDtoContainer) toStruct() (*MosaicRestrictionEntry, error) {
	return ref.MosaicRestrictionEntry.toStruct()
}

func (ref *MosaicRestrictionEntryDto) toStruct() (*MosaicRestrictionEntry, error) {
	mosaicRestriction := MosaicRestrictionEntry{}
	mosaicId, err := ref.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	mosaicRestriction.MosaicId = mosaicId
	mosaicRestriction.CompositeHash = ref.CompositeHash
	mosaicRestriction.Address = ref.Address
	mosaicRestriction.Version = ref.Version
	mosaicRestriction.EntryType = ref.EntryType
	mosaicRestriction.Restrictions = make([]interface{}, len(ref.Restrictions))
	for i, value := range ref.Restrictions {

		if ref.EntryType == MosaicRestrictionEntryType_Address {
			restriction := AddressMosaicRestrictionDto{
				Key: uint64DTO{
					uint32((value["key"].([]interface{}))[0].(float64)),
					uint32((value["key"].([]interface{}))[1].(float64)),
				},
				Value: uint64DTO{
					uint32((value["value"].([]interface{}))[0].(float64)),
					uint32((value["value"].([]interface{}))[1].(float64)),
				},
			}
			restrictionptr := (&restriction).toStruct()
			mosaicRestriction.Restrictions[i] = restrictionptr
		} else {
			val := value["restriction"].(map[string]interface{})
			restriction := GlobalMosaicRestrictionDto{
				Key: uint64DTO{
					uint32((value["key"].([]interface{}))[0].(float64)),
					uint32((value["key"].([]interface{}))[1].(float64)),
				},
				Value: GlobalRestrictionValueDto{
					ReferenceMosaicId: mosaicIdDTO{
						uint32((val["referenceMosaicId"].([]interface{}))[0].(float64)),
						uint32((val["referenceMosaicId"].([]interface{}))[1].(float64)),
					},
					RestrictionValue: uint64DTO{
						uint32((val["restrictionValue"].([]interface{}))[0].(float64)),
						uint32((val["restrictionValue"].([]interface{}))[1].(float64)),
					},
					RestrictionType: MosaicRestrictionType(val["restrictionType"].(float64)),
				},
			}
			restrictionptr, err := restriction.toStruct()
			if err != nil {
				return nil, err
			}
			mosaicRestriction.Restrictions[i] = restrictionptr
		}
	}

	return &mosaicRestriction, nil
}

type MosaicRestrictionsPageDto struct {
	MosaicRestrictions []MosaicRestrictionEntryDtoContainer `json:"data"`
	Pagination         struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (t *MosaicRestrictionsPageDto) toStruct(networkType NetworkType) (*MosaicRestrictionsPage, error) {
	page := &MosaicRestrictionsPage{
		Restrictions: make([]MosaicRestrictionEntry, len(t.MosaicRestrictions)),
		Pagination: Pagination{
			TotalEntries: t.Pagination.TotalEntries,
			PageNumber:   t.Pagination.PageNumber,
			PageSize:     t.Pagination.PageSize,
			TotalPages:   t.Pagination.TotalPages,
		},
	}
	for i, t := range t.MosaicRestrictions {
		restrictions, err := t.toStruct()
		if err != nil {
			return nil, err
		}
		page.Restrictions[i] = *restrictions
	}

	return page, nil
}
