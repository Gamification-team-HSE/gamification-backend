package http

import (
	"context"
	"encoding/json"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"gitlab.com/krespix/gamification-api/internal/api/graphql/directives"
	"gitlab.com/krespix/gamification-api/internal/api/http/middlewares"
	"gitlab.com/krespix/gamification-api/internal/core/metrics"
	"gitlab.com/krespix/gamification-api/internal/services/auth"
	"gitlab.com/krespix/gamification-api/pkg/graphql/server"
	"net/http"
)

// Server represents an HTTP server that can handle requests for this microservice.
type Server struct {
	resolver    server.ResolverRoot
	authService auth.Service
}

// New will instantiate a new instance of Server.
func New(resolver server.ResolverRoot, authSvc auth.Service) *Server {
	return &Server{
		resolver:    resolver,
		authService: authSvc,
	}
}

// AddRoutes will add the routes this server supports to the router.
func (s *Server) AddRoutes(baseRouter *mux.Router) error {
	healthHandler := http.HandlerFunc(s.healthCheck)

	baseRouter.Use(middleware.Logger)
	baseRouter.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{
		Registry: metrics.Registry,
	}))

	schema := server.NewExecutableSchema(server.Config{
		Resolvers:  s.resolver,
		Directives: server.DirectiveRoot{Auth: directives.Auth},
	})
	srv := handler.NewDefaultServer(schema)

	baseRouter.Use(middlewares.EnableCORS)
	baseRouter.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	authmw := middlewares.NewAuth(s.authService)

	graphqlRouter := baseRouter.PathPrefix("/gapi/v1").Subrouter()
	apiSubRouter := baseRouter.PathPrefix("/api").Subrouter()
	v1SubRouter := apiSubRouter.PathPrefix("/v1").Subrouter()

	graphqlRouter.Use(authmw.AuthMiddleware)

	pgHandler := playground.Handler("gamification-api", "/gapi/v1/query")
	graphqlRouter.Handle("/query", srv)
	graphqlRouter.Handle("/playground", pgHandler)

	v1SubRouter.Use(middlewares.IncrementIncomingRequestsMiddleware)
	v1SubRouter.Handle("/health", healthHandler).Methods(http.MethodGet)

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
