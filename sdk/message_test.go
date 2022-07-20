package sdk

import (
	"testing"

	crypto "github.com/proximax-storage/go-xpx-crypto"
	"github.com/stretchr/testify/assert"
)

func TestPlaintTextToSecureMessageAndBackEdSha2(t *testing.T) {
	const message = "Hello guys, let's do this!"
	sender, err := crypto.NewKeyPairByEngine(crypto.CryptoEngines.Ed25519Sha2Engine)
	assert.Nil(t, err)
	recipient, err := crypto.NewKeyPairByEngine(crypto.CryptoEngines.Ed25519Sha2Engine)
	assert.Nil(t, err)

	secureMessage, err := NewSecureMessageFromPlaintText(message, sender.PrivateKey, recipient.PublicKey, crypto.CryptoEngines.Ed25519Sha2Engine)
	assert.Nil(t, err)

	plainMessage, err := NewPlainMessageFromEncodedData(secureMessage.Payload(), recipient.PrivateKey, sender.PublicKey, crypto.CryptoEngines.Ed25519Sha2Engine)
	assert.Nil(t, err)

	assert.Equal(t, message, plainMessage.Message())
}

func TestPlaintTextToSecureMessageAndBackEdSha3(t *testing.T) {
	const message = "Hello guys, let's do this!"
	sender, err := crypto.NewKeyPairByEngine(crypto.CryptoEngines.Ed25519Sha3Engine)
	assert.Nil(t, err)
	recipient, err := crypto.NewKeyPairByEngine(crypto.CryptoEngines.Ed25519Sha3Engine)
	assert.Nil(t, err)

	secureMessage, err := NewSecureMessageFromPlaintText(message, sender.PrivateKey, recipient.PublicKey, crypto.CryptoEngines.Ed25519Sha3Engine)
	assert.Nil(t, err)

	plainMessage, err := NewPlainMessageFromEncodedData(secureMessage.Payload(), recipient.PrivateKey, sender.PublicKey, crypto.CryptoEngines.Ed25519Sha3Engine)
	assert.Nil(t, err)

	assert.Equal(t, message, plainMessage.Message())
}
