package sdk

import (
	"fmt"
)

type TransactionChannelHandle struct {
	Address         *Address
	TransactionType *EntityType
}

func (ref *TransactionChannelHandle) String() string {
	if ref.TransactionType != nil {
		return fmt.Sprintf("%x", ref.TransactionType)
	}
	return ref.Address.Address
}

func (ref *TransactionChannelHandle) IsAddress() bool {
	return ref.TransactionType == nil
}

func (ref *TransactionChannelHandle) IsTransactionType() bool {
	return ref.TransactionType != nil
}

func NewTransactionChannelHandleFromAddress(address *Address) *TransactionChannelHandle {
	handle := TransactionChannelHandle{
		Address:         address,
		TransactionType: nil,
	}
	return &handle
}

func NewTransactionChannelHandleFromTransactionType(entityType EntityType) *TransactionChannelHandle {
	handle := TransactionChannelHandle{
		Address:         nil,
		TransactionType: &entityType,
	}
	return &handle
}
