// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/net"
	"net/http"
)

type StorageService service

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
	AllRoles   DriveParticipantFilter = ""
	Owner      DriveParticipantFilter = "/owner"
	Replicator DriveParticipantFilter = "/replicator"
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
