package infra

import (
	"example.com/stocker-back/internal/stock"
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
	Ticker string `db:"ticker" json:"ticker"`
	Name   string `db:"name" json:"name"`
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

func (r RecordDailyData) ToMap() map[string]any {
	return map[string]any{
		"ticker": r.Ticker,
		"date":   r.Date,

		"open":  r.Open,
		"high":  r.High,
		"low":   r.Low,
		"close": r.Close,

		"volume":     r.Volume,
		"value":      r.Value,
		"volatility": r.Volatility,
		"pchange":    r.Pchange,
		"change":     r.Change,
		"turnover":   r.Turnover,
	}
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

func (repo *StockRepositoryPB) GetStocksAll() ([]stock.Stock, error) {
	records, err := repo.pb.Dao().FindRecordsByExpr("stocks")
	if err != nil {
		return []stock.Stock{}, err
	}

	stocks := make([]stock.Stock, len(records))
	for idx := range records {
		stocks[idx].Ticker = records[idx].GetString("code")
		stocks[idx].Name = records[idx].GetString("name")
	}

	return stocks, nil
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

func (repo *StockRepositoryPB) SetDailyData(dailyData []stock.DailyData) error {
	collection, err := repo.pb.Dao().FindCollectionByNameOrId("daily")
	if err != nil {
		return err
	}

	err = repo.pb.Dao().RunInTransaction(func(txDao *daos.Dao) error {
		for _, data := range dailyData {
			recordData := convertDailyDataToRecord(data).ToMap()
			record := models.NewRecord(collection)
			record.Load(recordData)

			err = txDao.SaveRecord(record)
			if err != nil {
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

func convertDailyDataToRecord(dailyData stock.DailyData) RecordDailyData {
	return RecordDailyData{
		Ticker: dailyData.Ticker,
		Date:   dailyData.Date,

		Open:  dailyData.Open,
		High:  dailyData.High,
		Low:   dailyData.Low,
		Close: dailyData.Close,

		Volume:     dailyData.Volume,
		Value:      dailyData.Value,
		Volatility: dailyData.Volatility,
		Pchange:    dailyData.Pchange,
		Change:     dailyData.Change,
		Turnover:   dailyData.Turnover,
	}
}
