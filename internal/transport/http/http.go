//go:generate mockgen -destination=./mocks/http_mock.go -package mocks github.com/speakeasy-api/rest-template-go/internal/transport/http Users,DB

package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
func (s *Server) AddRoutes(r *mux.Router) error {
	r.HandleFunc("/health", s.healthCheck).Methods(http.MethodGet)

	_ = r.PathPrefix("/v1").Subrouter()

	return nil
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.db.PingContext(r.Context()); err != nil {
		handleError(r.Context(), w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("app healthy"))
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
