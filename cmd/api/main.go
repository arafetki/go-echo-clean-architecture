package main

import (
	"flag"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/arafetki/go-echo-clean-architecture/internal/config"
	database "github.com/arafetki/go-echo-clean-architecture/internal/db"
	"github.com/arafetki/go-echo-clean-architecture/internal/env"
	"github.com/labstack/echo/v4"
	"github.com/lmittmann/tint"
)

type application struct {
	echo   *echo.Echo
	config config.Config
	logger *slog.Logger
	wg     sync.WaitGroup
}

func main() {

	var cfg config.Config
	flag.IntVar(&cfg.Server.Port, "port", 8080, "http server port")
	flag.Parse()

	cfg.App.Env = env.GetString("APP_ENV", "development")
	cfg.App.BaseURL = env.GetString("APP_BASE_URL", "http://localhost:8080")
	cfg.App.Version = env.GetString("APP_VERSION", "0.1.0")
	cfg.App.Debug = env.GetBool("APP_DEBUG", true)

	cfg.Server.ReadTimeoutDuration = env.GetInt("SERVER_READ_TIMEOUT_SECONDS", 10)
	cfg.Server.WriteTimeoutDuration = env.GetInt("SERVER_WRITE_TIMEOUT_SECONDS", 20)
	cfg.Server.IdleTimeoutDuration = env.GetInt("SERVER_IDLE_TIMEOUT_SECONDS", 60)
	cfg.Server.ShutdownPeriod = env.GetInt("SERVER_SHUTDOWN_PERIOD_SECONDS", 60)

	cfg.Database.Dsn = env.GetString("DATABASE_DSN", "local:localrootsecret@localhost:5432/postgres?sslmode=disable")
	cfg.Database.ConnTimeoutDuration = env.GetInt("DATABASE_CONN_TIMEOUT_SECONDS", 10)
	cfg.Database.AutoMigrate = env.GetBool("DATABASE_AUTO_MIGRATE", true)
	cfg.Database.MaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 25)
	cfg.Database.MaxIdleConns = env.GetInt("DATABASE_MAX_IDLE_CONNS", 25)
	cfg.Database.ConnMaxIdleDuration = env.GetInt("DATABASE_MAX_IDLE_DURATION_SECONDS", 300)
	cfg.Database.ConnMaxLifeDuration = env.GetInt("DATABASE_MAX_LIFE_DURATION_SECONDS", 3600)

	cfg.JWT.SecretKey = env.GetString("JWT_SECRET_KEY", "YK+8I3SYlMEs0cwv6CZ1QA5mYiltXcupiTJCEn/+7c4=")

	logLevel := slog.LevelInfo
	if cfg.App.Debug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel}))

	db, err := database.Postgres(cfg.Database.Dsn, database.Options{
		AutoMigrate:     cfg.Database.AutoMigrate,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnTimeout:     time.Duration(cfg.Database.ConnTimeoutDuration) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.Database.ConnMaxIdleDuration) * time.Second,
		ConnMaxLifetime: time.Duration(cfg.Database.ConnMaxLifeDuration) * time.Second,
	})

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection has been established successfully")

	app := &application{
		echo:   echo.New(),
		config: cfg,
		logger: logger,
	}

	app.registerRoutes()

	err = app.run()
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

}
