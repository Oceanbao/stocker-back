package apieastmoney

import (
	"bytes"
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

func (raw *RawStockCrawl) ToModel(rawRank RawRankCrawl) (stock.Stock, error) {
	dateOfPublic, _ := time.Parse(common.DateLayoutNewOriental, fmt.Sprintf("%v", raw.Data.DateOfPublic))

	rankTotalCap, ok := rawRank.Data.Diff[0]["f1020"].(float64)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `rankTotalCap` as float64")
		return stock.NewEmptyStock(), err
	}

	rankNetAsset, ok := rawRank.Data.Diff[0]["f1135"].(float64)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `rankNetAsset` as float64")
		return stock.NewEmptyStock(), err
	}

	rankNetProfit, ok := rawRank.Data.Diff[0]["f1045"].(float64)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `rankNetProfit` as float64")
		return stock.NewEmptyStock(), err
	}

	rankGrossMargin, ok := rawRank.Data.Diff[0]["f1049"].(float64)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `rankGrossMargin` as float64")
		return stock.NewEmptyStock(), err
	}

	rankPER, ok := rawRank.Data.Diff[0]["f1009"].(float64)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `rankPER` as float64")
		return stock.NewEmptyStock(), err
	}

	rankPBR, ok := rawRank.Data.Diff[0]["f1023"].(float64)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `rankPBR` as float64")
		return stock.NewEmptyStock(), err
	}

	rankNetMargin, ok := rawRank.Data.Diff[0]["f1129"].(float64)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `rankNetMargin` as float64")
		return stock.NewEmptyStock(), err
	}

	rankROE, ok := rawRank.Data.Diff[0]["f1037"].(float64)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `rankROE` as float64")
		return stock.NewEmptyStock(), err
	}

	sector, ok := rawRank.Data.Diff[1]["f14"].(string)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `sector` as string")
		return stock.NewEmptyStock(), err
	}

	sectorTotal, ok := rawRank.Data.Diff[1]["f134"].(float64)
	if !ok {
		err := fmt.Errorf("error: failed to parse rank data `sectorTotal` as float64")
		return stock.NewEmptyStock(), err
	}

	// DELE: fix possible overflow in big float64 and esp. NetAsset by log2 value.
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
		RankTotalCap:       int(rankTotalCap),
		RankNetAsset:       int(rankNetAsset),
		RankNetProfit:      int(rankNetProfit),
		RankGrossMargin:    int(rankGrossMargin),
		RankPER:            int(rankPER),
		RankPBR:            int(rankPBR),
		RankNetMargin:      int(rankNetMargin),
		RankROE:            int(rankROE),
		Sector:             sector,
		SectorTotal:        int(sectorTotal),
	}, nil
}

type RawRankCrawl struct {
	Data struct {
		Diff []map[string]interface{} `json:"diff"`
	} `json:"data"`
}

// ValidateStockByTicker checks from API if ticker exists.
func ValidateStockByTicker(ticker string) (map[string]any, error) {
	rawStock, err := crawlStock(ticker)
	if err != nil {
		return nil, err
	}

	if rawStock.Data.Name == "" {
		return make(map[string]any), nil
	}

	var m map[string]any
	b, err := json.Marshal(rawStock)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return m, nil
}

// CrawlStocks concurrently crawls and produces stock.Stock given tickers.
func (s *APIServiceEastmoney) CrawlStock(ticker string) (stock.Stock, error) {
	rawStock, err := crawlStock(ticker)
	if err != nil {
		s.logger.Errorf("CRAWL", "failed", ticker, "error", err.Error())
		return stock.NewEmptyStock(), err
	}

	// DELE
	if rawStock.Data.Name == "" {
		err := errors.New("ticker does not exists")
		s.logger.Errorf("CRAWL", "failed", ticker, "error", err.Error())
		return stock.NewEmptyStock(), err
	}

	rawRank, err := crawlRank(ticker)
	if err != nil {
		s.logger.Errorf("CRAWL", "failed", ticker, "error", err.Error())
		return stock.NewEmptyStock(), err
	}

	s.logger.Infof("CRAWL done", "ticker", ticker)

	model, err := rawStock.ToModel(rawRank)
	if err != nil {
		s.logger.Errorf("CRAWL", "failed", ticker, "error", err.Error())
		return stock.NewEmptyStock(), err
	}

	return model, nil
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
				s.logger.Infof("CrawlStocks", "ticker", ticker)
				rawStock, err := crawlStock(ticker)
				if err != nil {
					chanResults <- stock.NewEmptyStock()
					s.logger.Errorf("CrawlStocks", "failed", ticker, "error", err.Error())
					continue
				}

				// DELE
				rawRank, err := crawlRank(ticker)
				if err != nil {
					chanResults <- stock.NewEmptyStock()
					s.logger.Errorf("CrawlStocks", "failed", ticker, "error", err.Error())
					continue
				}

				s.logger.Infof("CrawlStocks", "ok", ticker, "data", rawStock)

				model, err := rawStock.ToModel(rawRank)
				if err != nil {
					s.logger.Errorf("CrawlStocks", "failed", ticker, "error", err.Error())
					chanResults <- stock.NewEmptyStock()
					continue
				}

				chanResults <- model
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
			output = append(output, stock)
		}
	}

	return output
}

// crawlStock crawls the Eastmoney endpoint for stock meta.
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

	// DELE: need to check if `stock` actually exsits in returned text.
	var output RawStockCrawl
	// Handle NaN in JSON.
	b := bytes.ReplaceAll([]byte(text), []byte(":NaN"), []byte(":null"))
	if err := json.Unmarshal(b, &output); err != nil {
		return RawStockCrawl{}, err
	}

	return output, nil
}

// crawlRank crawls the Eastmoney endpoint for stock sector rank meta.
func crawlRank(ticker string) (RawRankCrawl, error) {
	url := fmt.Sprintf(
		"https://push2.eastmoney.com/api/qt/slist/get?fltt=1&"+
			"invt=2&cb=jQuery35105571137681219451_1708499614794&"+
			"fields=f12%%2Cf13%%2Cf14%%2Cf20%%2Cf58%%2Cf45%%2Cf132%%2Cf9%%2Cf152%%2Cf23%%2Cf49%%2Cf131%%2Cf137%%2Cf133%%2Cf134%%2Cf135%%2Cf129%%2Cf37%%2Cf1000%%2Cf3000%%2Cf2000&"+ //nolint:lll
			"secid=%s"+
			"&ut=fa5fd1943c7b386f172d6893dbfba10b&pn=1&np=1&spt=1&wbp2u=%%7C0%%7C0%%7C0%%7Cweb&_=1708499614795", ticker,
	)

	timeout := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	body, err := common.Fetch(ctx, url)
	if err != nil {
		return RawRankCrawl{}, err
	}

	text := sliceStringByChar(string(body), "(", ")")

	// DELE: need to check if `stock` actually exsits in returned text.
	var output RawRankCrawl
	// Handle NaN in JSON.
	b := bytes.ReplaceAll([]byte(text), []byte(":NaN"), []byte(":null"))
	if err := json.Unmarshal(b, &output); err != nil {
		return RawRankCrawl{}, err
	}

	return output, nil
}
