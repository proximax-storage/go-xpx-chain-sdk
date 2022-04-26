// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type Rate struct {
	CurrencyAmount Amount
	MosaicAmount   Amount
}

type Turnover struct {
	Rate     *Rate
	Turnover Amount
}

type LiquidityProvider struct {
	MosaicId           *MosaicId
	ProviderKey        *PublicAccount
	Owner              *PublicAccount
	AdditionallyMinted Amount
	SlashingAccount    *PublicAccount
	SlashingPeriod     uint32
	WindowSize         uint16
	CreationHeight     Height
	Alpha              uint32
	Beta               uint32
	TurnoverHistory    []*Turnover
	RecentTurnover     *Turnover
}

type LiquidityProviderPage struct {
	LiquidityProviders []*LiquidityProvider
	Pagination         Pagination
}

type LiquidityProviderPageOptions struct {
	PaginationOrderingOptions
	MosaicId        string `url:"mosaicId,omitempty"`
	SlashingAccount string `url:"slashingAccount,omitempty"`
	Owner           string `url:"owner,omitempty"`
}

type CreateLiquidityProviderTransaction struct {
	AbstractTransaction
	ProviderMosaicId      *MosaicId
	CurrencyDeposit       Amount
	InitialMosaicsMinting Amount
	SlashingPeriod        uint32
	WindowSize            uint16
	SlashingAccount       *PublicAccount
	Alpha                 uint32
	Beta                  uint32
}

type ManualRateChangeTransaction struct {
	AbstractTransaction
	ProviderMosaicId        *MosaicId
	CurrencyBalanceIncrease bool
	CurrencyBalanceChange   Amount
	MosaicBalanceIncrease   bool
	MosaicBalanceChange     Amount
}
