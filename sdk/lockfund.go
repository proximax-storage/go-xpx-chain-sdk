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

type LockFundService service

func (s *LockFundService) GetLockFundKeyRecords(ctx context.Context, accountKey *PublicAccount) ([]*LockFundKeyRecord, error) {
	if accountKey == nil {
		return nil, ErrNilAccount
	}

	url := net.NewUrl(fmt.Sprintf(lockFundKeyRecordGroupRoute, accountKey.PublicKey))

	dto := &LockFundKeyRecordGroupDtos{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *LockFundService) GetLockFundHeightRecords(ctx context.Context, height Height) ([]*LockFundHeightRecord, error) {
	if height == 0 {
		return nil, ErrArgumentNotValid
	}

	url := net.NewUrl(fmt.Sprintf(lockFundHeightRecordGroupRoute, height))

	dto := &LockFundHeightRecordGroupDtos{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}
