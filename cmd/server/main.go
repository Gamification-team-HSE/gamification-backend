package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/pressly/goose/v3"
	"gitlab.com/krespix/gamification-api/internal/api/graphql/resolvers"
	httpAPI "gitlab.com/krespix/gamification-api/internal/api/http"
	"gitlab.com/krespix/gamification-api/internal/clients/smtp"
	"gitlab.com/krespix/gamification-api/internal/config"
	"gitlab.com/krespix/gamification-api/internal/core/app"
	"gitlab.com/krespix/gamification-api/internal/core/listeners/http"
	"gitlab.com/krespix/gamification-api/internal/core/logging"
	"gitlab.com/krespix/gamification-api/internal/repositories/cache"
	authRepository "gitlab.com/krespix/gamification-api/internal/repositories/cache/auth"
	"gitlab.com/krespix/gamification-api/internal/repositories/postgres"
	userRepository "gitlab.com/krespix/gamification-api/internal/repositories/postgres/user"
	authService "gitlab.com/krespix/gamification-api/internal/services/auth"
	userService "gitlab.com/krespix/gamification-api/internal/services/user"
	"go.uber.org/zap"
	"os"
	"time"
)

const defaultConfigPath = "config/config.yaml"

var (
	// configName - имя файла конфигурации.
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", defaultConfigPath, "Configuration file path")
	flag.Parse()
}

func main() {
	app.Start(appStart)
}

func appStart(ctx context.Context, a *app.App) ([]app.Listener, error) {
	cfg, err := config.Load(ctx, configPath)
	if err != nil {
		return nil, err
	}

	// Connect to the postgres DB
	db, err := initDatabase(ctx, cfg, a)
	if err != nil {
		return nil, err
	}
	//init sentry
	if cfg.Sentry.Enabled {
		err = sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Sentry.DSN,
			TracesSampleRate: 1.0,
		})
		if err != nil {
			return nil, err
		}
	}

	validate := validator.New()

	//init clients
	smtpClient := smtp.New(cfg.SMTP)
	cacheClient := cache.New(time.Second*60, time.Second*70)

	//init repositories
	userRepo := userRepository.New(db)
	authRepo := authRepository.New(cacheClient)

	//init services
	userSrc := userService.New(userRepo, validate)
	authSrc := authService.New(smtpClient, userRepo, authRepo, validate, cfg.JWT.Secret, time.Hour*24)

	resolver := resolvers.New(userSrc, authSrc)
	httpServer := httpAPI.New(resolver, authSrc)

	//init super admin
	err = userSrc.InitSuperAdmin(ctx, cfg.SuperAdmin)
	if err != nil {
		return nil, err
	}

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

func initDatabase(ctx context.Context, cfg *config.Config, a *app.App) (*postgres.Client, error) {
	db, err := postgres.New(cfg.DB)
	if err != nil {
		return nil, err
	}

	if err := db.Connect(ctx); err != nil {
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}

	currDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if err := goose.Up(db.GetDB(), fmt.Sprintf("%s/migrations", currDir)); err != nil {
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
