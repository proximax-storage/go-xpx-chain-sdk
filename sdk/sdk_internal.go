// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type uint64DTO [2]uint32

func (dto uint64DTO) toUint64() uint64 {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b[:4], dto[1])
	binary.BigEndian.PutUint32(b[4:], dto[0])
	return binary.BigEndian.Uint64(b)
}

func uint32ToHex(u uint32) string {
	return fmt.Sprintf("%08x", u)
}

func uint64ToHex(id uint64) string {
	u := uint64ToArray(id)
	return strings.ToUpper(uint32ToHex(u[1]) + uint32ToHex(u[0]))
}

func uint64ToArray(int uint64) [2]uint32 {
	l := uint32(int & 0xFFFFFFFF)
	r := uint32(int>>32) & 0xFFFFFFFF
	return [2]uint32{l, r}
}

type amountDTO uint64DTO

func (dto amountDTO) toStruct() *Amount {
	return NewAmount(uint64DTO(dto).toUint64())
}

type blockchainIdDTO uint64DTO

func (dto blockchainIdDTO) toStruct() (BlockchainId, error) {
	id := uint64DTO(dto).toUint64()

	if hasBits(id, NamespaceBit) {
		return (*namespaceIdDTO)(&dto).toStruct()
	} else {
		return (*mosaicIdDTO)(&dto).toStruct()
	}
}

type heightDTO uint64DTO

func (dto heightDTO) toStruct() *Height {
	return NewHeight(uint64DTO(dto).toUint64())
}

type durationDTO uint64DTO

func (dto durationDTO) toStruct() *Duration {
	return NewDuration(uint64DTO(dto).toUint64())
}

type difficultyDTO uint64DTO

func (dto difficultyDTO) toStruct() *Difficulty {
	return NewDifficulty(uint64DTO(dto).toUint64())
}

type blockchainTimestampDTO uint64DTO

func (dto blockchainTimestampDTO) toStruct() *BlockchainTimestamp {
	return NewBlockchainTimestamp(int64(uint64DTO(dto).toUint64()))
}

func hasBits(number uint64, bits uint64) bool {
	return (number & bits) == bits
}
