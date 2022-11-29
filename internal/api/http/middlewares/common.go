package middlewares

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/krespix/gamification-api/internal/core/metrics"
	"net/http"
)

func IncrementIncomingRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.IncomingHTTPRequestsTotal.With(prometheus.Labels{"method": r.Method, "uri": r.RequestURI}).Inc()
		next.ServeHTTP(w, r)
	})
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}
