// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
)

type activeDataModificationDTO struct {
	Id              hashDto   `json:"id"`
	Owner           string    `json:"owner"`
	DownloadDataCdi hashDto   `json:"downloadDataCdi"`
	UploadSize      uint64DTO `json:"uploadSize"`
}

func (ref *activeDataModificationDTO) toStruct(networkType NetworkType) (*ActiveDataModification, error) {
	id, err := ref.Id.Hash()
	if err != nil {
		return nil, err
	}

	owner, err := NewAccountFromPublicKey(ref.Owner, networkType)
	if err != nil {
		return nil, err
	}

	downloadDataCdi, err := ref.DownloadDataCdi.Hash()
	if err != nil {
		return nil, err
	}

	return &ActiveDataModification{
		Id:              id,
		Owner:           owner,
		DownloadDataCdi: downloadDataCdi,
		UploadSize:      ref.UploadSize.toStruct(),
	}, nil
}

type activeDataModificationsDTOs []*activeDataModificationDTO

func (ref *activeDataModificationsDTOs) toStruct(networkType NetworkType) ([]*ActiveDataModification, error) {
	var (
		dtos                    = *ref
		activeDataModifications = make([]*ActiveDataModification, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		activeDataModifications = append(activeDataModifications, info)
	}

	return activeDataModifications, nil
}

type completedDataModificationDTO struct {
	ActiveDataModification *activeDataModificationDTO `json:"activeDataModification"`
	State                  DataModificationState      `json:"state"`
}

func (ref *completedDataModificationDTO) toStruct(networkType NetworkType) (*CompletedDataModification, error) {
	activeDataModification, err := ref.ActiveDataModification.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &CompletedDataModification{
		ActiveDataModification: activeDataModification,
		State:                  ref.State,
	}, nil
}

type completedDataModificationsDTOs []*completedDataModificationDTO

func (ref *completedDataModificationsDTOs) toStruct(networkType NetworkType) ([]*CompletedDataModification, error) {
	var (
		dtos                       = *ref
		completedDataModifications = make([]*CompletedDataModification, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		completedDataModifications = append(completedDataModifications, info)
	}

	return completedDataModifications, nil
}

type activeDownloadsDTOs []*hashDto

func (ref *activeDownloadsDTOs) toStruct(networkType NetworkType) ([]*Hash, error) {
	var (
		dtos            = *ref
		activeDownloads = make([]*Hash, 0, len(dtos))
	)

	for _, dto := range dtos {
		activeId, err := dto.Hash()
		if err != nil {
			return nil, err
		}

		activeDownloads = append(activeDownloads, activeId)
	}

	return activeDownloads, nil
}

type completedDownloadsDTOs []*hashDto

func (ref *completedDownloadsDTOs) toStruct(networkType NetworkType) ([]*Hash, error) {
	var (
		dtos               = *ref
		completedDownloads = make([]*Hash, 0, len(dtos))
	)

	for _, dto := range dtos {
		completedId, err := dto.Hash()
		if err != nil {
			return nil, err
		}

		completedDownloads = append(completedDownloads, completedId)
	}

	return completedDownloads, nil
}

type replicatorDTOs []*PublicAccount

func (ref *replicatorDTOs) toStruct(networkType NetworkType) ([]*PublicAccount, error) {
	var (
		dtos        = *ref
		replicators = make([]*PublicAccount, 0, len(dtos))
	)

	for i, dto := range dtos {
		replicatorKey, err := NewAccountFromPublicKey(dto.PublicKey, networkType)
		if err != nil {
			return nil, err
		}

		replicators[i] = replicatorKey
	}

	return replicators, nil
}

type bcDriveDTO struct {
	BcDrive struct {
		DriveKey                   string                         `json:"driveKey"`
		Owner                      string                         `json:"owner"`
		RootHash                   hashDto                        `json:"rootHash"`
		DriveSize                  uint64                         `json:"driveSize"`
		ReplicatorCount            uint16                         `json:"replicatorCount"`
		ActiveDataModifications    activeDataModificationsDTOs    `json:"activeDataModifications"`
		CompletedDataModifications completedDataModificationsDTOs `json:"completedDataModifications"`
		ActiveDownloads            activeDownloadsDTOs            `json:"activeDownloads"`
		CompletedDownloads         completedDownloadsDTOs         `json:"completedDownloads"`
		Replicators                replicatorDTOs                 `json:"replicators"`
	}
}

func (ref *bcDriveDTO) toStruct(networkType NetworkType) (*BcDrive, error) {
	bcDrive := BcDrive{}

	bcDriveAccount, err := NewAccountFromPublicKey(ref.BcDrive.DriveKey, networkType)
	if err != nil {
		return nil, err
	}

	ownerAccount, err := NewAccountFromPublicKey(ref.BcDrive.Owner, networkType)
	if err != nil {
		return nil, err
	}

	rootHash, err := ref.BcDrive.RootHash.Hash()
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.RootHash.Hash: %v", err)
	}

	bcDrive.BcDriveAccount = bcDriveAccount
	bcDrive.OwnerAccount = ownerAccount
	bcDrive.RootHash = rootHash
	bcDrive.DriveSize = ref.BcDrive.DriveSize
	bcDrive.ReplicatorCount = ref.BcDrive.ReplicatorCount

	activeDataModifications, err := ref.BcDrive.ActiveDataModifications.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.ActiveDataModifications.toStruct: %v", err)
	}

	bcDrive.ActiveDataModifications = activeDataModifications

	completedDataModifications, err := ref.BcDrive.CompletedDataModifications.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.CompletedDataModifications.toStruct: %v", err)
	}

	bcDrive.CompletedDataModifications = completedDataModifications

	activeDownloads, err := ref.BcDrive.ActiveDownloads.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.ActiveDownloads.toStruct: %v", err)
	}

	bcDrive.ActiveDownloads = activeDownloads

	completedDownloads, err := ref.BcDrive.CompletedDownloads.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.CompletedDownloads.toStruct: %v", err)
	}

	bcDrive.CompletedDownloads = completedDownloads

	replicators, err := ref.BcDrive.Replicators.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.Replicators.toStruct: %v", err)
	}

	bcDrive.Replicators = replicators

	return &bcDrive, nil
}
