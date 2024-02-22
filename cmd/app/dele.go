package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *Application) deleHandler(c echo.Context) error {
	err := app.command.UpdateDailyScreen()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"dele": 0})
}

func (app *Application) updateDailyData(c echo.Context) error {
	go func() {
		err := app.command.UpdateDailyData()
		if err != nil {
			app.logger.Error("updateDailyData", "error", err.Error())
			app.notifier.Sendf("updateDailyData", fmt.Sprintf("error: %v", err.Error()))
		}
	}()
	return c.JSON(http.StatusOK, map[string]any{"updateDailyData": "ok"})
}
