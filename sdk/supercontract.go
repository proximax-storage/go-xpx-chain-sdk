// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
)

type SuperContractService service

func (s *SuperContractService) GetSuperContract(ctx context.Context, contractKey *PublicAccount) (*SuperContract, error) {
	if contractKey == nil {
		return nil, ErrNilAddress
	}

	//url := net.NewUrl(fmt.Sprintf(driveRoute, driveKey.PublicKey))
	//
	//dto := &driveDTO{}
	//
	//resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
	//	return nil, err
	//}
	//
	//return dto.toStruct(s.client.NetworkType())

	return nil, nil
}

type SuperContractParticipantFilter string

const (
	AllRolesContract   SuperContractParticipantFilter = ""
	OwnerContact       SuperContractParticipantFilter = "/owner"
	ExecutorContact    SuperContractParticipantFilter = "/executor"
)

func (s *SuperContractService) GetAccountSuperContracts(ctx context.Context, accountKey *PublicAccount, filter SuperContractParticipantFilter) ([]*SuperContract, error) {
	if accountKey == nil {
		return nil, ErrNilAddress
	}
	//
	//url := net.NewUrl(fmt.Sprintf(drivesOfAccountRoute, driveKey.PublicKey, filter))
	//
	//dto := &driveDTOs{}
	//
	//resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	//if err != nil {
	//	// Skip ErrResourceNotFound
	//	// not return err
	//	return nil, nil
	//}
	//
	//if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
	//	return nil, err
	//}
	//
	//return dto.toStruct(s.client.NetworkType())

	return nil, nil
}

func (s *SuperContractService) GetExectutionStatus(ctx context.Context, operationHash *Hash) (*Operation, error) {
	return nil, nil
}