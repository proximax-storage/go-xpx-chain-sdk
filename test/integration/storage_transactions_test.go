// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by a BSD-style
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/proximax-storage/go-xpx-chain-sdk/sdk"
)

func TestStorageDrivePrepareTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStoragePrepareDriveTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.Duration(100),
			100,
			2,
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}

func TestNewStorageDriveVerificationTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageDriveVerificationTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}

func TestStorageDriveProlongationTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageDriveProlongationTransaction(
			sdk.NewDeadline(time.Hour),
			sdk.Duration(100),
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageFileDepositTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)
	hash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageFileDepositTransaction(
			sdk.NewDeadline(time.Hour),
			hash,
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageDriveDepositTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)
	hash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageDriveDepositTransaction(
			sdk.NewDeadline(time.Hour),
			hash,
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageFileDepositReturnTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)
	hash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageFileDepositReturnTransaction(
			sdk.NewDeadline(time.Hour),
			hash,
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageDriveDepositReturnTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)
	hash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageDriveDepositReturnTransaction(
			sdk.NewDeadline(time.Hour),
			hash,
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageFilePaymentTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)
	hash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageFilePaymentTransaction(
			sdk.NewDeadline(time.Hour),
			hash,
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageDrivePaymentTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)
	hash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageDrivePaymentTransaction(
			sdk.NewDeadline(time.Hour),
			hash,
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageCreateDirectoryTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	parentHash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	hash, err := sdk.StringToHash("b2b4469f696f517977f1224218e49194358c155c30c2dbb2e38351354400db16")
	fileName := "test"

	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageCreateDirectoryTransaction(
			sdk.NewDeadline(time.Hour),
			&sdk.StorageFile{
				Hash:       hash,
				ParentHash: parentHash,
				Name:       fileName,
			},
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageRemoveDirectoryTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	parentHash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	hash, err := sdk.StringToHash("b2b4469f696f517977f1224218e49194358c155c30c2dbb2e38351354400db16")
	fileName := "test"

	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageRemoveDirectoryTransaction(
			sdk.NewDeadline(time.Hour),
			&sdk.StorageFile{
				Hash:       hash,
				ParentHash: parentHash,
				Name:       fileName,
			},
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageUploadFileTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	parentHash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	hash, err := sdk.StringToHash("b2b4469f696f517977f1224218e49194358c155c30c2dbb2e38351354400db16")
	fileName := "test"

	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageUploadFileTransaction(
			sdk.NewDeadline(time.Hour),
			&sdk.StorageFile{
				Hash:       hash,
				ParentHash: parentHash,
				Name:       fileName,
			},
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageDownloadFileTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	parentHash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	hash, err := sdk.StringToHash("b2b4469f696f517977f1224218e49194358c155c30c2dbb2e38351354400db16")
	fileName := "test"

	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageDownloadFileTransaction(
			sdk.NewDeadline(time.Hour),
			&sdk.StorageFile{
				Hash:       hash,
				ParentHash: parentHash,
				Name:       fileName,
			},
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageDeleteFileTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	parentHash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	hash, err := sdk.StringToHash("b2b4469f696f517977f1224218e49194358c155c30c2dbb2e38351354400db16")
	fileName := "test"

	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageDeleteFileTransaction(
			sdk.NewDeadline(time.Hour),
			&sdk.StorageFile{
				Hash:       hash,
				ParentHash: parentHash,
				Name:       fileName,
			},
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
func TestStorageMoveFileTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	parentHash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	hash, err := sdk.StringToHash("b2b4469f696f517977f1224218e49194358c155c30c2dbb2e38351354400db16")
	destinationHash, err := sdk.StringToHash("a38e57051364b4363353acb21bb60b9cdbfd4ebe0d90433569bc6a7adf4aaee4")
	fileName := "test"

	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageMoveFileTransaction(
			sdk.NewDeadline(time.Hour),
			&sdk.StorageFile{
				Hash:       hash,
				ParentHash: parentHash,
				Name:       fileName,
			},
			&sdk.StorageFile{
				Hash:       destinationHash,
				ParentHash: parentHash,
				Name:       fileName,
			},
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}

func TestStorageCopyFileTransaction(t *testing.T) {
	testAccount, err := client.NewAccount()
	assert.Nil(t, err)

	fmt.Println(testAccount)

	parentHash, err := sdk.StringToHash("895e166388ae24d7f18c6a7a7b271c01182608697044d7f65c013f89a2d2a8b4")
	hash, err := sdk.StringToHash("b2b4469f696f517977f1224218e49194358c155c30c2dbb2e38351354400db16")
	destinationHash, err := sdk.StringToHash("a38e57051364b4363353acb21bb60b9cdbfd4ebe0d90433569bc6a7adf4aaee4")
	fileName := "test"

	assert.Nil(t, err)
	result := sendTransaction(t, func() (sdk.Transaction, error) {
		return client.NewStorageCopyFileTransaction(
			sdk.NewDeadline(time.Hour),
			&sdk.StorageFile{
				Hash:       hash,
				ParentHash: parentHash,
				Name:       fileName,
			},
			&sdk.StorageFile{
				Hash:       destinationHash,
				ParentHash: parentHash,
				Name:       fileName,
			},
			sdk.MijinTest,
		)
	}, testAccount)
	assert.Nil(t, result.error)
}
