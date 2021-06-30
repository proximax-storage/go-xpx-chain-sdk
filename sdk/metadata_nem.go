package sdk

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
)

// TODO: Add NEM's metadata routes. But the problem is that one route can return all 3 types of metadata
//type MetadataNemService service

func CalculateUniqueAccountMetadataId(sourceAddress *Address, targetAccount *PublicAccount, key ScopedMetadataKey) (*Hash, error) {
	return calculate(sourceAddress, targetAccount, key, 0, 0)
}

func CalculateUniqueMosaicMetadataId(sourceAddress *Address, targetAccount *PublicAccount, key ScopedMetadataKey, mosaic *MosaicId) (*Hash, error) {
	return calculate(sourceAddress, targetAccount, key, mosaic.baseInt64, 1)
}

func CalculateUniqueNamespaceMetadataId(sourceAddress *Address, targetAccount *PublicAccount, key ScopedMetadataKey, namespace *NamespaceId) (*Hash, error) {
	return calculate(sourceAddress, targetAccount, key, namespace.baseInt64, 2)
}

func calculate(sourceAddress *Address, targetAccount *PublicAccount, key ScopedMetadataKey, targetId baseInt64, metadataType uint8) (*Hash, error) {
	result := sha3.New256()
	source, err := sourceAddress.Decode()
	if err != nil {
		return nil, err
	}

	targetKey, err := hex.DecodeString(targetAccount.PublicKey)
	if err != nil {
		return nil, err
	}

	if _, err := result.Write(source[:]); err != nil {
		return nil, err
	}

	if _, err := result.Write(targetKey[:]); err != nil {
		return nil, err
	}

	if _, err := result.Write(key.toLittleEndian()); err != nil {
		return nil, err
	}

	if _, err := result.Write(targetId.toLittleEndian()); err != nil {
		return nil, err
	}

	if _, err := result.Write([]byte{metadataType}); err != nil {
		return nil, err
	}

	hash, err := bytesToHash(result.Sum(nil))
	if err != nil {
		return nil, err
	}

	return hash, nil
}
