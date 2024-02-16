package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *Application) deleHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": 0})
}
