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

func cronDailyPriceUpdate(app *pocketbase.PocketBase) { //nolint:funlen,gocognit // ignore for now
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
	var tempDaily []dailyRecord
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
	groupedDaily := lo.GroupBy(tempDaily, func(d dailyRecord) string {
		return d.Code
	})

	// 3.5 For each, compute RSI and KDJ and MACD.
	tempAlertUpsert := make([]struct {
		Code string
		Rsi  float64
		K    float64
		D    float64
		J    float64
		Diff float64
		Dea  float64
		Name string
		Cap  float64
	}, len(groupedDaily))
	tempCounter := 0
	for code, v := range groupedDaily {
		currIndex := tempCounter
		tempCounter++
		rsi, _ := RSI(lo.Map(v, func(d dailyRecord, _ int) float64 {
			return d.Close
		}))
		k, d, j := KDJ(
			lo.Map(v, func(d dailyRecord, _ int) float64 {
				return d.High
			}),
			lo.Map(v, func(d dailyRecord, _ int) float64 {
				return d.Low
			}),
			lo.Map(v, func(d dailyRecord, _ int) float64 {
				return d.Close
			}),
		)
		diff, dea := MACD(lo.Map(v, func(d dailyRecord, _ int) float64 {
			return d.Close
		}))
		// Get "name" and "cap" from `stocks`.
		record, errGetStockCode := app.Dao().FindFirstRecordByData("stocks", "code", code)
		if errGetStockCode != nil {
			log.Printf("error in finding record in 'stocks': (code: %v) (error: %v)\n", code, errGetStockCode)
			continue
		}
		tempAlertUpsert[currIndex] = struct {
			Code string
			Rsi  float64
			K    float64
			D    float64
			J    float64
			Diff float64
			Dea  float64
			Name string
			Cap  float64
		}{
			code,
			rsi[len(rsi)-1],
			k[len(k)-1],
			d[len(d)-1],
			j[len(j)-1],
			diff[len(diff)-1],
			dea[len(dea)-1],
			record.Get("name").(string),
			record.Get("cap").(float64),
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
				return fmt.Errorf("error in deleting record with ID: %s : %v", x.ID, err)
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

// func cronDailyTrackUpdate(app *pocketbase.PocketBase) {
// 	// Collection `track`: 'code, name, started, days, change'.
// 	// When front create new record, the latter 3 fields are set to zero value.
// 	// 1. Get all records.
// 	var tempRecords = []struct {
// 		ID      string `db:"id" json:"id"`
// 		Code    string `db:"code" json:"code"`
// 		Name    string `db:"name" json:"name"`
// 		Started string `db:"started" json:"started"`
// 		Days    string `db:"days" json:"days"`
// 		Change  string `db:"change" json:"change"`
// 	}{}
// 	err := app.Dao().DB().
// 		Select("id", "code", "name", "started", "days", "change").
// 		From("track").
// 		All(&tempRecords)
// 	if err != nil {
// 		log.Println("error in getting all records from database `track`")
// 		return
// 	}
// 	if len(tempRecords) == 0 {
// 		log.Println("no record in collection `track`, skip cron job.")
// 		return
// 	}

// 	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
// 		for _, x := range tempRecords {
// 			filter := fmt.Sprintf("code = '%s' && date > '%s'", x.Code, x.Started)
// 			recordsDaily, errFindDailyByDate := app.Dao().FindRecordsByFilter("daily", filter, "date", -1, 0)
// 			if errFindDailyByDate != nil {
// 				return fmt.Errorf("error in finding record in `daily`: %w", errFindDailyByDate)
// 			}

// 			dateStarted, errParseTime := time.Parse(DefaultDateLayout, x.Started)
// 			if errParseTime != nil {
// 				return fmt.Errorf("error in parsing started time: %w", errParseTime)
// 			}

// 			daysPast := int(time.Since(dateStarted).Hours() / float64(HoursPerDay))
// 			priceLatest := recordsDaily[len(recordsDaily)-1].Get("close").(float64) //nolint: errcheck //what
// 			priceStarted := recordsDaily[0].Get("close").(float64)                  //nolint: errcheck  //what
// 			change := (priceLatest - priceStarted) / priceStarted

// 			recordTrack, errFindTrackRecord := app.Dao().FindRecordById("track", x.ID)
// 			if errFindTrackRecord != nil {
// 				return fmt.Errorf("error in finding record in `track`: %w", errFindTrackRecord)
// 			}

// 			recordTrack.Set("days", daysPast)
// 			recordTrack.Set("change", change)

// 			if err = txDao.SaveRecord(recordTrack); err != nil {
// 				return fmt.Errorf("error in updating record: %v (%w)", x.ID, err)
// 			}
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		log.Println("error in transaction of updating stocks records: ", err)
// 		return
// 	}
// }
