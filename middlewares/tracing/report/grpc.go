package report

import (
	"github.com/go-kit/kit/tracing/zipkin"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	opzipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
)

/// GRPC
func NewGrpcServerTracer(r reporter.Reporter) kitgrpc.ServerOption {
	zkTracer, err := opzipkin.NewTracer(r)
	if nil != err {
		panic(err)
	}
	return zipkin.GRPCServerTrace(zkTracer)
}

func NewGrpcClientTracer(r reporter.Reporter) kitgrpc.ClientOption {
	zkTracer, err := opzipkin.NewTracer(r)
	if nil != err {
		panic(err)
	}

	return zipkin.GRPCClientTrace(zkTracer) // zipkin.Name(cfg.ServiceName),

}
