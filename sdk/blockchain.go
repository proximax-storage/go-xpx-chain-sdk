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

type BlockchainService service

// returns BlockInfo for passed block's height
func (b *BlockchainService) GetBlockByHeight(ctx context.Context, height Height) (*BlockInfo, error) {
	if height == 0 {
		return nil, ErrNilOrZeroHeight
	}

	u := fmt.Sprintf(blockByHeightRoute, height)

	dto := &blockInfoDTO{}

	resp, err := b.client.doNewRequest(ctx, http.MethodGet, u, nil, &dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

// returns BlockInfo's for range block height - (block height + limit)
// Example: GetBlocksByHeightWithLimit(ctx, 1, 25) => [BlockInfo25, BlockInfo24, ..., BlockInfo1]
func (b *BlockchainService) GetBlocksByHeightWithLimit(ctx context.Context, height Height, limit Amount) ([]*BlockInfo, error) {
	if height == 0 {
		return nil, ErrNilOrZeroHeight
	}

	if limit == 0 {
		return nil, ErrNilOrZeroLimit
	}

	url := net.NewUrl(fmt.Sprintf(blockInfoRoute, height, limit))

	dtos := blockInfoDTOs(make([]*blockInfoDTO, 0))

	resp, err := b.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dtos.toStruct()
}

func (b *BlockchainService) GetBlockchainHeight(ctx context.Context) (Height, error) {
	bh := &struct {
		Height uint64DTO `json:"height"`
	}{}

	resp, err := b.client.doNewRequest(ctx, http.MethodGet, blockHeightRoute, nil, &bh)
	if err != nil {
		return 0, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return 0, err
	}

	return bh.Height.toStruct(), nil
}

func (b *BlockchainService) GetBlockchainScore(ctx context.Context) (*ChainScore, error) {
	cs := &chainScoreDTO{}
	resp, err := b.client.doNewRequest(ctx, http.MethodGet, blockScoreRoute, nil, &cs)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return nil, err
	}

	return cs.toStruct(), nil
}

func (b *BlockchainService) GetBlockchainStorage(ctx context.Context) (*BlockchainStorageInfo, error) {
	bstorage := &BlockchainStorageInfo{}
	resp, err := b.client.doNewRequest(ctx, http.MethodGet, blockStorageRoute, nil, &bstorage)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return nil, err
	}

	return bstorage, nil
}
