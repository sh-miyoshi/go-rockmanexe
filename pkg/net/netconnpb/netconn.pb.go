// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.12
// source: pkg/net/netconnpb/netconn.proto

package netconnpb

import (
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

type Request_MsgType int32

const (
	Request_AUTHENTICATE Request_MsgType = 0
	Request_SENDSIGNAL   Request_MsgType = 1
	Request_ACTION       Request_MsgType = 2
)

// Enum value maps for Request_MsgType.
var (
	Request_MsgType_name = map[int32]string{
		0: "AUTHENTICATE",
		1: "SENDSIGNAL",
		2: "ACTION",
	}
	Request_MsgType_value = map[string]int32{
		"AUTHENTICATE": 0,
		"SENDSIGNAL":   1,
		"ACTION":       2,
	}
)

func (x Request_MsgType) Enum() *Request_MsgType {
	p := new(Request_MsgType)
	*p = x
	return p
}

func (x Request_MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Request_MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_net_netconnpb_netconn_proto_enumTypes[0].Descriptor()
}

func (Request_MsgType) Type() protoreflect.EnumType {
	return &file_pkg_net_netconnpb_netconn_proto_enumTypes[0]
}

func (x Request_MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Request_MsgType.Descriptor instead.
func (Request_MsgType) EnumDescriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{0, 0}
}

type Request_SignalType int32

const (
	Request_CHIPSELECT   Request_SignalType = 0
	Request_GOCHIPSELECT Request_SignalType = 1
	Request_INITPARAMS   Request_SignalType = 2
	Request_CUTIN        Request_SignalType = 3
)

// Enum value maps for Request_SignalType.
var (
	Request_SignalType_name = map[int32]string{
		0: "CHIPSELECT",
		1: "GOCHIPSELECT",
		2: "INITPARAMS",
		3: "CUTIN",
	}
	Request_SignalType_value = map[string]int32{
		"CHIPSELECT":   0,
		"GOCHIPSELECT": 1,
		"INITPARAMS":   2,
		"CUTIN":        3,
	}
)

func (x Request_SignalType) Enum() *Request_SignalType {
	p := new(Request_SignalType)
	*p = x
	return p
}

func (x Request_SignalType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Request_SignalType) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_net_netconnpb_netconn_proto_enumTypes[1].Descriptor()
}

func (Request_SignalType) Type() protoreflect.EnumType {
	return &file_pkg_net_netconnpb_netconn_proto_enumTypes[1]
}

func (x Request_SignalType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Request_SignalType.Descriptor instead.
func (Request_SignalType) EnumDescriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{0, 1}
}

type Request_ActionType int32

const (
	Request_MOVE    Request_ActionType = 0
	Request_BUSTER  Request_ActionType = 1
	Request_CHIPUSE Request_ActionType = 2
)

// Enum value maps for Request_ActionType.
var (
	Request_ActionType_name = map[int32]string{
		0: "MOVE",
		1: "BUSTER",
		2: "CHIPUSE",
	}
	Request_ActionType_value = map[string]int32{
		"MOVE":    0,
		"BUSTER":  1,
		"CHIPUSE": 2,
	}
)

func (x Request_ActionType) Enum() *Request_ActionType {
	p := new(Request_ActionType)
	*p = x
	return p
}

func (x Request_ActionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Request_ActionType) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_net_netconnpb_netconn_proto_enumTypes[2].Descriptor()
}

func (Request_ActionType) Type() protoreflect.EnumType {
	return &file_pkg_net_netconnpb_netconn_proto_enumTypes[2]
}

func (x Request_ActionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Request_ActionType.Descriptor instead.
func (Request_ActionType) EnumDescriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{0, 2}
}

type Response_MsgType int32

const (
	Response_AUTHRESPONSE Response_MsgType = 0
	Response_UPDATESTATUS Response_MsgType = 1
	Response_DATA         Response_MsgType = 2
	Response_SYSTEM       Response_MsgType = 3
)

// Enum value maps for Response_MsgType.
var (
	Response_MsgType_name = map[int32]string{
		0: "AUTHRESPONSE",
		1: "UPDATESTATUS",
		2: "DATA",
		3: "SYSTEM",
	}
	Response_MsgType_value = map[string]int32{
		"AUTHRESPONSE": 0,
		"UPDATESTATUS": 1,
		"DATA":         2,
		"SYSTEM":       3,
	}
)

