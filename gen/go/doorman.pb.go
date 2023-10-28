// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.6
// source: doorman.proto

package _go

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Change struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *Change) Reset() {
	*x = Change{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Change) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Change) ProtoMessage() {}

func (x *Change) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Change.ProtoReflect.Descriptor instead.
func (*Change) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{0}
}

func (x *Change) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type Relation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject string   `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Verb    string   `protobuf:"bytes,2,opt,name=verb,proto3" json:"verb,omitempty"`
	Object  string   `protobuf:"bytes,3,opt,name=object,proto3" json:"object,omitempty"`
	Path    []string `protobuf:"bytes,4,rep,name=path,proto3" json:"path,omitempty"`
}

func (x *Relation) Reset() {
	*x = Relation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Relation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Relation) ProtoMessage() {}

func (x *Relation) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Relation.ProtoReflect.Descriptor instead.
func (*Relation) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{1}
}

func (x *Relation) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *Relation) GetVerb() string {
	if x != nil {
		return x.Verb
	}
	return ""
}

func (x *Relation) GetObject() string {
	if x != nil {
		return x.Object
	}
	return ""
}

func (x *Relation) GetPath() []string {
	if x != nil {
		return x.Path
	}
	return nil
}

type Role struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Verbs []string `protobuf:"bytes,2,rep,name=verbs,proto3" json:"verbs,omitempty"`
}

func (x *Role) Reset() {
	*x = Role{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Role) ProtoMessage() {}

func (x *Role) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Role.ProtoReflect.Descriptor instead.
func (*Role) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{2}
}

func (x *Role) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Role) GetVerbs() []string {
	if x != nil {
		return x.Verbs
	}
	return nil
}

type CheckRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject string `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Verb    string `protobuf:"bytes,2,opt,name=verb,proto3" json:"verb,omitempty"`
	Object  string `protobuf:"bytes,3,opt,name=object,proto3" json:"object,omitempty"`
}

func (x *CheckRequest) Reset() {
	*x = CheckRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckRequest) ProtoMessage() {}

func (x *CheckRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckRequest.ProtoReflect.Descriptor instead.
func (*CheckRequest) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{3}
}

func (x *CheckRequest) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *CheckRequest) GetVerb() string {
	if x != nil {
		return x.Verb
	}
	return ""
}

func (x *CheckRequest) GetObject() string {
	if x != nil {
		return x.Object
	}
	return ""
}

type CheckResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *CheckResponse) Reset() {
	*x = CheckResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckResponse) ProtoMessage() {}

func (x *CheckResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckResponse.ProtoReflect.Descriptor instead.
func (*CheckResponse) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{4}
}

