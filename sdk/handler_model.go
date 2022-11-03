package sdk

import (
	"fmt"
)

type CompoundChannelHandle struct {
	Address         *Address
	TransactionType *EntityType
}

func (ref *CompoundChannelHandle) String() string {
	if ref.TransactionType != nil {
		return fmt.Sprintf("%x", ref.TransactionType)
	}
	return ref.Address.Address
}

func (ref *CompoundChannelHandle) IsAddress() bool {
	return ref.TransactionType == nil
}

func (ref *CompoundChannelHandle) IsTransactionType() bool {
	return ref.TransactionType != nil
}

func NewCompoundChannelHandleFromAddress(address *Address) *CompoundChannelHandle {
	handle := CompoundChannelHandle{
		Address:         address,
		TransactionType: nil,
	}
	return &handle
}

func NewCompoundChannelHandleFromEntityType(entityType EntityType) *CompoundChannelHandle {
	handle := CompoundChannelHandle{
		Address:         nil,
		TransactionType: &entityType,
	}
	return &handle
}
