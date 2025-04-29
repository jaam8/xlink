package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var (
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "route"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(httpRequests)
}

func RecordRequest(method string, path string, durationSeconds float64) {
	httpDuration.WithLabelValues(method, path).Observe(durationSeconds)
	httpRequests.WithLabelValues(method, path).Inc()
}

func MetricsMiddlewareHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()

		RecordRequest(r.Method, r.URL.Path, duration)
	})
}
