package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *Application) deleHandler(c echo.Context) error {
	// err := app.command.UpdateDailyScreen()
	// err := app.command.UpdateStocks()
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	// }
	return c.JSON(http.StatusOK, map[string]any{"dele": 0})
}

func (app *Application) onceHandler(c echo.Context) error {
	go func() {
		err := app.command.UpdateStocks()
		if err != nil {
			app.logger.Error("onceHandler", "error", err.Error())
			app.notifier.Sendf("onceHandler", fmt.Sprintf("error: %v", err.Error()))
		}
	}()
	return c.JSON(http.StatusOK, map[string]any{"once": "ok"})
}
