// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/proximax-storage/go-xpx-utils/net"
)

type LiquidityProviderService service

func (lp *LiquidityProviderService) GetLiquidityProviders(ctx context.Context, lpOptions *LiquidityProviderPageOptions) (*LiquidityProviderPage, error) {
	lpsDTO := &liquidityProvidersPageDTO{}

	u, err := addOptions(liquidityProvidersRoute, lpOptions)
	if err != nil {
		return nil, err
	}

	resp, err := lp.client.doNewRequest(ctx, http.MethodGet, u, nil, &lpsDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return lpsDTO.toStruct(lp.client.NetworkType())
}

func (lp *LiquidityProviderService) GetLiquidityProvider(ctx context.Context, provider *PublicAccount) (*LiquidityProvider, error) {
	if provider == nil {
		return nil, ErrNilAccount
	}

	url := net.NewUrl(fmt.Sprintf(liquidityProviderRoute, provider.PublicKey))

	dto := &liquidityProviderDTO{}

	resp, err := lp.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(lp.client.NetworkType())
}
