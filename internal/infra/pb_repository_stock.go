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

// RecordStock is type of PB database collection schema.
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
	RankTotalCap       int     `db:"ranktotalcap" json:"ranktotalcap"`
	RankNetAsset       int     `db:"ranknetasset" json:"ranknetasset"`
	RankNetProfit      int     `db:"ranknetprofit" json:"ranknetprofit"`
	RankGrossMargin    int     `db:"rankgrossmargin" json:"rankgrossmargin"`
	RankPER            int     `db:"rankper" json:"rankper"`
	RankPBR            int     `db:"rankpbr" json:"rankpbr"`
	RankNetMargin      int     `db:"ranknetmargin" json:"ranknetmargin"`
	RankROE            int     `db:"rankroe" json:"rankroe"`
	Sector             string  `db:"sector" json:"sector"`
	SectorTotal        int     `db:"sectortotal" json:"sectortotal"`
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
		RankTotalCap:       r.RankTotalCap,
		RankNetAsset:       r.RankNetAsset,
		RankNetProfit:      r.RankNetProfit,
		RankGrossMargin:    r.RankGrossMargin,
		RankPER:            r.RankPER,
		RankPBR:            r.RankPBR,
		RankNetMargin:      r.RankNetMargin,
		RankROE:            r.RankROE,
		Sector:             r.Sector,
		SectorTotal:        r.SectorTotal,
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

func (repo *StockRepositoryPB) GetStocksBySector(sector string) ([]stock.Stock, error) {
	var records []RecordStock

	err := repo.pb.Dao().DB().
		Select().
		From("stocks").
		Where(dbx.NewExp("sector = {:sector}", dbx.Params{"sector": sector})).
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
	expr := dbx.NewExp("ticker = {:ticker}", dbx.Params{"ticker": stock.Ticker})
	records, err := repo.pb.Dao().FindRecordsByExpr("stocks", expr)
	if err != nil {
		repo.pb.Logger().Error("SetStock: fail to find record", "error", err.Error(), "ticker", stock.Ticker)
		return err
	}

	recordUnique := records[0]
	newRecordData, err := stock.ToMap()
	if err != nil {
		return err
	}
	recordUnique.Load(newRecordData)

	err = repo.pb.Dao().SaveRecord(recordUnique)
	if err != nil {
		repo.pb.Logger().Error("SetStock: cannot write to `stocks`", "error", err.Error())
		return err
	}

	return nil
}

// SetStocks upsert pb databse with given stocks.
func (repo *StockRepositoryPB) SetStocks(stocks []stock.Stock) error {
	for _, stock := range stocks {
		expr := dbx.NewExp("ticker = {:ticker}", dbx.Params{"ticker": stock.Ticker})
		records, err := repo.pb.Dao().FindRecordsByExpr("stocks", expr)
		if err != nil {
			repo.pb.Logger().Error(
				"SetStocks - failed to find record from `stocks` - skip",
				"error", err.Error(),
				"ticker", stock.Ticker,
			)
			continue
		}

		recordUnique := records[0]
		recordUnique.MarkAsNotNew()

		newRecordData, err := stock.ToMap()
		if err != nil {
			repo.pb.Logger().Error(
				"SetStocks - failed to ToMap stock - skip",
				"error", err.Error(),
				"ticker", stock.Ticker,
			)
			continue
		}
		recordUnique.Load(newRecordData)

		if err = repo.pb.Dao().SaveRecord(recordUnique); err != nil {
			repo.pb.Logger().Error(
				"SetStocks - cannot write to `stocks` - skip",
				"error", err.Error(),
				"ticker", stock.Ticker,
			)
			continue
		}
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
			recordData, err := data.ToMap()
			if err != nil {
				return err
			}
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

// DeleteStockByTicker deletes `ticker` records everywhere in the database.
func (repo *StockRepositoryPB) DeleteStockByTicker(ticker string) error {
	repo.pb.Logger().Info("DeleteStockByTicker", "ticker", ticker)
	collections, err := repo.pb.Dao().FindCollectionsByType(models.CollectionTypeBase)
	if err != nil {
		return err
	}

	for _, c := range collections {
		if !fieldInCollection(repo.pb.Dao(), "ticker", c) {
			continue
		}
		if err := deleteRecords(repo.pb.Dao(), c.Name, "ticker", ticker); err != nil {
			repo.pb.Logger().Info("deleteRecords", "error", err.Error())
			return err
		}
	}

	return nil
}

// deleteRecords deletes all records in `collectionName` where `fieldName` = `fieldValue`.
func deleteRecords(dao *daos.Dao, collectionName, fieldName, fieldValue string) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE %v = '%v'", collectionName, fieldName, fieldValue)
	if _, err := dao.DB().NewQuery(query).Execute(); err != nil {
		return err
	}

	return nil
}

// fieldInCollection checks if `field` is in `collection`.
func fieldInCollection(dao *daos.Dao, field string, collection *models.Collection) bool {
	record := models.Record{
		BaseModel: collection.BaseModel,
	}
	query := dao.RecordQuery(collection.Name).Limit(1)
	if err := query.One(&record); err != nil {
		return false
	}

	if record.Get(field) != nil {
		return true
	}

	return false
}
