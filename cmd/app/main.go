package main

import (
	_ "embed"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
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
		routeHello(e)

		return nil
	})

	// ----------------- Cron ----------------------
	// Daily price and alert update.
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		err := cronDaily(app)
		if err != nil {
			log.Println("error in hooking cronDaily: ", err)
			return err
		}

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
