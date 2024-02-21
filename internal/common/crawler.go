package common

import (
	"context"
	"io"
	"net/http"

	"example.com/stocker-back/internal/stock"
)

type DailyCrawler interface {
	CrawlStock(tickers []string) []stock.Stock
	// CrawlRank(tickers []string) []stock.Rank
	CrawlDaily(dailyData []stock.DailyData) []stock.DailyData
}

func Fetch(ctx context.Context, url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
