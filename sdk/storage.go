// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-utils/net"
	"net/http"
)

type StorageService struct {
	*service
	LockService *LockService
}

func (s *StorageService) GetDrive(ctx context.Context, driveKey *PublicAccount) (*Drive, error) {
	if driveKey == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(driveRoute, driveKey.PublicKey))

	dto := &driveDTO{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

type DriveParticipantFilter string

const (
	AllDriveRoles   DriveParticipantFilter = ""
	OwnerDrive      DriveParticipantFilter = "/owner"
	ReplicatorDrive DriveParticipantFilter = "/replicator"
)

func (s *StorageService) GetAccountDrives(ctx context.Context, driveKey *PublicAccount, filter DriveParticipantFilter) ([]*Drive, error) {
	if driveKey == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(drivesOfAccountRoute, driveKey.PublicKey, filter))

	dto := &driveDTOs{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		// Skip ErrResourceNotFound
		// not return err
		return nil, nil
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *StorageService) GetVerificationStatus(ctx context.Context, driveKey *PublicAccount) (*VerificationStatus, error) {
	if driveKey == nil {
		return nil, ErrNilAddress
	}

	compositeHash, err := CalculateCompositeHash(&Hash{}, driveKey.Address)
	if err != nil {
		return nil, err
	}

	lockInfo, err := s.LockService.GetSecretLockInfo(ctx, compositeHash)
	if err != nil {
		switch e := err.(type) {
		case *HttpError:
			if e.StatusCode == 404 {
				return &VerificationStatus{
					Active:    false,
					Available: true,
				}, nil
			} else {
				return nil, err
			}
		default:
			return nil, err
		}

		return nil, err
	}

	if lockInfo.HashAlgorithm != Internal_Hash_Type {
		return nil, errors.New("wrong type of drive secret lock")
	}

	return &VerificationStatus{
		Active:    lockInfo.Status == Unused,
		Available: false,
	}, nil
}

func (s *StorageService) GetDownloadInfo(ctx context.Context, operationToken *Hash) (*DownloadInfo, error) {
	if operationToken == nil {
		return nil, ErrNilHash
	}

	url := net.NewUrl(fmt.Sprintf(downloadInfoRoute, operationToken.String()))

	dto := &downloadInfoDTO{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *StorageService) GetAccountDownloadInfos(ctx context.Context, recipient *PublicAccount) ([]*DownloadInfo, error) {
	if recipient == nil {
		return nil, ErrNilAccount
	}

	url := net.NewUrl(fmt.Sprintf(accountDownloadInfosRoute, recipient.PublicKey))

	dto := &downloadInfoDTOs{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *StorageService) GetDriveDownloadInfos(ctx context.Context, drive *PublicAccount) ([]*DownloadInfo, error) {
	if drive == nil {
		return nil, ErrNilAccount
	}

	url := net.NewUrl(fmt.Sprintf(driveDownloadInfosRoute, drive.PublicKey))

	dto := &downloadInfoDTOs{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}
