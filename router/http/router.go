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

package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/leewckk/go-kit-micro-service/middlewares/transport/http"
)

var (
	DEFAULT_HTTP_VERSION_V1 = "/v1"
	DEFAULT_HTTP_HEART_BEAT = "/actuator"
	DEFAULT_HTTP_HEALTH_API = DEFAULT_HTTP_HEART_BEAT + "/info"
)

type Router struct {
	apiver        string
	healthApi     string
	healthGrp     string
	healthSvc     *http.RouterObject
	engine        *gin.Engine
	registerProcs []RegisterRouteProc
	serverOpts    []http.ServerOption
}

func NewRouter(bing string, routerOpts ...RouterOption) *Router {
	r := Router{
		apiver:        DEFAULT_HTTP_VERSION_V1,
		healthApi:     DEFAULT_HTTP_HEALTH_API,
		engine:        gin.Default(),
		serverOpts:    make([]http.ServerOption, 0, 0),
		registerProcs: make([]RegisterRouteProc, 0, 0),
	}

	for _, opt := range routerOpts {
		opt(&r)
	}
	return &r
}

type RouterOption func(r *Router)
type RegisterRouteProc func(r *Router, opts ...http.ServerOption)

//// config import routers
func RouterProcs(procs ...RegisterRouteProc) RouterOption {
	return func(r *Router) {
		r.registerProcs = append(r.registerProcs, procs...)
	}
}

/// server options
func RouterServerOptions(serverOpts ...http.ServerOption) RouterOption {
	return func(r *Router) {
		r.serverOpts = append(r.serverOpts, serverOpts...)
	}
}

//// heartbeat config
func RouterHealth(group string, healthAPI string, svc *http.RouterObject) RouterOption {
	return func(r *Router) {
		r.healthApi = healthAPI
		r.healthGrp = group
		r.healthSvc = svc
	}
}

/// api version config
func RouterVersion(ver string) RouterOption {
	return func(r *Router) {
		r.apiver = ver
	}
}

func (r *Router) HealthAPI() string {
	return r.healthGrp + r.healthApi
}

func (r *Router) Register(obj *http.RouterObject) error {

	uri := r.apiver + obj.Pattern

	if obj.HttpMethod == nethttp.MethodGet {
		r.engine.GET(uri, obj.Handler)
		return nil
	}

	if obj.HttpMethod == nethttp.MethodPut {
		r.engine.PUT(uri, obj.Handler)
		return nil
	}

	if obj.HttpMethod == nethttp.MethodPost {
		r.engine.POST(uri, obj.Handler)
		return nil
	}

	if obj.HttpMethod == nethttp.MethodPatch {
		r.engine.PATCH(uri, obj.Handler)
		return nil
	}

	if obj.HttpMethod == nethttp.MethodDelete {
		r.engine.DELETE(uri, obj.Handler)
		return nil
	}
	return nil
}

func (r *Router) Run(bind string) error {

	for _, proc := range r.registerProcs {
		proc(r, r.serverOpts...)
	}

	//// register health api
	if nil != r.healthSvc {
		health := r.engine.Group(r.healthGrp)
		health.GET(r.healthApi, r.healthSvc.Handler)
	}
	r.engine.Use(gin.Logger())
	return r.engine.Run(bind)
}
