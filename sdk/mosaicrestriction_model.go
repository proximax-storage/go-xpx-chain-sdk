// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import "errors"

type MosaicRestrictionType uint8

type MosaicRestrictionEntryType uint8

const (
	MosaicRestrictionType_NONE MosaicRestrictionType = iota
	MosaicRestrictionType_EQ
	MosaicRestrictionType_NE
	MosaicRestrictionType_LT
	MosaicRestrictionType_LE
	MosaicRestrictionType_GT
	MosaicRestrictionType_GE
)

const (
	MosaicRestrictionEntryType_Address MosaicRestrictionEntryType = iota
	MosaicRestrictionEntryType_Global
)

type MosaicRestrictionsPageOptions struct {
	MosaicId  *MosaicId                   `json:"mosaicId" url:"mosaicId,omitempty"`
	EntryType *MosaicRestrictionEntryType `json:"entryType" url:"entryType,omitempty"`
	Address   *Address                    `json:"targetAddress" url:"targetAddress,omitempty"`
	PaginationOrderingOptions
}

type GlobalRestrictionValue struct {
	ReferenceMosaicId *MosaicId
	RestrictionValue  uint64
	RestrictionType   MosaicRestrictionType
}

type AddressRestrictionValue uint64

type MosaicRestriction struct {
	Key   uint64
	Value interface{}
}

type AddressMosaicRestriction struct {
	Key   uint64
	Value uint64
}

type GlobalMosaicRestriction struct {
	Key   uint64
	Value GlobalRestrictionValue
}

type AddressMosaicRestrictionEntry struct {
	Version       uint32
	CompositeHash string
	MosaicId      *MosaicId
	Address       *Address
	Restrictions  []AddressMosaicRestriction
}

type GlobalMosaicRestrictionEntry struct {
	Version       uint32
	CompositeHash string
	MosaicId      *MosaicId
	Restrictions  []GlobalMosaicRestriction
}

type MosaicRestrictionEntry struct {
	Version       uint32
	CompositeHash string
	EntryType     MosaicRestrictionEntryType
	MosaicId      *MosaicId
	Address       *Address
	Restrictions  []interface{}
}

func (ref *MosaicRestrictionEntry) ToAddressMosaicRestrictionEntry() (*AddressMosaicRestrictionEntry, error) {
	if ref.EntryType != MosaicRestrictionEntryType_Address {
		return nil, errors.New("Entry type does not support this conversion")
	}
	mosaicRestriction := AddressMosaicRestrictionEntry{}

	mosaicRestriction.MosaicId = ref.MosaicId
	mosaicRestriction.CompositeHash = ref.CompositeHash
	mosaicRestriction.Address = ref.Address
	mosaicRestriction.Version = ref.Version
	mosaicRestriction.Restrictions = make([]AddressMosaicRestriction, len(ref.Restrictions))
	for i, value := range ref.Restrictions {
		restriction := value.(*AddressMosaicRestriction)
		mosaicRestriction.Restrictions[i] = *restriction
	}
	return &mosaicRestriction, nil
}

func (ref *MosaicRestrictionEntry) ToGlobalMosaicRestrictionEntry() (*GlobalMosaicRestrictionEntry, error) {
	if ref.EntryType != MosaicRestrictionEntryType_Global {
		return nil, errors.New("Entry type does not support this conversion")
	}
	mosaicRestriction := GlobalMosaicRestrictionEntry{}

	mosaicRestriction.MosaicId = ref.MosaicId
	mosaicRestriction.CompositeHash = ref.CompositeHash
	mosaicRestriction.Version = ref.Version
	mosaicRestriction.Restrictions = make([]GlobalMosaicRestriction, len(ref.Restrictions))
	for i, value := range ref.Restrictions {
		restriction := value.(*GlobalMosaicRestriction)
		mosaicRestriction.Restrictions[i] = *restriction
	}
	return &mosaicRestriction, nil
}

type MosaicRestrictionsPage struct {
	Restrictions []MosaicRestrictionEntry
	Pagination   Pagination
}

type MosaicAddressRestrictionTransaction struct {
	AbstractTransaction
	MosaicId                 AssetId
	RestrictionKey           uint64
	PreviousRestrictionValue uint64
	NewRestrictionValue      uint64
	TargetAddress            *Address
}

type MosaicGlobalRestrictionTransaction struct {
	AbstractTransaction
	MosaicId                 AssetId
	ReferenceMosaicId        AssetId
	RestrictionKey           uint64
	PreviousRestrictionValue uint64
	PreviousRestrictionType  MosaicRestrictionType
	NewRestrictionValue      uint64
	NewRestrictionType       MosaicRestrictionType
}
