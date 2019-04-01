package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reqsStatus = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "The total number of requests which were performed.",
	}, []string{"method", "path", "status"})

	reqsCurrent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_requests_current",
		Help: "The current number of requests in course.",
	}, []string{"method", "path"})

	reqsDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "The duration of the requests in seconds.",
	}, []string{"method", "path"})
)

func PrometheusHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func MetricsMiddleware(path string) gin.HandlerFunc {
	return func(g *gin.Context) {
		labels := prometheus.Labels{"method": g.Request.Method, "path": path}
		reqsCurrent.With(labels).Inc()
		startTime := time.Now()
		g.Next()

		reqsDuration.With(labels).Observe(time.Since(startTime).Seconds())
		reqsCurrent.With(labels).Dec()

		labels["status"] = strconv.Itoa(g.Writer.Status())
		reqsStatus.With(labels).Inc()
	}
}

func RegisterCustomMetrics() {
	prometheus.MustRegister(reqsStatus)
	prometheus.MustRegister(reqsCurrent)
	prometheus.MustRegister(reqsDuration)
}
