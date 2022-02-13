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

type activeDataModificationDTOs []*activeDataModificationDTO

func (ref *activeDataModificationDTOs) toStruct(networkType NetworkType) ([]*ActiveDataModification, error) {
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
	activeDataModificationDTO
	State DataModificationState `json:"state"`
}

func (ref *completedDataModificationDTO) toStruct(networkType NetworkType) (*CompletedDataModification, error) {
	activeDataModifications, err := ref.activeDataModificationDTO.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &CompletedDataModification{
		activeDataModifications,
		ref.State,
	}, nil
}

type completedDataModificationDTOs []*completedDataModificationDTO

func (ref *completedDataModificationDTOs) toStruct(networkType NetworkType) ([]*CompletedDataModification, error) {
	var (
		dtos                       = *ref
		completedDataModifications = make([]*CompletedDataModification, 0, len(dtos))
	)

	for _, dto := range dtos {
		completed, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		completedDataModifications = append(completedDataModifications, completed)
	}

	return completedDataModifications, nil
}

type confirmedUsedSizeDTO struct {
	Replicator string    `json:"replicator"`
	Size       uint64DTO `json:"size"`
}

func (ref *confirmedUsedSizeDTO) toStruct(networkType NetworkType) (*ConfirmedUsedSize, error) {
	replicatorAccount, err := NewAccountFromPublicKey(ref.Replicator, networkType)
	if err != nil {
		return nil, err
	}

	return &ConfirmedUsedSize{
		Replicator: replicatorAccount,
		Size:       ref.Size.toStruct(),
	}, nil
}

type confirmedUsedSizeDTOs []*confirmedUsedSizeDTO

func (ref *confirmedUsedSizeDTOs) toStruct(networkType NetworkType) ([]*ConfirmedUsedSize, error) {
	var (
		dtos               = *ref
		confirmedUsedSizes = make([]*ConfirmedUsedSize, 0, len(dtos))
	)

	for _, dto := range dtos {
		confirmed, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		confirmedUsedSizes = append(confirmedUsedSizes, confirmed)
	}

	return confirmedUsedSizes, nil
}

type replicatorsListDTOs []string

