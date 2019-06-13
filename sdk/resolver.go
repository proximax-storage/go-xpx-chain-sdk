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

// returns MosaicInfo from blockchain identifier
func (ref *ResolverService) GetMosaicInfoByBlockchainId(ctx context.Context, blockchainId BlockchainId) (*MosaicInfo, error) {
	if blockchainId == nil {
		return nil, ErrNilBlockchainId
	}

	switch blockchainId.Type() {
	case NamespaceBlockchainIdType:
		namespaceId := blockchainId.(*NamespaceId)
		namespaceInfo, err := ref.NamespaceService.GetNamespaceInfo(ctx, namespaceId)

		if err != nil {
			return nil, err
		}

		if namespaceInfo.Alias == nil || namespaceInfo.Alias.MosaicId() == nil {
			return nil, errors.New("Namespace is not aliased to Mosaic")
		}

		return ref.MosaicService.GetMosaicInfo(ctx, namespaceInfo.Alias.MosaicId())
	case MosaicBlockchainIdType:
		mosaicId := blockchainId.(*MosaicId)
		return ref.MosaicService.GetMosaicInfo(ctx, mosaicId)
	}

	return nil, ErrUnknownBlockchainType
}

// returns an array of MosaicInfo from blockchain identifiers
func (ref *ResolverService) GetMosaicInfosByBlockchainIds(ctx context.Context, blockchainIds ...BlockchainId) ([]*MosaicInfo, error) {
	if len(blockchainIds) == 0 {
		return nil, ErrEmptyBlockchainIds
	}

	var err error = nil

	mosaicInfos := make([]*MosaicInfo, len(blockchainIds))
	namespaceIds := make([]*NamespaceId, 0)
	mosaicIds := make([]*MosaicId, 0)

	for _, blockchainId := range blockchainIds {
		if blockchainId == nil {
			return nil, ErrNilBlockchainId
		}

		switch blockchainId.Type() {
		case NamespaceBlockchainIdType:
			namespaceId := blockchainId.(*NamespaceId)
			namespaceIds = append(namespaceIds, namespaceId)
		case MosaicBlockchainIdType:
			mosaicId := blockchainId.(*MosaicId)
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
		mosaicInfo, err := ref.GetMosaicInfoByBlockchainId(ctx, namespaceId)

		if err != nil {
			return nil, err
		}

		mosaicInfos = append(mosaicInfos, mosaicInfo)
	}

	return mosaicInfos, nil
}
