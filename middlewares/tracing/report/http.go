package report

import (
	"github.com/go-kit/kit/tracing/zipkin"
	kithttp "github.com/go-kit/kit/transport/http"
	gintracing "github.com/leewckk/go-kit-micro-service/middlewares/tracing/gin"
	ginhttp "github.com/leewckk/go-kit-micro-service/middlewares/transport/http/gin"
	opzipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
)

/// HTTP
func NewHttpServerTracer(r reporter.Reporter) kithttp.ServerOption {
	zkTracer, err := opzipkin.NewTracer(r)
	if nil != err {
		panic(err)
	}
	return zipkin.HTTPServerTrace(zkTracer)
}

func NewHttpClientTracer(r reporter.Reporter) kithttp.ClientOption {
	zkTracer, err := opzipkin.NewTracer(r)
	if nil != err {
		panic(err)
	}
	return zipkin.HTTPClientTrace(zkTracer)
}

/// Gin HTTP
func NewGinServerTracer(r reporter.Reporter) ginhttp.ServerOption {
	zkTracer, err := opzipkin.NewTracer(r)
	if nil != err {
		panic(err)
	}
	return gintracing.HTTPServerTrace(zkTracer)
}
