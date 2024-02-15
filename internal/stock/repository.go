package stock

// Repository is the persistence interface for stock module.
type Repository interface {
	GetStocksAll() ([]Stock, error)
	GetDailyDataLastAll() ([]DailyData, error)
	// GetStockByTicker(ticker string) (Stock, error)
	// CreateStock(stock Stock) error
	// GetDailyAllByTicker(ticker string) ([]DailyData, error)
	// GetDailyByDays(ticker string, days int) ([]DailyData, error)
	SetDailyData(dailydata []DailyData) error
}
