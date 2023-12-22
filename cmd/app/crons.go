package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/samber/lo"
)

const DefaultDateLayout = "2006-01-02 15:04:05.000Z"
const HoursPerDay = 24

func cronDailyPriceUpdate(app *pocketbase.PocketBase) {
	// 0. Clear entire collection of `alert`.
	// ---------------------------------------------------
	err := clearCollection(app, "alert")
	if err != nil {
		log.Println("error in clearCollection of `alert`: ", err)
	}

	// 1. Update `daily` collection.
	updateDailyCollection(app)

	// 3.4 Get latest 180-day `daily` record and groupby code.
	xDaysAgo := time.Now().AddDate(0, 0, -180).Format("2006-01-02 15:04:05.000Z")
	var tempDaily []recordDaily
	err = app.Dao().DB().
		Select("code", "date", "open", "high", "low", "close").
		From("daily").
		Where(dbx.NewExp(fmt.Sprintf("date >= \"%s\"", xDaysAgo))).
		OrderBy("date ASC").
		All(&tempDaily)
	if err != nil {
		log.Println("error in reading daily collection: ", err)
		return
	}
	groupedDaily := lo.GroupBy(tempDaily, func(d recordDaily) string {
		return d.Code
	})

	// 3.5 For each, compute RSI and KDJ and MACD.
	tempAlertUpsert := make([]recordAlert, len(groupedDaily))
	tempCounter := 0
	for code, v := range groupedDaily {
		currIndex := tempCounter
		tempCounter++
		rsi := RSI(lo.Map(v, func(d recordDaily, _ int) float64 {
			return d.Close
		}))
		k, d, j := KDJ(
			lo.Map(v, func(d recordDaily, _ int) float64 {
				return d.High
			}),
			lo.Map(v, func(d recordDaily, _ int) float64 {
				return d.Low
			}),
			lo.Map(v, func(d recordDaily, _ int) float64 {
				return d.Close
			}),
		)
		diff, dea := MACD(lo.Map(v, func(d recordDaily, _ int) float64 {
			return d.Close
		}))
		// Get "name" and "cap" from `stocks`.
		record, errGetStockCode := app.Dao().FindFirstRecordByData("stocks", "code", code)
		if errGetStockCode != nil {
			log.Printf("error in finding record in 'stocks': (code: %v) (error: %v)\n", code, errGetStockCode)
			continue
		}
		tempAlertUpsert[currIndex] = recordAlert{
			Code: code,
			Name: record.Get("name").(string),
			Cap:  record.Get("cap").(float64),
			Rsi:  rsi[len(rsi)-1],
			K:    k[len(k)-1],
			D:    d[len(d)-1],
			J:    j[len(j)-1],
			Diff: diff[len(diff)-1],
			Dea:  dea[len(dea)-1],
		}
	}

	// 3.6 Check for target, insert into `alert` (code, rsi, name, cap).
	collection, err := app.Dao().FindCollectionByNameOrId("alert")
	if err != nil {
		log.Println("error in finding collection 'alert': ", err)
		return
	}
	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, x := range tempAlertUpsert {
			record := models.NewRecord(collection)
			record.Load(map[string]any{
				"code": x.Code,
				"name": x.Name,
				"cap":  x.Cap,
				"rsi":  x.Rsi,
				"k":    x.K,
				"d":    x.D,
				"j":    x.J,
				"diff": x.Diff,
				"dea":  x.Dea,
			})

			if err = txDao.SaveRecord(record); err != nil {
				log.Println("error in writing record to alert: ", x, err)
			}
		}

		return nil
	})
	if err != nil {
		log.Println("error in transaction of writing new alert records: ", err)
		return
	}
}

