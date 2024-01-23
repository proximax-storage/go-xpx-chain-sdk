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

type RemoveDbrbProcessTransaction struct {
	AbstractTransaction
}

func (a *RemoveDbrbProcessTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &a.AbstractTransaction
}

func (a *RemoveDbrbProcessTransaction) Size() int {
	return TransactionHeaderSize + 2
}

func (a *RemoveDbrbProcessTransaction) Bytes() ([]byte, error) {
	return nil, errors.New("cannot get bytes of RemoveDbrbProcessTransaction")
}
