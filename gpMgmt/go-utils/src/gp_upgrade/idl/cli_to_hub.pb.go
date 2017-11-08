// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cli_to_hub.proto

/*
Package idl is a generated protocol buffer package.

It is generated from these files:
	cli_to_hub.proto
	command.proto

It has these top-level messages:
	StatusUpgradeRequest
	StatusUpgradeReply
	UpgradeStepStatus
	CheckUpgradeStatusRequest
	CheckUpgradeStatusReply
	FileSysUsage
	CheckDiskUsageRequest
	CheckDiskUsageReply
*/
package idl

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

type UpgradeSteps int32

const (
	UpgradeSteps_UNKNOWN_STEP UpgradeSteps = 0
	UpgradeSteps_CHECK_CONFIG UpgradeSteps = 1
	UpgradeSteps_SEGINSTALL   UpgradeSteps = 2
)

var UpgradeSteps_name = map[int32]string{
	0: "UNKNOWN_STEP",
	1: "CHECK_CONFIG",
	2: "SEGINSTALL",
}
var UpgradeSteps_value = map[string]int32{
	"UNKNOWN_STEP": 0,
	"CHECK_CONFIG": 1,
	"SEGINSTALL":   2,
}

func (x UpgradeSteps) String() string {
	return proto.EnumName(UpgradeSteps_name, int32(x))
}
func (UpgradeSteps) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type StepStatus int32

const (
	StepStatus_UNKNOWN_STATUS StepStatus = 0
	StepStatus_PENDING        StepStatus = 1
	StepStatus_RUNNING        StepStatus = 2
	StepStatus_COMPLETE       StepStatus = 3
	StepStatus_FAILED         StepStatus = 4
)

var StepStatus_name = map[int32]string{
	0: "UNKNOWN_STATUS",
	1: "PENDING",
	2: "RUNNING",
	3: "COMPLETE",
	4: "FAILED",
}
var StepStatus_value = map[string]int32{
	"UNKNOWN_STATUS": 0,
	"PENDING":        1,
	"RUNNING":        2,
	"COMPLETE":       3,
	"FAILED":         4,
}

func (x StepStatus) String() string {
	return proto.EnumName(StepStatus_name, int32(x))
}
func (StepStatus) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type StatusUpgradeRequest struct {
}

