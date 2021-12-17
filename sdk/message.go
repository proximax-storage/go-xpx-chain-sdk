package sdk

import (
	"encoding/hex"
	"errors"

	xpxcrypto "github.com/proximax-storage/go-xpx-crypto"
	"github.com/proximax-storage/go-xpx-utils/str"
)

type MessageType uint8

const (
	PlainMessageType MessageType = iota
	SecureMessageType
	PersistentHarvestingDelegationMessageType = 0xFE
)

type MessageMarker [8]uint8

const (
	PersistentDelegationUnlockMarker = "2A8061577301E2"
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
	case PersistentHarvestingDelegationMessageType:
		return NewPersistentHarvestingDelegationMessage(b), nil
	default:
		return nil, errors.New("Not supported MessageType")
	}
}

type PersistentHarvestingDelegationMessage struct {
	payload []byte
}

func (m *PersistentHarvestingDelegationMessage) String() string {
	return str.StructToString(
		"PersistentHarvestingDelegationMessage",
		str.NewField("Type", str.IntPattern, m.Type()),
		str.NewField("Payload", str.StringPattern, m.Payload()),
	)
}

func (m *PersistentHarvestingDelegationMessage) Type() MessageType {
	return PersistentHarvestingDelegationMessageType
}

func (m *PersistentHarvestingDelegationMessage) Payload() []byte {
	return m.payload
}

func (m *PersistentHarvestingDelegationMessage) Message() string {
	return string(m.payload)
}

func NewPersistentHarvestingDelegationMessage(payload []byte) *PersistentHarvestingDelegationMessage {
	return &PersistentHarvestingDelegationMessage{payload}
}

func NewPersistentHarvestingDelegationMessageFromPlainText(harvesterPrivateKey *xpxcrypto.PrivateKey, vrfPrivateKey *xpxcrypto.PrivateKey, recipient *xpxcrypto.PublicKey) (*PersistentHarvestingDelegationMessage, error) {
	// Ephemeral keys always use Ed25519 Sha3 Engine
	ephemeralKeyPair, err := xpxcrypto.NewKeyPairByEngine(Ephemeral_Key_Derivation_Scheme)
	if err != nil {
		return nil, err
	}
	concat := append(harvesterPrivateKey.Raw, vrfPrivateKey.Raw...)
	recipientKeyPair, err := xpxcrypto.NewKeyPair(nil, recipient, Node_Boot_Key_Derivation_Scheme)
	if err != nil {
		return nil, err
	}
	blockCypher := xpxcrypto.NewEd25519Sha3BlockCipher(ephemeralKeyPair, recipientKeyPair, nil)
	encoded, err := blockCypher.EncryptGCMNacl(concat, nil)
	if err != nil {
		return nil, err
	}
	marker, err := hex.DecodeString(PersistentDelegationUnlockMarker)
	if err != nil {
		return nil, err
	}
	encrypted := append(append(marker, ephemeralKeyPair.PublicKey.Raw...), encoded...)

	return &PersistentHarvestingDelegationMessage{encrypted}, nil
}
