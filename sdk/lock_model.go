// Copyright 2019 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"golang.org/x/crypto/sha3"
)

type LockStatusType uint8

const (
	/// Lock is unused.
	Unused LockStatusType = iota

	/// Lock was already used.
	Used
)

type CommonLockInfo struct {
	Account     *PublicAccount
	MosaicId    *MosaicId
	Amount      Amount
	Height      Height
	Status      LockStatusType
}

type HashLockInfo struct {
	CommonLockInfo
	Hash    *Hash
}

type SecretLockInfo struct {
	CommonLockInfo
	HashAlgorithm   HashType
	CompositeHash   *Hash
	Secret          *Hash
	Recipient       *PublicAccount
}

func CalculateCompositeHash(secret *Hash, recipient *Address) (*Hash, error) {
	if secret == nil {
		return nil, ErrNilSecret
	}
	if recipient == nil {
		return nil, ErrNilAddress
	}

	result := sha3.New256()

	if _, err := result.Write(secret[:]); err != nil {
		return nil, err
	}

	recipientB, err := recipient.Decode()
	if err != nil {
		return nil, err
	}

	if _, err := result.Write(recipientB); err != nil {
		return nil, err
	}

	hash ,err := bytesToHash(result.Sum(nil))
	if err != nil {
		return nil, err
	}

	return hash, nil
}