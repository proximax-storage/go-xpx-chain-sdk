// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/sha3"
)

type mosaicIdDTO uint64DTO

func (dto *mosaicIdDTO) toStruct() (*MosaicId, error) {
	return NewMosaicId(uint64DTO(*dto).toUint64())
}

type mosaicIdDTOs []*mosaicIdDTO

func (dto *mosaicIdDTOs) toStruct() ([]*MosaicId, error) {
	ids := make([]*MosaicId, len(*dto))
	var err error

	for i, m := range *dto {
		ids[i], err = m.toStruct()
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}

func generateMosaicId(nonce uint32, ownerPublicKey string) (*MosaicId, error) {
	result := sha3.New256()
	nonceB := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonceB, nonce)

	if _, err := result.Write(nonceB); err != nil {
		return nil, err
	}

	ownerBytes, err := hex.DecodeString(ownerPublicKey)

	if err != nil {
		return nil, err
	}

	if _, err := result.Write(ownerBytes); err != nil {
		return nil, err
	}

	t := result.Sum(nil)
	return NewMosaicId(binary.LittleEndian.Uint64(t) & (^NamespaceBit))
}

type mosaicDTO struct {
	AssetId assetIdDTO `json:"id"`
	Amount  uint64DTO  `json:"amount"`
}

func (dto *mosaicDTO) toStruct() (*Mosaic, error) {
	assetId, err := dto.AssetId.toStruct()
	if err != nil {
		return nil, err
	}

	return NewMosaic(assetId, dto.Amount.toStruct())
}

type MosaicPropertyId uint8

// MosaicPropertyId enums
const (
	MosaicPropertyFlagsId MosaicPropertyId = iota
	MosaicPropertyDivisibilityId
	MosaicPropertyDurationId
)

type mosaicPropertyDTO struct {
	Id    MosaicPropertyId `json:"id"`
	Value uint64DTO        `json:"value"`
}

type mosaicPropertiesDTO []*mosaicPropertyDTO

// namespaceMosaicMetaDTO
type namespaceMosaicMetaDTO struct {
	Active bool
	Index  int
	Id     string
}

type mosaicDefinitionDTO struct {
	MosaicId   mosaicIdDTO
	Supply     uint64DTO
	Height     uint64DTO
	Owner      string
	Revision   uint32
	Properties mosaicPropertiesDTO
	Levy       interface{}
}

// mosaicInfoDTO is temporary struct for reading response & fill MosaicInfo
type mosaicInfoDTO struct {
	Mosaic mosaicDefinitionDTO
}

func (dto *mosaicPropertiesDTO) toStruct() (*MosaicProperties, error) {
	flags := uint64(0)
	divisibility := uint8(0)
	duration := Duration(0)
	for _, property := range *dto {
		switch property.Id {
		case MosaicPropertyFlagsId:
			flags = property.Value.toUint64()
		case MosaicPropertyDivisibilityId:
			divisibility = byte(property.Value.toUint64())
		case MosaicPropertyDurationId:
			duration = Duration(property.Value.toUint64())
		default:
			return nil, errors.New("Unknown Property Id")
		}
	}

	return NewMosaicProperties(
		hasBits(flags, Supply_Mutable),
		hasBits(flags, Transferable),
		divisibility,
		duration,
	), nil
}

func (ref *mosaicInfoDTO) toStruct(networkType NetworkType) (*MosaicInfo, error) {
	publicAcc, err := NewAccountFromPublicKey(ref.Mosaic.Owner, networkType)
	if err != nil {
		return nil, err
	}

	if len(ref.Mosaic.Properties) < 3 {
		return nil, errors.New("mosaic Properties is not valid")
	}

	mosaicId, err := ref.Mosaic.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	properties, err := ref.Mosaic.Properties.toStruct()
	if err != nil {
		return nil, err
	}

	mscInfo := &MosaicInfo{
		MosaicId:   mosaicId,
		Supply:     ref.Mosaic.Supply.toStruct(),
		Height:     ref.Mosaic.Height.toStruct(),
		Owner:      publicAcc,
		Revision:   ref.Mosaic.Revision,
		Properties: properties,
	}

	return mscInfo, nil
}

type mosaicInfoDTOs []*mosaicInfoDTO

func (m *mosaicInfoDTOs) toStruct(networkType NetworkType) ([]*MosaicInfo, error) {
	dtos := *m

	mscInfos := make([]*MosaicInfo, 0, len(dtos))

	for _, dto := range dtos {
		mscInfo, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		mscInfos = append(mscInfos, mscInfo)
	}

	return mscInfos, nil
}

type mosaicIds struct {
	MosaicIds []*MosaicId `json:"mosaicIds"`
}

func (ref *mosaicIds) MarshalJSON() ([]byte, error) {
	buf := []byte(`{"mosaicIds": [`)

	for i, nsId := range ref.MosaicIds {
		if i > 0 {
			buf = append(buf, ',')
		}

		buf = append(buf, []byte(`"`+nsId.toHexString()+`"`)...)
	}

	buf = append(buf, ']', '}')

	return buf, nil
}

type mosaicNameDTO struct {
	MosaicId mosaicIdDTO `json:"mosaicId"`
	Names    []string    `json:"names"`
}

func (m *mosaicNameDTO) toStruct() (*MosaicName, error) {
	mosaicId, err := m.MosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	return &MosaicName{
		MosaicId: mosaicId,
		Names:    m.Names,
	}, nil
}

type mosaicNameDTOs []*mosaicNameDTO

func (m *mosaicNameDTOs) toStruct() ([]*MosaicName, error) {
	dtos := *m
	mscNames := make([]*MosaicName, 0, len(dtos))

	for _, dto := range dtos {
		mscName, err := dto.toStruct()
		if err != nil {
			return nil, err
		}

		mscNames = append(mscNames, mscName)
	}

	return mscNames, nil
}

type mosaicLevyDTO struct {
	Type      uint8        `json:"type"`
	Recipient string       `json:"recipient"`
	Fee       *uint64DTO   `json:"fee"`
	MosaicId  *mosaicIdDTO `json:"mosaicId"`
}

func (dto *mosaicLevyDTO) toStruct() (MosaicLevy, error) {

	mosaicId, err := dto.MosaicId.toStruct()
	if err != nil {
		return MosaicLevy{}, err
	}

	var f Amount = Amount(0)
	if dto.Fee != nil {
		f = dto.Fee.toStruct()
	}

	a, err := NewAddressFromBase32(dto.Recipient)
	if err != nil {
		a = NewAddress("", NotSupportedNet)
	}

	return MosaicLevy{
		Type:      dto.Type,
		Recipient: a,
		MosaicId:  mosaicId,
		Fee:       f,
	}, nil
}
