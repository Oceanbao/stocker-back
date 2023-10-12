package main

type recordTrack struct {
	ID      string `db:"id" json:"id"`
	Code    string `db:"code" json:"code"`
	Name    string `db:"name" json:"name"`
	Started string `db:"started" json:"started"`
}

type recordDaily struct {
	ID    string  `db:"id" json:"id"`
	Code  string  `db:"code" json:"code"`
	Date  string  `db:"date" json:"date"`
	Open  float64 `db:"open" json:"open"`
	High  float64 `db:"high" json:"high"`
	Low   float64 `db:"low" json:"low"`
	Close float64 `db:"close" json:"close"`
}

type recordAlert struct {
	Code string  `db:"code" json:"code"`
	Name string  `db:"name" json:"name"`
	Cap  float64 `db:"cap" json:"cap"`
	Rsi  float64 `db:"rsi" json:"rsi"`
	K    float64 `db:"k" json:"k"`
	D    float64 `db:"d" json:"d"`
	J    float64 `db:"j" json:"j"`
	Diff float64 `db:"diff" json:"diff"`
	Dea  float64 `db:"dea" json:"dea"`
}
