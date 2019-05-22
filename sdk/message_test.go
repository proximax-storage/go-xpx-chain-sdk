package sdk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlaintTexToSecureMessageAndBack(t *testing.T) {
	const networkType = MijinTest
	const message = "Hello guys, let's do this!"
	sender, err := NewAccount(networkType)
	assert.Nil(t, err)
	recipient, err := NewAccount(networkType)
	assert.Nil(t, err)

	secureMessage, err := NewSecureMessageFromPlaintText(message, sender.PrivateKey, recipient.KeyPair.PublicKey)
	assert.Nil(t, err)

	plainMessage, err := NewPlainMessageFromEncodedData(secureMessage.Payload(), recipient.PrivateKey, sender.KeyPair.PublicKey)
	assert.Nil(t, err)

	assert.Equal(t, message, string(plainMessage.Payload()))
}
