// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"math/big"
	"regexp"
)

// mosaic id for XEM mosaic
var XemMosaicId, _ = NewMosaicId(big.NewInt(0x0DC67FBE1CAD29E3))
// mosaic id for XPX mosaic
var XpxMosaicId, _ = NewMosaicId(big.NewInt(0x0DC67FBE1CAD29E3))

// routes for account service
const (
	accountsRoute                 = "/account"
	accountRoute                  = "/account/%s"
	multisigAccountRoute          = "/account/%s/multisig"
	multisigAccountGraphInfoRoute = "/account/%s/multisig/graph"
	transactionsByAccountRoute    = "/account/%s/%s"
	accountTransactionsRoute      = "/transactions"
	incomingTransactionsRoute     = "/transactions/incoming"
	outgoingTransactionsRoute     = "/transactions/outgoing"
	unconfirmedTransactionsRoute  = "/transactions/unconfirmed"
	aggregateTransactionsRoute    = "/transactions/partial"
)

// routes for namespace service
const (
	namespaceRoute              = "/namespace/%s"
	namespacesFromAccountsRoute = "/account/namespaces"
	namespaceNamesRoute         = "/namespace/names"
	namespacesFromAccountRoutes = "/account/%s/namespaces"
)

// routes for mosaic service
const (
	mosaicsRoute = "/mosaic"
	mosaicRoute  = "/mosaic/%s"
)

// routes for blockchain service
const (
	blockHeightRoute         = "/chain/height"
	blockByHeightRoute       = "/block/%d"
	blockScoreRoute          = "/chain/score"
	blockGetTransactionRoute = "/block/%d/transactions"
	blockInfoRoute           = "/blocks/%d/limit/%d"
	blockStorageRoute        = "/diagnostic/storage"
)

// routes for contracts service
const (
	contractsInfoRoute      = "/contract"
	contractsByAccountRoute = "/account/%s/contracts"
)

// routes for metadata service
const (
	metadatasInfoRoute       = "/metadata"
	metadataInfoRoute        = "/metadata/%s"
	metadataByAccountRoute   = "/account/%s/metadata"
	metadataByMosaicRoute    = "/mosaic/%s/metadata"
	metadataByNamespaceRoute = "/namespace/%s/metadata"
)

// routes for network service
const (
	networkRoute = "/network"
)

// routes for transaction service
const (
	transactionsRoute                 = "/transaction"
	transactionRoute                  = "/transaction/%s"
	transactionStatusRoute            = "/transaction/%s/status"
	transactionsStatusRoute           = "/transaction/statuses"
	announceAggregateRoute            = "/transaction/partial"
	announceAggregateCosignatureRoute = "/transaction/cosignature"
)

type NamespaceType uint8

const (
	Root NamespaceType = iota
	Sub
)

// regValidNamespace check namespace on valid symbols
var (
	regValidNamespace = regexp.MustCompile(`^[a-z0-9][a-z0-9-_]*$`)
)
