// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/net"
	"net/http"
)

type SuperContractService service

func (s *SuperContractService) GetSuperContract(ctx context.Context, contractKey *PublicAccount) (*SuperContract, error) {
	if contractKey == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(superContractRoute, contractKey.PublicKey))

	dto := &superContractDTO{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *SuperContractService) GetDriveSuperContracts(ctx context.Context, driveKey *PublicAccount) ([]*SuperContract, error) {
	if driveKey == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(driveSuperContractsRoute, driveKey.PublicKey))

	dto := &superContractDTOs{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *SuperContractService) GetExecutionStatus(ctx context.Context, operationHash *Hash) (*Operation, error) {
	return nil, nil
}