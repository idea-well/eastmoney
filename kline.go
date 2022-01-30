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
		map_[ss[0]] = &KLineData{}
		map_[ss[0]].Open, _ = strconv.ParseFloat(ss[1], 64)
		map_[ss[0]].Close, _ = strconv.ParseFloat(ss[2], 64)
		map_[ss[0]].High, _ = strconv.ParseFloat(ss[3], 64)
		map_[ss[0]].Low, _ = strconv.ParseFloat(ss[4], 64)
		map_[ss[0]].Volume, _ = strconv.ParseFloat(ss[5], 64)
		map_[ss[0]].Amount, _ = strconv.ParseFloat(ss[6], 64)
		map_[ss[0]].Amplitude, _ = strconv.ParseFloat(ss[7], 64)
		map_[ss[0]].Change, _ = strconv.ParseFloat(ss[8], 64)
		map_[ss[0]].Turnover, _ = strconv.ParseFloat(ss[9], 64)
	}
	return map_
}

const kLineApi = "http://push2his.eastmoney.com/api/qt/stock/kline/get" +
	"?fields1=f5&fields2=f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61" +
	"&klt=101&fqt=1&end=20500101&lmt=1000000&secid=%d.%s"

func KLine(code string, market int) (map[string]*KLineData, error) {
	var res = new(kLineRes)
	resp, err := http.Get(fmt.Sprintf(kLineApi, market, code))
	err = callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
	return res.Data.KLines.toData(), err
}