func (x *CheckResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type GrantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject string `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Role    string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	Object  string `protobuf:"bytes,3,opt,name=object,proto3" json:"object,omitempty"`
}

func (x *GrantRequest) Reset() {
	*x = GrantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrantRequest) ProtoMessage() {}

func (x *GrantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrantRequest.ProtoReflect.Descriptor instead.
func (*GrantRequest) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{5}
}

func (x *GrantRequest) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *GrantRequest) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *GrantRequest) GetObject() string {
	if x != nil {
		return x.Object
	}
	return ""
}

type GrantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GrantResponse) Reset() {
	*x = GrantResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrantResponse) ProtoMessage() {}

func (x *GrantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrantResponse.ProtoReflect.Descriptor instead.
func (*GrantResponse) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{6}
}

type RevokeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject string `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Role    string `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
	Object  string `protobuf:"bytes,3,opt,name=object,proto3" json:"object,omitempty"`
}

func (x *RevokeRequest) Reset() {
	*x = RevokeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RevokeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevokeRequest) ProtoMessage() {}

func (x *RevokeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevokeRequest.ProtoReflect.Descriptor instead.
func (*RevokeRequest) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{7}
}

func (x *RevokeRequest) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *RevokeRequest) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *RevokeRequest) GetObject() string {
	if x != nil {
		return x.Object
	}
	return ""
}

type RevokeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RevokeResponse) Reset() {
	*x = RevokeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RevokeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevokeResponse) ProtoMessage() {}

func (x *RevokeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevokeResponse.ProtoReflect.Descriptor instead.
func (*RevokeResponse) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{8}
}

type RemoveRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RemoveRoleRequest) Reset() {
	*x = RemoveRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveRoleRequest) ProtoMessage() {}

func (x *RemoveRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveRoleRequest.ProtoReflect.Descriptor instead.
func (*RemoveRoleRequest) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{9}
}

func (x *RemoveRoleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type UpsertRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Verbs []string `protobuf:"bytes,2,rep,name=verbs,proto3" json:"verbs,omitempty"`
}

func (x *UpsertRoleRequest) Reset() {
	*x = UpsertRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertRoleRequest) ProtoMessage() {}

func (x *UpsertRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertRoleRequest.ProtoReflect.Descriptor instead.
func (*UpsertRoleRequest) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{10}
}

func (x *UpsertRoleRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpsertRoleRequest) GetVerbs() []string {
	if x != nil {
		return x.Verbs
	}
	return nil
}

type ListRelationsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Subject string  `protobuf:"bytes,1,opt,name=subject,proto3" json:"subject,omitempty"`
	Verb    *string `protobuf:"bytes,2,opt,name=verb,proto3,oneof" json:"verb,omitempty"`
}

func (x *ListRelationsRequest) Reset() {
	*x = ListRelationsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRelationsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRelationsRequest) ProtoMessage() {}

func (x *ListRelationsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRelationsRequest.ProtoReflect.Descriptor instead.
func (*ListRelationsRequest) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{11}
}

func (x *ListRelationsRequest) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *ListRelationsRequest) GetVerb() string {
	if x != nil && x.Verb != nil {
		return *x.Verb
	}
	return ""
}

type ListRelationsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items           []*Relation `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	PaginationToken *string     `protobuf:"bytes,2,opt,name=pagination_token,json=paginationToken,proto3,oneof" json:"pagination_token,omitempty"`
}

func (x *ListRelationsResponse) Reset() {
	*x = ListRelationsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRelationsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRelationsResponse) ProtoMessage() {}

func (x *ListRelationsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRelationsResponse.ProtoReflect.Descriptor instead.
func (*ListRelationsResponse) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{12}
}

func (x *ListRelationsResponse) GetItems() []*Relation {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *ListRelationsResponse) GetPaginationToken() string {
	if x != nil && x.PaginationToken != nil {
		return *x.PaginationToken
	}
	return ""
}

type ChangesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type            *string `protobuf:"bytes,1,opt,name=type,proto3,oneof" json:"type,omitempty"`
	PaginationToken *string `protobuf:"bytes,2,opt,name=pagination_token,json=paginationToken,proto3,oneof" json:"pagination_token,omitempty"`
}

func (x *ChangesRequest) Reset() {
	*x = ChangesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangesRequest) ProtoMessage() {}

func (x *ChangesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangesRequest.ProtoReflect.Descriptor instead.
func (*ChangesRequest) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{13}
}

func (x *ChangesRequest) GetType() string {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return ""
}

func (x *ChangesRequest) GetPaginationToken() string {
	if x != nil && x.PaginationToken != nil {
		return *x.PaginationToken
	}
	return ""
}

type ChangesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items           []*Change `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	PaginationToken *string   `protobuf:"bytes,2,opt,name=pagination_token,json=paginationToken,proto3,oneof" json:"pagination_token,omitempty"`
}

func (x *ChangesResponse) Reset() {
	*x = ChangesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_doorman_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangesResponse) ProtoMessage() {}

func (x *ChangesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_doorman_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangesResponse.ProtoReflect.Descriptor instead.
func (*ChangesResponse) Descriptor() ([]byte, []int) {
	return file_doorman_proto_rawDescGZIP(), []int{14}
}

func (x *ChangesResponse) GetItems() []*Change {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *ChangesResponse) GetPaginationToken() string {
	if x != nil && x.PaginationToken != nil {
		return *x.PaginationToken
	}
	return ""
}

var File_doorman_proto protoreflect.FileDescriptor

var file_doorman_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1c, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x22, 0x64, 0x0a, 0x08, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x65,
	0x72, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x76, 0x65, 0x72, 0x62, 0x12, 0x16,
	0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x22, 0x2c, 0x0a, 0x04, 0x52, 0x6f,
	0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x65, 0x72, 0x62, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x65, 0x72, 0x62, 0x73, 0x22, 0x54, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x65, 0x72, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x76, 0x65, 0x72, 0x62, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x29,
	0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x54, 0x0a, 0x0c, 0x47, 0x72, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22,
	0x0f, 0x0a, 0x0d, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x55, 0x0a, 0x0d, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x72,
	0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x22, 0x10, 0x0a, 0x0e, 0x52, 0x65, 0x76, 0x6f, 0x6b,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x23, 0x0a, 0x11, 0x52, 0x65, 0x6d,
	0x6f, 0x76, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x39,
	0x0a, 0x11, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x65, 0x72, 0x62, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x65, 0x72, 0x62, 0x73, 0x22, 0x52, 0x0a, 0x14, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x17, 0x0a, 0x04, 0x76,
	0x65, 0x72, 0x62, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x76, 0x65, 0x72,
	0x62, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x76, 0x65, 0x72, 0x62, 0x22, 0x85, 0x01,
	0x0a, 0x15, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e,
	0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x12, 0x2e, 0x0a, 0x10, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0f, 0x70, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01,
	0x42, 0x13, 0x0a, 0x11, 0x5f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x77, 0x0a, 0x0e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01,
	0x12, 0x2e, 0x0a, 0x10, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0f, 0x70, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01,
	0x42, 0x07, 0x0a, 0x05, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x70, 0x61,
	0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x7d,
	0x0a, 0x0f, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x25, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x2e, 0x0a, 0x10, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x0f, 0x70, 0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x70, 0x61, 0x67,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xc1, 0x04,
	0x0a, 0x07, 0x44, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x12, 0x49, 0x0a, 0x05, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x12, 0x15, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x64, 0x6f, 0x6f, 0x72,
	0x6d, 0x61, 0x6e, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x11, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0b, 0x22, 0x06, 0x2f, 0x63, 0x68, 0x65, 0x63,
	0x6b, 0x3a, 0x01, 0x2a, 0x12, 0x49, 0x0a, 0x05, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x12, 0x15, 0x2e,
	0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x47, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x47,
	0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x11, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x0b, 0x22, 0x06, 0x2f, 0x67, 0x72, 0x61, 0x6e, 0x74, 0x3a, 0x01, 0x2a, 0x12,
	0x4d, 0x0a, 0x06, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x12, 0x16, 0x2e, 0x64, 0x6f, 0x6f, 0x72,
	0x6d, 0x61, 0x6e, 0x2e, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x52, 0x65, 0x76, 0x6f,
	0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x12, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x0c, 0x22, 0x07, 0x2f, 0x72, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x3a, 0x01, 0x2a, 0x12, 0x4c,
	0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x1a, 0x2e, 0x64,
	0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x52, 0x6f, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d,
	0x61, 0x6e, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x2a,
	0x0b, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x4f, 0x0a, 0x0a,
	0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x1a, 0x2e, 0x64, 0x6f, 0x6f,
	0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e,
	0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x1a, 0x0b, 0x2f,
	0x72, 0x6f, 0x6c, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0x62, 0x0a,
	0x0d, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1d,
	0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e,
	0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x12, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x4e, 0x0a, 0x07, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x12, 0x17, 0x2e, 0x64,
	0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2e,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x10, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0a, 0x12, 0x08, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x73, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x74, 0x64, 0x30, 0x6d, 0x2f, 0x64, 0x6f, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x67, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_doorman_proto_rawDescOnce sync.Once
	file_doorman_proto_rawDescData = file_doorman_proto_rawDesc
)

func file_doorman_proto_rawDescGZIP() []byte {
	file_doorman_proto_rawDescOnce.Do(func() {
		file_doorman_proto_rawDescData = protoimpl.X.CompressGZIP(file_doorman_proto_rawDescData)
	})
	return file_doorman_proto_rawDescData
}

var file_doorman_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_doorman_proto_goTypes = []interface{}{
	(*Change)(nil),                // 0: doorman.Change
	(*Relation)(nil),              // 1: doorman.Relation
	(*Role)(nil),                  // 2: doorman.Role
	(*CheckRequest)(nil),          // 3: doorman.CheckRequest
	(*CheckResponse)(nil),         // 4: doorman.CheckResponse
	(*GrantRequest)(nil),          // 5: doorman.GrantRequest
	(*GrantResponse)(nil),         // 6: doorman.GrantResponse
	(*RevokeRequest)(nil),         // 7: doorman.RevokeRequest
	(*RevokeResponse)(nil),        // 8: doorman.RevokeResponse
	(*RemoveRoleRequest)(nil),     // 9: doorman.RemoveRoleRequest
	(*UpsertRoleRequest)(nil),     // 10: doorman.UpsertRoleRequest
	(*ListRelationsRequest)(nil),  // 11: doorman.ListRelationsRequest
	(*ListRelationsResponse)(nil), // 12: doorman.ListRelationsResponse
	(*ChangesRequest)(nil),        // 13: doorman.ChangesRequest
	(*ChangesResponse)(nil),       // 14: doorman.ChangesResponse
}
var file_doorman_proto_depIdxs = []int32{
	1,  // 0: doorman.ListRelationsResponse.items:type_name -> doorman.Relation
	0,  // 1: doorman.ChangesResponse.items:type_name -> doorman.Change
	3,  // 2: doorman.Doorman.Check:input_type -> doorman.CheckRequest
	5,  // 3: doorman.Doorman.Grant:input_type -> doorman.GrantRequest
	7,  // 4: doorman.Doorman.Revoke:input_type -> doorman.RevokeRequest
	9,  // 5: doorman.Doorman.RemoveRole:input_type -> doorman.RemoveRoleRequest
	10, // 6: doorman.Doorman.UpsertRole:input_type -> doorman.UpsertRoleRequest
	11, // 7: doorman.Doorman.ListRelations:input_type -> doorman.ListRelationsRequest
	13, // 8: doorman.Doorman.Changes:input_type -> doorman.ChangesRequest
	4,  // 9: doorman.Doorman.Check:output_type -> doorman.CheckResponse
	6,  // 10: doorman.Doorman.Grant:output_type -> doorman.GrantResponse
	8,  // 11: doorman.Doorman.Revoke:output_type -> doorman.RevokeResponse
	2,  // 12: doorman.Doorman.RemoveRole:output_type -> doorman.Role
	2,  // 13: doorman.Doorman.UpsertRole:output_type -> doorman.Role
	12, // 14: doorman.Doorman.ListRelations:output_type -> doorman.ListRelationsResponse
	14, // 15: doorman.Doorman.Changes:output_type -> doorman.ChangesResponse
	9,  // [9:16] is the sub-list for method output_type
	2,  // [2:9] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_doorman_proto_init() }
func file_doorman_proto_init() {
	if File_doorman_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_doorman_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Change); i {
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
		file_doorman_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Relation); i {
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
		file_doorman_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Role); i {
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
		file_doorman_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckRequest); i {
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
		file_doorman_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckResponse); i {
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
		file_doorman_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrantRequest); i {
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
		file_doorman_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrantResponse); i {
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
		file_doorman_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RevokeRequest); i {
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
		file_doorman_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RevokeResponse); i {
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
		file_doorman_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveRoleRequest); i {
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
		file_doorman_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertRoleRequest); i {
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
		file_doorman_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRelationsRequest); i {
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
		file_doorman_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRelationsResponse); i {
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
		file_doorman_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangesRequest); i {
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
		file_doorman_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangesResponse); i {
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
	file_doorman_proto_msgTypes[11].OneofWrappers = []interface{}{}
	file_doorman_proto_msgTypes[12].OneofWrappers = []interface{}{}
	file_doorman_proto_msgTypes[13].OneofWrappers = []interface{}{}
	file_doorman_proto_msgTypes[14].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_doorman_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_doorman_proto_goTypes,
		DependencyIndexes: file_doorman_proto_depIdxs,
		MessageInfos:      file_doorman_proto_msgTypes,
	}.Build()
	File_doorman_proto = out.File
	file_doorman_proto_rawDesc = nil
	file_doorman_proto_goTypes = nil
	file_doorman_proto_depIdxs = nil
}
