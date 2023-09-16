package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
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

	// wg := sync.WaitGroup{}

	defer rootCtxCancelFunc()
}
