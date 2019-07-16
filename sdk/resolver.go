// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"errors"
)

// TODO: Implement resolving namespace to account
type ResolverService struct {
	*service
	NamespaceService *NamespaceService
	MosaicService    *MosaicService
}

func (ref *ResolverService) GetMosaicInfoByAssetId(ctx context.Context, assetId AssetId) (*MosaicInfo, error) {
	if assetId == nil {
		return nil, ErrNilAssetId
	}

	switch assetId.Type() {
	case NamespaceAssetIdType:
		namespaceId := assetId.(*NamespaceId)
		namespaceInfo, err := ref.NamespaceService.GetNamespaceInfo(ctx, namespaceId)

		if err != nil {
			return nil, err
		}

		if namespaceInfo.Alias == nil || namespaceInfo.Alias.MosaicId() == nil {
			return nil, errors.New("Namespace is not aliased to Mosaic")
		}

		return ref.MosaicService.GetMosaicInfo(ctx, namespaceInfo.Alias.MosaicId())
	case MosaicAssetIdType:
		mosaicId := assetId.(*MosaicId)
		return ref.MosaicService.GetMosaicInfo(ctx, mosaicId)
	}

	return nil, ErrUnknownBlockchainType
}

func (ref *ResolverService) GetMosaicInfosByAssetIds(ctx context.Context, assetIds ...AssetId) ([]*MosaicInfo, error) {
	if len(assetIds) == 0 {
		return nil, ErrEmptyAssetIds
	}

	var err error = nil

	mosaicInfos := make([]*MosaicInfo, len(assetIds))
	namespaceIds := make([]*NamespaceId, 0)
	mosaicIds := make([]*MosaicId, 0)

	for _, assetId := range assetIds {
		if assetId == nil {
			return nil, ErrNilAssetId
		}

		switch assetId.Type() {
		case NamespaceAssetIdType:
			namespaceId := assetId.(*NamespaceId)
			namespaceIds = append(namespaceIds, namespaceId)
		case MosaicAssetIdType:
			mosaicId := assetId.(*MosaicId)
			mosaicIds = append(mosaicIds, mosaicId)
		}
	}

	if len(mosaicIds) > 0 {
		mosaicInfos, err = ref.MosaicService.GetMosaicInfos(ctx, mosaicIds)

		if err != nil {
			return nil, err
		}
	}

	for _, namespaceId := range namespaceIds {
		mosaicInfo, err := ref.GetMosaicInfoByAssetId(ctx, namespaceId)

		if err != nil {
			return nil, err
		}

		mosaicInfos = append(mosaicInfos, mosaicInfo)
	}

	return mosaicInfos, nil
}
