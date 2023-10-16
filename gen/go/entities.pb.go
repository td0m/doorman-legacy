// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.6
// source: entities.proto

package _go

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Entity struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Attrs *structpb.Struct `protobuf:"bytes,2,opt,name=attrs,proto3" json:"attrs,omitempty"`
}

func (x *Entity) Reset() {
	*x = Entity{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entities_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entity) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entity) ProtoMessage() {}

func (x *Entity) ProtoReflect() protoreflect.Message {
	mi := &file_entities_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entity.ProtoReflect.Descriptor instead.
func (*Entity) Descriptor() ([]byte, []int) {
	return file_entities_proto_rawDescGZIP(), []int{0}
}

func (x *Entity) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Entity) GetAttrs() *structpb.Struct {
	if x != nil {
		return x.Attrs
	}
	return nil
}

type EntitiesCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Attrs *structpb.Struct `protobuf:"bytes,2,opt,name=attrs,proto3" json:"attrs,omitempty"`
}

func (x *EntitiesCreateRequest) Reset() {
	*x = EntitiesCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entities_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitiesCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitiesCreateRequest) ProtoMessage() {}

func (x *EntitiesCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_entities_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitiesCreateRequest.ProtoReflect.Descriptor instead.
func (*EntitiesCreateRequest) Descriptor() ([]byte, []int) {
	return file_entities_proto_rawDescGZIP(), []int{1}
}

func (x *EntitiesCreateRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EntitiesCreateRequest) GetAttrs() *structpb.Struct {
	if x != nil {
		return x.Attrs
	}
	return nil
}

type EntitiesRetrieveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *EntitiesRetrieveRequest) Reset() {
	*x = EntitiesRetrieveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entities_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitiesRetrieveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitiesRetrieveRequest) ProtoMessage() {}

