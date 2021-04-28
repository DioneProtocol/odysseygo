// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gsharedmemory.proto

package gsharedmemoryproto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type BatchPut struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatchPut) Reset()         { *m = BatchPut{} }
func (m *BatchPut) String() string { return proto.CompactTextString(m) }
func (*BatchPut) ProtoMessage()    {}
func (*BatchPut) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{0}
}

func (m *BatchPut) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BatchPut.Unmarshal(m, b)
}
func (m *BatchPut) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BatchPut.Marshal(b, m, deterministic)
}
func (m *BatchPut) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchPut.Merge(m, src)
}
func (m *BatchPut) XXX_Size() int {
	return xxx_messageInfo_BatchPut.Size(m)
}
func (m *BatchPut) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchPut.DiscardUnknown(m)
}

var xxx_messageInfo_BatchPut proto.InternalMessageInfo

func (m *BatchPut) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *BatchPut) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type BatchDelete struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BatchDelete) Reset()         { *m = BatchDelete{} }
func (m *BatchDelete) String() string { return proto.CompactTextString(m) }
func (*BatchDelete) ProtoMessage()    {}
func (*BatchDelete) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{1}
}

func (m *BatchDelete) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BatchDelete.Unmarshal(m, b)
}
func (m *BatchDelete) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BatchDelete.Marshal(b, m, deterministic)
}
func (m *BatchDelete) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchDelete.Merge(m, src)
}
func (m *BatchDelete) XXX_Size() int {
	return xxx_messageInfo_BatchDelete.Size(m)
}
func (m *BatchDelete) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchDelete.DiscardUnknown(m)
}

var xxx_messageInfo_BatchDelete proto.InternalMessageInfo

func (m *BatchDelete) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

