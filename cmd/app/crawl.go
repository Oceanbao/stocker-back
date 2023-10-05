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

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/samber/lo"
)

type DailyData struct {
	Code   string   `json:"code"`
	Market int      `json:"market"`
	Klines []string `json:"klines"`
}

type DailyResponse struct {
	Data DailyData `json:"data"`
}

func updateDailyCollection(app *pocketbase.PocketBase) {
	var tempStocks = []struct {
		Code string `db:"code" json:"code"`
	}{}
	err := app.Dao().DB().
		Select("code").
		From("stocks").
		All(&tempStocks)
	if err != nil {
		log.Println("error in reading database `stocks`: ", err)
		return
	}

	resultsGood, resultsBad := crawlDaily(
		lo.Map(tempStocks, func(d struct {
			Code string `db:"code" json:"code"`
		}, _ int) string {
			return d.Code
		}),
	)

	// Write in `fail_daily` for log.
	if len(resultsBad) > 0 {
		collection, err := app.Dao().FindCollectionByNameOrId("fail_daily")
		if err != nil {
			log.Println("error in finding collection 'fail_daily': ", err)
			return
		}
		err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
			for _, x := range resultsBad {
				dataToEnter := map[string]any{
					"url": x,
				}
				record := models.NewRecord(collection)
				record.Load(dataToEnter)

				if err = txDao.SaveRecord(record); err != nil {
					log.Println("error in writing record: ", x, err)
				}
			}
			return nil
		})
		if err != nil {
			log.Println("error in transaction of writing fail_daily records: ", err)
			return
		}
	}

	// 3.3 Write valid responses to `daily` collection.
	collection, err := app.Dao().FindCollectionByNameOrId("daily")
	if err != nil {
		log.Println("error in finding collection 'daily': ", err)
		return
	}
	err = app.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, x := range resultsGood {
			for _, entry := range x.Data.Klines {
				parts := strings.Split(entry, ",")
				dataToEnter := map[string]any{
					"code":  fmt.Sprintf("%d.%s", x.Data.Market, x.Data.Code),
					"date":  parts[0],
					"open":  parts[1],
					"high":  parts[3],
					"low":   parts[4],
					"close": parts[2],
				}
				record := models.NewRecord(collection)
				record.Load(dataToEnter)

				if err = txDao.SaveRecord(record); err != nil {
					log.Println("error in writing record to `daily`: ", x.Data.Code, parts[0], err)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Println("error in transaction of writing new daily records: ", err)
		return
	}
}

func crawlDaily(codes []string) (resultGood []DailyResponse, resultBad []string) {
	numJobs := len(codes)
	// numJobs := 20
	urls := make(chan string, numJobs)
	results := make(chan string, numJobs)

	today := time.Now()
	pastXDays := 14
	fiveDaysAgo := today.AddDate(0, 0, -pastXDays).Format("20060102")

	numWorkers := 3
	for w := 1; w <= numWorkers; w++ {
		go requestWorker(w, urls, results)
	}
	for _, x := range codes[:numJobs] {
		url := fmt.Sprintf(
			"https://54.push2his.eastmoney.com/api/qt/stock/kline/get?"+
				"cb=jQuery35106707668456928451_1695010059469"+
				"&secid=%s"+
				"&ut=fa5fd1943c7b386f172d6893dbfba10b"+
				"&fields1=f1%%2Cf2%%2Cf3%%2Cf4%%2Cf5%%2Cf6"+
				"&fields2=f51%%2Cf52%%2Cf53%%2Cf54%%2Cf55%%2Cf56%%2Cf57%%2Cf58%%2Cf59%%2Cf60%%2Cf61"+
				"&klt=101&fqt=1"+
				"&beg=%s&end=20500101"+
				"&lmt=10&_=1695010059524", x, fiveDaysAgo,
		)
		urls <- url
	}
	close(urls)

	output := make([]string, 0)
	for idx := 1; idx <= numJobs; idx++ {
		output = append(output, <-results)
	}
	close(results)

	outputGood := make([]string, 0)
	outputBad := make([]string, 0)
	for _, v := range output {
		if !strings.HasPrefix(v, "fail") {
			outputGood = append(outputGood, v)
		} else {
			outputBad = append(outputBad, v)
		}
	}

	validResults := make([]DailyResponse, len(outputGood))
	for idx, x := range outputGood {
		if err := json.Unmarshal([]byte(x), &validResults[idx]); err != nil {
			log.Println("error in unmarshalling json: ", err)
			continue
		}
	}

	return validResults, outputBad
}

func requestWorker(id int, urls <-chan string, results chan<- string) {
	for url := range urls {
		sleepSecond := 3
		time.Sleep(time.Second * time.Duration(rand.Intn(sleepSecond)+1)) //nolint:gosec // no need
		log.Println("worker", id, "started")

		resp, err := http.Get(url) //nolint:gosec,noctx // just ignore
		if err != nil {
			results <- fmt.Sprintf("fail: %v (%v)", err, url)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			results <- fmt.Sprintf("fail: %v (%v)", err, url)
			continue
		}

		bodyText := string(body)
		meat := sliceStringByChar(bodyText, "(", ")")
		log.Println("worker", id, "done")
		results <- meat
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
