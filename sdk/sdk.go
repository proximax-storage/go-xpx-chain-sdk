// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package sdk provides a client library for the Catapult REST API.
package sdk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/json-iterator/go"
)

const (
	DefaultWebsocketReconnectionTimeout = time.Second * 5
	DefaultFeeCalculationStrategy       = MiddleCalculationStrategy
	DefaultMaxFee                       = 5 * 1000000
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type HttpError struct {
	error
	StatusCode int
}

type FeeCalculationStrategy uint32

// FeeCalculationStrategy enums
const (
	HighCalculationStrategy   FeeCalculationStrategy = 2500
	MiddleCalculationStrategy FeeCalculationStrategy = 250
	LowCalculationStrategy    FeeCalculationStrategy = 25
)

// Provides service configuration
type Config struct {
	reputationConfig      *reputationConfig
	BaseURLs              []*url.URL
	UsedBaseUrl           *url.URL
	WsReconnectionTimeout time.Duration
	GenerationHash        *Hash
	NetworkType
	FeeCalculationStrategy
}

type reputationConfig struct {
	minInteractions   uint64
	defaultReputation float64
}

var defaultRepConfig = reputationConfig{
	minInteractions:   10,
	defaultReputation: 0.9,
}

func NewReputationConfig(minInter uint64, defaultRep float64) (*reputationConfig, error) {
	if defaultRep < 0 || defaultRep > 1 {
		return nil, ErrInvalidReputationConfig
	}

	return &reputationConfig{minInteractions: minInter, defaultReputation: defaultRep}, nil
}

// returns config for HTTP Client from passed node url, filled by information from remote blockchain node
func NewConfig(ctx context.Context, baseUrls []string) (*Config, error) {
	// We want to fill config from remote node(get network type, generationHash of network and etc.).
	// To fill config we need to create a Client. But to create a Client we need to create config =)
	// So we create temporary config with information about connection to server, then create temporary client,
	// which requests information and after that we create a final fully filled config
	tempConf, err := NewConfigWithReputation(
		baseUrls,
		NotSupportedNet,
		&defaultRepConfig,
		DefaultWebsocketReconnectionTimeout,
		nil,
		DefaultFeeCalculationStrategy,
	)

	if err != nil {
		return nil, err
	}

	tempClient := NewClient(nil, tempConf)

	block, err := tempClient.Blockchain.GetBlockByHeight(ctx, Height(1))
	if err != nil {
		return nil, err
	}

	networkType, err := tempClient.Network.GetNetworkType(ctx)
	if err != nil {
		return nil, err
	}

	return NewConfigWithReputation(
		baseUrls,
		networkType,
		&defaultRepConfig,
		DefaultWebsocketReconnectionTimeout,
		block.GenerationHash,
		DefaultFeeCalculationStrategy,
	)
}

func NewConfigWithReputation(
	baseUrls []string,
	networkType NetworkType,
	repConf *reputationConfig,
	wsReconnectionTimeout time.Duration,
	generationHash *Hash,
	strategy FeeCalculationStrategy) (*Config, error) {
	if len(baseUrls) == 0 {
		return nil, errors.New("empty base urls")
	}
	urls := make([]*url.URL, 0, len(baseUrls))

	for _, singleUrlStr := range baseUrls {
		u, err := url.Parse(singleUrlStr)
		if err != nil {
			return nil, err
		}

		urls = append(urls, u)
	}

	c := &Config{
		BaseURLs:               urls,
		UsedBaseUrl:            urls[0],
		WsReconnectionTimeout:  wsReconnectionTimeout,
		NetworkType:            networkType,
		reputationConfig:       repConf,
		GenerationHash:         generationHash,
		FeeCalculationStrategy: strategy,
	}

	return c, nil
}

// Catapult API Client configuration
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.
	config *Config
	common service // Reuse a single struct instead of allocating one for each service on the heap.
	// Services for communicating to the Catapult REST APIs
	Blockchain    *BlockchainService
	Exchange      *ExchangeService
	Mosaic        *MosaicService
	Namespace     *NamespaceService
	Network       *NetworkService
	Transaction   *TransactionService
	Resolve       *ResolverService
	Account       *AccountService
	Storage       *StorageService
	SuperContract *SuperContractService
	Lock          *LockService
	Contract      *ContractService
	Metadata      *MetadataService
}

