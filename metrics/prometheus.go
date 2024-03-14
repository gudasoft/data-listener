package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "handler"},
	)

	requestSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "Size of HTTP requests.",
		},
		[]string{"method", "handler"},
	)

	errorCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "app_errors_total",
			Help: "Total number of errors logged.",
		},
	)

	startTime time.Time
)

func init() {
	prometheus.MustRegister(requestCount, requestSize, errorCount)
}

func RecordRequestMetrics(method, handler string, size int) {
	requestCount.WithLabelValues(method, handler).Inc()
	requestSize.WithLabelValues(method, handler).Observe(float64(size))
}

func RecordError() {
	errorCount.Inc()
}

func RunMetricsServer(addr string, port int, path string) {
	http.Handle(path, promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
}
