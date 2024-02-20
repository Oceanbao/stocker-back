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

func (raw *RawDailyCrawl) ToDailyData() FormedDailyCrawl {
	ticker := fmt.Sprintf("%v.%v", raw.Data.Market, raw.Data.Code)
	dailyDataAll := make([]stock.DailyData, 0, len(raw.Data.Klines))

	for _, data := range raw.Data.Klines {
		parts := strings.Split(data, ",")

		var dailyData stock.DailyData

		dailyData.Ticker = ticker

		dailyData.Date = parts[0]

		var num float64
		// Open
		num, _ = strconv.ParseFloat(parts[1], 64)
		dailyData.Open = num
		// Close
		num, _ = strconv.ParseFloat(parts[2], 64)
		dailyData.Close = num
		// High
		num, _ = strconv.ParseFloat(parts[3], 64)
		dailyData.High = num
		// Low
		num, _ = strconv.ParseFloat(parts[4], 64)
		dailyData.Low = num
		// Volume
		num, _ = strconv.ParseFloat(parts[5], 64)
		dailyData.Volume = num
		// Value
		num, _ = strconv.ParseFloat(parts[6], 64)
		dailyData.Value = num
		// Volatility
		num, _ = strconv.ParseFloat(parts[7], 64)
		dailyData.Volatility = num
		// % Change
		num, _ = strconv.ParseFloat(parts[8], 64)
		dailyData.Pchange = num
		// Change
		num, _ = strconv.ParseFloat(parts[9], 64)
		dailyData.Change = num
		// % Turnover
		num, _ = strconv.ParseFloat(parts[10], 64)
		dailyData.Turnover = num

		dailyDataAll = append(dailyDataAll, dailyData)
	}
	return FormedDailyCrawl{
		Ticker:    ticker,
		DailyData: dailyDataAll,
	}
}

type FormedDailyCrawl struct {
	Ticker    string
	DailyData []stock.DailyData
}

type CrawlService struct {
	logger Logger
}

func NewCrawlService(logger Logger) *CrawlService {
	return &CrawlService{
		logger: logger,
	}
}

// CrawlDailyDataToDate concurrently crawls and produces DailyData for each up to date.
func (s *CrawlService) CrawlDailyDataToDate(dailyDataToCrawl []stock.DailyData) []stock.DailyData {
	numJobs := len(dailyDataToCrawl)
	chanJobs := make(chan stock.DailyData, numJobs)
	chanResults := make(chan FormedDailyCrawl, numJobs)
	concurrency := 3
	secondThrottled := 3

	for range lo.Range(concurrency) {
		go func() {
			for stockToCrawl := range chanJobs {
				time.Sleep(time.Second * time.Duration(secondThrottled))
				lastDate, _ := time.Parse(DateLayoutPocketbase, stockToCrawl.Date)
				startDate := lastDate.AddDate(0, 0, 1)
				rawDaily, err := crawlDailyByTicker(stockToCrawl.Ticker, startDate)
				if err != nil {
					chanResults <- FormedDailyCrawl{
						Ticker:    "",
						DailyData: nil,
					}
					s.logger.Debugf("CRAWL", "fail", stockToCrawl.Ticker, "error", err.Error())
					continue
				}
				// No new klines for this stock and startDate.
				if len(rawDaily.Data.Klines) == 0 {
					chanResults <- FormedDailyCrawl{
						Ticker:    "",
						DailyData: nil,
					}
					s.logger.Debugf("CRAWL", "nil", stockToCrawl.Ticker, "message", "no new daily")
					continue
				}

				s.logger.Debugf("CRAWL", "ok", stockToCrawl.Ticker, "len of daily", len(rawDaily.Data.Klines))

				chanResults <- rawDaily.ToDailyData()
			}
		}()
	}

	for _, job := range dailyDataToCrawl[:numJobs] {
		chanJobs <- job
	}
	close(chanJobs)

	var output []stock.DailyData
	for range numJobs {
		formedDailyCrawl := <-chanResults
		if formedDailyCrawl.Ticker != "" {
			output = append(output, formedDailyCrawl.DailyData...)
		}
	}

	return output
}

// crawlDailyByTicker crawls the last days raw daily data for a given ticker.
func crawlDailyByTicker(ticker string, startDate time.Time) (RawDailyCrawl, error) {
	startDateFormated := startDate.Format(DateLayoutNewOriental)
	url := fmt.Sprintf(
		"https://push2his.eastmoney.com/api/qt/stock/kline/get?"+
			"cb=jQuery35104990802373722225_1708415137417"+
			"&secid=%s"+
			"&ut=fa5fd1943c7b386f172d6893dbfba10b"+
			"&fields1=f1%%2Cf2%%2Cf3%%2Cf4%%2Cf5%%2Cf6"+
			"&fields2=f51%%2Cf52%%2Cf53%%2Cf54%%2Cf55%%2Cf56%%2Cf57%%2Cf58%%2Cf59%%2Cf60%%2Cf61"+
			"&klt=101&fqt=1"+
			"&beg=%s&end=21000101", ticker, startDateFormated,
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
