package stock

// Repository is the persistence interface for stock domain.
type Repository interface {
	GetStockByTicker(ticker string) (Stock, error)
	GetStocks() ([]Stock, error)
	GetStocksBySector(sector string) ([]Stock, error)
	GetDailyDataAll() (map[string][]DailyData, error)
	GetDailyDataLastByTicker(ticker string) (DailyData, error)
	GetDailyDataLastAll() ([]DailyData, error)

	CreateStock(stock Stock) error

	UpdateStock(stock Stock) error
	UpdateStocks(stocks []Stock) error

	CreateDailyData(dailydata []DailyData) error

	DeleteStockByTicker(ticker string) error
}
