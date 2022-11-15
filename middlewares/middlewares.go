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

package middlewares

import "github.com/go-kit/kit/endpoint"

const (
	OPTIONAL_MIDDLEWARE_AUTH = "auth"
)

var (
	endpointMiddlewares []endpoint.Middleware
	optionalMiddlewares map[string]endpoint.Middleware
)

func init() {
	optionalMiddlewares = make(map[string]endpoint.Middleware)
}

func ImportEndpointMiddlewares(wares ...endpoint.Middleware) error {
	endpointMiddlewares = append(endpointMiddlewares, wares...)
	return nil
}

func ImportOptionalMiddleWares(name string, mw endpoint.Middleware) error {
	optionalMiddlewares[name] = mw
	return nil
}

func InjectMiddleWares(e endpoint.Endpoint) endpoint.Endpoint {
	for _, m := range endpointMiddlewares {
		e = m(e)
	}
	return e
}

func InjectOptionalMiddlewares(ep endpoint.Endpoint, list []string) endpoint.Endpoint {
	for _, name := range list {
		if m, ok := optionalMiddlewares[name]; ok {
			ep = m(ep)
		}
	}
	return ep
}
