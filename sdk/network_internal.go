// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type networkDTO struct {
	Name        string
	Description string
}

type networkConfigDTO struct {
	DTO struct {
		StartedHeight           uint64DTO `json:"height"`
		BlockChainConfig        string    `json:"blockChainConfig"`
		SupportedEntityVersions string    `json:"supportedEntityVersions"`
	} `json:"catapultConfig"`
}

func (dto *networkConfigDTO) toStruct() *NetworkConfig {
	return &NetworkConfig{
		Height(dto.DTO.StartedHeight.toUint64()),
		dto.DTO.BlockChainConfig,
		dto.DTO.SupportedEntityVersions,
	}
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
