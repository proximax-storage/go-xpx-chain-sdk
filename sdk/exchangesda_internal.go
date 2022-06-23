// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import "fmt"

type sdaOfferBalanceDTO struct {
	MosaicIdGive      uint64DTO `json:"mosaicIdGive"`
	MosaicIdGet       uint64DTO `json:"mosaicIdGet"`
	CurrentAmountGive uint64DTO `json:"currentMosaicGiveAmount"`
	CurrentAmountGet  uint64DTO `json:"currentMosaicGetAmount"`
	InitialAmountGive uint64DTO `json:"initialMosaicGiveAmount"`
	InitialAmountGet  uint64DTO `json:"initialMosaicGetAmount"`
	Owner             string    `json:"owner"`
	Deadline          uint64DTO `json:"deadline"`
}

type sdaOfferBalanceDTOs []*sdaOfferBalanceDTO

func (ref *sdaOfferBalanceDTOs) toStruct(networkType NetworkType) ([]*SdaOfferBalance, error) {
	var (
		dtos   = *ref
		offers = make([]*SdaOfferBalance, len(*ref))
	)

	for i, dto := range dtos {
		owner, err := NewAccountFromPublicKey(dto.Owner, networkType)
		if err != nil {
			return nil, err
		}

		mosaicIdGive, err := NewMosaicId(dto.MosaicIdGive.toUint64())
		if err != nil {
			return nil, err
		}

		mosaicIdGet, err := NewMosaicId(dto.MosaicIdGet.toUint64())
		if err != nil {
			return nil, err
		}

		offers[i] = &SdaOfferBalance{
			Owner:             owner,
			MosaicGive:        newMosaicPanic(mosaicIdGive, dto.CurrentAmountGive.toStruct()),
			MosaicGet:         newMosaicPanic(mosaicIdGet, dto.CurrentAmountGet.toStruct()),
			InitialAmountGive: dto.InitialAmountGive.toStruct(),
			InitialAmountGet:  dto.InitialAmountGet.toStruct(),
			Deadline:          dto.Deadline.toStruct(),
		}
	}

	return offers, nil
}

type sdaExchangeDTO struct {
	ExchangeSda struct {
		Owner            string              `json:"owner"`
		SdaOfferBalances sdaOfferBalanceDTOs `json:"sdaOfferBalances"`
	} `json:"exchangesda"`
}

func (ref *sdaExchangeDTO) toStruct(networkType NetworkType) (*UserSdaExchangeInfo, error) {
	sdaExchangeInfo := UserSdaExchangeInfo{}

	owner, err := NewAccountFromPublicKey(ref.ExchangeSda.Owner, networkType)
	if err != nil {
		return nil, err
	}

	for _, dto := range ref.ExchangeSda.SdaOfferBalances {
		dto.Owner = owner.PublicKey
	}

	offers, err := ref.ExchangeSda.SdaOfferBalances.toStruct(networkType)
	if err != nil {
		return nil, fmt.Errorf("sdk.sdaExchangeDTO.toStruct ExchangeSda.SdaOfferBalances.toStruct: %v", err)
	}

	sdaExchangeInfo.Owner = owner
	sdaExchangeInfo.SdaOfferBalances = offers

	return &sdaExchangeInfo, nil
}
