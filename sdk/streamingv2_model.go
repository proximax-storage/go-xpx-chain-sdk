// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

type StreamStartTransaction struct {
	AbstractTransaction
	DriveKey           string
	ExpectedUploadSize StorageSize
	Folder             string
	FeedbackFeeAmount  Amount
}