func (x Response_MsgType) Enum() *Response_MsgType {
	p := new(Response_MsgType)
	*p = x
	return p
}

func (x Response_MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Response_MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_net_netconnpb_netconn_proto_enumTypes[3].Descriptor()
}

func (Response_MsgType) Type() protoreflect.EnumType {
	return &file_pkg_net_netconnpb_netconn_proto_enumTypes[3]
}

func (x Response_MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Response_MsgType.Descriptor instead.
func (Response_MsgType) EnumDescriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{1, 0}
}

type Response_Status int32

const (
	Response_CONNECTWAIT    Response_Status = 0
	Response_CHIPSELECTWAIT Response_Status = 1
	Response_ACTING         Response_Status = 2
	Response_GAMEEND        Response_Status = 3
	Response_CUTIN          Response_Status = 4
)

// Enum value maps for Response_Status.
var (
	Response_Status_name = map[int32]string{
		0: "CONNECTWAIT",
		1: "CHIPSELECTWAIT",
		2: "ACTING",
		3: "GAMEEND",
		4: "CUTIN",
	}
	Response_Status_value = map[string]int32{
		"CONNECTWAIT":    0,
		"CHIPSELECTWAIT": 1,
		"ACTING":         2,
		"GAMEEND":        3,
		"CUTIN":          4,
	}
)

func (x Response_Status) Enum() *Response_Status {
	p := new(Response_Status)
	*p = x
	return p
}

func (x Response_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Response_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_net_netconnpb_netconn_proto_enumTypes[4].Descriptor()
}

func (Response_Status) Type() protoreflect.EnumType {
	return &file_pkg_net_netconnpb_netconn_proto_enumTypes[4]
}

func (x Response_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Response_Status.Descriptor instead.
func (Response_Status) EnumDescriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{1, 1}
}

type Response_System_SystemType int32

const (
	Response_System_CUTIN Response_System_SystemType = 0
)

// Enum value maps for Response_System_SystemType.
var (
	Response_System_SystemType_name = map[int32]string{
		0: "CUTIN",
	}
	Response_System_SystemType_value = map[string]int32{
		"CUTIN": 0,
	}
)

func (x Response_System_SystemType) Enum() *Response_System_SystemType {
	p := new(Response_System_SystemType)
	*p = x
	return p
}

func (x Response_System_SystemType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Response_System_SystemType) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_net_netconnpb_netconn_proto_enumTypes[5].Descriptor()
}

func (Response_System_SystemType) Type() protoreflect.EnumType {
	return &file_pkg_net_netconnpb_netconn_proto_enumTypes[5]
}

func (x Response_System_SystemType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Response_System_SystemType.Descriptor instead.
func (Response_System_SystemType) EnumDescriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{1, 1, 0}
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SessionID string          `protobuf:"bytes,1,opt,name=sessionID,proto3" json:"sessionID,omitempty"`
	ClientID  string          `protobuf:"bytes,2,opt,name=clientID,proto3" json:"clientID,omitempty"`
	Type      Request_MsgType `protobuf:"varint,3,opt,name=type,proto3,enum=netconn.Request_MsgType" json:"type,omitempty"`
	// Types that are assignable to Data:
	//
	//	*Request_Req
	//	*Request_Signal_
	//	*Request_Act
	Data isRequest_Data `protobuf_oneof:"data"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetSessionID() string {
	if x != nil {
		return x.SessionID
	}
	return ""
}

func (x *Request) GetClientID() string {
	if x != nil {
		return x.ClientID
	}
	return ""
}

func (x *Request) GetType() Request_MsgType {
	if x != nil {
		return x.Type
	}
	return Request_AUTHENTICATE
}

func (m *Request) GetData() isRequest_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *Request) GetReq() *Request_AuthRequest {
	if x, ok := x.GetData().(*Request_Req); ok {
		return x.Req
	}
	return nil
}

func (x *Request) GetSignal() *Request_Signal {
	if x, ok := x.GetData().(*Request_Signal_); ok {
		return x.Signal
	}
	return nil
}

func (x *Request) GetAct() *Request_Action {
	if x, ok := x.GetData().(*Request_Act); ok {
		return x.Act
	}
	return nil
}

type isRequest_Data interface {
	isRequest_Data()
}

type Request_Req struct {
	Req *Request_AuthRequest `protobuf:"bytes,4,opt,name=req,proto3,oneof"`
}

type Request_Signal_ struct {
	Signal *Request_Signal `protobuf:"bytes,5,opt,name=signal,proto3,oneof"`
}

type Request_Act struct {
	Act *Request_Action `protobuf:"bytes,6,opt,name=act,proto3,oneof"`
}

func (*Request_Req) isRequest_Data() {}

func (*Request_Signal_) isRequest_Data() {}

func (*Request_Act) isRequest_Data() {}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type Response_MsgType `protobuf:"varint,1,opt,name=type,proto3,enum=netconn.Response_MsgType" json:"type,omitempty"`
	// Types that are assignable to Data:
	//
	//	*Response_AuthRes
	//	*Response_Status_
	//	*Response_RawData
	//	*Response_System_
	Data isResponse_Data `protobuf_oneof:"data"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{1}
}

