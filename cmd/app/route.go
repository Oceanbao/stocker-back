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
