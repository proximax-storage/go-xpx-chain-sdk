// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

const (
	address = "http://10.32.150.136:3000"
)

var (
	ctx        = context.Background()
	mockServer = newSdkMock(5 * time.Minute)
)

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// Uint64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Uint64(v uint64) *uint64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

func TestBigIntegerToHex_bigIntegerNEMAndXEMToHex(t *testing.T) {
	testBigInt(t, "15358872602548358953", "d525ad41d95fcf29")
	testBigInt(t, "9562080086528621131", "84b3552d375ffa4b")
	testBigInt(t, "153588726025483589", "0221a821f040f545")
	testBigInt(t, "-7680974160236284465", "9567b2b2622975cf")
	testBigInt(t, "23160236284465", "0000151069a81a31")
}

func testBigInt(t *testing.T, str, hexStr string) {
	i, ok := (&big.Int{}).SetString(str, 10)
	assert.True(t, ok)
	result := bigIntegerToHex(i)
	assert.Equal(t, hexStr, result)
}

type sdkMock struct {
	*mock.Mock
}

func newSdkMock(closeAfter time.Duration) *sdkMock {
	return &sdkMock{mock.NewMock(closeAfter)}
}

func newSdkMockWithRouter(router *mock.Router) *sdkMock {
	sdkMock := &sdkMock{mock.NewMock(0)}

	sdkMock.AddRouter(router)

	return sdkMock
}

func (m sdkMock) getClientByNetworkType(networkType NetworkType) (*Client, error) {
	conf, err := NewConfig(m.GetServerURL(), networkType)

	if err != nil {
		return nil, err
	}

	client := NewClient(nil, conf)

	return client, nil
}

func (m *sdkMock) getPublicTestClient() (*Client, error) {
	return m.getClientByNetworkType(PublicTest)
}

func (m *sdkMock) getPublicTestClientUnsafe() *Client {
	client, _ := m.getPublicTestClient()

	return client
}