type service struct {
	client *Client
}

// returns catapult http.Client from passed existing client and configuration
// if passed client is nil, http.DefaultClient will be used
func NewClient(httpClient *http.Client, conf *Config) *Client {
	if httpClient == nil {
		var netTransport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
		}

		httpClient = &http.Client{
			Timeout:   10 * time.Second,
			Transport: netTransport,
		}
	}

	c := &Client{client: httpClient, config: conf}
	c.common.client = c
	c.Blockchain = (*BlockchainService)(&c.common)
	c.Mosaic = (*MosaicService)(&c.common)
	c.Namespace = (*NamespaceService)(&c.common)
	c.Network = &NetworkService{&c.common, c.Blockchain}
	c.Resolve = &ResolverService{&c.common, c.Namespace, c.Mosaic}
	c.Transaction = &TransactionService{&c.common, c.Blockchain}
	c.Exchange = &ExchangeService{&c.common, c.Resolve}
	c.Account = (*AccountService)(&c.common)
	c.Lock = (*LockService)(&c.common)
	c.Storage = &StorageService{&c.common, c.Lock}
	c.SuperContract = (*SuperContractService)(&c.common)
	c.Contract = (*ContractService)(&c.common)
	c.Metadata = (*MetadataService)(&c.common)

	return c
}

func (c *Client) NetworkType() NetworkType {
	return c.config.NetworkType
}

func (c *Client) GenerationHash() *Hash {
	return c.config.GenerationHash
}

//BlockGenerationTime gets value from config. If value not found returns default value - 15s
func (c *Client) BlockGenerationTime(ctx context.Context) (time.Duration, error) {
	cfg, err := c.Network.GetNetworkConfig(ctx)
	if err != nil {
		return 0, err
	}

	if pl, ok := cfg.NetworkConfig.Sections["chain"]; ok {
		if v, ok := pl.Fields["blockGenerationTargetTime"]; ok {
			return time.ParseDuration(v.Value)
		}
	}

	return time.Second * 15, nil
}

// AdaptAccount returns a new account with the same network type and generation hash like a Client
func (c *Client) AdaptAccount(account *Account) (*Account, error) {
	return c.NewAccountFromPrivateKey(account.PrivateKey.String())
}

// doNewRequest creates new request, Do it & return result in V
func (c *Client) doNewRequest(ctx context.Context, method string, path string, body interface{}, v interface{}) (*http.Response, error) {
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req, v)
	if err != nil {
		switch err.(type) {
		case *url.Error:
			for _, url := range c.config.BaseURLs {
				if c.config.UsedBaseUrl == url {
					continue
				}

				req.URL.Host = url.Host
				resp, err = c.do(ctx, req, v)
				if err != nil {
					continue
				}

				c.config.UsedBaseUrl = url
				return resp, nil
			}

			return nil, err
		default:
			return nil, err
		}
	}

	return resp, nil
}

