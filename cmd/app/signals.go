package main

const (
	BIGNUM  = 9999
	HUNDRED = 100
	FIFTY   = 50
)

func MACD(closing []float64) ([]float64, []float64) {
	ema12 := make([]float64, 0)
	ema26 := make([]float64, 0)
	diff := make([]float64, 0)
	dea := make([]float64, 0)

	for idx, c := range closing {
		if idx == 0 {
			thisEma12 := (c*float64(11) + c*float64(2)) / float64(13) //nolint:gomnd //ignore
			thisEma26 := (c*float64(25) + c*float64(2)) / float64(27) //nolint:gomnd //ignore
			thisDiff := thisEma12 - thisEma26
			thisDea := (thisDiff*float64(8) + thisDiff*float64(2)) / float64(10) //nolint:gomnd //ignore

			ema12 = append(ema12, thisEma12)
			ema26 = append(ema26, thisEma26)
			diff = append(diff, thisDiff)
			dea = append(dea, thisDea)
			continue
		}

		thisEma12 := (ema12[idx-1]*float64(11) + c*float64(2)) / float64(13) //nolint:gomnd //ignore
		thisEma26 := (ema26[idx-1]*float64(25) + c*float64(2)) / float64(27) //nolint:gomnd //ignore
		thisDiff := thisEma12 - thisEma26
		thisDea := (dea[idx-1]*float64(8) + thisDiff*float64(2)) / float64(10) //nolint:gomnd //ignore

		ema12 = append(ema12, thisEma12)
		ema26 = append(ema26, thisEma26)
		diff = append(diff, thisDiff)
		dea = append(dea, thisDea)
	}

	return diff, dea
}

func KDJ(high, low, closing []float64) ([]float64, []float64, []float64) { //nolint:gocognit //ignore
	allRSV := make([]float64, 0)
	rsv := func(idx int) float64 {
		getRange := func(input []float64) []float64 {
			diff := idx + 1 - 9 //nolint:gomnd // ignore
			var start int
			if diff > 0 {
				start = diff
			} else {
				start = 0
			}
			return input[start : idx+1]
		}
		reduceLowest := func(input []float64) float64 {
			acc := float64(BIGNUM)
			for _, v := range input {
				if acc < v {
					continue
				}
				acc = v
			}
			return acc
		}
		reduceHighest := func(input []float64) float64 {
			acc := float64(-1)
			for _, v := range input {
				if acc > v {
					continue
				}
				acc = v
			}
			return acc
		}

		C := closing[idx]
		L := reduceLowest(getRange(low))
		H := reduceHighest(getRange(high))

		result := ((C - L) / (H - L)) * float64(HUNDRED)

		if result != result {
			result = float64(FIFTY)
		}

		return result
	}

	for idx := range closing {
		allRSV = append(allRSV, rsv(idx))
	}

	allK := make([]float64, 0)
	allD := make([]float64, 0)
	allJ := make([]float64, 0)

	for idx, rsv := range allRSV {
		if idx == 0 {
			allK = append(allK, FIFTY)
			allD = append(allD, FIFTY)
			allJ = append(allJ, FIFTY)
			continue
		}

		thisK := (float64(2)*allK[idx-1] + rsv) / float64(3)   //nolint:gomnd //ignore
		thisD := (float64(2)*allD[idx-1] + thisK) / float64(3) //nolint:gomnd //ignore
		thisJ := float64(3)*thisK - float64(2)*thisD           //nolint:gomnd //ignore

		allK = append(allK, thisK)
		allD = append(allD, thisD)
		allJ = append(allJ, thisJ)
	}

	return allK, allD, allJ
}

func RSI(closing []float64) []float64 {
	period := 6
	pastAvgPeriod := period - 1
	rsi := make([]float64, len(closing))
	gains := make([]float64, len(closing))
	losses := make([]float64, len(closing))

	var rsGain float64
	var rsLoss float64

	for idx := range closing {
		if idx == 0 {
			gains[idx] = 0
			losses[idx] = 0
			rsi[idx] = 0
			rsGain = 0.0
			rsLoss = 0.0
			continue
		}
		diff := closing[idx] - closing[idx-1]
		if diff > 0 {
			gains[idx] = diff
			losses[idx] = 0
		} else {
			gains[idx] = 0
			losses[idx] = -diff
		}

		rsGain = (rsGain*float64(pastAvgPeriod) + gains[idx]) / float64(period)
		rsLoss = (rsLoss*float64(pastAvgPeriod) + losses[idx]) / float64(period)

		if rsGain == 0.0 || rsLoss == 0.0 {
			rsi[idx] = 0.0
			continue
		}

		rs := rsGain / rsLoss
		rsi[idx] = (rs / (1.0 + rs)) * float64(HUNDRED)
	}

	// sma := SMA(period, rsi)

	return rsi
}

// Simple Moving Average (SMA).
// func SMA(period int, values []float64) []float64 {
// 	result := make([]float64, 0)

// 	for i, value := range values {
// 		if i == 0 {
// 			result = append(result, value)
// 			continue
// 		}
// 		thisSma := (result[i-1]*float64(period-1) + value) / float64(period)
// 		result = append(result, thisSma)
// 	}

// 	return result
// }
