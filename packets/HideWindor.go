// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package packets

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type HideWindorT struct {
}

func (t *HideWindorT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil { return 0 }
	HideWindorStart(builder)
	return HideWindorEnd(builder)
}

func (rcv *HideWindor) UnPackTo(t *HideWindorT) {
}

func (rcv *HideWindor) UnPack() *HideWindorT {
	if rcv == nil { return nil }
	t := &HideWindorT{}
	rcv.UnPackTo(t)
	return t
}

type HideWindor struct {
	_tab flatbuffers.Table
}

func GetRootAsHideWindor(buf []byte, offset flatbuffers.UOffsetT) *HideWindor {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &HideWindor{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsHideWindor(buf []byte, offset flatbuffers.UOffsetT) *HideWindor {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &HideWindor{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *HideWindor) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *HideWindor) Table() flatbuffers.Table {
	return rcv._tab
}

func HideWindorStart(builder *flatbuffers.Builder) {
	builder.StartObject(0)
}
func HideWindorEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
