// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"unsafe"

	"github.com/pkg/errors"

	"github.com/json-iterator/go"
	"github.com/proximax-storage/go-xpx-utils/str"
)

const NamespaceBit uint64 = 1 << 63

type NamespaceId struct {
	baseInt64
}

// returns new NamespaceId from passed namespace identifier
func NewNamespaceId(id uint64) (*NamespaceId, error) {
	if id != 0 && !hasBits(id, NamespaceBit) {
		return nil, ErrWrongBitNamespaceId
	}

	return newNamespaceIdPanic(id), nil
}

// returns new NamespaceId from passed namespace identifier
// TODO
func newNamespaceIdPanic(id uint64) *NamespaceId {
	namespaceId := NamespaceId{baseInt64(id)}
	return &namespaceId
}

func (m *NamespaceId) UnmarshalJSON(data []byte) error {
	var id uint64
	err := binary.Read(bytes.NewBuffer(data[:]), binary.LittleEndian, &id)
	if err != nil {
		return err
	}

	ns, err := NewNamespaceId(id)
	if err != nil {
		return err
	}

	*m = *ns
	return nil
}

func (m *NamespaceId) MarshalJSON() ([]byte, error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, m.Id())
	return data, nil
}

func (m *NamespaceId) Type() AssetIdType {
	return NamespaceAssetIdType
}

func (m *NamespaceId) Id() uint64 {
	return uint64(m.baseInt64)
}

func (m *NamespaceId) String() string {
	return m.toHexString()
}

func (m *NamespaceId) toHexString() string {
	return uint64ToHex(m.Id())
}

func (m *NamespaceId) Equals(id AssetId) (bool, error) {
	if id.Type() != m.Type() {
		return false, errors.New("Mismatch asset types")
	}
	return m.Id() == id.Id(), nil
}

// returns namespace id from passed namespace name
// should be used for creating root, child and grandchild namespace ids
// to create root namespace pass namespace name in format like 'rootname'
// to create child namespace pass namespace name in format like 'rootname.childname'
// to create grand child namespace pass namespace name in format like 'rootname.childname.grandchildname'
func NewNamespaceIdFromName(namespaceName string) (*NamespaceId, error) {
	if list, err := GenerateNamespacePath(namespaceName); err != nil {
		return nil, err
	} else {
		l := len(list)

		if l == 0 {
			return nil, ErrInvalidNamespaceName
		}

		return list[l-1], nil
	}
}

type namespaceIds struct {
	List []*NamespaceId
}

func (ref *namespaceIds) MarshalJSON() (buf []byte, err error) {
	buf = []byte(`{"namespaceIds": [`)

	for i, nsId := range ref.List {
		if i > 0 {
			buf = append(buf, ',')
		}

		buf = append(buf, []byte(`"`+nsId.toHexString()+`"`)...)
	}

	buf = append(buf, ']', '}')

	return
}

func (ref *namespaceIds) IsEmpty(ptr unsafe.Pointer) bool {
	return len((*namespaceIds)(ptr).List) == 0
}

func (ref *namespaceIds) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if (*namespaceIds)(ptr) == nil {
		ptr = (unsafe.Pointer)(&namespaceIds{})
	}

	if iter.ReadNil() {
		*((*unsafe.Pointer)(ptr)) = nil
	} else {
		if iter.WhatIsNext() == jsoniter.ArrayValue {
			iter.Skip()
			newIter := iter.Pool().BorrowIterator([]byte("{}"))
			defer iter.Pool().ReturnIterator(newIter)
			v := newIter.Read()
			list := make([]*NamespaceId, 0)
			for _, val := range v.([]*NamespaceId) {
				list = append(list, val)
			}
			(*namespaceIds)(ptr).List = list
		}
	}
}

func (ref *namespaceIds) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	buf, err := (*namespaceIds)(ptr).MarshalJSON()
	if err == nil {
		_, err = stream.Write(buf)
		//	todo: log error in future
	}

}

