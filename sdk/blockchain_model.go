// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"github.com/proximax-storage/go-xpx-utils/str"
)

type BlockInfo struct {
	NetworkType
	BlockHash              Hash
	GenerationHash         Hash
	TotalFee               Amount
	NumTransactions        uint64
	Signature              Signature
	Signer                 *PublicAccount
	Version                uint8
	Type                   uint64
	Height                 Height
	Timestamp              *Timestamp
	Difficulty             Difficulty
	FeeMultiplier          uint32
	PreviousBlockHash      Hash
	BlockTransactionsHash  Hash
	BlockReceiptsHash      Hash
	StateHash              Hash
	Beneficiary            *PublicAccount
	FeeInterest            uint32
	FeeInterestDenominator uint32
}

func (b *BlockInfo) String() string {
	return str.StructToString(
		"BlockInfo",
		str.NewField("NetworkType", str.IntPattern, b.NetworkType),
		str.NewField("BlockHash", str.StringPattern, b.BlockHash),
		str.NewField("GenerationHash", str.StringPattern, b.GenerationHash),
		str.NewField("TotalFee", str.StringPattern, b.TotalFee),
		str.NewField("NumTransactions", str.IntPattern, b.NumTransactions),
		str.NewField("Signature", str.StringPattern, b.Signature),
		str.NewField("Signer", str.StringPattern, b.Signer),
		str.NewField("Version", str.IntPattern, b.Version),
		str.NewField("Type", str.IntPattern, b.Type),
		str.NewField("Height", str.StringPattern, b.Height),
		str.NewField("Timestamp", str.StringPattern, b.Timestamp),
		str.NewField("Difficulty", str.StringPattern, b.Difficulty),
		str.NewField("FeeMultiplier", str.IntPattern, b.FeeMultiplier),
		str.NewField("PreviousBlockHash", str.StringPattern, b.PreviousBlockHash),
		str.NewField("BlockTransactionsHash", str.StringPattern, b.BlockTransactionsHash),
		str.NewField("BlockReceiptsHash", str.StringPattern, b.BlockReceiptsHash),
		str.NewField("StateHash", str.StringPattern, b.StateHash),
		str.NewField("Beneficiary", str.StringPattern, b.Beneficiary),
		str.NewField("FeeInterest", str.IntPattern, b.FeeInterest),
		str.NewField("FeeInterestDenominator", str.IntPattern, b.FeeInterestDenominator),
	)
}

type BlockchainStorageInfo struct {
	NumBlocks       int `json:"numBlocks"`
	NumTransactions int `json:"numTransactions"`
	NumAccounts     int `json:"numAccounts"`
}

func (b *BlockchainStorageInfo) String() string {
	return str.StructToString(
		"BlockchainStorageInfo",
		str.NewField("NumBlocks", str.IntPattern, b.NumBlocks),
		str.NewField("NumTransactions", str.IntPattern, b.NumTransactions),
		str.NewField("NumAccounts", str.IntPattern, b.NumAccounts),
	)
}
