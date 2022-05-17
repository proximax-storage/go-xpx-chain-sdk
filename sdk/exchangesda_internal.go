// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type sdaOfferBalanceDTO struct {
	MosaicIdGive      uint64DTO `json:"mosaicIdGive"`
	MosaicIdGet       uint64DTO `json:"mosaicIdGet"`
	CurrentAmountGive uint64DTO `json:"currentAmountGive"`
	CurrentAmountGet  uint64DTO `json:"currentAmountGet"`
	InitialAmountGive uint64DTO `json:"initialAmountGive"`
	InitialAmountGet  uint64DTO `json:"initialAmountGet"`
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
	owner, err := NewAccountFromPublicKey(ref.ExchangeSda.Owner, networkType)
	if err != nil {
		return nil, err
	}

	for _, dto := range ref.ExchangeSda.SdaOfferBalances {
		dto.Owner = owner.PublicKey
	}

	offersMap := make(map[MosaicId]map[MosaicId]*SdaOfferBalance)

	offers, err := ref.ExchangeSda.SdaOfferBalances.toStruct(networkType)
	for _, offer := range offers {
		mosaicIdGive, err := NewMosaicId(offer.MosaicGive.AssetId.Id())
		if err != nil {
			return nil, err
		}

		mosaicIdGet, err := NewMosaicId(offer.MosaicGet.AssetId.Id())
		if err != nil {
			return nil, err
		}

		offersMap[*mosaicIdGive][*mosaicIdGet] = offer
	}

	return &UserSdaExchangeInfo{
		Owner:     owner,
		SdaOffers: offersMap,
	}, nil
}
