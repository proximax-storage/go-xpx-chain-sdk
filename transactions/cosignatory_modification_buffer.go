// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
// File is Auto-Generated

package transactions

import (
	"github.com/google/flatbuffers/go"
)

type CosignatoryModificationBuffer struct {
	_tab flatbuffers.Table
}

func GetRootAsCosignatoryModificationBuffer(buf []byte, offset flatbuffers.UOffsetT) *CosignatoryModificationBuffer {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &CosignatoryModificationBuffer{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *CosignatoryModificationBuffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *CosignatoryModificationBuffer) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *CosignatoryModificationBuffer) Type() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *CosignatoryModificationBuffer) MutateType(n byte) bool {
	return rcv._tab.MutateByteSlot(4, n)
}

func (rcv *CosignatoryModificationBuffer) CosignatoryPublicKey(j int) byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetByte(a + flatbuffers.UOffsetT(j*1))
	}
	return 0
}

func (rcv *CosignatoryModificationBuffer) CosignatoryPublicKeyLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *CosignatoryModificationBuffer) CosignatoryPublicKeyBytes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *CosignatoryModificationBuffer) MutateCosignatoryPublicKey(j int, n byte) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.MutateByte(a+flatbuffers.UOffsetT(j*1), n)
	}
	return false
}

func CosignatoryModificationBufferStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func CosignatoryModificationBufferAddType(builder *flatbuffers.Builder, type_ byte) {
	builder.PrependByteSlot(0, type_, 0)
}
func CosignatoryModificationBufferAddCosignatoryPublicKey(builder *flatbuffers.Builder, cosignatoryPublicKey flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(cosignatoryPublicKey), 0)
}
func CosignatoryModificationBufferStartCosignatoryPublicKeyVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(1, numElems, 1)
}
func CosignatoryModificationBufferEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
