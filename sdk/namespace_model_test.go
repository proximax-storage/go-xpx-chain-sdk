// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGenerateNamespacePath_GeneratesCorrectWellKnownRootPath(t *testing.T) {
	ids, err := GenerateNamespacePath("nem")
	assert.Nil(t, err)

	assert.Equal(t, len(ids), 1, `ids.size() and 1 must by equal !`)

	assert.Equal(t, uint64(0x84B3552D375FFA4B), ids[0].Id())
}

// @Test
func TestNamespacePath_GeneratesCorrectWellKnownChildPath(t *testing.T) {
	ids, err := GenerateNamespacePath("nem.xem")
	assert.Nil(t, err)
	assert.Equal(t, len(ids), 2, `ids.size() and 2 must be equal !`)

	assert.Equal(t, uint64(0x84B3552D375FFA4B), ids[0].Id())
	assert.Equal(t, uint64(0xD525AD41D95FCF29), ids[1].Id())
}

// @Test
func TestNamespacePathSupportsMultiLevelNamespaces(t *testing.T) {
	ids := make([]*NamespaceId, 3)
	var err error
	ids[0], err = generateNamespaceId("foo", NewNamespaceIdNoCheck(0))
	assert.Nil(t, err)
	ids[1], err = generateNamespaceId("bar", ids[0])
	assert.Nil(t, err)
	ids[2], err = generateNamespaceId("baz", ids[1])
	assert.Nil(t, err)
	ids1, err := GenerateNamespacePath("foo.bar.baz")
	assert.Nil(t, err)
	assert.Equal(t, ids1, ids, `GenerateNamespacePath("foo.bar.baz") and ids must by equal !`)
}

// @Test
func TestNamespacePathRejectsNamesWithTooManyParts(t *testing.T) {
	_, err := GenerateNamespacePath("a.b.c.d")
	assert.Equal(t, ErrNamespaceTooManyPart, err, "Err 'too many parts' must return")
	_, err = GenerateNamespacePath("a.b.c.d.e")
	assert.Equal(t, ErrNamespaceTooManyPart, err, "Err 'too many parts' must return")

}

// @Test
func TestMosaicIdGeneratesCorrectWellKnowId(t *testing.T) {
	account, err := NewAccountFromPrivateKey("C06B2CC5D7B66900B2493CF68BE10B7AA8690D973B7F0B65D0DAE4F7AA464716", MijinTest)
	assert.Nil(t, err)
	id, err := generateMosaicId(0, account.PublicAccount.PublicKey)
	assert.Nil(t, err)
	assert.Equal(t, uint64(992621222383397347), id.Id())
}

// @Test
func TestNewAddressFromNamespace(t *testing.T) {
	namespaceId := NewNamespaceIdNoCheck(0x85bbea6cc462b244)
	address, err := NewAddressFromNamespace(namespaceId)
	assert.Nil(t, err)

	fmt.Println(address)

	pH, err := base32.StdEncoding.DecodeString(address.Address)
	assert.Nil(t, err)
	parsed := strings.ToUpper(hex.EncodeToString(pH))

	assert.Equal(t, address.Type, AliasAddress)
	assert.Equal(t, "9144B262C46CEABB8500000000000000000000000000000000", parsed)
}
