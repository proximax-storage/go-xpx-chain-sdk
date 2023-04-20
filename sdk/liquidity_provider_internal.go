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
	LiquidityProvider struct {
		MosaicId           uint64DTO           `json:"mosaicId"`
		Provider           string              `json:"providerKey"`
		Owner              string              `json:"owner"`
		AdditionallyMinted uint64DTO           `json:"additionallyMinted"`
		SlashingAccount    string              `json:"slashingAccount"`
		SlashingPeriod     uint32              `json:"slashingPeriod"`
		WindowSize         uint16              `json:"windowSize"`
		CreationHeight     uint64DTO           `json:"creationHeight"`
		Alpha              uint32              `json:"alpha"`
		Beta               uint32              `json:"beta"`
		TurnoverHistory    turnoverHistoryDTOs `json:"turnoverHistory"`
		RecentTurnover     turnoverDTO         `json:"recentTurnover"`
	} `json:"liquidityProvider"`
}

func (ref *liquidityProviderDTO) toStruct(networkType NetworkType) (*LiquidityProvider, error) {
	mosaicId, err := NewMosaicId(ref.LiquidityProvider.MosaicId.toUint64())
	if err != nil {
		return nil, err
	}

	provider, err := NewAccountFromPublicKey(ref.LiquidityProvider.Provider, networkType)
	if err != nil {
		return nil, err
	}

	owner, err := NewAccountFromPublicKey(ref.LiquidityProvider.Owner, networkType)
	if err != nil {
		return nil, err
	}

	slashingAccount, err := NewAccountFromPublicKey(ref.LiquidityProvider.SlashingAccount, networkType)
	if err != nil {
		return nil, err
	}

	turnoverHistory, err := ref.LiquidityProvider.TurnoverHistory.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	recentTurnover, err := ref.LiquidityProvider.RecentTurnover.toStruct(networkType)
	if err != nil {
		return nil, err
	}

	return &LiquidityProvider{
		MosaicId:           mosaicId,
		ProviderKey:        provider,
		Owner:              owner,
		AdditionallyMinted: ref.LiquidityProvider.AdditionallyMinted.toStruct(),
		SlashingAccount:    slashingAccount,
		SlashingPeriod:     ref.LiquidityProvider.SlashingPeriod,
		WindowSize:         ref.LiquidityProvider.WindowSize,
		CreationHeight:     ref.LiquidityProvider.CreationHeight.toStruct(),
		Alpha:              ref.LiquidityProvider.Alpha,
		Beta:               ref.LiquidityProvider.Beta,
		TurnoverHistory:    turnoverHistory,
		RecentTurnover:     recentTurnover,
	}, nil
}

type liquidityProvidersPageDTO struct {
	LiquidityProviders []liquidityProviderDTO `json:"data"`

	Pagination struct {
		TotalEntries uint64 `json:"totalEntries"`
		PageNumber   uint64 `json:"pageNumber"`
		PageSize     uint64 `json:"pageSize"`
		TotalPages   uint64 `json:"totalPages"`
	} `json:"pagination"`
}

func (lpp *liquidityProvidersPageDTO) toStruct(networkType NetworkType) (*LiquidityProviderPage, error) {
	page := &LiquidityProviderPage{
		LiquidityProviders: make([]*LiquidityProvider, len(lpp.LiquidityProviders)),
		Pagination: Pagination{
			TotalEntries: lpp.Pagination.TotalEntries,
			PageNumber:   lpp.Pagination.PageNumber,
			PageSize:     lpp.Pagination.PageSize,
			TotalPages:   lpp.Pagination.TotalPages,
		},
	}

	var err error
	for i, lp := range lpp.LiquidityProviders {
		page.LiquidityProviders[i], err = lp.toStruct(networkType)
		if err != nil {
			return nil, err
		}
	}

	return page, nil
}

type createLiquidityProviderTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ProviderMosaicId      mosaicIdDTO `json:"providerMosaicId"`
		CurrencyDeposit       uint64DTO   `json:"currencyDeposit"`
		InitialMosaicsMinting uint64DTO   `json:"initialMosaicsMinting"`
		SlashingPeriod        uint32      `json:"slashingPeriod"`
		WindowSize            uint16      `json:"windowSize"`
		SlashingAccount       string      `json:"slashingAccount"`
		Alpha                 uint32      `json:"alpha"`
		Beta                  uint32      `json:"beta"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *createLiquidityProviderTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	providerMosaicId, err := dto.Tx.ProviderMosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	slashingAccount, err := NewAccountFromPublicKey(dto.Tx.SlashingAccount, atx.NetworkType)
	if err != nil {
		return nil, err
	}

	return &CreateLiquidityProviderTransaction{
		*atx,
		providerMosaicId,
		dto.Tx.CurrencyDeposit.toStruct(),
		dto.Tx.InitialMosaicsMinting.toStruct(),
		dto.Tx.SlashingPeriod,
		dto.Tx.WindowSize,
		slashingAccount,
		dto.Tx.Alpha,
		dto.Tx.Beta,
	}, nil
}

type manualRateChangeTransactionDTO struct {
	Tx struct {
		abstractTransactionDTO
		ProviderMosaicId        mosaicIdDTO `json:"providerMosaicId"`
		CurrencyBalanceIncrease uint8       `json:"currencyBalanceIncrease,omitempty"`
		CurrencyBalanceChange   uint64DTO   `json:"currencyBalanceChange"`
		MosaicBalanceIncrease   uint8       `json:"mosaicBalanceIncrease,omitempty"`
		MosaicBalanceChange     uint64DTO   `json:"mosaicBalanceChange"`
	} `json:"transaction"`
	TDto transactionInfoDTO `json:"meta"`
}

func (dto *manualRateChangeTransactionDTO) toStruct(*Hash) (Transaction, error) {
	info, err := dto.TDto.toStruct()
	if err != nil {
		return nil, err
	}

	atx, err := dto.Tx.abstractTransactionDTO.toStruct(info)
	if err != nil {
		return nil, err
	}

	providerMosaicId, err := dto.Tx.ProviderMosaicId.toStruct()
	if err != nil {
		return nil, err
	}

	return &ManualRateChangeTransaction{
		*atx,
		providerMosaicId,
		true,
		dto.Tx.CurrencyBalanceChange.toStruct(),
		true,
		dto.Tx.MosaicBalanceChange.toStruct(),
	}, nil
}
