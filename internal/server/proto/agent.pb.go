// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/server/proto/agent.proto

package proto

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ExecuteTaskRequest struct {
	Prompt               string            `protobuf:"bytes,1,opt,name=prompt,proto3" json:"prompt,omitempty"`
	Context              map[string]string `protobuf:"bytes,2,rep,name=context,proto3" json:"context,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Tools                []string          `protobuf:"bytes,3,rep,name=tools,proto3" json:"tools,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ExecuteTaskRequest) Reset()         { *m = ExecuteTaskRequest{} }
func (m *ExecuteTaskRequest) String() string { return proto.CompactTextString(m) }
func (*ExecuteTaskRequest) ProtoMessage()    {}
func (*ExecuteTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6286a68aa6b2b2f2, []int{0}
}

func (m *ExecuteTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecuteTaskRequest.Unmarshal(m, b)
}
func (m *ExecuteTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecuteTaskRequest.Marshal(b, m, deterministic)
}
func (m *ExecuteTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuteTaskRequest.Merge(m, src)
}
func (m *ExecuteTaskRequest) XXX_Size() int {
	return xxx_messageInfo_ExecuteTaskRequest.Size(m)
}
func (m *ExecuteTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuteTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuteTaskRequest proto.InternalMessageInfo

func (m *ExecuteTaskRequest) GetPrompt() string {
	if m != nil {
		return m.Prompt
	}
	return ""
}

func (m *ExecuteTaskRequest) GetContext() map[string]string {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *ExecuteTaskRequest) GetTools() []string {
	if m != nil {
		return m.Tools
	}
	return nil
}

type ExecuteTaskResponse struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Message              string   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExecuteTaskResponse) Reset()         { *m = ExecuteTaskResponse{} }
func (m *ExecuteTaskResponse) String() string { return proto.CompactTextString(m) }
func (*ExecuteTaskResponse) ProtoMessage()    {}
func (*ExecuteTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6286a68aa6b2b2f2, []int{1}
}

func (m *ExecuteTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExecuteTaskResponse.Unmarshal(m, b)
}
func (m *ExecuteTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExecuteTaskResponse.Marshal(b, m, deterministic)
}
func (m *ExecuteTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExecuteTaskResponse.Merge(m, src)
}
func (m *ExecuteTaskResponse) XXX_Size() int {
	return xxx_messageInfo_ExecuteTaskResponse.Size(m)
}
func (m *ExecuteTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ExecuteTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ExecuteTaskResponse proto.InternalMessageInfo

func (m *ExecuteTaskResponse) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *ExecuteTaskResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *ExecuteTaskResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type GetTaskStatusRequest struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTaskStatusRequest) Reset()         { *m = GetTaskStatusRequest{} }
func (m *GetTaskStatusRequest) String() string { return proto.CompactTextString(m) }
func (*GetTaskStatusRequest) ProtoMessage()    {}
func (*GetTaskStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6286a68aa6b2b2f2, []int{2}
}

func (m *GetTaskStatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTaskStatusRequest.Unmarshal(m, b)
}
func (m *GetTaskStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTaskStatusRequest.Marshal(b, m, deterministic)
}
func (m *GetTaskStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTaskStatusRequest.Merge(m, src)
}
func (m *GetTaskStatusRequest) XXX_Size() int {
	return xxx_messageInfo_GetTaskStatusRequest.Size(m)
}
func (m *GetTaskStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTaskStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTaskStatusRequest proto.InternalMessageInfo

func (m *GetTaskStatusRequest) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

type GetTaskStatusResponse struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Result               string   `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
	Events               []string `protobuf:"bytes,4,rep,name=events,proto3" json:"events,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTaskStatusResponse) Reset()         { *m = GetTaskStatusResponse{} }
func (m *GetTaskStatusResponse) String() string { return proto.CompactTextString(m) }
func (*GetTaskStatusResponse) ProtoMessage()    {}
func (*GetTaskStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6286a68aa6b2b2f2, []int{3}
}

func (m *GetTaskStatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTaskStatusResponse.Unmarshal(m, b)
}
func (m *GetTaskStatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTaskStatusResponse.Marshal(b, m, deterministic)
}
func (m *GetTaskStatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTaskStatusResponse.Merge(m, src)
}
func (m *GetTaskStatusResponse) XXX_Size() int {
	return xxx_messageInfo_GetTaskStatusResponse.Size(m)
}
func (m *GetTaskStatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTaskStatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetTaskStatusResponse proto.InternalMessageInfo

func (m *GetTaskStatusResponse) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *GetTaskStatusResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *GetTaskStatusResponse) GetResult() string {
	if m != nil {
		return m.Result
	}
	return ""
}

func (m *GetTaskStatusResponse) GetEvents() []string {
	if m != nil {
		return m.Events
	}
	return nil
}

type CancelTaskRequest struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelTaskRequest) Reset()         { *m = CancelTaskRequest{} }
func (m *CancelTaskRequest) String() string { return proto.CompactTextString(m) }
func (*CancelTaskRequest) ProtoMessage()    {}
func (*CancelTaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6286a68aa6b2b2f2, []int{4}
}

func (m *CancelTaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancelTaskRequest.Unmarshal(m, b)
}
func (m *CancelTaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancelTaskRequest.Marshal(b, m, deterministic)
}
func (m *CancelTaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancelTaskRequest.Merge(m, src)
}
func (m *CancelTaskRequest) XXX_Size() int {
	return xxx_messageInfo_CancelTaskRequest.Size(m)
}
func (m *CancelTaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CancelTaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CancelTaskRequest proto.InternalMessageInfo

func (m *CancelTaskRequest) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

type CancelTaskResponse struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Message              string   `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelTaskResponse) Reset()         { *m = CancelTaskResponse{} }
func (m *CancelTaskResponse) String() string { return proto.CompactTextString(m) }
func (*CancelTaskResponse) ProtoMessage()    {}
func (*CancelTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6286a68aa6b2b2f2, []int{5}
}

func (m *CancelTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancelTaskResponse.Unmarshal(m, b)
}
func (m *CancelTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancelTaskResponse.Marshal(b, m, deterministic)
}
func (m *CancelTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancelTaskResponse.Merge(m, src)
}
func (m *CancelTaskResponse) XXX_Size() int {
	return xxx_messageInfo_CancelTaskResponse.Size(m)
}
func (m *CancelTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CancelTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CancelTaskResponse proto.InternalMessageInfo

func (m *CancelTaskResponse) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *CancelTaskResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *CancelTaskResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*ExecuteTaskRequest)(nil), "agent.ExecuteTaskRequest")
	proto.RegisterMapType((map[string]string)(nil), "agent.ExecuteTaskRequest.ContextEntry")
	proto.RegisterType((*ExecuteTaskResponse)(nil), "agent.ExecuteTaskResponse")
	proto.RegisterType((*GetTaskStatusRequest)(nil), "agent.GetTaskStatusRequest")
	proto.RegisterType((*GetTaskStatusResponse)(nil), "agent.GetTaskStatusResponse")
	proto.RegisterType((*CancelTaskRequest)(nil), "agent.CancelTaskRequest")
	proto.RegisterType((*CancelTaskResponse)(nil), "agent.CancelTaskResponse")
}

func init() {
	proto.RegisterFile("internal/server/proto/agent.proto", fileDescriptor_6286a68aa6b2b2f2)
}

var fileDescriptor_6286a68aa6b2b2f2 = []byte{
	// 413 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0xc1, 0x6e, 0xd4, 0x30,
	0x10, 0x55, 0x36, 0x74, 0x57, 0x9d, 0x16, 0x09, 0x4c, 0x29, 0x6e, 0xe0, 0xb0, 0xe4, 0x80, 0xf6,
	0x80, 0x12, 0xa9, 0x5c, 0x50, 0x85, 0x10, 0xa5, 0x54, 0x08, 0x8e, 0x29, 0x27, 0x2e, 0x8b, 0x37,
	0x1d, 0x2d, 0xd1, 0x26, 0x76, 0xb0, 0xc7, 0x61, 0xf7, 0xef, 0xf8, 0x17, 0x7e, 0x04, 0x25, 0x71,
	0x44, 0x96, 0xcd, 0x5e, 0x50, 0x4f, 0xf6, 0xf3, 0xbc, 0x79, 0x33, 0x7e, 0x63, 0xc3, 0xf3, 0x4c,
	0x12, 0x6a, 0x29, 0xf2, 0xd8, 0xa0, 0xae, 0x50, 0xc7, 0xa5, 0x56, 0xa4, 0x62, 0xb1, 0x44, 0x49,
	0x51, 0xb3, 0x67, 0x07, 0x0d, 0x08, 0x7f, 0x79, 0xc0, 0xae, 0xd7, 0x98, 0x5a, 0xc2, 0x2f, 0xc2,
	0xac, 0x12, 0xfc, 0x61, 0xd1, 0x10, 0x3b, 0x85, 0x71, 0xa9, 0x55, 0x51, 0x12, 0xf7, 0xa6, 0xde,
	0xec, 0x30, 0x71, 0x88, 0xbd, 0x83, 0x49, 0xaa, 0x24, 0xe1, 0x9a, 0xf8, 0x68, 0xea, 0xcf, 0x8e,
	0xce, 0x5f, 0x44, 0xad, 0xe8, 0xae, 0x46, 0x74, 0xd5, 0x12, 0xaf, 0x25, 0xe9, 0x4d, 0xd2, 0xa5,
	0xb1, 0x13, 0x38, 0x20, 0xa5, 0x72, 0xc3, 0xfd, 0xa9, 0x3f, 0x3b, 0x4c, 0x5a, 0x10, 0x5c, 0xc0,
	0x71, 0x9f, 0xce, 0x1e, 0x80, 0xbf, 0xc2, 0x8d, 0x2b, 0x5e, 0x6f, 0xeb, 0xbc, 0x4a, 0xe4, 0x16,
	0xf9, 0xa8, 0x39, 0x6b, 0xc1, 0xc5, 0xe8, 0xb5, 0x17, 0x7e, 0x83, 0x47, 0x5b, 0xd5, 0x4d, 0xa9,
	0xa4, 0x41, 0xf6, 0x04, 0x26, 0x24, 0xcc, 0x6a, 0x9e, 0xdd, 0x76, 0x77, 0xa8, 0xe1, 0xa7, 0xdb,
	0xfa, 0x6e, 0x86, 0x04, 0x59, 0xe3, 0xa4, 0x1c, 0x62, 0x1c, 0x26, 0x05, 0x1a, 0x23, 0x96, 0xc8,
	0xfd, 0x26, 0xd0, 0xc1, 0x30, 0x86, 0x93, 0x8f, 0x48, 0xb5, 0xfa, 0x4d, 0x43, 0xed, 0x5c, 0xda,
	0x57, 0x22, 0x5c, 0xc3, 0xe3, 0x7f, 0x12, 0xfe, 0xb7, 0xa9, 0x53, 0x18, 0x6b, 0x34, 0x36, 0x27,
	0xd7, 0x93, 0x43, 0xf5, 0x39, 0x56, 0x28, 0xc9, 0xf0, 0x7b, 0x8d, 0x8f, 0x0e, 0x85, 0x2f, 0xe1,
	0xe1, 0x95, 0x90, 0x29, 0xe6, 0xfd, 0x69, 0xee, 0xed, 0x73, 0x0e, 0xac, 0xcf, 0xbe, 0x73, 0xe7,
	0xce, 0x7f, 0x7b, 0x70, 0x7c, 0x59, 0x3f, 0x90, 0x1b, 0xd4, 0x55, 0x96, 0x22, 0xfb, 0x00, 0x47,
	0xbd, 0x61, 0xb1, 0xb3, 0xbd, 0xcf, 0x27, 0x08, 0x86, 0x42, 0xae, 0xc3, 0xcf, 0x70, 0x7f, 0xcb,
	0x5f, 0xf6, 0xd4, 0x91, 0x87, 0xc6, 0x14, 0x3c, 0x1b, 0x0e, 0x3a, 0xad, 0x4b, 0x80, 0xbf, 0x1e,
	0x30, 0xee, 0xb8, 0x3b, 0x26, 0x06, 0x67, 0x03, 0x91, 0x56, 0xe2, 0xfd, 0xdb, 0xaf, 0x6f, 0x96,
	0x19, 0x7d, 0xb7, 0x8b, 0x28, 0x55, 0x45, 0x6c, 0x4a, 0x4c, 0x49, 0xdb, 0xe2, 0x27, 0x2e, 0x52,
	0xf7, 0xe7, 0xe6, 0xda, 0x4a, 0xca, 0x0a, 0x8c, 0x07, 0xff, 0xe5, 0x62, 0xdc, 0x2c, 0xaf, 0xfe,
	0x04, 0x00, 0x00, 0xff, 0xff, 0x38, 0x71, 0x16, 0x87, 0xb7, 0x03, 0x00, 0x00,
}
