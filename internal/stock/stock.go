package stock

import (
	"github.com/samber/lo"
)

// Stock is entity for stock object.
type Stock struct {
	Ticker   string  `db:"ticker" json:"ticker"`
	Name     string  `db:"name" json:"name"`
	TotalCap float64 `db:"totalcap" json:"totcalcap"`
	TradeCap float64 `db:"tradecap" json:"tradecap"`
	Etf      bool    `db:"etf" json:"etf"`
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
	Ticker string  `db:"ticker" json:"ticker"`
	Date   string  `db:"date" json:"date"`
	Open   float64 `db:"open" json:"open"`
	High   float64 `db:"high" json:"high"`
	Low    float64 `db:"low" json:"low"`
	Close  float64 `db:"close" json:"close"`
}
