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

// namespace id for storage mosaic
var StorageNamespaceId, _ = NewNamespaceIdFromName("prx.so")

// namespace id for streaming mosaic
var StreamingNamespaceId, _ = NewNamespaceIdFromName("prx.sm")

// namespace id for suepr contract mosaic
var SuperContractNamespaceId, _ = NewNamespaceIdFromName("prx.sc")

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
	blockHeightRoute   = "/chain/height"
	blockByHeightRoute = "/block/%s"
	blockScoreRoute    = "/chain/score"
	blockInfoRoute     = "/blocks/%s/limit/%s"
	blockStorageRoute  = "/diagnostic/storage"
)

// routes for ContractsService
const (
	contractsInfoRoute      = "/contract"
	contractsByAccountRoute = "/account/%s/contracts"
)

// routes for LockService
const (
	hashLocksRoute            = "/account/%s/lock/hash"
	secretLocksByAccountRoute = "/account/%s/lock/secret"
	hashLockRoute             = "/lock/hash/%s"
	secretLockRoute           = "/lock/compositeHash/%s"
	secretLocksBySecretRoute  = "/lock/secret/%s"
)

// routes for MetadataService
const (
	metadatasInfoRoute       = "/metadata"
	metadataInfoRoute        = "/metadata/%s"
	metadataByAccountRoute   = "/account/%s/metadata"
	metadataByMosaicRoute    = "/mosaic/%s/metadata"
	metadataByNamespaceRoute = "/namespace/%s/metadata"
)

// routes for MetadataNemService
const (
	metadataEntriesRoute = "/metadata_nem"
	// POST and GET
	metadataEntryHashRoute = "/metadata_nem/%s"
)

// routes for NodeService
const (
	nodeInfoRoute  = "/node/info"
	nodeTimeRoute  = "/node/time"
	nodePeersRoute = "/node/peers"
)

// routes for NetworkService
const (
	networkRoute = "/network"
	configRoute  = "/config/%s"
	upgradeRoute = "/upgrade/%s"
)

// routes for StorageService
const (
	drivesRoute               = "/drives"
	driveRoute                = "/drive/%s"
	drivesOfAccountRoute      = "/account/%s/drive%s"
	downloadInfoRoute         = "/downloads/%s"
	driveDownloadInfosRoute   = "/drive/%s/downloads"
	accountDownloadInfosRoute = "/account/%s/downloads"
)

// routes for SuperContractService
const (
	driveSuperContractsRoute = "/drive/%s/supercontracts"
	superContractRoute       = "/supercontract/%s"
	accountOperationsRoute   = "/account/%s/operations"
	operationRoute           = "/operation/%s"
)

// routes for ExchangeService
const (
	exchangeRoute       = "/account/%s/exchange"
	offersByMosaicRoute = "/exchange/%s/%s"
)

// routes for TransactionService
const (
	transactionsRoute                 = "/transactions"
	transactionsByGroupRoute          = "/transactions/%s"
	transactionsByIdRoute             = "/transactions/%s/%s"
	transactionStatusRoute            = "/transactionStatus"
	transactionStatusByIdRoute        = "/transactionStatus/%s"
	announceAggregateRoute            = "/transactions/partial"
	announceAggregateCosignatureRoute = "/transactions/cosignature"
)

type TransactionGroup string

const (
	Confirmed   TransactionGroup = "confirmed"
	Unconfirmed TransactionGroup = "unconfirmed"
	Partial     TransactionGroup = "partial"
)

type NamespaceType uint8

const (
	Root NamespaceType = iota
	Sub
)

var (
	regValidNamespace = regexp.MustCompile(`^[a-z0-9][a-z0-9-_]*$`)
)
