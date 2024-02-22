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

func (q *Query) GetStockByTicker(ticker string) (stock.Stock, error) {
	stockFound, err := q.repoStock.GetStockByTicker(ticker)
	if err != nil {
		return stock.NewEmptyStock(), err
	}

	return stockFound, nil
}

func (q *Query) GetScreens() ([]map[string]interface{}, error) {
	screens, err := q.repoScreen.GetScreens()
	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, s := range screens {
		if s.Kdj <= 30 { //nolint:gomnd // ignore
			m := make(map[string]interface{})
			m["kdj"] = s.Kdj

			stock, err := q.repoStock.GetStockByTicker(s.Ticker)
			if err != nil {
				continue
			}
			m["stock"] = stock

			dailyData, err := q.repoStock.GetDailyDataLastByTicker(s.Ticker)
			if err != nil {
				continue
			}
			m["daily"] = dailyData

			output = append(output, m)
		}
	}

	return output, nil
}
