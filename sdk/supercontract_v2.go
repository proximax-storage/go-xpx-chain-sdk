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

type SuperContractV2Service service

func (s *SuperContractV2Service) GetSuperContractV2(ctx context.Context, superContractKey *PublicAccount) (*SuperContractV2, error) {
	if superContractKey == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(superContractRouteV2, superContractKey.PublicKey))

	dto := &superContractV2DTO{}


	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *SuperContractV2Service) GetSuperContractsV2(ctx context.Context, scPageOpts *SuperContractsV2PageOptions) (*SuperContractsV2Page, error) {
	scPageDTO := &superContractV2PageDTO{}

	u, err := addOptions(superContractsRouteV2, scPageOpts)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, u, nil, &scPageDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return scPageDTO.toStruct(s.client.NetworkType())
}