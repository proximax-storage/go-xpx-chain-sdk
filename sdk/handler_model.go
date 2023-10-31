package sdk

import (
	"encoding/binary"
	"fmt"
)

type CompoundChannelHandle struct {
	Address         *Address
	TransactionType *EntityType
}

func (ref *CompoundChannelHandle) String() string {
	if ref.TransactionType != nil {
		bytes := make([]byte, 2)
		binary.BigEndian.PutUint16(bytes, uint16(*(ref.TransactionType)))
		return fmt.Sprintf("%02x%02x", bytes[0], bytes[1])
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
