package main

func (app *Application) cronDailyDataUpdate() {
	if err := app.command.UpdateDailyData(); err != nil {
		app.logger.Error("cronDailyDataUpdate", "error", err.Error())
	}
}
