// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
	"github.com/proximax-storage/go-xpx-catapult-sdk/utils"
	"github.com/proximax-storage/go-xpx-utils/str"
	"math/big"
)

// MosaicId
type MosaicId big.Int

func (m *MosaicId) String() string {
	return (*big.Int)(m).String()
}

func (m *MosaicId) Equals(id *MosaicId) bool {
	return (*big.Int)(m).Uint64() == (*big.Int)(id).Uint64()
}

// returns mosaic id for passed nonce and public key of owner
func NewMosaicIdFromNonceAndOwner(nonce uint32, ownerPublicKey string) (*MosaicId, error) {
	if len(ownerPublicKey) != 64 {
		return nil, ErrInvalidOwnerPublicKey
	}

	id, err := generateMosaicId(nonce, ownerPublicKey)

	return bigIntToMosaicId(id), err
}

// returns mosaic id corresponding passed big int
func NewMosaicId(id *big.Int) (*MosaicId, error) {
	if id == nil {
		return nil, ErrNilMosaicId
	}

	return bigIntToMosaicId(id), nil
}

func (m *MosaicId) toHexString() string {
	return BigIntegerToHex(mosaicIdToBigInt(m))
}

type Mosaic struct {
	MosaicId *MosaicId
	Amount   *big.Int
}

// returns a mosaic for passed mosaic id and amount
func NewMosaic(mosaicId *MosaicId, amount *big.Int) (*Mosaic, error) {
	if mosaicId == nil {
		return nil, ErrNilMosaicId
	}

	if amount == nil {
		return nil, ErrNilMosaicAmount
	}

	if utils.EqualsBigInts(amount, big.NewInt(0)) {
		return nil, ErrNilMosaicAmount
	}

	return &Mosaic{
		MosaicId: mosaicId,
		Amount:   amount,
	}, nil
}

func (m *Mosaic) String() string {
	return str.StructToString(
		"MosaicId",
		str.NewField("MosaicId", str.StringPattern, m.MosaicId),
		str.NewField("Amount", str.IntPattern, m.Amount),
	)
}

// MosaicInfo info structure contains its properties, the owner and the namespace to which it belongs to.
type MosaicInfo struct {
	MosaicId   *MosaicId
	Supply     *big.Int
	Height     *big.Int
	Owner      *PublicAccount
	Revision   uint32
	Properties *MosaicProperties
}

func (m *MosaicInfo) String() string {
	return str.StructToString(
		"MosaicInfo",
		str.NewField("MosaicId", str.StringPattern, m.MosaicId),
		str.NewField("Supply", str.StringPattern, m.Supply),
		str.NewField("Height", str.StringPattern, m.Height),
		str.NewField("Owner", str.StringPattern, m.Owner),
		str.NewField("Revision", str.IntPattern, m.Revision),
		str.NewField("Properties", str.StringPattern, m.Properties),
	)
}

// structure which includes several properties for defining mosaic
// `SupplyMutable` - is supply of defined mosaic can be changed in future
// `Transferable` - if this property is set to "false", only transfer transactions having the creator as sender or as recipient can transfer mosaics of that type. If set to "true" the mosaics can be transferred to and from arbitrary accounts
// `LevyMutable` - if this property is set to "true", whenever other users transact with your mosaic, owner gets a levy fee from them
// `Divisibility` - divisibility determines up to what decimal place the mosaic can be divided into
// `Duration` - duration in blocks mosaic will be available. After the renew mosaic is inactive and can be renewed
type MosaicProperties struct {
	SupplyMutable bool
	Transferable  bool
	LevyMutable   bool
	Divisibility  uint8
	Duration      *big.Int
}

func NewMosaicProperties(supplyMutable bool, transferable bool, levyMutable bool, divisibility uint8, duration *big.Int) *MosaicProperties {
	ref := &MosaicProperties{
		supplyMutable,
		transferable,
		levyMutable,
		divisibility,
		duration,
	}

	return ref
}

func (mp *MosaicProperties) String() string {
	return str.StructToString(
		"MosaicProperties",
		str.NewField("SupplyMutable", str.BooleanPattern, mp.SupplyMutable),
		str.NewField("Transferable", str.BooleanPattern, mp.Transferable),
		str.NewField("LevyMutable", str.BooleanPattern, mp.LevyMutable),
		str.NewField("Divisibility", str.IntPattern, mp.Divisibility),
		str.NewField("Duration", str.StringPattern, mp.Duration),
	)
}

type MosaicSupplyType uint8

const (
	Decrease MosaicSupplyType = iota
	Increase
)

func (tx MosaicSupplyType) String() string {
	return fmt.Sprintf("%d", tx)
}

// returns XEM mosaic with passed amount
func Xem(amount int64) *Mosaic {
	return &Mosaic{XemMosaicId, big.NewInt(amount)}
}

// returns XPX mosaic with passed amount
func Xpx(amount int64) *Mosaic {
	return &Mosaic{XpxMosaicId, big.NewInt(amount)}
}

// returns XEM with actual passed amount
func XemRelative(amount int64) *Mosaic {
	return Xem(big.NewInt(0).Mul(big.NewInt(1000000), big.NewInt(amount)).Int64())
}

// returns XPX with actual passed amount
func XpxRelative(amount int64) *Mosaic {
	return Xpx(big.NewInt(0).Mul(big.NewInt(1000000), big.NewInt(amount)).Int64())
}
