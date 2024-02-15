//nolint:gomnd //ignore
package stock

import (
	"github.com/samber/lo"
)

// SMA is valueobject holding all SMA values for a given day.
type SMA struct {
	Sma5   float64 `db:"sma5" json:"sma5"`
	Sma10  float64 `db:"sma10" json:"sma10"`
	Sma20  float64 `db:"sma20" json:"sma20"`
	Sma30  float64 `db:"sma30" json:"sma30"`
	Sma90  float64 `db:"sma90" json:"sma90"`
	Sma120 float64 `db:"sma120" json:"sma120"`
}

// ComputeSMA wraps computeSMA to handle typed objects.
func ComputeSMA(candles []OHLC) ([]SMA, error) {
	result := make([]SMA, len(candles))
	closes := OHLC2Close(candles)
	nSMA := map[int][]float64{}

	if len(candles) >= 5 {
		nSMA[5] = nil
	}
	if len(candles) >= 10 {
		nSMA[10] = nil
	}
	if len(candles) >= 20 {
		nSMA[20] = nil
	}
	if len(candles) >= 30 {
		nSMA[30] = nil
	}
	if len(candles) >= 90 {
		nSMA[90] = nil
	}
	if len(candles) >= 120 {
		nSMA[120] = nil
	}

	for n := range nSMA {
		sma, err := computeSMA(closes, float64(n))
		if err != nil {
			return nil, err
		}
		nSMA[n] = sma
	}

	for idx := range lo.Range(len(candles)) {
		for n, sma := range nSMA {
			switch n {
			case 5: //nolint:gomnd //ignore
				result[idx].Sma5 = sma[idx]
			case 10: //nolint:gomnd //ignore
				result[idx].Sma10 = sma[idx]
			case 20: //nolint:gomnd //ignore
				result[idx].Sma20 = sma[idx]
			case 30: //nolint:gomnd //ignore
				result[idx].Sma30 = sma[idx]
			case 90: //nolint:gomnd //ignore
				result[idx].Sma90 = sma[idx]
			case 120: //nolint:gomnd //ignore
				result[idx].Sma120 = sma[idx]
			}
		}
	}

	return result, nil
}

// computeSMA calculates all SMA given a time serie close prices.
func computeSMA(closes []float64, n float64) ([]float64, error) {
	sma := make([]float64, 0, len(closes))
	initSma := lo.Sum(lo.Slice(closes, 0, int(n))) / n
	for range lo.Range(int(n)) {
		sma = append(sma, initSma)
	}

	for _, idx := range lo.RangeFrom(int(n), len(closes)-int(n)) {
		lastSMA, err := lo.Last(sma)
		if err != nil {
			return nil, err
		}
		thisSMA := lastSMA + (closes[idx]-closes[idx-int(n)])/n
		sma = append(sma, thisSMA)
	}

	return sma, nil
}

// ComputeSMAOne wraps computeSMAOne to hanlde typed objects.
func ComputeSMAOne(candles []OHLC, lastSMA SMA) SMA {
	var sma SMA
	closes := OHLC2Close(candles)

	nSMA := map[int]float64{
		5:   0.0,
		10:  0.0,
		20:  0.0,
		30:  0.0,
		90:  0.0,
		120: 0.0,
	}

	for n := range nSMA {
		switch n {
		case 5:
			sma.Sma5 = computeSMAOne(closes, lastSMA.Sma5, float64(n))
		case 10:
			sma.Sma10 = computeSMAOne(closes, lastSMA.Sma10, float64(n))
		case 20:
			sma.Sma20 = computeSMAOne(closes, lastSMA.Sma20, float64(n))
		case 30:
			sma.Sma30 = computeSMAOne(closes, lastSMA.Sma30, float64(n))
		case 90:
			sma.Sma90 = computeSMAOne(closes, lastSMA.Sma90, float64(n))
		case 120:
			sma.Sma120 = computeSMAOne(closes, lastSMA.Sma120, float64(n))
		}
	}

	return sma
}

// computeSMAOne calculates SMA for one single day.
func computeSMAOne(closes []float64, lastSMA float64, n float64) float64 {
	return lastSMA + (closes[len(closes)-1]-closes[len(closes)-int(n)-1])/n
}

// MACD is valueobject holding all values for MACD for a given day.
type MACD struct {
	Ema12 float64 `db:"ema12" json:"ema12"`
	Ema26 float64 `db:"ema26" json:"ema26"`
	Diff  float64 `db:"diff" json:"diff"`
	Dea   float64 `db:"dea" json:"dea"`
	Hist  float64 `db:"hist" json:"hist"`
}

// ComputeMACD wraps computeMACD to take typed input output.
func ComputeMACD(candles []OHLC) []MACD {
	closes := OHLC2Close(candles)
	macd := make([]MACD, len(closes))

	ema12, ema26, diff, dea, hist := computeMACD(closes)

	for idx := range closes {
		macd[idx].Ema12 = ema12[idx]
		macd[idx].Ema26 = ema26[idx]
		macd[idx].Diff = diff[idx]
		macd[idx].Dea = dea[idx]
		macd[idx].Hist = hist[idx]
	}

	return macd
}

