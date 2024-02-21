package stock

// Repository is the persistence interface for stock module.
type Repository interface {
	GetStockByTicker(ticker string) (Stock, error)
	GetStocksAll() ([]Stock, error)
	GetDailyDataAll() (map[string][]DailyData, error)
	GetDailyDataLastAll() ([]DailyData, error)
	// GetStockByTicker(ticker string) (Stock, error)
	// CreateStock(stock Stock) error
	// GetDailyByDays(ticker string, days int) ([]DailyData, error)
	SetStocksAll(stocks []Stock) error
	SetDailyData(dailydata []DailyData) error
}
