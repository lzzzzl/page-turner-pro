package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/gin-gonic/gin"
	"github.com/lzzzzl/page-turner-pro/internal/app"
	"github.com/lzzzzl/page-turner-pro/internal/app/handlers"
	"github.com/lzzzzl/page-turner-pro/internal/app/middleware"
	"github.com/rs/zerolog"
)

var (
	AppName    = "Page-Turner-Pro"
	AppVersion = "unknown_version"
	AppBuild   = "unknown_build"
)

const (
	defaultEnv      = "staging"
	defaultLogLevel = "info"
	defaultPort     = "9000"
)

type AppConfig struct {
	// General configuration
	Env      *string
	LogLevel *string

	// Database configuration
	DatabaseDSN *string

	// HTTP configuration
	Port *int
}

func runHTTPServer(rootCtx context.Context, wg *sync.WaitGroup, port int, app *app.Application) {
	// Set to release mode to disabled Gin logger
	gin.SetMode(gin.ReleaseMode)

	// Create gin router
	ginRouter := gin.New()

	// Set general middleware
	middleware.SetGeneralMiddlewares(rootCtx, ginRouter)

	// Register all handlers
	handlers.RegisterHandlers(ginRouter, app)

	// Build HTTP server
	httpAddr := fmt.Sprintf("0.0.0.0:%d", port)
	server := &http.Server{
		Addr:    httpAddr,
		Handler: ginRouter,
	}

	// Run the server in a goroutine
	go func() {
		zerolog.Ctx(rootCtx).Info().Msgf("HTTP server is on http://%s", httpAddr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			zerolog.Ctx(rootCtx).Panic().Err(err).Str("addr", httpAddr).Msg("fail to satrt HTTP server")
		}
	}()

	// Wait for rootCtx done
	go func() {
		<-rootCtx.Done()

		// Graceful shutdown http server with a timeout
		zerolog.Ctx(rootCtx).Info().Msgf("HTTP server is closing")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("fail to shutdown HTTP server")
		}

		// Notify when server is closed
		zerolog.Ctx(rootCtx).Info().Msgf("HTTP server is closed")
		wg.Done()
	}()
}

func initAppConfig() AppConfig {
	// Setup basic application information
	app := kingpin.New(AppName, "Page Turner PRO Server").Version(fmt.Sprintf("version: %s, build: %s", AppVersion, AppBuild))
	var config AppConfig

	config.Env = app.
		Flag("env", "The running environment").
		Envar("ENV").Default(defaultEnv).Enum("staging", "production")

	config.LogLevel = app.
		Flag("log_level", "Log filtering level").
		Envar("LOG_LEVEL").Default(defaultLogLevel).Enum("error", "warn", "info", "debug", "disabled")

	config.Port = app.
		Flag("port", "The HTTP server port").
		Envar("PORT").Default(defaultPort).Int()

	config.DatabaseDSN = app.
		Flag("database_dsn", "The database DSN").
		Envar("DATABASE_DSN").Required().String()

	kingpin.MustParse(app.Parse(os.Args[1:]))

	return config
}

func initRootLogger(levelStr, env string) zerolog.Logger {
	// Set global Log Level
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Set logger time format
	const rfc3339Micro = "2006-01-02T15:04:05.000000Z07:00"
	zerolog.TimeFieldFormat = rfc3339Micro

	serviceName := fmt.Sprintf("%s-%s", AppName, env)
	rootLogger := zerolog.New(os.Stdout).With().
		Timestamp().Str("service", serviceName).Logger()

	return rootLogger
}

func main() {
	// Setup app configuration
	cfg := initAppConfig()

	// Create root Logger
	rootLogger := initRootLogger(*cfg.LogLevel, *cfg.Env)

	// Create root context
	rootCtx, rootCtxCancelFunc := context.WithCancel(context.Background())
	rootCtx = rootLogger.WithContext(rootCtx)

	rootLogger.Info().
		Str("version", AppVersion).
		Str("build", AppBuild).
		Msgf("Lanuching %s", AppName)

	wg := sync.WaitGroup{}

	// Create application
	app := app.MustNewApplication(rootCtx, &wg, app.ApplicationParams{
		Env:         *cfg.Env,
		DatabaseDSN: *cfg.DatabaseDSN,
	})

	// Run server
	wg.Add(1)
	runHTTPServer(rootCtx, &wg, *cfg.Port, app)

	// Listen to SIGTERM/SIGINT to close
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	<-gracefulStop
	rootCtxCancelFunc()

	// Wait for all services to close with a specific timeout
	var waitUnitlDone = make(chan struct{})
	go func() {
		wg.Wait()
		close(waitUnitlDone)
	}()
	select {
	case <-waitUnitlDone:
		rootLogger.Info().Msg("success to close all services")
	case <-time.After(10 * time.Second):
		rootLogger.Err(context.DeadlineExceeded).Msg("fail to close all services")
	}
}
