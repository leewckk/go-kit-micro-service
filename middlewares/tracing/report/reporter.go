package report

import (
	"github.com/openzipkin/zipkin-go/reporter"
	"github.com/openzipkin/zipkin-go/reporter/http"
)

func init() {
}

func NewZipkinReporter(reporter string) reporter.Reporter {
	if reporter != "" {
		return http.NewReporter(reporter)
	}
	return nil
}
