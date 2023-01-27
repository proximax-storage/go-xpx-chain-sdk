package sdk

import (
	"bytes"
	crypto "github.com/proximax-storage/go-xpx-crypto"

	"github.com/pkg/errors"
)

type mapTransactionFunc func(b *bytes.Buffer, generationHash *Hash) (Transaction, error)

//======================================================================================================================

func MapBlock(m []byte) (*BlockInfo, error) {
	dto := &blockInfoDTO{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

type BlockMapper interface {
	MapBlock(m []byte) (*BlockInfo, error)
}

type BlockMapperFn func(m []byte) (*BlockInfo, error)

func (p BlockMapperFn) MapBlock(m []byte) (*BlockInfo, error) {
	return p(m)
}

//======================================================================================================================

func MapReceipt(m []byte) (*AnonymousReceipt, error) {
	dto := &anonymousReceiptDto{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

type ReceiptMapper interface {
	MapReceipt(m []byte) (*AnonymousReceipt, error)
}

type ReceiptMapperFn func(m []byte) (*AnonymousReceipt, error)

func (p ReceiptMapperFn) MapReceipt(m []byte) (*AnonymousReceipt, error) {
	return p(m)
}

//======================================================================================================================

func NewConfirmedAddedMapper(mapTransactionFunc mapTransactionFunc, generationHash *Hash) ConfirmedAddedMapper {
	return &confirmedAddedMapperImpl{
		mapTransactionFunc: mapTransactionFunc,
		generationHash:     generationHash,
	}
}

type ConfirmedAddedMapper interface {
	MapConfirmedAdded(m []byte) (Transaction, error)
}

type confirmedAddedMapperImpl struct {
	mapTransactionFunc mapTransactionFunc
	generationHash     *Hash
}

func (ref *confirmedAddedMapperImpl) MapConfirmedAdded(m []byte) (Transaction, error) {
	buf := bytes.NewBuffer(m)
	return ref.mapTransactionFunc(buf, ref.generationHash)
}

//======================================================================================================================

func NewUnconfirmedAddedMapper(mapTransactionFunc mapTransactionFunc, generationHash *Hash) UnconfirmedAddedMapper {
	return &unconfirmedAddedMapperImpl{
		mapTransactionFunc: mapTransactionFunc,
		generationHash:     generationHash,
	}
}

type UnconfirmedAddedMapper interface {
	MapUnconfirmedAdded(m []byte) (Transaction, error)
}

type unconfirmedAddedMapperImpl struct {
	mapTransactionFunc mapTransactionFunc
	generationHash     *Hash
}

func (p unconfirmedAddedMapperImpl) MapUnconfirmedAdded(m []byte) (Transaction, error) {
	buf := bytes.NewBuffer(m)
	return p.mapTransactionFunc(buf, p.generationHash)
}

//======================================================================================================================

func MapUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error) {
	dto := &unconfirmedRemovedDto{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

type UnconfirmedRemovedMapper interface {
	MapUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error)
}
type UnconfirmedRemovedMapperFn func(m []byte) (*UnconfirmedRemoved, error)

func (p UnconfirmedRemovedMapperFn) MapUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error) {
	return p(m)
}

//======================================================================================================================

func MapStatus(m []byte) (*StatusInfo, error) {
	statusInfoDto := &statusInfoDto{}
	if err := json.Unmarshal(m, statusInfoDto); err != nil {
		return nil, err
	}

	hash, err := statusInfoDto.Hash.Hash()
	if err != nil {
		return nil, err
	}

	return &StatusInfo{
		statusInfoDto.Status,
		hash,
	}, nil
}

type StatusMapper interface {
	MapStatus(m []byte) (*StatusInfo, error)
}

type StatusMapperFn func(m []byte) (*StatusInfo, error)

func (p StatusMapperFn) MapStatus(m []byte) (*StatusInfo, error) {
	return p(m)
}

//======================================================================================================================

func NewPartialAddedMapper(mapTransactionFunc mapTransactionFunc, generationHash *Hash) PartialAddedMapper {
	return &partialAddedMapperImpl{
		mapTransactionFunc: mapTransactionFunc,
		generationHash:     generationHash,
	}
}

type PartialAddedMapper interface {
	MapPartialAdded(m []byte) (Transaction, error)
}

type partialAddedMapperImpl struct {
	mapTransactionFunc mapTransactionFunc
	generationHash     *Hash
}

func (p partialAddedMapperImpl) MapPartialAdded(m []byte) (Transaction, error) {
	buf := bytes.NewBuffer(m)
	tr, err := p.mapTransactionFunc(buf, p.generationHash)
	if err != nil {
		return nil, err
	}

	v, ok := tr.(*AggregateTransactionV2)
	if !ok {
		o, ok := tr.(*AggregateTransactionV1)
		if !ok {
			return nil, errors.New("error cast types")
		}
		return o, nil
	}

	return v, nil
}

//======================================================================================================================

func MapPartialRemoved(m []byte) (*PartialRemovedInfo, error) {
	dto := &partialRemovedInfoDTO{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

type PartialRemovedMapper interface {
	MapPartialRemoved(m []byte) (*PartialRemovedInfo, error)
}

type PartialRemovedMapperFn func(m []byte) (*PartialRemovedInfo, error)

func (p PartialRemovedMapperFn) MapPartialRemoved(m []byte) (*PartialRemovedInfo, error) {
	return p(m)
}

//======================================================================================================================

func MapCosignature(m []byte) (*SignerInfo, error) {
	signerInfoDto := &signerInfoDto{}
	if err := json.Unmarshal(m, signerInfoDto); err != nil {
		return nil, err
	}

	signature, err := signerInfoDto.Signature.Signature()
	if err != nil {
		return nil, err
	}

	parentHash, err := signerInfoDto.ParentHash.Hash()
	if err != nil {
		return nil, err
	}

	return &SignerInfo{
		signerInfoDto.Signer,
		signature,
		crypto.DerivationScheme(signerInfoDto.DerivationScheme),
		parentHash,
	}, nil
}

type CosignatureMapper interface {
	MapCosignature(m []byte) (*SignerInfo, error)
}

type CosignatureMapperFn func(m []byte) (*SignerInfo, error)

func (p CosignatureMapperFn) MapCosignature(m []byte) (*SignerInfo, error) {
	return p(m)
}

//======================================================================================================================

func MapDriveState(m []byte) (*DriveStateInfo, error) {
	driveStateDto := &driveStateDto{}
	if err := json.Unmarshal(m, driveStateDto); err != nil {
		return nil, err
	}

	return driveStateDto.toStruct()
}

type DriveStateMapper interface {
	MapDriveState(m []byte) (*DriveStateInfo, error)
}

type DriveStateMapperFn func(m []byte) (*DriveStateInfo, error)

func (p DriveStateMapperFn) MapDriveState(m []byte) (*DriveStateInfo, error) {
	return p(m)
}
