// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/proximax-storage/go-xpx-utils/net"
)

type SdaExchangeService struct {
	*service
	ResolveService *ResolverService
}

func (e *SdaExchangeService) GetAccountSdaExchangeInfo(ctx context.Context, account *PublicAccount) (*UserSdaExchangeInfo, error) {
	if account == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(exchangeSdaRoute, account.PublicKey))

	dto := &sdaExchangeDTO{}

	resp, err := e.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(e.client.NetworkType())
}

// Return offers with same mosaic id give or mosaic id get.
// offerType = give OR offerType = get ONLY
func (e *SdaExchangeService) GetSdaExchangeOfferByAssetId(ctx context.Context, assetId AssetId, offerType string) ([]*SdaOfferBalance, error) {
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
	default:
		return nil, errors.New("unknown assetID type")
	}

	url := net.NewUrl(fmt.Sprintf(sdaOffersByMosaicRoute, offerType, mosaicId.toHexString()))

	dto := &sdaOfferBalanceDTOs{}

	resp, err := e.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(e.client.NetworkType())
}
