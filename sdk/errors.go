// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import "errors"

type RespErr struct {
	msg string
}

func newRespError(msg string) error {
	return &RespErr{msg: msg}
}

func (r *RespErr) Error() string {
	return r.msg
}

// Catapult REST API errors
var (
	ErrResourceNotFound              = newRespError("resource is not found")
	ErrArgumentNotValid              = newRespError("argument is not valid")
	ErrInvalidRequest                = newRespError("request is not valid")
	ErrInternalError                 = newRespError("response is nil")
	ErrNotAcceptedResponseStatusCode = newRespError("not accepted response status code")
)

// Metadata errors
var (
	ErrMetadataEmptyAddresses    = errors.New("list adresses ids must not by empty")
	ErrMetadataNilAdress         = errors.New("adress must not be blank")
	ErrMetadataEmptyMosaicIds    = errors.New("list mosaics ids must not by empty")
	ErrMetadataNilMosaicId       = errors.New("mosaicId must not be nil")
	ErrMetadataEmptyNamespaceIds = errors.New("list namespaces ids must not by empty")
	ErrMetadataNilNamespaceId    = errors.New("namespaceId must not be nil")
)

// Common errors
var (
	ErrNilAssetId             = errors.New("AssetId should not be nil")
	ErrEmptyAssetIds          = errors.New("AssetId's array should not be empty")
	ErrUnknownBlockchainType  = errors.New("Not supported Blockchain Type")
	ErrInvalidHashLength      = errors.New("The length of Hash is invalid")
	ErrInvalidSignatureLength = errors.New("The length of Signature is invalid")
)

// Mosaic errors
var (
	ErrEmptyMosaicIds        = errors.New("list mosaics ids must not by empty")
	ErrNilMosaicId           = errors.New("mosaicId must not be nil")
	ErrWrongBitMosaicId      = errors.New("mosaicId has 64th bit")
	ErrInvalidOwnerPublicKey = errors.New("public owner key is invalid")
	ErrNilMosaicProperties   = errors.New("mosaic properties must not be nil")
)

// Namespace errors
var (
	ErrNamespaceTooManyPart = errors.New("too many parts")
	ErrNilNamespaceId       = errors.New("namespaceId is nil or zero")
	ErrWrongBitNamespaceId  = errors.New("namespaceId doesn't have 64th bit")
	ErrEmptyNamespaceIds    = errors.New("list namespace ids must not by empty")
	ErrInvalidNamespaceName = errors.New("namespace name is invalid")
)

// Blockchain errors
var (
	ErrNilOrZeroHeight = errors.New("block height should not be nil or zero")
	ErrNilOrZeroLimit  = errors.New("limit should not be nil or zero")
)

// Lock errors
var (
	ErrNilSecret = errors.New("Secret should not be nil")
	ErrNilProof  = errors.New("Proof should not be nil")
)

// plain errors
var (
	ErrEmptyAddressesIds = errors.New("list of addresses should not be empty")
	ErrNilAddress        = errors.New("address is nil")
	ErrBlankAddress      = errors.New("address is blank")
	ErrNilAccount        = errors.New("account should not be nil")
	ErrInvalidAddress    = errors.New("wrong address")
	ErrNoChanges    	 = errors.New("transaction should contain changes")
)

// reputations error
var (
	ErrInvalidReputationConfig = errors.New("default reputation should be greater than 0 and less than 1")
)
