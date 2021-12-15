// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type Signature [64]byte

func (s Signature) String() string {
	return hex.EncodeToString(s[:])
}

type Hash [32]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func (h Hash) Empty() bool {
	return h.Equal(&Hash{})
}

func (h Hash) Equal(other *Hash) bool {
	return bytes.Compare(h[:], other[:]) == 0
}

func (h Hash) Xor(other *Hash) *Hash {
	temp := Hash{}

	for i := range other {
		temp[i] = other[i] ^ h[i]
	}

	return &temp
}

type blockchainInt64 interface {
	toArray() [2]uint32
	toLittleEndian() []byte
}

type AssetIdType uint8

// AssetIdType enums
const (
	NamespaceAssetIdType AssetIdType = iota
	MosaicAssetIdType
)

type AssetId interface {
	blockchainInt64
	fmt.Stringer
	Type() AssetIdType
	Id() uint64
	Equals(AssetId) (bool, error)
}

func NewAssetIdFromId(id uint64) (AssetId, error) {
	if hasBits(id, NamespaceBit) {
		return NewNamespaceId(id)
	} else {
		return NewMosaicId(id)
	}
}

type TransactionOrder string

const (
	TRANSACTION_ORDER_ASC  TransactionOrder = "id"
	TRANSACTION_ORDER_DESC TransactionOrder = "-id"
)

type SortOptions struct {
	SortField string
	Direction SortDirection
}

type SortDirection string

const (
	ASC  SortDirection = "asc"
	DESC SortDirection = "desc"
)

func (sD SortDirection) String() string {
	return string(sD)
}

type PaginationOrderingOptions struct {
	PageSize      uint64 `url:"pageSize,omitempty"`
	PageNumber    uint64 `url:"pageNumber,omitempty"`
	Offset        string `url:"offset,omitempty"`
	SortField     string `url:"sortField,omitempty"`
	SortDirection string `url:"order,omitempty"`
}

type baseInt64 int64

func (m baseInt64) String() string {
	return fmt.Sprintf("%d", m)
}

func (m baseInt64) ToHexString() string {
	return uint64ToHex(uint64(m))
}

func (m baseInt64) toArray() [2]uint32 {
	return uint64ToArray(uint64(m))
}

func (m baseInt64) toLittleEndian() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(m))
	return b
}

type Amount = baseInt64
type Height = baseInt64
type Duration = baseInt64
type Difficulty = baseInt64
type StorageSize = baseInt64
type ScopedMetadataKey = baseInt64

type BlockChainVersion uint64

func NewBlockChainVersion(major uint16, minor uint16, revision uint16, build uint16) BlockChainVersion {
	version := BlockChainVersion(0)
	version |= BlockChainVersion(major) << 48
	version |= BlockChainVersion(minor) << 32
	version |= BlockChainVersion(revision) << 16
	version |= BlockChainVersion(build) << 0
	return version
}

func (m BlockChainVersion) String() string {
	getTwoBytesByShift := func(number BlockChainVersion, shift uint) uint16 {
		return uint16(number>>shift) & 0xFF
	}

	return fmt.Sprintf(
		"%d [Major %d, Minor %d, Revision %d, Build %d]",
		m,
		getTwoBytesByShift(m, 48),
		getTwoBytesByShift(m, 32),
		getTwoBytesByShift(m, 16),
		getTwoBytesByShift(m, 0),
	)
}

func (m BlockChainVersion) toArray() [2]uint32 {
	return uint64ToArray(uint64(m))
}

func (m BlockChainVersion) toLittleEndian() []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(m))
	return bytes
}

type ChainScore [2]uint64

func (m *ChainScore) String() string {
	return fmt.Sprintf("[ %d, %d ]", m[0], m[1])
}

// returns new ChainScore from passed low and high score
func NewChainScore(scoreLow uint64, scoreHigh uint64) *ChainScore {
	chainScore := ChainScore([2]uint64{scoreLow, scoreHigh})
	return &chainScore
}

const TimestampNemesisBlockMilliseconds int64 = 1459468800 * 1000

type BlockchainTimestamp struct {
	baseInt64
}

// returns new BlockchainTimestamp from passed milliseconds value
func NewBlockchainTimestamp(milliseconds int64) *BlockchainTimestamp {
	timestamp := BlockchainTimestamp{baseInt64(milliseconds)}
	return &timestamp
}

func (t *BlockchainTimestamp) ToTimestamp() *Timestamp {
	return NewTimestamp(int64(t.baseInt64) + TimestampNemesisBlockMilliseconds)
}

type Timestamp struct {
	time.Time
}

// returns new Timestamp from passed milliseconds value
func NewTimestamp(milliseconds int64) *Timestamp {
	return &Timestamp{time.Unix(0, milliseconds*int64(time.Millisecond))}
}

func (t *Timestamp) ToBlockchainTimestamp() *BlockchainTimestamp {
	return NewBlockchainTimestamp(t.Time.UnixNano()/int64(time.Millisecond) - TimestampNemesisBlockMilliseconds)
}

type Deadline struct {
	Timestamp
}

// returns new Deadline from passed duration
func NewDeadline(delta time.Duration) *Deadline {
	return &Deadline{Timestamp{time.Now().Add(delta)}}
}

// returns new Deadline from passed BlockchainTimestamp
func NewDeadlineFromBlockchainTimestamp(timestamp *BlockchainTimestamp) *Deadline {
	return &Deadline{*timestamp.ToTimestamp()}
}
