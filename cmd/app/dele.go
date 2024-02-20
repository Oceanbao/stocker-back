package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *Application) deleHandler(c echo.Context) error {
	err := app.command.UpdateDailyData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"data": 0})
}

func (app *Application) onceHandler(c echo.Context) error {
	go func() {
		err := app.command.UpdateDailyData()
		if err != nil {
			app.logger.Error("onceHandler", "error", err.Error())
		}
	}()
	return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
}
