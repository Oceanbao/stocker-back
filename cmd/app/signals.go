package main

import "github.com/cinar/indicator"

const HUNDRED = 100

func MACD(closing []float64) ([]float64, []float64) {
	macd, signal := indicator.Macd(closing)
	return macd, signal
}

func KDJ(high, low, closing []float64) ([]float64, []float64, []float64) {
	return indicator.DefaultKdj(high, low, closing)
}

func RSI(period int, closing []float64) ([]float64, []float64) {
	rsi := make([]float64, len(closing))
	diff := make([]float64, len(closing))
	gains := make([]float64, len(closing))
	losses := make([]float64, len(closing))

	for idx := range closing {
		if idx == 0 {
			diff[idx] = 0
			continue
		}
		percentDelta := (closing[idx] - closing[idx-1]) / closing[idx-1]
		diff[idx] = percentDelta
	}

	for idx, val := range diff {
		if idx == 0 {
			gains[idx] = 0
			losses[idx] = 0
			continue
		}
		if val > 0 {
			gains[idx] = val
			losses[idx] = 0
		} else {
			gains[idx] = 0
			losses[idx] = -val
		}
	}

	for idx := range closing {
		if idx >= period {
			avgGain := sumT(gains[idx-period+1:idx+1]) / float64(period)
			avgLoss := sumT(losses[idx-period+1:idx+1]) / float64(period)
			rs := avgGain / avgLoss
			rsi[idx] = HUNDRED - (HUNDRED / (1 + rs))
		} else {
			rsi[idx] = 0
		}
	}

	sma := SMA(period, rsi)

	return rsi, sma
}

func sumT[T int | float64](i []T) T {
	var o T
	for _, v := range i {
		o += v
	}
	return o
}

// Simple Moving Average (SMA).
func SMA(period int, values []float64) []float64 {
	result := make([]float64, len(values))
	sum := float64(0)

	for i, value := range values {
		count := i + 1
		sum += value

		if i >= period {
			sum -= values[i-period]
			count = period
		}

		result[i] = sum / float64(count)
	}

	return result
}
