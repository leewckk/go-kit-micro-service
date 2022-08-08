package gin

import (
	"net/http"

	"github.com/go-kit/log"
)

// TracerOption allows for functional options to our Zipkin tracing middleware.
type TracerOption func(o *tracerOptions)

// Name sets the name for an instrumented transport endpoint. If name is omitted
// at tracing middleware creation, the method of the transport or transport rpc
// name is used.
func Name(name string) TracerOption {
	return func(o *tracerOptions) {
		o.name = name
	}
}

// Tags adds default tags to our Zipkin transport spans.
func Tags(tags map[string]string) TracerOption {
	return func(o *tracerOptions) {
		for k, v := range tags {
			o.tags[k] = v
		}
	}
}

// Logger adds a Go kit logger to our Zipkin Middleware to log SpanContext
// extract / inject errors if they occur. Default is Noop.
func Logger(logger log.Logger) TracerOption {
	return func(o *tracerOptions) {
		if logger != nil {
			o.logger = logger
		}
	}
}

type tracerOptions struct {
	tags           map[string]string
	name           string
	logger         log.Logger
	propagate      bool
	requestSampler func(r *http.Request) bool
}
