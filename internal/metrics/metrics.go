package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	IncomingHTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "incoming_http_requests_total",
		Help: "The total number of incoming http requests",
	}, []string{"uri", "method"})
)
