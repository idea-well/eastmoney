package eastmoney

type kLineRes struct {
	Data struct {
		KLines kLineStrings `json:"klines"`
	} `json:"data"`
}

type kLineStrings []string

func (kss kLineStrings) toData() []*KLineData {
	var ds = make([]*KLineData, len(kss))
	return ds
}

type KLineData struct {
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume int     `json:"volume"`
	Amount float64 `json:"amount"`
}
