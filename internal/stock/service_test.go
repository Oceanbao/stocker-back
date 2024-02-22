//nolint:testpackage,lll //ignore
package stock

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func loadOHLC() ([]OHLC, error) {
	jsonFile, err := os.Open("testdata/test_ohlc.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	var ohlc []OHLC
	if err = json.Unmarshal(bytes, &ohlc); err != nil {
		return nil, err
	}

	return ohlc, nil
}

func loadJSON(filename string) (map[string][]float64, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	var output map[string][]float64
	if err = json.Unmarshal(bytes, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func loadSMAStruct(filename string) ([]SMA, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	var sma []SMA
	if err = json.Unmarshal(bytes, &sma); err != nil {
		return nil, err
	}

	return sma, nil
}

func float2string(number float64) string {
	return fmt.Sprintf("%f.4", number)
}

func TestComputeSMAInternal(t *testing.T) {
	ohlc, err := loadOHLC()
	if err != nil {
		t.Fatalf("fail to loadOHLC: %v", err)
	}
	smaGold, err := loadJSON("testdata/test_sma.json")
	if err != nil {
		t.Fatalf("fail to loadJSON")
	}

	tests := []struct {
		name string
		n    int
		want []float64
	}{
		{name: "SMA5", n: 5, want: smaGold["5"]},
		{name: "SMA10", n: 10, want: smaGold["10"]},
		{name: "SMA20", n: 20, want: smaGold["20"]},
		{name: "SMA30", n: 30, want: smaGold["30"]},
		{name: "SMA90", n: 90, want: smaGold["90"]},
		{name: "SMA120", n: 120, want: smaGold["120"]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computeSMA(OHLC2Close(ohlc), float64(tt.n))
			if err != nil {
				t.Fatal("fail to run computeSMA()")
			}
			var gotRounded []string
			for _, val := range got {
				gotRounded = append(gotRounded, float2string(val))
			}
			var wantRounded []string
			for _, val := range tt.want {
				wantRounded = append(wantRounded, float2string(val))
			}
			assert.Equal(t, wantRounded, gotRounded)
		})
	}
}

func TestComputeSMAPubic(t *testing.T) {
	ohlc, err := loadOHLC()
	if err != nil {
		t.Fatalf("fail to run loadOHLC()")
	}
	want, err := loadSMAStruct("testdata/test_sma_struct.json")
	if err != nil {
		t.Fatalf("fail to run loadSMAStruct()")
	}

	got, err := ComputeSMA(ohlc)
	if err != nil {
		t.Fatal("fail to run ComputeSMA()")
	}

	var wantStr [][]string
	var gotStr [][]string

	for _, val := range want {
		var arr []string
		arr = append(arr, float2string(val.Sma5))
		arr = append(arr, float2string(val.Sma10))
		arr = append(arr, float2string(val.Sma20))
		arr = append(arr, float2string(val.Sma30))
		arr = append(arr, float2string(val.Sma90))
		arr = append(arr, float2string(val.Sma120))
		wantStr = append(wantStr, arr)
	}

	for _, val := range got {
		var arr []string
		arr = append(arr, float2string(val.Sma5))
		arr = append(arr, float2string(val.Sma10))
		arr = append(arr, float2string(val.Sma20))
		arr = append(arr, float2string(val.Sma30))
		arr = append(arr, float2string(val.Sma90))
		arr = append(arr, float2string(val.Sma120))
		gotStr = append(gotStr, arr)
	}

	assert.Equal(t, wantStr, gotStr)
}

func TestComputeSMAOneInternal(t *testing.T) {
	ohlc, err := loadOHLC()
	if err != nil {
		t.Fatalf("fail to loadOHLC")
	}
	smaGold, err := loadJSON("testdata/test_sma_one.json")
	if err != nil {
		t.Fatalf("fail to loadJSON")
	}

	tests := []struct {
		name string
		n    int
		want []float64
	}{
		{name: "SMA5", n: 5, want: smaGold["5"]},
		{name: "SMA10", n: 10, want: smaGold["10"]},
		{name: "SMA20", n: 20, want: smaGold["20"]},
		{name: "SMA30", n: 30, want: smaGold["30"]},
		{name: "SMA90", n: 90, want: smaGold["90"]},
		{name: "SMA120", n: 120, want: smaGold["120"]},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []float64
			allSma, err := computeSMA(OHLC2Close(ohlc), float64(tt.n))
			if err != nil {
				t.Fatal("fail to run computeSMA()")
			}
			for idx := range tt.want {
				closes := OHLC2Close(ohlc)
				sma := computeSMAOne(closes[:len(closes)-idx], allSma[len(allSma)-(idx+2)], float64(tt.n))
				got = append(got, sma)
			}

			var gotRounded []string
			for _, val := range got {
				gotRounded = append(gotRounded, float2string(val))
			}
			var wantRounded []string
			for _, val := range tt.want {
				wantRounded = append(wantRounded, float2string(val))
			}

			assert.Equal(t, wantRounded, gotRounded)
		})
	}
}

func TestComputeMACD(t *testing.T) {
	ohlc, err := loadOHLC()
	if err != nil {
		t.Fatalf("fail to loadOHLC")
	}
	gold, err := loadJSON("testdata/test_macd.json")
	if err != nil {
		t.Fatalf("fail to loadJSON")
	}

	gotEma12, gotEma26, gotDiff, gotDea, gotHist := computeMACD(OHLC2Close(ohlc))

	assert.Equal(t, gold["ema12"], gotEma12)
	assert.Equal(t, gold["ema26"], gotEma26)
	assert.Equal(t, gold["diff"], gotDiff)
	assert.Equal(t, gold["dea"], gotDea)
	assert.Equal(t, gold["hist"], gotHist)
}

