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

type MosaicRestrictionService service

func (s *MosaicRestrictionService) SearchMosaicRestrictions(ctx context.Context, tpOpts *MosaicRestrictionsPageOptions) (*MosaicRestrictionsPage, error) {
	accResDTO := MosaicRestrictionsPageDto{}
	u, err := addOptions(fmt.Sprintf(mosaicRestrictionsRoute, ""), tpOpts)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.doNewRequest(ctx, http.MethodPost, u, nil, &accResDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}
	return accResDTO.toStruct(s.client.NetworkType())
}

func (s *MosaicRestrictionService) GetMosaicRestrictions(ctx context.Context, compositeHash string) (*MosaicRestrictionEntry, error) {
	mosaicResDTO := MosaicRestrictionEntryDtoContainer{}

	url := net.NewUrl(fmt.Sprintf(mosaicRestrictionsRoute, compositeHash))

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, &mosaicResDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}
	return mosaicResDTO.toStruct()
}
