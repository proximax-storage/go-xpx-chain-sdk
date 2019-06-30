package sdk

import "encoding/hex"

func StringToHash(hash string) (Hash, error) {
	return hex.DecodeString(hash)
}

func StringToHashNoCheck(hash string) Hash {
	bytes, _ := hex.DecodeString(hash)

	return bytes
}

func StringToSignature(hash string) (Signature, error) {
	return hex.DecodeString(hash)
}

func StringToSignatureNoCheck(hash string) Signature {
	bytes, _ := hex.DecodeString(hash)

	return bytes
}
