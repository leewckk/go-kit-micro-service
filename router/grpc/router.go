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
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	// "google.golang.org/grpc"
	"github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Router struct {
	name          string
	grpcOpts      []grpc.ServerOption
	routerOpts    []RouterOption
	registerProcs []RegisterRouteProc
	engine        *grpc.Server
}

type RouterOption func(r *Router)
type RegisterRouteProc func(s *grpc.Server, opts ...grpc.ServerOption)

func NewRouter(name string, opts ...RouterOption) *Router {

	r := &Router{
		name:          name + ".grpc",
		grpcOpts:      make([]grpc.ServerOption, 0, 0),
		routerOpts:    make([]RouterOption, 0, 0),
		registerProcs: make([]RegisterRouteProc, 0, 0),
	}
	r.grpcOpts = append(r.grpcOpts, grpc_middleware.WithUnaryServerChain(RecoveryInterceptor))

	for _, opt := range opts {
		opt(r)
	}
	r.engine = grpc.NewServer(r.grpcOpts...)
	return r
}

func RouterProcs(procs ...RegisterRouteProc) RouterOption {
	return func(r *Router) {
		r.registerProcs = append(r.registerProcs, procs...)
	}
}

/// server options
func RouterServerOptions(serverOpts ...grpc.ServerOption) RouterOption {
	return func(r *Router) {
		r.grpcOpts = append(r.grpcOpts, serverOpts...)
	}
}

func (r *Router) HealthAPI() string {
	return r.name
}

func (r *Router) Run(lis net.Listener) error {

	for _, reg := range r.registerProcs {
		if nil != reg {
			reg(r.engine, r.grpcOpts...)
		}
	}

	/// health
	healthServer := health.NewServer()
	healthServer.SetServingStatus(r.name, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(r.engine, healthServer)
	reflection.Register(r.engine)
	return r.engine.Serve(lis)
}
