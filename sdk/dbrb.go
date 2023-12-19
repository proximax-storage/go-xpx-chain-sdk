package sdk

import (
	"errors"
)

type AddDbrbProcessTransaction struct {
	AbstractTransaction
}

func (a *AddDbrbProcessTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &a.AbstractTransaction
}

func (a *AddDbrbProcessTransaction) Size() int {
	return TransactionHeaderSize + 2
}

func (a *AddDbrbProcessTransaction) Bytes() ([]byte, error) {
	return nil, errors.New("cannot get bytes of AddDbrbProcessTransaction")
}
