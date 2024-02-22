package stock

// Repository is the persistence interface for stock domain.
type Repository interface {
	GetStockByTicker(ticker string) (Stock, error)
	GetStocks() ([]Stock, error)
	GetDailyDataAll() (map[string][]DailyData, error)
	GetDailyDataLastAll() ([]DailyData, error)
	SetStock(stock Stock) error
	SetStocks(stocks []Stock) error
	SetDailyData(dailydata []DailyData) error
	DeleteStockByTicker(ticker string) error
	DeleteDailyDataByTicker(ticker string) error
}
