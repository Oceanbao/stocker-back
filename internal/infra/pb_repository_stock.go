package infra

import (
	"fmt"

	"example.com/stocker-back/internal/stock"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/samber/lo"
)

type StockRepositoryPB struct {
	pb *pocketbase.PocketBase
}

func NewStockRepositoryPB(pb *pocketbase.PocketBase) *StockRepositoryPB {
	return &StockRepositoryPB{
		pb: pb,
	}
}

type RecordStock struct {
	Ticker             string  `db:"ticker" json:"ticker"`
	Name               string  `db:"name" json:"name"`
	ETF                bool    `db:"etf" json:"etf"`
	DateOfPublic       string  `db:"dateofpublic" json:"dateofpublic"`
	EPS                float64 `db:"eps" json:"eps"`
	UndistProfit       float64 `db:"undistprofit" json:"undistprofit"`
	TotalShare         float64 `db:"totalshare" json:"totalshare"`
	TotalShareOut      float64 `db:"totalshareout" json:"totalshareout"`
	TotalCap           float64 `db:"totalcap" json:"totcalcap"`
	TradeCap           float64 `db:"tradecap" json:"tradecap"`
	NetAsset           float64 `db:"netasset" json:"netasset"`
	NetAssetPerShare   float64 `db:"netassetpershare" json:"netassetpershare"`
	NetProfit          float64 `db:"netprofit" json:"netprofit"`
	NetProfitChange    float64 `db:"netprofitchange" json:"netprofitchange"`
	ProfitMargin       float64 `db:"profitmargin" json:"profitmargin"`
	PricePerEarning    float64 `db:"priceperearning" json:"priceperearning"`
	PricePerBook       float64 `db:"priceperbook" json:"priceperbook"`
	ROE                float64 `db:"roe" json:"roe"`
	TotalRevenue       float64 `db:"totalrevenue" json:"totalrevenue"`
	TotalRevenueChange float64 `db:"totalrevenuechange" json:"totalrevenuechange"`
	GrossProfitMargin  float64 `db:"grossprofitmargin" json:"grossprofitmargin"`
	DebtRatio          float64 `db:"debtratio" json:"debtratio"`
}

func (r RecordStock) ToModel() stock.Stock {
	return stock.Stock{
		Ticker:             r.Ticker,
		Name:               r.Name,
		ETF:                r.ETF,
		DateOfPublic:       r.DateOfPublic,
		EPS:                r.EPS,
		UndistProfit:       r.UndistProfit,
		TotalShare:         r.TotalShare,
		TotalShareOut:      r.TotalShareOut,
		TotalCap:           r.TotalCap,
		TradeCap:           r.TradeCap,
		NetAsset:           r.NetAsset,
		NetAssetPerShare:   r.NetAssetPerShare,
		NetProfit:          r.NetProfit,
		NetProfitChange:    r.NetProfitChange,
		ProfitMargin:       r.ProfitMargin,
		PricePerEarning:    r.PricePerEarning,
		PricePerBook:       r.PricePerBook,
		ROE:                r.ROE,
		TotalRevenue:       r.TotalRevenue,
		TotalRevenueChange: r.TotalRevenueChange,
		GrossProfitMargin:  r.GrossProfitMargin,
		DebtRatio:          r.DebtRatio,
	}
}

type RecordDailyData struct {
	Ticker string `db:"ticker" json:"ticker"`
	Date   string `db:"date" json:"date"`

	Open  float64 `db:"open" json:"open"`
	High  float64 `db:"high" json:"high"`
	Low   float64 `db:"low" json:"low"`
	Close float64 `db:"close" json:"close"`

	Volume     float64 `db:"volume" json:"volume"`
	Value      float64 `db:"value" json:"value"`
	Volatility float64 `db:"volatility" json:"volatility"`
	Pchange    float64 `db:"pchange" json:"pchange"`
	Change     float64 `db:"change" json:"change"`
	Turnover   float64 `db:"turnover" json:"turnover"`
}

