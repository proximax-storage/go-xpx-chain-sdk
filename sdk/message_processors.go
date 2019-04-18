package sdk

import (
	"bytes"
	"github.com/pkg/errors"
)

type mapTransactionFunc func(b *bytes.Buffer) (Transaction, error)

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

func NewConfirmedAddedMapper(mapTransactionFunc mapTransactionFunc) ConfirmedAddedMapper {
	return &confirmedAddedMapperImpl{
		mapTransactionFunc: mapTransactionFunc,
	}
}

type ConfirmedAddedMapper interface {
	MapConfirmedAdded(m []byte) (Transaction, error)
}

type confirmedAddedMapperImpl struct {
	mapTransactionFunc mapTransactionFunc
}

func (ref *confirmedAddedMapperImpl) MapConfirmedAdded(m []byte) (Transaction, error) {
	buf := bytes.NewBuffer(m)
	return ref.mapTransactionFunc(buf)
}

//======================================================================================================================

func NewUnconfirmedAddedMapper(mapTransactionFunc mapTransactionFunc) UnconfirmedAddedMapper {
	return &unconfirmedAddedMapperImpl{
		mapTransactionFunc: mapTransactionFunc,
	}
}

type UnconfirmedAddedMapper interface {
	MapUnconfirmedAdded(m []byte) (Transaction, error)
}

type unconfirmedAddedMapperImpl struct {
	mapTransactionFunc mapTransactionFunc
}

func (p unconfirmedAddedMapperImpl) MapUnconfirmedAdded(m []byte) (Transaction, error) {
	buf := bytes.NewBuffer(m)
	return p.mapTransactionFunc(buf)
}

//======================================================================================================================

func MapUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error) {
	dto := &unconfirmedRemovedDto{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct(), nil
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
	statusInfo := &StatusInfo{}
	if err := json.Unmarshal(m, statusInfo); err != nil {
		return nil, err
	}

	return statusInfo, nil
}

type StatusMapper interface {
	MapStatus(m []byte) (*StatusInfo, error)
}

type StatusMapperFn func(m []byte) (*StatusInfo, error)

func (p StatusMapperFn) MapStatus(m []byte) (*StatusInfo, error) {
	return p(m)
}

//======================================================================================================================

func NewPartialAddedMapper(mapTransactionFunc mapTransactionFunc) PartialAddedMapper {
	return &partialAddedMapperImpl{
		mapTransactionFunc: mapTransactionFunc,
	}
}

type PartialAddedMapper interface {
	MapPartialAdded(m []byte) (*AggregateTransaction, error)
}

type partialAddedMapperImpl struct {
	mapTransactionFunc mapTransactionFunc
}

func (p partialAddedMapperImpl) MapPartialAdded(m []byte) (*AggregateTransaction, error) {
	buf := bytes.NewBuffer(m)
	tr, err := p.mapTransactionFunc(buf)
	if err != nil {
		return nil, err
	}

	v, ok := tr.(*AggregateTransaction)
	if !ok {
		return nil, errors.New("error cast types")
	}

	return v, nil
}

//======================================================================================================================

func MapPartialRemoved(m []byte) (*PartialRemovedInfo, error) {
	dto := &partialRemovedInfoDTO{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct(), nil
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
	signerInfo := &SignerInfo{}
	if err := json.Unmarshal(m, signerInfo); err != nil {
		return nil, err
	}

	return signerInfo, nil
}

type CosignatureMapper interface {
	MapCosignature(m []byte) (*SignerInfo, error)
}

type CosignatureMapperFn func(m []byte) (*SignerInfo, error)

func (p CosignatureMapperFn) MapCosignature(m []byte) (*SignerInfo, error) {
	return p(m)
}
