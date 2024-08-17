package sdk

import (
	"bytes"

	"github.com/pkg/errors"
)

type MapperFunc[T any] func(generationHash *Hash, payload []byte) (T, error)

type Mapper[T any] interface {
	Map([]byte) (T, error)
}

type mapper[T any] struct {
	fn             MapperFunc[T]
	generationHash *Hash
}

func (m *mapper[T]) Map(payload []byte) (T, error) {
	return m.fn(m.generationHash, payload)
}

func NewMapper[T any](generationHash *Hash, fn MapperFunc[T]) Mapper[T] {
	return &mapper[T]{
		fn:             fn,
		generationHash: generationHash,
	}
}

var TransactionMapperFunc MapperFunc[Transaction] = func(generationHash *Hash, payload []byte) (Transaction, error) {
	buf := bytes.NewBuffer(payload)
	return MapTransaction(buf, generationHash)
}

var BlockMapperFunc MapperFunc[*BlockInfo] = func(generationHash *Hash, payload []byte) (*BlockInfo, error) {
	dto := &blockInfoDTO{}
	if err := json.Unmarshal(payload, dto); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

var CosignatureMapperFunc MapperFunc[*SignerInfo] = func(generationHash *Hash, payload []byte) (*SignerInfo, error) {
	signerInfoDto := &signerInfoDto{}
	if err := json.Unmarshal(payload, signerInfoDto); err != nil {
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
		parentHash,
	}, nil
}

var UnconfirmedRemovedMapperFunc MapperFunc[*UnconfirmedRemoved] = func(generationHash *Hash, payload []byte) (*UnconfirmedRemoved, error) {
	dto := &unconfirmedRemovedDto{}
	if err := json.Unmarshal(payload, dto); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

var DriveStateMapperFunc MapperFunc[*DriveStateInfo] = func(generationHash *Hash, payload []byte) (*DriveStateInfo, error) {
	driveStateDto := &driveStateDto{}
	if err := json.Unmarshal(payload, driveStateDto); err != nil {
		return nil, err
	}

	return driveStateDto.toStruct()
}

var PartialRemovedMapperFunc MapperFunc[*PartialRemovedInfo] = func(generationHash *Hash, payload []byte) (*PartialRemovedInfo, error) {
	dto := &partialRemovedInfoDTO{}
	if err := json.Unmarshal(payload, dto); err != nil {
		return nil, err
	}

	return dto.toStruct()
}

var AggregateTransactionMapperFunc MapperFunc[*AggregateTransaction] = func(generationHash *Hash, payload []byte) (*AggregateTransaction, error) {
	buf := bytes.NewBuffer(payload)
	tr, err := MapTransaction(buf, generationHash)
	if err != nil {
		return nil, err
	}

	v, ok := tr.(*AggregateTransaction)
	if !ok {
		return nil, errors.New("error cast types")
	}

	return v, nil
}

var StatusMapperFunc MapperFunc[*StatusInfo] = func(generationHash *Hash, payload []byte) (*StatusInfo, error) {
	statusInfoDto := &statusInfoDto{}
	if err := json.Unmarshal(payload, statusInfoDto); err != nil {
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
