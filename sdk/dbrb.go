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

type RemoveDbrbProcessByNetworkTransaction struct {
	AbstractTransaction
}

func (a *RemoveDbrbProcessByNetworkTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &a.AbstractTransaction
}

func (a *RemoveDbrbProcessByNetworkTransaction) Size() int {
	return TransactionHeaderSize + 2
}

func (a *RemoveDbrbProcessByNetworkTransaction) Bytes() ([]byte, error) {
	return nil, errors.New("cannot get bytes of RemoveDbrbProcessByNetworkTransaction")
}

type AddOrUpdateDbrbProcessTransaction struct {
	AbstractTransaction
}

func (a *AddOrUpdateDbrbProcessTransaction) GetAbstractTransaction() *AbstractTransaction {
	return &a.AbstractTransaction
}

func (a *AddOrUpdateDbrbProcessTransaction) Size() int {
	return TransactionHeaderSize + 2
}

func (a *AddOrUpdateDbrbProcessTransaction) Bytes() ([]byte, error) {
	return nil, errors.New("cannot get bytes of AddOrUpdateDbrbProcessTransaction")
}