func (x *Response) GetType() Response_MsgType {
	if x != nil {
		return x.Type
	}
	return Response_AUTHRESPONSE
}

func (m *Response) GetData() isResponse_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *Response) GetAuthRes() *Response_AuthResponse {
	if x, ok := x.GetData().(*Response_AuthRes); ok {
		return x.AuthRes
	}
	return nil
}

func (x *Response) GetStatus() Response_Status {
	if x, ok := x.GetData().(*Response_Status_); ok {
		return x.Status
	}
	return Response_CONNECTWAIT
}

func (x *Response) GetRawData() []byte {
	if x, ok := x.GetData().(*Response_RawData); ok {
		return x.RawData
	}
	return nil
}

func (x *Response) GetSystem() *Response_System {
	if x, ok := x.GetData().(*Response_System_); ok {
		return x.System
	}
	return nil
}

type isResponse_Data interface {
	isResponse_Data()
}

type Response_AuthRes struct {
	AuthRes *Response_AuthResponse `protobuf:"bytes,2,opt,name=authRes,proto3,oneof"`
}

type Response_Status_ struct {
	Status Response_Status `protobuf:"varint,3,opt,name=status,proto3,enum=netconn.Response_Status,oneof"`
}

type Response_RawData struct {
	RawData []byte `protobuf:"bytes,4,opt,name=rawData,proto3,oneof"`
}

type Response_System_ struct {
	System *Response_System `protobuf:"bytes,5,opt,name=system,proto3,oneof"`
}

func (*Response_AuthRes) isResponse_Data() {}

func (*Response_Status_) isResponse_Data() {}

func (*Response_RawData) isResponse_Data() {}

func (*Response_System_) isResponse_Data() {}

type Request_AuthRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Key     string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *Request_AuthRequest) Reset() {
	*x = Request_AuthRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request_AuthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request_AuthRequest) ProtoMessage() {}

func (x *Request_AuthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request_AuthRequest.ProtoReflect.Descriptor instead.
func (*Request_AuthRequest) Descriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Request_AuthRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Request_AuthRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Request_AuthRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

type Request_Action struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    Request_ActionType `protobuf:"varint,1,opt,name=type,proto3,enum=netconn.Request_ActionType" json:"type,omitempty"`
	RawData []byte             `protobuf:"bytes,2,opt,name=rawData,proto3" json:"rawData,omitempty"`
}

func (x *Request_Action) Reset() {
	*x = Request_Action{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request_Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request_Action) ProtoMessage() {}

func (x *Request_Action) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request_Action.ProtoReflect.Descriptor instead.
func (*Request_Action) Descriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Request_Action) GetType() Request_ActionType {
	if x != nil {
		return x.Type
	}
	return Request_MOVE
}

func (x *Request_Action) GetRawData() []byte {
	if x != nil {
		return x.RawData
	}
	return nil
}

type Request_Signal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    Request_SignalType `protobuf:"varint,1,opt,name=type,proto3,enum=netconn.Request_SignalType" json:"type,omitempty"`
	RawData []byte             `protobuf:"bytes,2,opt,name=rawData,proto3" json:"rawData,omitempty"`
}

func (x *Request_Signal) Reset() {
	*x = Request_Signal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request_Signal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request_Signal) ProtoMessage() {}

func (x *Request_Signal) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request_Signal.ProtoReflect.Descriptor instead.
func (*Request_Signal) Descriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{0, 2}
}

