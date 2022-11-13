package metrics

import (
	"net/http"
	"strconv"

	"allandeng.cn/allandeng/go-blog/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var TotalRequests *prometheus.CounterVec
var RequestDurations *prometheus.HistogramVec

func Init() {
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "The total number of handled HTTP requests.",
		},
		[]string{"path", "method", "code"},
	)

	RequestDurations = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "A histogram of the HTTP request durations in seconds.",
			// Bucket 配置：第一个 bucket 包括所有在 0.5ms 内完成的请求，最后一个包括所有在10s内完成的请求。
			Buckets: []float64{0.5, 1, 2.5, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 100000},
		},
		[]string{"path", "method", "code"},
	)
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(RequestDurations)

	go http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(config.GlobalConfig.WebServer.MetricsPort)), promhttp.Handler())
}
