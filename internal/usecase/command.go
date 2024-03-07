package usecase

import (
	"fmt"

	"example.com/stocker-back/internal/infra"
	apieastmoney "example.com/stocker-back/internal/infra/api_eastmoney"
	"example.com/stocker-back/internal/screener"
	"example.com/stocker-back/internal/stock"
	"example.com/stocker-back/internal/tracking"
	"github.com/samber/lo"
)

type Command struct {
	repoStock    stock.Repository
	repoScreen   screener.Repository
	repoTracking tracking.Repository
	logger       infra.Logger
	notifier     infra.Notifier
}

func NewCommand(repoStock stock.Repository, repoScreen screener.Repository, repoTracking tracking.Repository, logger infra.Logger, notifier infra.Notifier) *Command { //nolint:lll
	return &Command{
		repoStock:    repoStock,
		repoScreen:   repoScreen,
		repoTracking: repoTracking,
		logger:       logger,
		notifier:     notifier,
	}
}

func (c *Command) UpdateStocks() error {
	c.logger.Infof("UpdateStocks - starting...")
	stocksAll, err := c.repoStock.GetStocks()
	if err != nil {
		return err
	}

	tickers := lo.Map(stocksAll, func(stock stock.Stock, _ int) string {
		return stock.Ticker
	})

	apiServiceEastmoney := apieastmoney.NewAPIServiceEastmoney(c.logger)
	c.logger.Infof("UpdateStocks - starting crawling...")

	failedTickers := make([]string, 0, len(tickers))
	for _, ticker := range tickers[:3] {
		stockNew, err := apiServiceEastmoney.CrawlStock(ticker)
		if err != nil {
			c.logger.Errorf("CrawlStock", "error", err.Error(), "ticker", ticker)
			failedTickers = append(failedTickers, ticker)
			continue
		}

		err = c.repoStock.SetStock(stockNew)
		if err != nil {
			c.logger.Errorf("SetStock", "error", err.Error(), "ticker", ticker)
			failedTickers = append(failedTickers, ticker)
			continue
		}
	}

	c.logger.Infof("UpdateStocks - DONE", "failed stocks", len(failedTickers), "tickers", failedTickers)
	c.notifier.Sendf(
		"UpdateStocks DONE",
		fmt.Sprintf("failed stocks len: %d tickers: %v", len(failedTickers), failedTickers),
	)

	return nil
}

func (c *Command) UpdateDailyData() error {
	c.logger.Infof("UpdateDailyData()", "message", "start...")
	dailyDataToCrawl, err := c.repoStock.GetDailyDataLastAll()
	if err != nil {
		return err
	}

	c.logger.Infof("UpdateDailyData()", "message", "crawl...")
	apiServiceEastmoney := apieastmoney.NewAPIServiceEastmoney(c.logger)
	dailyDataNew := apiServiceEastmoney.CrawlDailyToDate(dailyDataToCrawl)

	if err = c.repoStock.SetDailyData(dailyDataNew); err != nil {
		c.logger.Errorf("SetDailyData()", "error", err.Error())
		c.notifier.Sendf("SetDailyData()", err.Error())
		return err
	}

	c.logger.Infof("total crawled: [%d]", "len", len(dailyDataNew))
	c.notifier.Sendf("Stocker - total crawled", fmt.Sprintf("%d", len(dailyDataNew)))

	return nil
}

func (c *Command) UpdateDailyScreen() error {
	c.logger.Infof("UpdateDailyScreen", "message", "start...")

	dailyDataLastAll, err := c.repoStock.GetDailyDataAll()
	if err != nil {
		return err
	}

	screens := make([]screener.Screen, 0, len(dailyDataLastAll))
	for key, data := range dailyDataLastAll {
		candles := make([]stock.OHLC, len(data))
		for i, c := range data {
			candles[i].Date = c.Date
			candles[i].Open = c.Open
			candles[i].High = c.High
			candles[i].Low = c.Low
			candles[i].Close = c.Close
		}

		kdj := stock.ComputeKDJ(candles)

		lastJ := kdj[len(kdj)-1].J

		screens = append(screens, screener.Screen{
			Ticker: key,
			Kdj:    lastJ,
		})
	}

	err = c.repoScreen.SetScreens(screens)
	if err != nil {
		return err
	}

	return nil
}

func (c *Command) CreateStock(ticker string) error {
	// Crawl ticker stock.
	apiServiceEastmoney := apieastmoney.NewAPIServiceEastmoney(c.logger)
	stockNew, err := apiServiceEastmoney.CrawlStock(ticker)
	if err != nil {
		return err
	}

	// Write db
	if err := c.repoStock.SetStock(stockNew); err != nil {
		return err
	}

	c.logger.Infof("CreateStock", "ok", ticker)

	return nil
}

func (c *Command) CreateStockAndDailyData(ticker string) error {
	// Crawl ticker stock.
	apiServiceEastmoney := apieastmoney.NewAPIServiceEastmoney(c.logger)
	stockNew, err := apiServiceEastmoney.CrawlStock(ticker)
	if err != nil {
		return err
	}
	// Write db
	if err := c.repoStock.SetStock(stockNew); err != nil {
		return err
	}

	// Crawl ticker dailydata.
	dailyData := apiServiceEastmoney.CrawlDailyOne(ticker, 200) //nolint:gomnd // ignore
	// Write db
	if err := c.repoStock.SetDailyData(dailyData); err != nil {
		return err
	}

	return nil
}

// DeleteStockByTicker handles stock deletion by ticker, from all collections concerned.
func (c *Command) DeleteStockByTicker(ticker string) error {
	if err := c.repoStock.DeleteStockByTicker(ticker); err != nil {
		return err
	}

	return nil
}

// CreateTracking creates tracking entry for ticker.
func (c *Command) CreateTracking(ticker string) error {
	stock, err := c.repoStock.GetStockByTicker(ticker)
	if err != nil {
		return err
	}

	tracking := tracking.Tracking{
		Ticker: stock.Ticker,
		Name:   stock.Name,
	}
	if err = c.repoTracking.SetTracking(tracking); err != nil {
		return err
	}

	return nil
}

// DeleteTracking deletes tracking entry from collection.
func (c *Command) DeleteTracking(ticker string) error {
	if err := c.repoTracking.DeleteTracking(ticker); err != nil {
		return err
	}

	return nil
}