func (x *Request_Signal) GetType() Request_SignalType {
	if x != nil {
		return x.Type
	}
	return Request_CHIPSELECT
}

func (x *Request_Signal) GetRawData() []byte {
	if x != nil {
		return x.RawData
	}
	return nil
}

type Response_AuthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success    bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ErrMsg     string   `protobuf:"bytes,2,opt,name=errMsg,proto3" json:"errMsg,omitempty"`
	SessionID  string   `protobuf:"bytes,3,opt,name=sessionID,proto3" json:"sessionID,omitempty"`
	AllUserIDs []string `protobuf:"bytes,4,rep,name=allUserIDs,proto3" json:"allUserIDs,omitempty"`
}

func (x *Response_AuthResponse) Reset() {
	*x = Response_AuthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response_AuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response_AuthResponse) ProtoMessage() {}

func (x *Response_AuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response_AuthResponse.ProtoReflect.Descriptor instead.
func (*Response_AuthResponse) Descriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Response_AuthResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *Response_AuthResponse) GetErrMsg() string {
	if x != nil {
		return x.ErrMsg
	}
	return ""
}

func (x *Response_AuthResponse) GetSessionID() string {
	if x != nil {
		return x.SessionID
	}
	return ""
}

func (x *Response_AuthResponse) GetAllUserIDs() []string {
	if x != nil {
		return x.AllUserIDs
	}
	return nil
}

type Response_System struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    Response_System_SystemType `protobuf:"varint,1,opt,name=type,proto3,enum=netconn.Response_System_SystemType" json:"type,omitempty"`
	RawData []byte                     `protobuf:"bytes,2,opt,name=rawData,proto3" json:"rawData,omitempty"`
}

func (x *Response_System) Reset() {
	*x = Response_System{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response_System) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response_System) ProtoMessage() {}

func (x *Response_System) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_net_netconnpb_netconn_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response_System.ProtoReflect.Descriptor instead.
func (*Response_System) Descriptor() ([]byte, []int) {
	return file_pkg_net_netconnpb_netconn_proto_rawDescGZIP(), []int{1, 1}
}

func (x *Response_System) GetType() Response_System_SystemType {
	if x != nil {
		return x.Type
	}
	return Response_System_CUTIN
}

func (x *Response_System) GetRawData() []byte {
	if x != nil {
		return x.RawData
	}
	return nil
}

var File_pkg_net_netconnpb_netconn_proto protoreflect.FileDescriptor

