// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
)

type StorageService service

func (s *StorageService) GetDrive(ctx context.Context, driveKey *PublicAccount) (*Drive, error) {
	//if address == nil {
	//	return nil, ErrNilAddress
	//}
	//
	//if len(address.Address) == 0 {
	//	return nil, ErrBlankAddress
	//}
	//
	//url := net.NewUrl(fmt.Sprintf(accountPropertiesRoute, address.Address))
	//
	//dto := &accountPropertiesDTO{}
	//
	//resp, err := a.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
	//	return nil, err
	//}

	return nil, nil
}

type DriveParticipantFilter string

const (
	AllRoles   DriveParticipantFilter = ""
	Owner      DriveParticipantFilter = "owner"
	Replicator DriveParticipantFilter = "replicator"
)

func (s *StorageService) GetAccountDrives(ctx context.Context, driveKey *PublicAccount, filter DriveParticipantFilter) ([]*Drive, error) {
	//if address == nil {
	//	return nil, ErrNilAddress
	//}
	//
	//if len(address.Address) == 0 {
	//	return nil, ErrBlankAddress
	//}
	//
	//url := net.NewUrl(fmt.Sprintf(accountPropertiesRoute, address.Address))
	//
	//dto := &accountPropertiesDTO{}
	//
	//resp, err := a.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
	//	return nil, err
	//}

	return nil, nil
}
