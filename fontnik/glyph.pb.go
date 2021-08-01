// Code generated by protoc-gen-go. DO NOT EDIT.
// source: glyph.proto

package fontnik

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Stores a glyph with metrics and optional SDF bitmap information.
type Glyph struct {
	Id *uint32 `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	// A signed distance field of the glyph with a border of 3 pixels.
	Bitmap []byte `protobuf:"bytes,2,opt,name=bitmap" json:"bitmap,omitempty"`
	// Glyph metrics.
	Width                *uint32  `protobuf:"varint,3,req,name=width" json:"width,omitempty"`
	Height               *uint32  `protobuf:"varint,4,req,name=height" json:"height,omitempty"`
	Left                 *int32   `protobuf:"zigzag32,5,req,name=left" json:"left,omitempty"`
	Top                  *int32   `protobuf:"zigzag32,6,req,name=top" json:"top,omitempty"`
	Advance              *uint32  `protobuf:"varint,7,req,name=advance" json:"advance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Glyph) Reset()         { *m = Glyph{} }
func (m *Glyph) String() string { return proto.CompactTextString(m) }
func (*Glyph) ProtoMessage()    {}
func (*Glyph) Descriptor() ([]byte, []int) {
	return fileDescriptor_f812381304ca7993, []int{0}
}

func (m *Glyph) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Glyph.Unmarshal(m, b)
}
func (m *Glyph) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Glyph.Marshal(b, m, deterministic)
}
func (m *Glyph) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Glyph.Merge(m, src)
}
func (m *Glyph) XXX_Size() int {
	return xxx_messageInfo_Glyph.Size(m)
}
func (m *Glyph) XXX_DiscardUnknown() {
	xxx_messageInfo_Glyph.DiscardUnknown(m)
}

var xxx_messageInfo_Glyph proto.InternalMessageInfo

func (m *Glyph) GetId() uint32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Glyph) GetBitmap() []byte {
	if m != nil {
		return m.Bitmap
	}
	return nil
}

func (m *Glyph) GetWidth() uint32 {
	if m != nil && m.Width != nil {
		return *m.Width
	}
	return 0
}

func (m *Glyph) GetHeight() uint32 {
	if m != nil && m.Height != nil {
		return *m.Height
	}
	return 0
}

func (m *Glyph) GetLeft() int32 {
	if m != nil && m.Left != nil {
		return *m.Left
	}
	return 0
}

func (m *Glyph) GetTop() int32 {
	if m != nil && m.Top != nil {
		return *m.Top
	}
	return 0
}

func (m *Glyph) GetAdvance() uint32 {
	if m != nil && m.Advance != nil {
		return *m.Advance
	}
	return 0
}

// Stores fontstack information and a list of faces.
type Fontstack struct {
	Name                 *string  `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	Range                *string  `protobuf:"bytes,2,req,name=range" json:"range,omitempty"`
	Glyphs               []*Glyph `protobuf:"bytes,3,rep,name=glyphs" json:"glyphs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Fontstack) Reset()         { *m = Fontstack{} }
func (m *Fontstack) String() string { return proto.CompactTextString(m) }
func (*Fontstack) ProtoMessage()    {}
func (*Fontstack) Descriptor() ([]byte, []int) {
	return fileDescriptor_f812381304ca7993, []int{1}
}

func (m *Fontstack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Fontstack.Unmarshal(m, b)
}
func (m *Fontstack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Fontstack.Marshal(b, m, deterministic)
}
func (m *Fontstack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Fontstack.Merge(m, src)
}
func (m *Fontstack) XXX_Size() int {
	return xxx_messageInfo_Fontstack.Size(m)
}
func (m *Fontstack) XXX_DiscardUnknown() {
	xxx_messageInfo_Fontstack.DiscardUnknown(m)
}

var xxx_messageInfo_Fontstack proto.InternalMessageInfo

func (m *Fontstack) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Fontstack) GetRange() string {
	if m != nil && m.Range != nil {
		return *m.Range
	}
	return ""
}

func (m *Fontstack) GetGlyphs() []*Glyph {
	if m != nil {
		return m.Glyphs
	}
	return nil
}

type Glyphs struct {
	Stacks                       []*Fontstack `protobuf:"bytes,1,rep,name=stacks" json:"stacks,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}     `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *Glyphs) Reset()         { *m = Glyphs{} }
func (m *Glyphs) String() string { return proto.CompactTextString(m) }
func (*Glyphs) ProtoMessage()    {}
func (*Glyphs) Descriptor() ([]byte, []int) {
	return fileDescriptor_f812381304ca7993, []int{2}
}