func (r RecordDailyData) ToModel() stock.DailyData {
	return stock.DailyData{
		Ticker:     r.Ticker,
		Date:       r.Date,
		Open:       r.Open,
		High:       r.High,
		Low:        r.Low,
		Close:      r.Close,
		Volume:     r.Volume,
		Value:      r.Value,
		Volatility: r.Volatility,
		Pchange:    r.Pchange,
		Change:     r.Change,
		Turnover:   r.Turnover,
	}
}

func (repo *StockRepositoryPB) GetStockByTicker(ticker string) (stock.Stock, error) {
	var stockFound stock.Stock
	expr := dbx.NewExp("ticker = {:ticker}", dbx.Params{"ticker": ticker})
	err := repo.pb.Dao().DB().Select().From("stocks").Where(expr).One(&stockFound)
	if err != nil {
		return stock.NewEmptyStock(), err
	}

	return stockFound, nil
}

func (repo *StockRepositoryPB) GetStocks() ([]stock.Stock, error) {
	var records []RecordStock

	err := repo.pb.Dao().DB().
		Select().
		From("stocks").
		All(&records)
	if err != nil {
		return nil, err
	}

	output := make([]stock.Stock, 0, len(records))
	for _, s := range records {
		output = append(output, s.ToModel())
	}

	return output, nil
}

func (repo *StockRepositoryPB) GetDailyDataAll() (map[string][]stock.DailyData, error) {
	var recordDailyData []RecordDailyData

	err := repo.pb.Dao().DB().
		Select().
		From("daily").
		All(&recordDailyData)
	if err != nil {
		return nil, err
	}

	recordGrouped := lo.GroupBy(recordDailyData, func(rec RecordDailyData) string {
		return rec.Ticker
	})

	output := make(map[string][]stock.DailyData, len(recordGrouped))
	for key, val := range recordGrouped {
		dd := make([]stock.DailyData, 0, len(val))
		for _, d := range val {
			dd = append(dd, d.ToModel())
		}
		output[key] = dd
	}

	return output, nil
}
func (repo *StockRepositoryPB) GetDailyDataLastByTicker(ticker string) (stock.DailyData, error) {
	var records []RecordDailyData

	err := repo.pb.Dao().DB().
		Select().
		From("daily").
		Where(dbx.NewExp("ticker = {:ticker}", dbx.Params{"ticker": ticker})).
		OrderBy("date DESC").
		All(&records)
	if err != nil {
		return stock.NewEmptyDailyData(), err
	}

	return records[0].ToModel(), nil
}

func (repo *StockRepositoryPB) GetDailyDataLastAll() ([]stock.DailyData, error) {
	var recordDailyData []RecordDailyData

	err := repo.pb.Dao().DB().
		Select().
		From("daily").
		OrderBy("date DESC").
		All(&recordDailyData)
	if err != nil {
		return nil, err
	}

	recordGrouped := lo.GroupBy(recordDailyData, func(rec RecordDailyData) string {
		return rec.Ticker
	})

	output := make([]stock.DailyData, 0, len(recordGrouped))
	for _, val := range recordGrouped {
		// Take first item since ordered by "date DESC" to get last daily.
		output = append(output, val[0].ToModel())
	}

	return output, nil
}

func (repo *StockRepositoryPB) SetStock(stock stock.Stock) error {
	collection, err := repo.pb.Dao().FindCollectionByNameOrId("stocks")
	if err != nil {
		return err
	}

	model := models.NewRecord(collection)
	record := convertStockToMap(stock)
	model.Load(record)

	err = repo.pb.Dao().SaveRecord(model)
	if err != nil {
		repo.pb.Logger().Error("cannot write to `stocks`", "error", err.Error())
		return nil
	}

	return nil
}

