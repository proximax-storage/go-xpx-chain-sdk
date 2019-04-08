// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"math/big"
	"regexp"
)

// 0x0DC67FBE1CAD29E3 = 992621222383397347
var XemMosaicId, _ = NewMosaicId(big.NewInt(0x0DC67FBE1CAD29E3))
var XpxMosaicId, _ = NewMosaicId(big.NewInt(0x0DC67FBE1CAD29E3))

// const routers path for methods AccountService
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

// const routers path for methods NamespaceService
const (
	namespaceRoute              = "/namespace/%s"
	namespacesFromAccountsRoute = "/account/namespaces"
	namespaceNamesRoute         = "/namespace/names"
	namespacesFromAccountRoutes = "/account/%s/namespaces"
)

// const routers path for methods MosaicService
const (
	mosaicsRoute = "/mosaic"
	mosaicRoute  = "/mosaic/%s"
)

// const routers path for methods BlockchainService
const (
	blockHeightRoute         = "/chain/height"
	blockByHeightRoute       = "/block/%d"
	blockScoreRoute          = "/chain/score"
	blockGetTransactionRoute = "/block/%d/transactions"
	blockInfoRoute           = "/blocks/%d/limit/%d"
	blockStorageRoute        = "/diagnostic/storage"
)

const (
	contractsInfoRoute      = "/contract"
	contractsByAccountRoute = "/account/%s/contracts"
)

// const routers path for methods MosaicService
const (
	networkRoute = "/network"
)

// const routers path for methods TransactionService
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
