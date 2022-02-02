package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type kLineRes struct {
	Data struct {
		KLines    kLineStrings `json:"klines"`
		PreKPrice float64      `json:"preKPrice"`
	} `json:"data"`
}

type kLineStrings []string

func (ks kLineStrings) toData(lastPre float64) map[string]*KLineData {
	map_ := make(map[string]*KLineData)
	for _, kStr := range ks {
		ss := strings.Split(kStr, ",")
		date := strings.ReplaceAll(ss[0], "-", "")
		map_[date] = &KLineData{PreClose: lastPre}
		map_[date].Open = ParseFloat(ss[1])
		map_[date].Close = ParseFloat(ss[2])
		map_[date].High = ParseFloat(ss[3])
		map_[date].Low = ParseFloat(ss[4])
		map_[date].Volume = ParseInt(ss[5])
		map_[date].Amount = ParseFloat(ss[6])
		map_[date].Amplitude = ParseFloat(ss[7])
		map_[date].Change = ParseFloat(ss[8])
		map_[date].Turnover = ParseFloat(ss[9])
		map_[date].AvgPrice = map_[date].Amount / (float64(map_[date].Volume) * 100)
		lastPre = map_[date].Close
	}
	return map_
}

// KLineData
type KLineData struct {
	Open      float64 `json:"open"`      // 开盘
	Close     float64 `json:"close"`     // 收盘
	High      float64 `json:"high"`      // 最高
	Low       float64 `json:"low"`       // 最低
	Volume    int     `json:"volume"`    // 成交量
	Amount    float64 `json:"amount"`    // 成交额
	Change    float64 `json:"change"`    // 日涨幅
	Turnover  float64 `json:"turnover"`  // 换手率
	Amplitude float64 `json:"amplitude"` // 日振幅
	PreClose  float64 `json:"pre_close"` // 昨收盘
	AvgPrice  float64 `json:"avg_price"` // 平均价
}

const kLineApi = "http://push2his.eastmoney.com/api/qt/stock/kline/get" +
	"?fields1=f6&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f61" +
	"&klt=101&fqt=1&beg=%s&end=%s&lmt=1000000&secid=%d.%s"

func KLine(code string, market int, date ...string) (map[string]*KLineData, error) {
	var beg, end string
	var res = new(kLineRes)
	if len(date) == 0 {
		beg, end = "0", "20500101"
	} else {
		beg = MinString(date[0], date...)
		end = MaxString(date[0], date...)
	}
	url := fmt.Sprintf(kLineApi, beg, end, market, code)
	resp, err := http.Get(url)
	err = callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
	return res.Data.KLines.toData(res.Data.PreKPrice), err
}

func KLineDate(code string, market int, date string) (*KLineData, error) {
	map_, err := KLine(code, market, date)
	return map_[date], err
}