// computeMACD calculates all MACD for a given time seris close prices.
func computeMACD(closes []float64) ([]float64, []float64, []float64, []float64, []float64) {
	var ema12 []float64
	var ema26 []float64
	var diff []float64
	var dea []float64
	var hist []float64

	var thisEma12 float64
	var thisEma26 float64
	var thisDiff float64
	var thisDea float64
	for idx, close := range closes {
		if idx == 0 {
			thisEma12 = (close*11.0 + close*2.0) / 13.0
			thisEma26 = (close*25.0 + close*2.0) / 27.0
			thisDiff = thisEma12 - thisEma26
			thisDea = (thisDiff*8.0 + thisDiff*2.0) / 10.0
		} else {
			thisEma12 = (ema12[idx-1]*11.0 + close*2.0) / 13.0
			thisEma26 = (ema26[idx-1]*25.0 + close*2.0) / 27.0
			thisDiff = thisEma12 - thisEma26
			thisDea = (dea[idx-1]*8.0 + thisDiff*2.0) / 10.0
		}

		ema12 = append(ema12, thisEma12)
		ema26 = append(ema26, thisEma26)
		diff = append(diff, thisDiff)
		dea = append(dea, thisDea)
		hist = append(hist, 2.0*(thisDiff-thisDea))
	}

	return ema12, ema26, diff, dea, hist
}

// ComputeMACDOne wraps computeMACDOne to return typed MACD object.
func ComputeMACDOne(thisClose float64, lastMACD MACD) MACD {
	ema12, ema26, diff, dea, hist := computeMACDOne(thisClose, lastMACD.Ema12, lastMACD.Ema26, lastMACD.Dea)

	return MACD{
		Ema12: ema12,
		Ema26: ema26,
		Diff:  diff,
		Dea:   dea,
		Hist:  hist,
	}
}

// computeMACDOne calculates a single MACD from previous values.
func computeMACDOne(thisClose, lastEma12, lastEma26, lastDea float64) (float64, float64, float64, float64, float64) {
	ema12 := (lastEma12*11.0 + thisClose*2.0) / 13.0
	ema26 := (lastEma26*25.0 + thisClose*2.0) / 27.0
	diff := ema12 - ema26
	dea := (lastDea*8.0 + diff*2.0) / 10.0
	hist := 2.0 * (diff - dea)

	return ema12, ema26, diff, dea, hist
}

// RSI is valueobject holding all values for RSI for a given day.
type RSI struct {
	Rsi    float64 `db:"rsi" json:"rsi"`
	RsGain float64 `db:"rsgain" json:"rsgain"`
	RsLoss float64 `db:"rsloss" json:"rsloss"`
}

// ComputeRSI wraps computeRSI to produce RSI typed objects.
func ComputeRSI(candles []OHLC) []RSI {
	closes := OHLC2Close(candles)
	output := make([]RSI, len(closes))

	rsi, rsGains, rsLosses := computeRSI(closes)

	for idx := range closes {
		output[idx].Rsi = rsi[idx]
		output[idx].RsGain = rsGains[idx]
		output[idx].RsLoss = rsLosses[idx]
	}

	return output
}

// computeRSI calculates all RSI for a given time seris close prices.
func computeRSI(closes []float64) ([]float64, []float64, []float64) {
	n := 6.0
	pastN := n - 1.0
	rsi := make([]float64, len(closes))
	rsGains := make([]float64, len(closes))
	rsLosses := make([]float64, len(closes))

	rsi[0] = 0.0
	rsGains[0] = 0.0
	rsLosses[0] = 0.0

	for _, idx := range lo.RangeFrom(1, len(closes)-1) {
		gain := 0.0
		loss := 0.0
		diff := closes[idx] - closes[idx-1]
		if diff > 0.0 {
			gain = diff
		} else {
			loss = -diff
		}
		thisRsGain := (rsGains[idx-1]*pastN + gain) / n
		thisRsLoss := (rsLosses[idx-1]*pastN + loss) / n
		rsGains[idx] = thisRsGain
		rsLosses[idx] = thisRsLoss

		if thisRsGain == 0.0 || thisRsLoss == 0.0 {
			rsi[idx] = 0.0
		} else {
			rs := thisRsGain / thisRsLoss
			rsi[idx] = (rs / (1.0 + rs)) * 100.0
		}
	}

	return rsi, rsGains, rsLosses
}

// ComputeRSIOne wraps computeRSIOne to return typed object.
func ComputeRSIOne(thisClose, lastClose float64, lastRSI RSI) RSI {
	lastRsi, lastRsGain, lastRsLoss := computeRSIOne(thisClose, lastClose, lastRSI.RsGain, lastRSI.RsLoss)

	return RSI{
		Rsi:    lastRsi,
		RsGain: lastRsGain,
		RsLoss: lastRsLoss,
	}
}

