package sdk

import (
	"github.com/proximax-storage/xpx-crypto-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlaintTexToSecureMessageAndBack(t *testing.T) {
	const message = "Hello guys, let's do this!"
	sender, err := crypto.NewKeyPairByEngine(crypto.CryptoEngines.DefaultEngine)
	assert.Nil(t, err)
	recipient, err := crypto.NewKeyPairByEngine(crypto.CryptoEngines.DefaultEngine)
	assert.Nil(t, err)

	secureMessage, err := NewSecureMessageFromPlaintText(message, sender.PrivateKey, recipient.PublicKey)
	assert.Nil(t, err)

	plainMessage, err := NewPlainMessageFromEncodedData(secureMessage.Payload(), recipient.PrivateKey, sender.PublicKey)
	assert.Nil(t, err)

	assert.Equal(t, message, plainMessage.Message())
}
