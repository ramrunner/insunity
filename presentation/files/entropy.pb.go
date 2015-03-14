// Code generated by protoc-gen-go.
// source: entropy.proto
// DO NOT EDIT!

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	entropy.proto

It has these top-level messages:
	EntropyRequest
	EntropyResponse
*/
package service

import proto "github.com/chai2010/protorpc/proto"
import math "math"

import "io"
import "log"
import "net"
import "net/rpc"
import "time"
import protorpc "github.com/chai2010/protorpc"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type EntropyRequest struct {
	Data             *string `protobuf:"bytes,1,req,name=data" json:"data,omitempty"`
	Id               *string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	Version          *uint32 `protobuf:"varint,3,opt,name=version" json:"version,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *EntropyRequest) Reset()         { *m = EntropyRequest{} }
func (m *EntropyRequest) String() string { return proto.CompactTextString(m) }
func (*EntropyRequest) ProtoMessage()    {}

func (m *EntropyRequest) GetData() string {
	if m != nil && m.Data != nil {
		return *m.Data
	}
	return ""
}

func (m *EntropyRequest) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *EntropyRequest) GetVersion() uint32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

type EntropyResponse struct {
	Entropy          *float32 `protobuf:"fixed32,1,req,name=entropy" json:"entropy,omitempty"`
	Id               *string  `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	Version          *uint32  `protobuf:"varint,3,opt,name=version" json:"version,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *EntropyResponse) Reset()         { *m = EntropyResponse{} }
func (m *EntropyResponse) String() string { return proto.CompactTextString(m) }
func (*EntropyResponse) ProtoMessage()    {}

func (m *EntropyResponse) GetEntropy() float32 {
	if m != nil && m.Entropy != nil {
		return *m.Entropy
	}
	return 0
}

func (m *EntropyResponse) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *EntropyResponse) GetVersion() uint32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func init() {
}

type EntropyService interface {
	Entropy(in *EntropyRequest, out *EntropyResponse) error
}

// AcceptEntropyServiceClient accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func AcceptEntropyServiceClient(lis net.Listener, x EntropyService) {
	srv := rpc.NewServer()
	if err := srv.RegisterName("EntropyService", x); err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeCodec(protorpc.NewServerCodec(conn))
	}
}

// RegisterEntropyService publish the given EntropyService implementation on the server.
func RegisterEntropyService(srv *rpc.Server, x EntropyService) error {
	if err := srv.RegisterName("EntropyService", x); err != nil {
		return err
	}
	return nil
}

// NewEntropyServiceServer returns a new EntropyService Server.
func NewEntropyServiceServer(x EntropyService) *rpc.Server {
	srv := rpc.NewServer()
	if err := srv.RegisterName("EntropyService", x); err != nil {
		log.Fatal(err)
	}
	return srv
}

// ListenAndServeEntropyService listen announces on the local network address laddr
// and serves the given EntropyService implementation.
func ListenAndServeEntropyService(network, addr string, x EntropyService) error {
	lis, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("EntropyService", x); err != nil {
		return err
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeCodec(protorpc.NewServerCodec(conn))
	}
}

type EntropyServiceClient struct {
	*rpc.Client
}

// NewEntropyServiceClient returns a EntropyService rpc.Client and stub to handle
// requests to the set of EntropyService at the other end of the connection.
func NewEntropyServiceClient(conn io.ReadWriteCloser) (*EntropyServiceClient, *rpc.Client) {
	c := rpc.NewClientWithCodec(protorpc.NewClientCodec(conn))
	return &EntropyServiceClient{c}, c
}

func (c *EntropyServiceClient) Entropy(in *EntropyRequest, out *EntropyResponse) error {
	return c.Call("EntropyService.Entropy", in, out)
}

// DialEntropyService connects to an EntropyService at the specified network address.
func DialEntropyService(network, addr string) (*EntropyServiceClient, *rpc.Client, error) {
	c, err := protorpc.Dial(network, addr)
	if err != nil {
		return nil, nil, err
	}
	return &EntropyServiceClient{c}, c, nil
}

// DialEntropyServiceTimeout connects to an EntropyService at the specified network address.
func DialEntropyServiceTimeout(network, addr string,
	timeout time.Duration) (*EntropyServiceClient, *rpc.Client, error) {
	c, err := protorpc.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, nil, err
	}
	return &EntropyServiceClient{c}, c, nil
}