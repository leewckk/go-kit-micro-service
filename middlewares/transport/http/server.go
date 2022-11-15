package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

type Server struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	enc          EncodeResponseFunc
	before       []kithttp.RequestFunc
	after        []kithttp.ServerResponseFunc
	errorEncoder kithttp.ErrorEncoder
	finalizer    []kithttp.ServerFinalizerFunc
	errorHandler transport.ErrorHandler
}

func NewServer(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	options ...ServerOption,
) *Server {

	s := &Server{
		e:            e,
		dec:          dec,
		enc:          enc,
		errorEncoder: kithttp.DefaultErrorEncoder,
		errorHandler: transport.NewLogErrorHandler(log.NewNopLogger()),
	}

	for _, option := range options {
		option(s)
	}
	return s
}

type ServerOption func(*Server)

// ServerBefore functions are executed on the HTTP request object before the
// request is decoded.
func ServerBefore(before ...kithttp.RequestFunc) ServerOption {
	return func(s *Server) { s.before = append(s.before, before...) }
}

// ServerAfter functions are executed on the HTTP response writer after the
// endpoint is invoked, but before anything is written to the client.
func ServerAfter(after ...kithttp.ServerResponseFunc) ServerOption {
	return func(s *Server) { s.after = append(s.after, after...) }
}

// ServerErrorEncoder is used to encode errors to the http.ResponseWriter
// whenever they're encountered in the processing of a request. Clients can
// use this to provide custom error formatting and response codes. By default,
// errors will be written with the DefaultErrorEncoder.
func ServerErrorEncoder(ee kithttp.ErrorEncoder) ServerOption {
	return func(s *Server) { s.errorEncoder = ee }
}

// ServerErrorLogger is used to log non-terminal errors. By default, no errors
// are logged. This is intended as a diagnostic measure. Finer-grained control
// of error handling, including logging in more detail, should be performed in a
// custom ServerErrorEncoder or ServerFinalizer, both of which have access to
// the context.
// Deprecated: Use ServerErrorHandler instead.
func ServerErrorLogger(logger log.Logger) ServerOption {
	return func(s *Server) { s.errorHandler = transport.NewLogErrorHandler(logger) }
}

// ServerErrorHandler is used to handle non-terminal errors. By default, non-terminal errors
// are ignored. This is intended as a diagnostic measure. Finer-grained control
// of error handling, including logging in more detail, should be performed in a
// custom ServerErrorEncoder or ServerFinalizer, both of which have access to
// the context.
func ServerErrorHandler(errorHandler transport.ErrorHandler) ServerOption {
	return func(s *Server) { s.errorHandler = errorHandler }
}

// ServerFinalizer is executed at the end of every HTTP request.
// By default, no finalizer is registered.
func ServerFinalizer(f ...kithttp.ServerFinalizerFunc) ServerOption {
	return func(s *Server) { s.finalizer = append(s.finalizer, f...) }
}

func (s Server) ServerHTTP(gctx *gin.Context) {
	w := http.ResponseWriter(gctx.Writer)
	r := gctx.Request
	ctx := gctx.Request.Context()
	// ctx := gctx

	if len(s.finalizer) > 0 {
		iw := &interceptingWriter{w, http.StatusOK, 0}
		defer func() {
			ctx = context.WithValue(ctx, kithttp.ContextKeyResponseHeaders, iw.Header())
			ctx = context.WithValue(ctx, kithttp.ContextKeyResponseSize, iw.written)
			for _, f := range s.finalizer {
				f(ctx, iw.code, r)
			}
		}()
		w = iw
	}

	for _, f := range s.before {
		ctx = f(ctx, r)
	}

	request, err := s.dec(gctx)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		s.errorEncoder(ctx, err, w)
		return
	}

	/// 将header信息保存到context中
	context.WithValue(ctx, "headers", gctx.Request.Header)
	response, err := s.e(ctx, request)
	// response, err := s.e(gctx, request)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		s.errorEncoder(ctx, err, w)
		return
	}

	for _, f := range s.after {
		ctx = f(ctx, w)
	}

	if err := s.enc(gctx, response); err != nil {
		s.errorHandler.Handle(ctx, err)
		s.errorEncoder(ctx, err, w)
		return
	}
}

type interceptingWriter struct {
	http.ResponseWriter
	code    int
	written int64
}

// WriteHeader may not be explicitly called, so care must be taken to
// initialize w.code to its default value of http.StatusOK.
func (w *interceptingWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *interceptingWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	w.written += int64(n)
	return n, err
}
