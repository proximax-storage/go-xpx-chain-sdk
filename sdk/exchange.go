// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
)

type ExchangeService struct {
	*service
	ResolveService *ResolverService
}

func (e *ExchangeService) GetAccountExchangeInfo(ctx context.Context, account *PublicAccount) (*UserExchangeInfo, error) {
	if account == nil {
		return nil, ErrNilAddress
	}
	//
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

	return nil, nil
}

// Return offers with same operation typr and mosaic id.
// Example: If you want to buy Storage units, you need to call GetExchangeOfferByAssetId(StorageMosaicId, SellOffer)
func (e *ExchangeService) GetExchangeOfferByAssetId(ctx context.Context, assetId AssetId, offerType OfferType) ([]*OfferInfo, error) {
	_, err := e.ResolveService.GetMosaicInfoByAssetId(ctx, assetId)
	if err != nil {
		return nil, err
	}

	//url := net.NewUrl(fmt.Sprintf(drivesOfAccountRoute, driveKey.PublicKey, filter))
	//
	//dto := &driveDTOs{}
	//
	//resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
	//	return nil, err
	//}

	return nil, nil
}
