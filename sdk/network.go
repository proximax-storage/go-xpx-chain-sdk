// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type NetworkService struct {
	*service
	BlockchainService *BlockchainService
}

func (ref *NetworkService) GetNetworkType(ctx context.Context) (NetworkType, error) {
	netDTO := &networkDTO{}

	resp, err := ref.client.doNewRequest(ctx, http.MethodGet, networkRoute, nil, netDTO)

	if err != nil {
		return NotSupportedNet, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return NotSupportedNet, err
	}

	networkType := NetworkTypeFromString(netDTO.Name)

	if networkType == NotSupportedNet {
		err = errors.New(fmt.Sprintf("network %s is not supported yet by the sdk", netDTO.Name))
	}

	return networkType, err
}

func (ref *NetworkService) GetNetworkConfigAtHeight(ctx context.Context, height Height) (*BlockchainConfig, error) {
	blockchainDTO := &blockchainConfigDTO{}

	url := fmt.Sprintf(configRoute, height)

	resp, err := ref.client.doNewRequest(ctx, http.MethodGet, url, nil, blockchainDTO)

	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return nil, err
	}

	return blockchainDTO.toStruct()
}

func (ref *NetworkService) GetNetworkConfig(ctx context.Context) (*BlockchainConfig, error) {
	height, err := ref.BlockchainService.GetBlockchainHeight(ctx)
	if err != nil {
		return nil, err
	}

	return ref.GetNetworkConfigAtHeight(ctx, height)
}

func (ref *NetworkService) GetNetworkVersionAtHeight(ctx context.Context, height Height) (*NetworkVersion, error) {
	netDTO := &networkVersionDTO{}

	url := fmt.Sprintf(upgradeRoute, height)

	resp, err := ref.client.doNewRequest(ctx, http.MethodGet, url, nil, netDTO)

	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, nil); err != nil {
		return nil, err
	}

	return netDTO.toStruct(), nil
}

func (ref *NetworkService) GetNetworkVersion(ctx context.Context) (*NetworkVersion, error) {
	height, err := ref.BlockchainService.GetBlockchainHeight(ctx)
	if err != nil {
		return nil, err
	}

	return ref.GetNetworkVersionAtHeight(ctx, height)
}
