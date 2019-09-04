// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBaseInt64_ToArray(t *testing.T) {
	base := baseInt64(9543417332823)
	want := [2]uint32{1111, 2222}
	assert.Equal(t, want, base.toArray())

	base = baseInt64(429492434645049)
	want = [2]uint32{12345, 99999}
	assert.Equal(t, want, base.toArray())
}

func TestBaseInt64_Bytes(t *testing.T) {
	base := baseInt64(9543417332823)
	assert.Equal(t, []byte{0x57, 0x4, 0x0, 0x0, 0xae, 0x8, 0x0, 0x0}, base.toLittleEndian())
}

func TestBlockchainTimestampConversion(t *testing.T) {
	deadline := NewDeadline(time.Hour)
	assert.Equal(t, deadline.Second(), deadline.ToBlockchainTimestamp().ToTimestamp().Second())
}

func TestNewDeadlineFromBlockchainTimestamp(t *testing.T) {
	deadline := NewDeadlineFromBlockchainTimestamp(NewBlockchainTimestamp(0))
	assert.Equal(t, time.Unix(0, TimestampNemesisBlockMilliseconds*int64(time.Millisecond)).String(), deadline.String())
}
