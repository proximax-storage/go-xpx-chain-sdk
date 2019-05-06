// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/sha3"
	"math/big"
)

func bigIntToMosaicId(bigInt *big.Int) *MosaicId {
	if bigInt == nil {
		return nil
	}

	mscId := MosaicId(*bigInt)

	return &mscId
}

func mosaicIdToBigInt(mscId *MosaicId) *big.Int {
	if mscId == nil {
		return nil
	}

	return (*big.Int)(mscId)
}

func generateMosaicId(nonce uint32, ownerPublicKey string) (*big.Int, error) {
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

	return uint64DTO{binary.LittleEndian.Uint32(t[0:4]), binary.LittleEndian.Uint32(t[4:8]) & 0x7FFFFFFF}.toBigInt(), nil
}

type mosaicDTO struct {
	MosaicId uint64DTO `json:"id"`
	Amount   uint64DTO `json:"amount"`
}

func (dto *mosaicDTO) toStruct() (*Mosaic, error) {
	mosaicId, err := NewMosaicId(dto.MosaicId.toBigInt())
	if err != nil {
		return nil, err
	}

	return &Mosaic{mosaicId, dto.Amount.toBigInt()}, nil
}

type mosaicPropertiesDTO []uint64DTO

// namespaceMosaicMetaDTO
type namespaceMosaicMetaDTO struct {
	Active bool
	Index  int
	Id     string
}

type mosaicDefinitionDTO struct {
	MosaicId   uint64DTO
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

func (dto *mosaicPropertiesDTO) toStruct() *MosaicProperties {
	flags := "00" + (*dto)[0].toBigInt().Text(2)
	bitMapFlags := flags[len(flags)-3:]

	return NewMosaicProperties(bitMapFlags[2] == '1',
		bitMapFlags[1] == '1',
		bitMapFlags[0] == '1',
		byte((*dto)[1].toBigInt().Int64()),
		(*dto)[2].toBigInt(),
	)
}

func (ref *mosaicInfoDTO) toStruct(networkType NetworkType) (*MosaicInfo, error) {
	publicAcc, err := NewAccountFromPublicKey(ref.Mosaic.Owner, networkType)
	if err != nil {
		return nil, err
	}

	if len(ref.Mosaic.Properties) < 3 {
		return nil, errors.New("mosaic Properties is not valid")
	}

	mosaicId, err := NewMosaicId(ref.Mosaic.MosaicId.toBigInt())

	mscInfo := &MosaicInfo{
		MosaicId:   mosaicId,
		Supply:     ref.Mosaic.Supply.toBigInt(),
		Height:     ref.Mosaic.Height.toBigInt(),
		Owner:      publicAcc,
		Revision:   ref.Mosaic.Revision,
		Properties: ref.Mosaic.Properties.toStruct(),
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
	MosaicId uint64DTO `json:"mosaicId"`
	Names    []string  `json:"names"`
}

func (m *mosaicNameDTO) toStruct() (*MosaicName, error) {
	mosaicId, err := NewMosaicId(m.MosaicId.toBigInt())
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
