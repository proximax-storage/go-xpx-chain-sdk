// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/proximax-storage/go-xpx-utils/net"
)

// NamespaceService provides a set of methods for obtaining information about the namespace
type NamespaceService service

func (ref *NamespaceService) GetNamespaceInfo(ctx context.Context, nsId *NamespaceId) (*NamespaceInfo, error) {
	if nsId == nil {
		return nil, ErrNilNamespaceId
	}

	nsInfoDTO := &namespaceInfoDTO{}

	url := net.NewUrl(fmt.Sprintf(namespaceRoute, nsId.toHexString()))

	resp, err := ref.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, nsInfoDTO)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	nsInfo, err := nsInfoDTO.toStruct()
	if err != nil {
		return nil, err
	}

	if err = ref.buildNamespaceHierarchy(ctx, nsInfo); err != nil {
		return nil, err
	}

	return nsInfo, nil
}

// returns NamespaceInfo's corresponding to passed Address and NamespaceId with maximum limit
// TODO: fix pagination
func (ref *NamespaceService) GetNamespaceInfosFromAccount(ctx context.Context, address *Address, nsId *NamespaceId,
	pageSize int) ([]*NamespaceInfo, error) {
	if address == nil {
		return nil, ErrNilAddress
	}

	url := net.NewUrl(fmt.Sprintf(namespacesFromAccountRoutes, address.Address))

	if nsId != nil {
		url.SetParam("id", nsId.toHexString())
	}

	if pageSize > 0 {
		url.SetParam("pageSize", strconv.Itoa(pageSize))
	}

	dtos := namespaceInfoDTOs(make([]*namespaceInfoDTO, 0))

	resp, err := ref.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	nsInfos, err := dtos.toStruct()
	if err != nil {
		return nil, err
	}

	if err = ref.buildNamespacesHierarchy(ctx, nsInfos); err != nil {
		return nil, err
	}

	return nsInfos, nil
}

// returns NamespaceInfo's corresponding to passed Address's and NamespaceId with maximum limit
// TODO: fix pagination
func (ref *NamespaceService) GetNamespaceInfosFromAccounts(ctx context.Context, addrs []*Address, nsId *NamespaceId,
	pageSize int) ([]*NamespaceInfo, error) {
	if len(addrs) == 0 {
		return nil, ErrEmptyAddressesIds
	}

	url := net.NewUrl(namespacesFromAccountsRoute)

	if nsId != nil {
		url.AddParam("id", nsId.toHexString())
	}

	if pageSize > 0 {
		url.AddParam("pageSize", strconv.Itoa(pageSize))
	}

	dtos := namespaceInfoDTOs(make([]*namespaceInfoDTO, 0))

	resp, err := ref.client.doNewRequest(ctx, http.MethodPost, url.Encode(), &addresses{addrs}, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{400: ErrInvalidRequest, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	nsInfos, err := dtos.toStruct()
	if err != nil {
		return nil, err
	}

	if err = ref.buildNamespacesHierarchy(ctx, nsInfos); err != nil {
		return nil, err
	}

	return nsInfos, nil
}

func (ref *NamespaceService) GetNamespaceNames(ctx context.Context, nsIds []*NamespaceId) ([]*NamespaceName, error) {
	if len(nsIds) == 0 {
		return nil, ErrEmptyNamespaceIds
	}

	dtos := namespaceNameDTOs(make([]*namespaceNameDTO, 0))

	resp, err := ref.client.doNewRequest(ctx, http.MethodPost, namespaceNamesRoute, &namespaceIds{nsIds}, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{400: ErrInvalidRequest, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dtos.toStruct()
}

// GetLinkedMosaicId
// @/namespace/%s
func (ref *NamespaceService) GetLinkedMosaicId(ctx context.Context, namespaceId *NamespaceId) (*MosaicId, error) {
	if namespaceId == nil {
		return nil, ErrNilAddress
	}

	info, err := ref.GetNamespaceInfo(ctx, namespaceId)

	if err != nil {
		return nil, err
	}

	return info.Alias.MosaicId(), nil
}

// GetLinkedAddress
// @/namespace/%s
func (ref *NamespaceService) GetLinkedAddress(ctx context.Context, namespaceId *NamespaceId) (*Address, error) {
	if namespaceId == nil {
		return nil, ErrNilAddress
	}

	info, err := ref.GetNamespaceInfo(ctx, namespaceId)

	if err != nil {
		return nil, err
	}

	return info.Alias.Address(), nil
}

func (ref *NamespaceService) buildNamespaceHierarchy(ctx context.Context, nsInfo *NamespaceInfo) error {
	if nsInfo == nil || nsInfo.Parent == nil {
		return nil
	}

	if nsInfo.Parent.NamespaceId == nil || nsInfo.Parent.NamespaceId.Id() == 0 {
		return nil
	}

	parentNsInfo, err := ref.GetNamespaceInfo(ctx, nsInfo.Parent.NamespaceId)
	if err != nil {
		return err
	}

	nsInfo.Parent = parentNsInfo

	return ref.buildNamespaceHierarchy(ctx, nsInfo.Parent)
}

func (ref *NamespaceService) buildNamespacesHierarchy(ctx context.Context, nsInfos []*NamespaceInfo) error {
	var err error

	for _, nsInfo := range nsInfos {
		if err = ref.buildNamespaceHierarchy(ctx, nsInfo); err != nil {
			return err
		}
	}

	return nil
}
