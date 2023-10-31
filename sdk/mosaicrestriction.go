// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/net"
	"golang.org/x/crypto/sha3"
	"net/http"
)

type MosaicRestrictionService service

func (s *MosaicRestrictionService) SearchMosaicRestrictions(ctx context.Context, tpOpts *MosaicRestrictionsPageOptions) (*MosaicRestrictionsPage, error) {
	accResDTO := MosaicRestrictionsPageDto{}
	u, err := addOptions(mosaicRestrictionsSimpleRoute, tpOpts)
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

func CalculateMosaicAddressRestrictionUniqueId(mosaic *Mosaic, targetAddress *Address) (*Hash, error) {
	result := sha3.New256()
	addressBytes, err := targetAddress.Decode()
	if err != nil {
		return nil, err
	}
	if _, err := result.Write(mosaic.AssetId.toLittleEndian()); err != nil {
		return nil, err
	}
	if _, err := result.Write(addressBytes[:]); err != nil {
		return nil, err
	}

	hash, err := bytesToHash(result.Sum(nil))
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func CalculateMosaicGlobalRestrictionUniqueId(mosaic *Mosaic) (*Hash, error) {
	result := sha3.New256()
	addressBytes := make([]byte, AddressSize)

	if _, err := result.Write(mosaic.AssetId.toLittleEndian()); err != nil {
		return nil, err
	}
	if _, err := result.Write(addressBytes[:]); err != nil {
		return nil, err
	}

	hash, err := bytesToHash(result.Sum(nil))
	if err != nil {
		return nil, err
	}

	return hash, nil
}
