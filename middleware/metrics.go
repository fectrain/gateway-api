package middleware

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ApiSum = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ApiDuration",
		Help: "api request duration ms",
	}, []string{"api"})
	ApiCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ApiCount",
		Help: "api request count",
	}, []string{"api"})
	ErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ApiErrorCount",
		Help: "api request error count: api/ws",
	}, []string{"type"})

)

func init() {
	prometheus.MustRegister(ApiCount, ApiSum, ErrorCount)
}
