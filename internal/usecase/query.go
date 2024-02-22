package usecase

import (
	"example.com/stocker-back/internal/infra"
	"example.com/stocker-back/internal/screener"
	"example.com/stocker-back/internal/stock"
)

type Query struct {
	repoStock  stock.Repository
	repoScreen screener.Repository
	logger     infra.Logger
	notifier   infra.Notifier
}

func NewQuery(repoStock stock.Repository, repoScreen screener.Repository, logger infra.Logger, notifier infra.Notifier) *Query { //nolint:lll
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
