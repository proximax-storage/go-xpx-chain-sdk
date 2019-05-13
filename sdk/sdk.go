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
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type HttpError struct {
	error
	StatusCode int
}

// Provides service configuration
type Config struct {
	reputationConfig *reputationConfig
	BaseURL          *url.URL
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

// returns config for HTTP Client from passed node url and NetworkType
func NewConfig(baseUrl string, networkType NetworkType) (*Config, error) {
	return NewConfigWithReputation(baseUrl, networkType, &defaultRepConfig)
}

func NewConfigWithReputation(baseUrl string, networkType NetworkType, repConf *reputationConfig) (*Config, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	c := &Config{BaseURL: u, NetworkType: networkType, reputationConfig: repConf}

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
	c.Transaction = (*TransactionService)(&c.common)
	c.Account = (*AccountService)(&c.common)
	c.Contract = (*ContractService)(&c.common)
	c.Metadata = (*MetadataService)(&c.common)

	return c
}

// doNewRequest creates new request, Do it & return result in V
func (s *Client) doNewRequest(ctx context.Context, method string, path string, body interface{}, v interface{}) (*http.Response, error) {
	req, err := s.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := s.do(ctx, req, v)
	if err != nil {
		return nil, err
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

	u, err := c.config.BaseURL.Parse(urlStr)
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
