package sdk

import (
	"bytes"
	"github.com/pkg/errors"
)

func newMessageProcessor(transactionFunc mapTransactionFunc) messageProcessor {
	return &catapultWebSocketMessageProcessor{
		mapTransactionFunc: transactionFunc,
	}
}

type messageProcessor interface {
	ProcessBlock(m []byte) (*BlockInfo, error)
	ProcessConfirmedAdded(m []byte) (Transaction, error)
	ProcessUnconfirmedAdded(m []byte) (Transaction, error)
	ProcessUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error)
	ProcessStatus(m []byte) (*StatusInfo, error)
	ProcessPartialAdded(m []byte) (*AggregateTransaction, error)
	ProcessPartialRemoved(m []byte) (*PartialRemovedInfo, error)
	ProcessCosignature(m []byte) (*SignerInfo, error)
}

type mapTransactionFunc func(b *bytes.Buffer) (Transaction, error)

type catapultWebSocketMessageProcessor struct {
	mapTransactionFunc mapTransactionFunc
}

func (*catapultWebSocketMessageProcessor) ProcessBlock(m []byte) (*BlockInfo, error) {
	dto := &blockInfoDTO{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

func (p *catapultWebSocketMessageProcessor) ProcessConfirmedAdded(m []byte) (Transaction, error) {
	buf := bytes.NewBuffer(m)
	return p.mapTransactionFunc(buf)
}

func (p *catapultWebSocketMessageProcessor) ProcessUnconfirmedAdded(m []byte) (Transaction, error) {
	buf := bytes.NewBuffer(m)
	return p.mapTransactionFunc(buf)
}

func (*catapultWebSocketMessageProcessor) ProcessUnconfirmedRemoved(m []byte) (*UnconfirmedRemoved, error) {
	dto := &unconfirmedRemovedDto{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct(), nil
}

func (*catapultWebSocketMessageProcessor) ProcessStatus(m []byte) (*StatusInfo, error) {
	statusInfo := &StatusInfo{}
	if err := json.Unmarshal(m, statusInfo); err != nil {
		return nil, err
	}

	return statusInfo, nil
}

func (p *catapultWebSocketMessageProcessor) ProcessPartialAdded(m []byte) (*AggregateTransaction, error) {
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

func (p *catapultWebSocketMessageProcessor) ProcessPartialRemoved(m []byte) (*PartialRemovedInfo, error) {
	dto := &partialRemovedInfoDTO{}
	if err := json.Unmarshal(m, dto); err != nil {
		return nil, err
	}

	return dto.toStruct(), nil
}

func (p *catapultWebSocketMessageProcessor) ProcessCosignature(m []byte) (*SignerInfo, error) {
	signerInfo := &SignerInfo{}
	if err := json.Unmarshal(m, signerInfo); err != nil {
		return nil, err
	}

	return signerInfo, nil
}
