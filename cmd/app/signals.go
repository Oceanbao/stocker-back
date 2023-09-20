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

func RSI(closing []float64) ([]float64, []float64) {
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

		rsGain = (rsGain*5.0 + gains[idx]) / 6.0
		rsLoss = (rsLoss*5.0 + losses[idx]) / 6.0

		if rsGain == 0.0 || rsLoss == 0.0 {
			rsi[idx] = 0.0
			continue
		}

		rs := rsGain / rsLoss
		rsi[idx] = (rs / (1.0 + rs)) * 100.0
	}

	sma := SMA(6, rsi)

	return rsi, sma
}

// func sumT[T int | float64](i []T) T {
// 	var o T
// 	for _, v := range i {
// 		o += v
// 	}
// 	return o
// }

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
