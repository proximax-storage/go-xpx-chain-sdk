// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/proximax-storage/go-xpx-utils/net"
)

type StorageV2Service service

func (s *StorageV2Service) GetDrive(ctx context.Context, driveKey *PublicAccount) (*BcDrive, error) {
	if driveKey == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(driveRouteV2, driveKey.PublicKey))

	dto := &bcDriveDTO{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *StorageV2Service) GetDrives(ctx context.Context, bdpOpts *BcDrivesPageOptions) (*BcDrivesPage, error) {
	bcdspDTO := &bcDrivesPageDTO{}

	u, err := addOptions(drivesRouteV2, bdpOpts)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, u, nil, &bcdspDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return bcdspDTO.toStruct(s.client.NetworkType())
}

func (s *StorageV2Service) GetReplicator(ctx context.Context, replicatorKey *PublicAccount) (*Replicator, error) {
	if replicatorKey == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(replicatorRouteV2, replicatorKey.PublicKey))

	dto := &replicatorV2DTO{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *StorageV2Service) GetReplicators(ctx context.Context, rpOpts *ReplicatorsPageOptions) (*ReplicatorsPage, error) {
	rspDTO := &replicatorsPageDTO{}

	u, err := addOptions(replicatorsRouteV2, rpOpts)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, u, nil, &rspDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return rspDTO.toStruct(s.client.NetworkType())
}

func (s *StorageV2Service) GetDownloadChannelInfo(ctx context.Context, downloadChannelId *Hash) ([]*DownloadChannel, error) {
	if downloadChannelId == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(downloadChannelRouteV2, downloadChannelId))

	dto := &downloadChannelDTOs{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *StorageV2Service) GetDownloadChannels(ctx context.Context, rpOpts *DownloadChannelsPageOptions) (*DownloadChannelsPage, error) {
	dcspDTO := &downloadChannelsPageDTO{}

	u, err := addOptions(downloadChannelsRouteV2, rpOpts)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, u, nil, &dcspDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dcspDTO.toStruct(s.client.NetworkType())
}