var extRange_Glyphs = []proto.ExtensionRange{
	{Start: 16, End: 8191},
}

func (*Glyphs) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_Glyphs
}

func (m *Glyphs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Glyphs.Unmarshal(m, b)
}
func (m *Glyphs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Glyphs.Marshal(b, m, deterministic)
}
func (m *Glyphs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Glyphs.Merge(m, src)
}
func (m *Glyphs) XXX_Size() int {
	return xxx_messageInfo_Glyphs.Size(m)
}
func (m *Glyphs) XXX_DiscardUnknown() {
	xxx_messageInfo_Glyphs.DiscardUnknown(m)
}

var xxx_messageInfo_Glyphs proto.InternalMessageInfo

func (m *Glyphs) GetStacks() []*Fontstack {
	if m != nil {
		return m.Stacks
	}
	return nil
}

func init() {
	proto.RegisterType((*Glyph)(nil), "fontnik.glyph")
	proto.RegisterType((*Fontstack)(nil), "fontnik.fontstack")
	proto.RegisterType((*Glyphs)(nil), "fontnik.glyphs")
}

func init() { proto.RegisterFile("glyph.proto", fileDescriptor_f812381304ca7993) }

var fileDescriptor_f812381304ca7993 = []byte{
	// 244 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x3c, 0x4f, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x55, 0xec, 0x26, 0xa1, 0x57, 0xa8, 0xc2, 0x09, 0xa1, 0x1b, 0xa3, 0x0c, 0x28, 0xea, 0x90,
	0x81, 0x95, 0x05, 0x31, 0x31, 0x7b, 0x67, 0x30, 0x8d, 0x9b, 0x58, 0x6d, 0x1d, 0xab, 0xb1, 0x40,
	0x6c, 0xfc, 0x08, 0xff, 0x8a, 0x7c, 0x75, 0xbb, 0xbd, 0xf7, 0xee, 0xbd, 0xbb, 0x77, 0xb0, 0x1a,
	0x0e, 0x3f, 0x7e, 0xec, 0xfc, 0x69, 0x0a, 0x13, 0x96, 0xbb, 0xc9, 0x05, 0x67, 0xf7, 0xcd, 0x5f,
	0x06, 0x39, 0x0f, 0x70, 0x0d, 0xc2, 0xf6, 0x94, 0xd5, 0xa2, 0xbd, 0x53, 0xc2, 0xf6, 0xf8, 0x08,
	0xc5, 0xa7, 0x0d, 0x47, 0xed, 0x49, 0xd4, 0x59, 0x7b, 0xab, 0x12, 0xc3, 0x07, 0xc8, 0xbf, 0x6d,
	0x1f, 0x46, 0x92, 0x6c, 0x3d, 0x93, 0xe8, 0x1e, 0x8d, 0x1d, 0xc6, 0x40, 0x0b, 0x96, 0x13, 0x43,
	0x84, 0xc5, 0xc1, 0xec, 0x02, 0xe5, 0xb5, 0x68, 0xef, 0x15, 0x63, 0xac, 0x40, 0x86, 0xc9, 0x53,
	0xc1, 0x52, 0x84, 0x48, 0x50, 0xea, 0xfe, 0x4b, 0xbb, 0xad, 0xa1, 0x92, 0xe3, 0x17, 0xda, 0x7c,
	0xc0, 0x32, 0x56, 0x9d, 0x83, 0xde, 0xee, 0xe3, 0x32, 0xa7, 0x8f, 0x86, 0x4b, 0x2e, 0x15, 0xe3,
	0x58, 0xe7, 0xa4, 0xdd, 0x60, 0x48, 0xb0, 0x78, 0x26, 0xf8, 0x04, 0x05, 0x7f, 0x35, 0x93, 0xac,
	0x65, 0xbb, 0x7a, 0x5e, 0x77, 0xe9, 0xe1, 0x8e, 0x65, 0x95, 0xa6, 0xcd, 0xcb, 0xc5, 0x87, 0x1b,
	0x28, 0xf8, 0xc8, 0x4c, 0x19, 0x27, 0xf0, 0x9a, 0xb8, 0xde, 0x57, 0xc9, 0xb1, 0xc9, 0x6f, 0xaa,
	0xea, 0xf7, 0xf5, 0x4d, 0xbc, 0xcb, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0x45, 0x86, 0xec, 0x07,
	0x56, 0x01, 0x00, 0x00,
}
