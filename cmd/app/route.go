package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (app *Application) stockSearchHandler(c echo.Context) error {
	ticker := c.QueryParam("ticker")
	if ticker == "" {
		return c.JSON(http.StatusNotFound, map[string]any{"message": "require param missing"})
	}

	stock, err := app.query.GetStockByTicker(ticker)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"data": stock})
}

func (app *Application) stockCreateHandler(c echo.Context) error {
	ticker := c.QueryParam("ticker")
	if ticker == "" {
		return c.JSON(http.StatusNotFound, map[string]any{"message": "require param missing"})
	}

	// Ignore err since it is not nil if record non-exist.
	stock, _ := app.query.GetStockByTicker(ticker)
	if stock.Ticker != "" {
		return c.JSON(http.StatusOK, map[string]any{"message": "ticker already created"})
	}

	if err := app.command.CreateStockAndDailyData(ticker); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
}

func (app *Application) stockDeleteHandler(c echo.Context) error {
	ticker := c.QueryParam("ticker")
	if ticker == "" {
		return c.JSON(http.StatusNotFound, map[string]any{"message": "require param missing"})
	}

	if err := app.command.DeleteStockByTicker(ticker); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
}