// computeRSIOne calculates a single RSI from previous values.
func computeRSIOne(thisClose, lastClose, lastRsGain, lastRsLoss float64) (float64, float64, float64) {
	n := 6.0
	pastN := n - 1.0

	gain := 0.0
	loss := 0.0
	rsi := 0.0

	diff := thisClose - lastClose
	if diff > 0.0 {
		gain = diff
	} else {
		loss = -diff
	}

	thisRsGain := (lastRsGain*pastN + gain) / n
	thisRsLoss := (lastRsLoss*pastN + loss) / n
	if thisRsGain != 0.0 || thisRsLoss != 0.0 {
		rs := thisRsGain / thisRsLoss
		rsi = (rs / (1.0 + rs)) * 100.0
	}

	return rsi, thisRsGain, thisRsLoss
}

// KDJ is valueobject holding all values for KDJ for a given day.
type KDJ struct {
	Rsv float64 `db:"rsv" json:"rsv"`
	K   float64 `db:"k" json:"k"`
	D   float64 `db:"d" json:"d"`
	J   float64 `db:"j" json:"j"`
}

// ComputeKDJ wraps computeKDJ to produce typed objects.
func ComputeKDJ(candles []OHLC) []KDJ {
	closes := OHLC2Close(candles)
	highs := OHLC2High(candles)
	lows := OHLC2Low(candles)
	kdj := make([]KDJ, len(closes))

	rsv, k, d, j := computeKDJ(closes, highs, lows)

	for idx := range closes {
		kdj[idx].Rsv = rsv[idx]
		kdj[idx].K = k[idx]
		kdj[idx].D = d[idx]
		kdj[idx].J = j[idx]
	}
	return kdj
}

// computeKDJ calculates all KDJ for a given time seris close prices.
func computeKDJ(closes, highs, lows []float64) ([]float64, []float64, []float64, []float64) {
	rsv := make([]float64, len(closes))
	k := make([]float64, len(closes))
	d := make([]float64, len(closes))
	j := make([]float64, len(closes))

	getRange := func(idx int, values []float64) []float64 {
		var start int
		diff := idx + 1 - 9
		if diff > 0 {
			start = diff
		} else {
			start = 0
		}
		end := idx + 1

		return values[start:end]
	}

	reduceLowest := func(values []float64) float64 {
		lowest := 9999.0
		for _, value := range values {
			if lowest < value {
				continue
			}
			lowest = value
		}
		return lowest
	}

	reduceHighest := func(values []float64) float64 {
		highest := -1.0
		for _, value := range values {
			if highest > value {
				continue
			}
			highest = value
		}
		return highest
	}

	for idx := range closes {
		c := closes[idx]
		l := reduceLowest(getRange(idx, lows))
		h := reduceHighest(getRange(idx, highs))
		thisRSV := ((c - l) / (h - l)) * 100.0
		rsv[idx] = thisRSV

		var thisK float64
		var thisD float64
		var thisJ float64

		if idx == 0 {
			thisK = 50.0
			thisD = 50.0
			thisJ = 50.0
		} else {
			thisK = (2.0*k[idx-1] + thisRSV) / 3.0
			thisD = (2.0*d[idx-1] + thisK) / 3.0
			thisJ = 3.0*thisK - 2.0*thisD
		}

		k[idx] = thisK
		d[idx] = thisD
		j[idx] = thisJ
	}

	return rsv, k, d, j
}

// ComputeKDJOne wraps computeKDJ to return typed object.
func ComputeKDJOne(candles []OHLC, lastK, lastD float64) KDJ {
	closes := OHLC2Close(candles)
	highs := OHLC2High(candles)
	lows := OHLC2Low(candles)

	rsv, k, d, j := computeKDJOne(closes, highs, lows, lastK, lastD)

	return KDJ{
		Rsv: rsv,
		K:   k,
		D:   d,
		J:   j,
	}
}

// computeKDJOne calculates a single KDJ from previous values.
func computeKDJOne(closes, highs, lows []float64, lastK, lastD float64) (float64, float64, float64, float64) {
	n := 9

	reduceLowest := func(values []float64) float64 {
		lowest := 9999.0
		for _, value := range values {
			if lowest < value {
				continue
			}
			lowest = value
		}
		return lowest
	}

	reduceHighest := func(values []float64) float64 {
		highest := -1.0
		for _, value := range values {
			if highest > value {
				continue
			}
			highest = value
		}
		return highest
	}

	total := len(closes)

	c := closes[total-1]
	l := reduceLowest(lows[total-n:])
	h := reduceHighest(highs[total-n:])
	thisRSV := ((c - l) / (h - l)) * 100.0
	thisK := (2.0*lastK + thisRSV) / 3.0
	thisD := (2.0*lastD + thisK) / 3.0
	thisJ := 3.0*thisK - 2.0*thisD

	return thisRSV, thisK, thisD, thisJ
}
