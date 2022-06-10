// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package packets

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type RegisterT struct {
	Banner *BannerT
	Ips []string
}

func (t *RegisterT) Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	if t == nil { return 0 }
	bannerOffset := t.Banner.Pack(builder)
	ipsOffset := flatbuffers.UOffsetT(0)
	if t.Ips != nil {
		ipsLength := len(t.Ips)
		ipsOffsets := make([]flatbuffers.UOffsetT, ipsLength)
		for j := 0; j < ipsLength; j++ {
			ipsOffsets[j] = builder.CreateString(t.Ips[j])
		}
		RegisterStartIpsVector(builder, ipsLength)
		for j := ipsLength - 1; j >= 0; j-- {
			builder.PrependUOffsetT(ipsOffsets[j])
		}
		ipsOffset = builder.EndVector(ipsLength)
	}
	RegisterStart(builder)
	RegisterAddBanner(builder, bannerOffset)
	RegisterAddIps(builder, ipsOffset)
	return RegisterEnd(builder)
}

func (rcv *Register) UnPackTo(t *RegisterT) {
	t.Banner = rcv.Banner(nil).UnPack()
	ipsLength := rcv.IpsLength()
	t.Ips = make([]string, ipsLength)
	for j := 0; j < ipsLength; j++ {
		t.Ips[j] = string(rcv.Ips(j))
	}
}

func (rcv *Register) UnPack() *RegisterT {
	if rcv == nil { return nil }
	t := &RegisterT{}
	rcv.UnPackTo(t)
	return t
}

type Register struct {
	_tab flatbuffers.Table
}

func GetRootAsRegister(buf []byte, offset flatbuffers.UOffsetT) *Register {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Register{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsRegister(buf []byte, offset flatbuffers.UOffsetT) *Register {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Register{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *Register) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Register) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Register) Banner(obj *Banner) *Banner {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Banner)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *Register) Ips(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j*4))
	}
	return nil
}

func (rcv *Register) IpsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func RegisterStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func RegisterAddBanner(builder *flatbuffers.Builder, banner flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(banner), 0)
}
func RegisterAddIps(builder *flatbuffers.Builder, ips flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(ips), 0)
}
func RegisterStartIpsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(4, numElems, 4)
}
func RegisterEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
