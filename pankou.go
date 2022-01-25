package eastmoney

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type panKouRes struct {
	Data struct {
		Pkyd panKouStrings `json:"pkyd"`
	} `json:"data"`
}

type panKouStrings []string

func (pks panKouStrings) toData(start, end string) []*PanKouData {
	ds := make([]*PanKouData, 0)
	for _, pk := range pks {
		ss := strings.Split(pk, ",")
		if ss[0] < start || ss[0] >= end {
			continue
		}
		ds = append(ds, &PanKouData{
			Time: ss[0], Code: ss[1], Market: ss[2], Name: ss[3],
			Type: ss[4], Desc: ss[5], Direct: ss[6],
		})
	}
	return ds
}

type PanKouData struct {
	Time   string `json:"time"`   // 时间
	Code   string `json:"code"`   // 代码
	Name   string `json:"name"`   // 名称
	Type   string `json:"type"`   // 类型代码
	Desc   string `json:"desc"`   // 类型描述
	Market string `json:"market"` // 市场编号
	Direct string `json:"direct"` // 方向 1涨 2跌
}

func PanKou(d time.Duration, limit int, handler func([]*PanKouData) error, logger func(error)) {
	var start, end string
	tck := time.NewTicker(d)
	for now := range tck.C {
		ss, err := doFetchPanKou(limit)
		logger(callWithoutErr(err, func() error {
			end = now.Format("15:04:05")
			return handler(ss.toData(start, end))
		}))
		start = end // swap start
	}
}

const panKouApi = "https://push2.eastmoney.com/api/qt/pkyd/get?fields=f1,f2,f3,f4,f5,f6,f7"

func doFetchPanKou(limit int) (panKouStrings, error) {
	var res = new(panKouRes)
	resp, err := http.Get(panKouApi + fmt.Sprintf("&lmt=%d", limit))
	return res.Data.Pkyd, callWithoutErr(err, func() error {
		defer resp.Body.Close()
		bts, err := ioutil.ReadAll(resp.Body)
		return callWithoutErr(err, func() error {
			return json.Unmarshal(bts, res)
		})
	})
}
