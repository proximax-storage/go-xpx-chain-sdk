package sdk

import (
	"encoding/hex"
	"errors"

	"github.com/proximax-storage/go-xpx-utils/str"
	xpxcrypto "github.com/proximax-storage/go-xpx-crypto"
)

type MessageType uint8

const (
	PlainMessageType MessageType = iota
	SecureMessageType
)

type Message interface {
	Type() MessageType
	Payload() []byte
	String() string
}

type PlainMessage struct {
	payload []byte
}

func (m *PlainMessage) String() string {
	return str.StructToString(
		"PlainMessage",
		str.NewField("Type", str.IntPattern, m.Type()),
		str.NewField("Payload", str.StringPattern, m.Payload()),
	)
}

func (m *PlainMessage) Type() MessageType {
	return PlainMessageType
}

func (m *PlainMessage) Payload() []byte {
	return m.payload
}

func (m *PlainMessage) Message() string {
	return string(m.payload)
}

func NewPlainMessage(payload string) *PlainMessage {
	return &PlainMessage{[]byte(payload)}
}

func NewPlainMessageFromEncodedData(encodedData []byte, recipient *xpxcrypto.PrivateKey, sender *xpxcrypto.PublicKey) (*PlainMessage, error) {
	rkp, err := xpxcrypto.NewKeyPair(recipient, nil, nil)
	if err != nil {
		return nil, err
	}

	skp, err := xpxcrypto.NewKeyPair(nil, sender, nil)
	if err != nil {
		return nil, err
	}

	plaintText, err := xpxcrypto.NewBlockCipher(skp, rkp, nil).Decrypt(encodedData)
	if err != nil {
		return nil, err
	}

	return NewPlainMessage(string(plaintText)), nil
}

type SecureMessage struct {
	encodedData []byte
}

func (m *SecureMessage) String() string {
	return str.StructToString(
		"SecureMessage",
		str.NewField("Type", str.IntPattern, m.Type()),
		str.NewField("EncodedData", str.StringPattern, m.Payload()),
	)
}

func (m *SecureMessage) Payload() []byte {
	return m.encodedData
}

func (m *SecureMessage) Type() MessageType {
	return SecureMessageType
}

func NewSecureMessage(encodedData []byte) *SecureMessage {
	return &SecureMessage{encodedData}
}

func NewSecureMessageFromPlaintText(plaintText string, sender *xpxcrypto.PrivateKey, recipient *xpxcrypto.PublicKey) (*SecureMessage, error) {
	skp, err := xpxcrypto.NewKeyPair(sender, nil, nil)
	if err != nil {
		return nil, err
	}

	rkp, err := xpxcrypto.NewKeyPair(nil, recipient, nil)
	if err != nil {
		return nil, err
	}

	encodedData, err := xpxcrypto.NewBlockCipher(skp, rkp, nil).Encrypt([]byte(plaintText))
	if err != nil {
		return nil, err
	}

	return NewSecureMessage(encodedData), nil
}

type messageDTO struct {
	Type    MessageType `json:"type"`
	Payload string      `json:"payload"`
}

func (m *messageDTO) toStruct() (Message, error) {
	b := make([]byte, 0)

	if len(m.Payload) != 0 {
		var err error
		b, err = hex.DecodeString(m.Payload)

		if err != nil {
			return nil, err
		}
	}

	switch m.Type {
	case PlainMessageType:
		return NewPlainMessage(string(b)), nil
	case SecureMessageType:
		return NewSecureMessage(b), nil
	default:
		return nil, errors.New("Not supported MessageType")
	}
}
