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

type MosaicService service

func (ref *MosaicService) GetMosaicInfo(ctx context.Context, mosaicId *MosaicId) (*MosaicInfo, error) {
	if mosaicId == nil {
		return nil, ErrNilMosaicId
	}

	url := net.NewUrl(fmt.Sprintf(mosaicRoute, mosaicId.toHexString()))

	dto := &mosaicInfoDTO{}

	resp, err := ref.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	mscInfo, err := dto.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return mscInfo, nil
}

func (ref *MosaicService) GetMosaicInfos(ctx context.Context, mscIds []*MosaicId) ([]*MosaicInfo, error) {
	if len(mscIds) == 0 {
		return nil, ErrEmptyMosaicIds
	}

	dtos := mosaicInfoDTOs(make([]*mosaicInfoDTO, 0))

	resp, err := ref.client.doNewRequest(ctx, http.MethodPost, mosaicsRoute, &mosaicIds{mscIds}, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{400: ErrInvalidRequest, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	mscInfos, err := dtos.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return mscInfos, nil
}

// GetMosaicsNames Get readable names for a set of mosaics
// post @/mosaic/names
func (ref *MosaicService) GetMosaicsNames(ctx context.Context, mscIds ...*MosaicId) ([]*MosaicName, error) {
	if len(mscIds) == 0 {
		return nil, ErrEmptyMosaicIds
	}

	dtos := mosaicNameDTOs{}

	resp, err := ref.client.doNewRequest(ctx, http.MethodPost, mosaicNamesRoute, &mosaicIds{mscIds}, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{400: ErrInvalidRequest, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dtos.toStruct()
}

// GetMosaicLevy returns mosaic levy
// get @/mosaic/%s/levy
func (ref *MosaicService) GetMosaicLevy(ctx context.Context, mosaicId *MosaicId) (*MosaicLevy, error) {
	if mosaicId == nil {
		return nil, ErrNilMosaicId
	}

	url := net.NewUrl(fmt.Sprintf(mosaicLevyRoute, mosaicId.toHexString()))
	mosaicLevyDTO := &mosaicLevyDTO{}

	resp, err := ref.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, mosaicLevyDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	mosaicLevy, err := mosaicLevyDTO.toStruct()
	if err != nil {
		return nil, err
	}

	return mosaicLevy, nil
}
