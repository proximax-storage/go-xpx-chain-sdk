// Copyright 2022 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/hex"
	"fmt"

	"github.com/google/flatbuffers/go"
	"github.com/proximax-storage/go-xpx-chain-sdk/transactions"
)

func NewCreateLiquidityProviderTransaction(
	deadline *Deadline,
	providerMosaicId *MosaicId,
	currencyDeposit Amount,
	initialMosaicsMinting Amount,
	slashingPeriod uint32,
	windowSize uint16,
	slashingAccount *PublicAccount,
	alpha uint32,
	beta uint32,
	networkType NetworkType,
) (*CreateLiquidityProviderTransaction, error) {
	tx := CreateLiquidityProviderTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     CreateLiquidityProviderVersion,
			Deadline:    deadline,
			Type:        CreateLiquidityProvider,
			NetworkType: networkType,
		},
		ProviderMosaicId:      providerMosaicId,
		CurrencyDeposit:       currencyDeposit,
		InitialMosaicsMinting: initialMosaicsMinting,
		SlashingPeriod:        slashingPeriod,
		WindowSize:            windowSize,
		SlashingAccount:       slashingAccount,
		Alpha:                 alpha,
		Beta:                  beta,
	}

	return &tx, nil
}

func (tx *CreateLiquidityProviderTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *CreateLiquidityProviderTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"ProviderMosaicId": %s,
			"CurrencyDeposit": %s,
			"InitialMosaicsMinting": %s,
			"SlashingPeriod": %d,
			"WindowSize": %d,
			"SlashingAccount": %s,
			"Alpha": %d,
			"Beta": %d,
		`,
		tx.AbstractTransaction.String(),
		tx.ProviderMosaicId.String(),
		tx.CurrencyDeposit.String(),
		tx.InitialMosaicsMinting.String(),
		tx.SlashingPeriod,
		tx.WindowSize,
		tx.SlashingAccount.String(),
		tx.Alpha,
		tx.Beta,
	)
}

func (tx *CreateLiquidityProviderTransaction) Size() int {
	return CreateLiquidityProviderHeaderSize
}

func (tx *CreateLiquidityProviderTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	slashingAccountB, err := hex.DecodeString(tx.SlashingAccount.PublicKey)
	if err != nil {
		return nil, err
	}

	slashingAccountV := transactions.TransactionBufferCreateByteVector(builder, slashingAccountB)
	providerMosaicIdV := transactions.TransactionBufferCreateUint32Vector(builder, tx.ProviderMosaicId.toArray())
	currencyDepositV := transactions.TransactionBufferCreateUint32Vector(builder, tx.CurrencyDeposit.toArray())
	initialMosaicsMintingV := transactions.TransactionBufferCreateUint32Vector(builder, tx.InitialMosaicsMinting.toArray())

	transactions.CreateLiquidityProviderBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.CreateLiquidityProviderBufferAddProviderMosaicId(builder, providerMosaicIdV)
	transactions.CreateLiquidityProviderBufferAddCurrencyDeposit(builder, currencyDepositV)
	transactions.CreateLiquidityProviderBufferAddInitialMosaicsMinting(builder, initialMosaicsMintingV)
	transactions.CreateLiquidityProviderBufferAddSlashingPeriod(builder, tx.SlashingPeriod)
	transactions.CreateLiquidityProviderBufferAddWindowSize(builder, tx.WindowSize)
	transactions.CreateLiquidityProviderBufferAddSlashingAccount(builder, slashingAccountV)
	transactions.CreateLiquidityProviderBufferAddAlpha(builder, tx.Alpha)
	transactions.CreateLiquidityProviderBufferAddBeta(builder, tx.Beta)

	t := transactions.CreateLiquidityProviderBufferEnd(builder)
	builder.Finish(t)

	return createLiquidityProviderTransactionSchema().serialize(builder.FinishedBytes()), nil
}

func NewManualRateChangeTransaction(
	deadline *Deadline,
	providerMosaicId *MosaicId,
	currencyBalanceIncrease bool,
	currencyBalanceChange Amount,
	mosaicBalanceIncrease bool,
	mosaicBalanceChange Amount,
	networkType NetworkType,
) (*ManualRateChangeTransaction, error) {
	tx := ManualRateChangeTransaction{
		AbstractTransaction: AbstractTransaction{
			Version:     ManualRateChangeVersion,
			Deadline:    deadline,
			Type:        ManualRateChange,
			NetworkType: networkType,
		},
		ProviderMosaicId:        providerMosaicId,
		CurrencyBalanceIncrease: currencyBalanceIncrease,
		CurrencyBalanceChange:   currencyBalanceChange,
		MosaicBalanceIncrease:   mosaicBalanceIncrease,
		MosaicBalanceChange:     mosaicBalanceChange,
	}

	return &tx, nil
}

func (tx *ManualRateChangeTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &tx.AbstractTransaction
}

func (tx *ManualRateChangeTransaction) String() string {
	return fmt.Sprintf(
		`
			"AbstractTransaction": %s,
			"ProviderMosaicId": %s,
			"CurrencyBalanceIncrease": %t,
			"CurrencyBalanceChange": %s,
			"MosaicBalanceIncrease": %t,
			"MosaicBalanceChange": %s,
		`,
		tx.AbstractTransaction.String(),
		tx.ProviderMosaicId.String(),
		tx.CurrencyBalanceIncrease,
		tx.CurrencyBalanceChange.String(),
		tx.MosaicBalanceIncrease,
		tx.MosaicBalanceChange.String(),
	)
}

func (tx *ManualRateChangeTransaction) Size() int {
	return ManualRateChangeHeaderSize
}

func (tx *ManualRateChangeTransaction) Bytes() ([]byte, error) {
	builder := flatbuffers.NewBuilder(0)

	v, signatureV, signerV, deadlineV, fV, err := tx.AbstractTransaction.generateVectors(builder)
	if err != nil {
		return nil, err
	}

	providerMosaicIdV := transactions.TransactionBufferCreateUint32Vector(builder, tx.ProviderMosaicId.toArray())
	currencyBalanceChangeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.CurrencyBalanceChange.toArray())
	mosaicBalanceChangeV := transactions.TransactionBufferCreateUint32Vector(builder, tx.MosaicBalanceChange.toArray())

	transactions.ManualRateChangeBufferStart(builder)
	transactions.TransactionBufferAddSize(builder, tx.Size())
	tx.AbstractTransaction.buildVectors(builder, v, signatureV, signerV, deadlineV, fV)

	transactions.ManualRateChangeBufferAddProviderMosaicId(builder, providerMosaicIdV)
	transactions.ManualRateChangeBufferAddCurrencyBalanceIncrease(builder, tx.CurrencyBalanceIncrease)
	transactions.ManualRateChangeBufferAddCurrencyBalanceChange(builder, currencyBalanceChangeV)
	transactions.ManualRateChangeBufferAddMosaicBalanceIncrease(builder, tx.MosaicBalanceIncrease)
	transactions.ManualRateChangeBufferAddMosaicBalanceChange(builder, mosaicBalanceChangeV)

	t := transactions.ManualRateChangeBufferEnd(builder)
	builder.Finish(t)

	return manualRateChangeTransactionSchema().serialize(builder.FinishedBytes()), nil
}
