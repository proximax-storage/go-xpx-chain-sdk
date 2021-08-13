// Copyright 2021 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
// File is Auto-Generated

package transactions

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type MetadataTransactionBuffer struct {
	_tab flatbuffers.Table
}

func GetRootAsMetadataTransactionBuffer(buf []byte, offset flatbuffers.UOffsetT) *MetadataTransactionBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &MetadataTransactionBuffer{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsMetadataTransactionBuffer(buf []byte, offset flatbuffers.UOffsetT) *MetadataTransactionBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &MetadataTransactionBuffer{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *MetadataTransactionBuffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *MetadataTransactionBuffer) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *MetadataTransactionBuffer) Size() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) MutateSize(n uint32) bool {
	return rcv._tab.MutateUint32Slot(4, n)
}

func (rcv *MetadataTransactionBuffer) Signature(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) SignatureLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) SignatureBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MetadataTransactionBuffer) MutateSignature(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *MetadataTransactionBuffer) Signer(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) SignerLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) SignerBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MetadataTransactionBuffer) MutateSigner(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *MetadataTransactionBuffer) Version() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) MutateVersion(n uint32) bool {
	return rcv._tab.MutateUint32Slot(10, n)
}

func (rcv *MetadataTransactionBuffer) Type() uint16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) MutateType(n uint16) bool {
	return rcv._tab.MutateUint16Slot(12, n)
}

func (rcv *MetadataTransactionBuffer) MaxFee(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) MaxFeeLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) MutateMaxFee(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *MetadataTransactionBuffer) Deadline(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) DeadlineLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) MutateDeadline(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

func (rcv *MetadataTransactionBuffer) TargetKey(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) TargetKeyLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) TargetKeyBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MetadataTransactionBuffer) MutateTargetKey(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *MetadataTransactionBuffer) ScopedMetadataKey(j int) uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetUint32(a + flatbuffers.UOffsetT(j*4))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) ScopedMetadataKeyLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) MutateScopedMetadataKey(j int, n uint32) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateUint32(a+flatbuffers.UOffsetT(j*4), n)
	}
	return false
}

/// In case of address it is empty array. In case of mosaic or namespace id it is 8 byte array(or 2 uint32 array)
func (rcv *MetadataTransactionBuffer) TargetId(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) TargetIdLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) TargetIdBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

/// In case of address it is empty array. In case of mosaic or namespace id it is 8 byte array(or 2 uint32 array)
func (rcv *MetadataTransactionBuffer) MutateTargetId(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *MetadataTransactionBuffer) ValueSizeDelta(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) ValueSizeDeltaLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) ValueSizeDeltaBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MetadataTransactionBuffer) MutateValueSizeDelta(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *MetadataTransactionBuffer) ValueSize(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) ValueSizeLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) ValueSizeBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MetadataTransactionBuffer) MutateValueSize(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func (rcv *MetadataTransactionBuffer) Value(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) ValueLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *MetadataTransactionBuffer) ValueBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *MetadataTransactionBuffer) MutateValue(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func MetadataTransactionBufferStart(builder *flatbuffers.Builder) {
	builder.StartObject(13)
}
func MetadataTransactionBufferAddSize(builder *flatbuffers.Builder, size uint32) {
	builder.PrependUint32Slot(0, size, 0)
}
func MetadataTransactionBufferAddSignature(builder *flatbuffers.Builder, signature flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(signature), 0)
}
func MetadataTransactionBufferStartSignatureVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func MetadataTransactionBufferAddSigner(builder *flatbuffers.Builder, signer flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(signer), 0)
}
func MetadataTransactionBufferStartSignerVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func MetadataTransactionBufferAddVersion(builder *flatbuffers.Builder, version uint32) {
	builder.PrependUint32Slot(3, version, 0)
}
func MetadataTransactionBufferAddType(builder *flatbuffers.Builder, type_ uint16) {
	builder.PrependUint16Slot(4, type_, 0)
}
func MetadataTransactionBufferAddMaxFee(builder *flatbuffers.Builder, maxFee flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(maxFee), 0)
}
func MetadataTransactionBufferStartMaxFeeVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MetadataTransactionBufferAddDeadline(builder *flatbuffers.Builder, deadline flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(deadline), 0)
}
func MetadataTransactionBufferStartDeadlineVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MetadataTransactionBufferAddTargetKey(builder *flatbuffers.Builder, targetKey flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(7, flatbuffers.UOffsetT(targetKey), 0)
}
func MetadataTransactionBufferStartTargetKeyVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func MetadataTransactionBufferAddScopedMetadataKey(builder *flatbuffers.Builder, scopedMetadataKey flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(scopedMetadataKey), 0)
}
func MetadataTransactionBufferStartScopedMetadataKeyVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func MetadataTransactionBufferAddTargetId(builder *flatbuffers.Builder, targetId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(9, flatbuffers.UOffsetT(targetId), 0)
}
func MetadataTransactionBufferStartTargetIdVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func MetadataTransactionBufferAddValueSizeDelta(builder *flatbuffers.Builder, valueSizeDelta flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(10, flatbuffers.UOffsetT(valueSizeDelta), 0)
}
func MetadataTransactionBufferStartValueSizeDeltaVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func MetadataTransactionBufferAddValueSize(builder *flatbuffers.Builder, valueSize flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(11, flatbuffers.UOffsetT(valueSize), 0)
}
func MetadataTransactionBufferStartValueSizeVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func MetadataTransactionBufferAddValue(builder *flatbuffers.Builder, value flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(12, flatbuffers.UOffsetT(value), 0)
}
func MetadataTransactionBufferStartValueVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func MetadataTransactionBufferEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
