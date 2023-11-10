// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/base32"
	"encoding/hex"
	jsoniter "github.com/json-iterator/go"
)

type MosaicRestrictionDto struct {
	Key   AccountRestrictionFlags
	Value interface{}
}

type GlobalRestrictionValueDto struct {
	ReferenceMosaicId mosaicIdHexDTO        `json:"referenceMosaicId"`
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

type MosaicRestrictionEntryPartialDto struct {
	Version       uint32                     `json:"version"`
	CompositeHash string                     `json:"compositeHash"`
	EntryType     MosaicRestrictionEntryType `json:"entryType"`
	MosaicId      mosaicIdDTO                `json:"mosaicId"`
	Address       string                     `json:"targetAddress"`
	Restrictions  jsoniter.RawMessage        `json:"restrictions"`
}
type MosaicRestrictionEntryDto struct {
	Version       uint32                     `json:"version"`
	CompositeHash string                     `json:"compositeHash"`
	EntryType     MosaicRestrictionEntryType `json:"entryType"`
	MosaicId      mosaicIdDTO                `json:"mosaicId"`
	Address       string                     `json:"targetAddress"`
	Restrictions  []interface{}              `json:"restrictions"`
}

// /TODO: Optimize unmarshal and toStruct to duplicate less work
func (m *MosaicRestrictionEntryDto) UnmarshalJSON(data []byte) error {
	var partialDto MosaicRestrictionEntryPartialDto
	if err := json.Unmarshal(data, &partialDto); err != nil {
		return err
	}
	if partialDto.EntryType == MosaicRestrictionEntryType_Address {
		restrictions := make([]AddressMosaicRestrictionDto, 0)
		if err := json.Unmarshal(partialDto.Restrictions, &restrictions); err != nil {
			return err
		}
		m.Restrictions = make([]interface{}, len(restrictions))
		for i, _ := range restrictions {
			m.Restrictions[i] = &restrictions[i]
		}
	} else {
		restrictions := make([]GlobalMosaicRestrictionDto, 0)
		if err := json.Unmarshal(partialDto.Restrictions, &restrictions); err != nil {
			return err
		}
		m.Restrictions = make([]interface{}, len(restrictions))
		for i, _ := range restrictions {
			m.Restrictions[i] = &restrictions[i]
		}
	}
	m.Version = partialDto.Version
	m.Address = partialDto.Address
	m.CompositeHash = partialDto.CompositeHash
	m.MosaicId = partialDto.MosaicId
	m.EntryType = partialDto.EntryType
	return nil
}

type MosaicRestrictionEntryContainerDto struct {
	MosaicRestrictionEntry MosaicRestrictionEntryDto `json:"mosaicRestrictionEntry"`
}

func (ref *MosaicRestrictionEntryContainerDto) toStruct(networkType NetworkType) (*MosaicRestrictionEntry, error) {
	return ref.MosaicRestrictionEntry.toStruct(networkType)
}

func (ref *MosaicRestrictionEntryDto) toStruct(networkType NetworkType) (*MosaicRestrictionEntry, error) {
	mosaicRestriction := MosaicRestrictionEntry{}
	mosaicId, err := ref.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	mosaicRestriction.MosaicId = mosaicId
	mosaicRestriction.CompositeHash = ref.CompositeHash
	bytes, err := hex.DecodeString(ref.Address)
	mosaicRestriction.Address = NewAddress(base32.StdEncoding.EncodeToString(bytes), networkType)
	mosaicRestriction.Version = ref.Version
	mosaicRestriction.EntryType = ref.EntryType
	mosaicRestriction.Restrictions = make([]interface{}, len(ref.Restrictions))
	for i, restriction := range ref.Restrictions {
		if mosaicRestriction.EntryType == MosaicRestrictionEntryType_Address {
			mosaicRestriction.Restrictions[i] = restriction.(*AddressMosaicRestrictionDto).toStruct()
		} else {
			mosaicRestrictionPtr, err := restriction.(*GlobalMosaicRestrictionDto).toStruct()
			if err != nil {
				return nil, err
			}
			mosaicRestriction.Restrictions[i] = mosaicRestrictionPtr
		}
	}
	return &mosaicRestriction, nil
}

type MosaicRestrictionsPageDto struct {
	MosaicRestrictions []MosaicRestrictionEntryContainerDto `json:"data"`
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
		restrictions, err := t.toStruct(networkType)
		if err != nil {
			return nil, err
		}
		page.Restrictions[i] = *restrictions
	}

	return page, nil
}
