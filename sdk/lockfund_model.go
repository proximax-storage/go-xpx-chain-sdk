// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"fmt"
)

type LockFundAction uint8

type LockFundHeightRecord struct {
	Identifier Height
	Records    map[string]*LockFundRecord
}

type LockFundKeyRecord struct {
	Identifier *PublicAccount
	Records    map[Height]*LockFundRecord
}

type LockFundRecord struct {
	ActiveRecord    []*Mosaic
	InactiveRecords []*([]*Mosaic)
}

func (s *LockFundHeightRecord) String() string {
	return fmt.Sprintf(
		`
			"Identifier": %s,
			"Records": %T
		`,
		s.Identifier,
		s.Records,
	)
}

type LockFundTransferTransaction struct {
	AbstractTransaction
	Duration Duration
	Action   LockFundAction
	Mosaics  []*Mosaic
}

type LockFundCancelUnlockTransaction struct {
	AbstractTransaction
	TargetHeight Height
}
