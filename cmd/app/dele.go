package main

import (
	"net/http"

	"example.com/stocker-back/internal/infra"
	"example.com/stocker-back/internal/usecase"
	"github.com/labstack/echo/v5"
)

func (app *Application) deleHandler(c echo.Context) error {
	repoStock := infra.NewStockRepositoryPB(app.pb)
	loggerSlog := infra.NewLoggerSlog(app.logger)
	usecaseCommand := usecase.NewCommand(repoStock, loggerSlog)
	if err := usecaseCommand.UpdateDailyData(); err != nil {
		return c.JSON(http.StatusOK, map[string]any{"error": err})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
