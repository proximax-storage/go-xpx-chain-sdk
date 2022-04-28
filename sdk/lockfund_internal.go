// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type LockFundHeightRecordGroupDto struct {
	LockFundRecordGroup struct {
		Identifier uint64DTO                  `json:"identifier"`
		Records    []*LockFundHeightRecordDto `json:"records"`
	}
}
type InactiveRecordsDto []*InactiveRecordDto

type InactiveRecordDto struct {
	Mosaics []*mosaicDTO `json:"mosaics"`
}

type LockFundHeightRecordDto struct {
	Key             string
	ActiveMosaics   []*mosaicDTO       `json:"activeMosaics"`
	InactiveRecords InactiveRecordsDto `json:"inactiveRecords"`
}

type LockFundKeyRecordDto struct {
	Key             uint64DTO          `json:"key"`
	ActiveMosaics   []*mosaicDTO       `json:"activeMosaics"`
	InactiveRecords InactiveRecordsDto `json:"inactiveRecords"`
}
type LockFundKeyRecordGroupDto struct {
	LockFundRecordGroup struct {
		Identifier string                  `json:"identifier"`
		Records    []*LockFundKeyRecordDto `json:"records"`
	}
}

func (ref *LockFundHeightRecordDto) toStruct(networkType NetworkType) (*LockFundRecord, *string, error) {
	lockFundRecord := LockFundRecord{}
	lockFundRecord.ActiveRecord = make([]*Mosaic, len(ref.ActiveMosaics))
	for i, mosaic := range ref.ActiveMosaics {
		msc, err := mosaic.toStruct()
		if err != nil {
			return nil, nil, err
		}

		lockFundRecord.ActiveRecord[i] = msc
	}

	inactiveRecords := make([]*([]*Mosaic), 0)
	for _, inactiveRecord := range ref.InactiveRecords {
		record := make([]*Mosaic, len(inactiveRecord.Mosaics))
		for i, mosaic := range inactiveRecord.Mosaics {
			msc, err := mosaic.toStruct()
			if err != nil {
				return nil, nil, err
			}

			record[i] = msc
		}
		inactiveRecords = append(inactiveRecords, &record)
	}
	lockFundRecord.InactiveRecords = inactiveRecords

	return &lockFundRecord, &ref.Key, nil
}
func (ref *LockFundKeyRecordDto) toStruct(networkType NetworkType) (*LockFundRecord, Height, error) {
	lockFundRecord := LockFundRecord{}
	lockFundRecord.ActiveRecord = make([]*Mosaic, len(ref.ActiveMosaics))
	for i, mosaic := range ref.ActiveMosaics {
		msc, err := mosaic.toStruct()
		if err != nil {
			return nil, 0, err
		}

		lockFundRecord.ActiveRecord[i] = msc
	}

	inactiveRecords := make([]*([]*Mosaic), 0)
	for _, inactiveRecord := range ref.InactiveRecords {
		record := make([]*Mosaic, len(inactiveRecord.Mosaics))
		for i, mosaic := range inactiveRecord.Mosaics {
			msc, err := mosaic.toStruct()
			if err != nil {
				return nil, 0, err
			}

			record[i] = msc
		}
		inactiveRecords = append(inactiveRecords, &record)
	}
	lockFundRecord.InactiveRecords = inactiveRecords

	return &lockFundRecord, Height(ref.Key.toUint64()), nil
}
func (ref *LockFundHeightRecordGroupDto) toStruct(networkType NetworkType) (*LockFundHeightRecord, error) {
	lockFundHeightRecord := LockFundHeightRecord{}

	records := make(map[string]*LockFundRecord)
	for _, record := range ref.LockFundRecordGroup.Records {

		detail, indexKey, err := record.toStruct(networkType)
		if err != nil {
			return nil, err
		}
		records[*indexKey] = detail
		if err != nil {
			return nil, err
		}
	}
	lockFundHeightRecord.Identifier = Height(ref.LockFundRecordGroup.Identifier.toUint64())
	lockFundHeightRecord.Records = records

	return &lockFundHeightRecord, nil
}

func (ref *LockFundKeyRecordGroupDto) toStruct(networkType NetworkType) (*LockFundKeyRecord, error) {
	lockFundKeyRecord := LockFundKeyRecord{}
	key, err := NewAccountFromPublicKey(ref.LockFundRecordGroup.Identifier, networkType)
	if err != nil {
		return nil, err
	}
	records := make(map[Height]*LockFundRecord)
	for _, record := range ref.LockFundRecordGroup.Records {
		detail, indexKey, err := record.toStruct(networkType)
		records[indexKey] = detail
		if err != nil {
			return nil, err
		}
	}
	lockFundKeyRecord.Identifier = key
	lockFundKeyRecord.Records = records

	return &lockFundKeyRecord, nil
}

type LockFundKeyRecordGroupDtos []*LockFundKeyRecordGroupDto

func (ref *LockFundKeyRecordGroupDtos) toStruct(networkType NetworkType) ([]*LockFundKeyRecord, error) {
	var (
		dtos    = *ref
		records = make([]*LockFundKeyRecord, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		records = append(records, info)
	}

	return records, nil
}

type LockFundHeightRecordGroupDtos []*LockFundHeightRecordGroupDto

func (ref *LockFundHeightRecordGroupDtos) toStruct(networkType NetworkType) ([]*LockFundHeightRecord, error) {
	var (
		dtos    = *ref
		records = make([]*LockFundHeightRecord, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		records = append(records, info)
	}

	return records, nil
}
