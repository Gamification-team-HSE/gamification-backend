//go:generate mockgen -destination=./mocks/http_mock.go -package mocks github.com/speakeasy-api/rest-template-go/internal/transport/http Users,DB

package http

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/krespix/gamification-api/internal/core/metrics"
	"net/http"
)

// DB represents a type that can be used to interact with the database.
type DB interface {
	PingContext(ctx context.Context) error
}

// Server represents an HTTP server that can handle requests for this microservice.
type Server struct {
	db DB
}

// New will instantiate a new instance of Server.
func New(db DB) *Server {
	return &Server{
		db: db,
	}
}

// AddRoutes will add the routes this server supports to the router.
func (s *Server) AddRoutes(r *chi.Mux) error {
	healthHandler := http.HandlerFunc(s.healthCheck)

	r.Use(middleware.Logger)

	r.Handle("/health", incrementIncomingRequestsMiddleware(healthHandler))
	r.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{
		Registry: metrics.Registry,
	}))

	return nil
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.db.PingContext(r.Context()); err != nil {
		handleError(r.Context(), w, err)
		return
	}
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
