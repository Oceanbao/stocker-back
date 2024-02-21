package usecase

import (
	"example.com/stocker-back/internal/common"
	"example.com/stocker-back/internal/screener"
	"example.com/stocker-back/internal/stock"
)

type Query struct {
	repoStock  stock.Repository
	repoScreen screener.Repository
	logger     common.Logger
	notifier   common.Notifier
}

func NewQuery(repoStock stock.Repository, repoScreen screener.Repository, logger common.Logger, notifier common.Notifier) *Query { //nolint:lll
	return &Query{
		repoStock:  repoStock,
		repoScreen: repoScreen,
		logger:     logger,
		notifier:   notifier,
	}
}

func (c *Query) GetStockByTicker(ticker string) (stock.Stock, error) {
	stockFound, err := c.repoStock.GetStockByTicker(ticker)
	if err != nil {
		return stock.NewEmptyStock(), err
	}

	return stockFound, nil
}
