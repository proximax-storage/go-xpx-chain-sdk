// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
)

type activeDataModificationDTO struct {
	Id                 hashDto   `json:"id"`
	Owner              string    `json:"owner"`
	DownloadDataCdi    hashDto   `json:"downloadDataCdi"`
	ExpectedUploadSize uint64DTO `json:"expectedUploadSize"`
	ActualUploadSize   uint64DTO `json:"actualUploadSize"`
	FolderName         string    `json:"folderName"`
	ReadyForApproval   bool      `json:"readyForApproval"`
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
		Id:                 id,
		Owner:              owner,
		DownloadDataCdi:    downloadDataCdi,
		ExpectedUploadSize: ref.ExpectedUploadSize.toStruct(),
		ActualUploadSize:   ref.ActualUploadSize.toStruct(),
		FolderName:         ref.FolderName,
		ReadyForApproval:   ref.ReadyForApproval,
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
	ActiveDataModifications activeDataModificationsDTOs `json:"activeDataModifications"`
	State                   DataModificationState       `json:"state"`
}

func (ref *completedDataModificationDTO) toStruct(networkType NetworkType) (*CompletedDataModification, error) {
	activeDataModifications, err := ref.ActiveDataModifications.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &CompletedDataModification{
		ActiveDataModification: activeDataModifications,
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

type confirmedUsedSizeDTO struct {
	Replicator hashDto   `json:"replicator"`
	Size       uint64DTO `json:"size"`
}

func (ref *confirmedUsedSizeDTO) toStruct(networkType NetworkType) (*ConfirmedUsedSize, error) {
	replicator, err := ref.Replicator.Hash()
	if err != nil {
		return nil, err
	}

	return &ConfirmedUsedSize{
		Replicator: replicator,
		Size:       ref.Size.toStruct(),
	}, nil
}

type confirmedUsedSizesDTOs []*confirmedUsedSizeDTO

func (ref *confirmedUsedSizesDTOs) toStruct(networkType NetworkType) ([]*ConfirmedUsedSize, error) {
	var (
		dtos               = *ref
		confirmedUsedSizes = make([]*ConfirmedUsedSize, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		confirmedUsedSizes = append(confirmedUsedSizes, info)
	}

	return confirmedUsedSizes, nil
}

type replicatorsListDTOs []*hashDto

func (ref *replicatorsListDTOs) toStruct() ([]*Hash, error) {
	var (
		dtos        = *ref
		replicators = make([]*Hash, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.Hash()
		if err != nil {
			return nil, err
		}

		replicators = append(replicators, info)
	}

	return replicators, nil
}

type verificationOpinionDTO struct {
	Prover hashDto `json:"prover"`
	Result uint8   `json:"result"`
}

func (ref *verificationOpinionDTO) toStruct(networkType NetworkType) (*VerificationOpinion, error) {
	prover, err := ref.Prover.Hash()
	if err != nil {
		return nil, err
	}

	return &VerificationOpinion{
		Prover: prover,
		Result: ref.Result,
	}, nil
}

type verificationOpinionsDTOs []*verificationOpinionDTO

func (ref *verificationOpinionsDTOs) toStruct(networkType NetworkType) ([]*VerificationOpinion, error) {
	var (
		dtos                 = *ref
		verificationOpinions = make([]*VerificationOpinion, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		verificationOpinions = append(verificationOpinions, info)
	}

	return verificationOpinions, nil
}

type verificationDTO struct {
	VerificationTrigger  hashDto                  `json:"verificationTrigger"`
	State                VerificationState        `json:"state"`
	VerificationOpinions verificationOpinionsDTOs `json:"verificationOpinions"`
}

func (ref *verificationDTO) toStruct(networkType NetworkType) (*Verification, error) {
	verificationTrigger, err := ref.VerificationTrigger.Hash()
	if err != nil {
		return nil, err
	}

	verificationOpinions, err := ref.VerificationOpinions.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.verificationDTO.toStruct VerificationOpinions.toStruct: %v", err)
	}

	return &Verification{
		VerificationTrigger:  verificationTrigger,
		State:                ref.State,
		VerificationOpinions: verificationOpinions,
	}, nil
}

type verificationsDTOs []*verificationDTO

func (ref *verificationsDTOs) toStruct(networkType NetworkType) ([]*Verification, error) {
	var (
		dtos          = *ref
		verifications = make([]*Verification, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		verifications = append(verifications, info)
	}

	return verifications, nil
}

type bcDriveDTO struct {
	BcDrive struct {
		DriveKey                   string                         `json:"multisig"`
		Owner                      string                         `json:"owner"`
		RootHash                   hashDto                        `json:"rootHash"`
		DriveSize                  uint64DTO                      `json:"size"`
		UsedSize                   uint64DTO                      `json:"usedSize"`
		MetaFilesSize              uint64DTO                      `json:"metaFilesSize"`
		ReplicatorCount            uint16                         `json:"replicatorCount"`
		OwnerCumulativeUploadSize  uint64DTO                      `json:"ownerCumulativeUploadSize"`
		ActiveDataModifications    activeDataModificationsDTOs    `json:"activeDataModifications"`
		CompletedDataModifications completedDataModificationsDTOs `json:"completedDataModifications"`
		ConfirmedUsedSizes         confirmedUsedSizesDTOs         `json:"confirmedUsedSizes"`
		Replicators                replicatorsListDTOs            `json:"replicators"`
		Verifications              verificationsDTOs              `json:"verifications"`
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
		return nil, err
	}

	bcDrive.BcDriveAccount = bcDriveAccount
	bcDrive.OwnerAccount = ownerAccount
	bcDrive.RootHash = rootHash
	bcDrive.DriveSize = ref.BcDrive.DriveSize.toStruct()
	bcDrive.UsedSize = ref.BcDrive.UsedSize.toStruct()
	bcDrive.MetaFilesSize = ref.BcDrive.MetaFilesSize.toStruct()
	bcDrive.ReplicatorCount = ref.BcDrive.ReplicatorCount
	bcDrive.OwnerCumulativeUploadSize = ref.BcDrive.OwnerCumulativeUploadSize.toStruct()

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

	confirmedUsedSizes, err := ref.BcDrive.ConfirmedUsedSizes.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.ConfirmedUsedSizes.toStruct: %v", err)
	}

	bcDrive.ConfirmedUsedSizes = confirmedUsedSizes

	replicators, err := ref.BcDrive.Replicators.toStruct()
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.Replicators.toStruct: %v", err)
	}

	bcDrive.Replicators = replicators

	verifications, err := ref.BcDrive.Verifications.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.Verifications.toStruct: %v", err)
	}

	bcDrive.Verifications = verifications

	return &bcDrive, nil
}

type bcDriveDTOs []*bcDriveDTO

func (ref *bcDriveDTOs) toStruct(networkType NetworkType) ([]*BcDrive, error) {
	var (
		dtos     = *ref
		bcDrives = make([]*BcDrive, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		bcDrives = append(bcDrives, info)
	}

	return bcDrives, nil
}

type driveV2DTO struct {
	Drive                          string    `json:"drive"`
	LastApprovedDataModificationId *hashDto  `json:"lastApprovedDataModificationId"`
	DataModificationIdIsValid      bool      `json:"dataModificationIdIsValid"`
	InitialDownloadWork            uint64DTO `json:"initialDownloadWork"`
}

type driveV2DTOs []*driveV2DTO

func (ref *driveV2DTOs) toStruct(networkType NetworkType) (map[string]*DriveInfo, error) {
	var (
		dtos      = *ref
		driveInfo = make(map[string]*DriveInfo)
	)

	for i, dto := range dtos {
		drive, err := NewAccountFromPublicKey(dto.Drive, networkType)
		if err != nil {
			return nil, err
		}

		lastApprovedDataModificationId, err := dto.LastApprovedDataModificationId.Hash()
		if err != nil {
			return nil, err
		}

		info := DriveInfo{
			LastApprovedDataModificationId: lastApprovedDataModificationId,
			DataModificationIdIsValid:      dto.DataModificationIdIsValid,
			InitialDownloadWork:            dto.InitialDownloadWork.toStruct(),
			Index:                          i,
		}

		driveInfo[drive.PublicKey] = &info
	}

	return driveInfo, nil
}

type replicatorV2DTO struct {
	Replicator struct {
		ReplicatorKey string      `json:"key"`
		Version       uint32      `json:"version"`
		Capacity      uint64DTO   `json:"capacity"`
		Drives        driveV2DTOs `json:"drives"`
	}
}

func (ref *replicatorV2DTO) toStruct(networkType NetworkType) (*Replicator, error) {
	replicator := Replicator{}

	replicatorAccount, err := NewAccountFromPublicKey(ref.Replicator.ReplicatorKey, networkType)
	if err != nil {
		return nil, err
	}

	replicator.ReplicatorAccount = replicatorAccount
	replicator.Version = ref.Replicator.Version
	replicator.Capacity = ref.Replicator.Capacity.toStruct()

	drives, err := ref.Replicator.Drives.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.replicatorV2DTO.toStruct Replicator.Drives.toStruct: %v", err)
	}

	replicator.Drives = drives

	return &replicator, nil
}

type replicatorV2DTOs []*replicatorV2DTO

func (ref *replicatorV2DTOs) toStruct(networkType NetworkType) ([]*Replicator, error) {
	var (
		dtos        = *ref
		replicators = make([]*Replicator, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		replicators = append(replicators, info)
	}

	return replicators, nil
}

type bcDrivesPageDTO struct {
	BcDrives []bcDriveDTO `json:"data"`

	Pagination struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (t *bcDrivesPageDTO) toStruct(networkType NetworkType) (*BcDrivesPage, error) {
	page := &BcDrivesPage{
		BcDrives: make([]*BcDrive, len(t.BcDrives)),
		Pagination: Pagination{
			TotalEntries: t.Pagination.TotalEntries,
			PageNumber:   t.Pagination.PageNumber,
			PageSize:     t.Pagination.PageSize,
			TotalPages:   t.Pagination.TotalPages,
		},
	}

	errs := make([]error, len(t.BcDrives))
	for i, t := range t.BcDrives {
		currDr, currErr := t.toStruct(networkType)
		page.BcDrives[i], errs[i] = currDr, currErr
	}

	for _, err := range errs {
		if err != nil {
			return page, err
		}
	}

	return page, nil
}

type replicatorsPageDTO struct {
	Replicators []replicatorV2DTO `json:"data"`

	Pagination struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (t *replicatorsPageDTO) toStruct(networkType NetworkType) (*ReplicatorsPage, error) {
	page := &ReplicatorsPage{
		Replicators: make([]*Replicator, len(t.Replicators)),
		Pagination: Pagination{
			TotalEntries: t.Pagination.TotalEntries,
			PageNumber:   t.Pagination.PageNumber,
			PageSize:     t.Pagination.PageSize,
			TotalPages:   t.Pagination.TotalPages,
		},
	}

	for i, t := range t.Replicators {
		currDr, err := t.toStruct(networkType)
		page.Replicators[i] = currDr
		if err != nil {
			return page, err
		}
	}

	return page, nil
}
