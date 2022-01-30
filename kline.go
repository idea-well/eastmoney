package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type kLineRes struct {
	Data struct {
		KLines kLineStrings `json:"klines"`
	} `json:"data"`
}

type kLineStrings []string

func (ks kLineStrings) toData() map[string]*KLineData {
	map_ := make(map[string]*KLineData)
	for _, kStr := range ks {
		ss := strings.Split(kStr, ",")
		date := strings.ReplaceAll(ss[0], "-", "")
		map_[date] = &KLineData{}
		map_[date].Open, _ = strconv.ParseFloat(ss[1], 64)
		map_[date].Close, _ = strconv.ParseFloat(ss[2], 64)
		map_[date].High, _ = strconv.ParseFloat(ss[3], 64)
		map_[date].Low, _ = strconv.ParseFloat(ss[4], 64)
		map_[date].Volume, _ = strconv.ParseFloat(ss[5], 64)
		map_[date].Amount, _ = strconv.ParseFloat(ss[6], 64)
		map_[date].Amplitude, _ = strconv.ParseFloat(ss[7], 64)
		map_[date].Change, _ = strconv.ParseFloat(ss[8], 64)
		map_[date].Turnover, _ = strconv.ParseFloat(ss[9], 64)
	}
	return map_
}

const kLineApi = "http://push2his.eastmoney.com/api/qt/stock/kline/get" +
	"?fields1=f5&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61" +
	"&klt=101&fqt=1&beg=%s&end=%s&lmt=1000000&secid=%d.%s"

func KLine(code string, market int, date ...string) (map[string]*KLineData, error) {
	var beg, end string
	var res = new(kLineRes)
	if len(date) == 0 {
		beg, end = "0", "20500101"
	} else {
		beg = minString(date[0], date...)
		end = maxString(date[0], date...)
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
	return res.Data.KLines.toData(), err
}

func KLineDate(code string, market int, date string) (*KLineData, error) {
	map_, err := KLine(code, market, date)
	return map_[date], err
}

func minString(s string, ss ...string) string {
	if len(ss) == 0 {
		return s
	}
	if s <= ss[0] {
		return minString(s, ss[1:]...)
	}
	return minString(ss[0], ss[1:]...)
}

func maxString(s string, ss ...string) string {
	if len(ss) == 0 {
		return s
	}
	if s >= ss[0] {
		return maxString(s, ss[1:]...)
	}
	return maxString(ss[0], ss[1:]...)
}
