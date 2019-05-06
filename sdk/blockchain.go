// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/net"
	"math/big"
	"net/http"
)

type BlockchainService service

// returns info for block with passed height
func (b *BlockchainService) GetBlockByHeight(ctx context.Context, height *big.Int) (*BlockInfo, error) {
	if height == nil || height.Int64() == 0 {
		return nil, ErrNilOrZeroHeight
	}

	u := fmt.Sprintf(blockByHeightRoute, height)

	dto := &blockInfoDTO{}

	resp, err := b.client.DoNewRequest(ctx, http.MethodGet, u, nil, &dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

// get transactions inside of block with passed height
func (b *BlockchainService) GetBlockTransactions(ctx context.Context, height *big.Int) ([]Transaction, error) {
	if height == nil || height.Int64() == 0 {
		return nil, ErrNilOrZeroHeight
	}

	url := net.NewUrl(fmt.Sprintf(blockGetTransactionRoute, height))

	var data bytes.Buffer

	resp, err := b.client.DoNewRequest(ctx, http.MethodGet, url.Encode(), nil, &data)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return MapTransactions(&data)
}

func (b *BlockchainService) GetBlocksByHeightWithLimit(ctx context.Context, height, limit *big.Int) ([]*BlockInfo, error) {
	if height == nil || height.Int64() == 0 {
		return nil, ErrNilOrZeroHeight
	}

	if limit == nil || limit.Int64() == 0 {
		return nil, ErrNilOrZeroLimit
	}

	url := net.NewUrl(fmt.Sprintf(blockInfoRoute, height, limit))

	dtos := blockInfoDTOs(make([]*blockInfoDTO, 0))

	resp, err := b.client.DoNewRequest(ctx, http.MethodGet, url.Encode(), nil, &dtos)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dtos.toStruct()
}

// returns blockchain height
func (b *BlockchainService) GetBlockchainHeight(ctx context.Context) (*big.Int, error) {
	bh := &struct {
		Height uint64DTO `json:"height"`
	}{}

	resp, err := b.client.DoNewRequest(ctx, http.MethodGet, blockHeightRoute, nil, &bh)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return nil, err
	}

	return bh.Height.toBigInt(), nil
}

// returns blockchain score
func (b *BlockchainService) GetBlockchainScore(ctx context.Context) (*big.Int, error) {
	cs := &chainScoreDTO{}
	resp, err := b.client.DoNewRequest(ctx, http.MethodGet, blockScoreRoute, nil, &cs)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return nil, err
	}

	return cs.toStruct(), nil
}

// returns blockchain storage information
func (b *BlockchainService) GetBlockchainStorage(ctx context.Context) (*BlockchainStorageInfo, error) {
	bstorage := &BlockchainStorageInfo{}
	resp, err := b.client.DoNewRequest(ctx, http.MethodGet, blockStorageRoute, nil, &bstorage)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return nil, err
	}

	return bstorage, nil
}
