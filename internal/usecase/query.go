package usecase

import (
	"encoding/json"
	"slices"

	"example.com/stocker-back/internal/infra"
	"example.com/stocker-back/internal/screener"
	"example.com/stocker-back/internal/stock"
	"example.com/stocker-back/internal/tracking"
	"github.com/samber/lo"
)

type Query struct {
	repoStock    stock.Repository
	repoScreen   screener.Repository
	repoTracking tracking.Repository
	logger       infra.Logger
	notifier     infra.Notifier
}

// DELE: fix this into config.
func NewQuery(repoStock stock.Repository, repoScreen screener.Repository, repoTracking tracking.Repository, logger infra.Logger, notifier infra.Notifier) *Query { //nolint:lll
	return &Query{
		repoStock:    repoStock,
		repoScreen:   repoScreen,
		repoTracking: repoTracking,
		logger:       logger,
		notifier:     notifier,
	}
}

func (q *Query) GetStockByTicker(ticker string) (stock.Stock, error) {
	stockFound, err := q.repoStock.GetStockByTicker(ticker)
	if err != nil {
		// DELE: err means likely no row found, therefore pass empty Stock.
		return stock.NewEmptyStock(), err
	}

	return stockFound, nil
}

func (q *Query) GetRandomTickers(num int) ([]string, error) {
	stocksAll, err := q.repoStock.GetStocks()
	if err != nil {
		return nil, err
	}

	tickers := lo.Map(stocksAll, func(stock stock.Stock, _ int) string {
		return stock.Ticker
	})

	return lo.Samples(tickers, num), nil
}

// GetScreens queries screens data augmented with necessary meta.
func (q *Query) GetScreens() ([]map[string]interface{}, error) {
	screens, err := q.repoScreen.GetScreens()
	if err != nil {
		return nil, err
	}

	// DELE: better shape
	trackings, err := q.repoTracking.GetTrackings()
	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, s := range screens {
		// DELE: better shape
		if s.Kdj > 30 { //nolint:gomnd // ignore
			continue
		}

		var m map[string]interface{}

		stock, err := q.repoStock.GetStockByTicker(s.Ticker)
		if err != nil {
			return nil, err
		}
		b, err := json.Marshal(stock)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, err
		}
		m["screenkdj"] = s.Kdj

		isTracked := slices.ContainsFunc(trackings, func(t tracking.Tracking) bool {
			return t.Ticker == stock.Ticker
		})
		if isTracked {
			m["tracking"] = true
		} else {
			m["tracking"] = false
		}

		dailyData, err := q.repoStock.GetDailyDataLastByTicker(s.Ticker)
		if err != nil {
			continue
		}
		m["dailyvalue"] = dailyData.Value

		output = append(output, m)
	}

	return output, nil
}

func (q *Query) GetTrackings() ([]map[string]interface{}, error) {
	trackings, err := q.repoTracking.GetTrackings()
	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, s := range trackings {
		stock, err := q.repoStock.GetStockByTicker(s.Ticker)
		if err != nil {
			continue
		}
		var m map[string]interface{}
		b, err := json.Marshal(stock)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, &m); err != nil {
			continue
		}

		m["tracking"] = true

		output = append(output, m)
	}

	return output, nil
}

func (q *Query) GetStocksBySector(sector string) ([]map[string]any, error) {
	stocks, err := q.repoStock.GetStocksBySector(sector)
	if err != nil {
		return nil, err
	}

	// DELE: better shape
	trackings, err := q.repoTracking.GetTrackings()
	if err != nil {
		return nil, err
	}

	var output []map[string]interface{}
	for _, s := range stocks {
		var m map[string]interface{}
		b, err := json.Marshal(s)
		if err != nil {
			continue
		}
		if err := json.Unmarshal(b, &m); err != nil {
			continue
		}

		isTracked := slices.ContainsFunc(trackings, func(t tracking.Tracking) bool {
			return t.Ticker == s.Ticker
		})
		if isTracked {
			m["tracking"] = true
		} else {
			m["tracking"] = false
		}

		output = append(output, m)
	}

	return output, nil
}
