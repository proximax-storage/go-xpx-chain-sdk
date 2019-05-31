// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testNEMPublicKey = "b4f12e7c9f6946091e2cb8b6d3a12b50d17ccbbf646386ea27ce2946a7423dcf"
)

var testAddressesForEncoded = map[NetworkType]string{
	MijinTest:   "SARNASAS2BIAB6LMFA3FPMGBPGIJGK6IJETM3ZSP",
	Mijin:       "MARNASAS2BIAB6LMFA3FPMGBPGIJGK6IJE5K5RYU",
	Public:      "XARNASAS2BIAB6LMFA3FPMGBPGIJGK6IJF6CHIGW",
	PublicTest:  "VARNASAS2BIAB6LMFA3FPMGBPGIJGK6IJGOH3FCE",
	Private:     "ZARNASAS2BIAB6LMFA3FPMGBPGIJGK6IJF2S3UOQ",
	PrivateTest: "WARNASAS2BIAB6LMFA3FPMGBPGIJGK6IJHPRCU4F",
}

func TestGenerateNewAccount_NEM(t *testing.T) {
	acc, err := NewAccount(MijinTest)
	if err != nil {
		t.Fatal("Error")
	}
	a := acc.KeyPair.PrivateKey.String()
	t.Log("Private Key: " + a)

	assert.NotNil(t, acc.KeyPair.PrivateKey.String(), "Error generating new KeyPair")
}

func TestGenerateEncodedAddress_NEM(t *testing.T) {

	for nType, testAddress := range testAddressesForEncoded {

		res, err := generateEncodedAddress(testNEMPublicKey, nType)
		if err != nil {
			t.Fatal("Error")
		}

		assert.Equal(t, testAddress, res, "Wrong address")
	}
}

func TestGenerateEncodedAddress(t *testing.T) {
	res, err := generateEncodedAddress("321DE652C4D3362FC2DDF7800F6582F4A10CFEA134B81F8AB6E4BE78BBA4D18E", 144)
	if err != nil {
		t.Fatal("Error")
	}

	assert.Equal(t, "SBFBW6TUGLEWQIBCMTBMXXQORZKUP3WTVVTOKK5M", res, "Wrong address %s", res)
}

func TestEncryptMessageAndDecryptMessage(t *testing.T) {
	const networkType = MijinTest
	const message = "Hello guys, let's do this!"
	sender, err := NewAccount(networkType)
	assert.Nil(t, err)
	recipient, err := NewAccount(networkType)
	assert.Nil(t, err)

	secureMessage, err := sender.EncryptMessage(message, recipient.PublicAccount)
	assert.Nil(t, err)

	plainMessage, err := recipient.DecryptMessage(secureMessage, sender.PublicAccount)
	assert.Nil(t, err)

	assert.Equal(t, message, plainMessage.Message())
}
