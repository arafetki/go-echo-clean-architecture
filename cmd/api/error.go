package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) httpErrorHandler(err error, c echo.Context) {

	code := http.StatusInternalServerError
	message := "the server encountered a problem and could not process your request"

	httpError, ok := err.(*echo.HTTPError)
	if ok {
		code = httpError.Code
		switch code {
		case http.StatusNotFound:
			message = "the requested resource could not be found"
		case http.StatusMethodNotAllowed:
			message = fmt.Sprintf("the %s method is not supported for this resource", c.Request().Method)
		case http.StatusBadRequest:
			message = "The request could not be understood by the server due to malformed syntax or incorrect parameter type"
		case http.StatusInternalServerError:
			message = "the server encountered a problem and could not process your request"
		case http.StatusUnauthorized:
			c.Response().Header().Set("WWW-Authenticate", `Bearer realm="restricted", charset="UTF-8"`)
			message = "You must be authenticated to access this resource"
		default:
			message = httpError.Message.(string)
		}
	}

	if !c.Response().Committed {
		c.JSON(code, echo.Map{"error": message})
	}
}
