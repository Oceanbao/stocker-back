package screener

import "encoding/json"

type Screen struct {
	Ticker string  `json:"ticker"`
	Kdj    float64 `json:"kdj"`
}

func (s *Screen) ToMap() (map[string]interface{}, error) {
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
