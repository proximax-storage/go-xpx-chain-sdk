package sdk

import "bytes"

type transactionMapperDependencies struct {
}

func newTransactionMapper(dependencies *transactionMapperDependencies) transactionMapper {
	return &transactionMapperImpl{}
}

type transactionMapper interface {
	MapTransaction(b *bytes.Buffer) (Transaction, error)
	MapTransactions(b *bytes.Buffer) ([]Transaction, error)
}

type transactionMapperImpl struct {
}

func (m *transactionMapperImpl) MapTransaction(b *bytes.Buffer) (Transaction, error) {
	panic("implement me")
}

func (m *transactionMapperImpl) MapTransactions(b *bytes.Buffer) ([]Transaction, error) {
	panic("implement me")
}
