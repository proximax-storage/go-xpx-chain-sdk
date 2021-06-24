package sdk

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	blst "github.com/supranational/blst/bindings/go"
	"io"
	"reflect"
	"strings"
)

const FILECOIN_DST = string("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_NUL_")
const ETH2_DST = string("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

// It must be a 32 byte sequence.
type BLSPrivateKey string

var ZeroBLSPrivateKey = BLSPrivateKey(make([]byte, 32))

// It must be a 48 byte sequence.
type BLSPublicKey string

var ZeroBLSPublicKey = BLSPublicKey(make([]byte, 48))

// It must be a 96 byte sequence.
type BLSSignature string

var ZeroBLSSignature = BLSSignature(make([]byte, 96))

func isNilFixed(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func GeneratePrivateKey(seed io.Reader) BLSPrivateKey {
	if isNilFixed(seed) {
		seed = rand.Reader
	}
	var ikm [32]byte
	_, _ = seed.Read(ikm[:])
	return GeneratePrivateKeyFromIKM(ikm)
}

func GeneratePrivateKeyFromIKM(ikm [32]byte) BLSPrivateKey {
	sk := blst.KeyGen(ikm[:])
	return BLSPrivateKey(sk.ToLEndian())
}

func (priv BLSPrivateKey) sk() *blst.SecretKey {
	return new(blst.SecretKey).FromLEndian([]byte(priv))
}

func (priv BLSPrivateKey) Public() BLSPublicKey {
	return BLSPublicKey(new(blst.P1Affine).From(priv.sk()).Compress())
}

func (priv BLSPrivateKey) Sign(msg string) BLSSignature {
	return BLSSignature(new(blst.P2Affine).Sign(priv.sk(), []byte(msg), []byte(FILECOIN_DST)).Compress())
}

func (priv BLSPrivateKey) HexString() string {
	return strings.ToUpper(hex.EncodeToString([]byte(priv)))
}

func (pub BLSPublicKey) Verify(msg string, signature BLSSignature) bool {
	pk := new(blst.P1Affine).Uncompress([]byte(pub))
	sig := new(blst.P2Affine).Uncompress([]byte(signature))
	return sig.Verify(true, pk, true, []byte(msg), []byte(FILECOIN_DST))
}

func (pub BLSPublicKey) HexString() string {
	return strings.ToUpper(hex.EncodeToString([]byte(pub)))
}

func (signature BLSSignature) Verify(msg string, pub BLSPublicKey) bool {
	pk := new(blst.P1Affine).Uncompress([]byte(pub))
	sig := new(blst.P2Affine).Uncompress([]byte(signature))
	return sig.Verify(true, pk, true, []byte(msg), []byte(FILECOIN_DST))
}

func (signature BLSSignature) HexString() string {
	return strings.ToUpper(hex.EncodeToString([]byte(signature)))
}

type KeyPair struct {
	PublicKey  BLSPublicKey
	PrivateKey BLSPrivateKey
}

func GenerateKeyPair(seed io.Reader) *KeyPair {
	sk := GeneratePrivateKey(seed)
	pk := sk.Public()
	return &KeyPair{
		pk, sk,
	}
}

func GenerateKeyPairFromIKM(ikm [32]byte) *KeyPair {
	sk := GeneratePrivateKeyFromIKM(ikm)
	pk := sk.Public()
	return &KeyPair{
		pk, sk,
	}
}

func (p *KeyPair) Sign(msg string) BLSSignature {
	return p.PrivateKey.Sign(msg)
}

func AggregateSignatures(signatures ...BLSSignature) (BLSSignature, error) {
	if len(signatures) == 0 {
		return ZeroBLSSignature, nil
	}
	temp := make([][]byte, len(signatures))
	aggregator := new(blst.P2Aggregate)
	for i, signature := range signatures {
		temp[i] = []byte(signature)
	}
	if !aggregator.AggregateCompressed(temp, true) {
		return ZeroBLSSignature, errors.New("unable to aggregate signature")
	}

	return BLSSignature(aggregator.ToAffine().Compress()), nil
}

func AggregatePublicKeys(keys ...BLSPublicKey) (BLSPublicKey, error) {
	if len(keys) == 0 {
		return ZeroBLSPublicKey, nil
	}
	temp := make([][]byte, len(keys))
	aggregator := new(blst.P1Aggregate)
	for i, key := range keys {
		temp[i] = []byte(key)
	}
	if !aggregator.AggregateCompressed(temp, true) {
		return ZeroBLSPublicKey, errors.New("unable to aggregate public keys")
	}
	return BLSPublicKey(aggregator.ToAffine().Compress()), nil
}

func AggregateVerify(keys []BLSPublicKey, messages []string, signature BLSSignature) bool {
	if len(keys) != len(messages) {
		return false
	}
	pks := make([][]byte, len(keys))
	msgs := make([]blst.Message, len(messages))
	for i, key := range keys {
		pks[i] = []byte(key)
		msgs[i] = blst.Message(messages[i])
	}
	return new(blst.P2Affine).AggregateVerifyCompressed([]byte(signature), true, pks, true, msgs, []byte(FILECOIN_DST))
}

func FastAggregateVerify(keys []BLSPublicKey, message string, signature BLSSignature) bool {
	aggregatedKey, err := AggregatePublicKeys(keys...)
	if err != nil {
		return false
	}

	return aggregatedKey.Verify(message, signature)
}
