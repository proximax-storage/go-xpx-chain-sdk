// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"github.com/proximax-storage/go-xpx-utils/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	address = "http://127.0.0.1:3000"
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
	testHexConversion(t, 15358872602548358953, "D525AD41D95FCF29")
	testHexConversion(t, 9562080086528621131, "84B3552D375FFA4B")
	testHexConversion(t, 153588726025483589, "0221A821F040F545")
	testHexConversion(t, 0x9567B2B2622975CF, "9567B2B2622975CF")
	testHexConversion(t, 23160236284465, "0000151069A81A31")
}

func testHexConversion(t *testing.T, id uint64, hexStr string) {
	assert.Equal(t, hexStr, newNamespaceIdPanic(id).toHexString())
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
	conf, err := NewConfigWithReputation([]string{m.GetServerURL()}, networkType, &defaultRepConfig, DefaultWebsocketReconnectionTimeout, nil)

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

func TestClient_AdaptAccount(t *testing.T) {
	var stockHash = &Hash{1}
	var defaultHash = &Hash{2}
	account, err := NewAccount(PublicTest, stockHash)
	assert.Nil(t, err)

	config, err := NewConfigWithReputation([]string{""}, MijinTest, &defaultRepConfig, DefaultWebsocketReconnectionTimeout, defaultHash)
	assert.Nil(t, err)

	client := NewClient(nil, config)

	adaptedAccount, err := client.AdaptAccount(account)
	assert.Equal(t, MijinTest, adaptedAccount.PublicAccount.Address.Type)
	assert.Equal(t, defaultHash, adaptedAccount.generationHash)

	assert.Equal(t, PublicTest, account.PublicAccount.Address.Type)
	assert.Equal(t, stockHash, account.generationHash)
}
