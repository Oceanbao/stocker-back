package usecase

import (
	"fmt"

	"example.com/stocker-back/internal/common"
	"example.com/stocker-back/internal/stock"
)

type Command struct {
	repoStock stock.Repository
	logger    common.Logger
	notifier  common.Notifier
}

func NewCommand(repoStock stock.Repository, logger common.Logger, notifier common.Notifier) *Command {
	return &Command{
		repoStock: repoStock,
		logger:    logger,
		notifier:  notifier,
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
