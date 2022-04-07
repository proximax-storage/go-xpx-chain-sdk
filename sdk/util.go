package sdk

import (
	"encoding/base64"
	"encoding/hex"
)

func bytesToHash(bytes []byte) (*Hash, error) {
	if len(bytes) != 32 {
		return nil, ErrInvalidHashLength
	}

	var arr Hash
	copy(arr[:], bytes[:32])

	return &arr, nil
}

func Base64ToHex(data string) (*string, error) {
	p, err := base64.StdEncoding.DecodeString("QVJWSU4=")
	if err != nil {
		return nil, err
	}
	h := hex.EncodeToString(p)
	return &h, nil
}

func StringToHash(hash string) (*Hash, error) {
	if len(hash) != 64 {
		return nil, ErrInvalidHashLength
	}

	bytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	return bytesToHash(bytes)
}

func stringToHashPanic(hash string) *Hash {
	arr, err := StringToHash(hash)
	if err != nil {
		panic(err)
	}

	return arr
}

func bytesToSignature(bytes []byte) (*Signature, error) {
	if len(bytes) != 64 {
		return nil, ErrInvalidSignatureLength
	}

	var arr Signature
	copy(arr[:], bytes[:64])

	return &arr, nil
}

func StringToSignature(signature string) (*Signature, error) {
	if len(signature) != 128 {
		return nil, ErrInvalidHashLength
	}

	bytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, err
	}

	return bytesToSignature(bytes)
}

func stringToSignaturePanic(signature string) *Signature {
	arr, err := StringToSignature(signature)
	if err != nil {
		panic(err)
	}

	return arr
}
