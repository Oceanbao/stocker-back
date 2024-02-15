package main

import (
	_ "embed"
	"log"
	"log/slog"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type Application struct {
	pb     *pocketbase.PocketBase
	logger *slog.Logger
}

func main() {
	isDevMode := false
	if os.Getenv("APP_ENV") == "dev" {
		isDevMode = true
	}

	var loggingLevel = new(slog.LevelVar)
	if isDevMode {
		loggingLevel.Set(slog.LevelDebug)
	}
	loggingOpt := &slog.HandlerOptions{
		Level: loggingLevel,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggingOpt))

	app := Application{
		pb:     pocketbase.New(),
		logger: logger,
	}
	app.logger.Info("starting app...")

	// // loosely check if it was executed using "go run".
	// isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	// // register migration command
	// migratecmd.MustRegister(app.pb, app.pb.RootCmd, migratecmd.Config{
	// 	Automigrate: isGoRun,
	// })

	// repoStock := infra.NewStockRepositoryPB(app.pb)
	// loggerSlog := infra.NewLoggerSlog(logger)
	// usecaseCommand := usecase.NewCommand(repoStock, loggerSlog)

	// ----------------- Route ----------------------
	app.pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Global middleware.
		// e.Router.Use(apis.RequireRecordAuth("user"))

		e.Router.GET("/dele", app.deleHandler)

		// e.Router.GET("/track", app.getTrackHandler, apis.RequireRecordAuth("users"))
		// e.Router.GET("/update-daily", app.updateDailyHandler, apis.RequireRecordAuth("users"))
		// e.Router.GET("/update-daily-etf", app.updateDailyETFHandler, apis.RequireRecordAuth("users"))

		return nil
	})

	// ----------------- Cron ----------------------
	// app.pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
	// 	if isDevMode {
	// 		app.logger.Debug("running in dev mode, turning off CRONs")
	// 		return nil
	// 	}

	// 	scheduler := cron.New()

	// 	// Every week Mon-Fri at 10:00 UTC (18:00 Beijing Time)
	// 	cronSignalDailyPriceUpdate := "0 10 * * 1-5"
	// 	err := scheduler.Add("daily", cronSignalDailyPriceUpdate, app.cronDailyPriceUpdate)
	// 	if err != nil {
	// 		return fmt.Errorf("error in adding cron job `dailyPrice`: %w", err)
	// 	}

	// 	// Every week Mon-Fri at 10:15 UTC (18:15 Beijing Time)
	// 	cronSignalDailySelectETFUpdate := "15 10 * * 1-5"
	// 	err = scheduler.Add("daily-etf", cronSignalDailySelectETFUpdate, app.cronDailySelectETFUpdate)
	// 	if err != nil {
	// 		return fmt.Errorf("error in adding cron job `dailyPriceETF`: %w", err)
	// 	}

	// 	// Every week Mon-Fri at 10:20 UTC (18:20 Beijing Time)
	// 	cronSignalDailyTallyUpdate := "20 10 * * 1-5"
	// 	err = scheduler.Add("daily-tally", cronSignalDailyTallyUpdate, app.cronDailyTallyUpdate)
	// 	if err != nil {
	// 		return fmt.Errorf("error in adding cron job `dailyTally`: %w", err)
	// 	}

	// 	scheduler.Start()

	// 	return nil
	// })

	if err := app.pb.Start(); err != nil {
		log.Fatal(err)
	}
}