func (x *EntitiesRetrieveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_entities_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitiesRetrieveRequest.ProtoReflect.Descriptor instead.
func (*EntitiesRetrieveRequest) Descriptor() ([]byte, []int) {
	return file_entities_proto_rawDescGZIP(), []int{2}
}

func (x *EntitiesRetrieveRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type EntitiesListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PaginationToken *string `protobuf:"bytes,2,opt,name=pagination_token,json=paginationToken,proto3,oneof" json:"pagination_token,omitempty"`
}

func (x *EntitiesListRequest) Reset() {
	*x = EntitiesListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entities_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitiesListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitiesListRequest) ProtoMessage() {}

func (x *EntitiesListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_entities_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitiesListRequest.ProtoReflect.Descriptor instead.
func (*EntitiesListRequest) Descriptor() ([]byte, []int) {
	return file_entities_proto_rawDescGZIP(), []int{3}
}

func (x *EntitiesListRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EntitiesListRequest) GetPaginationToken() string {
	if x != nil && x.PaginationToken != nil {
		return *x.PaginationToken
	}
	return ""
}

type EntitiesListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items           []*Entity `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	PaginationToken *string   `protobuf:"bytes,2,opt,name=pagination_token,json=paginationToken,proto3,oneof" json:"pagination_token,omitempty"`
}

func (x *EntitiesListResponse) Reset() {
	*x = EntitiesListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entities_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitiesListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitiesListResponse) ProtoMessage() {}

func (x *EntitiesListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_entities_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitiesListResponse.ProtoReflect.Descriptor instead.
func (*EntitiesListResponse) Descriptor() ([]byte, []int) {
	return file_entities_proto_rawDescGZIP(), []int{4}
}

func (x *EntitiesListResponse) GetItems() []*Entity {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *EntitiesListResponse) GetPaginationToken() string {
	if x != nil && x.PaginationToken != nil {
		return *x.PaginationToken
	}
	return ""
}

type EntitiesUpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Attrs *structpb.Struct `protobuf:"bytes,2,opt,name=attrs,proto3" json:"attrs,omitempty"`
}

func (x *EntitiesUpdateRequest) Reset() {
	*x = EntitiesUpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entities_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitiesUpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitiesUpdateRequest) ProtoMessage() {}

func (x *EntitiesUpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_entities_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitiesUpdateRequest.ProtoReflect.Descriptor instead.
func (*EntitiesUpdateRequest) Descriptor() ([]byte, []int) {
	return file_entities_proto_rawDescGZIP(), []int{5}
}

func (x *EntitiesUpdateRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EntitiesUpdateRequest) GetAttrs() *structpb.Struct {
	if x != nil {
		return x.Attrs
	}
	return nil
}

type EntitiesDeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *EntitiesDeleteRequest) Reset() {
	*x = EntitiesDeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entities_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitiesDeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitiesDeleteRequest) ProtoMessage() {}

func (x *EntitiesDeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_entities_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitiesDeleteRequest.ProtoReflect.Descriptor instead.
func (*EntitiesDeleteRequest) Descriptor() ([]byte, []int) {
	return file_entities_proto_rawDescGZIP(), []int{6}
}

func (x *EntitiesDeleteRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type EntitiesDeleteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EntitiesDeleteResponse) Reset() {
	*x = EntitiesDeleteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entities_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntitiesDeleteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntitiesDeleteResponse) ProtoMessage() {}

func (x *EntitiesDeleteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_entities_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntitiesDeleteResponse.ProtoReflect.Descriptor instead.
func (*EntitiesDeleteResponse) Descriptor() ([]byte, []int) {
	return file_entities_proto_rawDescGZIP(), []int{7}
}

var File_entities_proto protoreflect.FileDescriptor

var file_entities_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x07, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x47, 0x0a, 0x06, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x2d, 0x0a, 0x05, 0x61, 0x74, 0x74, 0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x05, 0x61, 0x74, 0x74, 0x72, 0x73, 0x22, 0x56,
	0x0a, 0x15, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2d, 0x0a, 0x05, 0x61, 0x74, 0x74, 0x72, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52,
	0x05, 0x61, 0x74, 0x74, 0x72, 0x73, 0x22, 0x29, 0x0a, 0x17, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69,
	0x65, 0x73, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x6a, 0x0a, 0x13, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x10, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x0f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x70, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x82, 0x01,
	0x0a, 0x14, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x2e, 0x0a,
	0x10, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0f, 0x70, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x13, 0x0a,
	0x11, 0x5f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x56, 0x0a, 0x15, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2d, 0x0a, 0x05, 0x61,
	0x74, 0x74, 0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x52, 0x05, 0x61, 0x74, 0x74, 0x72, 0x73, 0x22, 0x27, 0x0a, 0x15, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x18, 0x0a, 0x16, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xc3, 0x03,
	0x0a, 0x08, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x4f, 0x0a, 0x06, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x14, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0e, 0x22, 0x09, 0x2f,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x3a, 0x01, 0x2a, 0x12, 0x55, 0x0a, 0x08, 0x52,
	0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x12, 0x20, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61,
	0x6e, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x64, 0x6f, 0x6f, 0x72,
	0x6d, 0x61, 0x6e, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x10, 0x12, 0x0e, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69,
	0x64, 0x7d, 0x12, 0x56, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1c, 0x2e, 0x64, 0x6f, 0x6f,
	0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d,
	0x61, 0x6e, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x12,
	0x09, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x54, 0x0a, 0x06, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x32, 0x0e, 0x2f,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x3a, 0x01, 0x2a,
	0x12, 0x61, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x1e, 0x2e, 0x64, 0x6f, 0x6f,
	0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x64, 0x6f, 0x6f,
	0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x10, 0x2a, 0x0e, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x7b,
	0x69, 0x64, 0x7d, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x74, 0x64, 0x30, 0x6d, 0x2f, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_entities_proto_rawDescOnce sync.Once
	file_entities_proto_rawDescData = file_entities_proto_rawDesc
)

func file_entities_proto_rawDescGZIP() []byte {
	file_entities_proto_rawDescOnce.Do(func() {
		file_entities_proto_rawDescData = protoimpl.X.CompressGZIP(file_entities_proto_rawDescData)
	})
	return file_entities_proto_rawDescData
}

var file_entities_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_entities_proto_goTypes = []interface{}{
	(*Entity)(nil),                  // 0: doorman.Entity
	(*EntitiesCreateRequest)(nil),   // 1: doorman.EntitiesCreateRequest
	(*EntitiesRetrieveRequest)(nil), // 2: doorman.EntitiesRetrieveRequest
	(*EntitiesListRequest)(nil),     // 3: doorman.EntitiesListRequest
	(*EntitiesListResponse)(nil),    // 4: doorman.EntitiesListResponse
	(*EntitiesUpdateRequest)(nil),   // 5: doorman.EntitiesUpdateRequest
	(*EntitiesDeleteRequest)(nil),   // 6: doorman.EntitiesDeleteRequest
	(*EntitiesDeleteResponse)(nil),  // 7: doorman.EntitiesDeleteResponse
	(*structpb.Struct)(nil),         // 8: google.protobuf.Struct
}
var file_entities_proto_depIdxs = []int32{
	8, // 0: doorman.Entity.attrs:type_name -> google.protobuf.Struct
	8, // 1: doorman.EntitiesCreateRequest.attrs:type_name -> google.protobuf.Struct
	0, // 2: doorman.EntitiesListResponse.items:type_name -> doorman.Entity
	8, // 3: doorman.EntitiesUpdateRequest.attrs:type_name -> google.protobuf.Struct
	1, // 4: doorman.Entities.Create:input_type -> doorman.EntitiesCreateRequest
	2, // 5: doorman.Entities.Retrieve:input_type -> doorman.EntitiesRetrieveRequest
	3, // 6: doorman.Entities.List:input_type -> doorman.EntitiesListRequest
	5, // 7: doorman.Entities.Update:input_type -> doorman.EntitiesUpdateRequest
	6, // 8: doorman.Entities.Delete:input_type -> doorman.EntitiesDeleteRequest
	0, // 9: doorman.Entities.Create:output_type -> doorman.Entity
	0, // 10: doorman.Entities.Retrieve:output_type -> doorman.Entity
	4, // 11: doorman.Entities.List:output_type -> doorman.EntitiesListResponse
	0, // 12: doorman.Entities.Update:output_type -> doorman.Entity
	7, // 13: doorman.Entities.Delete:output_type -> doorman.EntitiesDeleteResponse
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_entities_proto_init() }
func file_entities_proto_init() {
	if File_entities_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_entities_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entity); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_entities_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntitiesCreateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_entities_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntitiesRetrieveRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_entities_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntitiesListRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_entities_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntitiesListResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_entities_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntitiesUpdateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_entities_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntitiesDeleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_entities_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EntitiesDeleteResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_entities_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_entities_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_entities_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_entities_proto_goTypes,
		DependencyIndexes: file_entities_proto_depIdxs,
		MessageInfos:      file_entities_proto_msgTypes,
	}.Build()
	File_entities_proto = out.File
	file_entities_proto_rawDesc = nil
	file_entities_proto_goTypes = nil
	file_entities_proto_depIdxs = nil
}