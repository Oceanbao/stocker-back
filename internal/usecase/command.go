package usecase

import (
	"fmt"

	"example.com/stocker-back/internal/infra"
	apieastmoney "example.com/stocker-back/internal/infra/api_eastmoney"
	"example.com/stocker-back/internal/screener"
	"example.com/stocker-back/internal/stock"
	"github.com/samber/lo"
)

type Command struct {
	repoStock  stock.Repository
	repoScreen screener.Repository
	logger     infra.Logger
	notifier   infra.Notifier
}

func NewCommand(repoStock stock.Repository, repoScreen screener.Repository, logger infra.Logger, notifier infra.Notifier) *Command { //nolint:lll
	return &Command{
		repoStock:  repoStock,
		repoScreen: repoScreen,
		logger:     logger,
		notifier:   notifier,
	}
}

func (c *Command) UpdateStocks() error {
	c.logger.Debugf("UpdateStocks", "message", "start...")
	stocksAll, err := c.repoStock.GetStocks()
	if err != nil {
		return err
	}

	tickers := lo.Map(stocksAll, func(stock stock.Stock, _ int) string {
		return stock.Ticker
	})

	apiServiceEastmoney := apieastmoney.NewAPIServiceEastmoney(c.logger)
	stocksAllNew := apiServiceEastmoney.CrawlStocks(tickers)

	err = c.repoStock.SetStocks(stocksAllNew)
	if err != nil {
		c.logger.Errorf("SetStocksAll", "error", err.Error())
		c.notifier.Sendf("SetStocksAll", err.Error())
		return err
	}

	return nil
}

func (c *Command) UpdateDailyData() error {
	c.logger.Debugf("UpdateDailyData()", "message", "start...")
	dailyDataToCrawl, err := c.repoStock.GetDailyDataLastAll()
	if err != nil {
		return err
	}

	c.logger.Debugf("UpdateDailyData()", "message", "crawl...")
	apiServiceEastmoney := apieastmoney.NewAPIServiceEastmoney(c.logger)
	dailyDataNew := apiServiceEastmoney.CrawlDailyToDate(dailyDataToCrawl)

	if err = c.repoStock.SetDailyData(dailyDataNew); err != nil {
		c.logger.Errorf("SetDailyData()", "error", err.Error())
		c.notifier.Sendf("SetDailyData()", err.Error())
		return err
	}

	c.logger.Debugf("total crawled: [%d]", "len", len(dailyDataNew))
	c.notifier.Sendf("Stocker - total crawled", fmt.Sprintf("%d", len(dailyDataNew)))

	return nil
}

func (c *Command) UpdateDailyScreen() error {
	c.logger.Debugf("UpdateDailyScreen", "message", "start...")

	dailyDataLastAll, err := c.repoStock.GetDailyDataAll()
	if err != nil {
		return err
	}

	// for key, val := range dailyDataLastAll {
	// 	c.logger.Debugf("UpdateDailyScreen", "ticker", key, "daily", val)
	// 	break
	// }

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

func (c *Command) CreateStockAndDailyData(ticker string) error {
	// Crawl ticker stock.
	apiServiceEastmoney := apieastmoney.NewAPIServiceEastmoney(c.logger)
	stockNew := apiServiceEastmoney.CrawlStock(ticker)
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

func (c *Command) DeleteStockByTicker(ticker string) error {
	// NOTE: need to delete all related collections.
	if err := c.repoStock.DeleteStockByTicker(ticker); err != nil {
		return err
	}

	if err := c.repoStock.DeleteDailyDataByTicker(ticker); err != nil {
		return err
	}

	return nil
}
