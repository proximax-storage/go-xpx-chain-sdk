// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/proximax-storage/go-xpx-utils/net"
)

type AccountService service

func (a *AccountService) GetAccountProperties(ctx context.Context, address *Address) (*AccountProperties, error) {
	if address == nil {
		return nil, ErrNilAddress
	}

	if len(address.Address) == 0 {
		return nil, ErrBlankAddress
	}

	url := net.NewUrl(fmt.Sprintf(accountPropertiesRoute, address.Address))

	dto := &accountPropertiesDTO{}

	resp, err := a.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

func (a *AccountService) GetAccountsProperties(ctx context.Context, addresses ...*Address) ([]*AccountProperties, error) {
	if len(addresses) == 0 {
		return nil, ErrEmptyAddressesIds
	}

	addrs := struct {
		Messages []string `json:"addresses"`
	}{
		Messages: make([]string, len(addresses)),
	}

	for i, address := range addresses {
		addrs.Messages[i] = address.Address
	}

	dtos := accountPropertiesDTOs(make([]*accountPropertiesDTO, 0))

	resp, err := a.client.doNewRequest(ctx, http.MethodPost, accountsPropertiesRoute, addrs, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dtos.toStruct()
}

func (a *AccountService) GetAccountInfo(ctx context.Context, address *Address) (*AccountInfo, error) {
	if address == nil {
		return nil, ErrNilAddress
	}

	if len(address.Address) == 0 {
		return nil, ErrBlankAddress
	}

	url := net.NewUrl(fmt.Sprintf(accountRoute, address.Address))

	dto := &accountInfoDTO{}

	resp, err := a.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(a.client.config.reputationConfig)
}

func (a *AccountService) GetAccountsInfo(ctx context.Context, addresses ...*Address) ([]*AccountInfo, error) {
	if len(addresses) == 0 {
		return nil, ErrEmptyAddressesIds
	}

	addrs := struct {
		Messages []string `json:"addresses"`
	}{
		Messages: make([]string, len(addresses)),
	}

	for i, address := range addresses {
		addrs.Messages[i] = address.Address
	}

	dtos := accountInfoDTOs(make([]*accountInfoDTO, 0))

	resp, err := a.client.doNewRequest(ctx, http.MethodPost, accountsRoute, addrs, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dtos.toStruct(a.client.config.reputationConfig)
}

func (a *AccountService) GetMultisigAccountInfo(ctx context.Context, address *Address) (*MultisigAccountInfo, error) {
	if address == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(multisigAccountRoute, address.Address))

	dto := &multisigAccountInfoDTO{}

	resp, err := a.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(a.client.config.NetworkType)
}

func (a *AccountService) GetMultisigAccountGraphInfo(ctx context.Context, address *Address) (*MultisigAccountGraphInfo, error) {
	if address == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(multisigAccountGraphInfoRoute, address.Address))

	dto := &multisigAccountGraphInfoDTOS{}

	resp, err := a.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(a.client.config.NetworkType)
}

// GetAccountNames Returns friendly names for accounts.
// post @/account/names
func (a *AccountService) GetAccountNames(ctx context.Context, addr ...*Address) ([]*AccountName, error) {

	if len(addr) == 0 {
		return nil, ErrEmptyAddressesIds
	}

	dtos := accountNamesDTOs(make([]*accountNamesDTO, 0))

	resp, err := a.client.doNewRequest(ctx, http.MethodPost, accountNamesRoute, &addresses{addr}, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{400: ErrInvalidRequest, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dtos.toStruct()
}
