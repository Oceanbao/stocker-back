package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"os"

	"example.com/stocker-back/internal/infra"
	"example.com/stocker-back/internal/usecase"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
)

type Application struct {
	pb       *pocketbase.PocketBase
	command  usecase.Command
	query    usecase.Query
	logger   *slog.Logger
	notifier infra.Notifier
}

func main() { //nolint:funlen //ignore
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
	repoTracking := infra.NewTrackingRepositoryPB(pb)
	loggerSlog := infra.NewLoggerSlog(logger)
	usecaseCommand := usecase.NewCommand(repoStock, repoScreen, repoTracking, loggerSlog, notifierPushbullet)
	usecaseQuery := usecase.NewQuery(repoStock, repoScreen, repoTracking, loggerSlog, notifierPushbullet)

	app := Application{
		pb:       pb,
		command:  *usecaseCommand,
		query:    *usecaseQuery,
		logger:   logger,
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

		gDele := e.Router.Group("/dele")
		gDele.Use(apis.RequireRecordAuth("users"))
		gDele.GET("/dev", app.deleHandler)
		gDele.GET("/updatedaily", app.updateDailyData)
		gDele.GET("/updatescreen", app.screenUpdateHandler)

		gStock := e.Router.Group("/stocks")
		gStock.Use(apis.RequireRecordAuth("users"))
		gStock.GET("/:ticker", app.stockSearchHandler)
		gStock.POST("/:ticker", app.stockCreateHandler)
		gStock.DELETE("/:ticker", app.stockDeleteHandler)

		gTracking := e.Router.Group("/tracking")
		gTracking.Use(apis.RequireRecordAuth("users"))
		gTracking.GET("/all", app.trackingSearchHandler)
		gTracking.POST("/:ticker", app.trackingCreateHandler)

		e.Router.GET("/screen", app.screenReadHandler, apis.RequireRecordAuth("users"))

		e.Router.GET("/random/:num", app.randomStocksHandler, apis.RequireRecordAuth("users"))

		return nil
	})

	// ----------------- Cron ----------------------
	app.pb.OnBeforeServe().Add(func(_ *core.ServeEvent) error {
		if isDevMode {
			app.logger.Debug("running in dev mode, turning off CRONs")
			return nil
		}

		scheduler := cron.New()

		// Every week Mon-Fri at 10:00 UTC (18:00 Beijing Time)
		err := scheduler.Add("dailydata", "0 10 * * 1-5", app.cronDailyDataUpdate)
		if err != nil {
			return fmt.Errorf("error in adding cron job `cronSignalDailyDataUpdate`: %w", err)
		}
		app.logger.Info("cron", "messge", "cronSignalDailyDataUpdate registered")

		// Every week Mon-Fri at 11:00 UTC (19:00 Beijing Time)
		err = scheduler.Add("dailyscreen", "0 11 * * 1-5", app.cronDailyScreening)
		if err != nil {
			return fmt.Errorf("error in adding cron job `cronDailyScreening`: %w", err)
		}
		app.logger.Info("cron", "messge", "cronDailyScreening registered")

		// Every week Fri at 12:00 UTC (20:00 Beijing Time)
		err = scheduler.Add("weeklystocks", "0 12 * * 5", app.cronWeeklyStocksUpdate)
		if err != nil {
			return fmt.Errorf("error in adding cron job `cronWeeklyStocksUpdate`: %w", err)
		}
		app.logger.Info("cron", "messge", "cronWeeklyStocksUpdate registered")

		scheduler.Start()

		return nil
	})

	if err := app.pb.Start(); err != nil {
		log.Fatal(err)
	}
}
