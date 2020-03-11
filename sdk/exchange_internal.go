// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type offerInfoDTO struct {
	MosaicId         uint64DTO `json:"mosaicId"`
	Amount           uint64DTO `json:"amount"`
	PriceNumerator   uint64DTO `json:"initialCost"`
	PriceDenominator uint64DTO `json:"initialAmount"`
	Deadline         uint64DTO `json:"deadline"`
	Owner            string    `json:"owner"`
	Type             OfferType `json:"type"`
}

type offerInfoDTOs []*offerInfoDTO

func (ref *offerInfoDTOs) toStruct(networkType NetworkType) ([]*OfferInfo, error) {
	var (
		dtos   = *ref
		offers = make([]*OfferInfo, len(*ref))
	)

	for i, dto := range dtos {
		owner, err := NewAccountFromPublicKey(dto.Owner, networkType)
		if err != nil {
			return nil, err
		}

		mosaicId, err := NewMosaicId(dto.MosaicId.toUint64())
		if err != nil {
			return nil, err
		}

		offers[i] = &OfferInfo{
			Type:             dto.Type,
			Owner:            owner,
			Mosaic:           newMosaicPanic(mosaicId, dto.Amount.toStruct()),
			PriceDenominator: dto.PriceDenominator.toStruct(),
			PriceNumerator:   dto.PriceNumerator.toStruct(),
			Deadline:         dto.Deadline.toStruct(),
		}
	}

	return offers, nil
}

type exchangeDTO struct {
	Exchange struct {
		Owner      string        `json:"owner"`
		BuyOffers  offerInfoDTOs `json:"buyOffers"`
		SellOffers offerInfoDTOs `json:"sellOffers"`
	} `json:"exchange"`
}

func (ref *exchangeDTO) toStruct(networkType NetworkType) (*UserExchangeInfo, error) {
	owner, err := NewAccountFromPublicKey(ref.Exchange.Owner, networkType)
	if err != nil {
		return nil, err
	}

	for _, dto := range ref.Exchange.BuyOffers {
		dto.Type = BuyOffer
		dto.Owner = owner.PublicKey
	}

	for _, dto := range ref.Exchange.SellOffers {
		dto.Type = SellOffer
		dto.Owner = owner.PublicKey
	}

	offersMap := make(map[OfferType]map[MosaicId]*OfferInfo)
	offersMap[BuyOffer] = make(map[MosaicId]*OfferInfo)
	offersMap[SellOffer] = make(map[MosaicId]*OfferInfo)

	offers, err := ref.Exchange.BuyOffers.toStruct(networkType)
	for _, offer := range offers {
		mosaicId, err := NewMosaicId(offer.Mosaic.AssetId.Id())
		if err != nil {
			return nil, err
		}

		offersMap[BuyOffer][*mosaicId] = offer
	}

	offers, err = ref.Exchange.SellOffers.toStruct(networkType)
	for _, offer := range offers {
		mosaicId, err := NewMosaicId(offer.Mosaic.AssetId.Id())
		if err != nil {
			return nil, err
		}

		offersMap[SellOffer][*mosaicId] = offer
	}

	return &UserExchangeInfo{
		Owner:  owner,
		Offers: offersMap,
	}, nil
}