// NamespaceAlias contains aliased mosaicId or address and type of alias
type NamespaceAlias struct {
	mosaicId *MosaicId
	address  *Address
	Type     AliasType
}

func (ref *NamespaceAlias) Address() *Address {
	return ref.address
}

func (ref *NamespaceAlias) MosaicId() *MosaicId {
	return ref.mosaicId
}

func (ref *NamespaceAlias) String() string {
	switch ref.Type {
	case AddressAliasType:
		return str.StructToString(
			"NamespaceAlias",
			str.NewField("Address", str.StringPattern, ref.Address()),
			str.NewField("Type", str.IntPattern, ref.Type),
		)
	case MosaicAliasType:
		return str.StructToString(
			"NamespaceAlias",
			str.NewField("MosaicId", str.StringPattern, ref.MosaicId()),
			str.NewField("Type", str.IntPattern, ref.Type),
		)
	}
	return str.StructToString(
		"NamespaceAlias",
		str.NewField("Type", str.IntPattern, ref.Type),
	)
}

type NamespaceInfo struct {
	NamespaceId *NamespaceId
	Active      bool
	TypeSpace   NamespaceType
	Depth       int
	Levels      []*NamespaceId
	Alias       *NamespaceAlias
	Parent      *NamespaceInfo
	Owner       *PublicAccount
	StartHeight Height
	EndHeight   Height
}

func (info NamespaceInfo) String() string {
	return fmt.Sprintf(
		`
			"NamespaceId": %s,
			"Active": %t,
			"TypeSpace": %d,
			"Depth": %d,
			"Levels": %s,
			"Alias": %s,
			"Parent": %v,
			"Owner": %s,
			"StartHeight": %s,
			"EndHeight": %s,
		`,
		info.NamespaceId,
		info.Active,
		info.TypeSpace,
		info.Depth,
		info.Levels,
		info.Alias,
		info.Parent,
		info.Owner,
		info.StartHeight,
		info.EndHeight,
	)
}

type NamespaceName struct {
	NamespaceId *NamespaceId
	FullName    string
}

func (n *NamespaceName) String() string {
	return str.StructToString(
		"NamespaceName",
		str.NewField("NamespaceId", str.StringPattern, n.NamespaceId),
		str.NewField("FullName", str.StringPattern, n.FullName),
	)
}

// returns an array of big ints representation if namespace ids from passed namespace path
// to create root namespace pass namespace name in format like 'rootname'
// to create child namespace pass namespace name in format like 'rootname.childname'
// to create grand child namespace pass namespace name in format like 'rootname.childname.grandchildname'
func GenerateNamespacePath(name string) ([]*NamespaceId, error) {
	parts := strings.Split(name, ".")

	if len(parts) == 0 {
		return nil, ErrInvalidNamespaceName
	}

	if len(parts) > 3 {
		return nil, ErrNamespaceTooManyPart
	}

	var (
		namespaceId = newNamespaceIdPanic(0)
		path        = make([]*NamespaceId, 0)
		err         error
	)

	for _, part := range parts {
		if !regValidNamespace.MatchString(part) {
			return nil, ErrInvalidNamespaceName
		}

		if namespaceId, err = generateNamespaceId(part, namespaceId); err != nil {
			return nil, err
		} else {
			path = append(path, namespaceId)
		}
	}

	return path, nil
}

// returns new Address from namespace identifier
func NewAddressFromNamespace(namespaceId *NamespaceId) (*Address, error) {
	// 0x91 | namespaceId on 8 bytes | 16 bytes 0-pad = 25 bytes
	a := fmt.Sprintf("%X", int(AliasAddress))

	namespaceB := make([]byte, 8)
	binary.LittleEndian.PutUint64(namespaceB, namespaceId.Id())

	a += hex.EncodeToString(namespaceB)
	a += strings.Repeat("00", 16)

	return NewAddressFromBase32(a)
}
