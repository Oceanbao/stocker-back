package main

func (app *Application) cronDailyDataUpdate() {
	if err := app.command.UpdateDailyData(); err != nil {
		app.logger.Error("cronDailyDataUpdate", "error", err.Error())
	}
}

func (app *Application) cronDailyScreening() {
	if err := app.command.UpdateDailyScreen(); err != nil {
		app.logger.Error("cronDailyScreening", "error", err.Error())
	}
}

func (app *Application) cronWeeklyStocksUpdate() {
	if err := app.command.UpdateStocks(); err != nil {
		app.logger.Error("cronWeeklyStocksUpdate", "error", err.Error())
	}
}
