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

func (ks kLineStrings) toData(lastPre float64, min bool) map[string]*KLineData {
	map_ := make(map[string]*KLineData)
	for _, kStr := range ks {
		var key string
		ss := strings.Split(kStr, ",")
		if min {
			key = strings.ReplaceAll(strings.Split(ss[0], " ")[1], ":", "")
		} else {
			key = strings.ReplaceAll(ss[0], "-", "")
		}
		map_[key] = &KLineData{PreClose: lastPre}
		map_[key].Open = ParseFloat(ss[1])
		map_[key].Close = ParseFloat(ss[2])
		map_[key].High = ParseFloat(ss[3])
		map_[key].Low = ParseFloat(ss[4])
		map_[key].Volume = ParseInt(ss[5])
		map_[key].Amount = ParseFloat(ss[6])
		map_[key].Amplitude = ParseFloat(ss[7])
		map_[key].Change = ParseFloat(ss[8])
		map_[key].Turnover = ParseFloat(ss[9])
		lastPre = map_[key].Close
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
}

const kLineApi = "http://push2his.eastmoney.com/api/qt/stock/kline/get" +
	"?fields1=f6&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f61" +
	"&klt=%d&fqt=1&beg=%s&end=%s&lmt=1000000&secid=%d.%s"

func MLine(code string, market int) (map[string]*KLineData, error) {
	var res = new(kLineRes)
	url := fmt.Sprintf(kLineApi, 1, "0", "20500101", market, code)
	resp, err := http.Get(url)
	err = callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
	return res.Data.KLines.toData(res.Data.PreKPrice, true), err
}

func KLine(code string, market int, date ...string) (map[string]*KLineData, error) {
	var beg, end string
	var res = new(kLineRes)
	if len(date) == 0 {
		beg, end = "0", "20500101"
	} else {
		beg = MinString(date[0], date...)
		end = MaxString(date[0], date...)
	}
	url := fmt.Sprintf(kLineApi, 101, beg, end, market, code)
	resp, err := http.Get(url)
	err = callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
	return res.Data.KLines.toData(res.Data.PreKPrice, false), err
}

func KLineDate(code string, market int, date string) (*KLineData, error) {
	map_, err := KLine(code, market, date)
	return map_[date], err
}