func (repo *StockRepositoryPB) SetStocks(stocks []stock.Stock) error {
	err := repo.pb.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, data := range stocks {
			expr := dbx.NewExp("ticker = {:ticker}", dbx.Params{"ticker": data.Ticker})
			records, err := txDao.FindRecordsByExpr("stocks", expr)
			if err != nil {
				repo.pb.Logger().Error("failed to find record from `stocks` - skip", "error", err.Error(), "ticker", data.Ticker)
				continue
			}

			recordUnique := records[0]
			newRecordData := convertStockToMap(data)
			recordUnique.Load(newRecordData)

			if err = txDao.SaveRecord(recordUnique); err != nil {
				repo.pb.Logger().Error("cannot write to `stocks`", "error", err.Error())
				continue
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (repo *StockRepositoryPB) SetDailyData(dailyData []stock.DailyData) error {
	collection, err := repo.pb.Dao().FindCollectionByNameOrId("daily")
	if err != nil {
		return err
	}

	err = repo.pb.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, data := range dailyData {
			recordData := convertDailyDataToRecord(data)
			record := models.NewRecord(collection)
			record.Load(recordData)

			err = txDao.SaveRecord(record)
			if err != nil {
				repo.pb.Logger().Error("cannot write to `daily`", "error", err.Error())
				continue
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (repo *StockRepositoryPB) DeleteStockByTicker(ticker string) error {
	expr := dbx.NewExp("ticker = {:ticker}", dbx.Params{"ticker": ticker})
	records, err := repo.pb.Dao().FindRecordsByExpr("stocks", expr)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("ticker %v not found in `stocks`", ticker)
	}

	recordUnique := records[0]

	if err := repo.pb.Dao().DeleteRecord(recordUnique); err != nil {
		return err
	}

	return nil
}

func (repo *StockRepositoryPB) DeleteDailyDataByTicker(ticker string) error {
	expr := dbx.NewExp("ticker = {:ticker}", dbx.Params{"ticker": ticker})
	records, err := repo.pb.Dao().FindRecordsByExpr("daily", expr)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("ticker %v not found in `daily`", ticker)
	}

	for _, rec := range records {
		if err := repo.pb.Dao().DeleteRecord(rec); err != nil {
			return err
		}
	}

	return nil
}

func convertStockToMap(stock stock.Stock) map[string]any {
	return map[string]any{
		"ticker":             stock.Ticker,
		"name":               stock.Name,
		"etf":                stock.ETF,
		"dateofpublic":       stock.DateOfPublic,
		"eps":                stock.EPS,
		"undistprofit":       stock.UndistProfit,
		"totalshare":         stock.TotalShare,
		"totalshareout":      stock.TotalShareOut,
		"totalcap":           stock.TotalCap,
		"tradecap":           stock.TradeCap,
		"netasset":           stock.NetAsset,
		"netassetpershare":   stock.NetAssetPerShare,
		"netprofit":          stock.NetProfit,
		"netprofitchange":    stock.NetProfitChange,
		"profitmargin":       stock.ProfitMargin,
		"priceperearning":    stock.PricePerEarning,
		"priceperbook":       stock.PricePerBook,
		"roe":                stock.ROE,
		"totalrevenue":       stock.TotalRevenue,
		"totalrevenuechange": stock.TotalRevenueChange,
		"grossprofitmargin":  stock.GrossProfitMargin,
		"debtratio":          stock.DebtRatio,
	}
}

func convertDailyDataToRecord(dailyData stock.DailyData) map[string]any {
	return map[string]any{
		"ticker": dailyData.Ticker,
		"date":   dailyData.Date,

		"open":  dailyData.Open,
		"high":  dailyData.High,
		"low":   dailyData.Low,
		"close": dailyData.Close,

		"volume":     dailyData.Volume,
		"value":      dailyData.Value,
		"volatility": dailyData.Volatility,
		"pchange":    dailyData.Pchange,
		"change":     dailyData.Change,
		"turnover":   dailyData.Turnover,
	}
}
