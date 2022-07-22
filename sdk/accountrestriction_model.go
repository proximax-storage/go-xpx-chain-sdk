// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
)

type AccountRestrictionFlags uint16

const (
	AccountRestrictionFlag_Address         AccountRestrictionFlags = 0x0001
	AccountRestrictionFlag_MosaicId        AccountRestrictionFlags = 0x0002
	AccountRestrictionFlag_TransactionType AccountRestrictionFlags = 0x0004
	AccountRestrictionFlag_Sentinel        AccountRestrictionFlags = 0x0008
	AccountRestrictionFlag_Outgoing        AccountRestrictionFlags = 0x4000
	AccountRestrictionFlag_Block           AccountRestrictionFlags = 0x8000
)

type AccountRestriction struct {
	RestrictionFlags AccountRestrictionFlags
	Values           []interface{}
}

func (s *AccountRestriction) ToAddressRestriction() AddressAccountRestriction {
	addresses := make([]*Address, len(s.Values))
	for index, val := range s.Values {
		addresses[index] = val.(*Address)
	}
	return AddressAccountRestriction{
		RestrictionFlags: s.RestrictionFlags,
		Values:           addresses,
	}
}
func (s *AccountRestriction) ToMosaicRestriction() MosaicAccountRestriction {
	mosaics := make([]*MosaicId, len(s.Values))
	for index, val := range s.Values {
		mosaics[index] = val.(*MosaicId)
	}
	return MosaicAccountRestriction{
		RestrictionFlags: s.RestrictionFlags,
		Values:           mosaics,
	}
}
func (s *AccountRestriction) ToOperationRestriction() OperationAccountRestriction {
	entityTypes := make([]EntityType, len(s.Values))
	for index, val := range s.Values {
		entityTypes[index] = val.(EntityType)
	}
	return OperationAccountRestriction{
		RestrictionFlags: s.RestrictionFlags,
		Values:           entityTypes,
	}
}

type AddressAccountRestriction struct {
	RestrictionFlags AccountRestrictionFlags
	Values           []*Address
}
type MosaicAccountRestriction struct {
	RestrictionFlags AccountRestrictionFlags
	Values           []*MosaicId
}
type OperationAccountRestriction struct {
	RestrictionFlags AccountRestrictionFlags
	Values           []EntityType
}
type AccountRestrictionsPage struct {
	Restrictions []AccountRestrictions
	Pagination   Pagination
}

type AccountRestrictionsPageOptions struct {
	Address *Address `json:"address"`
	PaginationOrderingOptions
}

func (s *AccountRestriction) String() string {
	return fmt.Sprintf(
		`
			"Flags": %x,
			"Values": %T
		`,
		s.RestrictionFlags,
		s.Values,
	)
}

type AccountRestrictions struct {
	Version      uint32
	Address      *Address
	Restrictions []AccountRestriction
}

func (s *AccountRestrictions) String() string {
	return fmt.Sprintf(
		`
			"Version": %d,
			"Address": %s
			"Restrictions": %T
		`,
		s.Version,
		s.Address.String(),
		s.Restrictions,
	)
}

type AccountAddressRestrictionTransaction struct {
	AbstractTransaction
	RestrictionFlags     uint16
	RestrictionAdditions []*Address
	RestrictionDeletions []*Address
}

type AccountMosaicRestrictionTransaction struct {
	AbstractTransaction
	RestrictionFlags     uint16
	RestrictionAdditions []AssetId
	RestrictionDeletions []AssetId
}

type AccountOperationRestrictionTransaction struct {
	AbstractTransaction
	RestrictionFlags     uint16
	RestrictionAdditions []EntityType
	RestrictionDeletions []EntityType
}
