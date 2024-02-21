package usecase

import (
	"fmt"

	"example.com/stocker-back/internal/common"
	"example.com/stocker-back/internal/screener"
	"example.com/stocker-back/internal/stock"
)

type Command struct {
	repoStock  stock.Repository
	repoScreen screener.Repository
	logger     common.Logger
	notifier   common.Notifier
}

func NewCommand(repoStock stock.Repository, repoScreen screener.Repository, logger common.Logger, notifier common.Notifier) *Command { //nolint:lll
	return &Command{
		repoStock:  repoStock,
		repoScreen: repoScreen,
		logger:     logger,
		notifier:   notifier,
	}
}

func (c *Command) UpdateDailyData() error {
	c.logger.Debugf("UpdateDailyData()", "message", "start...")
	dailyDataToCrawl, err := c.repoStock.GetDailyDataLastAll()
	if err != nil {
		return err
	}

	for _, val := range dailyDataToCrawl {
		c.logger.Debugf("CHECK", "ticker", val.Ticker, "date", val.Date)
	}

	c.logger.Debugf("UpdateDailyData()", "message", "crawl...")
	crawlService := common.NewCrawlService(c.logger)
	dailyDataNew := crawlService.CrawlDailyDataToDate(dailyDataToCrawl)

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

	err = c.repoScreen.SetScreenAll(screens)
	if err != nil {
		return err
	}

	return nil
}