func (ref *replicatorsListDTOs) toStruct(networkType NetworkType) ([]*PublicAccount, error) {
	var (
		dtos        = *ref
		replicators = make([]*PublicAccount, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := NewAccountFromPublicKey(dto, networkType)
		if err != nil {
			return nil, err
		}

		replicators = append(replicators, info)
	}

	return replicators, nil
}

type shardDTO struct {
	DownloadDataCdi hashDto
	Replicators     replicatorsListDTOs
}

func (ref *shardDTO) toStruct(networkType NetworkType) (*Shard, error) {
	downloadDataCdi, err := ref.DownloadDataCdi.Hash()
	if err != nil {
		return nil, err
	}

	replicators, err := ref.Replicators.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &Shard{
		DownloadChannelId: downloadDataCdi,
		Replicators:       replicators,
	}, nil
}

type shardDTOs []*shardDTO

func (ref *shardDTOs) toStruct(networkType NetworkType) ([]*Shard, error) {
	var (
		dtos   = *ref
		shards = make([]*Shard, 0, len(dtos))
	)

	for _, dto := range dtos {
		shard, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		shards = append(shards, shard)
	}

	return shards, nil
}

type verificationDTO struct {
	VerificationTrigger *Hash                  `json:"verification_trigger"`
	Expiration          blockchainTimestampDTO `json:"expiration"`
	Expired             bool                   `json:"expired"`
	Shards              shardDTOs              `json:"shards"`
}

func (ref *verificationDTO) toStruct(networkType NetworkType) (*Verification, error) {
	shards, err := ref.Shards.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &Verification{
		VerificationTrigger: ref.VerificationTrigger,
		Expiration:          ref.Expiration.toStruct().ToTimestamp(),
		Expired:             ref.Expired,
		Shards:              shards,
	}, nil
}

type verificationDTOs []*verificationDTO

func (ref *verificationDTOs) toStruct(networkType NetworkType) ([]*Verification, error) {
	var (
		dtos          = *ref
		verifications = make([]*Verification, 0, len(dtos))
	)

	for _, dto := range dtos {
		verification, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		verifications = append(verifications, verification)
	}

	return verifications, nil
}

type bcDriveDTO struct {
	Drive struct {
		DriveKey                   string                        `json:"multisig"`
		Owner                      string                        `json:"owner"`
		RootHash                   hashDto                       `json:"rootHash"`
		DriveSize                  uint64DTO                     `json:"size"`
		UsedSize                   uint64DTO                     `json:"usedSize"`
		MetaFilesSize              uint64DTO                     `json:"metaFilesSize"`
		ReplicatorCount            uint16                        `json:"replicatorCount"`
		OwnerCumulativeUploadSize  uint64DTO                     `json:"ownerCumulativeUploadSize"`
		ActiveDataModifications    activeDataModificationDTOs    `json:"activeDataModifications"`
		CompletedDataModifications completedDataModificationDTOs `json:"completedDataModifications"`
		ConfirmedUsedSizes         confirmedUsedSizeDTOs         `json:"confirmedUsedSizes"`
		Replicators                replicatorsListDTOs           `json:"replicators"`
		OffboardingReplicators     replicatorsListDTOs           `json:"offboardingReplicators"`
		Verifications              verificationDTOs              `json:"verifications"`
	}
}

func (ref *bcDriveDTO) toStruct(networkType NetworkType) (*BcDrive, error) {
	bcDrive := BcDrive{}

	bcDriveAccount, err := NewAccountFromPublicKey(ref.Drive.DriveKey, networkType)
	if err != nil {
		return nil, err
	}

	ownerAccount, err := NewAccountFromPublicKey(ref.Drive.Owner, networkType)
	if err != nil {
		return nil, err
	}

	rootHash, err := ref.Drive.RootHash.Hash()
	if err != nil {
		return nil, err
	}

	bcDrive.MultisigAccount = bcDriveAccount
	bcDrive.Owner = ownerAccount
	bcDrive.RootHash = rootHash
	bcDrive.Size = ref.Drive.DriveSize.toStruct()
	bcDrive.UsedSize = ref.Drive.UsedSize.toStruct()
	bcDrive.MetaFilesSize = ref.Drive.MetaFilesSize.toStruct()
	bcDrive.ReplicatorCount = ref.Drive.ReplicatorCount
	bcDrive.OwnerCumulativeUploadSize = ref.Drive.OwnerCumulativeUploadSize.toStruct()

	activeDataModifications, err := ref.Drive.ActiveDataModifications.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.ActiveDataModifications.toStruct: %v", err)
	}

	bcDrive.ActiveDataModifications = activeDataModifications

	completedDataModifications, err := ref.Drive.CompletedDataModifications.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.CompletedDataModifications.toStruct: %v", err)
	}

	bcDrive.CompletedDataModifications = completedDataModifications

	confirmedUsedSizes, err := ref.Drive.ConfirmedUsedSizes.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.ConfirmedUsedSizes.toStruct: %v", err)
	}

	bcDrive.ConfirmedUsedSizes = confirmedUsedSizes

	replicators, err := ref.Drive.Replicators.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.Replicators.toStruct: %v", err)
	}

	bcDrive.Replicators = replicators

	offboardingReplicators, err := ref.Drive.OffboardingReplicators.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.OffboardingReplicators.toStruct: %v", err)
	}

	bcDrive.OffboardingReplicators = offboardingReplicators

	verifications, err := ref.Drive.Verifications.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.bcDriveDTO.toStruct BcDrive.Verifications.toStruct: %v", err)
	}

	bcDrive.Verifications = verifications

	return &bcDrive, nil
}

type driveV2DTO struct {
	Drive                               string    `json:"drive"`
	LastApprovedDataModificationId      *hashDto  `json:"lastApprovedDataModificationId"`
	DataModificationIdIsValid           bool      `json:"dataModificationIdIsValid"`
	InitialDownloadWork                 uint64DTO `json:"initialDownloadWork"`
	LastCompletedCumulativeDownloadWork uint64DTO `json:"lastCompletedCumulativeDownloadWork"`
}

func (ref *driveV2DTO) toStruct(networkType NetworkType) (*DriveInfo, error) {
	driveKey, err := NewAccountFromPublicKey(ref.Drive, networkType)
	if err != nil {
		return nil, err
	}

	lastApprovedDataModificationId, err := ref.LastApprovedDataModificationId.Hash()
	if err != nil {
		return nil, err
	}

	return &DriveInfo{
		DriveKey:                            driveKey,
		LastApprovedDataModificationId:      lastApprovedDataModificationId,
		DataModificationIdIsValid:           ref.DataModificationIdIsValid,
		InitialDownloadWork:                 ref.InitialDownloadWork.toStruct(),
		LastCompletedCumulativeDownloadWork: ref.LastCompletedCumulativeDownloadWork.toStruct(),
	}, nil
}

