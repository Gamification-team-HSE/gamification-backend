//go:generate mockgen -destination=./mocks/http_mock.go -package mocks github.com/speakeasy-api/rest-template-go/internal/transport/http Users,DB

package http

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/krespix/gamification-api/internal/core/metrics"
	"net/http"
)

// Server represents an HTTP server that can handle requests for this microservice.
type Server struct{}

// New will instantiate a new instance of Server.
func New() *Server {
	return &Server{}
}

// AddRoutes will add the routes this server supports to the router.
func (s *Server) AddRoutes(baseRouter *mux.Router) error {
	healthHandler := http.HandlerFunc(s.healthCheck)

	baseRouter.Use(middleware.Logger)
	baseRouter.Handle("/health", incrementIncomingRequestsMiddleware(healthHandler))
	baseRouter.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{
		Registry: metrics.Registry,
	}))

	return nil
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	handleResponse(r.Context(), w, "healthy")
}

func handleResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	jsonRes := struct {
		Data interface{} `json:"data"`
	}{
		Data: data,
	}

	dataBytes, err := json.Marshal(jsonRes)
	if err != nil {
		handleError(ctx, w, err)
		return
	}

	if _, err := w.Write(dataBytes); err != nil {
		handleError(ctx, w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func incrementIncomingRequestsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.IncomingHTTPRequestsTotal.With(prometheus.Labels{"method": r.Method, "uri": r.RequestURI}).Inc()
		next.ServeHTTP(w, r)
	})
}
