// Copyright 2020 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"context"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/net"
	"net/http"
)

type ReceiptService service

func (s *ReceiptService) GetBlockStatementAtHeight(ctx context.Context, height Height) (*BlockStatement, error) {
	if height <= 1 {
		return nil, ErrArgumentNotValid
	}

	url := net.NewUrl(fmt.Sprintf(blockStatementsByHeight, height))

	dto := &BlockStatementDto{}

	resp, err := s.client.doNewRequest(ctx, http.MethodGet, url.Encode(), nil, dto)
	if err != nil {
		return nil, err
	}

	if err = handleResponseStatusCode(resp, map[int]error{404: ErrResourceNotFound, 409: ErrArgumentNotValid}); err != nil {
		return nil, err
	}

	return dto.toStruct(s.client.NetworkType())
}
