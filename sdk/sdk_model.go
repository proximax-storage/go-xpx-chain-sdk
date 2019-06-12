// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
	"fmt"
	"time"
)

type BlockchainInt64 interface {
	ToArray() [2]uint32
	Bytes() []byte
}

type BlockchainIdType uint8

type BlockchainId interface {
	BlockchainInt64
	fmt.Stringer
	Type() BlockchainIdType
	Id() uint64
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

type BaseInt64 uint64

func (m BaseInt64) String() string {
	return fmt.Sprintf("%d", m.Int64())
}

func (m BaseInt64) Int64() int64 {
	return int64(m)
}

func (m BaseInt64) Uint64() uint64 {
	return uint64(m)
}

func (m BaseInt64) ToArray() [2]uint32 {
	return uint64ToArray(m.Uint64())
}

func (m BaseInt64) Bytes() []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, m.Uint64())
	return bytes
}

type Amount struct {
	BaseInt64
}

func NewAmount(id uint64) *Amount {
	amount := Amount{BaseInt64(id)}
	return &amount
}

type Height struct {
	BaseInt64
}

func NewHeight(id uint64) *Height {
	height := Height{BaseInt64(id)}
	return &height
}

type Duration struct {
	BaseInt64
}

func NewDuration(id int64) *Duration {
	duration := Duration{BaseInt64(id)}
	return &duration
}

type Difficulty struct {
	BaseInt64
}

func NewDifficulty(id uint64) *Difficulty {
	difficulty := Difficulty{BaseInt64(id)}
	return &difficulty
}

type ChainScore [2]uint64

func (m *ChainScore) String() string {
	return fmt.Sprintf("[ %d, %d ]", m[0], m[1])
}

func NewChainScore(scoreLow uint64, scoreHigh uint64) *ChainScore {
	chainScore := ChainScore([2]uint64{scoreLow, scoreHigh})
	return &chainScore
}

const TimestampNemesisBlockMilliseconds int64 = 1459468800 * 1000

type BlockchainTimestamp struct {
	BaseInt64
}

func NewBlockchainTimestamp(milliseconds int64) *BlockchainTimestamp {
	timestamp := BlockchainTimestamp{BaseInt64(milliseconds)}
	return &timestamp
}

func (t *BlockchainTimestamp) ToTimestamp() *Timestamp {
	return NewTimestamp(int64(t.BaseInt64) + TimestampNemesisBlockMilliseconds)
}

type Timestamp struct {
	time.Time
}

func NewTimestamp(milliseconds int64) *Timestamp {
	return &Timestamp{time.Unix(0, milliseconds*int64(time.Millisecond))}
}

func (t *Timestamp) ToBlockchainTimestamp() *BlockchainTimestamp {
	return NewBlockchainTimestamp((t.Time.UnixNano()/int64(time.Millisecond) - TimestampNemesisBlockMilliseconds))
}

type Deadline struct {
	Timestamp
}

// Create deadline based on current time of system.
func NewDeadline(delta time.Duration) *Deadline {
	return &Deadline{Timestamp{time.Now().Add(delta)}}
}

// Create deadline from blockchain timestamp.
func NewDeadlineFromBlockchainTimestamp(timestamp *BlockchainTimestamp) *Deadline {
	return &Deadline{*timestamp.ToTimestamp()}
}