type Batch struct {
	Puts                 []*BatchPut    `protobuf:"bytes,1,rep,name=puts,proto3" json:"puts,omitempty"`
	Deletes              []*BatchDelete `protobuf:"bytes,2,rep,name=deletes,proto3" json:"deletes,omitempty"`
	Id                   int64          `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Batch) Reset()         { *m = Batch{} }
func (m *Batch) String() string { return proto.CompactTextString(m) }
func (*Batch) ProtoMessage()    {}
func (*Batch) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{2}
}

func (m *Batch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Batch.Unmarshal(m, b)
}
func (m *Batch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Batch.Marshal(b, m, deterministic)
}
func (m *Batch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Batch.Merge(m, src)
}
func (m *Batch) XXX_Size() int {
	return xxx_messageInfo_Batch.Size(m)
}
func (m *Batch) XXX_DiscardUnknown() {
	xxx_messageInfo_Batch.DiscardUnknown(m)
}

var xxx_messageInfo_Batch proto.InternalMessageInfo

func (m *Batch) GetPuts() []*BatchPut {
	if m != nil {
		return m.Puts
	}
	return nil
}

func (m *Batch) GetDeletes() []*BatchDelete {
	if m != nil {
		return m.Deletes
	}
	return nil
}

func (m *Batch) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type Element struct {
	Key                  []byte   `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
	Traits               [][]byte `protobuf:"bytes,3,rep,name=Traits,proto3" json:"Traits,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Element) Reset()         { *m = Element{} }
func (m *Element) String() string { return proto.CompactTextString(m) }
func (*Element) ProtoMessage()    {}
func (*Element) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{3}
}

func (m *Element) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Element.Unmarshal(m, b)
}
func (m *Element) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Element.Marshal(b, m, deterministic)
}
func (m *Element) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Element.Merge(m, src)
}
func (m *Element) XXX_Size() int {
	return xxx_messageInfo_Element.Size(m)
}
func (m *Element) XXX_DiscardUnknown() {
	xxx_messageInfo_Element.DiscardUnknown(m)
}

var xxx_messageInfo_Element proto.InternalMessageInfo

func (m *Element) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Element) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Element) GetTraits() [][]byte {
	if m != nil {
		return m.Traits
	}
	return nil
}

type PutRequest struct {
	PeerChainID          []byte     `protobuf:"bytes,1,opt,name=peerChainID,proto3" json:"peerChainID,omitempty"`
	Elems                []*Element `protobuf:"bytes,2,rep,name=elems,proto3" json:"elems,omitempty"`
	Batches              []*Batch   `protobuf:"bytes,3,rep,name=batches,proto3" json:"batches,omitempty"`
	Id                   int64      `protobuf:"varint,4,opt,name=id,proto3" json:"id,omitempty"`
	Continues            bool       `protobuf:"varint,5,opt,name=continues,proto3" json:"continues,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *PutRequest) Reset()         { *m = PutRequest{} }
func (m *PutRequest) String() string { return proto.CompactTextString(m) }
func (*PutRequest) ProtoMessage()    {}
func (*PutRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{4}
}

func (m *PutRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PutRequest.Unmarshal(m, b)
}
func (m *PutRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PutRequest.Marshal(b, m, deterministic)
}
func (m *PutRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PutRequest.Merge(m, src)
}
func (m *PutRequest) XXX_Size() int {
	return xxx_messageInfo_PutRequest.Size(m)
}
func (m *PutRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PutRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PutRequest proto.InternalMessageInfo

func (m *PutRequest) GetPeerChainID() []byte {
	if m != nil {
		return m.PeerChainID
	}
	return nil
}

func (m *PutRequest) GetElems() []*Element {
	if m != nil {
		return m.Elems
	}
	return nil
}

func (m *PutRequest) GetBatches() []*Batch {
	if m != nil {
		return m.Batches
	}
	return nil
}

func (m *PutRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *PutRequest) GetContinues() bool {
	if m != nil {
		return m.Continues
	}
	return false
}

type PutResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PutResponse) Reset()         { *m = PutResponse{} }
func (m *PutResponse) String() string { return proto.CompactTextString(m) }
func (*PutResponse) ProtoMessage()    {}
func (*PutResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{5}
}

func (m *PutResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PutResponse.Unmarshal(m, b)
}
func (m *PutResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PutResponse.Marshal(b, m, deterministic)
}
func (m *PutResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PutResponse.Merge(m, src)
}
func (m *PutResponse) XXX_Size() int {
	return xxx_messageInfo_PutResponse.Size(m)
}
func (m *PutResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PutResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PutResponse proto.InternalMessageInfo

type GetRequest struct {
	PeerChainID          []byte   `protobuf:"bytes,1,opt,name=peerChainID,proto3" json:"peerChainID,omitempty"`
	Keys                 [][]byte `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{6}
}

func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (m *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(m, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetPeerChainID() []byte {
	if m != nil {
		return m.PeerChainID
	}
	return nil
}

func (m *GetRequest) GetKeys() [][]byte {
	if m != nil {
		return m.Keys
	}
	return nil
}

type GetResponse struct {
	Values               [][]byte `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}
func (*GetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{7}
}

func (m *GetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetResponse.Unmarshal(m, b)
}
func (m *GetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetResponse.Marshal(b, m, deterministic)
}
func (m *GetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetResponse.Merge(m, src)
}
func (m *GetResponse) XXX_Size() int {
	return xxx_messageInfo_GetResponse.Size(m)
}
func (m *GetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetResponse proto.InternalMessageInfo

func (m *GetResponse) GetValues() [][]byte {
	if m != nil {
		return m.Values
	}
	return nil
}

type IndexedRequest struct {
	PeerChainID          []byte   `protobuf:"bytes,1,opt,name=peerChainID,proto3" json:"peerChainID,omitempty"`
	Traits               [][]byte `protobuf:"bytes,2,rep,name=traits,proto3" json:"traits,omitempty"`
	StartTrait           []byte   `protobuf:"bytes,3,opt,name=startTrait,proto3" json:"startTrait,omitempty"`
	StartKey             []byte   `protobuf:"bytes,4,opt,name=startKey,proto3" json:"startKey,omitempty"`
	Limit                int32    `protobuf:"varint,5,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IndexedRequest) Reset()         { *m = IndexedRequest{} }
func (m *IndexedRequest) String() string { return proto.CompactTextString(m) }
func (*IndexedRequest) ProtoMessage()    {}
func (*IndexedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{8}
}

func (m *IndexedRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IndexedRequest.Unmarshal(m, b)
}
func (m *IndexedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IndexedRequest.Marshal(b, m, deterministic)
}
func (m *IndexedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexedRequest.Merge(m, src)
}
func (m *IndexedRequest) XXX_Size() int {
	return xxx_messageInfo_IndexedRequest.Size(m)
}
func (m *IndexedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IndexedRequest proto.InternalMessageInfo

func (m *IndexedRequest) GetPeerChainID() []byte {
	if m != nil {
		return m.PeerChainID
	}
	return nil
}

func (m *IndexedRequest) GetTraits() [][]byte {
	if m != nil {
		return m.Traits
	}
	return nil
}

func (m *IndexedRequest) GetStartTrait() []byte {
	if m != nil {
		return m.StartTrait
	}
	return nil
}

func (m *IndexedRequest) GetStartKey() []byte {
	if m != nil {
		return m.StartKey
	}
	return nil
}

func (m *IndexedRequest) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type IndexedResponse struct {
	Values               [][]byte `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	LastTrait            []byte   `protobuf:"bytes,2,opt,name=lastTrait,proto3" json:"lastTrait,omitempty"`
	LastKey              []byte   `protobuf:"bytes,3,opt,name=lastKey,proto3" json:"lastKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IndexedResponse) Reset()         { *m = IndexedResponse{} }
func (m *IndexedResponse) String() string { return proto.CompactTextString(m) }
func (*IndexedResponse) ProtoMessage()    {}
func (*IndexedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{9}
}

func (m *IndexedResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IndexedResponse.Unmarshal(m, b)
}
func (m *IndexedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IndexedResponse.Marshal(b, m, deterministic)
}
func (m *IndexedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexedResponse.Merge(m, src)
}
func (m *IndexedResponse) XXX_Size() int {
	return xxx_messageInfo_IndexedResponse.Size(m)
}
func (m *IndexedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IndexedResponse proto.InternalMessageInfo

func (m *IndexedResponse) GetValues() [][]byte {
	if m != nil {
		return m.Values
	}
	return nil
}

func (m *IndexedResponse) GetLastTrait() []byte {
	if m != nil {
		return m.LastTrait
	}
	return nil
}

func (m *IndexedResponse) GetLastKey() []byte {
	if m != nil {
		return m.LastKey
	}
	return nil
}

type RemoveRequest struct {
	PeerChainID          []byte   `protobuf:"bytes,1,opt,name=peerChainID,proto3" json:"peerChainID,omitempty"`
	Keys                 [][]byte `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
	Batches              []*Batch `protobuf:"bytes,3,rep,name=batches,proto3" json:"batches,omitempty"`
	Id                   int64    `protobuf:"varint,4,opt,name=id,proto3" json:"id,omitempty"`
	Continues            bool     `protobuf:"varint,5,opt,name=continues,proto3" json:"continues,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveRequest) Reset()         { *m = RemoveRequest{} }
func (m *RemoveRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveRequest) ProtoMessage()    {}
func (*RemoveRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{10}
}

func (m *RemoveRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveRequest.Unmarshal(m, b)
}
func (m *RemoveRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveRequest.Marshal(b, m, deterministic)
}
func (m *RemoveRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveRequest.Merge(m, src)
}
func (m *RemoveRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveRequest.Size(m)
}
func (m *RemoveRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveRequest proto.InternalMessageInfo

func (m *RemoveRequest) GetPeerChainID() []byte {
	if m != nil {
		return m.PeerChainID
	}
	return nil
}

func (m *RemoveRequest) GetKeys() [][]byte {
	if m != nil {
		return m.Keys
	}
	return nil
}

func (m *RemoveRequest) GetBatches() []*Batch {
	if m != nil {
		return m.Batches
	}
	return nil
}

func (m *RemoveRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *RemoveRequest) GetContinues() bool {
	if m != nil {
		return m.Continues
	}
	return false
}

type RemoveResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveResponse) Reset()         { *m = RemoveResponse{} }
func (m *RemoveResponse) String() string { return proto.CompactTextString(m) }
func (*RemoveResponse) ProtoMessage()    {}
func (*RemoveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc30293c358724c5, []int{11}
}

func (m *RemoveResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveResponse.Unmarshal(m, b)
}
func (m *RemoveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveResponse.Marshal(b, m, deterministic)
}
func (m *RemoveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveResponse.Merge(m, src)
}
func (m *RemoveResponse) XXX_Size() int {
	return xxx_messageInfo_RemoveResponse.Size(m)
}
func (m *RemoveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*BatchPut)(nil), "gsharedmemoryproto.BatchPut")
	proto.RegisterType((*BatchDelete)(nil), "gsharedmemoryproto.BatchDelete")
	proto.RegisterType((*Batch)(nil), "gsharedmemoryproto.Batch")
	proto.RegisterType((*Element)(nil), "gsharedmemoryproto.Element")
	proto.RegisterType((*PutRequest)(nil), "gsharedmemoryproto.PutRequest")
	proto.RegisterType((*PutResponse)(nil), "gsharedmemoryproto.PutResponse")
	proto.RegisterType((*GetRequest)(nil), "gsharedmemoryproto.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "gsharedmemoryproto.GetResponse")
	proto.RegisterType((*IndexedRequest)(nil), "gsharedmemoryproto.IndexedRequest")
	proto.RegisterType((*IndexedResponse)(nil), "gsharedmemoryproto.IndexedResponse")
	proto.RegisterType((*RemoveRequest)(nil), "gsharedmemoryproto.RemoveRequest")
	proto.RegisterType((*RemoveResponse)(nil), "gsharedmemoryproto.RemoveResponse")
}

func init() { proto.RegisterFile("gsharedmemory.proto", fileDescriptor_cc30293c358724c5) }

var fileDescriptor_cc30293c358724c5 = []byte{
	// 533 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x96, 0xe3, 0xfc, 0x31, 0x76, 0x43, 0xb5, 0xa0, 0xca, 0x84, 0xaa, 0x31, 0x8b, 0x90, 0x72,
	0x8a, 0x20, 0x3d, 0x71, 0x2d, 0xa5, 0x55, 0x54, 0x21, 0xd0, 0x82, 0xb8, 0xbb, 0xf5, 0x88, 0x58,
	0xf5, 0x4f, 0xf0, 0x6e, 0x2a, 0x72, 0xe7, 0x31, 0x78, 0x00, 0xde, 0x03, 0x1e, 0x0c, 0xed, 0xec,
	0xba, 0x4e, 0x54, 0x27, 0x02, 0x21, 0xf5, 0xe6, 0x6f, 0xf6, 0x9b, 0xf9, 0xe6, 0x9b, 0x99, 0x04,
	0x1e, 0x7d, 0x91, 0xf3, 0xa8, 0xc4, 0x38, 0xc3, 0xac, 0x28, 0x57, 0x93, 0x45, 0x59, 0xa8, 0x82,
	0xb1, 0x8d, 0x20, 0xc5, 0xf8, 0x14, 0xfa, 0x27, 0x91, 0xba, 0x9a, 0x7f, 0x58, 0x2a, 0xb6, 0x0f,
	0xee, 0x35, 0xae, 0x02, 0x27, 0x74, 0xc6, 0xbe, 0xd0, 0x9f, 0xec, 0x31, 0x74, 0x6e, 0xa2, 0x74,
	0x89, 0x41, 0x8b, 0x62, 0x06, 0xf0, 0x11, 0x78, 0x94, 0x73, 0x8a, 0x29, 0x2a, 0xbc, 0x9b, 0xc6,
	0xbf, 0x3b, 0xd0, 0x21, 0x06, 0x7b, 0x09, 0xed, 0xc5, 0x52, 0xc9, 0xc0, 0x09, 0xdd, 0xb1, 0x37,
	0x3d, 0x9c, 0xdc, 0xed, 0x60, 0x52, 0xc9, 0x0b, 0x62, 0xb2, 0xd7, 0xd0, 0x8b, 0xa9, 0xae, 0x0c,
	0x5a, 0x94, 0x34, 0xda, 0x9a, 0x64, 0xf4, 0x45, 0xc5, 0x67, 0x03, 0x68, 0x25, 0x71, 0xe0, 0x86,
	0xce, 0xd8, 0x15, 0xad, 0x24, 0xe6, 0x33, 0xe8, 0xbd, 0x4d, 0x31, 0xc3, 0x9c, 0xac, 0x5d, 0xd4,
	0x3d, 0x5e, 0x18, 0x6b, 0x9f, 0xd7, 0xad, 0x11, 0x60, 0x07, 0xd0, 0xfd, 0x54, 0x46, 0x89, 0x92,
	0x81, 0x1b, 0xba, 0x63, 0x5f, 0x58, 0xc4, 0x7f, 0x3b, 0x00, 0xba, 0x47, 0xfc, 0xba, 0x44, 0xa9,
	0x58, 0x08, 0xde, 0x02, 0xb1, 0x7c, 0x33, 0x8f, 0x92, 0x7c, 0x76, 0x6a, 0xcb, 0xae, 0x87, 0xd8,
	0x2b, 0xe8, 0x60, 0x8a, 0x59, 0x65, 0xe2, 0x69, 0x93, 0x09, 0xdb, 0x9c, 0x30, 0x4c, 0x76, 0x0c,
	0xbd, 0x4b, 0x6d, 0x0b, 0x8d, 0xb8, 0x37, 0x7d, 0xb2, 0xd5, 0xb9, 0xa8, 0x98, 0xd6, 0x73, 0xbb,
	0xf2, 0xcc, 0x0e, 0xe1, 0xc1, 0x55, 0x91, 0xab, 0x24, 0x5f, 0xa2, 0x0c, 0x3a, 0xa1, 0x33, 0xee,
	0x8b, 0x3a, 0xc0, 0xf7, 0xc0, 0x23, 0x17, 0x72, 0x51, 0xe4, 0x12, 0xf9, 0x09, 0xc0, 0x39, 0xfe,
	0x83, 0x29, 0x06, 0xed, 0x6b, 0x5c, 0x19, 0x4f, 0xbe, 0xa0, 0x6f, 0xfe, 0x02, 0x3c, 0xaa, 0x61,
	0x4a, 0xea, 0x01, 0xd2, 0x91, 0x98, 0x95, 0xfb, 0xc2, 0x22, 0xfe, 0xc3, 0x81, 0xc1, 0x2c, 0x8f,
	0xf1, 0x1b, 0xc6, 0x7f, 0xaf, 0x77, 0x00, 0x5d, 0x65, 0xb6, 0x61, 0x14, 0x2d, 0x62, 0x47, 0x00,
	0x52, 0x45, 0xa5, 0xa2, 0xe5, 0xd0, 0xc2, 0x7d, 0xb1, 0x16, 0x61, 0x43, 0xe8, 0x13, 0xd2, 0x2b,
	0x6f, 0xd3, 0xeb, 0x2d, 0xd6, 0x7b, 0x4f, 0x93, 0x2c, 0x51, 0x34, 0x9c, 0x8e, 0x30, 0x80, 0x47,
	0xf0, 0xf0, 0xb6, 0xbb, 0xdd, 0x4e, 0xf4, 0x84, 0xd3, 0x48, 0x5a, 0x6d, 0x73, 0x3c, 0x75, 0x80,
	0x05, 0xd0, 0xd3, 0x40, 0x2b, 0x9b, 0xbe, 0x2a, 0xc8, 0x7f, 0x3a, 0xb0, 0x27, 0x30, 0x2b, 0x6e,
	0xf0, 0xbf, 0x06, 0x7e, 0x1f, 0x67, 0xb2, 0x0f, 0x83, 0xaa, 0x53, 0x33, 0x8c, 0xe9, 0xaf, 0x16,
	0xf8, 0x1f, 0x49, 0xe4, 0x1d, 0x89, 0xb0, 0x33, 0x70, 0xf5, 0x5f, 0xc6, 0x51, 0x93, 0x76, 0xfd,
	0x43, 0x19, 0x8e, 0xb6, 0xbe, 0xdb, 0x29, 0x9f, 0x81, 0x7b, 0x8e, 0x5b, 0xea, 0xd4, 0xb7, 0xd9,
	0x5c, 0x67, 0xfd, 0xee, 0x04, 0xf4, 0xec, 0x02, 0x19, 0x6f, 0xe2, 0x6e, 0xde, 0xde, 0xf0, 0xf9,
	0x4e, 0x8e, 0xad, 0xf9, 0x1e, 0xba, 0x66, 0x0c, 0xec, 0x59, 0x13, 0x7d, 0x63, 0x99, 0x43, 0xbe,
	0x8b, 0x62, 0x0a, 0x5e, 0x76, 0x29, 0x7a, 0xfc, 0x27, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x35, 0x59,
	0x3f, 0x9e, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// SharedMemoryClient is the client API for SharedMemory service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SharedMemoryClient interface {
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Indexed(ctx context.Context, in *IndexedRequest, opts ...grpc.CallOption) (*IndexedResponse, error)
	Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*RemoveResponse, error)
}

type sharedMemoryClient struct {
	cc grpc.ClientConnInterface
}

func NewSharedMemoryClient(cc grpc.ClientConnInterface) SharedMemoryClient {
	return &sharedMemoryClient{cc}
}

func (c *sharedMemoryClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*PutResponse, error) {
	out := new(PutResponse)
	err := c.cc.Invoke(ctx, "/gsharedmemoryproto.SharedMemory/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sharedMemoryClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/gsharedmemoryproto.SharedMemory/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sharedMemoryClient) Indexed(ctx context.Context, in *IndexedRequest, opts ...grpc.CallOption) (*IndexedResponse, error) {
	out := new(IndexedResponse)
	err := c.cc.Invoke(ctx, "/gsharedmemoryproto.SharedMemory/Indexed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sharedMemoryClient) Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*RemoveResponse, error) {
	out := new(RemoveResponse)
	err := c.cc.Invoke(ctx, "/gsharedmemoryproto.SharedMemory/Remove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SharedMemoryServer is the server API for SharedMemory service.
type SharedMemoryServer interface {
	Put(context.Context, *PutRequest) (*PutResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Indexed(context.Context, *IndexedRequest) (*IndexedResponse, error)
	Remove(context.Context, *RemoveRequest) (*RemoveResponse, error)
}

// UnimplementedSharedMemoryServer can be embedded to have forward compatible implementations.
type UnimplementedSharedMemoryServer struct {
}

func (*UnimplementedSharedMemoryServer) Put(ctx context.Context, req *PutRequest) (*PutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (*UnimplementedSharedMemoryServer) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedSharedMemoryServer) Indexed(ctx context.Context, req *IndexedRequest) (*IndexedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Indexed not implemented")
}
func (*UnimplementedSharedMemoryServer) Remove(ctx context.Context, req *RemoveRequest) (*RemoveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Remove not implemented")
}

func RegisterSharedMemoryServer(s *grpc.Server, srv SharedMemoryServer) {
	s.RegisterService(&_SharedMemory_serviceDesc, srv)
}

func _SharedMemory_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SharedMemoryServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gsharedmemoryproto.SharedMemory/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SharedMemoryServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SharedMemory_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SharedMemoryServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gsharedmemoryproto.SharedMemory/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SharedMemoryServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SharedMemory_Indexed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SharedMemoryServer).Indexed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gsharedmemoryproto.SharedMemory/Indexed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SharedMemoryServer).Indexed(ctx, req.(*IndexedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SharedMemory_Remove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SharedMemoryServer).Remove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gsharedmemoryproto.SharedMemory/Remove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SharedMemoryServer).Remove(ctx, req.(*RemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SharedMemory_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gsharedmemoryproto.SharedMemory",
	HandlerType: (*SharedMemoryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _SharedMemory_Put_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _SharedMemory_Get_Handler,
		},
		{
			MethodName: "Indexed",
			Handler:    _SharedMemory_Indexed_Handler,
		},
		{
			MethodName: "Remove",
			Handler:    _SharedMemory_Remove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gsharedmemory.proto",
}
