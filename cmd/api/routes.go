package main

import (
	"log/slog"

	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

func (app *application) registerRoutes() {

	app.echo.HTTPErrorHandler = app.httpErrorHandler

	app.echo.Use(slogecho.NewWithConfig(app.logger, slogecho.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,
	}))
	app.echo.Use(middleware.Recover())
	app.echo.Use(middleware.BodyLimit("2M"))

	// For Load balancers
	app.echo.GET("/health", app.healthCheckHandler)
}
