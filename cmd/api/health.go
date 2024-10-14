package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (app *application) healthCheckHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, echo.Map{
		"timestamp": time.Now().Nanosecond(),
	})
}