func TestComputeMACDOne(t *testing.T) {
	ohlc, err := loadOHLC()
	if err != nil {
		t.Fatalf("fail to loadOHLC")
	}
	gold, err := loadJSON("testdata/test_macd_one.json")
	if err != nil {
		t.Fatalf("fail to loadJSON")
	}

	closes := OHLC2Close(ohlc)
	ema12, ema26, _, dea, _ := computeMACD(closes)

	var gotEma12 []float64
	var gotEma26 []float64
	var gotDiff []float64
	var gotDea []float64
	var gotHist []float64

	length := len(closes)

	for _, idx := range lo.RangeFrom(1, 10) {
		aEma12, aEma26, aDiff, aDea, aHist := computeMACDOne(closes[length-idx], ema12[length-(idx+1)], ema26[length-(idx+1)], dea[length-(idx+1)])
		gotEma12 = append(gotEma12, aEma12)
		gotEma26 = append(gotEma26, aEma26)
		gotDiff = append(gotDiff, aDiff)
		gotDea = append(gotDea, aDea)
		gotHist = append(gotHist, aHist)
	}

	assert.Equal(t, gold["ema12"], gotEma12)
	assert.Equal(t, gold["ema26"], gotEma26)
	assert.Equal(t, gold["diff"], gotDiff)
	assert.Equal(t, gold["dea"], gotDea)
	assert.Equal(t, gold["hist"], gotHist)
}

func TestComputeRSI(t *testing.T) {
	ohlc, err := loadOHLC()
	if err != nil {
		t.Fatalf("fail to loadOHLC")
	}
	gold, err := loadJSON("testdata/test_rsi.json")
	if err != nil {
		t.Fatalf("fail to loadJSON")
	}

	closes := OHLC2Close(ohlc)
	gotRsi, gotRsGains, gotRsLosses := computeRSI(closes)

	assert.Equal(t, gold["rsi"], gotRsi)
	assert.Equal(t, gold["rsGain"], gotRsGains)
	assert.Equal(t, gold["rsLoss"], gotRsLosses)
}

func TestComputeRSIOne(t *testing.T) {
	ohlc, err := loadOHLC()
	if err != nil {
		t.Fatalf("fail to loadOHLC")
	}
	gold, err := loadJSON("testdata/test_rsi_one.json")
	if err != nil {
		t.Fatalf("fail to loadJSON")
	}

	closes := OHLC2Close(ohlc)
	_, rsGains, rsLosses := computeRSI(closes)

	var gotRsi []float64
	var gotRsGain []float64
	var gotRsLoss []float64

	length := len(closes)

	for _, idx := range lo.RangeFrom(1, 10) {
		thisRsi, thisRsGain, thisRsLoss := computeRSIOne(closes[length-idx], closes[length-(idx+1)], rsGains[length-(idx+1)], rsLosses[length-(idx+1)])
		gotRsi = append(gotRsi, thisRsi)
		gotRsGain = append(gotRsGain, thisRsGain)
		gotRsLoss = append(gotRsLoss, thisRsLoss)
	}

	assert.Equal(t, gold["rsi"], gotRsi)
	assert.Equal(t, gold["rsGain"], gotRsGain)
	assert.Equal(t, gold["rsLoss"], gotRsLoss)
}

func TestComputeKDJ(t *testing.T) {
	ohlc, err := loadOHLC()
	if err != nil {
		t.Fatalf("fail to loadOHLC")
	}
	gold, err := loadJSON("testdata/test_kdj.json")
	if err != nil {
		t.Fatalf("fail to loadJSON")
	}

	closes := OHLC2Close(ohlc)
	highs := OHLC2High(ohlc)
	lows := OHLC2Low(ohlc)
	_, gotK, gotD, gotJ := computeKDJ(closes, highs, lows)

	assert.Equal(t, gold["k"], gotK)
	assert.Equal(t, gold["d"], gotD)
	assert.Equal(t, gold["j"], gotJ)
}

func TestComputeKDJOne(t *testing.T) {
	ohlc, err := loadOHLC()
	if err != nil {
		t.Fatalf("fail to loadOHLC")
	}
	gold, err := loadJSON("testdata/test_kdj_one.json")
	if err != nil {
		t.Fatalf("fail to loadJSON")
	}

	closes := OHLC2Close(ohlc)
	highs := OHLC2High(ohlc)
	lows := OHLC2Low(ohlc)

	_, k, d, _ := computeKDJ(closes, highs, lows)

	var gotK []float64
	var gotD []float64
	var gotJ []float64

	total := len(closes)

	for _, idx := range lo.RangeFrom(1, 29) {
		_, thisK, thisD, thisJ := computeKDJOne(closes[total-(12+idx):total-(idx-1)], highs[total-(12+idx):total-(idx-1)], lows[total-(12+idx):total-(idx-1)], k[total-(idx+1)], d[total-(idx+1)])
		gotK = append(gotK, thisK)
		gotD = append(gotD, thisD)
		gotJ = append(gotJ, thisJ)
	}

	assert.Equal(t, gold["k"], gotK)
	assert.Equal(t, gold["d"], gotD)
	assert.Equal(t, gold["j"], gotJ)
}
