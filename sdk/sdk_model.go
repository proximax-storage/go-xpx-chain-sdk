// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
)

type uint64DTO [2]uint32

func (dto uint64DTO) toBigInt() *big.Int {
	if dto[0] == 0 && dto[1] == 0 {
		return &big.Int{}
	}
	var int big.Int
	b := make([]byte, len(dto)*4)
	binary.BigEndian.PutUint32(b[:len(dto)*2], dto[1])
	binary.BigEndian.PutUint32(b[len(dto)*2:], dto[0])
	int.SetBytes(b)
	return &int
}

type uint64DTOs []*uint64DTO

func (dto uint64DTOs) toBigInts() []*big.Int {
	result := make([]*big.Int, len(dto))

	for i, b := range dto {
		result[i] = b.toBigInt()
	}

	return result
}

func intToHex(u uint32) string {
	return fmt.Sprintf("%08x", u)
}

// analog JAVA Uint64.bigIntegerToHex
func bigIntegerToHex(id *big.Int) string {
	u := fromBigInt(id)
	return strings.ToUpper(intToHex(u[1]) + intToHex(u[0]))
}

// TODO why it is exported?
func fromBigInt(int *big.Int) []uint32 {
	if int == nil {
		return []uint32{0, 0}
	}

	var u64 = uint64(int.Int64())
	l := uint32(u64 & 0xFFFFFFFF)
	r := uint32(u64 >> 32)
	return []uint32{l, r}
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
