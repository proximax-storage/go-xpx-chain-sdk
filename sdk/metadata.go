package sdk

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/proximax-storage/go-xpx-utils/net"
	"net/http"
)

type MetadataService service

func (ref *MetadataService) GetAddressMetadatasInfo(ctx context.Context, addresses ...string) ([]*AddressMetadataInfo, error) {
	if len(addresses) == 0 {
		return nil, ErrMetadataEmptyAddresses
	}

	addressesDto := struct {
		Addresses []string `json:"metadataIds"`
	}{
		Addresses: addresses,
	}

	dtos := addressMetadataInfoDTOs(make([]*addressMetadataInfoDTO, 0))

	resp, err := ref.client.DoNewRequest(ctx, http.MethodPost, metadatasInfoRoute, addressesDto, &dtos)
	if err != nil {
		return nil, errors.Wrapf(err, "within POST request %s", metadatasInfoRoute)
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	infos, err := dtos.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (ref *MetadataService) GetMosaicMetadatasInfo(ctx context.Context, mosaicIds ...*MosaicId) ([]*MosaicMetadataInfo, error) {
	if len(mosaicIds) == 0 {
		return nil, ErrMetadataEmptyMosaicIds
	}

	mosaicsDto := struct {
		MosaicIds []string `json:"metadataIds"`
	}{
		MosaicIds: make([]string, len(mosaicIds)),
	}

	for i, m := range mosaicIds {
		mosaicsDto.MosaicIds[i] = m.toHexString()
	}

	dtos := mosaicMetadataInfoDTOs(make([]*mosaicMetadataInfoDTO, 0))

	resp, err := ref.client.DoNewRequest(ctx, http.MethodPost, metadatasInfoRoute, mosaicsDto, &dtos)
	if err != nil {
		return nil, errors.Wrapf(err, "within POST request %s", metadatasInfoRoute)
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	infos, err := dtos.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (ref *MetadataService) GetNamespaceMetadatasInfo(ctx context.Context, namespaceIds ...*NamespaceId) ([]*NamespaceMetadataInfo, error) {
	if len(namespaceIds) == 0 {
		return nil, ErrMetadataEmptyNamespaceIds
	}

	namespacesDto := struct {
		NamespaceIds []string `json:"metadataIds"`
	}{
		NamespaceIds: make([]string, len(namespaceIds)),
	}

	for i, n := range namespaceIds {
		namespacesDto.NamespaceIds[i] = n.toHexString()
	}

	dtos := namespaceMetadataInfoDTOs(make([]*namespaceMetadataInfoDTO, 0))

	resp, err := ref.client.DoNewRequest(ctx, http.MethodPost, metadatasInfoRoute, namespacesDto, &dtos)
	if err != nil {
		return nil, errors.Wrapf(err, "within POST request %s", metadatasInfoRoute)
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	infos, err := dtos.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

func (ref *MetadataService) GetMetadataByAddress(ctx context.Context, address string) (*AddressMetadataInfo, error) {
	if len(address) == 0 {
		return nil, ErrMetadataNilAdress
	}

	url := net.NewUrl(fmt.Sprintf(metadataByAccountRoute, address))

	dto := addressMetadataInfoDTO{}

	resp, err := ref.client.DoNewRequest(ctx, http.MethodGet, url.Encode(), nil, &dto)
	if err != nil {
		return nil, errors.Wrapf(err, "within GET request %s", url.Encode())
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	info, err := dto.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (ref *MetadataService) GetMetadataByMosaicId(ctx context.Context, mosaicId *MosaicId) (*MosaicMetadataInfo, error) {
	if mosaicId == nil {
		return nil, ErrMetadataNilMosaicId
	}

	url := net.NewUrl(fmt.Sprintf(metadataByMosaicRoute, mosaicId.toHexString()))

	dto := mosaicMetadataInfoDTO{}

	resp, err := ref.client.DoNewRequest(ctx, http.MethodGet, url.Encode(), nil, &dto)
	if err != nil {
		return nil, errors.Wrapf(err, "within GET request %s", url.Encode())
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	info, err := dto.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (ref *MetadataService) GetMetadataByNamespaceId(ctx context.Context, namespaceId *NamespaceId) (*NamespaceMetadataInfo, error) {
	if namespaceId == nil {
		return nil, ErrMetadataNilNamespaceId
	}

	url := net.NewUrl(fmt.Sprintf(metadataByNamespaceRoute, namespaceId.toHexString()))

	dto := namespaceMetadataInfoDTO{}

	resp, err := ref.client.DoNewRequest(ctx, http.MethodGet, url.Encode(), nil, &dto)
	if err != nil {
		return nil, errors.Wrapf(err, "within GET request %s", url.Encode())
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	info, err := dto.toStruct(ref.client.config.NetworkType)
	if err != nil {
		return nil, err
	}

	return info, nil
}
