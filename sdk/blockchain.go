// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/net"
	"net/http"
)

type BlockchainService service

// returns BlockInfo for passed block's height
func (b *BlockchainService) GetBlockByHeight(ctx context.Context, height *Height) (*BlockInfo, error) {
	if height == nil || height.Int64() == 0 {
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

// returns Transaction's inside of block at passed height
func (b *BlockchainService) GetBlockTransactions(ctx context.Context, height *Height) ([]Transaction, error) {
	if height == nil || height.Int64() == 0 {
		return nil, ErrNilOrZeroHeight
	}

	url := net.NewUrl(fmt.Sprintf(blockGetTransactionRoute, height))

	var data bytes.Buffer

	resp, err := b.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, &data)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return MapTransactions(&data)
}

// returns BlockInfo's for range block height - (block height + limit)
// Example: GetBlocksByHeightWithLimit(ctx, 1, 25) => [BlockInfo25, BlockInfo24, ..., BlockInfo1]
func (b *BlockchainService) GetBlocksByHeightWithLimit(ctx context.Context, height *Height, limit *Amount) ([]*BlockInfo, error) {
	if height == nil || height.Int64() == 0 {
		return nil, ErrNilOrZeroHeight
	}

	if limit == nil || limit.Int64() == 0 {
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

func (b *BlockchainService) GetBlockchainHeight(ctx context.Context) (*Height, error) {
	bh := &struct {
		Height heightDTO `json:"height"`
	}{}

	resp, err := b.client.doNewRequest(ctx, http.MethodGet, blockHeightRoute, nil, &bh)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return nil, err
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
