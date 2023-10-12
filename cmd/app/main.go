package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
)

func main() {
	app := pocketbase.New()

	// ----------------- Route ----------------------
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Global middleware.
		// e.Router.Use(apis.RequireRecordAuth("user"))

		routeTrack(e, app)
		routeUpdateDaily(e, app)
		routeUpdateDailyETF(e, app)

		return nil
	})

	// ----------------- Cron ----------------------
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		scheduler := cron.New()

		// Every week Mon-Fri at 10:00 UTC (18:00 Beijing Time)
		err := scheduler.Add("daily", "0 10 * * 1-5", func() {
			cronDailyPriceUpdate(app)
		})
		if err != nil {
			return fmt.Errorf("error in adding cron job `dailyPrice`: %w", err)
		}

		// Every week Mon-Fri at 10:15 UTC (18:15 Beijing Time)
		err = scheduler.Add("daily-etf", "15 10 * * 1-5", func() {
			cronDailySelectETFUpdate(app)
		})
		if err != nil {
			return fmt.Errorf("error in adding cron job `dailyPriceETF`: %w", err)
		}

		// Every week Mon-Fri at 10:20 UTC (18:20 Beijing Time)
		err = scheduler.Add("daily-tally", "20 10 * * 1-5", func() {
			cronDailyTallyUpdate(app)
		})
		if err != nil {
			return fmt.Errorf("error in adding cron job `dailyTally`: %w", err)
		}

		scheduler.Start()

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
