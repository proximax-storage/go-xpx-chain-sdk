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

type LockService service

func (s *LockService) GetHashLockInfosByAccount(ctx context.Context, account *PublicAccount) ([]*HashLockInfo, error) {
	if account == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(hashLocksRoute, account.PublicKey))

	dto := &hashLockInfoDTOs{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *LockService) GetHashLockInfo(ctx context.Context, hash *Hash) (*HashLockInfo, error) {
	if hash == nil {
		return nil, ErrNilHash
	}

	url := net.NewUrl(fmt.Sprintf(hashLockRoute, hash.String()))

	dto := &hashLockInfoDTO{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *LockService) GetSecretLockInfosByAccount(ctx context.Context, account *PublicAccount) ([]*SecretLockInfo, error) {
	if account == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(secretLocksByAccountRoute, account.PublicKey))

	dto := &secretLockInfoDTOs{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *LockService) GetSecretLockInfo(ctx context.Context, compositeHash *Hash) (*SecretLockInfo, error) {
	if compositeHash == nil {
		return nil, ErrNilHash
	}

	url := net.NewUrl(fmt.Sprintf(secretLockRoute, compositeHash.String()))

	dto := &secretLockInfoDTO{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}

func (s *LockService) GetSecretLockInfosBySecret(ctx context.Context, secret *Hash) ([]*SecretLockInfo, error) {
	if secret == nil {
		return nil, ErrNilSecret
	}

	url := net.NewUrl(fmt.Sprintf(secretLocksBySecretRoute, secret.String()))

	dto := &secretLockInfoDTOs{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}
