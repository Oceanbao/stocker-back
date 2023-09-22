package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/samber/lo"
)

func cronDaily(app *pocketbase.PocketBase) error {
	scheduler := cron.New()

	// Every week Mon-Fri at 00:00
	// err := scheduler.Add("daily", "*/1 * * * *", func() {
	err := scheduler.Add("daily", "0 0 * * *", func() {
		dailyCronJob(app)
	})

	if err != nil {
		log.Println("ERR: cron job `daily`")
		return err
	}

	scheduler.Start()

	return nil
}

func dailyCronJob(app *pocketbase.PocketBase) {
	// 0. Clear entire collection of `alert`.
	// ---------------------------------------------------
	err := clearCollection(app, "alert")
	if err != nil {
		log.Println("error in clearCollection of: alert")
	}

	// 1. Get all stock code from `stocks` collection.
	// ---------------------------------------------------
	var tempStocks = []struct {
		Code string `db:"code" json:"code"`
	}{}
	err = app.Dao().DB().
		Select("code").
		From("stocks").
		All(&tempStocks)
	if err != nil {
		log.Println("error in reading database `stocks`")
		return
	}

	// 3. Concurrently request all stocks for latest daily.
	// ---------------------------------------------------
	// 3.1 Build all urls into chan urls.
	numJobs := len(tempStocks)
	urls := make(chan string, numJobs)
	results := make(chan string, numJobs)

	today := time.Now()
	pastXDays := 3
	fiveDaysAgo := today.AddDate(0, 0, -pastXDays).Format("20060102")

	// 3.2 Fire off workers to request daily.
	numWorkers := 3
	for w := 1; w <= numWorkers; w++ {
		go requestWorker(w, urls, results)
	}
	for _, code := range tempStocks {
		urls <- fmt.Sprintf(
			"https://54.push2his.eastmoney.com/api/qt/stock/kline/get?"+
				"cb=jQuery35106707668456928451_1695010059469"+
				"&secid=%s"+
				"&ut=fa5fd1943c7b386f172d6893dbfba10b"+
				"&fields1=f1%%2Cf2%%2Cf3%%2Cf4%%2Cf5%%2Cf6"+
				"&fields2=f51%%2Cf52%%2Cf53%%2Cf54%%2Cf55%%2Cf56%%2Cf57%%2Cf58%%2Cf59%%2Cf60%%2Cf61"+
				"&klt=101&fqt=1"+
				"&beg=%s&end=20500101"+
				"&lmt=10&_=1695010059524", code, fiveDaysAgo,
		)
	}
	close(urls)

	var output []string
	for idx := 1; idx <= numJobs; idx++ {
		output = append(output, <-results)
	}
	close(results)

	// 3.3 Check results and write to `daily` collection.
	validResults := make([]struct {
		Data struct {
			Code   string
			Market int
			Klines []string
		}
	}, len(output))
	for idx, x := range output {
		if err = json.Unmarshal([]byte(x), &validResults[idx]); err != nil {
			log.Println("error in unmarshalling json: ", err)
			return
		}
	}
	collection, err := app.Dao().FindCollectionByNameOrId("daily")
	if err != nil {
		log.Println("error in finding collection 'daily': ", err)
		return
	}
	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, x := range validResults {
			for _, entry := range x.Data.Klines {
				parts := strings.Split(entry, ",")
				record := models.NewRecord(collection)
				record.Load(map[string]any{
					"code":  fmt.Sprintf("%d.%s", x.Data.Market, x.Data.Code),
					"date":  parts[0],
					"open":  parts[1],
					"high":  parts[3],
					"low":   parts[4],
					"close": parts[2],
				})

				if err = txDao.SaveRecord(record); err != nil {
					log.Println("error in writing record: ", x.Data.Code, parts[0], err)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Println("error in transaction of writing new daily records: ", err)
		return
	}

	// 3.4 Get latest 180-day `daily` record and groupby code.
	xDaysAgo := today.AddDate(0, 0, -180).Format("2006-01-02 15:04:05.000Z")
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

	// 3.5 For each, compute RSI and KDJ.
	tempAlertUpsert := make([]struct {
		Code string
		Rsi  float64
		Name string
		Cap  float64
	}, len(groupedDaily))
	tempCounter := 0
	for k, v := range groupedDaily {
		rsi, _ := RSI(lo.Map(v, func(d dailyRecord, _ int) float64 {
			return d.Close
		}))
		// Get "name" and "cap" from `stocks`.
		record, err := app.Dao().FindFirstRecordByData("stocks", "code", k)
		if err != nil {
			log.Println("error in finding record in 'stocks': ", err)
			return
		}
		tempAlertUpsert[tempCounter] = struct {
			Code string
			Rsi  float64
			Name string
			Cap  float64
		}{
			k,
			rsi[len(rsi)-1],
			record.Get("name").(string),
			record.Get("cap").(float64),
		}
		tempCounter++
	}

	// 3.6 Check for target, insert into `alert` (code, rsi, name, cap).
	collection, err = app.Dao().FindCollectionByNameOrId("alert")
	if err != nil {
		log.Println("error in finding collection 'alert': ", err)
		return
	}
	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, x := range tempAlertUpsert {
			record := models.NewRecord(collection)
			record.Load(map[string]any{
				"code": x.Code,
				"rsi":  x.Rsi,
				"name": x.Name,
				"cap":  x.Cap,
			})

			if err = txDao.SaveRecord(record); err != nil {
				log.Println("error in writing record to alert: ", x.Code)
			}
		}

		return nil
	})
	if err != nil {
		log.Println("error in transaction of writing new alert records: ", err)
		return
	}
}

func requestWorker(id int, urls <-chan string, results chan<- string) {
	for url := range urls {
		time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
		log.Println("worker", id, "started")
		resp, err := http.Get(url) //nolint:gosec,noctx // just ignore
		if err != nil {
			results <- url
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- url
			continue
		}
		bodyText := string(body)
		meat := sliceStringByChar(bodyText, "(", ")")
		log.Println("worker", id, "finished job")
		results <- meat

		resp.Body.Close()
	}
}

func sliceStringByChar(input, startChar, endChar string) string {
	startIndex := strings.Index(input, startChar)
	if startIndex == -1 {
		return ""
	}

	endIndex := strings.LastIndex(input, endChar)
	if endIndex == -1 {
		return ""
	}

	return input[startIndex+1 : endIndex]
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
		log.Println("error in reading records from collection: ", collection)
		return err
	}

	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, x := range tempAlert {
			record, errFindRecord := app.Dao().FindRecordById(collection, x.ID)
			if errFindRecord != nil {
				// Quite possible this is unlikely, so ignore.
				return fmt.Errorf("error in finding record by ID: %s", x.ID)
			}
			if err = txDao.DeleteRecord(record); err != nil {
				return fmt.Errorf("error in deleting record with ID: %s", x.ID)
			}
		}
		return nil
	})
	if err != nil {
		log.Println("error in transaction - clear collection: ", collection)
		return err
	}

	return nil
}