type driveV2DTOs []*driveV2DTO

func (ref *driveV2DTOs) toStruct(networkType NetworkType) ([]*DriveInfo, error) {
	var (
		dtos      = *ref
		driveInfo = make([]*DriveInfo, 0, len(dtos))
	)

	for _, dto := range dtos {
		info, err := dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}

		driveInfo = append(driveInfo, info)
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

type paymentV2DTO struct {
	Replicator string    `json:"replicator"`
	Payment    uint64DTO `json:"payment"`
}

type paymentsV2DTOs []*paymentV2DTO

func (ref *paymentsV2DTOs) toStruct(networkType NetworkType) ([]*Payment, error) {
	var (
		dtos     = *ref
		payments = make([]*Payment, 0, len(dtos))
	)

	for _, dto := range dtos {
		replicator, err := NewAccountFromPublicKey(dto.Replicator, networkType)
		if err != nil {
			return nil, err
		}

		payment := &Payment{
			Replicator: replicator,
			Payment:    dto.Payment.toStruct(),
		}

		payments = append(payments, payment)
	}

	return payments, nil
}

type downloadChannelDTO struct {
	DownloadChannelInfo struct {
		Id                    hashDto             `json:"id"`
		Consumer              string              `json:"consumer"`
		Drive                 string              `json:"drive"`
		DownloadSize          uint64DTO           `json:"downloadSize"`
		DownloadApprovalCount uint16              `json:"downloadApprovalCount"`
		ListOfPublicKeys      replicatorsListDTOs `json:"listOfPublicKeys"`
		CumulativePayments    paymentsV2DTOs      `json:"cumulativePayments"`
	}
}

func (ref *downloadChannelDTO) toStruct(networkType NetworkType) (*DownloadChannel, error) {
	downloadChannel := DownloadChannel{}

	id, err := ref.DownloadChannelInfo.Id.Hash()
	if err != nil {
		return nil, err
	}

	consumer, err := NewAccountFromPublicKey(ref.DownloadChannelInfo.Consumer, networkType)
	if err != nil {
		return nil, err
	}

	drive, err := NewAccountFromPublicKey(ref.DownloadChannelInfo.Drive, networkType)
	if err != nil {
		return nil, err
	}

	downloadChannel.Id = id
	downloadChannel.Consumer = consumer
	downloadChannel.Drive = drive
	downloadChannel.DownloadSize = ref.DownloadChannelInfo.DownloadSize.toStruct()

	listOfPublicKeys, err := ref.DownloadChannelInfo.ListOfPublicKeys.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.downloadChannelDTO.toStruct DownloadChannel.ListOfPublicKeys.toStruct: %v", err)
	}

	downloadChannel.ListOfPublicKeys = listOfPublicKeys

	cumulativePayments, err := ref.DownloadChannelInfo.CumulativePayments.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.downloadChannelDTO.toStruct DownloadChannel.CumulativePayments.toStruct: %v", err)
	}

	downloadChannel.CumulativePayments = cumulativePayments

	return &downloadChannel, nil
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

	var err error
	for i, t := range t.BcDrives {
		page.BcDrives[i], err = t.toStruct(networkType)
		if err != nil {
			return nil, err
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

	var err error
	for i, t := range t.Replicators {
		page.Replicators[i], err = t.toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}

	return page, nil
}

type downloadChannelsPageDTO struct {
	DownloadChannels []downloadChannelDTO `json:"data"`

	Pagination struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (t *downloadChannelsPageDTO) toStruct(networkType NetworkType) (*DownloadChannelsPage, error) {
	page := &DownloadChannelsPage{
		DownloadChannels: make([]*DownloadChannel, len(t.DownloadChannels)),
		Pagination: Pagination{
			TotalEntries: t.Pagination.TotalEntries,
			PageNumber:   t.Pagination.PageNumber,
			PageSize:     t.Pagination.PageSize,
			TotalPages:   t.Pagination.TotalPages,
		},
	}

	var err error
	for i, t := range t.DownloadChannels {
		page.DownloadChannels[i], err = t.toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}

	return page, nil
}

type opinionDTO struct {
	Opinion []uint64DTO
}

func (ref *opinionDTO) toStruct() (*Opinion, error) {
	opinion := make([]OpinionSize, 0)
	for _, dto := range ref.Opinion {
		info := dto.toStruct()
		opinion = append(opinion, info)
	}

	return &Opinion{
		Opinion: opinion,
	}, nil
}
