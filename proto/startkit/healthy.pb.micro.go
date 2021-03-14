// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/startkit/healthy.proto

package xtc_ogm_startkit

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/v2/client"
	server "github.com/micro/go-micro/v2/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Healthy service

type HealthyService interface {
	Echo(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error)
	PingPong(ctx context.Context, opts ...client.CallOption) (Healthy_PingPongService, error)
}

type healthyService struct {
	c    client.Client
	name string
}

func NewHealthyService(name string, c client.Client) HealthyService {
	return &healthyService{
		c:    c,
		name: name,
	}
}

func (c *healthyService) Echo(ctx context.Context, in *Request, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Healthy.Echo", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *healthyService) PingPong(ctx context.Context, opts ...client.CallOption) (Healthy_PingPongService, error) {
	req := c.c.NewRequest(c.name, "Healthy.PingPong", &Ping{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &healthyServicePingPong{stream}, nil
}

type Healthy_PingPongService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Ping) error
	Recv() (*Pong, error)
}

type healthyServicePingPong struct {
	stream client.Stream
}

func (x *healthyServicePingPong) Close() error {
	return x.stream.Close()
}

func (x *healthyServicePingPong) Context() context.Context {
	return x.stream.Context()
}

func (x *healthyServicePingPong) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *healthyServicePingPong) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *healthyServicePingPong) Send(m *Ping) error {
	return x.stream.Send(m)
}

func (x *healthyServicePingPong) Recv() (*Pong, error) {
	m := new(Pong)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Healthy service

type HealthyHandler interface {
	Echo(context.Context, *Request, *Response) error
	PingPong(context.Context, Healthy_PingPongStream) error
}

func RegisterHealthyHandler(s server.Server, hdlr HealthyHandler, opts ...server.HandlerOption) error {
	type healthy interface {
		Echo(ctx context.Context, in *Request, out *Response) error
		PingPong(ctx context.Context, stream server.Stream) error
	}
	type Healthy struct {
		healthy
	}
	h := &healthyHandler{hdlr}
	return s.Handle(s.NewHandler(&Healthy{h}, opts...))
}

type healthyHandler struct {
	HealthyHandler
}

func (h *healthyHandler) Echo(ctx context.Context, in *Request, out *Response) error {
	return h.HealthyHandler.Echo(ctx, in, out)
}

func (h *healthyHandler) PingPong(ctx context.Context, stream server.Stream) error {
	return h.HealthyHandler.PingPong(ctx, &healthyPingPongStream{stream})
}

type Healthy_PingPongStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
	Send(*Pong) error
	Recv() (*Ping, error)
}

type healthyPingPongStream struct {
	stream server.Stream
}

func (x *healthyPingPongStream) Close() error {
	return x.stream.Close()
}

func (x *healthyPingPongStream) Context() context.Context {
	return x.stream.Context()
}

func (x *healthyPingPongStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *healthyPingPongStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *healthyPingPongStream) Send(m *Pong) error {
	return x.stream.Send(m)
}

func (x *healthyPingPongStream) Recv() (*Ping, error) {
	m := new(Ping)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}