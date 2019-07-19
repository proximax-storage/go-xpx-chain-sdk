// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"strconv"
)

type networkDTO struct {
	Name        string
	Description string
}

type entityDTO struct {
	Name              string          `json:"name"`
	Type              string          `json:"type"`
	SupportedVersions []EntityVersion `json:"supportedVersions"`
}

func (dto *entityDTO) toStruct() (*Entity, error) {
	entityType, err := strconv.ParseUint(dto.Type, 10, 16)
	if err != nil {
		return nil, err
	}

	return &Entity{
		Name:              dto.Name,
		Type:              EntityType(entityType),
		SupportedVersions: dto.SupportedVersions,
	}, nil
}

type supportedEntitiesDTO struct {
	Entities []*entityDTO `json:"entities"`
}

func (dto *supportedEntitiesDTO) toStruct(ref *SupportedEntities) error {
	for _, dto := range dto.Entities {
		entity, err := dto.toStruct()
		if err != nil {
			return err
		}

		ref.Entities[entity.Type] = entity
	}

	return nil
}

type networkConfigDTO struct {
	DTO struct {
		StartedHeight           uint64DTO `json:"height"`
		BlockChainConfig        string    `json:"blockChainConfig"`
		SupportedEntityVersions string    `json:"supportedEntityVersions"`
	} `json:"catapultConfig"`
}

func (dto *networkConfigDTO) toStruct() (*NetworkConfig, error) {
	s := NewSupportedEntities()

	err := s.UnmarshalBinary([]byte(dto.DTO.SupportedEntityVersions))
	if err != nil {
		return nil, err
	}

	c := NewBlockChainConfig()

	err = c.UnmarshalBinary([]byte(dto.DTO.BlockChainConfig))
	if err != nil {
		return nil, err
	}

	return &NetworkConfig{
		Height(dto.DTO.StartedHeight.toUint64()),
		c,
		s,
	}, nil
}

type networkVersionDTO struct {
	DTO struct {
		StartedHeight   uint64DTO `json:"height"`
		CatapultVersion uint64DTO `json:"catapultVersion"`
	} `json:"catapultUpgrade"`
}

func (dto *networkVersionDTO) toStruct() *NetworkVersion {
	return &NetworkVersion{
		Height(dto.DTO.StartedHeight.toUint64()),
		CatapultVersion(dto.DTO.CatapultVersion.toUint64()),
	}
}
