package sdk

import (
	"bytes"
	"github.com/pkg/errors"
)

type mapTransactionFunc func(b *bytes.Buffer) (Transaction, error)

//======================================================================================================================

func ProcessBlock(m []byte) (*BlockInfo, error) {
	dto := &blockInfoDTO{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

type BlockProcessor interface {
	ProcessBlock(m []byte) (*BlockInfo, error)
}

type BlockProcessorFn func(m []byte) (*BlockInfo, error)

func (p BlockProcessorFn) ProcessBlock(m []byte) (*BlockInfo, error) {
	return p(m)
}

//======================================================================================================================

func NewConfirmedAddedProcessor(mapTransactionFunc mapTransactionFunc) ConfirmedAddedProcessor {
	return &confirmedAddedProcessorImpl{
		mapTransactionFunc: mapTransactionFunc,
	}
}

type ConfirmedAddedProcessor interface {
	ProcessConfirmedAdded(m []byte) (Transaction, error)
}

type confirmedAddedProcessorImpl struct {
	mapTransactionFunc mapTransactionFunc
}

func (ref *confirmedAddedProcessorImpl) ProcessConfirmedAdded(m []byte) (Transaction, error) {
	buf := bytes.NewBuffer(m)
	return ref.mapTransactionFunc(buf)
}

//======================================================================================================================

func NewUnconfirmedAddedProcessor(mapTransactionFunc mapTransactionFunc) UnconfirmedAddedProcessor {
	return &unconfirmedAddedProcessorImpl{
		mapTransactionFunc: mapTransactionFunc,
	}
}

type UnconfirmedAddedProcessor interface {
	ProcessUnconfirmedAdded(m []byte) (Transaction, error)
}

type unconfirmedAddedProcessorImpl struct {
	mapTransactionFunc mapTransactionFunc
}

func (p unconfirmedAddedProcessorImpl) ProcessUnconfirmedAdded(m []byte) (Transaction, error) {
	buf := bytes.NewBuffer(m)
	return p.mapTransactionFunc(buf)
}

//======================================================================================================================

func ProcessUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error) {
	dto := &unconfirmedRemovedDto{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct(), nil
}

type UnconfirmedRemovedProcessor interface {
	ProcessUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error)
}
type UnconfirmedRemovedProcessorFn func(m []byte) (*UnconfirmedRemoved, error)

func (p UnconfirmedRemovedProcessorFn) ProcessUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error) {
	return p(m)
}

//======================================================================================================================

func ProcessStatus(m []byte) (*StatusInfo, error) {
	statusInfo := &StatusInfo{}
	if err := json.Unmarshal(m, statusInfo); err != nil {
		return nil, err
	}

	return statusInfo, nil
}

type StatusProcessor interface {
	ProcessStatus(m []byte) (*StatusInfo, error)
}

type StatusProcessorFn func(m []byte) (*StatusInfo, error)

func (p StatusProcessorFn) ProcessStatus(m []byte) (*StatusInfo, error) {
	return p(m)
}

//======================================================================================================================

func NewPartialAddedProcessor(mapTransactionFunc mapTransactionFunc) PartialAddedProcessor {
	return &partialAddedProcessorImpl{
		mapTransactionFunc: mapTransactionFunc,
	}
}

type PartialAddedProcessor interface {
	ProcessPartialAdded(m []byte) (*AggregateTransaction, error)
}

type partialAddedProcessorImpl struct {
	mapTransactionFunc mapTransactionFunc
}

func (p partialAddedProcessorImpl) ProcessPartialAdded(m []byte) (*AggregateTransaction, error) {
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

func ProcessPartialRemoved(m []byte) (*PartialRemovedInfo, error) {
	dto := &partialRemovedInfoDTO{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct(), nil
}

type PartialRemovedProcessor interface {
	ProcessPartialRemoved(m []byte) (*PartialRemovedInfo, error)
}

type PartialRemovedProcessorFn func(m []byte) (*PartialRemovedInfo, error)

func (p PartialRemovedProcessorFn) ProcessPartialRemoved(m []byte) (*PartialRemovedInfo, error) {
	return p(m)
}

//======================================================================================================================

func ProcessCosignature(m []byte) (*SignerInfo, error) {
	signerInfo := &SignerInfo{}
	if err := json.Unmarshal(m, signerInfo); err != nil {
		return nil, err
	}

	return signerInfo, nil
}

type CosignatureProcessor interface {
	ProcessCosignature(m []byte) (*SignerInfo, error)
}

type CosignatureProcessorFn func(m []byte) (*SignerInfo, error)

func (p CosignatureProcessorFn) ProcessCosignature(m []byte) (*SignerInfo, error) {
	return p(m)
}
