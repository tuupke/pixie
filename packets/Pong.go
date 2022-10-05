// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package packets

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type PongT struct {
	Identifier string
}

func (t *PongT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil { return 0 }
	identifierOffset := builder.CreateString(t.Identifier)
	PongStart(builder)
	PongAddIdentifier(builder, identifierOffset)
	return PongEnd(builder)
}

func (rcv *Pong) UnPackTo(t *PongT) {
	t.Identifier = string(rcv.Identifier())
}

func (rcv *Pong) UnPack() *PongT {
	if rcv == nil { return nil }
	t := &PongT{}
	rcv.UnPackTo(t)
	return t
}

type Pong struct {
	_tab flatbuffers.Table
}

func GetRootAsPong(buf []byte, offset flatbuffers.UOffsetT) *Pong {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Pong{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsPong(buf []byte, offset flatbuffers.UOffsetT) *Pong {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Pong{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Pong) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Pong) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Pong) Identifier() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func PongStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func PongAddIdentifier(builder *flatbuffers.Builder, identifier flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(identifier), 0)
}
func PongEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}