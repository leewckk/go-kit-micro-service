package gin

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-kit/log"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"

	kithttp "github.com/go-kit/kit/transport/http"
	ginhttp "github.com/leewckk/go-kit-micro-service/middlewares/transport/http"
)

// HTTPServerTrace enables native Zipkin tracing of a Go kit HTTP transport
// Server.
//
// Go kit creates HTTP transport servers per HTTP endpoint. This middleware can
// be set-up individually by adding the method name for each of the Go kit
// method servers using the Name() TracerOption.
// If wanting to use the HTTP method (Get, Post, Put, etc.) as Span name you can
// create a global server tracer omitting the Name() TracerOption, which you can
// then feed to each Go kit method server.
//
// If instrumenting a service to external (not on your platform) clients, you
// will probably want to disallow propagation of a client SpanContext using
// the AllowPropagation TracerOption and setting it to false.
func HTTPServerTrace(tracer *zipkin.Tracer, options ...TracerOption) ginhttp.ServerOption {
	config := tracerOptions{
		tags:      make(map[string]string),
		name:      "",
		logger:    log.NewNopLogger(),
		propagate: true,
	}

	for _, option := range options {
		option(&config)
	}

	serverBefore := ginhttp.ServerBefore(
		func(ctx context.Context, req *http.Request) context.Context {
			var (
				spanContext model.SpanContext
				name        string
			)

			if config.name != "" {
				name = config.name
			} else {
				name = req.Method
			}

			if config.propagate {
				spanContext = tracer.Extract(b3.ExtractHTTP(req))

				if spanContext.Sampled == nil && config.requestSampler != nil {
					sample := config.requestSampler(req)
					spanContext.Sampled = &sample
				}

				if spanContext.Err != nil {
					config.logger.Log("err", spanContext.Err)
				}
			}

			tags := map[string]string{
				string(zipkin.TagHTTPMethod): req.Method,
				string(zipkin.TagHTTPPath):   req.URL.Path,
			}

			span := tracer.StartSpan(
				name,
				zipkin.Kind(model.Server),
				zipkin.Tags(config.tags),
				zipkin.Tags(tags),
				zipkin.Parent(spanContext),
				zipkin.FlushOnFinish(false),
			)

			return zipkin.NewContext(ctx, span)
		},
	)

	serverAfter := ginhttp.ServerAfter(
		func(ctx context.Context, _ http.ResponseWriter) context.Context {
			if span := zipkin.SpanFromContext(ctx); span != nil {
				span.Finish()
			}

			return ctx
		},
	)

	serverFinalizer := ginhttp.ServerFinalizer(
		func(ctx context.Context, code int, r *http.Request) {
			if span := zipkin.SpanFromContext(ctx); span != nil {
				zipkin.TagHTTPStatusCode.Set(span, strconv.Itoa(code))
				if code > 399 {
					// set http status as error tag (if already set, this is a noop)
					zipkin.TagError.Set(span, http.StatusText(code))
				}
				if rs, ok := ctx.Value(kithttp.ContextKeyResponseSize).(int64); ok {
					zipkin.TagHTTPResponseSize.Set(span, strconv.FormatInt(rs, 10))
				}

				// calling span.Finish() a second time is a noop, if we didn't get to
				// ServerAfter we can at least time the early bail out by calling it
				// here.
				span.Finish()
				// send span to the Reporter
				span.Flush()
			}
		},
	)

	return func(s *ginhttp.Server) {
		serverBefore(s)
		serverAfter(s)
		serverFinalizer(s)
	}
}
