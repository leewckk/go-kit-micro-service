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
	"fmt"

	"github.com/leewckk/go-kit-micro-service/discovery"
	"github.com/leewckk/go-kit-micro-service/middlewares/tracing/report"
	"github.com/leewckk/go-kit-micro-service/middlewares/transport/http"
	router "github.com/leewckk/go-kit-micro-service/router/http"
)

func Run(r *router.Router, port int, serverName string, center string, centerPort int, sdapi string) chan error {
	errch := make(chan error)
	go func() {
		if port <= 0 {
			errch <- fmt.Errorf("Invalid server port: %v", port)
		}
		opts := make([]http.ServerOption, 0, 0)
		if "" != sdapi {
			tr := report.NewZipkinReporter(sdapi)
			opts = append(opts, report.NewGinServerTracer(tr))
		}

		healthAPI := ""
		pub, err := Publish(port, serverName, center, centerPort, healthAPI)
		if nil == err && nil != pub {
			defer pub.UnRegister()
		}
		errch <- r.Run(fmt.Sprintf(":%v", port))
	}()
	return errch
}

func Publish(port int, serverName string, center string, centerPort int, healthApi string) (discovery.Publisher, error) {
	return discovery.Publish(center, centerPort, serverName, port, healthApi, "http")
}
