package main

func (app *Application) cronDailyDataUpdate() {
	if err := app.command.UpdateDailyData(); err != nil {
		app.pb.Logger().Error("cronDailyDataUpdate", "error", err.Error())
	}
}

func (app *Application) cronDailyScreening() {
	if err := app.command.UpdateDailyScreen(); err != nil {
		app.pb.Logger().Error("cronDailyScreening", "error", err.Error())
	}
}

func (app *Application) cronWeeklyStocksUpdate() {
	if err := app.command.UpdateStocks(); err != nil {
		app.pb.Logger().Error("cronWeeklyStocksUpdate", "error", err.Error())
	}
}
