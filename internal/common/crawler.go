package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/stocker-back/internal/stock"
	"github.com/samber/lo"
)

type RawDailyCrawl struct {
	Data struct {
		Code   string   `json:"code"`
		Market int      `json:"market"`
		Klines []string `json:"klines"`
	} `json:"data"`
}

type FormedDailyCrawl struct {
	Ticker string
	Ohlc   []stock.OHLC
}

// CrawlDailyDataToDate concurrently crawls and produces DailyData for each up to date.
func CrawlDailyDataToDate(dailyDataToCrawl []stock.DailyData) []stock.DailyData {
	// numJobs := len(dailyDataToCrawl)
	numJobs := 3
	chanJobs := make(chan stock.DailyData, numJobs)
	chanResults := make(chan FormedDailyCrawl, numJobs)
	concurrency := 3
	secondThrottled := 3

	for range lo.Range(concurrency) {
		go func(in <-chan stock.DailyData, out chan<- FormedDailyCrawl) {
			for stockToCrawl := range in {
				time.Sleep(time.Second * time.Duration(secondThrottled))
				lastDate, _ := time.Parse(DateLayoutPocketbase, stockToCrawl.Date)
				startDate := lastDate.AddDate(0, 0, 1)
				rawDaily, err := crawlDailyByTicker(stockToCrawl.Ticker, startDate)
				if err != nil {
					out <- FormedDailyCrawl{
						Ticker: "",
						Ohlc:   nil,
					}
					continue
				}
				// No new klines for this stock and startDate.
				if len(rawDaily.Data.Klines) == 0 {
					out <- FormedDailyCrawl{
						Ticker: "",
						Ohlc:   nil,
					}
					continue
				}

				candles := rawDailyDataToOHLC(rawDaily)

				out <- FormedDailyCrawl{
					Ticker: stockToCrawl.Ticker,
					Ohlc:   candles,
				}
			}
		}(chanJobs, chanResults)
	}

	for _, job := range dailyDataToCrawl[:numJobs] {
		chanJobs <- job
	}
	close(chanJobs)

	var output []stock.DailyData
	for range numJobs {
		result := <-chanResults
		if result.Ticker != "" {
			for _, val := range result.Ohlc {
				dailyData := stock.DailyData{
					Ticker: result.Ticker,
					Date:   val.Date,
					Open:   val.Open,
					High:   val.High,
					Low:    val.Low,
					Close:  val.Close,
				}
				output = append(output, dailyData)
			}
		}
	}

	return output
}

// crawlDailyByTicker crawls the last days raw daily data for a given ticker.
func crawlDailyByTicker(ticker string, startDate time.Time) (RawDailyCrawl, error) {
	startDateFormated := startDate.Format(DateLayoutNewOriental)

	url := fmt.Sprintf(
		"https://54.push2his.eastmoney.com/api/qt/stock/kline/get?"+
			"cb=jQuery35106707668456928451_1695010059469"+
			"&secid=%s"+
			"&ut=fa5fd1943c7b386f172d6893dbfba10b"+
			"&fields1=f1%%2Cf2%%2Cf3%%2Cf4%%2Cf5%%2Cf6"+
			"&fields2=f51%%2Cf52%%2Cf53%%2Cf54%%2Cf55%%2Cf56%%2Cf57%%2Cf58%%2Cf59%%2Cf60%%2Cf61"+
			"&klt=101&fqt=1"+
			"&beg=%s&end=21000101"+
			"&lmt=1000&_=1695010059524", ticker, startDateFormated,
	)

	timeout := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	body, err := fetch(ctx, url)
	if err != nil {
		return RawDailyCrawl{}, err
	}

	text := sliceStringByChar(string(body), "(", ")")

	var output RawDailyCrawl
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		return RawDailyCrawl{}, err
	}

	return output, nil
}

func fetch(ctx context.Context, url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

// sliceStringByChar slice input string by startChar and endChar if they are valid.
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

func rawDailyDataToOHLC(rawDaily RawDailyCrawl) []stock.OHLC {
	ohlc := make([]stock.OHLC, len(rawDaily.Data.Klines))

	for idx, data := range rawDaily.Data.Klines {
		parts := strings.Split(data, ",")

		ohlc[idx].Date = parts[0]

		var num float64
		num, _ = strconv.ParseFloat(parts[1], 64)
		ohlc[idx].Open = num

		num, _ = strconv.ParseFloat(parts[2], 64)
		ohlc[idx].Close = num

		num, _ = strconv.ParseFloat(parts[3], 64)
		ohlc[idx].High = num

		num, _ = strconv.ParseFloat(parts[4], 64)
		ohlc[idx].Low = num
	}

	return ohlc
}
