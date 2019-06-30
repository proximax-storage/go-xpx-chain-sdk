// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package sdk provides a client library for the Catapult REST API.
package sdk

import (
	"bytes"
	"context"
	"errors"
	"github.com/google/go-querystring/query"
	"github.com/json-iterator/go"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const (
	DefaultNetworkType                  = NotSupportedNet
	DefaultWebsocketReconnectionTimeout = time.Second * 5
	DefaultGenerationHash               = "0000000000000000000000000000000000000000000000000000000000000000"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type HttpError struct {
	error
	StatusCode int
}

// Provides service configuration
type Config struct {
	reputationConfig      *reputationConfig
	BaseURLs              []*url.URL
	UsedBaseUrl           *url.URL
	WsReconnectionTimeout time.Duration
	GenerationHash        *Hash
	NetworkType
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

// returns default config for HTTP Client from passed node url
func NewDefaultConfig(baseUrls []string) (*Config, error) {
	return NewConfigWithReputation(
		baseUrls,
		DefaultNetworkType,
		&defaultRepConfig,
		DefaultWebsocketReconnectionTimeout,
		DefaultGenerationHash,
	)
}

// returns config for HTTP Client from passed node url and network type
func NewConfig(baseUrls []string, networkType NetworkType, wsReconnectionTimeout time.Duration, GenerationHash string) (*Config, error) {
	if wsReconnectionTimeout == 0 {
		wsReconnectionTimeout = DefaultWebsocketReconnectionTimeout
	}

	return NewConfigWithReputation(baseUrls, networkType, &defaultRepConfig, wsReconnectionTimeout, GenerationHash)
}

func NewConfigWithReputation(
	baseUrls []string,
	networkType NetworkType,
	repConf *reputationConfig,
	wsReconnectionTimeout time.Duration,
	GenerationHash string) (*Config, error) {
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

	hash, err := StringToHash(GenerationHash)

	if err != nil {
		return nil, err
	}

	c := &Config{
		BaseURLs:              urls,
		UsedBaseUrl:           urls[0],
		WsReconnectionTimeout: wsReconnectionTimeout,
		NetworkType:           networkType,
		reputationConfig:      repConf,
		GenerationHash:        hash,
	}

	return c, nil
}

// Catapult API Client configuration
type Client struct {
	client *http.Client // HTTP client used to communicate with the API.
	config *Config
	common service // Reuse a single struct instead of allocating one for each service on the heap.
	// Services for communicating to the Catapult REST APIs
	Blockchain  *BlockchainService
	Mosaic      *MosaicService
	Namespace   *NamespaceService
	Network     *NetworkService
	Transaction *TransactionService
	Resolve     *ResolverService
	Account     *AccountService
	Contract    *ContractService
	Metadata    *MetadataService
}

type service struct {
	client *Client
}

// returns catapult http.Client from passed existing client and configuration
// if passed client is nil, http.DefaultClient will be used
func NewClient(httpClient *http.Client, conf *Config) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{client: httpClient, config: conf}
	c.common.client = c
	c.Blockchain = (*BlockchainService)(&c.common)
	c.Mosaic = (*MosaicService)(&c.common)
	c.Namespace = (*NamespaceService)(&c.common)
	c.Network = (*NetworkService)(&c.common)
	c.Resolve = &ResolverService{&c.common, c.Namespace, c.Mosaic}
	c.Transaction = &TransactionService{&c.common, c.Blockchain}
	c.Account = (*AccountService)(&c.common)
	c.Contract = (*ContractService)(&c.common)
	c.Metadata = (*MetadataService)(&c.common)

	return c
}

// NetworkType returns network type of config
func (c *Client) NetworkType() NetworkType {
	return c.config.NetworkType
}

// GenerationHash returns generation hash of config
func (c *Client) GenerationHash() *Hash {
	return c.config.GenerationHash
}

// AdaptAccount returns a new account with the same network type and generation hash like a Client
func (c *Client) AdaptAccount(account *Account) (*Account, error) {
	return c.NewAccountFromPrivateKey(account.PrivateKey.String())
}

// UpdateConfig takes information about network from rest server and updates config of Client
func (c *Client) SetupConfigFromRest(ctx context.Context) error {
	block, err := c.Blockchain.GetBlockByHeight(ctx, Height(1))
	if err != nil {
		return err
	}
	c.config.GenerationHash = block.GenerationHash

	networkType, err := c.Network.GetNetworkType(ctx)
	if err != nil {
		return err
	}
	c.config.NetworkType = networkType

	return nil
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
			errors.New(b.String()),
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
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// returns new Account
func (c *Client) NewAccount() (*Account, error) {
	return NewAccount(c.config.NetworkType, c.config.GenerationHash)
}

// returns new Account from private key
func (c *Client) NewAccountFromPrivateKey(pKey string) (*Account, error) {
	return NewAccountFromPrivateKey(pKey, c.config.NetworkType, c.config.GenerationHash)
}

// returns a PublicAccount from public key
func (c *Client) NewPublicAccountFromPublicKey(pKey string) (*PublicAccount, error) {
	return NewPublicAccountFromPublicKey(pKey, c.config.NetworkType)
}

// region transactions

func (c *Client) NewAddressAliasTransaction(deadline *Deadline, address *Address, namespaceId *NamespaceId, actionType AliasActionType) (*AddressAliasTransaction, error) {
	return NewAddressAliasTransaction(deadline, address, namespaceId, actionType, c.config.NetworkType)
}

func (c *Client) NewMosaicAliasTransaction(deadline *Deadline, mosaicId *MosaicId, namespaceId *NamespaceId, actionType AliasActionType) (*MosaicAliasTransaction, error) {
	return NewMosaicAliasTransaction(deadline, mosaicId, namespaceId, actionType, c.config.NetworkType)
}

func (c *Client) NewAccountLinkTransaction(deadline *Deadline, remoteAccount *PublicAccount, linkAction AccountLinkAction) (*AccountLinkTransaction, error) {
	return NewAccountLinkTransaction(deadline, remoteAccount, linkAction, c.config.NetworkType)
}

func (c *Client) NewAccountPropertiesAddressTransaction(deadline *Deadline, propertyType PropertyType, modifications []*AccountPropertiesAddressModification) (*AccountPropertiesAddressTransaction, error) {
	return NewAccountPropertiesAddressTransaction(deadline, propertyType, modifications, c.config.NetworkType)
}

func (c *Client) NewAccountPropertiesMosaicTransaction(deadline *Deadline, propertyType PropertyType, modifications []*AccountPropertiesMosaicModification) (*AccountPropertiesMosaicTransaction, error) {
	return NewAccountPropertiesMosaicTransaction(deadline, propertyType, modifications, c.config.NetworkType)
}

func (c *Client) NewAccountPropertiesEntityTypeTransaction(deadline *Deadline, propertyType PropertyType, modifications []*AccountPropertiesEntityTypeModification) (*AccountPropertiesEntityTypeTransaction, error) {
	return NewAccountPropertiesEntityTypeTransaction(deadline, propertyType, modifications, c.config.NetworkType)
}

func (c *Client) NewCompleteAggregateTransaction(deadline *Deadline, innerTxs []Transaction) (*AggregateTransaction, error) {
	return NewCompleteAggregateTransaction(deadline, innerTxs, c.config.NetworkType)
}

func (c *Client) NewBondedAggregateTransaction(deadline *Deadline, innerTxs []Transaction) (*AggregateTransaction, error) {
	return NewBondedAggregateTransaction(deadline, innerTxs, c.config.NetworkType)
}

func (c *Client) NewModifyMetadataAddressTransaction(deadline *Deadline, address *Address, modifications []*MetadataModification) (*ModifyMetadataAddressTransaction, error) {
	return NewModifyMetadataAddressTransaction(deadline, address, modifications, c.config.NetworkType)
}

func (c *Client) NewModifyMetadataMosaicTransaction(deadline *Deadline, mosaicId *MosaicId, modifications []*MetadataModification) (*ModifyMetadataMosaicTransaction, error) {
	return NewModifyMetadataMosaicTransaction(deadline, mosaicId, modifications, c.config.NetworkType)
}

func (c *Client) NewModifyMetadataNamespaceTransaction(deadline *Deadline, namespaceId *NamespaceId, modifications []*MetadataModification) (*ModifyMetadataNamespaceTransaction, error) {
	return NewModifyMetadataNamespaceTransaction(deadline, namespaceId, modifications, c.config.NetworkType)
}

func (c *Client) NewModifyMultisigAccountTransaction(deadline *Deadline, minApprovalDelta int8, minRemovalDelta int8, modifications []*MultisigCosignatoryModification) (*ModifyMultisigAccountTransaction, error) {
	return NewModifyMultisigAccountTransaction(deadline, minApprovalDelta, minRemovalDelta, modifications, c.config.NetworkType)
}

func (c *Client) NewModifyContractTransaction(
	deadline *Deadline, durationDelta Duration, hash *Hash,
	customers []*MultisigCosignatoryModification,
	executors []*MultisigCosignatoryModification,
	verifiers []*MultisigCosignatoryModification) (*ModifyContractTransaction, error) {
	return NewModifyContractTransaction(deadline, durationDelta, hash, customers, executors, verifiers, c.config.NetworkType)
}

func (c *Client) NewMosaicDefinitionTransaction(deadline *Deadline, nonce uint32, ownerPublicKey string, mosaicProps *MosaicProperties) (*MosaicDefinitionTransaction, error) {
	return NewMosaicDefinitionTransaction(deadline, nonce, ownerPublicKey, mosaicProps, c.config.NetworkType)
}

func (c *Client) NewMosaicSupplyChangeTransaction(deadline *Deadline, assetId AssetId, supplyType MosaicSupplyType, delta Duration) (*MosaicSupplyChangeTransaction, error) {
	return NewMosaicSupplyChangeTransaction(deadline, assetId, supplyType, delta, c.config.NetworkType)
}

func (c *Client) NewTransferTransaction(deadline *Deadline, recipient *Address, mosaics []*Mosaic, message Message) (*TransferTransaction, error) {
	return NewTransferTransaction(deadline, recipient, mosaics, message, c.config.NetworkType)
}

func (c *Client) NewTransferTransactionWithNamespace(deadline *Deadline, recipient *NamespaceId, mosaics []*Mosaic, message Message) (*TransferTransaction, error) {
	return NewTransferTransactionWithNamespace(deadline, recipient, mosaics, message, c.config.NetworkType)
}

func (c *Client) NewRegisterRootNamespaceTransaction(deadline *Deadline, namespaceName string, duration Duration) (*RegisterNamespaceTransaction, error) {
	return NewRegisterRootNamespaceTransaction(deadline, namespaceName, duration, c.config.NetworkType)
}

func (c *Client) NewRegisterSubNamespaceTransaction(deadline *Deadline, namespaceName string, parentId *NamespaceId) (*RegisterNamespaceTransaction, error) {
	return NewRegisterSubNamespaceTransaction(deadline, namespaceName, parentId, c.config.NetworkType)
}

func (c *Client) NewLockFundsTransaction(deadline *Deadline, mosaic *Mosaic, duration Duration, signedTx *SignedTransaction) (*LockFundsTransaction, error) {
	return NewLockFundsTransaction(deadline, mosaic, duration, signedTx, c.config.NetworkType)
}

func (c *Client) NewSecretLockTransaction(deadline *Deadline, mosaic *Mosaic, duration Duration, secret *Secret, recipient *Address) (*SecretLockTransaction, error) {
	return NewSecretLockTransaction(deadline, mosaic, duration, secret, recipient, c.config.NetworkType)
}

func (c *Client) NewSecretProofTransaction(deadline *Deadline, hashType HashType, proof *Proof, recipient *Address) (*SecretProofTransaction, error) {
	return NewSecretProofTransaction(deadline, hashType, proof, recipient, c.config.NetworkType)
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
