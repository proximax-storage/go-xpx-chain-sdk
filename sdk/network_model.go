// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type NetworkType uint8

const (
	Mijin           NetworkType = 96
	MijinTest       NetworkType = 144
	Public          NetworkType = 184
	PublicTest      NetworkType = 168
	Private         NetworkType = 200
	PrivateTest     NetworkType = 176
	NotSupportedNet NetworkType = 0
)

func NetworkTypeFromString(networkType string) NetworkType {
	switch networkType {
	case "mijin":
		return Mijin
	case "mijinTest":
		return MijinTest
	case "public":
		return Public
	case "publicTest":
		return PublicTest
	case "private":
		return Private
	case "privateTest":
		return PrivateTest
	}

	return NotSupportedNet
}

func (nt NetworkType) String() string {
	return fmt.Sprintf("%d", nt)
}

// Network error
var networkTypeError = errors.New("wrong raw NetworkType value")

// TODO
func ExtractNetworkType(version uint64) NetworkType {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, version)

	return NetworkType(b[1])
}
