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

package client

import (
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/leewckk/go-kit-micro-service/log"
)

func NewFactory(
	method string,
	api string,
	enc kithttp.EncodeRequestFunc,
	dec kithttp.DecodeResponseFunc,
	opts ...kithttp.ClientOption,
) sd.Factory {

	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}
		target, err := url.Parse(instance)
		if nil != err {
			panic(err)
		}
		target.Path = api

		// opts := make([]kithttp.ClientOption, 0, 0)
		// reporter := tracing.NewReporter()
		// if nil != reporter {
		// 	opts = append(opts, tracing.NewHttpClientTracer(reporter))
		// }

		return kithttp.NewClient(
			method,
			target,
			enc, dec,
			opts...,
		).Endpoint(), nil, nil
	}
}

func MakeClientEndpoint(
	client consul.Client,
	serverName string,
	methodName string,
	api string,
	enc kithttp.EncodeRequestFunc,
	dec kithttp.DecodeResponseFunc,
	opts ...kithttp.ClientOption) endpoint.Endpoint {

	logger := log.NewLogger()
	// client := discovery.NewClient(sdapi)
	serverName = serverName + ".http"

	passingOnly := true
	duration := 120 * time.Second
	instancer := consul.NewInstancer(client, logger, serverName, []string{}, passingOnly)

	factory := NewFactory(methodName, api, enc, dec, opts...)
	ep := sd.NewEndpointer(instancer,
		factory,
		logger,
	)
	balancer := lb.NewRoundRobin(ep)
	retry := lb.Retry(3, duration, balancer)
	return retry
}
