package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestDurationHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Request duration distribution",
		Buckets: []float64{1, 2, 5, 10, 20, 60},
	},
	[]string{"method"},
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {

	start := time.Now()

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	duration := time.Since(start)
	requestDurationHistogram.WithLabelValues(req.Method).Observe(duration.Seconds())
}

func main() {

	prometheus.MustRegister(requestDurationHistogram)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)

	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8090", nil)
}
