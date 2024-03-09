package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// stockSearchHandler is controller handling stock search of single ticker.
func (app *Application) stockSearchHandler(c echo.Context) error {
	ticker := c.PathParam("ticker")

	stock, err := app.query.GetStockByTicker(ticker)
	if err != nil {
		return c.JSON(http.StatusOK, ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseData(stock))
}

// stockCreateHandler is controller handling stock creation of single ticker.
func (app *Application) stockCreateHandler(c echo.Context) error {
	ticker := c.PathParam("ticker")

	// Ignore err since it == nil if record exists.
	stock, _ := app.query.GetStockByTicker(ticker)
	if stock.Ticker != "" {
		return c.JSON(http.StatusOK, ResponseErr("ticker already exists"))
	}

	if err := app.command.CreateStockAndDailyData(ticker); err != nil {
		return c.JSON(http.StatusOK, ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseOk())
}

// stockDeleteHandler is controller handling stock deletion of single ticker.
func (app *Application) stockDeleteHandler(c echo.Context) error {
	ticker := c.PathParam("ticker")

	if err := app.command.DeleteStockByTicker(ticker); err != nil {
		return c.JSON(http.StatusOK, ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseOk())
}

// screenUpdateHandler is controller handling update daily screening.
func (app *Application) screenUpdateHandler(c echo.Context) error {
	err := app.command.UpdateDailyScreen()
	if err != nil {
		return c.JSON(http.StatusOK, ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseOk())
}

// screenReadHandler is controller handling retrieval of daily screens.
func (app *Application) screenReadHandler(c echo.Context) error {
	data, err := app.query.GetScreens()
	if err != nil {
		return c.JSON(http.StatusOK, ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseData(data))
}

// trackingSearchHandler is controller getting all trackings.
func (app *Application) trackingSearchHandler(c echo.Context) error {
	data, err := app.query.GetTrackings()
	if err != nil {
		return c.JSON(http.StatusOK, ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseData(data))
}

// trackStockHandler is controller adding stock to track collection.
func (app *Application) trackingCreateHandler(c echo.Context) error {
	ticker := c.PathParam("ticker")

	if err := app.command.CreateTracking(ticker); err != nil {
		return c.JSON(http.StatusOK, ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseOk())
}

// trackingDeleteHandler is controller deleting stock from tracking.
func (app *Application) trackingDeleteHandler(c echo.Context) error {
	ticker := c.PathParam("ticker")

	if err := app.command.DeleteTracking(ticker); err != nil {
		return c.JSON(http.StatusOK, ResponseErr(err.Error()))
	}

	return c.JSON(http.StatusOK, ResponseOk())
}

// sectorReadHandler is controller handling retrieval of stocks by given sector.
func (app *Application) sectorReadHandler(c echo.Context) error {
	sector := c.PathParam("sector")

	stocks, err := app.query.GetStocksBySector(sector)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ResponseData(stocks))
}

// randomStocksHandler is controller handling retrieval of random stocks given number.
func (app *Application) randomStocksHandler(c echo.Context) error {
	numStr := c.PathParam("num")
	numInt, _ := strconv.Atoi(numStr)

	tickers, err := app.query.GetRandomTickers(numInt)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ResponseData(tickers))
}
