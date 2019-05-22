package sdk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSecret_SHA3_256(t *testing.T) {
	proof, err := NewProofFromHexString("B778A39A3663719DFC5E48C9D78431B1E45C2AF9DF538782BF199C189DABEAC7")
	assert.Nil(t, err)

	secret, err := proof.Secret(SHA3_256)
	assert.Nil(t, err)
	assert.Equal(t, "9B3155B37159DA50AA52D5967C509B410F5A36A3B1E31ECB5AC76675D79B4A5E", secret.HashString())
	assert.Equal(t, SHA3_256, secret.Type)
}

func TestGenerateSecret_KECCAK_256(t *testing.T) {
	proof, err := NewProofFromHexString("B778A39A3663719DFC5E48C9D78431B1E45C2AF9DF538782BF199C189DABEAC7")
	assert.Nil(t, err)

	secret, err := proof.Secret(KECCAK_256)
	assert.Nil(t, err)
	assert.Equal(t, "241C1D54C18C8422DEF03AA16B4B243A8BA491374295A1A6965545E6AC1AF314", secret.HashString())
	assert.Equal(t, KECCAK_256, secret.Type)
}

func TestGenerateSecret_HASH_160(t *testing.T) {
	proof := NewProofFromUint8(97)

	secret, err := proof.Secret(HASH_160)
	assert.Nil(t, err)
	assert.Equal(t, "994355199E516FF76C4FA4AAB39337B9D84CF12B000000000000000000000000", secret.HashString())
	assert.Equal(t, HASH_160, secret.Type)
}

func TestGenerateSecret_SHA_256(t *testing.T) {
	proof, err := NewProofFromHexString("DE188941A3375D3A8A061E67576E926D")
	assert.Nil(t, err)

	secret, err := proof.Secret(SHA_256)
	assert.Nil(t, err)
	assert.Equal(t, "2182D3FE9882FD597D25DAF6A85E3A574E5A9861DBC75C13CE3F47FE98572246", secret.HashString())
	assert.Equal(t, SHA_256, secret.Type)
}
