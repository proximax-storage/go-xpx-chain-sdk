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

type ExchangeService struct {
	*service
	ResolveService *ResolverService
}

func (e *ExchangeService) GetAccountExchangeInfo(ctx context.Context, account *PublicAccount) (*UserExchangeInfo, error) {
	if account == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(exchangeRoute, account.PublicKey))

	dto := &exchangeDTO{}

	resp, err := e.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(e.client.NetworkType())
}

// Return offers with same operation type and mosaic id.
// Example: If you want to buy Storage units, you need to call GetExchangeOfferByAssetId(StorageMosaicId, SellOffer)
func (e *ExchangeService) GetExchangeOfferByAssetId(ctx context.Context, assetId AssetId, offerType OfferType) ([]*OfferInfo, error) {
	var mosaicId *MosaicId

	switch assetId.Type() {
	case NamespaceAssetIdType:
		mosaicInfo, err := e.ResolveService.GetMosaicInfoByAssetId(ctx, assetId)
		if err != nil {
			return nil, err
		}
		mosaicId = mosaicInfo.MosaicId
	case MosaicAssetIdType:
		mosaicId = assetId.(*MosaicId)
	}

	url := net.NewUrl(fmt.Sprintf(offersByMosaicRoute, offerType.String(), mosaicId.toHexString()))

	dto := &offerInfoDTOs{}

	resp, err := e.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(e.client.NetworkType())
}
