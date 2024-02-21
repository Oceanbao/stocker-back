package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"os"

	"example.com/stocker-back/internal/common"
	"example.com/stocker-back/internal/infra"
	"example.com/stocker-back/internal/usecase"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
)

type Application struct {
	pb       *pocketbase.PocketBase
	logger   *slog.Logger
	command  usecase.Command
	notifier common.Notifier
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

	pb := pocketbase.New()

	notifierPushbullet, err := infra.NewNotifierPushbullet()
	if err != nil {
		log.Fatal(err)
	}

	repoStock := infra.NewStockRepositoryPB(pb)
	repoScreen := infra.NewScreenRepositoryPB(pb)
	loggerSlog := infra.NewLoggerSlog(logger)
	usecaseCommand := usecase.NewCommand(repoStock, repoScreen, loggerSlog, notifierPushbullet)

	app := Application{
		pb:       pb,
		logger:   logger,
		command:  *usecaseCommand,
		notifier: notifierPushbullet,
	}

	app.logger.Info("starting app...")

	// // loosely check if it was executed using "go run".
	// isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	// // register migration command
	// migratecmd.MustRegister(app.pb, app.pb.RootCmd, migratecmd.Config{
	// 	Automigrate: isGoRun,
	// })

	// ----------------- Route ----------------------
	app.pb.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Global middleware.
		// e.Router.Use(apis.RequireRecordAuth("user"))

		e.Router.GET("/dele", app.deleHandler)
		e.Router.GET("/once", app.onceHandler)

		// e.Router.GET("/track", app.getTrackHandler, apis.RequireRecordAuth("users"))
		// e.Router.GET("/update-daily", app.updateDailyHandler, apis.RequireRecordAuth("users"))
		// e.Router.GET("/update-daily-etf", app.updateDailyETFHandler, apis.RequireRecordAuth("users"))

		return nil
	})

	// ----------------- Cron ----------------------
	app.pb.OnBeforeServe().Add(func(_ *core.ServeEvent) error {
		// if isDevMode {
		// 	app.logger.Debug("running in dev mode, turning off CRONs")
		// 	return nil
		// }

		scheduler := cron.New()

		// Every week Mon-Fri at 10:00 UTC (18:00 Beijing Time)
		cronSignalDailyDataUpdate := "0 10 * * 1-5"
		err := scheduler.Add("dailydata", cronSignalDailyDataUpdate, app.cronDailyDataUpdate)
		if err != nil {
			return fmt.Errorf("error in adding cron job `cronSignalDailyDataUpdate`: %w", err)
		}
		app.logger.Info("cron", "messge", "cronSignalDailyDataUpdate registered")

		scheduler.Start()

		return nil
	})

	if err := app.pb.Start(); err != nil {
		log.Fatal(err)
	}
}
