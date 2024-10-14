package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) run() error {

	app.echo.Server.ReadTimeout = time.Duration(app.config.Server.ReadTimeoutDuration) * time.Second
	app.echo.Server.WriteTimeout = time.Duration(app.config.Server.WriteTimeoutDuration) * time.Second
	app.echo.Server.IdleTimeout = time.Duration(app.config.Server.IdleTimeoutDuration) * time.Second
	app.echo.Server.ErrorLog = slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn)

	shutdownErrChan := make(chan error)

	go func() {

		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(app.config.Server.ShutdownPeriod)*time.Second)
		defer cancel()

		shutdownErrChan <- app.echo.Shutdown(ctx)

	}()

	err := app.echo.Start(fmt.Sprintf(":%d", app.config.Server.Port))
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrChan
	if err != nil {
		return err
	}

	app.logger.Warn("server has stopped")

	app.wg.Wait()

	return nil
}