var file_pkg_net_netconnpb_netconn_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x6b, 0x67, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e,
	0x6e, 0x70, 0x62, 0x2f, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x07, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x22, 0xb5, 0x05, 0x0a, 0x07, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44,
	0x12, 0x2c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18,
	0x2e, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x30,
	0x0a, 0x03, 0x72, 0x65, 0x71, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6e, 0x65,
	0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x75,
	0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x03, 0x72, 0x65, 0x71,
	0x12, 0x31, 0x0a, 0x06, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x48, 0x00, 0x52, 0x06, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x6c, 0x12, 0x2b, 0x0a, 0x03, 0x61, 0x63, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x03, 0x61, 0x63, 0x74,
	0x1a, 0x49, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x53, 0x0a, 0x06, 0x41,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x61, 0x77, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x72, 0x61, 0x77, 0x44, 0x61, 0x74, 0x61,
	0x1a, 0x53, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x12, 0x2f, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x6e, 0x65, 0x74, 0x63, 0x6f,
	0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x61,
	0x6c, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72,
	0x61, 0x77, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x72, 0x61,
	0x77, 0x44, 0x61, 0x74, 0x61, 0x22, 0x37, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x10, 0x0a, 0x0c, 0x41, 0x55, 0x54, 0x48, 0x45, 0x4e, 0x54, 0x49, 0x43, 0x41, 0x54, 0x45,
	0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x45, 0x4e, 0x44, 0x53, 0x49, 0x47, 0x4e, 0x41, 0x4c,
	0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02, 0x22, 0x49,
	0x0a, 0x0a, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a,
	0x43, 0x48, 0x49, 0x50, 0x53, 0x45, 0x4c, 0x45, 0x43, 0x54, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c,
	0x47, 0x4f, 0x43, 0x48, 0x49, 0x50, 0x53, 0x45, 0x4c, 0x45, 0x43, 0x54, 0x10, 0x01, 0x12, 0x0e,
	0x0a, 0x0a, 0x49, 0x4e, 0x49, 0x54, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x53, 0x10, 0x02, 0x12, 0x09,
	0x0a, 0x05, 0x43, 0x55, 0x54, 0x49, 0x4e, 0x10, 0x03, 0x22, 0x2f, 0x0a, 0x0a, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4d, 0x4f, 0x56, 0x45, 0x10,
	0x00, 0x12, 0x0a, 0x0a, 0x06, 0x42, 0x55, 0x53, 0x54, 0x45, 0x52, 0x10, 0x01, 0x12, 0x0b, 0x0a,
	0x07, 0x43, 0x48, 0x49, 0x50, 0x55, 0x53, 0x45, 0x10, 0x02, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x22, 0x8f, 0x05, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2d, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e,
	0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x3a,
	0x0a, 0x07, 0x61, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48,
	0x00, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x18, 0x2e, 0x6e, 0x65, 0x74,
	0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x48, 0x00, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1a,
	0x0a, 0x07, 0x72, 0x61, 0x77, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x48,
	0x00, 0x52, 0x07, 0x72, 0x61, 0x77, 0x44, 0x61, 0x74, 0x61, 0x12, 0x32, 0x0a, 0x06, 0x73, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x6e, 0x65, 0x74,
	0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x48, 0x00, 0x52, 0x06, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x1a, 0x7e,
	0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x4d,
	0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x72, 0x72, 0x4d, 0x73, 0x67,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x1e,
	0x0a, 0x0a, 0x61, 0x6c, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0a, 0x61, 0x6c, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x73, 0x1a, 0x74,
	0x0a, 0x06, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x37, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e,
	0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x2e, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x61, 0x77, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x72, 0x61, 0x77, 0x44, 0x61, 0x74, 0x61, 0x22, 0x17, 0x0a, 0x0a, 0x53,
	0x79, 0x73, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x55, 0x54,
	0x49, 0x4e, 0x10, 0x00, 0x22, 0x43, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x10, 0x0a, 0x0c, 0x41, 0x55, 0x54, 0x48, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x10,
	0x00, 0x12, 0x10, 0x0a, 0x0c, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x53, 0x54, 0x41, 0x54, 0x55,
	0x53, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x41, 0x54, 0x41, 0x10, 0x02, 0x12, 0x0a, 0x0a,
	0x06, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x10, 0x03, 0x22, 0x51, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x57, 0x41,
	0x49, 0x54, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x43, 0x48, 0x49, 0x50, 0x53, 0x45, 0x4c, 0x45,
	0x43, 0x54, 0x57, 0x41, 0x49, 0x54, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54, 0x49,
	0x4e, 0x47, 0x10, 0x02, 0x12, 0x0b, 0x0a, 0x07, 0x47, 0x41, 0x4d, 0x45, 0x45, 0x4e, 0x44, 0x10,
	0x03, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x55, 0x54, 0x49, 0x4e, 0x10, 0x04, 0x42, 0x06, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x32, 0x41, 0x0a, 0x07, 0x4e, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x12,
	0x36, 0x0a, 0x09, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x2e, 0x6e,
	0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11,
	0x2e, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x37, 0x5a, 0x35, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x68, 0x2d, 0x6d, 0x69, 0x79, 0x6f, 0x73, 0x68, 0x69,
	0x2f, 0x67, 0x6f, 0x2d, 0x72, 0x6f, 0x63, 0x6b, 0x6d, 0x61, 0x6e, 0x65, 0x78, 0x65, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x6e, 0x65, 0x74, 0x2f, 0x6e, 0x65, 0x74, 0x63, 0x6f, 0x6e, 0x6e, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_net_netconnpb_netconn_proto_rawDescOnce sync.Once
	file_pkg_net_netconnpb_netconn_proto_rawDescData = file_pkg_net_netconnpb_netconn_proto_rawDesc
)

func file_pkg_net_netconnpb_netconn_proto_rawDescGZIP() []byte {
	file_pkg_net_netconnpb_netconn_proto_rawDescOnce.Do(func() {
		file_pkg_net_netconnpb_netconn_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_net_netconnpb_netconn_proto_rawDescData)
	})
	return file_pkg_net_netconnpb_netconn_proto_rawDescData
}