// do sends an API Request and returns a parsed response
func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {

	// set the Context for this request
	req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode > 226 || resp.StatusCode < 200 {
		b := &bytes.Buffer{}
		b.ReadFrom(resp.Body)
		httpError := HttpError{
			fmt.Errorf("sdk do request: %s", b.String()),
			resp.StatusCode,
		}
		return nil, &httpError
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

func (c *Client) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.config.UsedBaseUrl.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("sdk.newRequest config.UsedBaseUrl.Parse: %v", err)
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, fmt.Errorf("sdk.newRequest io.ReadWriter.Encode: %v", err)
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("sdk.newRequest http.NewRequest: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) NewAccount() (*Account, error) {
	return NewAccount(c.config.NetworkType, c.config.GenerationHash)
}

func (c *Client) NewAccountFromPrivateKey(pKey string) (*Account, error) {
	return NewAccountFromPrivateKey(pKey, c.config.NetworkType, c.config.GenerationHash)
}

func (c *Client) NewAccountFromPublicKey(pKey string) (*PublicAccount, error) {
	return NewAccountFromPublicKey(pKey, c.config.NetworkType)
}

// region transactions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (c *Client) modifyTransaction(tx Transaction) {
	// We don't change MaxFee for versioning transactions
	switch tx.GetAbstractTransaction().Type {
	case NetworkConfigEntityType, BlockchainUpgrade:
	default:
		tx.GetAbstractTransaction().MaxFee = Amount(min(tx.Size()*int(c.config.FeeCalculationStrategy), DefaultMaxFee))
	}
}

func (c *Client) NewAddressAliasTransaction(deadline *Deadline, address *Address, namespaceId *NamespaceId, actionType AliasActionType) (*AddressAliasTransaction, error) {
	tx, err := NewAddressAliasTransaction(deadline, address, namespaceId, actionType, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewMosaicAliasTransaction(deadline *Deadline, mosaicId *MosaicId, namespaceId *NamespaceId, actionType AliasActionType) (*MosaicAliasTransaction, error) {
	tx, err := NewMosaicAliasTransaction(deadline, mosaicId, namespaceId, actionType, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewAccountLinkTransaction(deadline *Deadline, remoteAccount *PublicAccount, linkAction AccountLinkAction) (*AccountLinkTransaction, error) {
	tx, err := NewAccountLinkTransaction(deadline, remoteAccount, linkAction, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewAccountPropertiesAddressTransaction(deadline *Deadline, propertyType PropertyType, modifications []*AccountPropertiesAddressModification) (*AccountPropertiesAddressTransaction, error) {
	tx, err := NewAccountPropertiesAddressTransaction(deadline, propertyType, modifications, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewAccountPropertiesMosaicTransaction(deadline *Deadline, propertyType PropertyType, modifications []*AccountPropertiesMosaicModification) (*AccountPropertiesMosaicTransaction, error) {
	tx, err := NewAccountPropertiesMosaicTransaction(deadline, propertyType, modifications, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewAccountPropertiesEntityTypeTransaction(deadline *Deadline, propertyType PropertyType, modifications []*AccountPropertiesEntityTypeModification) (*AccountPropertiesEntityTypeTransaction, error) {
	tx, err := NewAccountPropertiesEntityTypeTransaction(deadline, propertyType, modifications, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewAddExchangeOfferTransaction(deadline *Deadline, addOffers []*AddOffer) (*AddExchangeOfferTransaction, error) {
	tx, err := NewAddExchangeOfferTransaction(deadline, addOffers, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewExchangeOfferTransaction(deadline *Deadline, confirmations []*ExchangeConfirmation) (*ExchangeOfferTransaction, error) {
	tx, err := NewExchangeOfferTransaction(deadline, confirmations, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewRemoveExchangeOfferTransaction(deadline *Deadline, removeOffers []*RemoveOffer) (*RemoveExchangeOfferTransaction, error) {
	tx, err := NewRemoveExchangeOfferTransaction(deadline, removeOffers, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewNetworkConfigTransaction(deadline *Deadline, delta Duration, config *NetworkConfig, entities *SupportedEntities) (*NetworkConfigTransaction, error) {
	tx, err := NewNetworkConfigTransaction(deadline, delta, config, entities, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewBlockchainUpgradeTransaction(deadline *Deadline, upgradePeriod Duration, newBlockChainVersion BlockChainVersion) (*BlockchainUpgradeTransaction, error) {
	tx, err := NewBlockchainUpgradeTransaction(deadline, upgradePeriod, newBlockChainVersion, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewCompleteAggregateTransaction(deadline *Deadline, innerTxs []Transaction) (*AggregateTransaction, error) {
	tx, err := NewCompleteAggregateTransaction(deadline, innerTxs, c.config.NetworkType)
	if err != nil {
		return nil, err
	}
	c.modifyTransaction(tx)

	return tx, tx.UpdateUniqueAggregateHash(c.config.GenerationHash)
}

func (c *Client) NewBondedAggregateTransaction(deadline *Deadline, innerTxs []Transaction) (*AggregateTransaction, error) {
	tx, err := NewBondedAggregateTransaction(deadline, innerTxs, c.config.NetworkType)
	if err != nil {
		return nil, err
	}
	c.modifyTransaction(tx)

	return tx, tx.UpdateUniqueAggregateHash(c.config.GenerationHash)
}

func (c *Client) NewModifyMetadataAddressTransaction(deadline *Deadline, address *Address, modifications []*MetadataModification) (*ModifyMetadataAddressTransaction, error) {
	tx, err := NewModifyMetadataAddressTransaction(deadline, address, modifications, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewModifyMetadataMosaicTransaction(deadline *Deadline, mosaicId *MosaicId, modifications []*MetadataModification) (*ModifyMetadataMosaicTransaction, error) {
	tx, err := NewModifyMetadataMosaicTransaction(deadline, mosaicId, modifications, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewModifyMetadataNamespaceTransaction(deadline *Deadline, namespaceId *NamespaceId, modifications []*MetadataModification) (*ModifyMetadataNamespaceTransaction, error) {
	tx, err := NewModifyMetadataNamespaceTransaction(deadline, namespaceId, modifications, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewModifyMultisigAccountTransaction(deadline *Deadline, minApprovalDelta int8, minRemovalDelta int8, modifications []*MultisigCosignatoryModification) (*ModifyMultisigAccountTransaction, error) {
	tx, err := NewModifyMultisigAccountTransaction(deadline, minApprovalDelta, minRemovalDelta, modifications, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewModifyContractTransaction(
	deadline *Deadline, durationDelta Duration, hash *Hash,
	customers []*MultisigCosignatoryModification,
	executors []*MultisigCosignatoryModification,
	verifiers []*MultisigCosignatoryModification) (*ModifyContractTransaction, error) {
	tx, err := NewModifyContractTransaction(deadline, durationDelta, hash, customers, executors, verifiers, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewMosaicDefinitionTransaction(deadline *Deadline, nonce uint32, ownerPublicKey string, mosaicProps *MosaicProperties) (*MosaicDefinitionTransaction, error) {
	tx, err := NewMosaicDefinitionTransaction(deadline, nonce, ownerPublicKey, mosaicProps, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewMosaicSupplyChangeTransaction(deadline *Deadline, assetId AssetId, supplyType MosaicSupplyType, delta Duration) (*MosaicSupplyChangeTransaction, error) {
	tx, err := NewMosaicSupplyChangeTransaction(deadline, assetId, supplyType, delta, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewTransferTransaction(deadline *Deadline, recipient *Address, mosaics []*Mosaic, message Message) (*TransferTransaction, error) {
	tx, err := NewTransferTransaction(deadline, recipient, mosaics, message, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewTransferTransactionWithNamespace(deadline *Deadline, recipient *NamespaceId, mosaics []*Mosaic, message Message) (*TransferTransaction, error) {
	tx, err := NewTransferTransactionWithNamespace(deadline, recipient, mosaics, message, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewRegisterRootNamespaceTransaction(deadline *Deadline, namespaceName string, duration Duration) (*RegisterNamespaceTransaction, error) {
	tx, err := NewRegisterRootNamespaceTransaction(deadline, namespaceName, duration, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewRegisterSubNamespaceTransaction(deadline *Deadline, namespaceName string, parentId *NamespaceId) (*RegisterNamespaceTransaction, error) {
	tx, err := NewRegisterSubNamespaceTransaction(deadline, namespaceName, parentId, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewLockFundsTransaction(deadline *Deadline, mosaic *Mosaic, duration Duration, signedTx *SignedTransaction) (*LockFundsTransaction, error) {
	tx, err := NewLockFundsTransaction(deadline, mosaic, duration, signedTx, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewSecretLockTransaction(deadline *Deadline, mosaic *Mosaic, duration Duration, secret *Secret, recipient *Address) (*SecretLockTransaction, error) {
	tx, err := NewSecretLockTransaction(deadline, mosaic, duration, secret, recipient, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewSecretProofTransaction(deadline *Deadline, hashType HashType, proof *Proof, recipient *Address) (*SecretProofTransaction, error) {
	tx, err := NewSecretProofTransaction(deadline, hashType, proof, recipient, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}
func (c *Client) NewPrepareDriveTransaction(deadline *Deadline, owner *PublicAccount,
	duration Duration, billingPeriod Duration, billingPrice Amount, driveSize StorageSize,
	replicas uint16, minReplicators uint16, percentApprovers uint8) (*PrepareDriveTransaction, error) {

	tx, err := NewPrepareDriveTransaction(deadline, owner, duration, billingPeriod, billingPrice, driveSize, replicas, minReplicators, percentApprovers, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewJoinToDriveTransaction(deadline *Deadline, driveKey *PublicAccount) (*JoinToDriveTransaction, error) {
	tx, err := NewJoinToDriveTransaction(deadline, driveKey, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewDriveFileSystemTransaction(deadline *Deadline, driveKey string, newRootHash *Hash, oldRootHash *Hash, addActions []*Action, removeActions []*Action) (*DriveFileSystemTransaction, error) {
	tx, err := NewDriveFileSystemTransaction(deadline, driveKey, newRootHash, oldRootHash, addActions, removeActions, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewFilesDepositTransaction(deadline *Deadline, driveKey *PublicAccount, files []*File) (*FilesDepositTransaction, error) {
	tx, err := NewFilesDepositTransaction(deadline, driveKey, files, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewEndDriveTransaction(deadline *Deadline, driveKey *PublicAccount) (*EndDriveTransaction, error) {
	tx, err := NewEndDriveTransaction(deadline, driveKey, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewDriveFilesRewardTransaction(deadline *Deadline, uploadInfos []*UploadInfo) (*DriveFilesRewardTransaction, error) {
	tx, err := NewDriveFilesRewardTransaction(deadline, uploadInfos, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewStartDriveVerificationTransaction(deadline *Deadline, driveKey *PublicAccount) (*StartDriveVerificationTransaction, error) {
	tx, err := NewStartDriveVerificationTransaction(deadline, driveKey, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewEndDriveVerificationTransaction(deadline *Deadline, failures []*FailureVerification) (*EndDriveVerificationTransaction, error) {
	tx, err := NewEndDriveVerificationTransaction(deadline, failures, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewDeployTransaction(deadline *Deadline, drive, owner *PublicAccount, fileHash *Hash, vmVersion uint64) (*DeployTransaction, error) {
	tx, err := NewDeployTransaction(deadline, drive, owner, fileHash, vmVersion, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewStartExecuteTransaction(deadline *Deadline, supercontract *PublicAccount, mosaics []*Mosaic,
	function string, functionParameters []int64) (*StartExecuteTransaction, error) {

	tx, err := NewStartExecuteTransaction(deadline, supercontract, mosaics, function, functionParameters, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewEndExecuteTransaction(deadline *Deadline, mosaics []*Mosaic, token *Hash, status OperationStatus) (*EndExecuteTransaction, error) {
	tx, err := NewEndExecuteTransaction(deadline, mosaics, token, status, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewOperationIdentifyTransaction(deadline *Deadline, hash *Hash) (*OperationIdentifyTransaction, error) {
	tx, err := NewOperationIdentifyTransaction(deadline, hash, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewEndOperationTransaction(deadline *Deadline, mosaics []*Mosaic, token *Hash, status OperationStatus) (*EndOperationTransaction, error) {
	tx, err := NewEndOperationTransaction(deadline, mosaics, token, status, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewStartFileDownloadTransaction(deadline *Deadline, drive *PublicAccount, files []*DownloadFile) (*StartFileDownloadTransaction, error) {
	tx, err := NewStartFileDownloadTransaction(deadline, drive, files, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewEndFileDownloadTransaction(deadline *Deadline, recipient *PublicAccount, operationToken *Hash, files []*DownloadFile) (*EndFileDownloadTransaction, error) {
	tx, err := NewEndFileDownloadTransaction(deadline, recipient, operationToken, files, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewSuperContractFileSystemTransaction(deadline *Deadline, driveKey string, newRootHash *Hash, oldRootHash *Hash, addActions []*Action, removeActions []*Action) (*SuperContractFileSystemTransaction, error) {
	tx, err := NewSuperContractFileSystemTransaction(deadline, driveKey, newRootHash, oldRootHash, addActions, removeActions, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func (c *Client) NewDeactivateTransaction(deadline *Deadline, sc string, driveKey string) (*DeactivateTransaction, error) {
	tx, err := NewDeactivateTransaction(deadline, sc, driveKey, c.config.NetworkType)
	if tx != nil {
		c.modifyTransaction(tx)
	}

	return tx, err
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func handleResponseStatusCode(resp *http.Response, codeToErrs map[int]error) error {
	if resp == nil {
		return ErrInternalError
	}

	if codeToErrs != nil {
		if err, ok := codeToErrs[resp.StatusCode]; ok {
			return err
		}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return ErrNotAcceptedResponseStatusCode
	}

	return nil
}
