package stock

import (
	"encoding/json"

	"github.com/samber/lo"
)

// Stock is entity for stock object.
type Stock struct {
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

func NewEmptyStock() Stock {
	return Stock{
		Ticker:             "",
		Name:               "",
		ETF:                false,
		DateOfPublic:       "",
		EPS:                0.0,
		UndistProfit:       0.0,
		TotalShare:         0.0,
		TotalShareOut:      0.0,
		TotalCap:           0.0,
		TradeCap:           0.0,
		NetAsset:           0.0,
		NetAssetPerShare:   0.0,
		NetProfit:          0.0,
		NetProfitChange:    0.0,
		ProfitMargin:       0.0,
		PricePerEarning:    0.0,
		PricePerBook:       0.0,
		ROE:                0.0,
		TotalRevenue:       0.0,
		TotalRevenueChange: 0.0,
		GrossProfitMargin:  0.0,
		DebtRatio:          0.0,
		RankTotalCap:       0,
		RankNetAsset:       0,
		RankNetProfit:      0,
		RankGrossMargin:    0,
		RankPER:            0,
		RankPBR:            0,
		RankNetMargin:      0,
		RankROE:            0,
		Sector:             "",
		SectorTotal:        0,
	}
}

func (s *Stock) ToMap() (map[string]interface{}, error) {
	var m map[string]interface{}
	b, err := json.Marshal(*s)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// OHLC is valueobject for daily open,high,low,close prices.
type OHLC struct {
	Date  string  `db:"date" json:"date"`
	Open  float64 `db:"open" json:"open"`
	High  float64 `db:"high" json:"high"`
	Low   float64 `db:"low" json:"low"`
	Close float64 `db:"close" json:"close"`
}

func OHLC2Close(candles []OHLC) []float64 {
	return lo.Map(candles, func(candle OHLC, _ int) float64 {
		return candle.Close
	})
}

func OHLC2High(candles []OHLC) []float64 {
	return lo.Map(candles, func(candle OHLC, _ int) float64 {
		return candle.High
	})
}

func OHLC2Low(candles []OHLC) []float64 {
	return lo.Map(candles, func(candle OHLC, _ int) float64 {
		return candle.Low
	})
}

// DailyData is the aggregate valueobject for all daily time serie data.
type DailyData struct {
	Ticker     string  `db:"ticker" json:"ticker"`
	Date       string  `db:"date" json:"date"`
	Open       float64 `db:"open" json:"open"`
	High       float64 `db:"high" json:"high"`
	Low        float64 `db:"low" json:"low"`
	Close      float64 `db:"close" json:"close"`
	Volume     float64 `db:"volume" json:"volume"`
	Value      float64 `db:"value" json:"value"`
	Volatility float64 `db:"volatility" json:"volatility"`
	Pchange    float64 `db:"pchange" json:"pchange"`
	Change     float64 `db:"change" json:"change"`
	Turnover   float64 `db:"turnover" json:"turnover"`
}

func NewEmptyDailyData() DailyData {
	return DailyData{
		Ticker:     "",
		Date:       "",
		Open:       0.0,
		High:       0.0,
		Low:        0.0,
		Close:      0.0,
		Volume:     0.0,
		Value:      0.0,
		Volatility: 0.0,
		Pchange:    0.0,
		Change:     0.0,
		Turnover:   0.0,
	}
}

func (d *DailyData) ToMap() (map[string]interface{}, error) {
	var m map[string]interface{}
	b, err := json.Marshal(*d)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}