func clearCollection(app *pocketbase.PocketBase, collection string) error {
	var tempAlert = []struct {
		ID string `db:"id" json:"id"`
	}{}

	err := app.Dao().DB().
		Select("id").
		From(collection).
		All(&tempAlert)
	if err != nil {
		log.Println("error in reading records from collection: ", collection, err)
		return err
	}

	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, x := range tempAlert {
			record, errFindRecord := app.Dao().FindRecordById(collection, x.ID)
			if errFindRecord != nil {
				// Quite possible this is unlikely, so ignore.
				// return fmt.Errorf("error in finding record by ID: %s", x.ID)
				return errFindRecord
			}
			if err = txDao.DeleteRecord(record); err != nil {
				return fmt.Errorf("error in deleting record with ID: %s : %w", x.ID, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Println("error in transaction - clear collection: ", collection, err)
		return err
	}

	return nil
}

func cronDailySelectETFUpdate(app *pocketbase.PocketBase) {
	// 0. Clear entire collection of `alert`.
	// ---------------------------------------------------
	err := clearCollection(app, "alert_etf")
	if err != nil {
		log.Println("error in clearCollection of `alert_etf`: ", err)
	}

	// 1. Update `daily` collection.
	updatePeriod := 14
	updateDailyCollectionETF(app, updatePeriod)

	// 3.4 Get latest 180-day `daily` record and groupby code.
	// xDaysAgo := time.Now().AddDate(0, 0, -180).Format("2006-01-02 15:04:05.000Z")
	var tempDaily []recordDaily
	err = app.Dao().DB().
		Select("code", "date", "open", "high", "low", "close").
		From("daily_etf").
		// Where(dbx.NewExp(fmt.Sprintf("date >= \"%s\"", xDaysAgo))).
		OrderBy("date ASC").
		All(&tempDaily)
	if err != nil {
		log.Println("error in reading collection `daily_etf`: ", err)
		return
	}
	groupedDaily := lo.GroupBy(tempDaily, func(d recordDaily) string {
		return d.Code
	})

	// 3.5 For each, compute RSI and KDJ and MACD.
	tempAlertUpsert := make([]recordAlert, len(groupedDaily))
	tempCounter := 0
	for code, v := range groupedDaily {
		currIndex := tempCounter
		tempCounter++
		rsi := RSI(lo.Map(v, func(d recordDaily, _ int) float64 {
			return d.Close
		}))
		k, d, j := KDJ(
			lo.Map(v, func(d recordDaily, _ int) float64 {
				return d.High
			}),
			lo.Map(v, func(d recordDaily, _ int) float64 {
				return d.Low
			}),
			lo.Map(v, func(d recordDaily, _ int) float64 {
				return d.Close
			}),
		)
		diff, dea := MACD(lo.Map(v, func(d recordDaily, _ int) float64 {
			return d.Close
		}))
		// Get "name" and "cap" from `stocks`.
		record, errGetStockCode := app.Dao().FindFirstRecordByData("etf", "code", code)
		if errGetStockCode != nil {
			log.Printf("error in finding record in 'etf': (code: %v) (error: %v)\n", code, errGetStockCode)
			continue
		}
		tempAlertUpsert[currIndex] = recordAlert{
			Code: code,
			Name: record.Get("name").(string),
			Cap:  0,
			Rsi:  rsi[len(rsi)-1],
			K:    k[len(k)-1],
			D:    d[len(d)-1],
			J:    j[len(j)-1],
			Diff: diff[len(diff)-1],
			Dea:  dea[len(dea)-1],
		}
	}

	// 3.6 Check for target, insert into `alert` (code, rsi, name, cap).
	collection, err := app.Dao().FindCollectionByNameOrId("alert_etf")
	if err != nil {
		log.Println("error in finding collection 'alert_etf': ", err)
		return
	}
	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, x := range tempAlertUpsert {
			record := models.NewRecord(collection)
			record.Load(map[string]any{
				"code": x.Code,
				"name": x.Name,
				"rsi":  x.Rsi,
				"k":    x.K,
				"d":    x.D,
				"j":    x.J,
				"diff": x.Diff,
				"dea":  x.Dea,
			})

			if err = txDao.SaveRecord(record); err != nil {
				log.Println("error in writing record to `alert_etf`: ", x, err)
			}
		}

		return nil
	})
	if err != nil {
		log.Println("error in transaction of writing new `alert_etf` records: ", err)
		return
	}
}

func cronDailyTallyUpdate(app *pocketbase.PocketBase) { //nolint:gocognit //ignore
	// 1. Get all stock ID from `track`
	recordTracks := []recordTrack{}
	err := app.Dao().DB().
		Select("id", "code", "name", "started").
		From("track").
		All(&recordTracks)
	if err != nil {
		log.Println("error in getting all records from database `track`")
		return
	}

	recordTracksETF := []recordTrack{}
	err = app.Dao().DB().
		Select("id", "code", "name", "started").
		From("track_etf").
		All(&recordTracksETF)
	if err != nil {
		log.Println("error in getting all records from database `track_etf`")
		return
	}

	if len(recordTracks) == 0 && len(recordTracksETF) == 0 {
		log.Println("no record in collection `track` and `track_etf`, skip cron job.")
		return
	}

	allTracks := make([]struct {
		etf     bool
		records []recordTrack
	}, 2) //nolint:gomnd //ignore

	allTracks[0].etf = false
	allTracks[0].records = recordTracks
	allTracks[1].etf = true
	allTracks[1].records = recordTracksETF

	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, val := range allTracks {
			var collectionName string
			if val.etf {
				collectionName = "daily_etf"
			} else {
				collectionName = "daily"
			}

			for _, x := range val.records {
				filter := fmt.Sprintf("code = '%s' && date > '%s'", x.Code, x.Started)
				recordsDaily, err := app.Dao().FindRecordsByFilter(collectionName, filter, "date", -1, 0)
				if err != nil {
					log.Println(fmt.Errorf("error in finding record in `daily`: %w", err))
					continue
				}
				if len(recordsDaily) == 0 {
					log.Printf("skip [%v]: no day past since tracked\n", x.Code)
					continue
				}

				priceLatest, ok := recordsDaily[len(recordsDaily)-1].Get("close").(float64)
				if !ok {
					log.Printf("skip [%v]: fail to convert 'close' to float64\n", x.Code)
					continue
				}
				priceStarted, ok := recordsDaily[0].Get("close").(float64)
				if !ok {
					log.Printf("skip [%v]: fail to convert 'close' to float64\n", x.Code)
					continue
				}
				change := (priceLatest - priceStarted) / priceStarted

				recordUser, err := app.Dao().FindAuthRecordByUsername("users", "oceanbao")
				if err != nil {
					return fmt.Errorf("error in finding record in `users`: %w", err)
				}

				tally, ok := recordUser.Get("trade_tally").(float64)
				if !ok {
					log.Printf("skip [%v]: fail to convert 'tally' to float64\n", x.Code)
					continue
				}
				tally += change
				recordUser.Set("trade_tally", tally)

				if err = txDao.SaveRecord(recordUser); err != nil {
					return fmt.Errorf("error in updating record `users`: %v (%w)", recordUser.Id, err)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Println("error in transaction of updating `users` records (tally): ", err)
		return
	}
}
