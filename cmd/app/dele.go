package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *Application) deleUpdateStocksHandler(c echo.Context) error {
	go func() {
		err := app.command.UpdateStocks()
		if err != nil {
			return
		}
	}()

	return c.JSON(http.StatusOK, ResponseOk())
}

func (app *Application) updateDailyData(c echo.Context) error {
	go func() {
		err := app.command.UpdateDailyData()
		if err != nil {
			app.logger.Error("updateDailyData", "error", err.Error())
			app.notifier.Sendf("updateDailyData", fmt.Sprintf("error: %v", err.Error()))
		}
	}()
	return c.JSON(http.StatusOK, ResponseOk())
}
