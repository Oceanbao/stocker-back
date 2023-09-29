package main

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
)

type dailyRecord struct {
	ID    string  `db:"id" json:"id"`
	Code  string  `db:"code" json:"code"`
	Date  string  `db:"date" json:"date"`
	Open  float64 `db:"open" json:"open"`
	High  float64 `db:"high" json:"high"`
	Low   float64 `db:"low" json:"low"`
	Close float64 `db:"close" json:"close"`
}

func main() {
	app := pocketbase.New()

	// ----------------- Route ----------------------
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Global middleware.
		// e.Router.Use(apis.RequireRecordAuth("user"))

		routeDele(e, app)
		routeTrack(e, app)

		return nil
	})

	// ----------------- Cron ----------------------
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		scheduler := cron.New()

		// Every week Mon-Fri at 00:00
		// err := scheduler.Add("daily", "*/1 * * * *", func() {
		err := scheduler.Add("daily", "0 0 * * *", func() {
			cronDailyPriceUpdate(app)
		})
		if err != nil {
			return fmt.Errorf("error in adding cron job `dailyPrice`: %w", err)
		}

		scheduler.Start()

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
