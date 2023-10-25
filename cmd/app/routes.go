package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func routeUpdateDailyETF(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/update-daily-etf", func(c echo.Context) error {
		days := c.QueryParam("days")
		if days == "" {
			return c.JSON(http.StatusOK, map[string]any{"message": "require days param"})
		}

		daysInt, err := strconv.Atoi(days)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"message": "cannot conver query to int"})
		}

		go updateDailyCollectionETF(app, daysInt)

		return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
		// })
	}, /* optional middlewares */ apis.RequireRecordAuth("users"))
}

func routeUpdateDaily(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/update-daily", func(c echo.Context) error {
		go updateDailyCollection(app)

		return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
	}, /* optional middlewares */ apis.RequireRecordAuth("users"))
}

func routeGetTrack(e *core.ServeEvent, app *pocketbase.PocketBase) {
	e.Router.GET("/track", func(c echo.Context) error {
		// 1. Get all records from Collection `track`: 'code, name, started'.
		var tempRecords = []recordTrack{}
		err := app.Dao().DB().
			Select("id", "code", "name", "started").
			From("track").
			All(&tempRecords)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{"message": "error in getting all records from database `track`"},
			)
		}
		if len(tempRecords) == 0 {
			return c.JSON(http.StatusOK, map[string]string{"message": "ok", "data": ""})
		}

		results := make([]struct {
			Code    string  `json:"code"`
			Name    string  `json:"name"`
			Started string  `json:"started"`
			Change  float64 `json:"change"`
		}, 0)
		// 2. For each record, get its closes since 'started'.
		for _, rec := range tempRecords {
			var tempDaily []recordDaily
			err = app.Dao().DB().
				Select("code", "date", "open", "high", "low", "close").
				From("daily").
				Where(dbx.NewExp(fmt.Sprintf("code = \"%s\"", rec.Code))).
				AndWhere(dbx.NewExp(fmt.Sprintf("date >= \"%s\"", rec.Started))).
				OrderBy("date ASC").
				All(&tempDaily)
			if err != nil {
				log.Println("error in reading daily collection: ", err)
				continue
			}

			// If len() == 0, tracked stock has no data yet; make it to
			// take yesterday's value (same as though 1D past).
			if len(tempDaily) == 0 {
				results = append(results, struct {
					Code    string  "json:\"code\""
					Name    string  "json:\"name\""
					Started string  "json:\"started\""
					Change  float64 "json:\"change\""
				}{
					Code:    rec.Code,
					Name:    rec.Name,
					Started: rec.Started,
					Change:  float64(0),
				})
				continue
			}

			closeStart := tempDaily[0].Close
			closeEnd := tempDaily[len(tempDaily)-1].Close
			change := (closeEnd - closeStart) / closeStart

			// log.Println(rec.Code, rec.Started)
			// log.Println(change)
			// log.Println(tempDaily)
			// log.Println("--------------------")

			results = append(results, struct {
				Code    string  "json:\"code\""
				Name    string  "json:\"name\""
				Started string  "json:\"started\""
				Change  float64 "json:\"change\""
			}{
				Code:    rec.Code,
				Name:    rec.Name,
				Started: rec.Started,
				Change:  change,
			})
		}

		return c.JSON(http.StatusOK, map[string]any{"message": "ok", "data": results})
	}, /* optional middlewares */ apis.RequireRecordAuth("users"))
}

// dateStarted, errParseTime := time.Parse(DefaultDateLayout, x.Started)
// if errParseTime != nil {
// 	log.Println(fmt.Errorf("error in parsing started time: %w", errParseTime))
// 	continue
// }
// daysPast := int(time.Since(dateStarted).Hours() / float64(HoursPerDay))
