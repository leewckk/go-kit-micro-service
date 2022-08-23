// The MIT License (MIT)

// Copyright (c) 2022 leewckk@126.com

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package grpc

import (
	"io"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	gokitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/leewckk/go-kit-micro-service/log"
)

func NewFactory(
	c *Connection,
	serviceName string,
	method string,
	enc gokitgrpc.EncodeRequestFunc,
	dec gokitgrpc.DecodeResponseFunc,
	grpcReply interface{},
	opts ...gokitgrpc.ClientOption,
) sd.Factory {

	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		// opts := make([]gokitgrpc.ClientOption, 0, 0)
		// reporter := report.NewZipkinReporter(tracingAPI)
		// if nil != reporter {
		// 	opts = append(opts, report.NewGrpcClientTracer(reporter))
		// }

		conn := c.GetConnection(instance)
		if nil == conn {
			panic("failed get grpc conn")
		}
		return gokitgrpc.NewClient(
			conn,
			serviceName,
			method,
			enc, dec,
			grpcReply,
			opts...,
		).Endpoint(), nil, nil
	}
}

func MakeClientEndpoint(
	client consul.Client,
	serverName string,
	serviceName string,
	methodName string,
	enc gokitgrpc.EncodeRequestFunc,
	dec gokitgrpc.DecodeResponseFunc,
	// reply interface{},
	replyMaker func() interface{},
	opts ...gokitgrpc.ClientOption) endpoint.Endpoint {

	// client := discovery.NewClient(sdapi)
	serverName += ".grpc"
	passingOnly := true
	duration := 120 * time.Second

	reply := replyMaker()
	instancer := consul.NewInstancer(client, log.NewLogger(), serverName, []string{}, passingOnly)
	factory := NewFactory(NewConnection(), serviceName, methodName, enc, dec, reply, opts...)
	ep := sd.NewEndpointer(instancer, factory, log.NewLogger())
	bl := lb.NewRoundRobin(ep)
	return lb.Retry(3, duration, bl)
}
