// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type rateDTO struct {
	CurrencyAmount uint64DTO `json:"currencyAmount"`
	MosaicAmount   uint64DTO `json:"mosaicAmount"`
}

func (ref *rateDTO) toStruct(networkType NetworkType) (*Rate, error) {
	return &Rate{
		CurrencyAmount: ref.CurrencyAmount.toStruct(),
		MosaicAmount:   ref.MosaicAmount.toStruct(),
	}, nil
}

type turnoverDTO struct {
	Rate     rateDTO   `json:"rate"`
	Turnover uint64DTO `json:"turnover"`
}

func (ref *turnoverDTO) toStruct(networkType NetworkType) (*Turnover, error) {
	rate, err := ref.Rate.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &Turnover{
		Rate:     rate,
		Turnover: ref.Turnover.toStruct(),
	}, nil
}

type turnoverHistoryDTOs []*turnoverDTO

func (ref *turnoverHistoryDTOs) toStruct(networkType NetworkType) ([]*Turnover, error) {
	var (
		dtos      = *ref
		histories = make([]*Turnover, 0, len(dtos))
	)

	var err error
	for i, dto := range dtos {
		histories[i], err = dto.toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}

	return histories, nil
}

type liquidityProviderDTO struct {
	MosaicId           uint64DTO           `json:"mosaicId"`
	Provider           string              `json:"providerKey"`
	Owner              string              `json:"owner"`
	AdditionallyMinted uint64DTO           `json:"additionallyMinted"`
	SlashingAccount    string              `json:"slashingAccount"`
	SlashingPeriod     uint32              `json:"slashingPeriod"`
	WindowSize         uint32              `json:"windowSize"`
	CreationHeight     uint64DTO           `json:"creationHeight"`
	Alpha              uint32              `json:"alpha"`
	Beta               uint32              `json:"beta"`
	TurnoverHistory    turnoverHistoryDTOs `json:"turnoverHistory"`
	RecentTurnover     turnoverDTO         `json:"recentTurnover"`
}

func (ref *liquidityProviderDTO) toStruct(networkType NetworkType) (*LiquidityProvider, error) {
	mosaicId, err := NewMosaicId(ref.MosaicId.toUint64())
	if err != nil {
		return nil, err
	}

	provider, err := NewAccountFromPublicKey(ref.Provider, networkType)
	if err != nil {
		return nil, err
	}

	owner, err := NewAccountFromPublicKey(ref.Owner, networkType)
	if err != nil {
		return nil, err
	}

	slashingAccount, err := NewAccountFromPublicKey(ref.SlashingAccount, networkType)
	if err != nil {
		return nil, err
	}

	turnoverHistory, err := ref.TurnoverHistory.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	recentTurnover, err := ref.RecentTurnover.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &LiquidityProvider{
		MosaicId:           mosaicId,
		ProviderKey:        provider,
		Owner:              owner,
		AdditionallyMinted: ref.AdditionallyMinted.toStruct(),
		SlashingAccount:    slashingAccount,
		SlashingPeriod:     ref.SlashingPeriod,
		WindowSize:         ref.WindowSize,
		CreationHeight:     ref.CreationHeight.toStruct(),
		Alpha:              ref.Alpha,
		Beta:               ref.Beta,
		TurnoverHistory:    turnoverHistory,
		RecentTurnover:     recentTurnover,
	}, nil
}
