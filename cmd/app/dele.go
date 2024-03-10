package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *Application) deleUpdateStocksHandler(c echo.Context) error {
	go func() {
		app.pb.Logger().Info("deleUpdateStocksHandler running...")
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
			app.pb.Logger().Error("updateDailyData", "error", err.Error())
			app.notifier.Sendf("updateDailyData", fmt.Sprintf("error: %v", err.Error()))
		}
	}()
	return c.JSON(http.StatusOK, ResponseOk())
}

func (app *Application) deleDevHandler(c echo.Context) error {
	_, err := app.query.GetStocksBySector("dele")
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ResponseOk())
}

func (app *Application) deleteStocksHandler(c echo.Context) error {
	payload := struct {
		Tickers []string `json:"tickers"`
	}{
		Tickers: nil,
	}
	err := json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return c.JSON(http.StatusOK, ResponseErr(err.Error()))
	}

	if len(payload.Tickers) == 0 {
		return c.JSON(http.StatusOK, ResponseErr("missing tickers"))
	}

	failedTickers := make([]string, 0)
	for _, t := range payload.Tickers {
		if err := app.command.DeleteStockByTicker(t); err != nil {
			failedTickers = append(failedTickers, t)
		}
	}

	if len(failedTickers) != 0 {
		return c.JSON(http.StatusOK, ResponseErr(fmt.Sprintf("failed tickers: %v", failedTickers)))
	}

	return c.JSON(http.StatusOK, ResponseOk())
}
