// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
	"fmt"
	"time"
)

type blockchainInt64 interface {
	toArray() [2]uint32
	toLittleEndian() []byte
}

type BlockchainIdType uint8

// BlockchainIdType enums
const (
	NamespaceBlockchainIdType BlockchainIdType = iota
	MosaicBlockchainIdType
)

type BlockchainId interface {
	blockchainInt64
	fmt.Stringer
	Type() BlockchainIdType
	Id() uint64
	Equals(BlockchainId) bool
}

type TransactionOrder string

const (
	TRANSACTION_ORDER_ASC  TransactionOrder = "id"
	TRANSACTION_ORDER_DESC TransactionOrder = "-id"
)

type AccountTransactionsOption struct {
	PageSize int              `url:"pageSize,omitempty"`
	Id       string           `url:"id,omitempty"`
	Ordering TransactionOrder `url:"ordering,omitempty"`
}

type baseInt64 int64

func (m baseInt64) String() string {
	return fmt.Sprintf("%d", m)
}

func (m baseInt64) toArray() [2]uint32 {
	return uint64ToArray(uint64(m))
}

func (m baseInt64) toLittleEndian() []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(m))
	return bytes
}

type Amount = baseInt64
type Height = baseInt64
type Duration = baseInt64
type Difficulty = baseInt64

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
	return NewBlockchainTimestamp((t.Time.UnixNano()/int64(time.Millisecond) - TimestampNemesisBlockMilliseconds))
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
