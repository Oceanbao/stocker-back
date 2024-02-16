package usecase

import (
	"example.com/stocker-back/internal/common"
	"example.com/stocker-back/internal/stock"
)

type Command struct {
	repoStock stock.Repository
	logger    common.Logger
}

func NewCommand(repoStock stock.Repository, logger common.Logger) *Command {
	return &Command{
		repoStock: repoStock,
		logger:    logger,
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
	c.logger.Debugf("total crawled: [%d]", len(dailyDataNew))

	if err = c.repoStock.SetDailyData(dailyDataNew); err != nil {
		return err
	}

	c.logger.Debugf("done UpdateDailyData")

	return nil
}
