package main

import (
	"context"
	"github.com/getsentry/sentry-go"
	"gitlab.com/krespix/gamification-api/internal/config"
	"gitlab.com/krespix/gamification-api/internal/core/app"
	"gitlab.com/krespix/gamification-api/internal/core/drivers/psql"
	"gitlab.com/krespix/gamification-api/internal/core/listeners/http"
	"gitlab.com/krespix/gamification-api/internal/core/logging"
	httptransport "gitlab.com/krespix/gamification-api/internal/transport/http"
	"go.uber.org/zap"
)

func main() {
	app.Start(appStart)
}

func appStart(ctx context.Context, a *app.App) ([]app.Listener, error) {
	// Load configuration from config/config.yaml which contains details such as DB connection params
	cfg, err := config.Load(ctx)
	if err != nil {
		return nil, err
	}

	// Connect to the postgres DB
	db, err := initDatabase(ctx, cfg, a)
	if err != nil {
		return nil, err
	}

	//// Run our migrations which will update the DB or create it if it doesn't exist
	// if err := db.MigratePostgres(ctx, "./migrations"); err != nil {
	//	return nil, err
	// }
	// a.OnShutdown(func() {
	//	// Temp for development so database is cleared on shutdown
	//	if err := db.RevertMigrations(ctx, "./migrations"); err != nil {
	//		logging.From(ctx).Error("failed to revert migrations", zap.Error(err))
	//	}
	// })

	//init sentry
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		return nil, err
	}

	httpServer := httptransport.New(db.GetDB())

	// Create an HTTP server
	h, err := http.New(httpServer, cfg.HTTP)
	if err != nil {
		return nil, err
	}

	// Start listening for HTTP requests
	return []app.Listener{
		h,
	}, nil
}

func initDatabase(ctx context.Context, cfg *config.Config, a *app.App) (*psql.Driver, error) {
	db := psql.New(cfg.PSQL)

	if err := db.Connect(ctx); err != nil {
		return nil, err
	}

	a.OnShutdown(func() {
		// Shutdown connection when server terminated
		logging.From(ctx).Info("shutting down db connection")
		if err := db.Close(ctx); err != nil {
			logging.From(ctx).Error("failed to close db connection", zap.Error(err))
		}
	})

	return db, nil
}
