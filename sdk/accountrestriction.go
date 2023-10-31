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

type AccountRestrictionService service

func (s *AccountRestrictionService) SearchAccountRestrictions(ctx context.Context, tpOpts *AccountRestrictionsPageOptions) (*AccountRestrictionsPage, error) {
	accResDTO := AccountRestrictionsPageDTO{}
	u, err := addOptions(accountRestrictionsSimpleRoute, tpOpts)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, u, nil, &accResDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}
	return accResDTO.toStruct(s.client.NetworkType())
}

func (s *AccountRestrictionService) GetAccountRestrictions(ctx context.Context, address *Address) (*AccountRestrictions, error) {
	accResDTO := AccountRestrictionsDtoContainer{}

	url := net.NewUrl(fmt.Sprintf(accountRestrictionsRoute, address.Address))

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, &accResDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}
	return accResDTO.toStruct(s.client.NetworkType())
}
