package apieastmoney

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"example.com/stocker-back/internal/common"
	"example.com/stocker-back/internal/stock"
	"github.com/samber/lo"
)

type RawStockCrawl struct {
	Data struct {
		EPS                float64 `json:"f55"`
		Ticker             string  `json:"f57"`
		Name               string  `json:"f58"`
		Market             int     `json:"f107"`
		TotalShare         float64 `json:"f84"`
		TotalShareOut      float64 `json:"f85"`
		NetAssetPerShare   float64 `json:"f92"`
		NetProfit          float64 `json:"f105"`
		TotalCap           float64 `json:"f116"`
		TradeCap           float64 `json:"f117"`
		PricePerEarning    float64 `json:"f162"`
		PricePerBook       float64 `json:"f167"`
		ROE                float64 `json:"f173"`
		TotalRevenue       float64 `json:"f183"`
		TotalRevenueChange float64 `json:"f184"`
		NetProfitChange    float64 `json:"f185"`
		GrossProfit        float64 `json:"f186"`
		ProfitMargin       float64 `json:"f187"`
		DebtRatio          float64 `json:"f188"`
		DateOfPublic       int     `json:"f189"`
		UndistProfit       float64 `json:"f190"`
	} `json:"data"`
}

func (raw *RawStockCrawl) ToModel() stock.Stock {
	dateOfPublic, _ := time.Parse(common.DateLayoutNewOriental, fmt.Sprintf("%v", raw.Data.DateOfPublic))

	// NOTE: fix possible overflow in big float64 and esp. NetAsset by log2 value.
	return stock.Stock{
		Ticker:             fmt.Sprintf("%v.%v", raw.Data.Market, raw.Data.Ticker),
		Name:               raw.Data.Name,
		DateOfPublic:       dateOfPublic.Format(common.DateLayoutPocketbase),
		ETF:                strings.Contains(raw.Data.Name, "ETF"),
		EPS:                raw.Data.EPS,
		UndistProfit:       raw.Data.UndistProfit,
		TotalShare:         raw.Data.TotalShare,
		TotalShareOut:      raw.Data.TotalShareOut,
		TotalCap:           math.Log2(raw.Data.TotalCap),
		TradeCap:           math.Log2(raw.Data.TradeCap),
		NetAsset:           math.Log2(raw.Data.TotalCap * raw.Data.NetAssetPerShare),
		NetAssetPerShare:   raw.Data.NetAssetPerShare,
		NetProfit:          raw.Data.NetProfit,
		NetProfitChange:    raw.Data.NetProfitChange,
		ProfitMargin:       raw.Data.ProfitMargin,
		PricePerEarning:    raw.Data.PricePerEarning,
		PricePerBook:       raw.Data.PricePerBook,
		ROE:                raw.Data.ROE,
		TotalRevenue:       raw.Data.TotalRevenue,
		TotalRevenueChange: raw.Data.TotalRevenueChange,
		GrossProfitMargin:  raw.Data.GrossProfit,
		DebtRatio:          raw.Data.DebtRatio,
	}
}

// CrawlStocks concurrently crawls and produces stock.Stock given tickers.
func (s *APIServiceEastmoney) CrawlStock(ticker string) (stock.Stock, error) {
	rawStock, err := crawlStock(ticker)
	if err != nil {
		s.logger.Debugf("CRAWL", "failed", ticker, "error", err.Error())
		return stock.NewEmptyStock(), err
	}

	if rawStock.Data.Name == "" {
		err := errors.New("ticker does not exists")
		s.logger.Debugf("CRAWL", "failed", ticker, "error", err.Error())
		return stock.NewEmptyStock(), err
	}

	s.logger.Debugf("CRAWL done", "ticker", ticker)

	return rawStock.ToModel(), nil
}

// CrawlStocks concurrently crawls and produces stock.Stock given tickers.
func (s *APIServiceEastmoney) CrawlStocks(tickers []string) []stock.Stock {
	numJobs := len(tickers)
	chanJobs := make(chan string, numJobs)
	chanResults := make(chan stock.Stock, numJobs)
	concurrency := 3
	secondThrottled := 3

	for range lo.Range(concurrency) {
		go func() {
			for ticker := range chanJobs {
				time.Sleep(time.Second * time.Duration(secondThrottled))
				s.logger.Debugf("CRAWL crawling", "ticker", ticker)
				rawStock, err := crawlStock(ticker)
				if err != nil {
					chanResults <- stock.NewEmptyStock()
					s.logger.Debugf("CRAWL", "failed", ticker, "error", err.Error())
					continue
				}

				s.logger.Debugf("CRAWL", "ok", ticker, "data", rawStock)

				chanResults <- rawStock.ToModel()
			}
		}()
	}

	for _, job := range tickers[:numJobs] {
		chanJobs <- job
	}
	close(chanJobs)

	var output []stock.Stock
	for range numJobs {
		stock := <-chanResults
		if stock.Ticker != "" {
			s.logger.Debugf("CRAWL", "stock", stock)
			output = append(output, stock)
		}
	}

	return output
}

// crawlStock runs the implementation.
func crawlStock(ticker string) (RawStockCrawl, error) {
	url := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/stock/get?"+
			"invt=2&fltt=1&cb=jQuery35105571137681219451_1708499614785"+
			"&fields=f57%%2Cf58%%2Cf107%%2Cf162%%2Cf152%%2Cf167%%2Cf92%%2Cf59%%2Cf183%%2Cf184%%2Cf105%%2Cf185%%2Cf186%%2Cf187%%2Cf173%%2Cf188%%2Cf84%%2Cf116%%2Cf85%%2Cf117%%2Cf190%%2Cf189%%2Cf62%%2Cf55"+ //nolint:lll
			"&secid=%s"+
			"&ut=fa5fd1943c7b386f172d6893dbfba10b&wbp2u=%%7C0%%7C0%%7C0%%7Cweb&_=1708499614786", ticker,
	)

	timeout := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	body, err := common.Fetch(ctx, url)
	if err != nil {
		return RawStockCrawl{}, err
	}

	text := sliceStringByChar(string(body), "(", ")")

	// NOTE: need to check if `stock` actually exsits in returned text.
	var output RawStockCrawl
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		return RawStockCrawl{}, err
	}

	return output, nil
}