func (m *StatusUpgradeRequest) Reset()                    { *m = StatusUpgradeRequest{} }
func (m *StatusUpgradeRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusUpgradeRequest) ProtoMessage()               {}
func (*StatusUpgradeRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type StatusUpgradeReply struct {
	ListOfUpgradeStepStatuses []*UpgradeStepStatus `protobuf:"bytes,1,rep,name=listOfUpgradeStepStatuses" json:"listOfUpgradeStepStatuses,omitempty"`
}

func (m *StatusUpgradeReply) Reset()                    { *m = StatusUpgradeReply{} }
func (m *StatusUpgradeReply) String() string            { return proto.CompactTextString(m) }
func (*StatusUpgradeReply) ProtoMessage()               {}
func (*StatusUpgradeReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *StatusUpgradeReply) GetListOfUpgradeStepStatuses() []*UpgradeStepStatus {
	if m != nil {
		return m.ListOfUpgradeStepStatuses
	}
	return nil
}

type UpgradeStepStatus struct {
	Step   UpgradeSteps `protobuf:"varint,1,opt,name=step,enum=idl.UpgradeSteps" json:"step,omitempty"`
	Status StepStatus   `protobuf:"varint,2,opt,name=status,enum=idl.StepStatus" json:"status,omitempty"`
}

func (m *UpgradeStepStatus) Reset()                    { *m = UpgradeStepStatus{} }
func (m *UpgradeStepStatus) String() string            { return proto.CompactTextString(m) }
func (*UpgradeStepStatus) ProtoMessage()               {}
func (*UpgradeStepStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *UpgradeStepStatus) GetStep() UpgradeSteps {
	if m != nil {
		return m.Step
	}
	return UpgradeSteps_UNKNOWN_STEP
}

func (m *UpgradeStepStatus) GetStatus() StepStatus {
	if m != nil {
		return m.Status
	}
	return StepStatus_UNKNOWN_STATUS
}

func init() {
	proto.RegisterType((*StatusUpgradeRequest)(nil), "idl.StatusUpgradeRequest")
	proto.RegisterType((*StatusUpgradeReply)(nil), "idl.StatusUpgradeReply")
	proto.RegisterType((*UpgradeStepStatus)(nil), "idl.UpgradeStepStatus")
	proto.RegisterEnum("idl.UpgradeSteps", UpgradeSteps_name, UpgradeSteps_value)
	proto.RegisterEnum("idl.StepStatus", StepStatus_name, StepStatus_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for CliToHub service

type CliToHubClient interface {
	StatusUpgrade(ctx context.Context, in *StatusUpgradeRequest, opts ...grpc.CallOption) (*StatusUpgradeReply, error)
}

type cliToHubClient struct {
	cc *grpc.ClientConn
}

func NewCliToHubClient(cc *grpc.ClientConn) CliToHubClient {
	return &cliToHubClient{cc}
}

func (c *cliToHubClient) StatusUpgrade(ctx context.Context, in *StatusUpgradeRequest, opts ...grpc.CallOption) (*StatusUpgradeReply, error) {
	out := new(StatusUpgradeReply)
	err := grpc.Invoke(ctx, "/idl.CliToHub/StatusUpgrade", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for CliToHub service

type CliToHubServer interface {
	StatusUpgrade(context.Context, *StatusUpgradeRequest) (*StatusUpgradeReply, error)
}

func RegisterCliToHubServer(s *grpc.Server, srv CliToHubServer) {
	s.RegisterService(&_CliToHub_serviceDesc, srv)
}

func _CliToHub_StatusUpgrade_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusUpgradeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CliToHubServer).StatusUpgrade(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/idl.CliToHub/StatusUpgrade",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CliToHubServer).StatusUpgrade(ctx, req.(*StatusUpgradeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _CliToHub_serviceDesc = grpc.ServiceDesc{
	ServiceName: "idl.CliToHub",
	HandlerType: (*CliToHubServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StatusUpgrade",
			Handler:    _CliToHub_StatusUpgrade_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cli_to_hub.proto",
}

func init() { proto.RegisterFile("cli_to_hub.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 312 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xdf, 0x4f, 0xfa, 0x30,
	0x14, 0xc5, 0x19, 0x10, 0x20, 0x17, 0xbe, 0x7c, 0xcb, 0x8d, 0x41, 0xf0, 0x89, 0x2c, 0x31, 0x12,
	0x1e, 0x78, 0xc0, 0xbf, 0x00, 0x47, 0x81, 0x85, 0xd9, 0xe1, 0xd6, 0xc5, 0xc7, 0x85, 0x1f, 0x55,
	0x67, 0x9a, 0xac, 0xd2, 0xee, 0x81, 0xff, 0xde, 0x30, 0x49, 0x44, 0xd0, 0xc7, 0x9e, 0x73, 0x3e,
	0xb7, 0xed, 0x3d, 0x40, 0x36, 0x32, 0x89, 0x4d, 0x1a, 0xbf, 0x65, 0xeb, 0xa1, 0xda, 0xa5, 0x26,
	0xc5, 0x52, 0xb2, 0x95, 0x76, 0x1b, 0xae, 0x42, 0xb3, 0x32, 0x99, 0x8e, 0xd4, 0xeb, 0x6e, 0xb5,
	0x15, 0x81, 0xf8, 0xc8, 0x84, 0x36, 0xf6, 0x3b, 0xe0, 0x99, 0xae, 0xe4, 0x1e, 0x39, 0x74, 0x65,
	0xa2, 0x8d, 0xff, 0x72, 0x54, 0x43, 0x23, 0xd4, 0x57, 0x4c, 0xe8, 0x8e, 0xd5, 0x2b, 0xf5, 0xeb,
	0xa3, 0xf6, 0x30, 0xd9, 0xca, 0xe1, 0x85, 0x1f, 0xfc, 0x0d, 0xda, 0x1b, 0x68, 0x5d, 0xc8, 0x78,
	0x0b, 0x65, 0x6d, 0x84, 0xea, 0x58, 0x3d, 0xab, 0xdf, 0x1c, 0xb5, 0xce, 0xa7, 0xea, 0x20, 0xb7,
	0xf1, 0x0e, 0x2a, 0x3a, 0x07, 0x3a, 0xc5, 0x3c, 0xf8, 0x3f, 0x0f, 0x9e, 0xdc, 0x7b, 0xb4, 0x07,
	0x0f, 0xd0, 0x38, 0xc5, 0x91, 0x40, 0x23, 0x62, 0x0b, 0xe6, 0x3f, 0xb3, 0x38, 0xe4, 0x74, 0x49,
	0x0a, 0x07, 0xc5, 0x99, 0x53, 0x67, 0x11, 0x3b, 0x3e, 0x9b, 0xba, 0x33, 0x62, 0x61, 0x13, 0x20,
	0xa4, 0x33, 0x97, 0x85, 0x7c, 0xec, 0x79, 0xa4, 0x38, 0xe0, 0x00, 0x27, 0x2f, 0x44, 0x68, 0x7e,
	0x4f, 0x18, 0xf3, 0x28, 0x24, 0x05, 0xac, 0x43, 0x75, 0x49, 0xd9, 0xc4, 0x65, 0x07, 0xbc, 0x0e,
	0xd5, 0x20, 0x62, 0xec, 0x70, 0x28, 0x62, 0x03, 0x6a, 0x8e, 0xff, 0xb8, 0xf4, 0x28, 0xa7, 0xa4,
	0x84, 0x00, 0x95, 0xe9, 0xd8, 0xf5, 0xe8, 0x84, 0x94, 0x47, 0x4f, 0x50, 0x73, 0x64, 0xc2, 0xd3,
	0x79, 0xb6, 0x46, 0x0a, 0xff, 0x7e, 0xac, 0x1d, 0xbb, 0xc7, 0xff, 0x5c, 0x56, 0x74, 0x73, 0xfd,
	0x9b, 0xa5, 0xe4, 0xde, 0x2e, 0xac, 0x2b, 0x79, 0xc3, 0xf7, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x32, 0x35, 0x23, 0xe8, 0xf5, 0x01, 0x00, 0x00,
}