var file_pkg_net_netconnpb_netconn_proto_enumTypes = make([]protoimpl.EnumInfo, 6)
var file_pkg_net_netconnpb_netconn_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pkg_net_netconnpb_netconn_proto_goTypes = []interface{}{
	(Request_MsgType)(0),            // 0: netconn.Request.MsgType
	(Request_SignalType)(0),         // 1: netconn.Request.SignalType
	(Request_ActionType)(0),         // 2: netconn.Request.ActionType
	(Response_MsgType)(0),           // 3: netconn.Response.MsgType
	(Response_Status)(0),            // 4: netconn.Response.Status
	(Response_System_SystemType)(0), // 5: netconn.Response.System.SystemType
	(*Request)(nil),                 // 6: netconn.Request
	(*Response)(nil),                // 7: netconn.Response
	(*Request_AuthRequest)(nil),     // 8: netconn.Request.AuthRequest
	(*Request_Action)(nil),          // 9: netconn.Request.Action
	(*Request_Signal)(nil),          // 10: netconn.Request.Signal
	(*Response_AuthResponse)(nil),   // 11: netconn.Response.AuthResponse
	(*Response_System)(nil),         // 12: netconn.Response.System
}
var file_pkg_net_netconnpb_netconn_proto_depIdxs = []int32{
	0,  // 0: netconn.Request.type:type_name -> netconn.Request.MsgType
	8,  // 1: netconn.Request.req:type_name -> netconn.Request.AuthRequest
	10, // 2: netconn.Request.signal:type_name -> netconn.Request.Signal
	9,  // 3: netconn.Request.act:type_name -> netconn.Request.Action
	3,  // 4: netconn.Response.type:type_name -> netconn.Response.MsgType
	11, // 5: netconn.Response.authRes:type_name -> netconn.Response.AuthResponse
	4,  // 6: netconn.Response.status:type_name -> netconn.Response.Status
	12, // 7: netconn.Response.system:type_name -> netconn.Response.System
	2,  // 8: netconn.Request.Action.type:type_name -> netconn.Request.ActionType
	1,  // 9: netconn.Request.Signal.type:type_name -> netconn.Request.SignalType
	5,  // 10: netconn.Response.System.type:type_name -> netconn.Response.System.SystemType
	6,  // 11: netconn.NetConn.TransData:input_type -> netconn.Request
	7,  // 12: netconn.NetConn.TransData:output_type -> netconn.Response
	12, // [12:13] is the sub-list for method output_type
	11, // [11:12] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_pkg_net_netconnpb_netconn_proto_init() }
func file_pkg_net_netconnpb_netconn_proto_init() {
	if File_pkg_net_netconnpb_netconn_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_net_netconnpb_netconn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_pkg_net_netconnpb_netconn_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_pkg_net_netconnpb_netconn_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request_AuthRequest); i {
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
		file_pkg_net_netconnpb_netconn_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request_Action); i {
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
		file_pkg_net_netconnpb_netconn_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request_Signal); i {
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
		file_pkg_net_netconnpb_netconn_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response_AuthResponse); i {
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
		file_pkg_net_netconnpb_netconn_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response_System); i {
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
	file_pkg_net_netconnpb_netconn_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Request_Req)(nil),
		(*Request_Signal_)(nil),
		(*Request_Act)(nil),
	}
	file_pkg_net_netconnpb_netconn_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Response_AuthRes)(nil),
		(*Response_Status_)(nil),
		(*Response_RawData)(nil),
		(*Response_System_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_net_netconnpb_netconn_proto_rawDesc,
			NumEnums:      6,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_net_netconnpb_netconn_proto_goTypes,
		DependencyIndexes: file_pkg_net_netconnpb_netconn_proto_depIdxs,
		EnumInfos:         file_pkg_net_netconnpb_netconn_proto_enumTypes,
		MessageInfos:      file_pkg_net_netconnpb_netconn_proto_msgTypes,
	}.Build()
	File_pkg_net_netconnpb_netconn_proto = out.File
	file_pkg_net_netconnpb_netconn_proto_rawDesc = nil
	file_pkg_net_netconnpb_netconn_proto_goTypes = nil
	file_pkg_net_netconnpb_netconn_proto_depIdxs = nil
}
