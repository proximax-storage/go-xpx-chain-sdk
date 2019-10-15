// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"regexp"
)

// namespace id for XEM mosaic
var XemNamespaceId, _ = NewNamespaceIdFromName("nem.xem")

// namespace id for XPX mosaic
var XpxNamespaceId, _ = NewNamespaceIdFromName("prx.xpx")

// routes for AccountService
const (
	accountsRoute                 = "/account"
	accountRoute                  = "/account/%s"
	accountNamesRoute             = "/account/names"
	accountPropertiesRoute        = "/account/%s/properties"
	accountsPropertiesRoute       = "/account/properties"
	multisigAccountRoute          = "/account/%s/multisig"
	multisigAccountGraphInfoRoute = "/account/%s/multisig/graph"
	transactionsByAccountRoute    = "/account/%s/%s"
	accountTransactionsRoute      = "transactions"
	incomingTransactionsRoute     = "transactions/incoming"
	outgoingTransactionsRoute     = "transactions/outgoing"
	unconfirmedTransactionsRoute  = "transactions/unconfirmed"
	aggregateTransactionsRoute    = "transactions/partial"
)

// routes for NamespaceService
const (
	namespaceRoute              = "/namespace/%s"
	namespacesFromAccountsRoute = "/account/namespaces"
	namespaceNamesRoute         = "/namespace/names"
	namespacesFromAccountRoutes = "/account/%s/namespaces"
)

// routes for MosaicService
const (
	mosaicsRoute     = "/mosaic"
	mosaicRoute      = "/mosaic/%s"
	mosaicNamesRoute = "/mosaic/names"
)

// routes for BlockchainService
const (
	blockHeightRoute         = "/chain/height"
	blockByHeightRoute       = "/block/%s"
	blockScoreRoute          = "/chain/score"
	blockGetTransactionRoute = "/block/%s/transactions"
	blockInfoRoute           = "/blocks/%s/limit/%s"
	blockStorageRoute        = "/diagnostic/storage"
)

// routes for ContractsService
const (
	contractsInfoRoute      = "/contract"
	contractsByAccountRoute = "/account/%s/contracts"
)

// routes for MetadataService
const (
	metadatasInfoRoute       = "/metadata"
	metadataInfoRoute        = "/metadata/%s"
	metadataByAccountRoute   = "/account/%s/metadata"
	metadataByMosaicRoute    = "/mosaic/%s/metadata"
	metadataByNamespaceRoute = "/namespace/%s/metadata"
)

// routes for NetworkService
const (
	networkRoute = "/network"
	configRoute  = "/config/%s"
	upgradeRoute = "/upgrade/%s"
)

// routes for TransactionService
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

var (
	regValidNamespace = regexp.MustCompile(`^[a-z0-9][a-z0-9-_]*$`)
)
